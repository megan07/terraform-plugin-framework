package fromproto6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Plan returns the *tfsdk.Plan for a *tfprotov6.DynamicValue and
// *tfsdk.Schema.
func Plan(ctx context.Context, proto6DynamicValue *tfprotov6.DynamicValue, schema *tfsdk.Schema) (*tfsdk.Plan, diag.Diagnostics) {
	if proto6DynamicValue == nil {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Panic prevention here to simplify the calling implementations.
	// This should not happen, but just in case.
	if schema == nil {
		diags.AddError(
			"Unable to Convert Plan",
			"An unexpected error was encountered when converting the plan from the protocol type. "+
				"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+
				"Missing schema.",
		)

		return nil, diags
	}

	proto6Value, err := proto6DynamicValue.Unmarshal(schema.TerraformType(ctx))

	if err != nil {
		diags.AddError(
			"Unable to Convert Plan",
			"An unexpected error was encountered when converting the plan from the protocol type. "+
				"This is always an issue in the Terraform Provider SDK used to implement the provider and should be reported to the provider developers.\n\n"+
				"Please report this to the provider developer:\n\n"+err.Error(),
		)

		return nil, diags
	}

	fw := &tfsdk.Plan{
		Raw:    proto6Value,
		Schema: *schema,
	}

	return fw, nil
}
