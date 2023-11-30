package provider

import (
	"context"
	"fmt"
	"terraform-provider-netlify/internal/netlify"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CurrentUserDataSource struct {
	client *netlify.NetlifyClient
}

type CurrentUserDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Uid         types.String `tfsdk:"uid"`
	Slug        types.String `tfsdk:"slug"`
	FullName    types.String `tfsdk:"full_name"`
	AvatarUrl   types.String `tfsdk:"avatar_url"`
	Email       types.String `tfsdk:"email"`
	AffiliateId types.String `tfsdk:"affiliate_id"`
	SiteCount   types.Int64  `tfsdk:"site_count"`
	CreatedAt   types.String `tfsdk:"created_at"`
	LastLogin   types.String `tfsdk:"last_login"`
}

var (
	_ datasource.DataSource              = &CurrentUserDataSource{}
	_ datasource.DataSourceWithConfigure = &CurrentUserDataSource{}
)

func NewCurrentUserDataSource() datasource.DataSource {
	return &CurrentUserDataSource{}
}

func (d *CurrentUserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_current_user"
}

func (d *CurrentUserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netlify.NetlifyClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *NetlifyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *CurrentUserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Current user DataSource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"uid": schema.StringAttribute{
				Computed: true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
			},
			"full_name": schema.StringAttribute{
				Computed: true,
			},
			"avatar_url": schema.StringAttribute{
				Computed: true,
			},
			"email": schema.StringAttribute{
				Computed: true,
			},
			"affiliate_id": schema.StringAttribute{
				Computed: true,
			},
			"site_count": schema.Int64Attribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"last_login": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *CurrentUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CurrentUserDataSourceModel
	tflog.Debug(ctx, "Preparing to read CurrentUser data source")

	diag := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	currentUser, err := d.client.GetCurrentUser()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Netlify CurrentUser",
			err.Error(),
		)
		return
	}

	data.Id = types.StringValue(currentUser.Email)
	data.Uid = types.StringValue(currentUser.Uid)
	data.Slug = types.StringValue(currentUser.Slug)
	data.FullName = types.StringValue(currentUser.FullName)
	data.AvatarUrl = types.StringValue(currentUser.AvatarUrl)
	data.Email = types.StringValue(currentUser.Email)
	data.AffiliateId = types.StringValue(currentUser.AffiliateId)
	data.SiteCount = types.Int64Value(int64(currentUser.SiteCount))
	data.CreatedAt = types.StringValue(currentUser.CreatedAt)
	data.LastLogin = types.StringValue(currentUser.LastLogin)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
