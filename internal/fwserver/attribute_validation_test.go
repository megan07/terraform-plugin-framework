package fwserver

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	testtypes "github.com/hashicorp/terraform-plugin-framework/internal/testing/types"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAttributeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		req  tfsdk.ValidateAttributeRequest
		resp tfsdk.ValidateAttributeResponse
	}{
		"no-attributes-or-type": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Invalid Attribute Definition",
						"Attribute must define either Attributes or Type. This is always a problem with the provider and should be reported to the provider developer.",
					),
				},
			},
		},
		"both-attributes-and-type": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
									"testing": {
										Type:     types.StringType,
										Optional: true,
									},
								}),
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Invalid Attribute Definition",
						"Attribute cannot define both Attributes and Type. This is always a problem with the provider and should be reported to the provider developer.",
					),
				},
			},
		},
		"missing-required-optional-and-computed": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type: types.StringType,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Invalid Attribute Definition",
						"Attribute missing Required, Optional, or Computed definition. This is always a problem with the provider and should be reported to the provider developer.",
					),
				},
			},
		},
		"config-error": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.ListType{ElemType: types.StringType},
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeErrorDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Configuration Read Error",
						"An unexpected error was encountered trying to read an attribute from the configuration. This is always an error in the provider. Please report the following to the provider developer:\n\n"+
							"can't use tftypes.String<\"testvalue\"> as value of List with ElementType types.primitive, can only use tftypes.String values",
					),
				},
			},
		},
		"no-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"deprecation-message-known": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:               types.StringType,
								Optional:           true,
								DeprecationMessage: "Use something else instead.",
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeWarningDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Attribute Deprecated",
						"Use something else instead.",
					),
				},
			},
		},
		"deprecation-message-null": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, nil),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:               types.StringType,
								Optional:           true,
								DeprecationMessage: "Use something else instead.",
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"deprecation-message-unknown": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:               types.StringType,
								Optional:           true,
								DeprecationMessage: "Use something else instead.",
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					diag.NewAttributeWarningDiagnostic(
						tftypes.NewAttributePath().WithAttributeName("test"),
						"Attribute Deprecated",
						"Use something else instead.",
					),
				},
			},
		},
		"warnings": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								Validators: []tfsdk.AttributeValidator{
									testWarningAttributeValidator{},
									testWarningAttributeValidator{},
								},
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testWarningDiagnostic1,
					testWarningDiagnostic2,
				},
			},
		},
		"errors": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     types.StringType,
								Required: true,
								Validators: []tfsdk.AttributeValidator{
									testErrorAttributeValidator{},
									testErrorAttributeValidator{},
								},
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testErrorDiagnostic1,
					testErrorDiagnostic2,
				},
			},
		},
		"type-with-validate-error": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     testtypes.StringTypeWithValidateError{},
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testtypes.TestErrorDiagnostic(tftypes.NewAttributePath().WithAttributeName("test")),
				},
			},
		},
		"type-with-validate-warning": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(tftypes.Object{
						AttributeTypes: map[string]tftypes.Type{
							"test": tftypes.String,
						},
					}, map[string]tftypes.Value{
						"test": tftypes.NewValue(tftypes.String, "testvalue"),
					}),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Type:     testtypes.StringTypeWithValidateWarning{},
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testtypes.TestWarningDiagnostic(tftypes.NewAttributePath().WithAttributeName("test")),
				},
			},
		},
		"nested-attr-list-no-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.List{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.List{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								[]tftypes.Value{
									tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"nested-attr-list-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.List{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.List{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								[]tftypes.Value{
									tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
										Validators: []tfsdk.AttributeValidator{
											testErrorAttributeValidator{},
										},
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testErrorDiagnostic1,
				},
			},
		},
		"nested-attr-map-no-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Map{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Map{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								map[string]tftypes.Value{
									"testkey": tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"nested-attr-map-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Map{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Map{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								map[string]tftypes.Value{
									"testkey": tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.MapNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
										Validators: []tfsdk.AttributeValidator{
											testErrorAttributeValidator{},
										},
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testErrorDiagnostic1,
				},
			},
		},
		"nested-attr-set-no-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Set{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Set{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								[]tftypes.Value{
									tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"nested-attr-set-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Set{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Set{
									ElementType: tftypes.Object{
										AttributeTypes: map[string]tftypes.Type{
											"nested_attr": tftypes.String,
										},
									},
								},
								[]tftypes.Value{
									tftypes.NewValue(
										tftypes.Object{
											AttributeTypes: map[string]tftypes.Type{
												"nested_attr": tftypes.String,
											},
										},
										map[string]tftypes.Value{
											"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
										},
									),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
										Validators: []tfsdk.AttributeValidator{
											testErrorAttributeValidator{},
										},
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testErrorDiagnostic1,
				},
			},
		},
		"nested-attr-single-no-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"nested_attr": tftypes.String,
									},
								},
							},
						},
						map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"nested_attr": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{},
		},
		"nested-attr-single-validation": {
			req: tfsdk.ValidateAttributeRequest{
				AttributePath: tftypes.NewAttributePath().WithAttributeName("test"),
				Config: tfsdk.Config{
					Raw: tftypes.NewValue(
						tftypes.Object{
							AttributeTypes: map[string]tftypes.Type{
								"test": tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"nested_attr": tftypes.String,
									},
								},
							},
						}, map[string]tftypes.Value{
							"test": tftypes.NewValue(
								tftypes.Object{
									AttributeTypes: map[string]tftypes.Type{
										"nested_attr": tftypes.String,
									},
								},
								map[string]tftypes.Value{
									"nested_attr": tftypes.NewValue(tftypes.String, "testvalue"),
								},
							),
						},
					),
					Schema: tfsdk.Schema{
						Attributes: map[string]tfsdk.Attribute{
							"test": {
								Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
									"nested_attr": {
										Type:     types.StringType,
										Required: true,
										Validators: []tfsdk.AttributeValidator{
											testErrorAttributeValidator{},
										},
									},
								}),
								Required: true,
							},
						},
					},
				},
			},
			resp: tfsdk.ValidateAttributeResponse{
				Diagnostics: diag.Diagnostics{
					testErrorDiagnostic1,
				},
			},
		},
	}

	for name, tc := range testCases {
		name, tc := name, tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var got tfsdk.ValidateAttributeResponse
			attribute, err := tc.req.Config.Schema.AttributeAtPath(tc.req.AttributePath)

			if err != nil {
				t.Fatalf("Unexpected error getting %s", err)
			}

			AttributeValidate(context.Background(), attribute, tc.req, &got)

			if diff := cmp.Diff(got, tc.resp); diff != "" {
				t.Errorf("Unexpected response (+wanted, -got): %s", diff)
			}
		})
	}
}

