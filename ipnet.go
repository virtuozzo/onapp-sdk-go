package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const ipNetsBasePath string = "settings/networks/%d/ip_nets"

// IPNetsService is an interface for interfacing with the IPNet
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ip-nets
type IPNetsService interface {
  List(context.Context, int, *ListOptions) ([]IPNet, *Response, error)
  Get(context.Context, int, int) (*IPNet, *Response, error)
  Create(context.Context, int, *IPNetCreateRequest) (*IPNet, *Response, error)
  Delete(context.Context, int, int, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]IPNet, *Response, error)
}

// IPNetsServiceOp handles communication with the IPNet related methods of the
// OnApp API.
type IPNetsServiceOp struct {
  client *Client
}

var _ IPNetsService = &IPNetsServiceOp{}

// ID - 
type ID struct {
  ID int `json:"id"`
}

// IPNet - 
type IPNet struct {
  ID                  int       `json:"id,omitempty"`
  NetworkAddress      string    `json:"network_address,omitempty"`
  DefaultGateway      string    `json:"default_gateway,omitempty"`
  NetworkMask         int       `json:"network_mask,omitempty"`
  Ipv4                bool      `json:"ipv4,bool"`
  Label               string    `json:"label,omitempty"`
  CreatedAt           string    `json:"created_at,omitempty"`
  UpdatedAt           string    `json:"updated_at,omitempty"`
  OpenstackID         int       `json:"openstack_id,omitempty"`
  Kind                string    `json:"kind,omitempty"`
  GatewayOutsideIPNet bool      `json:"gateway_outside_ip_net,bool"`
  Enabled             bool      `json:"enabled,bool"`
  Network             ID        `json:"network"`
}

// IPNetCreateRequest - 
type IPNetCreateRequest struct {
  Label             string    `json:"label,omitempty"`
  AddDefaultIPRange int       `json:"add_default_ip_range,omitempty"`
  NetworkAddress    string    `json:"network_address,omitempty"`
  NetworkMask       string    `json:"network_mask,omitempty"`
}

type ipNetCreateRequestRoot struct {
  IPNetCreateRequest  *IPNetCreateRequest  `json:"ip_net"`
}

type ipNetRoot struct {
  IPNet  *IPNet  `json:"ip_net"`
}

func (d IPNetCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all IPNet.
func (s *IPNetsServiceOp) List(ctx context.Context, net int, opt *ListOptions) ([]IPNet, *Response, error) {
  if net < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipNetsBasePath, net) + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]IPNet
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]IPNet, len(out))
  for i := range arr {
    arr[i] = out[i]["ip_net"]
  }

  return arr, resp, err
}

// Get individual IPNet.
func (s *IPNetsServiceOp) Get(ctx context.Context, net int, id int) (*IPNet, *Response, error) {
  if net < 1 || id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipNetsBasePath, net)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(ipNetRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.IPNet, resp, err
}

// Create IPNet.
func (s *IPNetsServiceOp) Create(ctx context.Context, net int, createRequest *IPNetCreateRequest) (*IPNet, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("IPNet createRequest", "cannot be nil")
  }

  if net < 1 {
    return nil, nil, godo.NewArgError("net", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipNetsBasePath, net) + apiFormat
  rootRequest := &ipNetCreateRequestRoot {
    IPNetCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("IPNet [Create] req: ", req)

  root := new(ipNetRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.IPNet, resp, err
}

// Delete IPNet.
func (s *IPNetsServiceOp) Delete(ctx context.Context, net int, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipNetsBasePath, net)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("IPNet [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}
