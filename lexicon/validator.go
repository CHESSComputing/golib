package lexicon

// Lexicon validator module
// Copyright (c) 2024 - Valentin Kuznetsov <vkuznet@gmail.com>
//
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/CHESSComputing/golib/utils"
)

// Verbose defines module verbosity level
var Verbose int

// string parameters
var strParameters = []string{
	"did",
	"BTR",
	"cycle",
	"sample",
	"dataset",
	"parent",
	"release",
	"application",
	"tag",
	"processing",
	"file",
	"tier",
	"create_by",
	"user",
	"modify_by",
}

// integer parameters
var intParameters = []string{
	"cdate",
	"ldate",
	"min_cdate",
	"max_cdate",
	"min_ldate",
	"max_ldate",
	"dataset_id",
}

// mix type parameters
var mixParameters = []string{"run_num"}

// Lexicon represents single lexicon pattern structure
type Lexicon struct {
	Name     string   `json:"name"`
	Patterns []string `json:"patterns"`
	Length   int      `json:"length"`
}

func (r *Lexicon) String() string {
	data, err := json.MarshalIndent(r, "", "  ")
	if err == nil {
		return string(data)
	}
	return fmt.Sprintf("Lexicon: name=%s patters=%v length=%d", r.Name, r.Patterns, r.Length)
}

// LexiconPattern represents single lexicon compiled pattern structure
type LexiconPattern struct {
	Lexicon  Lexicon
	Patterns []*regexp.Regexp
}

// LexiconPatterns represents Lexicon patterns
var LexiconPatterns map[string]LexiconPattern

// LoadPatterns loads Lexion patterns from given file
// the format of the file is a list of the following dicts:
// [ {"name": <name>, "patterns": [list of patterns], "length": int},...]
func LoadPatterns(fname string) (map[string]LexiconPattern, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Printf("Unable to read, file '%s', error: %v\n", fname, err)
		return nil, Error(err, ReaderErrorCode, "", "lexicon.validator.LoadPatterns")
	}
	var records []Lexicon
	err = json.Unmarshal(data, &records)
	if err != nil {
		log.Printf("Unable to parse, file '%s', error: %v\n", fname, err)
		return nil, Error(err, UnmarshalErrorCode, "", "lexicon.validator.LoadPatterns")
	}
	// fetch and compile all patterns
	pmap := make(map[string]LexiconPattern)
	for _, rec := range records {
		var patterns []*regexp.Regexp
		for _, pat := range rec.Patterns {
			patterns = append(patterns, regexp.MustCompile(pat))
		}
		lex := LexiconPattern{Lexicon: rec, Patterns: patterns}
		key := rec.Name
		pmap[key] = lex
		if Verbose > 1 {
			log.Printf("regexp pattern\n%s", rec.String())
		}
	}
	return pmap, nil
}

// aux patterns
var UnixTimePattern = regexp.MustCompile(`^[1-9][0-9]{9}$`)
var IntPattern = regexp.MustCompile(`^\d+$`)
var RunRangePattern = regexp.MustCompile(`^\d+-\d+$`)

// ObjectPattern represents interface to check different objects
type ObjectPattern interface {
	Check(k string, v interface{}) error
}

// StrPattern represents string object pattern
type StrPattern struct {
	Patterns []*regexp.Regexp
	Len      int
}

