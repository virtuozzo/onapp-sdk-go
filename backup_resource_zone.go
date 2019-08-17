package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const backupResourceZonesBasePath = "settings/backups/resource_zones"

// BackupResourceZonesService is an interface for interfacing with the Backup Resource Zones
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/backup-server-zones
type BackupResourceZonesService interface {
  List(context.Context, *ListOptions) ([]BackupResourceZone, *Response, error)
  Get(context.Context, int) (*BackupResourceZone, *Response, error)
  Create(context.Context, *BackupResourceZoneCreateRequest) (*BackupResourceZone, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]BackupResourceZone, *Response, error)
}

// BackupResourceZonesServiceOp handles communication with the Backup Resource Zone related methods of the
// OnApp API.
type BackupResourceZonesServiceOp struct {
  client *Client
}

var _ BackupResourceZonesService = &BackupResourceZonesServiceOp{}

// BackupResourceZone represents a BackupResourceZone
type BackupResourceZone struct {
  ID              int     `json:"id,omitempty"`
  CreatedAt       string  `json:"created_at,omitempty"`
  Label           string  `json:"label,omitempty"`
  LocationGroupID int     `json:"location_group_id,omitempty"`
  UpdatedAt       string  `json:"updated_at,omitempty"`
}

// BackupResourceZoneCreateRequest represents a request to create a BackupResourceZone
type BackupResourceZoneCreateRequest struct {
  Label           string  `json:"label,omitempty"`
  LocationGroupID int     `json:"location_group_id,omitempty"`
}

type backupResourceZoneCreateRequestRoot struct {
  BackupResourceZoneCreateRequest  *BackupResourceZoneCreateRequest  `json:"backup_resource_zone"`
}

type backupResourceZoneRoot struct {
  BackupResourceZone  *BackupResourceZone  `json:"backup_resource_zone"`
}

func (d BackupResourceZoneCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *BackupResourceZonesServiceOp) List(ctx context.Context, opt *ListOptions) ([]BackupResourceZone, *Response, error) {
  path := backupResourceZonesBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]BackupResourceZone
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]BackupResourceZone, len(out))
  for i := range arr {
    arr[i] = out[i]["backup_resource_zone"]
  }

  return arr, resp, err
}

// Get individual BackupResourceZone.
func (s *BackupResourceZonesServiceOp) Get(ctx context.Context, id int) (*BackupResourceZone, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupResourceZonesBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(backupResourceZoneRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupResourceZone, resp, err
}

// Create BackupResourceZone.
func (s *BackupResourceZonesServiceOp) Create(ctx context.Context, createRequest *BackupResourceZoneCreateRequest) (*BackupResourceZone, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("BackupResourceZone createRequest", "cannot be nil")
  }

  path := backupResourceZonesBasePath + apiFormat
  rootRequest := &backupResourceZoneCreateRequestRoot{
    BackupResourceZoneCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  log.Println("BackupResourceZone [Create] req: ", req)

  root := new(backupResourceZoneRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.BackupResourceZone, resp, err
}

// Delete BackupResourceZone.
func (s *BackupResourceZonesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupResourceZonesBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }
  log.Println("BackupResourceZone [Delete] req: ", req)

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  filter := struct{
    ParentID    int
    ParentType  string
  }{
    ParentID    : id,
    ParentType  : "BackupResourceZone",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "BackupResourceZone")
}

// Debug - print formatted BackupResourceZone structure
func (obj BackupResourceZone) Debug() {
  fmt.Printf("             ID: %d\n", obj.ID)
  fmt.Printf("          Label: %s\n", obj.Label)
  fmt.Printf("LocationGroupID: %d\n", obj.LocationGroupID)
}
