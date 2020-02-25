package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const networkInterfacesBasePath string = "virtual_machines/%s/network_interfaces"

// NetworkInterfacesService is an interface for interfacing with the NetworkInterface
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/networks
type NetworkInterfacesService interface {
  List(context.Context, string, *ListOptions) ([]NetworkInterface, *Response, error)
  Get(context.Context, string, int) (*NetworkInterface, *Response, error)
  Create(context.Context, string, *NetworkInterfaceCreateRequest) (*NetworkInterface, *Response, error)
  Delete(context.Context, string, int, interface{}) (*Response, error)
  Edit(context.Context, string, int, *NetworkInterfaceEditRequest) (*Response, error)
}

// NetworkInterfacesServiceOp handles communication with the NetworkInterfaces related methods of the
// OnApp API.
type NetworkInterfacesServiceOp struct {
  client *Client
}

var _ NetworkInterfacesService = &NetworkInterfacesServiceOp{}

// NetworkInterface represents a NetworkInterface
type NetworkInterface Network

// NetworkInterfaceCreateRequest represents a request to create a NetworkInterface
type NetworkInterfaceCreateRequest struct {
  Label          string `json:"label,omitempty"`
  NetworkJoinID  int    `json:"network_join_id,omitempty"`
  RateLimit      int    `json:"rate_limit,omitempty"`
  Primary        int    `json:"primary,omitempty"`
}

// NetworkInterfaceEditRequest represents a request to edit a NetworkInterface
type NetworkInterfaceEditRequest struct {
  Label          string `json:"label,omitempty"`
  RateLimit      int    `json:"rate_limit,omitempty"`
}

type networkInterfaceCreateRequestRoot struct {
  NetworkInterfaceCreateRequest  *NetworkInterfaceCreateRequest  `json:"network"`
}

type networkInterfaceRoot struct {
  NetworkInterface  *NetworkInterface  `json:"network"`
}

func (d NetworkInterfaceCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all NetworkInterfaces.
func (s *NetworkInterfacesServiceOp) List(ctx context.Context, vm string, opt *ListOptions) ([]NetworkInterface, *Response, error) {
  path := fmt.Sprintf(networkInterfacesBasePath, vm) + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]NetworkInterface
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]NetworkInterface, len(out))
  for i := range arr {
    arr[i] = out[i]["network"]
  }

  return arr, resp, err
}

// Get individual NetworkInterface.
func (s *NetworkInterfacesServiceOp) Get(ctx context.Context, vm string, id int) (*NetworkInterface, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vm)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(networkInterfaceRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkInterface, resp, err
}

// Create NetworkInterface.
func (s *NetworkInterfacesServiceOp) Create(ctx context.Context, vm string, createRequest *NetworkInterfaceCreateRequest) (*NetworkInterface, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("NetworkInterface createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vm) + apiFormat
  rootRequest := &networkInterfaceCreateRequestRoot{
    NetworkInterfaceCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  fmt.Println("NetworkInterface [Create] req: ", req)

  root := new(networkInterfaceRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkInterface, resp, err
}

// Delete NetworkInterface.
func (s *NetworkInterfacesServiceOp) Delete(ctx context.Context, vm string, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vm)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  fmt.Println("NetworkInterface [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit NetworkInterface.
func (s *NetworkInterfacesServiceOp) Edit(ctx context.Context, vm string, id int, editRequest *NetworkInterfaceEditRequest) (*Response, error) {
  path := fmt.Sprintf(networkInterfacesBasePath, vm)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("NetworkInterface [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}