// Check implements ObjectPattern interface for StrPattern objects
func (o StrPattern) Check(key string, val interface{}) error {
	if Verbose > 0 {
		log.Printf("StrPatern check key=%s val=%v", key, val)
		log.Printf("patterns %v max length %v", o.Patterns, o.Len)
	}
	var v string
	switch vvv := val.(type) {
	case string:
		v = vvv
	default:
		msg := fmt.Sprintf(
			"invalid type of input parameter '%s' for value '%+v' type '%T'",
			key, val, val)
		return Error(PatternErr, PatternErrorCode, msg, "lexicon.validator.Check")
	}
	if len(o.Patterns) == 0 {
		// nothing to match in patterns
		if Verbose > 0 {
			log.Println("nothing to match since we do not have patterns")
		}
		return nil
	}
	if o.Len > 0 && len(v) > o.Len {
		if Verbose > 0 {
			log.Println("lexicon str pattern", o)
		}
		// check for list of LFNs
		if key == "file" {
			for _, lfn := range lfnList(v) {
				if len(lfn) > o.Len {
					msg := fmt.Sprintf("length of LFN %s exceed %d characters", lfn, o.Len)
					return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.validator.Check")
				}
			}
		} else {
			msg := fmt.Sprintf("length of %s exceed %d characters", v, o.Len)
			return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.validator.Check")
		}
	}
	if key == "file" {
		for _, vvv := range lfnList(v) {
			msg := fmt.Sprintf("unable to match '%s' value '%s' from LFN list", key, vvv)
			var pass bool
			for _, pat := range o.Patterns {
				if matched := pat.MatchString(vvv); matched {
					// if at least one pattern matched we'll return
					pass = true
					break
				}
			}
			if !pass {
				return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.validator.Check")
			}
		}
		return nil
	}
	msg := fmt.Sprintf("unable to match '%s' value '%s'", key, val)
	for _, pat := range o.Patterns {
		if matched := pat.MatchString(v); matched {
			// if at least one pattern matched we'll return
			return nil
		}
	}
	return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.validator.Check")
}

// helper function to convert input value into list of list
// we need it to properly match LFN list
func lfnList(v string) []string {
	fileList := strings.Replace(v, "[", "", -1)
	fileList = strings.Replace(fileList, "]", "", -1)
	fileList = strings.Replace(fileList, "'", "", -1)
	fileList = strings.Replace(fileList, "\"", "", -1)
	var lfns []string
	for _, val := range strings.Split(fileList, ",") {
		lfns = append(lfns, strings.Trim(val, " "))
	}
	return lfns
}

// helper function to validate string parameters
//
//gocyclo:ignore
func strType(key string, val interface{}) error {
	var v string
	switch vvv := val.(type) {
	case string:
		v = vvv
	default:
		msg := fmt.Sprintf(
			"invalid type of input parameter '%s' for value '%+v' type '%T'",
			key, val, val)
		return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.validator.strType")
	}
	mapKeys := make(map[string]string)
	mapKeys["did"] = "did"
	mapKeys["dataset"] = "dataset"
	mapKeys["file"] = "file"
	mapKeys["create_by"] = "user"
	mapKeys["modify_by"] = "user"
	mapKeys["processing"] = "processing"
	mapKeys["application"] = "application"
	mapKeys["tier"] = "tier"
	mapKeys["dataset"] = "dataset"
	mapKeys["release"] = "cmssw_version"
	var allowedWildCardKeys = []string{
		"processing",
		"application",
		"tier",
		"release",
	}

	var patterns []*regexp.Regexp
	var length int

	for k, lkey := range mapKeys {
		if key == k {
			if utils.InList(k, allowedWildCardKeys) {
				if v == "" && val == "*" { // when someone passed wildcard
					return nil
				}
			}
			if p, ok := LexiconPatterns[lkey]; ok {
				patterns = p.Patterns
				length = p.Lexicon.Length
			}
		}
		if key == "file" {
			if strings.Contains(v, "[") {
				if strings.Contains(v, "'") { // Python bad json, e.g. ['bla']
					v = strings.Replace(v, "'", "\"", -1)
				}
				var records []string
				err := json.Unmarshal([]byte(v), &records)
				if err != nil {
					return Error(err, UnmarshalErrorCode, "", "lexicon.validator.strType")
				}
				for _, r := range records {
					err := StrPattern{Patterns: patterns, Len: length}.Check(key, r)
					if err != nil {
						return Error(err, PatternErrorCode, "", "lexicon.validator.strType")
					}
				}
			}
		}
		if key == "block_name" {
			if strings.Contains(v, "[") {
				if strings.Contains(v, "'") { // Python bad json, e.g. ['bla']
					v = strings.Replace(v, "'", "\"", -1)
				}
				// split input into pieces
				input := strings.Replace(v, "[", "", -1)
				input = strings.Replace(input, "]", "", -1)
				for _, vvv := range strings.Split(input, ",") {
					err := checkBlockHash(strings.Trim(vvv, " "))
					if err != nil {
						return err
					}
				}
			} else {
				err := checkBlockHash(v)
				if err != nil {
					return err
				}
			}
		}
	}
	return StrPattern{Patterns: patterns, Len: length}.Check(key, val)
}

