// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-netlify/internal/netlify"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SiteResource{}
var _ resource.ResourceWithImportState = &SiteResource{}
var _ resource.ResourceWithConfigure = &SiteResource{}

func NewSiteResource() resource.Resource {
	return &SiteResource{}
}

// SiteResource defines the resource implementation.
type SiteResource struct {
	client *netlify.NetlifyClient
}

// SiteResourceModel describes the resource data model.
type SiteResourceModel struct {
	Id           types.String    `tfsdk:"id"`
	CustomDomain types.String    `tfsdk:"custom_domain"`
	Name         types.String    `tfsdk:"name"`
	Url          types.String    `tfsdk:"url"`
	CreatedAt    types.String    `tfsdk:"created_at"`
	UpdatedAt    types.String    `tfsdk:"updated_at"`
	State        types.String    `tfsdk:"state"`
	Repository   repositoryModel `tfsdk:"repository"`
	LastUpdated  types.String    `tfsdk:"last_updated"`
}

type repositoryModel struct {
	Provider    types.String `tfsdk:"provider"`
	DeployKeyId types.String `tfsdk:"deploy_key_id"`
	RepoPath    types.String `tfsdk:"repo_path"`
	RepoBranch  types.String `tfsdk:"repo_branch"`
	Cmd         types.String `tfsdk:"cmd"`
	Dir         types.String `tfsdk:"dir"`
}

func (r *SiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *SiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Site resource",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"custom_domain": schema.StringAttribute{
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
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
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"repository": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"provider": schema.StringAttribute{
						Required: true,
					},
					"deploy_key_id": schema.StringAttribute{
						Required: true,
					},
					"repo_path": schema.StringAttribute{
						Required: true,
					},
					"repo_branch": schema.StringAttribute{
						Required: true,
					},
					"cmd": schema.StringAttribute{
						Required: true,
					},
					"dir": schema.StringAttribute{
						Required: true,
					},
				},
			},
		},
	}
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netlify.NetlifyClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *NetlifyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tfRepo := data.Repository
	netlifyRepo := netlify.SiteRequest{
		Name:         data.Name.ValueString(),
		CustomDomain: data.Name.ValueString(),
		Repo: netlify.Repository{
			Provider:    tfRepo.Provider.ValueString(),
			Path:        tfRepo.RepoPath.ValueString(),
			Branch:      tfRepo.RepoBranch.ValueString(),
			DeployKeyId: tfRepo.DeployKeyId.ValueString(),
			Cmd:         tfRepo.Cmd.ValueString(),
			Dir:         tfRepo.Dir.ValueString(),
		},
	}
	site, err := r.client.CreateSite(netlifyRepo)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Netlify Site",
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
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	site, err := r.client.GetSite(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Netlify Site",
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

func (r *SiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data SiteResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tfRepo := data.Repository
	netlifyRepo := netlify.SiteRequest{
		Name:         data.Name.ValueString(),
		CustomDomain: data.Name.ValueString(),
		Repo: netlify.Repository{
			Provider:    tfRepo.Provider.ValueString(),
			Path:        tfRepo.RepoPath.ValueString(),
			Branch:      tfRepo.RepoBranch.ValueString(),
			DeployKeyId: tfRepo.DeployKeyId.ValueString(),
			Cmd:         tfRepo.Cmd.ValueString(),
			Dir:         tfRepo.Dir.ValueString(),
		},
	}

	site, err := r.client.UpdateSite(data.Id.ValueString(), netlifyRepo)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Netlify Site",
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
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *SiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SiteResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSite(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete Site, got error: %s", err))
		return
	}
}

func (r *SiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
