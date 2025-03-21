package doi

import (
	"os"
	"testing"
)

// TestRenderTemplateDefault tests the function using the default template.
func TestRenderTemplateDefault(t *testing.T) {
	data := DOIData{
		PI: "pi", Affiliation: "affiliation", Beamline: "beamline", StaffScientist: "scientist", Facility: "facility",
	}

	expected := `<html><body>
PI: pi
<br/>
Facility: facility
<br/>
Beamline: beamline
<br/>
Affiliation: affiliation
<br/>
StaffScientist: scientist
</body></html>`

	result, err := RenderTemplate("", data)
	if err != nil {
		t.Fatalf("RenderTemplate failed with error: %v", err)
	}

	if result != expected {
		t.Errorf("Unexpected output:\nGot:\n%q\nExpected:\n%q", result, expected)
	}
}

// TestRenderTemplateFile tests the function using a template file.
func TestRenderTemplateFile(t *testing.T) {
	// Create a temporary template file
	tmplContent := `Pi: {{.PI}} Affiliation: {{.Affiliation}}`
	tmpFile, err := os.CreateTemp("", "template_*.tmpl")
	if err != nil {
		t.Fatalf("Failed to create temporary template file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Write the template content to the file
	if _, err := tmpFile.WriteString(tmplContent); err != nil {
		t.Fatalf("Failed to write to temporary template file: %v", err)
	}
	tmpFile.Close()

	data := DOIData{
		PI: "PI", Affiliation: "University", Beamline: "beamline", StaffScientist: "scientist",
	}

	expected := "Pi: PI Affiliation: University"

	result, err := RenderTemplate(tmpFile.Name(), data)
	if err != nil {
		t.Fatalf("RenderTemplate failed with error: %v", err)
	}

	if result != expected {
		t.Errorf("Unexpected output:\nGot:\n%q\nExpected:\n%q", result, expected)
	}
}

// TestRenderTemplateInvalidFile tests the function with a nonexistent template file.
func TestRenderTemplateInvalidFile(t *testing.T) {
	_, err := RenderTemplate("nonexistent.tmpl", DOIData{})
	if err == nil {
		t.Fatal("Expected an error for a nonexistent template file, but got none")
	}
}