// helper function to check block hash
func checkBlockHash(blk string) error {
	arr := strings.Split(blk, "#")
	if len(arr) != 2 {
		msg := fmt.Sprintf("wrong parts in block name %s", blk)
		return Error(ValidationErr, PatternErrorCode, msg, "lexicon.validator.checkBlockHash")
	}
	if len(arr[1]) > 36 {
		msg := fmt.Sprintf("wrong length of block hash %s", blk)
		return Error(ValidationErr, PatternErrorCode, msg, "lexicon.validator.checkBlockHash")
	}
	return nil
}

// helper function to validate int parameters
func intType(k string, v interface{}) error {
	// to be implemented
	return nil
}

// helper function to validate mix parameters
func mixType(k string, v interface{}) error {
	// to be implemented
	return nil
}

// ValidateRecord validates given JSON record
func ValidateRecord(rec map[string]any) error {
	for key, val := range rec {
		if utils.InList(key, strParameters) {
			if err := strType(key, val); err != nil {
				return Error(err, ValidateErrorCode, "not str type", "lexicon.Validate")
			}
		}
	}
	return nil
}

// Validate provides validation of all input parameters of HTTP request
func Validate(r *http.Request) error {
	if r.Method == "GET" {
		for k, vvv := range r.URL.Query() {
			// vvv here is []string{} type since all HTTP parameters are treated
			// as list of strings
			for _, v := range vvv {
				if utils.InList(k, strParameters) {
					if err := strType(k, v); err != nil {
						return Error(err, ValidateErrorCode, "not str type", "lexicon.Validate")
					}
				}
				if utils.InList(k, intParameters) {
					if err := intType(k, v); err != nil {
						return Error(err, ValidateErrorCode, "not int type", "lexicon.Validate")
					}
				}
				if utils.InList(k, mixParameters) {
					if err := mixType(k, v); err != nil {
						return Error(err, ValidateErrorCode, "not mix type", "lexicon.Validate")
					}
				}
			}
			if Verbose > 0 {
				log.Printf("query parameter key=%s values=%+v\n", k, vvv)
			}
		}
	}
	return nil
}

// CheckPattern is a generic functino to check given key value within Lexicon map
func CheckPattern(key, value string) error {
	if p, ok := LexiconPatterns[key]; ok {
		for _, pat := range p.Patterns {
			if matched := pat.MatchString(value); matched {
				if Verbose > 1 {
					log.Printf("CheckPattern key=%s value='%s' found match %s", key, value, pat)
				}
				return nil
			}
			if Verbose > 1 {
				log.Printf("CheckPattern key=%s value='%s' does not match %s", key, value, pat)
			}
		}
		msg := fmt.Sprintf("invalid pattern for key=%s", key)
		return Error(InvalidParamErr, PatternErrorCode, msg, "lexicon.CheckPattern")
	}
	return nil
}

// ValidatePostPayload function to validate POST request
func ValidatePostPayload(rec map[string]any) error {
	for key, val := range rec {
		errMsg := fmt.Sprintf("unable to match '%s' value '%+v'", key, val)
		if key == "data_tier_name" {
			if vvv, ok := val.(string); ok {
				if err := CheckPattern("data_tier_name", vvv); err != nil {
					return Error(err, PatternErrorCode, "wrong data_tier_name pattern", "lexicon.ValidaatePostPayload")
				}
			}
		} else if key == "create_at" || key == "modify_at" {
			v, err := utils.CastInt(val)
			if err != nil {
				return Error(err, PatternErrorCode, errMsg, "lexicon.ValidaatePostPayload")
			} else if matched := UnixTimePattern.MatchString(fmt.Sprintf("%d", v)); !matched {
				return Error(InvalidParamErr, PatternErrorCode, errMsg, "lexicon.ValidaatePostPayload")
			}
		}
	}
	return nil
}
