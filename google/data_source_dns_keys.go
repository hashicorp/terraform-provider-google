package google

import (
	"context"
	"fmt"
	"google.golang.org/api/dns/v1"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ datasource.DataSource              = &GoogleDnsKeysDataSource{}
	_ datasource.DataSourceWithConfigure = &GoogleDnsKeysDataSource{}
)

func NewGoogleDnsKeysDataSource() datasource.DataSource {
	return &GoogleDnsKeysDataSource{}
}

// GoogleDnsKeysDataSource defines the data source implementation
type GoogleDnsKeysDataSource struct {
	client  *dns.Service
	project types.String
}

type GoogleDnsKeysModel struct {
	Id              types.String `tfsdk:"id"`
	ManagedZone     types.String `tfsdk:"managed_zone"`
	Project         types.String `tfsdk:"project"`
	KeySigningKeys  types.List   `tfsdk:"key_signing_keys"`
	ZoneSigningKeys types.List   `tfsdk:"zone_signing_keys"`
}

type GoogleZoneSigningKey struct {
	Algorithm    types.String `tfsdk:"algorithm"`
	CreationTime types.String `tfsdk:"creation_time"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	KeyLength    types.Int64  `tfsdk:"key_length"`
	KeyTag       types.Int64  `tfsdk:"key_tag"`
	PublicKey    types.String `tfsdk:"public_key"`
	Digests      types.List   `tfsdk:"digests"`
}

type GoogleKeySigningKey struct {
	Algorithm    types.String `tfsdk:"algorithm"`
	CreationTime types.String `tfsdk:"creation_time"`
	Description  types.String `tfsdk:"description"`
	Id           types.String `tfsdk:"id"`
	IsActive     types.Bool   `tfsdk:"is_active"`
	KeyLength    types.Int64  `tfsdk:"key_length"`
	KeyTag       types.Int64  `tfsdk:"key_tag"`
	PublicKey    types.String `tfsdk:"public_key"`
	Digests      types.List   `tfsdk:"digests"`

	DSRecord types.String `tfsdk:"ds_record"`
}

type GoogleZoneSigningKeyDigest struct {
	Digest types.String `tfsdk:"digest"`
	Type   types.String `tfsdk:"type"`
}

var (
	digestAttrTypes = map[string]attr.Type{
		"digest": types.StringType,
		"type":   types.StringType,
	}
)

func (d *GoogleDnsKeysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dns_keys"
}

func (d *GoogleDnsKeysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Get the DNSKEY and DS records of DNSSEC-signed managed zones",

		Attributes: map[string]schema.Attribute{
			"managed_zone": schema.StringAttribute{
				Description:         "The Name of the zone.",
				MarkdownDescription: "The Name of the zone.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				Description:         "The ID of the project for the Google Cloud.",
				MarkdownDescription: "The ID of the project for the Google Cloud.",
				Optional:            true,
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "DNS keys identifier",
				MarkdownDescription: "DNS keys identifier",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"zone_signing_keys": schema.ListNestedBlock{
				Description:         "A list of Zone-signing key (ZSK) records.",
				MarkdownDescription: "A list of Zone-signing key (ZSK) records.",
				NestedObject:        dnsKeyObject(),
			},
			"key_signing_keys": schema.ListNestedBlock{
				Description:         "A list of Key-signing key (KSK) records.",
				MarkdownDescription: "A list of Key-signing key (KSK) records.",
				NestedObject:        kskObject(),
			},
		},
	}
}

func (d *GoogleDnsKeysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*frameworkProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *frameworkProvider, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = p.NewDnsClient(p.userAgent, &resp.Diagnostics)
	d.project = p.project
}

func (d *GoogleDnsKeysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GoogleDnsKeysModel
	var metaData *ProviderMetaModel
	var diags diag.Diagnostics

	// Read Provider meta into the meta model
	resp.Diagnostics.Append(req.ProviderMeta.Get(ctx, &metaData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.client.UserAgent = generateFrameworkUserAgentString(metaData, d.client.UserAgent)

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fv := parseProjectFieldValueFramework("managedZones", data.ManagedZone.ValueString(), "project", data.Project, d.project, false, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Project = types.StringValue(fv.Project)
	data.ManagedZone = types.StringValue(fv.Name)

	data.Id = types.StringValue(fmt.Sprintf("projects/%s/managedZones/%s", data.Project.ValueString(), data.ManagedZone.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("fetching DNS keys from managed zone %s", data.ManagedZone.ValueString()))

	clientResp, err := d.client.DnsKeys.List(data.Project.ValueString(), data.ManagedZone.ValueString()).Do()
	if err != nil {
		if !IsGoogleApiErrorWithCode(err, 404) {
			resp.Diagnostics.AddError(fmt.Sprintf("Error when reading or editing dataSourceDnsKeys"), err.Error())
		}
		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	tflog.Trace(ctx, "read dns keys data source")

	zoneSigningKeys, keySigningKeys := flattenSigningKeys(ctx, clientResp.DnsKeys, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	zskObjType := types.ObjectType{}.WithAttributeTypes(getDnsKeyAttrs("zoneSigning"))
	data.ZoneSigningKeys, diags = types.ListValueFrom(ctx, zskObjType, zoneSigningKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	kskObjType := types.ObjectType{}.WithAttributeTypes(getDnsKeyAttrs("keySigning"))
	data.KeySigningKeys, diags = types.ListValueFrom(ctx, kskObjType, keySigningKeys)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// dnsKeyObject is a helper function for the zone_signing_keys schema and
// is also used by key_signing_keys schema (called in kskObject defined below)
func dnsKeyObject() schema.NestedBlockObject {
	return schema.NestedBlockObject{
		Attributes: map[string]schema.Attribute{
			"algorithm": schema.StringAttribute{
				Description: "String mnemonic specifying the DNSSEC algorithm of this key. Immutable after creation time. " +
					"Possible values are `ecdsap256sha256`, `ecdsap384sha384`, `rsasha1`, `rsasha256`, and `rsasha512`.",
				MarkdownDescription: "String mnemonic specifying the DNSSEC algorithm of this key. Immutable after creation time. " +
					"Possible values are `ecdsap256sha256`, `ecdsap384sha384`, `rsasha1`, `rsasha256`, and `rsasha512`.",
				Computed: true,
			},
			"creation_time": schema.StringAttribute{
				Description:         "The time that this resource was created in the control plane. This is in RFC3339 text format.",
				MarkdownDescription: "The time that this resource was created in the control plane. This is in RFC3339 text format.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description:         "A mutable string of at most 1024 characters associated with this resource for the user's convenience.",
				MarkdownDescription: "A mutable string of at most 1024 characters associated with this resource for the user's convenience.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				Description:         "Unique identifier for the resource; defined by the server.",
				MarkdownDescription: "Unique identifier for the resource; defined by the server.",
				Computed:            true,
			},
			"is_active": schema.BoolAttribute{
				Description: "Active keys will be used to sign subsequent changes to the ManagedZone. " +
					"Inactive keys will still be present as DNSKEY Resource Records for the use of resolvers validating existing signatures.",
				MarkdownDescription: "Active keys will be used to sign subsequent changes to the ManagedZone. " +
					"Inactive keys will still be present as DNSKEY Resource Records for the use of resolvers validating existing signatures.",
				Computed: true,
			},
			"key_length": schema.Int64Attribute{
				Description:         "Length of the key in bits. Specified at creation time then immutable.",
				MarkdownDescription: "Length of the key in bits. Specified at creation time then immutable.",
				Computed:            true,
			},
			"key_tag": schema.Int64Attribute{
				Description: "The key tag is a non-cryptographic hash of the a DNSKEY resource record associated with this DnsKey. " +
					"The key tag can be used to identify a DNSKEY more quickly (but it is not a unique identifier). " +
					"In particular, the key tag is used in a parent zone's DS record to point at the DNSKEY in this child ManagedZone. " +
					"The key tag is a number in the range [0, 65535] and the algorithm to calculate it is specified in RFC4034 Appendix B.",
				MarkdownDescription: "The key tag is a non-cryptographic hash of the a DNSKEY resource record associated with this DnsKey. " +
					"The key tag can be used to identify a DNSKEY more quickly (but it is not a unique identifier). " +
					"In particular, the key tag is used in a parent zone's DS record to point at the DNSKEY in this child ManagedZone. " +
					"The key tag is a number in the range [0, 65535] and the algorithm to calculate it is specified in RFC4034 Appendix B.",
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Description:         "Base64 encoded public half of this key.",
				MarkdownDescription: "Base64 encoded public half of this key.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"digests": schema.ListNestedBlock{
				Description:         "A list of cryptographic hashes of the DNSKEY resource record associated with this DnsKey. These digests are needed to construct a DS record that points at this DNS key.",
				MarkdownDescription: "A list of cryptographic hashes of the DNSKEY resource record associated with this DnsKey. These digests are needed to construct a DS record that points at this DNS key.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"digest": schema.StringAttribute{
							Description:         "The base-16 encoded bytes of this digest. Suitable for use in a DS resource record.",
							MarkdownDescription: "The base-16 encoded bytes of this digest. Suitable for use in a DS resource record.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							Description:         "Specifies the algorithm used to calculate this digest. Possible values are `sha1`, `sha256` and `sha384`",
							MarkdownDescription: "Specifies the algorithm used to calculate this digest. Possible values are `sha1`, `sha256` and `sha384`",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

// kskObject is a helper function for the key_signing_keys schema
func kskObject() schema.NestedBlockObject {
	nbo := dnsKeyObject()

	nbo.Attributes["ds_record"] = schema.StringAttribute{
		Description:         "The DS record based on the KSK record.",
		MarkdownDescription: "The DS record based on the KSK record.",
		Computed:            true,
	}

	return nbo
}

func flattenSigningKeys(ctx context.Context, signingKeys []*dns.DnsKey, diags *diag.Diagnostics) ([]types.Object, []types.Object) {
	var zoneSigningKeys []types.Object
	var keySigningKeys []types.Object
	var d diag.Diagnostics

	for _, signingKey := range signingKeys {
		if signingKey != nil {
			var digests []types.Object
			for _, dig := range signingKey.Digests {
				digest := GoogleZoneSigningKeyDigest{
					Digest: types.StringValue(dig.Digest),
					Type:   types.StringValue(dig.Type),
				}
				obj, d := types.ObjectValueFrom(ctx, digestAttrTypes, digest)
				diags.Append(d...)
				if diags.HasError() {
					return zoneSigningKeys, keySigningKeys
				}

				digests = append(digests, obj)
			}

			if signingKey.Type == "keySigning" && len(signingKey.Digests) > 0 {
				ksk := GoogleKeySigningKey{
					Algorithm:    types.StringValue(signingKey.Algorithm),
					CreationTime: types.StringValue(signingKey.CreationTime),
					Description:  types.StringValue(signingKey.Description),
					Id:           types.StringValue(signingKey.Id),
					IsActive:     types.BoolValue(signingKey.IsActive),
					KeyLength:    types.Int64Value(signingKey.KeyLength),
					KeyTag:       types.Int64Value(signingKey.KeyTag),
					PublicKey:    types.StringValue(signingKey.PublicKey),
				}

				objType := types.ObjectType{}.WithAttributeTypes(digestAttrTypes)
				ksk.Digests, d = types.ListValueFrom(ctx, objType, digests)
				diags.Append(d...)
				if diags.HasError() {
					return zoneSigningKeys, keySigningKeys
				}

				dsRecord, err := generateDSRecord(signingKey)
				if err != nil {
					diags.AddError("error generating ds record", err.Error())
					return zoneSigningKeys, keySigningKeys
				}

				ksk.DSRecord = types.StringValue(dsRecord)

				obj, d := types.ObjectValueFrom(ctx, getDnsKeyAttrs(signingKey.Type), ksk)
				diags.Append(d...)
				if diags.HasError() {
					return zoneSigningKeys, keySigningKeys
				}
				keySigningKeys = append(keySigningKeys, obj)
			} else {
				zsk := GoogleZoneSigningKey{
					Algorithm:    types.StringValue(signingKey.Algorithm),
					CreationTime: types.StringValue(signingKey.CreationTime),
					Description:  types.StringValue(signingKey.Description),
					Id:           types.StringValue(signingKey.Id),
					IsActive:     types.BoolValue(signingKey.IsActive),
					KeyLength:    types.Int64Value(signingKey.KeyLength),
					KeyTag:       types.Int64Value(signingKey.KeyTag),
					PublicKey:    types.StringValue(signingKey.PublicKey),
				}

				objType := types.ObjectType{}.WithAttributeTypes(digestAttrTypes)
				zsk.Digests, d = types.ListValueFrom(ctx, objType, digests)
				diags.Append(d...)
				if diags.HasError() {
					return zoneSigningKeys, keySigningKeys
				}

				obj, d := types.ObjectValueFrom(ctx, getDnsKeyAttrs("zoneSigning"), zsk)
				diags.Append(d...)
				if diags.HasError() {
					return zoneSigningKeys, keySigningKeys
				}
				zoneSigningKeys = append(zoneSigningKeys, obj)
			}

		}
	}

	return zoneSigningKeys, keySigningKeys
}

// DNSSEC Algorithm Numbers: https://www.iana.org/assignments/dns-sec-alg-numbers/dns-sec-alg-numbers.xhtml
// The following are algorithms that are supported by Cloud DNS
var dnssecAlgoNums = map[string]int{
	"rsasha1":         5,
	"rsasha256":       8,
	"rsasha512":       10,
	"ecdsap256sha256": 13,
	"ecdsap384sha384": 14,
}

// DS RR Digest Types: https://www.iana.org/assignments/ds-rr-types/ds-rr-types.xhtml
// The following are digests that are supported by Cloud DNS
var dnssecDigestType = map[string]int{
	"sha1":   1,
	"sha256": 2,
	"sha384": 4,
}

// generateDSRecord will generate the ds_record on key signing keys
func generateDSRecord(signingKey *dns.DnsKey) (string, error) {
	algoNum, found := dnssecAlgoNums[signingKey.Algorithm]
	if !found {
		return "", fmt.Errorf("DNSSEC Algorithm number for %s not found", signingKey.Algorithm)
	}

	digestType, found := dnssecDigestType[signingKey.Digests[0].Type]
	if !found {
		return "", fmt.Errorf("DNSSEC Digest type for %s not found", signingKey.Digests[0].Type)
	}

	return fmt.Sprintf("%d %d %d %s",
		signingKey.KeyTag,
		algoNum,
		digestType,
		signingKey.Digests[0].Digest), nil
}

func getDnsKeyAttrs(keyType string) map[string]attr.Type {
	dnsKeyAttrs := map[string]attr.Type{
		"algorithm":     types.StringType,
		"creation_time": types.StringType,
		"description":   types.StringType,
		"id":            types.StringType,
		"is_active":     types.BoolType,
		"key_length":    types.Int64Type,
		"key_tag":       types.Int64Type,
		"public_key":    types.StringType,
		"digests":       types.ListType{}.WithElementType(types.ObjectType{}.WithAttributeTypes(digestAttrTypes)),
	}

	if keyType == "keySigning" {
		dnsKeyAttrs["ds_record"] = types.StringType
	}

	return dnsKeyAttrs
}