var (
	testErrorDiagnostic1 = diag.NewErrorDiagnostic(
		"Error Diagnostic 1",
		"This is an error.",
	)
	testErrorDiagnostic2 = diag.NewErrorDiagnostic(
		"Error Diagnostic 2",
		"This is an error.",
	)
	testWarningDiagnostic1 = diag.NewWarningDiagnostic(
		"Warning Diagnostic 1",
		"This is a warning.",
	)
	testWarningDiagnostic2 = diag.NewWarningDiagnostic(
		"Warning Diagnostic 2",
		"This is a warning.",
	)
)

type testErrorAttributeValidator struct {
	tfsdk.AttributeValidator
}

func (v testErrorAttributeValidator) Description(ctx context.Context) string {
	return "validation that always returns an error"
}

func (v testErrorAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v testErrorAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if len(resp.Diagnostics) == 0 {
		resp.Diagnostics.Append(testErrorDiagnostic1)
	} else {
		resp.Diagnostics.Append(testErrorDiagnostic2)
	}
}

type testWarningAttributeValidator struct {
	tfsdk.AttributeValidator
}

func (v testWarningAttributeValidator) Description(ctx context.Context) string {
	return "validation that always returns a warning"
}

func (v testWarningAttributeValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v testWarningAttributeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if len(resp.Diagnostics) == 0 {
		resp.Diagnostics.Append(testWarningDiagnostic1)
	} else {
		resp.Diagnostics.Append(testWarningDiagnostic2)
	}
}
