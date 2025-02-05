package proto6server

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type testServeDataSourceTypeTwo struct{}

func (dt testServeDataSourceTypeTwo) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"family": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (dt testServeDataSourceTypeTwo) NewDataSource(_ context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, ok := p.(*testServeProvider)
	if !ok {
		prov, ok := p.(*testServeProviderWithMetaSchema)
		if !ok {
			panic(fmt.Sprintf("unexpected provider type %T", p))
		}
		provider = prov.testServeProvider
	}
	return testServeDataSourceTwo{
		provider: provider,
	}, nil
}

var testServeDataSourceTypeTwoType = tftypes.Object{
	AttributeTypes: map[string]tftypes.Type{
		"family": tftypes.String,
		"name":   tftypes.String,
		"id":     tftypes.String,
	},
}

type testServeDataSourceTwo struct {
	provider *testServeProvider
}

func (r testServeDataSourceTwo) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	r.provider.readDataSourceConfigValue = req.Config.Raw
	r.provider.readDataSourceConfigSchema = req.Config.Schema
	r.provider.readDataSourceProviderMetaValue = req.ProviderMeta.Raw
	r.provider.readDataSourceProviderMetaSchema = req.ProviderMeta.Schema
	r.provider.readDataSourceCalledDataSourceType = "test_two"
	r.provider.readDataSourceImpl(ctx, req, resp)
}
