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

func TestValidDataValue(t *testing.T) {
	tests := []struct {
		name     string
		rec      SchemaRecord
		value    any
		expected bool
	}{
		// === Direct type matches ===
		{"string ok", SchemaRecord{Type: "string"}, "hello", true},
		{"string with spaces ok", SchemaRecord{Type: "string", Value: "hello bla"}, "hello bla", true},
		{"int ok", SchemaRecord{Type: "int"}, 42, true},
		{"float ok", SchemaRecord{Type: "float64"}, 3.14, true},
		{"string vs int", SchemaRecord{Type: "string"}, 123, false},
		{"int vs string", SchemaRecord{Type: "int", Value: 1}, "1", false},
		{"int64 vs string", SchemaRecord{Type: "int64", Value: 1}, "1", false},

		// === Non list_str with allowed values ===
		{"enum single ok", SchemaRecord{Type: "string", Value: "yes"}, "yes", true},
		{"enum single fail", SchemaRecord{Type: "string", Value: "yes"}, "no", false},
		{"enum string ok", SchemaRecord{Type: "string", Value: []any{"a", "b"}}, "a", true},
		{"enum string fail", SchemaRecord{Type: "string", Value: []any{"a", "b"}}, "c", false},

		{"enum int ok", SchemaRecord{Type: "int", Value: []int{1, 2}}, 1, true},
		{"enum int fail", SchemaRecord{Type: "int", Value: []int{1, 2}}, 3, false},
		{"enum float64 ok", SchemaRecord{Type: "float64", Value: []float64{1., 2.}}, 1., true},
		{"enum float64 fail", SchemaRecord{Type: "float64", Value: []float64{1., 2.}}, 3., false},

		// === List type with no restriction ===
		{"list_str nil value", SchemaRecord{Type: "list_str", Value: nil}, []string{"a"}, true},

		// === list_str type with allowed values ===
		{"list_str single ok", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, "a", true},
		{"list_str single fail", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, "c", false},
		{"list_str slice ok", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, []string{"a", "b"}, true},
		{"list_str slice partial fail", SchemaRecord{Type: "list_str", Value: []any{"a", "b"}}, []string{"a", "c"}, false},
		{"list_str slice []any ok", SchemaRecord{Type: "list_str", Value: []any{"x", "y"}}, []any{"x"}, true},
		{"list_str slice []any fail", SchemaRecord{Type: "list_str", Value: []any{"x", "y"}}, []any{"z"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := validDataValue(tt.rec, tt.value, 0)
			if got != tt.expected {
				t.Errorf("validDataValue(%+v, %v) = %v, want %v",
					tt.rec, tt.value, got, tt.expected)
			}
		})
	}
}
