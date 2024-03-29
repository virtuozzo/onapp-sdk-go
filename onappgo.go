package onappgo

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	sdk "github.com/OnApp/onapp-sdk-go/version"

	"github.com/google/go-querystring/query"
)

var userAgent = "onappgo/" + sdk.String()

const (
	apiFormat          = ".json"
	defaultBaseURL     = "https://69.168.239.52"
	headerPage         = "X-Page"
	headerPerPage      = "X-Limit"
	headerRequestID    = "X-Request-Id"
	headerTotal        = "X-Total"
	mediaType          = "application/json"
	searchTransactions = 100
)

// Client manages communication with OnApp API.
type Client struct {
	// HTTP client used to communicate with the OnApp SDK API.
	client    *http.Client
	transport *http.Transport

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string

	apiUser     string
	apiPassword string

	// Services used for communicating with the API
	AccessControls            AccessControlsService
	BackupResources           BackupResourcesService
	BackupResourceZones       BackupResourceZonesService
	Backups                   BackupsService
	BackupServerGroups        BackupServerGroupsService
	BackupServerJoins         BackupServerJoinsService
	BackupServers             BackupServersService
	Buckets                   BucketsService
	CloudbootComputeResources CloudbootComputeResourcesService
	CloudbootIPAddresses      CloudbootIPAddressesService
	Configurations            ConfigurationsService
	DataStoreGroups           DataStoreGroupsService
	DataStoreJoins            DataStoreJoinsService
	DataStores                DataStoresService
	Disks                     DisksService
	Engines                   EnginesService
	FirewallRules             FirewallRulesService
	HypervisorGroups          HypervisorGroupsService
	Hypervisors               HypervisorsService
	HypervisorZones           HypervisorZonesService
	ImageTemplateGroups       ImageTemplateGroupsService
	ImageTemplates            ImageTemplatesService
	InstancePackages          InstancePackagesService
	IntegratedDataStores      IntegratedDataStoresService
	IPAddresses               IPAddressesService
	IPNets                    IPNetsService
	IPRanges                  IPRangesService
	Licenses                  LicensesService
	LocationGroups            LocationGroupsService
	NetworkGroups             NetworkGroupsService
	NetworkInterfaces         NetworkInterfacesService
	NetworkJoins              NetworkJoinsService
	Networks                  NetworksService
	RateCards                 RateCardsService
	RecipeGroups              RecipeGroupsService
	Recipes                   RecipesService
	RecipeSteps               RecipeStepsService
	RecipeJoins               RecipeJoinsService
	RemoteTemplates           RemoteTemplatesService
	Resolvers                 ResolversService
	Roles                     RolesService
	SoftwareLicenses          SoftwareLicensesService
	SSHKeys                   SSHKeysService
	Transactions              TransactionsService
	UserGroups                UserGroupsService
	Users                     UsersService
	UserWhiteLists            UserWhiteListsService
	VirtualMachineActions     VirtualMachineActionsService
	VirtualMachines           VirtualMachinesService

	// Optional function called after every successful request made to the OnApp APIs
	onRequestCompleted RequestCompletionCallback
}

// RequestCompletionCallback defines the type of the request callback function
type RequestCompletionCallback func(*http.Request, *http.Response)

// ListOptions specifies the optional parameters to various List methods that
// support pagination.
type ListOptions struct {
	// For paginated result sets, page of results to retrieve.
	Page int `url:"page,omitempty"`

	// For paginated result sets, the number of results to include per page.
	PerPage int `url:"per_page,omitempty"`
}

// Response is a OnApp response. This wraps the standard http.Response returned from OnApp.
type Response struct {
	*http.Response

	// Links that were returned with the response. These are parsed from
	// request body and not the header.
	Links *Links
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error messages
	Errors map[string][]string `json:"errors,omitempty"`
}

func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)

	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	origURL, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	origValues := origURL.Query()

	newValues, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	for k, v := range newValues {
		origValues[k] = v
	}

	origURL.RawQuery = origValues.Encode()
	return origURL.String(), nil
}

