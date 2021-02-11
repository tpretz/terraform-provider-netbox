package netbox

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"net/url"

	"github.com/go-openapi/strfmt"

	"github.com/tpretz/go-netbox/netbox/client"

	openapi_runtimeclient "github.com/go-openapi/runtime/client"
)

// Config provides the configuration for the Netbox provider.
type Config struct {
	// The application ID required for API requests. This needs to be created
	// in the NETBOX console (Admin->Users->Tokens). It can also be supplied via the NETBOX_APP_ID
	// environment variable.
	AppID string

	// The API endpoint. This defaults to http://localhost/api, and can also be
	// supplied via the NETBOX_ENDPOINT_ADDR environment variable.
	Endpoint string
}

type ProviderNetboxClient struct {
	client        *client.NetBox
	configuration Config
}

// Client does the heavy lifting of establishing a base Open API client to Netbox.
func (c *Config) Client() (interface{}, error) {
	cfg := Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
	}

	log.WithFields(
		log.Fields{
			"uri": cfg.Endpoint,
		},
	).Debug("Initializing Netbox client")

	parsedURI, uriParseError := url.ParseRequestURI(cfg.Endpoint)

	if uriParseError != nil {
		log.WithFields(
			log.Fields{
				"uri":   cfg.Endpoint,
				"error": uriParseError,
			},
		).Error("Failed to parse URI")

		return nil, uriParseError
	}

	parsedScheme := strings.ToLower(parsedURI.Scheme)

	if parsedScheme == "" {
		parsedScheme = "http"
	}

	desiredRuntimeClientSchemes := []string{parsedScheme}

	log.WithFields(
		log.Fields{
			"host":    parsedURI.Host,
			"schemes": desiredRuntimeClientSchemes,
		},
	).Debug("Initializing open API runtime client")

	runtimeClient := openapi_runtimeclient.New(parsedURI.Host, client.DefaultBasePath, desiredRuntimeClientSchemes)

	runtimeClient.DefaultAuthentication = openapi_runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.AppID))
	runtimeClient.SetLogger(log.StandardLogger())

	netboxClient := client.New(runtimeClient, strfmt.Default)

	terraformNetboxClient := ProviderNetboxClient{
		client:        netboxClient,
		configuration: cfg,
	}

	return &terraformNetboxClient, nil
}
