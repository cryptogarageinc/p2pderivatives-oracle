package cryptocompare

import (
	"fmt"
	"net/http"
	"p2pderivatives-oracle/internal/datafeed"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/go-resty/resty/v2"
)

const (
	priceRoute           = "/price"
	priceHistoricalRoute = "/pricehistorical"
)

// NewClient returns a new CryptoCompare Client (not initialized)
func NewClient(config *Config) *Client {
	return &Client{
		config:      config,
		initialized: false,
	}
}

type apiPriceResponse map[string]float64
type apiPastPriceResponse map[string]apiPriceResponse

// Client represents a CryptoCompare REST client
type Client struct {
	datafeed.DataFeed
	config      *Config
	httpClient  *resty.Client
	initialized bool
}

// Initialize initializes the http client
func (c *Client) Initialize() {
	c.httpClient = resty.New()
	c.httpClient.SetHostURL(c.config.APIBaseURL)
	c.httpClient.SetHeader("Accept", "application/json")
	c.httpClient.SetHeader("authorization", "Apikey "+c.config.APIKey)
	c.initialized = true
}

// IsInitialized returns true if the Client has been initialized
func (c *Client) IsInitialized() bool {
	return c.initialized
}

// FindCurrentAssetPrice sends a GET request to the CryptoCompare API to retrieve the current price of an asset
func (c *Client) FindCurrentAssetPrice(assetID string, currency string) (*float64, error) {
	route := fmt.Sprintf(priceRoute+"?fsym=%s&tsyms=%s", assetID, currency)
	resp, err := c.getAssetPrice(route, apiPriceResponse{})
	if err != nil {
		return nil, err
	}

	res := *(resp.Result().(*apiPriceResponse))
	val, ok := res[strings.ToUpper(currency)]

	// it should not happened if the request was well formed
	if !ok {
		return nil, errors.Errorf("error currency %s not found in cryptocompare response", currency)
	}

	return &val, nil
}

// FindPastAssetPrice sends a GET request to the CryptoCompare API to retrieve a past price of an asset
func (c *Client) FindPastAssetPrice(assetID string, currency string, date time.Time) (*float64, error) {
	if time.Now().Before(date) {
		return nil, errors.New("date should be before now")
	}
	route := fmt.Sprintf(priceHistoricalRoute+"?fsym=%s&tsyms=%s&ts=%d", assetID, currency, date.Unix())
	resp, err := c.getAssetPrice(route, apiPastPriceResponse{})
	if err != nil {
		return nil, err
	}

	res := *(resp.Result().(*apiPastPriceResponse))

	// check asset ID Object
	assetObj, ok := res[strings.ToUpper(assetID)]
	// it should not happened if the request was well formed
	if !ok {
		return nil, errors.Errorf("error asset %s not found in cryptocompare response", assetID)
	}

	// check currency
	val, ok := assetObj[strings.ToUpper(currency)]
	// it should not happened if the request was well formed
	if !ok {
		return nil, errors.Errorf("error currency %s not found in cryptocompare response", currency)
	}

	return &val, nil
}

func (c *Client) getAssetPrice(route string, resultType interface{}) (*resty.Response, error) {
	if !c.IsInitialized() {
		return nil, errors.New("crypto compare client is not initialized")
	}
	req := c.httpClient.R()
	req.SetResult(resultType)
	resp, err := req.Get(route)
	if err != nil {
		return nil, errors.WithMessagef(err, "error while sending a request to cryptocompare api %v", resp.String())
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New("error in the crypto compare response")
	}

	return resp, nil
}
