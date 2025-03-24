package beamlines

import (
	"fmt"
	"os"
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
