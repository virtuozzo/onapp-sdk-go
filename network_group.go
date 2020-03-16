package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const networkZonesBasePath string = "settings/network_zones"

// NetworkGroupsService is an interface for interfacing with the NetworkGroup
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/network-zones
type NetworkGroupsService interface {
  List(context.Context, *ListOptions) ([]NetworkGroup, *Response, error)
  Get(context.Context, int) (*NetworkGroup, *Response, error)
  Create(context.Context, *NetworkGroupCreateRequest) (*NetworkGroup, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  Edit(context.Context, int, *NetworkGroupEditRequest) (*Response, error)
}

// NetworkGroupsServiceOp handles communication with the NetworkGroups related methods of the
// OnApp API.
type NetworkGroupsServiceOp struct {
  client *Client
}

var _ NetworkGroupsService = &NetworkGroupsServiceOp{}

// NetworkGroup represents a NetworkGroup
type NetworkGroup struct {
  AdditionalFields  []AdditionalFields  `json:"additional_fields,omitempty"`
  Closed            bool                `json:"closed,bool"`
  CreatedAt         string              `json:"created_at,omitempty"`
  DatacenterID      int                 `json:"datacenter_id,omitempty"`
  DraasID           int                 `json:"draas_id,omitempty"`
  FederationEnabled bool                `json:"federation_enabled,bool"`
  FederationID      int                 `json:"federation_id,omitempty"`
  HypervisorID      int                 `json:"hypervisor_id,omitempty"`
  ID                int                 `json:"id,omitempty"`
  Identifier        string              `json:"identifier,omitempty"`
  Label             string              `json:"label,omitempty"`
  LocationGroupID   int                 `json:"location_group_id,omitempty"`
  PreconfiguredOnly bool                `json:"preconfigured_only,bool"`
  ProviderVdcID     int                 `json:"provider_vdc_id,omitempty"`
  ServerType        string              `json:"server_type,omitempty"`
  Traded            bool                `json:"traded,bool"`
  UpdatedAt         string              `json:"updated_at,omitempty"`
}

// NetworkGroupCreateRequest represents a request to create a NetworkGroup
type NetworkGroupCreateRequest struct {
  Label             string  `json:"label,omitempty"`
  LocationGroupID   int     `json:"location_group_id,omitempty"`
  PreconfiguredOnly bool    `json:"preconfigured_only,bool"`
  ServerType        string  `json:"server_type,omitempty"`
}

// NetworkGroupEditRequest represents a request to edit a NetworkGroup
type NetworkGroupEditRequest struct {
  Label             string  `json:"label,omitempty"`
  LocationGroupID   int     `json:"location_group_id,omitempty"`
  PreconfiguredOnly bool    `json:"preconfigured_only,bool"`
}

type networkZoneCreateRequestRoot struct {
  NetworkGroupCreateRequest  *NetworkGroupCreateRequest  `json:"network_group"`
}

type networkZoneRoot struct {
  NetworkGroup  *NetworkGroup  `json:"network_group"`
}

func (d NetworkGroupCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all NetworkGroups.
func (s *NetworkGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]NetworkGroup, *Response, error) {
  path := networkZonesBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]NetworkGroup
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]NetworkGroup, len(out))
  for i := range arr {
    arr[i] = out[i]["network_group"]
  }

  return arr, resp, err
}

// Get individual NetworkGroup.
func (s *NetworkGroupsServiceOp) Get(ctx context.Context, id int) (*NetworkGroup, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", networkZonesBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(networkZoneRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkGroup, resp, err
}

// Create NetworkGroup.
func (s *NetworkGroupsServiceOp) Create(ctx context.Context, createRequest *NetworkGroupCreateRequest) (*NetworkGroup, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("NetworkGroup createRequest", "cannot be nil")
  }

  path := networkZonesBasePath + apiFormat
  rootRequest := &networkZoneCreateRequestRoot{
    NetworkGroupCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("NetworkGroup [Create] req: ", req)

  root := new(networkZoneRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkGroup, resp, err
}

// Delete NetworkGroup.
func (s *NetworkGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", networkZonesBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("NetworkGroup [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit NetworkGroup.
func (s *NetworkGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *NetworkGroupEditRequest) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  if editRequest == nil {
    return nil, godo.NewArgError("NetworkGroup [Edit] editRequest", "cannot be nil")
  }

  path := fmt.Sprintf("%s/%d%s", networkZonesBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("NetworkGroup [Edit] req: ", req)

  return s.client.Do(ctx, req, nil)
}
