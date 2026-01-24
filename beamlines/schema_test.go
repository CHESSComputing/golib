package beamlines

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	srvConfig "github.com/CHESSComputing/golib/config"
)

// TestSchemaYaml tests schema yaml file
func TestSchemaYaml(t *testing.T) {
	config := os.Getenv("FOXDEN_CONFIG")
	if cobj, err := srvConfig.ParseConfig(config); err == nil {
		srvConfig.Config = &cobj
	}

	tmpFile, err := os.CreateTemp(os.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	yamlData := `
- key: pi
  optional: true
  type: string
- key: beam_energy
  optional: false
  type: int
`
	tmpFile.Write([]byte(yamlData))
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// load json data
	fname := tmpFile.Name()
	s := &Schema{FileName: fname}
	err = s.Load()
	if err != nil {
		t.Fatal(err)
	}

	keys, err := s.Keys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Schema keys", keys)
	okeys, err := s.OptionalKeys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Schema optional keys", okeys)

	rec := make(map[string]any)
	rec["pi"] = "person"
	rec["beam_energy"] = 123
	err = s.Validate(rec)
	if err != nil {
		t.Fatal(err)
	}
}

// TestSchemaJson tests schema json file
func TestSchemaJson(t *testing.T) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	jsonData := `[
    {"key": "pi", "type": "string", "optional": true},
    {"key": "beam_energy", "type": "int", "optional": false}
]`
	tmpFile.Write([]byte(jsonData))
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	// load json data
	fname := tmpFile.Name()
	s := &Schema{FileName: fname}
	err = s.Load()
	if err != nil {
		t.Fatal(err)
	}

	keys, err := s.Keys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Schema keys", keys)
	okeys, err := s.OptionalKeys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Schema optional keys", okeys)

	rec := make(map[string]any)
	rec["pi"] = "person"
	rec["beam_energy"] = 123
	err = s.Validate(rec)
	if err != nil {
		t.Fatal(err)
	}
}

// TestLoadSchemaWithInclude tests loading a schema that includes another file
func TestLoadSchemaWithInclude(t *testing.T) {
	// Create temporary directory for schema files
	tempDir := t.TempDir()
	firstPath := filepath.Join(tempDir, "first_schema.json")
	secondPath := filepath.Join(tempDir, "second_schema.json")

	// Content of first_schema.json
	firstSchema := `[
		{
			"key": "did",
			"type": "string",
			"optional": true,
			"multiple": false,
			"section": "User",
			"description": "Dataset IDentifier",
			"units": "",
			"placeholder": "/beamline=demo/btr=user-1234-a/cycle=2025-2/sample_name=testsample"
		}
	]`

	// Content of second_schema.json (includes first)
	secondSchema := fmt.Sprintf("[{ \"file\": \"%s\" },", firstPath)
	secondSchema += `
		{
			"key": "new",
			"type": "string",
			"optional": true,
			"multiple": false,
			"section": "User",
			"description": "New field",
			"units": "",
			"placeholder": "/new"
		}
	]`

	// Write first schema
	if err := os.WriteFile(firstPath, []byte(firstSchema), 0644); err != nil {
		t.Fatalf("Failed to write first schema: %v", err)
	}

	// Write second schema
	if err := os.WriteFile(secondPath, []byte(secondSchema), 0644); err != nil {
		t.Fatalf("Failed to write second schema: %v", err)
	}

	// Load second schema (which includes the first)
	schema := &Schema{FileName: secondPath}
	err := schema.Load()
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	// Expect 2 fields: "did" from first file, "new" from second
	keys, err := schema.Keys()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Schema keys", keys)
	if len(keys) != 2 {
		t.Errorf("Expected 2 schema fields, got %d", len(keys))
	}

	rec := make(map[string]any)
	rec["did"] = "did key"
	rec["new"] = "new key"

	found := map[string]bool{
		"did": false,
		"new": false,
	}

	for _, key := range keys {
		if _, ok := found[key]; ok {
			found[key] = true
		}
	}

	for key, wasFound := range found {
		if !wasFound {
			t.Errorf("Expected key %q not found in parsed schema", key)
		}
	}
}

