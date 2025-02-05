package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ProviderMeta returns the *tfsdk.Config for a *tfprotov6.DynamicValue and
// *tfsdk.Schema. This data handling is different than Config to simplify
// implementors, in that:
//
//     - Missing Schema will return nil, rather than an error
//     - Missing DynamicValue will return nil typed Value, rather than an error
func ProviderMeta(ctx context.Context, proto6DynamicValue *tfprotov6.DynamicValue, schema *tfsdk.Schema) (*tfsdk.Config, diag.Diagnostics) {
	if schema == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	fw := &tfsdk.Config{
		Raw:    tftypes.NewValue(schema.TerraformType(ctx), nil),
		Schema: *schema,
	}

	if proto6DynamicValue == nil {
		return fw, nil
	}

	proto6Value, err := proto6DynamicValue.Unmarshal(schema.TerraformType(ctx))

	if err != nil {
		diags.AddError(
			"Unable to Convert Provider Meta Configuration",
			"An unexpected error was encountered when converting the provider meta configuration from the protocol type. "+
				"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+err.Error(),
		)

		return nil, diags
	}

	fw.Raw = proto6Value

	return fw, nil
}
