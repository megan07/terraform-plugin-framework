package fwserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/logging"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// StateGetAttributeValue is a duplicate of tfsdk.State.getAttributeValue,
// except it calls a local duplicate to State.terraformValueAtPath as well.
// It is duplicated to prevent any oddities with trying to use
// tfsdk.State.GetAttribute, which has some potentially undesirable logic.
// Refer to the tfsdk package for the large amount of testing done there.
//
// TODO: Clean up this abstraction back into an internal State type method.
// The extra State parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/215
func StateGetAttributeValue(ctx context.Context, s tfsdk.State, path *tftypes.AttributePath) (attr.Value, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrType, err := s.Schema.AttributeTypeAtPath(path)
	if err != nil {
		err = fmt.Errorf("error getting attribute type in schema: %w", err)
		diags.AddAttributeError(
			path,
			"State Read Error",
			"An unexpected error was encountered trying to read an attribute from the state. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	// if the whole state is nil, the value of a valid attribute is also nil
	if s.Raw.IsNull() {
		return nil, nil
	}

	tfValue, err := StateTerraformValueAtPath(s, path)

	// Ignoring ErrInvalidStep will allow this method to return a null value of the type.
	if err != nil && !errors.Is(err, tftypes.ErrInvalidStep) {
		diags.AddAttributeError(
			path,
			"State Read Error",
			"An unexpected error was encountered trying to read an attribute from the state. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	// TODO: If ErrInvalidStep, check parent paths for unknown value.
	//       If found, convert this value to an unknown value.
	// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/186

	if attrTypeWithValidate, ok := attrType.(attr.TypeWithValidate); ok {
		logging.FrameworkTrace(ctx, "Type implements TypeWithValidate")
		logging.FrameworkDebug(ctx, "Calling provider defined Type Validate")
		diags.Append(attrTypeWithValidate.Validate(ctx, tfValue, path)...)
		logging.FrameworkDebug(ctx, "Called provider defined Type Validate")

		if diags.HasError() {
			return nil, diags
		}
	}

	attrValue, err := attrType.ValueFromTerraform(ctx, tfValue)

	if err != nil {
		diags.AddAttributeError(
			path,
			"State Read Error",
			"An unexpected error was encountered trying to read an attribute from the state. This is always an error in the provider. Please report the following to the provider developer:\n\n"+err.Error(),
		)
		return nil, diags
	}

	return attrValue, diags
}

// StateTerraformValueAtPath is a duplicate of
// tfsdk.State.terraformValueAtPath.
//
// TODO: Clean up this abstraction back into an internal State type method.
// The extra State parameter is a carry-over of creating the proto6server
// package from the tfsdk package and not wanting to export the method.
// Reference: https://github.com/hashicorp/terraform-plugin-framework/issues/215
func StateTerraformValueAtPath(s tfsdk.State, path *tftypes.AttributePath) (tftypes.Value, error) {
	rawValue, remaining, err := tftypes.WalkAttributePath(s.Raw, path)
	if err != nil {
		return tftypes.Value{}, fmt.Errorf("%v still remains in the path: %w", remaining, err)
	}
	attrValue, ok := rawValue.(tftypes.Value)
	if !ok {
		return tftypes.Value{}, fmt.Errorf("got non-tftypes.Value result %v", rawValue)
	}
	return attrValue, err
}
