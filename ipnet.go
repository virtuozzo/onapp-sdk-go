package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const ipNetsBasePath = "settings/networks/%d/ip_nets"

// IPNetsService is an interface for interfacing with the IPNet
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ip-nets
type IPNetsService interface {
  List(context.Context, int, *ListOptions) ([]IPNet, *Response, error)
  Get(context.Context, int, int) (*IPNet, *Response, error)
  Create(context.Context, int, *IPNetCreateRequest) (*IPNet, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]IPNet, *Response, error)
}

// IPNetsServiceOp handles communication with the IPNet related methods of the
// OnApp API.
type IPNetsServiceOp struct {
  client *Client
}

var _ IPNetsService = &IPNetsServiceOp{}

type ID struct {
  ID int `json:"id"`
}

type IPNet struct {
  ID             int       `json:"id,omitempty"`
  NetworkAddress string    `json:"network_address,omitempty"`
  DefaultGateway string    `json:"default_gateway,omitempty"`
  NetworkMask    int       `json:"network_mask,omitempty"`
  Ipv4           bool      `json:"ipv4,bool"`
  Label          string    `json:"label,omitempty"`
  CreatedAt      string    `json:"created_at,omitempty"`
  UpdatedAt      string    `json:"updated_at,omitempty"`
  OpenstackID    int       `json:"openstack_id,omitempty"`
  Kind           string    `json:"kind,omitempty"`
  Enabled        bool      `json:"enabled,bool"`
  Network        ID        `json:"network,omitempty"`
}

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
  if id < 1 {
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

  path := fmt.Sprintf(ipNetsBasePath, net) + apiFormat
  rootRequest := &ipNetCreateRequestRoot {
    IPNetCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nIPNet [Create] req: ", req)

  root := new(ipNetRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.IPNet, resp, err
}

// Delete IPNet.
func (s *IPNetsServiceOp) Delete(ctx context.Context, net int, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(ipNetsBasePath, net)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  filter := struct{
    ParentID    int
    ParentType  string
  }{
    ParentID    : id,
    ParentType  : "IPNet",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "IPNet")
}

// Debug - print formatted IPNet structure
func (obj IPNet) Debug() {
  fmt.Printf("            ID: %d\n", obj.ID)
  fmt.Printf("NetworkAddress: %s\n", obj.NetworkAddress)
  fmt.Printf("   NetworkMask: %d\n", obj.NetworkMask)
  fmt.Printf("DefaultGateway: %s\n", obj.DefaultGateway)
  fmt.Printf("         Label: %s\n", obj.Label)
  fmt.Printf("          Ipv4: %t\n", obj.Ipv4)
  fmt.Printf("          Kind: %s\n", obj.Kind)
  fmt.Printf("       Enabled: %t\n", obj.Enabled)
  fmt.Printf("    Network.ID: %d\n", obj.Network.ID)
}