// TestValidDataValue defines table driven unit tests for validateRecordValue function
func TestValidDataValue(t *testing.T) {
	var floatZero float64
	tests := []struct {
		name     string       // name of the test
		rec      SchemaRecord // schema record
		value    any          // value to check
		expected bool         // expected outcome of the function
	}{
		// === Direct type matches ===
		{"string_ok", SchemaRecord{Type: "string"}, "hello", true},
		{"string_with_spaces_ok", SchemaRecord{Type: "string", Value: "hello bla"}, "hello bla", true},
		{"int_ok", SchemaRecord{Type: "int"}, 42, true},
		{"float_ok", SchemaRecord{Type: "float64"}, 3.14, true},
		{"string_type_vs_int_value", SchemaRecord{Type: "string"}, 123, false},
		{"int_type_vs_string_value", SchemaRecord{Type: "int", Value: 1}, "1", false},
		{"int64_type_vs_string_value", SchemaRecord{Type: "int64", Value: 1}, "1", false},
		{"int64_type_vs_float_value", SchemaRecord{Type: "int64", Value: 0}, floatZero, true},

		// === Non list_str with allowed values ===
		{"string_match_ok", SchemaRecord{Type: "string", Value: "yes"}, "yes", true},
		{"string_match_fail", SchemaRecord{Type: "string", Value: "yes"}, "no", false},
		{"enum_string_ok", SchemaRecord{Type: "string", Value: []any{"a", "b"}}, "a", true},
		{"enum_string_fail", SchemaRecord{Type: "string", Value: []any{"a", "b"}}, "c", false},

		{"enum_int_ok", SchemaRecord{Type: "int", Value: []int{1, 2}}, 1, true},
		{"enum_int_fail", SchemaRecord{Type: "int", Value: []int{1, 2}}, 3, false},
		{"int_match_ok", SchemaRecord{Type: "int", Value: 1}, 1, true},
		{"int_match_fail", SchemaRecord{Type: "int", Value: 1}, 3, false},

		{"enum_float64_ok", SchemaRecord{Type: "float64", Value: []float64{1., 2.}}, 1., true},
		{"enum_float64_fail", SchemaRecord{Type: "float64", Value: []float64{1., 2.}}, 3., false},
		{"float64_match_ok", SchemaRecord{Type: "float64", Value: 1.}, 1., true},
		{"float64_match_fail", SchemaRecord{Type: "float64", Value: 1.}, 3., false},

		// === List type with no restriction ===
		{"list_str_nil_value", SchemaRecord{Type: "list_str", Value: nil}, []string{"a"}, true},

		// === list_str type with allowed values ===
		{"list_str_single_ok", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, "a", true},
		{"list_str_single_fail", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, "c", false},
		{"list_str_slice_ok", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, []string{"a", "b"}, true},
		{"list_str_slice_partial_fail", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, []string{"a", "c"}, false},
		{"list_str_slice_[]any_ok", SchemaRecord{Type: "list_str", Value: []any{"x", "y"}}, []any{"x"}, true},
		{"list_str_slice_[]any_fail", SchemaRecord{Type: "list_str", Value: []any{"x", "y"}}, []any{"z"}, false},

		// === any type ====
		{"any_map_string_float64_ok", SchemaRecord{Type: "any"}, map[string]float64{"m123": 1.23, "m456": 4.56}, true},
		{"any_map_string_string_ok", SchemaRecord{Type: "any"}, map[string]string{"key": "value"}, true},
	}

	// loop over all defined tests and validate function outcome
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validateRecordValue(tt.rec, tt.value, 0)
			if got != tt.expected {
				t.Errorf("validateRecordValue(%+v, %v) = %v, want %v",
					tt.rec, tt.value, got, tt.expected)
			}
		})
	}
}

