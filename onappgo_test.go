package onappgo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	mux                *http.ServeMux
	ctx                = context.TODO()
	client             *Client
	server             *httptest.Server
	allowUnverifiedSSL = true
	email              = "admin@example.com"
	token              = "91108764ba6758e7ef06ce7bb84631483b4d9607"
)

func setup(t *testing.T) {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = nil
	client, _ = New(nil, SetAllowUnverifiedSSL(allowUnverifiedSSL), SetBasicAuth(email, token))

	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v", r.Method, expected)
	}
}

type values map[string]string

func testFormValues(t *testing.T, r *http.Request, values values) {
	expected := url.Values{}
	for k, v := range values {
		expected.Add(k, v)
	}

	err := r.ParseForm()
	if err != nil {
		t.Fatalf("parseForm(): %v", err)
	}

	if !reflect.DeepEqual(expected, r.Form) {
		t.Errorf("Request parameters = %v, expected %v", r.Form, expected)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func testClientServices(t *testing.T, c *Client) {
	services := []string{
		"Buckets",
		"AccessControls",
		"RateCards",
		"Transactions",
		"InstancePackages",
		"VirtualMachines",
		"VirtualMachineActions",
		"Hypervisors",
		"HypervisorGroups",
		"DataStores",
		"DataStoreGroups",
		"RemoteTemplates",
		"ImageTemplates",
		"ImageTemplateGroups",
		"Disks",
		"Networks",
		"NetworkInterfaces",
		"NetworkJoins",
		"NetworkGroups",
		"IPNets",
		"IPRanges",
		"BackupServers",
		"BackupServerGroups",
		"BackupResources",
		"BackupResourceZones",
		"Users",
		"HypervisorZones",
		"Licenses",
		"Roles",
		"Configurations",
		"UserGroups",
		"FirewallRules",
	}

	cp := reflect.ValueOf(c)
	cv := reflect.Indirect(cp)

	for _, s := range services {
		if cv.FieldByName(s).IsNil() {
			t.Errorf("c.%s shouldn't be nil", s)
		}
	}
}

func testClientDefaultBaseURL(t *testing.T, c *Client) {
	if c.BaseURL == nil || c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, defaultBaseURL)
	}
}

func testClientDefaultUserAgent(t *testing.T, c *Client) {
	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}
}

func testClientDefaults(t *testing.T, c *Client) {
	testClientDefaultBaseURL(t, c)
	testClientDefaultUserAgent(t, c)
	testClientServices(t, c)
}

func TestNewClient(t *testing.T) {
	c, _ := New(nil, SetAllowUnverifiedSSL(allowUnverifiedSSL), SetBasicAuth(email, token))

	testClientDefaults(t, c)
}

func TestNew(t *testing.T) {
	c, err := New(nil)

	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	testClientDefaults(t, c)
}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"/foo"
	inBody, outBody := &VirtualMachineCreateRequest{Hostname: "l"},
		`{"acceleration_allowed":false,"hostname":"l",`+
			`"required_ip_address_assignment":false,"required_virtual_machine_build":false,`+
			`"required_virtual_machine_startup":false,"virsh_console":false}`+"\n"
	req, _ := c.NewRequest(ctx, http.MethodGet, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestNewRequest_badURL(t *testing.T) {
	c := NewClient(nil)
	_, err := c.NewRequest(ctx, http.MethodGet, ":", nil)
	testURLParseError(t, err)
}

func TestNewRequest_withCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	req, _ := c.NewRequest(ctx, http.MethodGet, "/foo", nil)

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := req.Header.Get("User-Agent"); got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

func TestDo(t *testing.T) {
	setup(t)
	if client == nil {
		t.Fatalf("client == nil")
	}
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}

	expected := &foo{"a"}
	if !reflect.DeepEqual(body, expected) {
		t.Errorf("Response body = %v, expected %v", body, expected)
	}
}

