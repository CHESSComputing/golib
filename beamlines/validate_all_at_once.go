package beamlines

//
// Copyright (c) 2026 - Valentin Kuznetsov <vkuznet AT gmail dot com>
//

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	srvConfig "github.com/CHESSComputing/golib/config"
	utils "github.com/CHESSComputing/golib/utils"
)

// ValidateTmplRecord validates a partial (tempalte) record against the schema
// collecting ALL errors rather than returning on the first one. Returns an empty
// string on success, or a human-readable multi-line report of every problem found.
func (s *Schema) ValidateTmplRecord(rec map[string]any) string {
	return s.validateAll(rec, false)
}

// ValidateAll validates a record against the schema collecting ALL errors
// rather than returning on the first one. Returns an empty string on success,
// or a human-readable multi-line report of every problem found.
func (s *Schema) ValidateAll(rec map[string]any) string {
	return s.validateAll(rec, true)
}

// validateAll validates a record against the schema collecting ALL errors
// rather than returning on the first one. Returns an empty string on success,
// or a human-readable multi-line report of every problem found.
func (s *Schema) validateAll(rec map[string]any, checkMandatoryKeys bool) string {
	var errs []string
	add := func(format string, a ...any) {
		errs = append(errs, fmt.Sprintf(format, a...))
	}

	if err := s.Load(); err != nil {
		return fmt.Sprintf("schema load error: %v", err)
	}
	keys, err := s.Keys()
	if err != nil {
		return fmt.Sprintf("schema keys error: %v", err)
	}

	var mkeys []string // mandatory keys actually present in the record

	// ── Pass 1: validate every key/value that the record contains ─────────────
	for k, v := range rec {
		// skip known meta-keys that are not part of the schema
		if utils.InList(k, srvConfig.Config.CHESSMetaData.SkipKeys) && !utils.InList(k, keys) {
			continue
		}

		// check the key is known at all
		if !utils.InList(k, keys) {
			if !checkSubKeys(k, v, keys) {
				add("unknown key %q (value %v, type %T) — not present in schema %s",
					k, v, v, s.FileName)
			}
			// unknown key: nothing more to check for it
			continue
		}

		check1 := false // matched a top-level schema record
		if m, ok := s.Map[k]; ok {
			// key name sanity (should always match, but guard anyway)
			if m.Key != k {
				add("key mismatch: got %q, schema has %q", k, m.Key)
			}
			// type check
			if !validateSchemaType(m.Type, v, s.Verbose) {
				add("wrong type for key %q: got %T (%v), schema expects %s",
					k, v, v, m.Type)
			}
			// value check
			if !validateRecordValue(m, v, s.Verbose) {
				add("invalid value for key %q: value=%v (type %T), schema type=%s, multiple=%v",
					k, v, v, m.Type, m.Multiple)
			}
			if !m.Optional {
				mkeys = append(mkeys, k)
			}
			check1 = true
		}

		// check composed sub-struct keys (e.g. "sample.name")
		check2 := false
		for sk := range s.Map {
			if strings.Contains(sk, ".") && strings.HasPrefix(sk, k) {
				if m, ok := s.Map[sk]; ok {
					if e := validateStructs(m.File, m, v, s.Verbose); e != nil {
						add("sub-schema validation failed for key %q: %v", k, e)
					}
					if !m.Optional {
						mkeys = append(mkeys, sk)
					}
					check2 = true
				} else {
					add("composed key %q not found in schema map", sk)
				}
			}
		}

		if !check1 && !check2 {
			add("key %q (value %v) matched no schema record", k, v)
		}
	}

	// ── Pass 2: check that all mandatory keys are present ─────────────────────
	if checkMandatoryKeys {
		mkeys = utils.List2Set(mkeys)

		smkeys, err := s.MandatoryKeys()
		if err != nil {
			add("could not retrieve mandatory keys: %v", err)
		} else if len(mkeys) != len(smkeys) {
			sort.Sort(utils.StringList(mkeys))
			for _, k := range smkeys {
				if !utils.InList(k, mkeys) {
					add("missing mandatory key %q", k)
				}
			}
		}
	}

	// ── Build report ──────────────────────────────────────────────────────────
	if len(errs) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("validation failed against schema %s (%d error(s)):\n",
		filepath.Base(s.FileName), len(errs)))
	for i, e := range errs {
		sb.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, e))
	}
	return strings.TrimRight(sb.String(), "\n")
}
