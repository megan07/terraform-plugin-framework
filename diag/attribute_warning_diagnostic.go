package diag

import "github.com/hashicorp/terraform-plugin-framework/attrpath"

// NewAttributeWarningDiagnostic returns a new warning severity diagnostic with the given summary, detail, and path.
func NewAttributeWarningDiagnostic(path attrpath.Path, summary string, detail string) DiagnosticWithPath {
	return withPath{
		Diagnostic: NewWarningDiagnostic(summary, detail),
		path:       path,
	}
}
