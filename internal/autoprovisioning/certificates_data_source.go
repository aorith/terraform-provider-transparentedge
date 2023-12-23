package autoprovisioning

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/TransparentEdge/terraform-provider-transparentedge/internal/teclient"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &certificatesDataSource{}
	_ datasource.DataSourceWithConfigure = &certificatesDataSource{}
)

// Helper function to simplify the provider implementation.
func NewCertificatesDataSource() datasource.DataSource {
	return &certificatesDataSource{}
}

// data source implementation.
type certificatesDataSource struct {
	client *teclient.Client
}

// Metadata returns the data source type name.
func (d *certificatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_certificates"
}

// Schema defines the schema for the data source.
func (d *certificatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Certificate listing.",
		MarkdownDescription: "Certificate listing.",

		Attributes: map[string]schema.Attribute{
			"certificates": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of all certificates.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed:            true,
							Description:         "ID of the certificate.",
							MarkdownDescription: "ID of the certificate.",
						},
						"company": schema.Int64Attribute{
							Computed:            true,
							Description:         "Company ID that owns this certificate.",
							MarkdownDescription: "Company ID that owns this certificate.",
						},
						"commonname": schema.StringAttribute{
							Computed:            true,
							Description:         "CN (Common Name) of the certificate.",
							MarkdownDescription: "CN (_Common Name_) of the certificate.",
						},
						"domains": schema.StringAttribute{
							Computed:            true,
							Description:         "SAN (Subject Alternative Name) domains included in the certificate, including the Common Name.",
							MarkdownDescription: "SAN (_Subject Alternative Name_) domains included in the certificate, including the Common Name.",
						},
						"expiration": schema.StringAttribute{
							Computed:            true,
							Description:         "Date when the certificate will expire.",
							MarkdownDescription: "Date when the certificate will expire.",
						},
						"autogenerated": schema.BoolAttribute{
							Computed:            true,
							Description:         "True if the certificate was autogenerated and is managed by TransparentEdge CDN (Certificates generated using HTTP or DNS challenge), autogenerated certificates cannot be modified.",
							MarkdownDescription: "`True` if the certificate was autogenerated and is managed by TransparentEdge CDN (Certificates generated using HTTP or DNS challenge), **autogenerated** certificates **cannot** be modified.",
						},
						"dnschallenge": schema.BoolAttribute{
							Computed:            true,
							Description:         "True if the certificate was autogenerated using the DNS challenge, if False and 'autogenerated' is True, it's a certificate generated using HTTP challenge.",
							MarkdownDescription: "`True` if the certificate was autogenerated using the DNS challenge, if `False` and **autogenerated** is `True`, it's a certificate generated using HTTP challenge.",
						},
						"standalone": schema.BoolAttribute{
							Computed:            true,
							Description:         "A standalone certificate will not be merged automatically on the SAN of other existing certificates for the same Company on creation or renewals.",
							MarkdownDescription: "A standalone certificate will not be merged automatically on the SAN of other existing certificates for the same Company on creation or renewals.",
						},
						"publickey": schema.StringAttribute{
							Computed:            true,
							Description:         "Public part of the certificate in PEM format, it's recommended to include the full chain.",
							MarkdownDescription: "Public part of the certificate in PEM format, it's recommended to include the full chain.",
						},
						"privatekey": schema.StringAttribute{
							Computed:            true,
							Description:         "Private key of the certificate in PEM format, it cannot be password protected.",
							MarkdownDescription: "Private key of the certificate in PEM format, it cannot be password protected.",
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *certificatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state Certificates

	certificates, err := d.client.GetCertificates()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Certificates info",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, certificate := range certificates {
		// Sort SAN domains
		sort.Strings(certificate.Domains)
		san_domains := strings.Join(certificate.Domains, ", ")

		// Parse expiration time
		expiration := certificate.Expiration
		exptoint, err := strconv.ParseFloat(certificate.Expiration, 64)
		if err == nil {
			expiration = time.Unix(int64(exptoint), 0).String()
		}

		certificateState := Certificate{
			ID:            types.Int64Value(int64(certificate.ID)),
			Company:       types.Int64Value(int64(certificate.Company)),
			CommonName:    types.StringValue(certificate.CommonName),
			Domains:       types.StringValue(san_domains),
			Expiration:    types.StringValue(expiration),
			Autogenerated: types.BoolValue(certificate.Autogenerated),
			Standalone:    types.BoolValue(certificate.Standalone),
			DNSChallenge:  types.BoolValue(certificate.DNSChallenge),
			PublicKey:     types.StringValue(certificate.PublicKey),
			PrivateKey:    types.StringValue(certificate.PrivateKey),
		}

		state.Certificates = append(state.Certificates, certificateState)
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Configure adds the provider configured client to the data source.
func (d *certificatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*teclient.Client)
}
