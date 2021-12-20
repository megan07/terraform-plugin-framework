package tfsdk

import (
	"context"
	"fmt"
	"math/big"
	goreflect "reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attrpath"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/internal/reflect"
	testtypes "github.com/hashicorp/terraform-plugin-framework/internal/testing/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func newStringPointer(in string) *string {
	return &in
}

func newStringPointerPointer(in string) **string {
	stringPointer := &in
	return &stringPointer
}

func newInt64Pointer(in int64) *int64 {
	return &in
}

func TestValueAs(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val           attr.Value
		target        interface{}
		expected      interface{}
		expectedDiags diag.Diagnostics
	}

	tests := map[string]testCase{
		"primitive": {
			val:      types.String{Value: "hello"},
			target:   newStringPointer(""),
			expected: newStringPointer("hello"),
		},
		"incompatible-type": {
			val:    types.String{Value: "hello"},
			target: newInt64Pointer(0),
			expectedDiags: diag.Diagnostics{
				diag.WithPath(
					attrpath.New(),
					reflect.DiagIntoIncompatibleType{
						Val:        tftypes.NewValue(tftypes.String, "hello"),
						TargetType: goreflect.TypeOf(int64(0)),
						Err:        fmt.Errorf("can't unmarshal %s into %T, expected *big.Float", tftypes.String, big.NewFloat(0)),
					},
				),
			},
		},
		"different-type": {
			val:    types.String{Value: "hello"},
			target: &testtypes.String{},
			expectedDiags: diag.Diagnostics{
				diag.WithPath(
					attrpath.New(),
					reflect.DiagNewAttributeValueIntoWrongType{
						ValType:    goreflect.TypeOf(types.String{Value: "hello"}),
						TargetType: goreflect.TypeOf(testtypes.String{}),
						SchemaType: types.StringType,
					},
				),
			},
		},
	}

	for name, tc := range tests {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := ValueAs(context.Background(), tc.val, tc.target)

			if diff := cmp.Diff(tc.expectedDiags, diags); diff != "" {
				t.Fatalf("Unexpected diff in diagnostics (-wanted, +got): %s", diff)
			}

			if diags.HasError() {
				return
			}

			if diff := cmp.Diff(tc.expected, tc.target); diff != "" {
				t.Fatalf("Unexpected diff in results (-wanted, +got): %s", diff)
			}
		})
	}
}

func TestValueAs_generic(t *testing.T) {
	t.Parallel()

	var target attr.Value
	val := types.String{Value: "hello"}
	diags := ValueAs(context.Background(), val, &target)
	if len(diags) > 0 {
		t.Fatalf("Unexpected diagnostics: %s", diags)
	}
	if !val.Equal(target.(attr.Value)) {
		t.Errorf("Expected target to be %v, got %v", val, target)
	}
}
