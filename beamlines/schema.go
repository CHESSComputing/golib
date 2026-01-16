package beamlines

// schema module
//
// Copyright (c) 2022 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
	yaml "gopkg.in/yaml.v2"
)

// Verbose control verbosity printout level
var Verbose int

// SkipKeys defines list of schema keys to skip validation process
var SkipKeys = []string{}

// SchemaKeys represents full collection of schema keys across all schemas
type SchemaKeys map[string]string

// schema keys map
var _schemaKeys SchemaKeys

// SchemaRenewInterval setup internal to update schema cache
var SchemaRenewInterval time.Duration

// SchemaObject holds current MetaData schema
type SchemaObject struct {
	Schema   *Schema
	LoadTime time.Time
}

// SchemaManager holds current map of MetaData schema objects
type SchemaManager struct {
	Map     map[string]*SchemaObject
	Verbose int
}

// Schema returns either cached schema map or load it from provided file
func (m *SchemaManager) String() string {
	var out string
	for k, v := range m.Map {
		out += fmt.Sprintf("\n%s %s, loaded %v\n", k, v.Schema, v.LoadTime)
	}
	return out
}

// Schema returns either cached schema map or load it from provided file
func (m *SchemaManager) Load(fname string) (*Schema, error) {
	// use full path of file name
	fname = utils.FullPath(fname)

	// check fname in our schema map
	if sobj, ok := m.Map[fname]; ok {
		if sobj.Schema != nil && time.Since(sobj.LoadTime) < SchemaRenewInterval {
			log.Println("schema taken from cache", fname)
			return sobj.Schema, nil
		}
	}
	schema := &Schema{FileName: fname, Verbose: m.Verbose}
	err := schema.Load()
	if err != nil {
		log.Println("unable to load schema from", fname, " error", err)
		return schema, err
	}
	if Verbose > 1 {
		log.Println("renew schema:", fname)
	}
	// reset map if it is expired
	if sobj, ok := m.Map[fname]; ok {
		if sobj.Schema != nil && time.Since(sobj.LoadTime) > SchemaRenewInterval {
			if Verbose > 1 {
				log.Println("reset schema manager")
			}
			m.Map = nil
		}
	}
	if m.Map == nil {
		m.Map = make(map[string]*SchemaObject)
	}
	m.Map[fname] = &SchemaObject{Schema: schema, LoadTime: time.Now()}

	return schema, nil
}

// MetaDetails returns list of schema unit maps, each map contains schema attribute and its units key-value pairs
func (m *SchemaManager) MetaDetails() []SchemaDetails {
	var out []SchemaDetails
	for _, sobj := range m.Map {
		smap := make(map[string]string)
		dmap := make(map[string]string)
		tmap := make(map[string]string)
		for _, rec := range sobj.Schema.Map {
			smap[rec.Key] = rec.Units
			dmap[rec.Key] = rec.Description
			tmap[rec.Key] = rec.Type
		}
		sname := SchemaName(sobj.Schema.FileName)
		sunits := SchemaDetails{Schema: sname, Units: smap, Descriptions: dmap, DataTypes: tmap}
		out = append(out, sunits)
	}
	return out
}

// SchemaDetails represents individual FOXDEN schema units dictionary
type SchemaDetails struct {
	Schema       string            `json:"schema"`
	Units        map[string]string `json:"units"`
	Descriptions map[string]string `json:"descriptions"`
	DataTypes    map[string]string `json:"types"`
}

// SchemaRecord provide schema record structure
type SchemaRecord struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Optional    bool   `json:"optional"`
	Multiple    bool   `json:"multiple"`
	Section     string `json:"section"`
	Value       any    `json:"value"`
	Schema      string `json:"schema"`
	Placeholder string `json:"placeholder"`
	Units       string `json:"units"`
	Description string `json:"description"`
	File        string `json:"file,omitempty"` // Used for inclusion
}

// SchemaCacheManager provides cache for schema maps
type SchemaCacheManager struct {
	mu    sync.RWMutex
	Cache map[string]*Schema
}

func NewSchemaCacheManager() *SchemaCacheManager {
	return &SchemaCacheManager{Cache: make(map[string]*Schema)}
}

func (s *SchemaCacheManager) Get(k string) (*Schema, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.Cache[k]
	return v, ok
}

func (s *SchemaCacheManager) Set(k string, v *Schema) {
	s.mu.Lock()
	s.Cache[k] = v
	s.mu.Unlock()
}

