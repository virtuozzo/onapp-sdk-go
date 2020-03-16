package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const networkInterfacesBasePath string = "virtual_machines/%d/network_interfaces"

// NetworkInterfacesService is an interface for interfacing with the NetworkInterface
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/network-interfaces
type NetworkInterfacesService interface {
  List(context.Context, int, *ListOptions) ([]NetworkInterface, *Response, error)
  Get(context.Context, int, int) (*NetworkInterface, *Response, error)
  Create(context.Context, int, *NetworkInterfaceCreateRequest) (*NetworkInterface, *Response, error)
  Delete(context.Context, int, int, interface{}) (*Response, error)
  Edit(context.Context, int, int, *NetworkInterfaceEditRequest) (*Response, error)
}

// NetworkInterfacesServiceOp handles communication with the NetworkInterfaces related methods of the
// OnApp API.
type NetworkInterfacesServiceOp struct {
  client *Client
}

var _ NetworkInterfacesService = &NetworkInterfacesServiceOp{}

// NetworkInterface represents a NetworkInterface
type NetworkInterface struct {
  AdapterType         string      `json:"adapter_type,omitempty"`
  Connected           bool        `json:"connected,bool"`
  CreatedAt           string      `json:"created_at,omitempty"`
  DefaultFirewallRule string      `json:"default_firewall_rule,omitempty"`
  EdgeGatewayID       int         `json:"edge_gateway_id,omitempty"`
  ID                  int         `json:"id,omitempty"`
  Identifier          string      `json:"identifier,omitempty"`
  Label               string      `json:"label,omitempty"`
  MacAddress          string      `json:"mac_address,omitempty"`
  NetworkJoinID       int         `json:"network_join_id"`
  OpenstackID         int         `json:"openstack_id,omitempty"`
  Primary             bool        `json:"primary,bool"`
  RateLimit           int         `json:"rate_limit,omitempty"`
  UpdatedAt           string      `json:"updated_at,omitempty"`
  Usage               bool        `json:"usage,bool"`
  UsageLastResetAt    bool        `json:"usage_last_reset_at,bool"`
  UsageMonthRolledAt  bool        `json:"usage_month_rolled_at,bool"`
  UseAsGateway        bool        `json:"use_as_gateway,bool"`
  VirtualMachineID    int         `json:"virtual_machine_id,omitempty"`
}

// NetworkInterfaceCreateRequest represents a request to create a NetworkInterface
type NetworkInterfaceCreateRequest struct {
  Label          string `json:"label,omitempty"`
  RateLimit      int    `json:"rate_limit,omitempty"`
  NetworkJoinID  int    `json:"network_join_id,omitempty"`
  Primary        bool   `json:"primary,bool"`
}

// NetworkInterfaceEditRequest represents a request to edit a NetworkInterface
type NetworkInterfaceEditRequest struct {
  Label          string `json:"label,omitempty"`
  RateLimit      int    `json:"rate_limit,omitempty"`
}

type networkInterfaceCreateRequestRoot struct {
  NetworkInterfaceCreateRequest  *NetworkInterfaceCreateRequest  `json:"network_interface"`
}

type networkInterfaceRoot struct {
  NetworkInterface  *NetworkInterface  `json:"network_interface"`
}

func (d NetworkInterfaceCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all NetworkInterfaces.
func (s *NetworkInterfacesServiceOp) List(ctx context.Context, vmID int, opt *ListOptions) ([]NetworkInterface, *Response, error) {
  path := fmt.Sprintf(networkInterfacesBasePath, vmID) + apiFormat
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
    arr[i] = out[i]["network_interface"]
  }

  return arr, resp, err
}

// Get individual NetworkInterface.
func (s *NetworkInterfacesServiceOp) Get(ctx context.Context, vmID int, id int) (*NetworkInterface, *Response, error) {
  if vmID < 1 || id < 1 {
    return nil, nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vmID)
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
func (s *NetworkInterfacesServiceOp) Create(ctx context.Context, vmID int, createRequest *NetworkInterfaceCreateRequest) (*NetworkInterface, *Response, error) {
  if vmID < 1 {
    return nil, nil, godo.NewArgError("vmID", "cannot be less than 1")
  }

  if createRequest == nil {
    return nil, nil, godo.NewArgError("NetworkInterface createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vmID) + apiFormat
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
func (s *NetworkInterfacesServiceOp) Delete(ctx context.Context, vmID int, id int, meta interface{}) (*Response, error) {
  if vmID < 1 || id < 1 {
    return nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vmID)
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
func (s *NetworkInterfacesServiceOp) Edit(ctx context.Context, vmID int, id int, editRequest *NetworkInterfaceEditRequest) (*Response, error) {
  if vmID < 1 || id < 1 {
    return nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  if editRequest == nil {
    return nil, godo.NewArgError("NetworkInterface [Edit] editRequest", "cannot be nil")
  }

  path := fmt.Sprintf(networkInterfacesBasePath, vmID)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("NetworkInterface [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}
