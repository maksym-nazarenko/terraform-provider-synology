package provider

import (
	"context"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/maksym-nazarenko/terraform-provider-synology/internal/provider/filestation"
	client "github.com/maksym-nazarenko/terraform-provider-synology/synology-go"
)

const (
	SYNOLOGY_HOST_ENV_VAR            = "SYNOLOGY_HOST"
	SYNOLOGY_USER_ENV_VAR            = "SYNOLOGY_USER"
	SYNOLOGY_PASSWORD_ENV_VAR        = "SYNOLOGY_PASSWORD"
	SYNOLOGY_SKIP_CERT_CHECK_ENV_VAR = "SYNOLOGY_SKIP_CERT_CHECK"
)

// Ensure SynologyProvider satisfies various provider interfaces.
var _ provider.Provider = &SynologyProvider{}

// SynologyProvider defines the provider implementation.
type SynologyProvider struct{}

// providerModel describes the provider data model.
type providerModel struct {
	Host          types.String `tfsdk:"host"`
	User          types.String `tfsdk:"user"`
	Password      types.String `tfsdk:"password"`
	SkipCertCheck types.Bool   `tfsdk:"skip_cert_check"`
}

func (p *SynologyProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "synology"
}

func (p *SynologyProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "Remote Synology station host in form of 'host:port'.",
				Optional:    true,
			},
			"user": schema.StringAttribute{
				Description: "User to connect to Synology station with.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password to use when connecting to Synology station.",
				Optional:    true,
				Sensitive:   true,
			},
			"skip_cert_check": schema.BoolAttribute{
				Description: "Whether to skip SSL certificate checks.",
				Optional:    true,
			},
		},
	}
}

func (p *SynologyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data providerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	host := data.Host.ValueString()
	if v := os.Getenv(SYNOLOGY_HOST_ENV_VAR); v != "" {
		host = v
	}

	user := data.User.ValueString()
	if v := os.Getenv(SYNOLOGY_USER_ENV_VAR); v != "" {
		user = v
	}
	password := data.Password.ValueString()
	if v := os.Getenv(SYNOLOGY_PASSWORD_ENV_VAR); v != "" {
		password = v
	}

	skipCertificateCheck := data.SkipCertCheck.ValueBool()
	if vString := os.Getenv(SYNOLOGY_SKIP_CERT_CHECK_ENV_VAR); vString != "" {
		if v, err := strconv.ParseBool(vString); err == nil {
			skipCertificateCheck = v
		}
	}

	if host == "" {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			path.Root("host"),
			"invalid provider configuration",
			"host information is not provided"))
	}
	if user == "" {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			path.Root("user"),
			"invalid provider configuration",
			"user information is not provided"))
	}
	if password == "" {
		resp.Diagnostics.Append(diag.NewAttributeErrorDiagnostic(
			path.Root("password"),
			"invalid provider configuration",
			"password information is not provided"))
	}
	// Example client configuration for data sources and resources
	client, err := client.New(host, skipCertificateCheck)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("synology client creation failed", err.Error()))
	}
	if err := client.Login(user, password, "webui"); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("login to Synology station failed", err.Error()))
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *SynologyProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *SynologyProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		filestation.NewInfoDataSource,
	}
}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &SynologyProvider{}
	}
}
