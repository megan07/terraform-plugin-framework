package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ attr.TypeWithValidate = NumberTypeWithValidateError{}
	_ attr.TypeWithValidate = NumberTypeWithValidateWarning{}
)

type NumberTypeWithValidateError struct {
	NumberType
}

type NumberTypeWithValidateWarning struct {
	NumberType
}

func (t NumberTypeWithValidateError) Validate(ctx context.Context, in tftypes.Value, path *tftypes.AttributePath) diag.Diagnostics {
	return diag.Diagnostics{TestErrorDiagnostic(path)}
}

func (t NumberTypeWithValidateWarning) Validate(ctx context.Context, in tftypes.Value, path *tftypes.AttributePath) diag.Diagnostics {
	return diag.Diagnostics{TestWarningDiagnostic(path)}
}

func (n NumberTypeWithValidateError) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	res, err := n.NumberType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	newNumber := res.(Number)
	newNumber.CreatedBy = n
	return newNumber, nil
}

func (n NumberTypeWithValidateWarning) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	res, err := n.NumberType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	newNumber := res.(Number)
	newNumber.CreatedBy = n
	return newNumber, nil
}
