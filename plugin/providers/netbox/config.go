package netbox

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"net/url"

	"github.com/go-openapi/strfmt"

	"github.com/digitalocean/go-netbox/netbox/client"

	openapi_runtimeclient "github.com/go-openapi/runtime/client"
)

// Config provides the configuration for the NETBOX providerr.
type Config struct {
	// The application ID required for API requests. This needs to be created
	// in the NETBOX console. It can also be supplied via the NETBOX_APP_ID
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

// ProviderNetboxClient is a structure that contains the client connections
// necessary to interface with the Go-Netbox API
//type ProviderNetboxClient struct {
//		client *Client
//}
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
	runtimeClient.SetLogger(&log.Logger{})
	runtimeClient.SetDebug(true)
	netboxClient := client.New(runtimeClient, strfmt.Default)

	// Validate that our connection is okay
	if err := c.ValidateConnection(netboxClient); err != nil {
		log.WithFields(
			log.Fields{
				"uri":   cfg.Endpoint,
				"error": err,
			},
		).Error("Failed to validate connection")

		return nil, err
	}

	terraformNetboxClient := ProviderNetboxClient{
		client:        netboxClient,
		configuration: cfg,
	}

	return &terraformNetboxClient, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *client.NetBox) error {
	log.Debug("Validating Netbox connection")

	_, err := sc.Dcim.DcimRacksList(nil, nil)

	if err != nil {
		log.Error("Failed to validate connection to Netbox")
	}

	log.Debug("Netbox connection validated")

	return nil
}
