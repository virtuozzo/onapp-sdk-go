package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const ipRangesBasePath string = "settings/networks/%d/ip_nets/%d/ip_ranges"

// IPRangesService is an interface for interfacing with the IPRange
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ip-ranges
type IPRangesService interface {
  List(context.Context, int, int, *ListOptions) ([]IPRange, *Response, error)
  Get(context.Context, int, int, int) (*IPRange, *Response, error)
  Create(context.Context, int, int, *IPRangeCreateRequest) (*IPRange, *Response, error)
  Delete(context.Context, int, int, int, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]Network, *Response, error)
}

// IPRangesServiceOp handles communication with the IPRange related methods of the
// OnApp API.
type IPRangesServiceOp struct {
  client *Client
}

var _ IPRangesService = &IPRangesServiceOp{}

// IPRange - 
type IPRange struct {
  ID                  int     `json:"id,omitempty"`
  StartAddress        string  `json:"start_address,omitempty"`
  EndAddress          string  `json:"end_address,omitempty"`
  DefaultGateway      string  `json:"default_gateway,omitempty"`
  Ipv4                bool    `json:"ipv4,bool"`
  CreatedAt           string  `json:"created_at,omitempty"`
  UpdatedAt           string  `json:"updated_at,omitempty"`
  Kind                string  `json:"kind,omitempty"`
  GatewayOutsideIPNet bool    `json:"gateway_outside_ip_net,bool"`
  IPNet               ID      `json:"ip_net"`
}

// IPRangeCreateRequest - 
type IPRangeCreateRequest struct {
  StartAddress      string    `json:"start_address,omitempty"`
  EndAddress        string    `json:"end_address,omitempty"`
  DefaultGateway    string    `json:"default_gateway,omitempty"`
}

type ipRangeCreateRequestRoot struct {
  IPRangeCreateRequest  *IPRangeCreateRequest  `json:"ip_range"`
}

type ipRangeRoot struct {
  IPRange  *IPRange  `json:"ip_range"`
}

func (d IPRangeCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all IPRanges.
func (s *IPRangesServiceOp) List(ctx context.Context, net int, ipnet int, opt *ListOptions) ([]IPRange, *Response, error) {
  path := fmt.Sprintf(ipRangesBasePath, net, ipnet) + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]IPRange
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]IPRange, len(out))
  for i := range arr {
    arr[i] = out[i]["ip_range"]
  }

  return arr, resp, err
}

// Get individual IPRange.
func (s *IPRangesServiceOp) Get(ctx context.Context, net int, ipnet int, id int) (*IPRange, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }
  
  path := fmt.Sprintf(ipRangesBasePath, net, ipnet)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(ipRangeRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.IPRange, resp, err
}

// Create IPRange.
func (s *IPRangesServiceOp) Create(ctx context.Context, net int, ipnet int, createRequest *IPRangeCreateRequest) (*IPRange, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("Network createRequest", "cannot be nil")
  }

  if net < 1 || ipnet < 1 {
    return nil, nil, godo.NewArgError("net || ipnet", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipRangesBasePath, net, ipnet)
  rootRequest := &ipRangeCreateRequestRoot {
    IPRangeCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("IPRange [Create] req: ", req)

  root := new(ipRangeRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.IPRange, resp, err
}

// Delete IPRange.
func (s *IPRangesServiceOp) Delete(ctx context.Context, net int, ipnet int, id int, meta interface{}) (*Response, error) {
  if net < 1 || ipnet < 1 || id < 1 {
    return nil, godo.NewArgError("net || ipnet || id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipRangesBasePath, net, ipnet)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("IPRange [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}
