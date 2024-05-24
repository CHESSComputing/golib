package ql

import (
	"testing"
)

// TestQLRecord
func TestQLRecord(t *testing.T) {
	Verbose = 1
	qrec := QLRecord{Key: "test"}
	if qrec.Details("key") != "test" {
		t.Errorf("unexpected ql key: %+v", qrec)
	}
	if qrec.Details("units") != "N/A" {
		t.Errorf("unexpected ql key: %+v", qrec)
	}
	qrec.Units = "units"
	if qrec.Details("units") != "units" {
		t.Errorf("unexpected ql key: %+v", qrec)
	}
}