// NewClient returns a new OnApp API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		c.transport = &http.Transport{
			DialContext:           t.DialContext,
			ExpectContinueTimeout: t.ExpectContinueTimeout,
			IdleConnTimeout:       t.IdleConnTimeout,
			MaxIdleConns:          t.MaxIdleConns,
			Proxy:                 t.Proxy,
			TLSHandshakeTimeout:   t.TLSHandshakeTimeout,
		}
	} else {
		c.transport = new(http.Transport)
	}

	c.AccessControls = &AccessControlsServiceOp{client: c}
	c.BackupResources = &BackupResourcesServiceOp{client: c}
	c.BackupResourceZones = &BackupResourceZonesServiceOp{client: c}
	c.Backups = &BackupsServiceOp{client: c}
	c.BackupServerGroups = &BackupServerGroupsServiceOp{client: c}
	c.BackupServerJoins = &BackupServerJoinsServiceOp{client: c}
	c.BackupServers = &BackupServersServiceOp{client: c}
	c.Buckets = &BucketsServiceOp{client: c}
	c.CloudbootComputeResources = &CloudbootComputeResourcesServiceOp{client: c}
	c.CloudbootIPAddresses = &CloudbootIPAddressesServiceOp{client: c}
	c.Configurations = &ConfigurationsServiceOp{client: c}
	c.DataStoreGroups = &DataStoreGroupsServiceOp{client: c}
	c.DataStoreJoins = &DataStoreJoinsServiceOp{client: c}
	c.DataStores = &DataStoresServiceOp{client: c}
	c.Disks = &DisksServiceOp{client: c}
	c.Engines = &EnginesServiceOp{client: c}
	c.FirewallRules = &FirewallRulesServiceOp{client: c}
	c.HypervisorGroups = &HypervisorGroupsServiceOp{client: c}
	c.Hypervisors = &HypervisorsServiceOp{client: c}
	c.HypervisorZones = &HypervisorZonesServiceOp{client: c}
	c.ImageTemplateGroups = &ImageTemplateGroupsServiceOp{client: c}
	c.ImageTemplates = &ImageTemplatesServiceOp{client: c}
	c.InstancePackages = &InstancePackagesServiceOp{client: c}
	c.IntegratedDataStores = &IntegratedDataStoresServiceOp{client: c}
	c.IPAddresses = &IPAddressesServiceOp{client: c}
	c.IPNets = &IPNetsServiceOp{client: c}
	c.IPRanges = &IPRangesServiceOp{client: c}
	c.Licenses = &LicensesServiceOp{client: c}
	c.LocationGroups = &LocationGroupsServiceOp{client: c}
	c.NetworkGroups = &NetworkGroupsServiceOp{client: c}
	c.NetworkInterfaces = &NetworkInterfacesServiceOp{client: c}
	c.NetworkJoins = &NetworkJoinsServiceOp{client: c}
	c.Networks = &NetworksServiceOp{client: c}
	c.RateCards = &RateCardsServiceOp{client: c}
	c.RecipeGroups = &RecipeGroupsServiceOp{client: c}
	c.Recipes = &RecipesServiceOp{client: c}
	c.RecipeSteps = &RecipeStepsServiceOp{client: c}
	c.RecipeJoins = &RecipeJoinsServiceOp{client: c}
	c.RemoteTemplates = &RemoteTemplatesServiceOp{client: c}
	c.Resolvers = &ResolversServiceOp{client: c}
	c.Roles = &RolesServiceOp{client: c}
	c.SoftwareLicenses = &SoftwareLicensesServiceOp{client: c}
	c.SSHKeys = &SSHKeysServiceOp{client: c}
	c.Transactions = &TransactionsServiceOp{client: c}
	c.UserGroups = &UserGroupsServiceOp{client: c}
	c.Users = &UsersServiceOp{client: c}
	c.UserWhiteLists = &UserWhiteListsServiceOp{client: c}
	c.VirtualMachineActions = &VirtualMachineActionsServiceOp{client: c}
	c.VirtualMachines = &VirtualMachinesServiceOp{client: c}

	return c
}

func (c *Client) certificate() *tls.Certificate {
	certs := c.transport.TLSClientConfig.Certificates
	if len(certs) == 0 {
		return nil
	}
	return &certs[0]
}

func (c *Client) setCertificate(cert tls.Certificate) {
	t := c.client.Transport.(*http.Transport)

	// Extension or HoK certificate
	t.TLSClientConfig.Certificates = []tls.Certificate{cert}
}

// ClientOpt are options for New.
type ClientOpt func(*Client) error

