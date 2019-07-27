package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const backupResourcesBasePath = "settings/backups/resources"

// BackupResourcesService is an interface for interfacing with the Data Store Zones
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/backup-server-zones
type BackupResourcesService interface {
  List(context.Context, *ListOptions) ([]BackupResource, *Response, error)
  Get(context.Context, int) (*BackupResource, *Response, error)
  Create(context.Context, *BackupResourceCreateRequest) (*BackupResource, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]BackupResource, *Response, error)
}

// BackupResourcesServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type BackupResourcesServiceOp struct {
  client *Client
}

var _ BackupResourcesService = &BackupResourcesServiceOp{}

// BackupResource represents a BackupResource
type BackupResource struct {
  ID              int     `json:"id,omitempty"`
  CreatedAt       string  `json:"created_at,omitempty"`
  Label           string  `json:"label,omitempty"`
  Enabled         bool    `json:"enabled,bool"`
  Plugin          string  `json:"plugin,omitempty"`
  PrimaryHost     string  `json:"primary_host,omitempty"`
  SecondaryHost   string  `json:"secondary_host,omitempty"`
  UpdatedAt       string  `json:"updated_at,omitempty,omitempty"`
  Username        string  `json:"username,omitempty"`
  Password        string  `json:"password,omitempty"`
  ResourceZoneID  int     `json:"resource_zone_id,omitempty"`
}

// BackupResourceCreateRequest represents a request to create a BackupResource
type BackupResourceCreateRequest struct {
  Label           string  `json:"label,omitempty"`
  Plugin          string  `json:"plugin,omitempty"`
  PrimaryHost     string  `json:"primary_host,omitempty"`
  Username        string  `json:"username,omitempty"`
  Password        string  `json:"password,omitempty"`
  ResourceZoneID  int     `json:"resource_zone_id,omitempty"`
}

type backupResourceCreateRequestRoot struct {
  BackupResourceCreateRequest  *BackupResourceCreateRequest  `json:"backup_resource"`
}

type backupResourceRoot struct {
  BackupResource  *BackupResource  `json:"backup_resource"`
}

func (d BackupResourceCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *BackupResourcesServiceOp) List(ctx context.Context, opt *ListOptions) ([]BackupResource, *Response, error) {
  path := backupResourcesBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]BackupResource
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]BackupResource, len(out))
  for i := range arr {
    arr[i] = out[i]["backup_resource"]
  }

  return arr, resp, err
}

// Get individual BackupResource.
func (s *BackupResourcesServiceOp) Get(ctx context.Context, id int) (*BackupResource, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupResourcesBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(backupResourceRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupResource, resp, err
}

// Create BackupResource.
func (s *BackupResourcesServiceOp) Create(ctx context.Context, createRequest *BackupResourceCreateRequest) (*BackupResource, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("BackupResource createRequest", "cannot be nil")
  }

  path := backupResourcesBasePath + apiFormat
  rootRequest := &backupResourceCreateRequestRoot{
    BackupResourceCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nBackupResource [Create] req: ", req)

  root := new(backupResourceRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.BackupResource, resp, err
}

// Delete BackupResource.
func (s *BackupResourcesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupResourcesBasePath, id, apiFormat)
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
    ParentType  : "BackupResource",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "BackupResource")
}

// Debug - print formatted BackupResource structure
func (obj BackupResource) Debug() {
  fmt.Printf("            ID: %d\n", obj.ID)
  fmt.Printf("         Label: %s\n", obj.Label)
  fmt.Printf("   PrimaryHost: %s\n", obj.PrimaryHost)
  fmt.Printf("      Username: %s\n", obj.Username)
  fmt.Printf("ResourceZoneID: %d\n", obj.ResourceZoneID)
  fmt.Printf("       Enabled: %t\n", obj.Enabled)
}
