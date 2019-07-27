package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const backupServerGroupsBasePath = "settings/backup_server_zones"

// BackupServerGroupsService is an interface for interfacing with the Data Store Zones
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/backup-server-zones
type BackupServerGroupsService interface {
  List(context.Context, *ListOptions) ([]BackupServerGroup, *Response, error)
  Get(context.Context, int) (*BackupServerGroup, *Response, error)
  Create(context.Context, *BackupServerGroupCreateRequest) (*BackupServerGroup, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]BackupServerGroup, *Response, error)
}

// BackupServerGroupsServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type BackupServerGroupsServiceOp struct {
  client *Client
}

var _ BackupServerGroupsService = &BackupServerGroupsServiceOp{}

// BackupServerGroup represents a BackupServerGroup
type BackupServerGroup struct {
  AdditionalFields  AdditionalFields `json:"additional_fields"`
  Closed            bool             `json:"closed"`
  CreatedAt         string           `json:"created_at"`
  DatacenterID      int              `json:"datacenter_id"`
  DraasID           int              `json:"draas_id"`
  FederationEnabled bool             `json:"federation_enabled"`
  FederationID      int              `json:"federation_id"`
  HypervisorID      int              `json:"hypervisor_id"`
  ID                int              `json:"id"`
  Identifier        string           `json:"identifier"`
  Label             string           `json:"label"`
  LocationGroupID   int              `json:"location_group_id"`
  ProviderVdcID     int              `json:"provider_vdc_id"`
  ServerType        string           `json:"server_type"`
  Traded            bool             `json:"traded"`
  UpdatedAt         string           `json:"updated_at"`
}

// BackupServerGroupCreateRequest represents a request to create a BackupServerGroup
type BackupServerGroupCreateRequest struct {
  Label             string  `json:"label,omitempty"`
  LocationGroupID   int     `json:"location_group_id,omitempty"`
  ServerType        string  `json:"server_type,omitempty"`
}

type backupServerGroupCreateRequestRoot struct {
  BackupServerGroupCreateRequest  *BackupServerGroupCreateRequest  `json:"backup_server_group"`
}

type backupServerGroupRoot struct {
  BackupServerGroup  *BackupServerGroup  `json:"backup_server_group"`
}

func (d BackupServerGroupCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *BackupServerGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]BackupServerGroup, *Response, error) {
  path := backupServerGroupsBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]BackupServerGroup
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]BackupServerGroup, len(out))
  for i := range arr {
    arr[i] = out[i]["backup_server_group"]
  }

  return arr, resp, err
}

// Get individual BackupServerGroup.
func (s *BackupServerGroupsServiceOp) Get(ctx context.Context, id int) (*BackupServerGroup, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupServerGroupsBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(backupServerGroupRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupServerGroup, resp, err
}

// Create BackupServerGroup.
func (s *BackupServerGroupsServiceOp) Create(ctx context.Context, createRequest *BackupServerGroupCreateRequest) (*BackupServerGroup, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("BackupServerGroup createRequest", "cannot be nil")
  }

  path := backupServerGroupsBasePath + apiFormat
  rootRequest := &backupServerGroupCreateRequestRoot{
    BackupServerGroupCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nDataStoreGroup [Create] req: ", req)

  root := new(backupServerGroupRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.BackupServerGroup, resp, err
}

// Delete BackupServerGroup.
func (s *BackupServerGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupServerGroupsBasePath, id, apiFormat)
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
    ParentType  : "BackupServerGroup",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "BackupServerGroup")
}

// Debug - print formatted BackupServerGroup structure
func (obj BackupServerGroup) Debug() {
  fmt.Printf("             ID: %d\n", obj.ID)
  fmt.Printf("          Label: %s\n", obj.Label)
  fmt.Printf("     Identifier: %s\n", obj.Identifier)
  fmt.Printf("     ServerType: %s\n", obj.ServerType)
  fmt.Printf("   HypervisorID: %d\n", obj.HypervisorID)
  fmt.Printf("   FederationID: %d\n", obj.FederationID)
  fmt.Printf("LocationGroupID: %d\n", obj.LocationGroupID)
  fmt.Printf("         Traded: %t\n", obj.Traded)
}
