// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-netlify/internal/netlify"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &EnvVarResource{}
	_ resource.ResourceWithImportState = &EnvVarResource{}
	_ resource.ResourceWithConfigure   = &EnvVarResource{}
)

func NewEnvVarRessource() resource.Resource {
	return &EnvVarResource{}
}

// EnvVarsRessource defines the resource implementation.
type EnvVarResource struct {
	client *netlify.NetlifyClient
}

// EnvVarsRessourceModel describes the resource data model.
type EnvVarResourceModel struct {
	AccountSlug types.String `tfsdk:"account_slug"`
	SiteId      types.String `tfsdk:"site_id"`
	Key         types.String `tfsdk:"key"`
	Scopes      types.List   `tfsdk:"scopes"`
	IsSecret    types.Bool   `tfsdk:"is_secret"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

func (r *EnvVarResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_env_var"
}

func (r *EnvVarResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "EnvVar resource",

		Attributes: map[string]schema.Attribute{
			"account_slug": schema.StringAttribute{
				Required: true,
			},
			"site_id": schema.StringAttribute{
				Required: true,
			},
			"key": schema.StringAttribute{
				Required: true,
			},
			"scopes": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"is_secret": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *EnvVarResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*netlify.NetlifyClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *netlify.NetlifyClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *EnvVarResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EnvVarResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	reqEnvVar := netlify.EnvVar{
		Key: data.Key.ValueString(),
		Values: []netlify.EnvVarValue{
			{Value: "Terraform Placeholder"},
		},
		IsSecret: data.IsSecret.ValueBool(),
	}

	for _, scope := range data.Scopes.Elements() {
		reqEnvVar.Scopes = append(reqEnvVar.Scopes, scope.String())
	}

	res, err := r.client.CreateEnvVar(data.AccountSlug.ValueString(), data.SiteId.ValueString(), reqEnvVar)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Netlify Env variable",
			err.Error())
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, res.Scopes)
	if diags.HasError() {
		return
	}

	data.Scopes = scopeList
	data.IsSecret = types.BoolValue(res.IsSecret)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EnvVarResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	res, err := r.client.GetEnvVar(data.AccountSlug.ValueString(), data.SiteId.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Netlify Env variable",
			err.Error())
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, res.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Scopes = scopeList
	data.IsSecret = types.BoolValue(res.IsSecret)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EnvVarResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var slug, siteId, key string
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("account_slug"), &slug)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("site_id"), &siteId)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("key"), &key)...)
	if resp.Diagnostics.HasError() {
		return
	}

	envVar, err := r.client.GetEnvVar(slug, siteId, key)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Netlify Env variable",
			err.Error())
	}

	for _, scope := range data.Scopes.Elements() {
		contain := false
		for _, v := range envVar.Scopes {
			if scope.String() == v {
				contain = true
			}
		}
		if !contain {
			envVar.Scopes = append(envVar.Scopes, scope.String())
		}
	}
	envVar.Key = data.Key.ValueString()
	envVar.IsSecret = data.IsSecret.ValueBool()

	res, err := r.client.UpdateEnvVar(slug, siteId, key, *envVar)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Update Netlify Env variable",
			err.Error())
	}

	scopeList, diags := types.ListValueFrom(ctx, types.StringType, res.Scopes)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Scopes = scopeList
	data.IsSecret = types.BoolValue(res.IsSecret)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EnvVarResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEnvVar(data.AccountSlug.ValueString(), data.SiteId.ValueString(), data.Key.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete EnvVarResource",
			err.Error(),
		)
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EnvVarResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