func TestDo_httpError(t *testing.T) {
	setup(t)
	if client == nil {
		t.Fatalf("client == nil")
	}
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do()
// function.
func TestDo_redirectLoop(t *testing.T) {
	setup(t)
	if client == nil {
		t.Fatalf("client == nil")
	}
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	_, err := client.Do(context.Background(), req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestCheckResponse(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{"message":"m",
			"errors": [{"resource": "r", "field": "f", "code": "c"}]}`)),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Fatalf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

// ensure that we properly handle API errors that do not contain a response
// body
func TestCheckResponse_noBody(t *testing.T) {
	res := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body:       ioutil.NopCloser(strings.NewReader("")),
	}
	err := CheckResponse(res).(*ErrorResponse)

	if err == nil {
		t.Errorf("Expected error response.")
	}

	expected := &ErrorResponse{
		Response: res,
	}
	if !reflect.DeepEqual(err, expected) {
		t.Errorf("Error = %#v, expected %#v", err, expected)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	res := &http.Response{Request: &http.Request{}}
	err := ErrorResponse{Response: res}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}

func checkCurrentPage(t *testing.T, resp *Response, expectedPage int) {
	links := resp.Links
	p, err := links.CurrentPage()
	if err != nil {
		t.Fatal(err)
	}

	if p != expectedPage {
		t.Fatalf("expected current page to be '%d', was '%d'", expectedPage, p)
	}
}

func TestDo_completion_callback(t *testing.T) {
	setup(t)
	if client == nil {
		t.Fatalf("client == nil")
	}
	defer teardown()

	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := http.MethodGet; m != r.Method {
			t.Errorf("Request method = %v, expected %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest(ctx, http.MethodGet, "/", nil)
	body := new(foo)
	var completedReq *http.Request
	var completedResp string
	client.OnRequestCompleted(func(req *http.Request, resp *http.Response) {
		completedReq = req
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			t.Errorf("Failed to dump response: %s", err)
		}
		completedResp = string(b)
	})
	_, err := client.Do(context.Background(), req, body)
	if err != nil {
		t.Fatalf("Do(): %v", err)
	}
	if !reflect.DeepEqual(req, completedReq) {
		t.Errorf("Completed request = %v, expected %v", completedReq, req)
	}
	expected := `{"A":"a"}`
	if !strings.Contains(completedResp, expected) {
		t.Errorf("expected response to contain %v, Response = %v", expected, completedResp)
	}
}

func TestAddOptions(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected string
		opts     *ListOptions
		isErr    bool
	}{
		{
			name:     "add options",
			path:     "/users",
			expected: "/users?page=1",
			opts:     &ListOptions{Page: 1},
			isErr:    false,
		},
		{
			name:     "add options with existing parameters",
			path:     "/users?scope=all",
			expected: "/users?page=1&scope=all",
			opts:     &ListOptions{Page: 1},
			isErr:    false,
		},
	}

	for _, c := range cases {
		got, err := addOptions(c.path, c.opts)
		if c.isErr && err == nil {
			t.Errorf("%q expected error but none was encountered", c.name)
			continue
		}

		if !c.isErr && err != nil {
			t.Errorf("%q unexpected error: %v", c.name, err)
			continue
		}

		gotURL, err := url.Parse(got)
		if err != nil {
			t.Errorf("%q unable to parse returned URL", c.name)
			continue
		}

		expectedURL, err := url.Parse(c.expected)
		if err != nil {
			t.Errorf("%q unable to parse expected URL", c.name)
			continue
		}

		if g, e := gotURL.Path, expectedURL.Path; g != e {
			t.Errorf("%q path = %q; expected %q", c.name, g, e)
			continue
		}

		if g, e := gotURL.Query(), expectedURL.Query(); !reflect.DeepEqual(g, e) {
			t.Errorf("%q query = %#v; expected %#v", c.name, g, e)
			continue
		}
	}
}

func TestCustomUserAgent(t *testing.T) {
	ua := "testing/0.0.1"
	c, err := New(nil, SetUserAgent(ua))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := fmt.Sprintf("%s %s", ua, userAgent)
	if got := c.UserAgent; got != expected {
		t.Errorf("New() UserAgent = %s; expected %s", got, expected)
	}
}

// func TestCustomBaseURL(t *testing.T) {
// 	baseURL := "http://localhost/foo"
// 	c, err := New(nil, SetBaseURL(baseURL))

// 	if err != nil {
// 		t.Fatalf("New() unexpected error: %v", err)
// 	}

// 	expected := baseURL
// 	if got := c.BaseURL.String(); got != expected {
// 		t.Errorf("New() BaseURL = %s; expected %s", got, expected)
// 	}
// }

func TestCustomBaseURL_badURL(t *testing.T) {
	baseURL := ":"
	_, err := New(nil, SetBaseURL(baseURL))

	testURLParseError(t, err)
}

func intPtr(val int) *int {
	return &val
}

func boolPtr(val bool) *bool {
	return &val
}