// New returns a new OnApp API client instance.
func New(httpClient *http.Client, opts ...ClientOpt) (*Client, error) {
	c := NewClient(httpClient)
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// SetBaseURL is a client option for setting the base URL.
func SetBaseURL(bu string) ClientOpt {
	return func(c *Client) error {
		u, err := url.Parse(bu)
		if err != nil {
			return err
		}

		c.BaseURL = u
		return nil
	}
}

// SetUserAgent is a client option for setting the user agent.
func SetUserAgent(ua string) ClientOpt {
	return func(c *Client) error {
		c.UserAgent = fmt.Sprintf("%s %s", ua, c.UserAgent)
		return nil
	}
}

// SetBasicAuth is a client option for setting the user and password for API call.
func SetBasicAuth(user, password string) ClientOpt {
	return func(c *Client) error {
		c.apiUser = user
		c.apiPassword = password
		return nil
	}
}

// SetAllowUnverifiedSSL is a client option for setting allowUnverifiedSSL.
func SetAllowUnverifiedSSL(isv bool) ClientOpt {
	return func(c *Client) error {
		c.transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: isv}

		// Don't bother setting DialTLS if InsecureSkipVerify=true
		if !isv {
			c.transport.DialTLSContext = nil
		}

		c.client.Transport = c.transport

		if cert := c.certificate(); cert != nil {
			c.setCertificate(*cert)
		}
		return nil
	}
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLS should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// maybe resolve problem with EOF error during POST request
	req.Close = true

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)

	req.SetBasicAuth(c.apiUser, c.apiPassword)

	return req, nil
}

// OnRequestCompleted sets the OnApp API request completion callback
func (c *Client) OnRequestCompleted(rc RequestCompletionCallback) {
	c.onRequestCompleted = rc
}

// newResponse creates a new Response for the provided http.Response
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	if response != nil {
		response.populateLinks()
	}

	return response
}

func (r *Response) populateLinks() {
	limit := r.Header.Get(headerPerPage)
	page := r.Header.Get(headerPage)
	total := r.Header.Get(headerTotal)

	if limit == "" && page == "" && total == "" {
		r.Links = nil
		return
	}

	r.Links = &Links{}
	r.Links.PerPage, _ = strconv.Atoi(limit)
	r.Links.CurPage, _ = strconv.Atoi(page)
	r.Links.Total, _ = strconv.Atoi(total)
	r.Links.NumPages = int(r.Links.Total/r.Links.PerPage) + 1
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*Response, error) {
	resp := &http.Response{}
	err := *new(error)

	resp, err = DoRequestWithClient(ctx, c.client, req)
	if err != nil {
		return nil, err
	}

	if c.onRequestCompleted != nil {
		c.onRequestCompleted(req, resp)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return nil, err
			}
		}
	}

	return response, err
}

// DoRequest submits an HTTP request.
func DoRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	return DoRequestWithClient(ctx, http.DefaultClient, req)
}

// DoRequestWithClient submits an HTTP request using the specified client.
func DoRequestWithClient(ctx context.Context, client *http.Client, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	return client.Do(req)
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.String())
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errorResponse)
	}

	return errorResponse
}

// Version return OnApp endpoint version
func (c *Client) Version() (string, *Response, error) {
	path := fmt.Sprintf("version%s", apiFormat)

	req, err := c.NewRequest(context.TODO(), http.MethodGet, path, nil)
	if err != nil {
		return "", nil, err
	}

	var res map[string]string
	resp, err := c.Do(context.TODO(), req, &res)
	if err != nil {
		return "", resp, err
	}

	return res["version"], resp, err
}

// String - convert ErrorResponse to the string
func (r *ErrorResponse) String() string {
	var str string

	for name, message := range r.Errors {
		str = str + fmt.Sprintf("\n%s %s", name, message)
	}

	return str
}

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}

// Int is a helper routine that allocates a new int32 value
// to store v and returns a pointer to it, but unlike Int32
// its argument value is an int.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// StreamToString converts a reader to a string
func StreamToString(stream io.Reader) string {
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(stream)
	return buf.String()
}

// StringInSlice -
func StringInSlice(valid []string, k string, ignoreCase bool) bool {
	for _, str := range valid {
		if k == str || (ignoreCase && strings.EqualFold(strings.ToLower(k), strings.ToLower(str))) {
			return true
		}
	}

	return false
}

// ServerTypesRestrictions -
func ServerTypesRestrictions(resourceType string, serverType string, ignoreCase bool) (string, bool) {
	slice := SERVER_TYPES_RESTRICTIONS[serverType]

	return strings.Join(slice, ", "), StringInSlice(slice, resourceType, ignoreCase)
}
