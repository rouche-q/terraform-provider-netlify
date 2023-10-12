package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type siteDataSource struct {
	client *NetlifyClient
}

type siteDataSourceModel struct {
	Id           types.String `tfsdk:"id"`
	CustomDomain types.String `tfsdk:"custom_domain"`
	Name         types.String `tfsdk:"name"`
	Url          types.String `tfsdk:"url"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
	State        types.String `tfsdk:"state"`
}

var (
	_ datasource.DataSource              = &siteDataSource{}
	_ datasource.DataSourceWithConfigure = &siteDataSource{}
)

func NewSiteDataSource() datasource.DataSource {
	return &siteDataSource{}
}

func (d *siteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (d *siteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*NetlifyClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NetlifyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *siteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed: true,
			},
			"custom_domain": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Required: true,
			},
			"url": schema.StringAttribute{
				Computed: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *siteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data siteDataSourceModel
	tflog.Debug(ctx, "Preparing to read site data source")

	diag := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := d.client.GetSite(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Netlify site",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(site.Id)
	data.Name = types.StringValue(site.Name)
	data.CustomDomain = types.StringValue(site.CustomDomain)
	data.Url = types.StringValue(site.Url)
	data.State = types.StringValue(site.State)
	data.CreatedAt = types.StringValue(site.CreatedAt)
	data.UpdatedAt = types.StringValue(site.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
