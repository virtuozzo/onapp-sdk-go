package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const ipRangeBasePath = "settings/networks/%d/ip_nets/%d/ip_ranges"

// IPRangesService is an interface for interfacing with the IPRange
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ip-ranges
type IPRangesService interface {
  List(context.Context, int, int, *ListOptions) ([]IPRange, *Response, error)
  Get(context.Context, int, int, int) (*IPRange, *Response, error)
  Create(context.Context, int, int, *IPRangeCreateRequest) (*IPRange, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, int, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]Network, *Response, error)
}

// IPRangesServiceOp handles communication with the IPRange related methods of the
// OnApp API.
type IPRangesServiceOp struct {
  client *Client
}

var _ IPRangesService = &IPRangesServiceOp{}

type IPRange struct {
  ID             int       `json:"id,omitempty"`
  StartAddress   string    `json:"start_address,omitempty"`
  EndAddress     string    `json:"end_address,omitempty"`
  DefaultGateway string    `json:"default_gateway,omitempty"`
  Ipv4           bool      `json:"ipv4,bool"`
  CreatedAt      string    `json:"created_at,omitempty"`
  UpdatedAt      string    `json:"updated_at,omitempty"`
  Kind           string    `json:"kind,omitempty"`
  IPNet          ID        `json:"ip_net,omitempty"`
}

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
  path := fmt.Sprintf(ipRangeBasePath, net, ipnet) + apiFormat
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
  
  path := fmt.Sprintf(ipRangeBasePath, net, ipnet)
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

  path := fmt.Sprintf(ipRangeBasePath, net, ipnet)
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
    return nil, nil, err
  }

  return root.IPRange, resp, err
}

// Delete IPRange.
func (s *IPRangesServiceOp) Delete(ctx context.Context, net int, ipnet int, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipRangeBasePath, net, ipnet)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }
  log.Println("IPRange [Delete] req: ", req)

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  filter := struct{
    ParentID    int
    ParentType  string
  }{
    ParentID    : id,
    ParentType  : "IPRange",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "IPRange")
}

// Debug - print formatted IPRange structure
func (obj IPRange) Debug() {
  fmt.Printf("            ID: %d\n", obj.ID)
  fmt.Printf("  StartAddress: %s\n", obj.StartAddress)
  fmt.Printf("    EndAddress: %s\n", obj.EndAddress)
  fmt.Printf("DefaultGateway: %s\n", obj.DefaultGateway)
  fmt.Printf("          Ipv4: %t\n", obj.Ipv4)
  fmt.Printf("          Kind: %s\n", obj.Kind)
  fmt.Printf("      IPNet.ID: %d\n", obj.IPNet.ID)
}
