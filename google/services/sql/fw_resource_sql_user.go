// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package sql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-provider-google/google/fwmodels"
	"github.com/hashicorp/terraform-provider-google/google/fwresource"
	"github.com/hashicorp/terraform-provider-google/google/fwtransport"
	"github.com/hashicorp/terraform-provider-google/google/transport"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

var (
	_ resource.Resource              = &SQLUserFWResource{}
	_ resource.ResourceWithConfigure = &SQLUserFWResource{}
)

func NewSQLUserFWResource() resource.Resource {
	return &SQLUserFWResource{}
}

type SQLUserFWResource struct {
	client         *sqladmin.Service
	providerConfig *transport_tpg.Config
}

type SQLUserModel struct {
	Id       types.String `tfsdk:"id"`
	Project  types.String `tfsdk:"project"`
	Name     types.String `tfsdk:"name"`
	Host     types.String `tfsdk:"host"`
	Instance types.String `tfsdk:"instance"`
	Password types.String `tfsdk:"password"`
	// PasswordWO       types.String `tfsdk:"password_wo"`
	// PasswordWOVersion       types.String `tfsdk:"password_wo_version"`
	Type types.String `tfsdk:"type"`
	// SqlServerUserDetails  types.List `tfsdk:"sql_server_user_details"`
	// PasswordPolicy  types.List `tfsdk:"password_policy"`
	// DeletionPolicy       types.String `tfsdk:"deletion_policy"`
	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

// Metadata returns the resource type name.
func (d *SQLUserFWResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fw_sql_user"
}

func (r *SQLUserFWResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*transport_tpg.Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *transport_tpg.Config, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = p.NewSqlAdminClient(p.UserAgent)
	if resp.Diagnostics.HasError() {
		return
	}
	r.providerConfig = p
}

func (d *SQLUserFWResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "A resource to represent a SQL User object.",

		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: `The name of the user. Changing this forces a new resource to be created.`,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					SQLUserNameIAMPlanModifier(),
				},
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"type": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					// TODO DiffSuppressFunc: tpgresource.EmptyOrDefaultStringSuppress("BUILT_IN"),
				},
			},
			// This is included for backwards compatibility with the original, SDK-implemented resource.
			"id": schema.StringAttribute{
				Description:         "Project identifier",
				MarkdownDescription: "Project identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
			}),
		},
	}
}

func (r *SQLUserFWResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SQLUserModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := fwresource.GetProjectFramework(data.Project, types.StringValue(r.providerConfig.Project), &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	nameData, diags := data.Name.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	instanceData, diags := data.Instance.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostData, diags := data.Host.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	typeData, diags := data.Type.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	passwordData, diags := data.Password.ToStringValue(ctx)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	user := &sqladmin.User{
		Name:     nameData.ValueString(),
		Instance: instanceData.ValueString(),
		Password: passwordData.ValueString(),
		Host:     hostData.ValueString(),
		Type:     typeData.ValueString(),
	}

	transport_tpg.MutexStore.Lock(instanceMutexKey(project.ValueString(), instanceData.ValueString()))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project.ValueString(), instanceData.ValueString()))

	r.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, r.client.UserAgent)

	// TODO host check logic

	var op *sqladmin.Operation
	var err error
	insertFunc := func() error {
		op, err = r.client.Users.Insert(project.ValueString(), instanceData.ValueString(),
			user).Do()
		return err
	}
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: insertFunc,
		Timeout:   createTimeout,
	})

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error, failed to insert "+
			"user %s into instance %s", nameData.ValueString(), instanceData.ValueString()), err.Error())
		return
	}

	err = SqlAdminOperationWaitTime(r.providerConfig, op, project.ValueString(), "Insert User", r.client.UserAgent, createTimeout)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error, failure waiting to insert "+
			"user %s into instance %s", nameData.ValueString(), instanceData.ValueString()), err.Error())
		return
	}

	tflog.Trace(ctx, "created sql user resource")

	// This will include a double-slash (//) for postgres instances,
	// for which user.Host is an empty string.  That's okay.
	data.Id = types.StringValue(fmt.Sprintf("%s/%s/%s", user.Name, user.Host, user.Instance))
	data.Project = project

	// read back sql user
	r.SQLUserRefresh(ctx, &data, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SQLUserFWResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data SQLUserModel
	var metaData *fwmodels.ProviderMetaModel

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	r.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, r.client.UserAgent)

	tflog.Trace(ctx, "read sql user resource")

	// read back sql user
	r.SQLUserRefresh(ctx, &data, &resp.State, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SQLUserFWResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var old, new SQLUserModel
	var metaData *fwmodels.ProviderMetaModel

	resp.Diagnostics.Append(req.State.Get(ctx, &old)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &new)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use provider_meta to set User-Agent
	r.client.UserAgent = fwtransport.GenerateFrameworkUserAgentString(metaData, r.client.UserAgent)

	if !old.Password.Equal(new.Password) {
		project := new.Project.ValueString()
		instance := new.Instance.ValueString()
		name := new.Name.ValueString()
		host := new.Host.ValueString()
		password := new.Password.ValueString()

		updateTimeout, diags := new.Timeouts.Update(ctx, 20*time.Minute)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		user := &sqladmin.User{
			Name:     name,
			Instance: instance,
			Password: password,
		}
		transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
		defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))
		var op *sqladmin.Operation
		var err error
		updateFunc := func() error {
			op, err = r.client.Users.Update(project, instance, user).Host(host).Name(name).Do()
			return err
		}
		err = transport_tpg.Retry(transport_tpg.RetryOptions{
			RetryFunc: updateFunc,
			Timeout:   updateTimeout,
		})

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failed to update"+
				"user %s in instance %s", name, instance), err.Error())
			return
		}

		err = SqlAdminOperationWaitTime(r.providerConfig, op, project, "Update User", r.client.UserAgent, updateTimeout)

		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("failure waiting for update"+
				"user %s in instance %s", name, instance), err.Error())
			return
		}

		// read back sql user
		r.SQLUserRefresh(ctx, &new, &resp.State, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &new)...)
}