// TestLoadSchemaWithStruct tests loading a schema that includes another sub-schema records
func TestLoadSchemaWithStruct(t *testing.T) {
	// Create temporary directory for schema files
	tempDir := t.TempDir()
	schemaFile := filepath.Join(tempDir, "schema.json")
	structFile := filepath.Join(tempDir, "struct.json")

	// Content of schema_with_subschema.json
	structRecords := `[
		{
			"key": "int_key",
			"type": "int"
		},
		{
			"key": "str_key",
			"type": "string"
		}
	]`

	// Write schema record to a file
	if err := os.WriteFile(structFile, []byte(structRecords), 0644); err != nil {
		t.Fatalf("Failed to write second schema: %v", err)
	}

	// Content of schema.json which will contain subschema
	schemaRecords := `[
		{
			"key": "did",
			"type": "string"
		},
		{
			"key": "sub",
			"type": "struct",
			"schema": "struct.json"
		}
	]`

	// Write schema record to a file
	if err := os.WriteFile(schemaFile, []byte(schemaRecords), 0644); err != nil {
		t.Fatalf("Failed to write second schema: %v", err)
	}

	// Load second schema (which includes the first)
	s := &Schema{FileName: schemaFile}
	err := s.Load()
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	sub := make(map[string]any)
	sub["int_key"] = 1
	sub["str_key"] = "string_value"
	rec := make(map[string]any)
	rec["did"] = "/path=1/foo=2"
	rec["sub"] = sub
	err = s.Validate(rec)
	if err != nil {
		t.Fatal(err)
	}

	// let's test that we create list of structs and pass it as value
	var subrecords []map[string]any
	sub1 := make(map[string]any)
	sub1["int_key"] = 1
	sub1["str_key"] = "string_value"
	subrecords = append(subrecords, sub1)
	sub2 := make(map[string]any)
	sub2["nokey"] = []int{1, 2, 3}
	subrecords = append(subrecords, sub2)
	nrec := make(map[string]any)
	nrec["did"] = "/path=1/foo=2"
	nrec["sub"] = subrecords
	err = s.Validate(nrec)
	t.Log("we should recieve ERROR from validataion")
	if err == nil {
		t.Logf("subrecords %+v", subrecords)
		t.Fatalf("Used record=%+v with schema=\n%v, fail validation of list_struct type", nrec, schemaRecords)
	}
}

// TestLoadSchemaWithListStruct tests loading a schema that includes another sub-schema records
func TestLoadSchemaWithListStruct(t *testing.T) {
	// Create temporary directory for schema files
	tempDir := t.TempDir()
	schemaFile := filepath.Join(tempDir, "schema.json")
	structFile := filepath.Join(tempDir, "struct.json")

	// Content of schema_with_subschema.json
	structRecords := `[
		{
			"key": "int_key",
			"type": "int"
		},
		{
			"key": "str_key",
			"type": "string"
		}
	]`

	// Write schema record to a file
	if err := os.WriteFile(structFile, []byte(structRecords), 0644); err != nil {
		t.Fatalf("Failed to write second schema: %v", err)
	}

	// Content of schema.json which will contain subschema
	schemaRecords := `[
		{
			"key": "did",
			"type": "string"
		},
		{
			"key": "sub",
			"type": "list_struct",
			"schema": "struct.json"
		}
	]`

	// Write schema record to a file
	if err := os.WriteFile(schemaFile, []byte(schemaRecords), 0644); err != nil {
		t.Fatalf("Failed to write second schema: %v", err)
	}

	// Load second schema (which includes the first)
	s := &Schema{FileName: schemaFile}
	err := s.Load()
	if err != nil {
		t.Fatalf("Failed to load schema: %v", err)
	}

	var subrecords []map[string]any
	sub1 := make(map[string]any)
	sub1["int_key"] = 1
	sub1["str_key"] = "string_value"
	subrecords = append(subrecords, sub1)
	nrec := make(map[string]any)
	nrec["did"] = "/path=1/foo=2"
	nrec["sub"] = subrecords
	err = s.Validate(nrec)
	if err != nil {
		t.Logf("record %+v", nrec)
		t.Fatal(err)
	}

	// assign wrong struct data-type
	rec := make(map[string]any)
	rec["did"] = "/path=1/foo=2"
	rec["sub"] = sub1
	rec["nokey"] = 123
	err = s.Validate(rec)
	t.Log("we should recieve ERROR from validataion")
	if err == nil {
		t.Fatalf("Used record=%+v with schema=\n%v, fail validation of list_struct type", rec, schemaRecords)
	}

}
