package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider                       = &HashProvider{}
	_ provider.ProviderWithFunctions          = &HashProvider{}
	_ provider.ProviderWithEphemeralResources = &HashProvider{}
)

// New returns a new provider implementation.
func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &HashProvider{
			version: version,
			commit:  commit,
		}
	}
}

// HashProviderData is the data available to the resource and data sources.
type HashProviderData struct {
	provider *HashProvider
	Model    *HashProviderModel
}

// HashProviderModel describes the provider data model.
type HashProviderModel struct{}

// HashProvider defines the provider implementation.
type HashProvider struct {
	version string
	commit  string
}

func (p *HashProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hash"
	resp.Version = p.version
}

func (p *HashProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider to support hash based functionality.",
	}
}

// Configure configures the provider.
func (p *HashProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if req.ClientCapabilities.DeferralAllowed && !req.Config.Raw.IsFullyKnown() {
		resp.Deferred = &provider.Deferred{
			Reason: provider.DeferredReasonProviderConfigUnknown,
		}
	}

	// Load the provider config
	model := &HashProviderModel{}
	diags := req.Config.Get(ctx, model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configure provider data
	providerData := &HashProviderData{
		provider: p,
		Model:    model,
	}

	resp.DataSourceData = providerData
	resp.EphemeralResourceData = providerData
	resp.ResourceData = providerData
}

func (p *HashProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewHashRingDataSource,
	}
}

func (p *HashProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *HashProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *HashProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