func (r *SQLUserFWResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data SQLUserModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project := data.Project.ValueString()
	instance := data.Instance.ValueString()
	name := data.Name.ValueString()
	host := data.Host.ValueString()

	deleteTimeout, diags := data.Timeouts.Delete(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	transport_tpg.MutexStore.Lock(instanceMutexKey(project, instance))
	defer transport_tpg.MutexStore.Unlock(instanceMutexKey(project, instance))
	var op *sqladmin.Operation
	var err error
	deleteFunc := func() error {
		op, err = r.client.Users.Delete(project, instance).Host(host).Name(name).Do()
		return err
	}
	err = transport_tpg.Retry(transport_tpg.RetryOptions{
		RetryFunc: deleteFunc,
		Timeout:   deleteTimeout,
	})

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("failed to delete"+
			"user %s in instance %s", name, instance), err.Error())
		return
	}

	err = SqlAdminOperationWaitTime(r.providerConfig, op, project, "Delete User", r.client.UserAgent, deleteTimeout)

	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error, failure waiting to delete "+
			"user %s", name), err.Error())
		return
	}
}

func (r *SQLUserFWResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, "/")

	// TODO recreate all import cases
	if len(idParts) != 4 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: project/instance/host/name. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("instance"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("host"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[3])...)
}

func (r *SQLUserFWResource) SQLUserRefresh(ctx context.Context, data *SQLUserModel, state *tfsdk.State, diag *diag.Diagnostics) {
	userReadResp, err := r.client.Users.Get(data.Project.ValueString(), data.Instance.ValueString(), data.Name.ValueString()).Host(data.Host.ValueString()).Do()
	if err != nil {
		// Treat HTTP 404 Not Found status as a signal to recreate resource
		// and return early
		if userReadResp != nil && transport.IsGoogleApiErrorWithCode(err, userReadResp.HTTPStatusCode) {
			tflog.Trace(ctx, "sql user resource not found, removing from state")
			state.RemoveResource(ctx)
			return
		}
		diag.AddError(fmt.Sprintf("Error, failure waiting to read "+
			"user %s", data.Name.ValueString()), err.Error())
		return
	}

	id := fmt.Sprintf("projects/%s/global/networks/%s", userReadResp.Project, userReadResp.Name)
	data.Id = types.StringValue(id)
	data.Project = types.StringValue(userReadResp.Project)
	data.Instance = types.StringValue(userReadResp.Instance)
	if userReadResp.Host != "" {
		data.Host = types.StringValue(userReadResp.Host)
	}
	if userReadResp.Type != "" {
		data.Type = types.StringValue(userReadResp.Type)
	}
}

// Plan Modifiers
func SQLUserNameIAMPlanModifier() planmodifier.String {
	return &sqlUserNameIAMPlanModifier{}
}

type sqlUserNameIAMPlanModifier struct {
}

func (d *sqlUserNameIAMPlanModifier) Description(ctx context.Context) string {
	return "Suppresses name diffs for IAM user types."
}
func (d *sqlUserNameIAMPlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

// Plan modifier to emulate the SDK diffSuppressIamUserName
func (d *sqlUserNameIAMPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Retrieve relevant fields
	var oldName types.String
	diags := req.State.GetAttribute(ctx, path.Root("name"), &oldName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var newName types.String
	diags = req.Plan.GetAttribute(ctx, path.Root("name"), &newName)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var userType types.String
	diags = req.Plan.GetAttribute(ctx, path.Root("type"), &userType)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Old diff suppress logic
	strippedNewName := strings.Split(newName.ValueString(), "@")[0]

	if oldName.ValueString() == strippedNewName && strings.Contains(userType.ValueString(), "IAM") {
		// Suppress the diff by setting the planned value to the old value
		resp.PlanValue = oldName
	}
}
