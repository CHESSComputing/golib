package utils

import (
	"testing"
)

// TestConvertCamelCaseKeys
func TestConvertCamelCaseKeys(t *testing.T) {
	record := make(map[string]any)
	record["CESRConditions"] = 1
	record["BTR"] = 1
	record["BeamEnergy"] = 1
	rec := ConvertCamelCaseKeys(record)
	expect := []string{"cesr_conditions", "btr", "beam_energy"}
	for k, _ := range rec {
		if !InList(k, expect) {
			t.Errorf("unexpected key '%s' does not belong to %v", k, expect)
		}
	}
}

// TestGetDid
func TestGetDid(t *testing.T) {
	record := make(map[string]any)
	record["Beamline"] = "ID3A"
	record["BTR"] = 1
	record["Cycle"] = 2001
	record["SampleName"] = "test"
	did := GetDid(record)
	if did != "/beamline=ID3A/btr=1/cycle=2001/sample=test" {
		t.Errorf("unable properly construct did='%s' from record %+v", did, record)
	}
}