var _smgr *SchemaCacheManager

// Schema provides structure of schema file
type Schema struct {
	FileName       string                  `json:"fileName`
	Map            map[string]SchemaRecord `json:"map"`
	WebSectionKeys map[string][]string     `json:"webSectionKeys"`
	Verbose        int                     `json:"verbose"`
}

// Load loads given schema file
func (s *Schema) String() string {
	if s.Map != nil {
		return fmt.Sprintf("<schema %s, map %d entries>", s.FileName, len(s.Map))
	}
	return fmt.Sprintf("<schema %s, map %v>", s.FileName, s.Map)
}

// Load loads given schema file
func (s *Schema) Load() error {
	fname := s.FileName
	if _smgr == nil {
		_smgr = NewSchemaCacheManager()
	}
	if sv, ok := _smgr.Get(fname); ok {
		// take schema from the cache
		s.FileName = sv.FileName
		s.Map = sv.Map
		s.WebSectionKeys = sv.WebSectionKeys
		s.Verbose = sv.Verbose
		if sv.Verbose > 1 {
			log.Printf("use cached schema %+v", s)
		}
		return nil
	}
	if s.Verbose > 1 {
		log.Printf("loading new schema %+v from file=%s", s, fname)
	}
	file, err := os.Open(fname)
	if err != nil {
		msg := fmt.Sprintf("Unable to open %s, error=%v", fname, err)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		msg := fmt.Sprintf("Unable to read %s, error=%v", fname, err)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	var records []SchemaRecord
	if strings.HasSuffix(fname, "json") {
		err = json.Unmarshal(data, &records)
		if err != nil {
			msg := fmt.Sprintf("fail to unmarshal json file %s, error=%v", fname, err)
			log.Printf("ERROR: %s", msg)
			return errors.New(msg)
		}
	} else if strings.HasSuffix(fname, "yaml") || strings.HasSuffix(fname, "yml") {
		var yrecords []map[interface{}]interface{}
		err = yaml.Unmarshal(data, &yrecords)
		if err != nil {
			msg := fmt.Sprintf("fail to unmarshal yaml file %s, error=%v", fname, err)
			log.Printf("ERROR: %s", msg)
			return errors.New(msg)
		}
		for _, yr := range yrecords {
			m := convertYaml(yr)
			smap := SchemaRecord{}
			for k, v := range m {
				if k == "key" {
					smap.Key = v.(string)
				} else if k == "type" {
					smap.Type = v.(string)
				} else if k == "optional" {
					smap.Optional = v.(bool)
				} else if k == "description" {
					smap.Description = v.(string)
				} else if k == "placeholder" {
					smap.Placeholder = v.(string)
				}
			}
			records = append(records, smap)
		}
	} else {
		msg := fmt.Sprintf("unsupported data format of schema file %s", fname)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	s.FileName = fname
	smap := make(map[string]SchemaRecord)
	for _, r := range records {
		if r.File != "" {
			// check if provided nested record file name is relative or absolute
			nestedFileName := r.File
			if _, err := os.Stat(nestedFileName); err == nil {
				// do nothing as file exsti
			} else if os.IsNotExist(err) {
				// use file directory of schema file fname as directory for embeded schema file
				fdir := filepath.Dir(fname)
				nestedFileName = fmt.Sprintf("%s/%s", fdir, r.File)
			}
			if nestedRecords, err := loadNestedRecords(nestedFileName); err == nil {
				for _, nr := range nestedRecords {
					if nr.Key != "" {
						smap[nr.Key] = nr
					}
				}
			} else {
				log.Printf("ERROR: unable to load nested schema from file %s, error=%v", r.File, err)
			}
		}
		smap[r.Key] = r
	}
	// discard map record with embeded schema file
	if r, ok := smap[""]; ok {
		if r.File != "" {
			delete(smap, "")
		}
	}
	// update schema map
	s.Map = smap

	// upload SchemaKeys object
	if _schemaKeys == nil {
		_schemaKeys = make(SchemaKeys)
	}
	for _, r := range smap {
		if _, ok := _schemaKeys[r.Key]; !ok {
			_schemaKeys[strings.ToLower(r.Key)] = r.Key
		}
	}

	filepath := srvConfig.Config.CHESSMetaData.WebSectionsFile
	if filepath == "" {
		// write schema to cache
		_smgr.Set(fname, s)
		return nil
	}
	if _, err := os.Stat(filepath); err == nil {
		file, err := os.Open(filepath)
		if err != nil {
			log.Println("unable to open", filepath, "error", err)
			return err
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			log.Println("unable to read file, error", err)
			return err
		}
		var rec map[string][]string
		err = json.Unmarshal(data, &rec)
		if err != nil {
			log.Fatal(err)
		}
		s.WebSectionKeys = rec
	}

	// write schema to cache
	_smgr.Set(fname, s)
	return nil
}

// Validate validates given record against schema
func (s *Schema) Validate(rec map[string]any) error {
	if err := s.Load(); err != nil {
		return err
	}
	keys, err := s.Keys()
	if err != nil {
		return err
	}
	// hidden mandatory keys we add to each form
	var mkeys []string
	for k, v := range rec {
		// skip user key if it does not belong to schema
		if utils.InList(k, SkipKeys) && !utils.InList(k, keys) {
			continue
		}
		// check if our record key belong to the schema keys
		if !utils.InList(k, keys) {
			msg := fmt.Sprintf("record key '%s' is not known", k)
			log.Printf("ERROR: %s, schema file %s, schema map %+v", msg, s.FileName, s.Map)
			return errors.New(msg)
		}

		if m, ok := s.Map[k]; ok {
			// check key name
			if m.Key != k {
				msg := fmt.Sprintf("invalid key=%s", k)
				log.Printf("ERROR: %s", msg)
				return errors.New(msg)
			}
			if m.Schema != "" {
				// we got record with another schema: it should be eitehr map or list of map records
				if e := validateStructs(s.FileName, m, v, s.Verbose); e != nil {
					return e
				}
				// collect mandatory keys
				if !m.Optional {
					mkeys = append(mkeys, k)
				}
			} else {
				// check data type
				if !validateSchemaType(m.Type, v, s.Verbose) {
					// check if provided data type can be converted to m.Type
					msg := fmt.Sprintf("invalid data type for key=%s, value=%v, type=%T, expect=%s", k, v, v, m.Type)
					log.Printf("ERROR: %s", msg)
					return errors.New(msg)
				}
				// check data value
				if !validateRecordValue(m, v, s.Verbose) {
					// check if provided data type can be converted to m.Type
					msg := fmt.Sprintf("invalid data value for key=%s, type=%s, multiple=%v, value=%v valuetype=%T", k, m.Type, m.Multiple, v, v)
					log.Printf("ERROR: %s", msg)
					return errors.New(msg)
				}
				// collect mandatory keys
				if !m.Optional {
					mkeys = append(mkeys, k)
				}
			}
		}
	}

	// check that we collected all mandatory keys
	smkeys, err := s.MandatoryKeys()
	if err != nil {
		return err
	}
	if len(mkeys) != len(smkeys) {
		sort.Sort(utils.StringList(mkeys))
		var missing []string
		for _, k := range smkeys {
			if !utils.InList(k, mkeys) {
				missing = append(missing, k)
			}
		}
		msg := fmt.Sprintf("Schema %s, mandatory keys %v, record keys %v, missing keys %v", s.FileName, smkeys, mkeys, missing)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	return nil
}

// Keys provides list of keys of the schema
func (s *Schema) Keys() ([]string, error) {
	var keys []string
	if err := s.Load(); err != nil {
		return keys, err
	}
	for k, _ := range s.Map {
		if k != "" {
			keys = append(keys, k)
		}
	}
	sort.Sort(utils.StringList(keys))
	return keys, nil
}

// OptionalKeys provides list of optional keys of the schema
func (s *Schema) OptionalKeys() ([]string, error) {
	var keys []string
	if err := s.Load(); err != nil {
		return keys, err
	}
	for k, _ := range s.Map {
		if m, ok := s.Map[k]; ok {
			if m.Optional {
				keys = append(keys, k)
			}
		}
	}
	sort.Sort(utils.StringList(keys))
	return keys, nil
}

// MandatoryKeys provides list of madatory keys of the schema
func (s *Schema) MandatoryKeys() ([]string, error) {
	var keys []string
	if err := s.Load(); err != nil {
		return keys, err
	}
	for k, _ := range s.Map {
		if m, ok := s.Map[k]; ok {
			if !m.Optional {
				keys = append(keys, k)
			}
		}
	}
	sort.Sort(utils.StringList(keys))
	return keys, nil
}

// Sections provides list of schema sections
func (s *Schema) Sections() ([]string, error) {
	if len(srvConfig.Config.CHESSMetaData.OrderedSections) > 0 {
		return srvConfig.Config.CHESSMetaData.OrderedSections, nil
	}
	var sections []string
	if err := s.Load(); err != nil {
		return sections, err
	}
	for k, _ := range s.Map {
		if m, ok := s.Map[k]; ok {
			if m.Section != "" {
				if !utils.InList(m.Section, sections) {
					sections = append(sections, m.Section)
				}
			}
		}
	}
	sort.Sort(utils.StringList(sections))
	return sections, nil
}

// SectionKeys provides map of section keys
func (s *Schema) SectionKeys() (map[string][]string, error) {
	smap := make(map[string][]string)
	sections, err := s.Sections()
	if err != nil {
		return smap, err
	}
	allKeys, err := s.Keys()
	if err != nil {
		return smap, err
	}
	// populate section map with keys defined in webSectionKeys
	if s.WebSectionKeys != nil {
		for k, v := range s.WebSectionKeys {
			smap[k] = v
		}
	}
	// loop over all sections and add section keys to the map
	for _, sect := range sections {
		for _, k := range allKeys {
			if r, ok := s.Map[k]; ok {
				if r.Section == sect {
					if skeys, ok := smap[sect]; ok {
						if !utils.InList(k, skeys) {
							skeys = append(skeys, k)
							smap[sect] = skeys
						}
					} else {
						smap[sect] = []string{k}
					}
				}
			}
		}
	}
	return smap, nil
}

// helper function to validate any structs within our record schema
func validateStructs(path string, rec SchemaRecord, v any, verbose int) error {
	var structType, listStructType bool
	valid := false
	switch vt := v.(type) {
	case []map[string]any:
		listStructType = true
		for _, r := range vt {
			valid = validateStruct(path, rec, r, verbose)
			if !valid {
				break
			}
		}
	case map[string]any:
		structType = true
		valid = validateStruct(path, rec, v.(map[string]any), verbose)
	default:
		msg := fmt.Sprintf("unsupported sub-schema record type=%T", v)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	if !valid {
		msg := fmt.Sprintf("invalid sub-struct record=%v, subschema=%s, expect=%s", v, rec.Schema, rec.Type)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	if rec.Type == "struct" && !structType {
		msg := fmt.Sprintf("mismatch of record type and record value, expected type struct (generic map) but received %T", v)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	if rec.Type == "list_struct" && !listStructType {
		msg := fmt.Sprintf("mismatch of record type and record value, expected type list_struct (list of maps) but received %T", v)
		log.Printf("ERROR: %s", msg)
		return errors.New(msg)
	}
	return nil
}

// helper function to validate sub-structure within schema record
func validateStruct(path string, rec SchemaRecord, v map[string]any, verbose int) bool {
	if verbose > 0 {
		log.Printf("validate subschema record %+v value=%+v", rec, v)
	}
	// load schema from given path and rec.Schema
	dir := filepath.Dir(path)
	s := &Schema{FileName: filepath.Join(dir, rec.Schema)}
	err := s.Load()
	if err != nil {
		log.Printf("ERROR: unable to load sub-schema from record %+v filename=%s err=%v", rec, path, err)
		return false
	}
	// loop over new schema records and validate our map value against them
	for key, record := range s.Map {
		if verbose > 0 {
			log.Printf("### subschema key=%s schema=%+v value=%+v", key, record, v)
		}
		if val, ok := v[key]; ok {
			// check data type
			if verbose > 0 {
				log.Printf("+++ validateSchemaType type=%s, value=%+v", record.Type, val)
			}
			if !validateSchemaType(record.Type, val, verbose) {
				log.Printf("struct %+v has invalid schema type='%s' expect %T", record, record.Type, val)
				return false
			}
			// check data value
			if verbose > 0 {
				log.Printf("+++ validateRecordValue record=%+v, value=%+v", record, val)
			}
			if !validateRecordValue(record, val, verbose) {
				log.Printf("struct %+v has invalid record value %+v", record, val)
				return false
			}
		}
	}
	/*
		// check data type
		if !validateSchemaType(rec.Type, v, verbose) {
			log.Printf("struct %+v has invalid schema type='%s' expect %T", rec, rec.Type, v)
			return false
		}
		// check data value
		if !validateRecordValue(rec, v, verbose) {
			log.Printf("struct %+v has invalid record value %+v", rec, v)
			return false
		}
	*/
	return true
}

// helper function to validate given value with respect to schema one
// only valid for value of list type
func validateRecordValue(rec SchemaRecord, v any, verbose int) bool {
	if rec.Type == "any" {
		return true
	}

	vtype := simpleType(v)
	if rec.Type == "struct" {
		log.Printf("validateRecordValue: rec=%+v value=%+v reassign type=%s", rec, v, rec.Type)
		vtype = rec.Type
	}
	if rec.Type == "list_struct" {
		log.Printf("validateRecordValue: rec=%+v value=%+v reassign type=%s", rec, v, rec.Type)
		vtype = rec.Type
	}
	// check for non list data-types
	if !strings.HasPrefix(rec.Type, "list") {
		// special case of zero float value and int schema record data-type
		if strings.Contains(rec.Type, "int") && strings.Contains(vtype, "float") {
			switch vvv := v.(type) {
			case float32:
				if vvv == 0 {
					return true
				}
				if Float64IsInt64Compatible(float64(vvv)) {
					return true
				}
			case float64:
				if vvv == 0 {
					return true
				}
				if Float64IsInt64Compatible(vvv) {
					return true
				}
			}
		}
		// check if data-types are different for non-list data-types
		if !strings.Contains(rec.Type, vtype) {
			log.Printf("ERROR: record type %s differ from value data-type %s", rec.Type, vtype)
			return false
		}
		// check if data-value are the same
		if fmt.Sprintf("%v", rec.Value) == fmt.Sprintf("%v", v) {
			// if data types of schema record and passed value are the same we declare that it is valid data
			return true
		}

		// compare values for schema fields with types other than "list_*"
		if rec.Value != nil {
			sv := fmt.Sprintf("%v", v)
			matched := false
			if verbose > 0 {
				log.Printf("rec=%+v, type(rec.Value)=%T v=%v type(v)=%T", rec, rec.Value, v, v)
			}
			switch vvv := rec.Value.(type) {
			case []any:
				for _, val := range vvv {
					sval := fmt.Sprintf("%v", val)
					if sv == sval {
						matched = true
					}
				}
			case []int:
				for _, val := range vvv {
					if v == val {
						matched = true
					}
				}
			case []float64:
				for _, val := range vvv {
					if v == val {
						matched = true
					}
				}
			case any:
				sval := fmt.Sprintf("%v", vvv)
				if sv == sval {
					matched = true
				}
			}
			if !matched {
				log.Printf("ERROR: no match found for record value %+v", rec)
				return false
			}
		}
	}

	// checks for list data type
	if strings.HasPrefix(rec.Type, "list") {
		var values []string
		if rec.Value == nil {
			return true
		}
		for _, v := range rec.Value.([]any) {
			vvv := strings.Trim(fmt.Sprintf("%v", v), " ")
			if !utils.InList(vvv, values) {
				values = append(values, vvv)
			}
		}
		matched := false
		if verbose > 0 {
			log.Printf("checking %v of type %T against %+v", v, v, rec)
		}
		if verbose > 0 {
			log.Printf("checking v=%v of type %T vtype=%v", v, v, vtype)
		}
		vtype := fmt.Sprintf("%T", v)
		if strings.HasPrefix(vtype, "[]") {
			// our input value is a list data-type and we should check all its values
			var matchArr []bool
			var rvalues []string
			switch rvals := v.(type) {
			case []string:
				for _, rv := range rvals {
					for _, vvv := range strings.Split(rv, " ") {
						rvalues = append(rvalues, vvv)
					}
				}
			case []any:
				for _, rv := range rvals {
					vvv := fmt.Sprintf("%v", rv)
					rvalues = append(rvalues, vvv)
				}
			}
			for _, rv := range rvalues {
				for _, vvv := range values {
					//                     if rv == vvv || (vvv != "" && strings.Contains(rv, vvv)) {
					if rv == vvv {
						matchArr = append(matchArr, true)
					}
				}
			}
			if verbose > 0 {
				log.Printf("values list %v type=%T total=%d", values, values, len(values))
				log.Printf("matched list %v type=%T total=%d", matchArr, matchArr, len(matchArr))
				log.Printf("expected list %v type=%T total=%d", rvalues, rvalues, len(rvalues))
			}
			// all matched values
			if len(matchArr) == len(rvalues) {
				matched = true
			}
		} else {
			for _, val := range values {
				vvv := fmt.Sprintf("%v", val)
				if v == vvv {
					matched = true
				}
			}
		}
		if !matched {
			return false
		}
	}
	return true
}

// helper function to return simple data-type, int64->int, float32 -> float
func simpleType(v any) string {
	vtype := fmt.Sprintf("%T", v)
	for _, s := range []string{"8", "16", "32", "64"} {
		vtype = strings.Replace(vtype, s, "", -1)
	}
	return vtype
}

// helper function to check if float64 can be converted to int64 without precision lost
func Float64IsInt64Compatible(v float64) bool {
	if v > float64(math.MaxInt64) || v < float64(math.MinInt64) {
		return false
	}
	return v == math.Trunc(v)
}

// helper function to validate schema type of given value with respect to schema
func validateSchemaType(stype string, v interface{}, verbose int) bool {
	// on web form 0 will be int type, but we can allow it for any int's float's
	if v == 0 || v == 0. {
		if strings.Contains(stype, "int") || strings.Contains(stype, "float") {
			return true
		}
	}

	// check if stype is struct it should be either map or list of maps
	if stype == "struct" {
		switch v.(type) {
		case map[string]any:
			return true
		default:
			return false
		}
	}
	if stype == "list_struct" {
		switch v.(type) {
		case []map[string]any:
			return true
		default:
			return false
		}
	}

	// check actual value type and compare it to given schema type
	var etype string
	switch v.(type) {
	case bool:
		etype = "bool"
	case int:
		etype = "int"
	case int8:
		etype = "int8"
	case int16:
		etype = "int16"
	case int32:
		etype = "int32"
	case int64:
		etype = "int64"
	case uint16:
		etype = "uint16"
	case uint32:
		etype = "uint32"
	case uint64:
		etype = "uint64"
	case float32:
		etype = "float"
	case float64:
		etype = "float64"
	case string:
		etype = "string"
	case []string:
		etype = "list_str"
	case []any:
		etype = "list_str"
	case []int:
		etype = "list_int"
	case []float64:
		etype = "list_float"
	case []float32:
		etype = "list_float"
	default:
		etype = "any"
	}
	sv := fmt.Sprintf("%v", v)
	vtype := fmt.Sprintf("%T", v)
	if verbose > 1 {
		log.Printf("### validateSchemaType schema type=%v value type=%T value=%v", stype, v, sv)
	}
	if stype == "int64" && vtype == "float64" && !strings.Contains(sv, ".") {
		return true
	}
	if stype == "list_float" && vtype == "[]interface {}" {
		return true
	}
	// empty list of floats
	if stype == "list_float" && vtype == "[]string" && sv == "[]" {
		return true
	}
	// check if we can reduce data-type
	if sv == "" && etype != "string" {
		return false
	}
	if stype != etype {
		return false
	}
	return true
}

// helper function to convert yaml map to json map interface
func convertYaml(m map[interface{}]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		switch v2 := v.(type) {
		case map[interface{}]interface{}:
			res[fmt.Sprint(k)] = convertYaml(v2)
		default:
			res[fmt.Sprint(k)] = v
		}
	}
	return res
}

// helper function to load nested schema records
func loadNestedRecords(filename string) ([]SchemaRecord, error) {
	var records []SchemaRecord
	_, err := os.Stat(filename)
	if err != nil {
		return records, err
	}
	file, err := os.Open(filename)
	if err != nil {
		return records, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return records, err
	}
	if strings.HasSuffix(filename, "json") {
		err = json.Unmarshal(data, &records)
		if err != nil {
			return records, err
		}
	} else if strings.HasSuffix(filename, "yaml") || strings.HasSuffix(filename, "yml") {
		var yrecords []map[interface{}]interface{}
		err = yaml.Unmarshal(data, &yrecords)
		if err != nil {
			return records, err
		}
		for _, yr := range yrecords {
			m := convertYaml(yr)
			smap := SchemaRecord{}
			for k, v := range m {
				if k == "key" {
					smap.Key = v.(string)
				} else if k == "type" {
					smap.Type = v.(string)
				} else if k == "optional" {
					smap.Optional = v.(bool)
				} else if k == "description" {
					smap.Description = v.(string)
				} else if k == "placeholder" {
					smap.Placeholder = v.(string)
				}
			}
			records = append(records, smap)
		}
	}
	if Verbose > 1 {
		log.Printf("### loading %s", filename)
		for _, r := range records {
			log.Printf("### record %+v", r)
		}
	}
	return records, nil
}
