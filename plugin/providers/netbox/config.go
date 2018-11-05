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
	log.Debugf("config.go Client() AppID: %s", c.AppID)
	log.Debugf("config.go Client() Endpoint: %s", c.Endpoint)
	cfg := Config{
		AppID:    c.AppID,
		Endpoint: c.Endpoint,
	}
	log.Debugf("Initializing Netbox controllers asdf asdfasfasdfasd")
	// sess := session.NewSession(cfg)
	// Create the Client
	// cli := api.NewNetboxWithAPIKey(cfg.Endpoint, cfg.AppID)

	parsedUri, uriParseError := url.ParseRequestURI(cfg.Endpoint)

	if uriParseError != nil {
		log.Debugf("Failed to parse URI %v into URL: %v", c.Endpoint, uriParseError)
		return nil, uriParseError
	}

	parsedScheme := strings.ToLower(parsedUri.Scheme)

	if parsedScheme == "" {
		parsedScheme = "http"
	}

	desiredRuntimeClientSchemes := []string{parsedScheme}

	log.Debugf("Initializing new openapi runtime client, host = %v, desired schemes = %v", parsedUri.Host, desiredRuntimeClientSchemes)
	runtimeClient := openapi_runtimeclient.New(parsedUri.Host, client.DefaultBasePath, desiredRuntimeClientSchemes)
	runtimeClient.DefaultAuthentication = openapi_runtimeclient.APIKeyAuth("Authorization", "header", fmt.Sprintf("Token %v", cfg.AppID))
	netboxClient := client.New(runtimeClient, strfmt.Default)

	// Validate that our connection is okay
	if err := c.ValidateConnection(netboxClient); err != nil {
		log.Debugf("config.go Client() Erro")
		return nil, err
	}
	cs := ProviderNetboxClient{
		client:        netboxClient,
		configuration: cfg,
	}
	return &cs, nil
}

// ValidateConnection ensures that we can connect to Netbox early, so that we
// do not fail in the middle of a TF run if it can be prevented.
func (c *Config) ValidateConnection(sc *client.NetBox) error {
	log.Debugf("config.go ValidateConnection() validando ")
	rs, err := sc.Dcim.DcimRacksList(nil, nil)
	log.Println(rs)
	return err
}
