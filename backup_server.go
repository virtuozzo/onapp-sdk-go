package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const backupServerBasePath = "settings/backup_servers"

// BackupServersService is an interface for interfacing with the BackupServer
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/backup-servers
type BackupServersService interface {
  List(context.Context, *ListOptions) ([]BackupServer, *Response, error)
  Get(context.Context, int) (*BackupServer, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]BackupServer, *Response, error)
}

// BackupServersServiceOp handles communication with the BackupServer related methods of the
// OnApp API.
type BackupServersServiceOp struct {
  client *Client
}

var _ BackupServersService = &BackupServersServiceOp{}

// BackupServer - represent a template of OnApp API
type BackupServer struct {
  BackupIPAddress     string      `json:"backup_ip_address,omitempty"`
  BackupServerGroupID int         `json:"backup_server_group_id,omitempty"`
  Capacity            int         `json:"capacity,omitempty"`
  CPUIdle             int         `json:"cpu_idle,omitempty"`
  CPUMhz              int         `json:"cpu_mhz,omitempty"`
  Cpus                int         `json:"cpus,omitempty"`
  CreatedAt           string      `json:"created_at,omitempty"`
  Distro              string      `json:"distro,omitempty"`
  Enabled             bool        `json:"enabled,bool"`
  ID                  int         `json:"id,omitempty"`
  IPAddress           string      `json:"ip_address,omitempty"`
  Label               string      `json:"label,omitempty"`
  OsVersion           int         `json:"os_version,omitempty"`
  OsVersionMinor      int         `json:"os_version_minor,omitempty"`
  Release             string      `json:"release,omitempty"`
  TotalMem            int         `json:"total_mem,omitempty"`
  UpdatedAt           string      `json:"updated_at,omitempty"`
  Uptime              string      `json:"uptime,omitempty"`
}

// BackupServerCreateRequest represents a request to create a BackupServer
type BackupServerCreateRequest struct {
  Label           string `json:"label"`
  Enabled         string `json:"enabled"`
  Capacity        string `json:"capacity"`
  IPAddress       string `json:"ip_address"`
  BackupIPAddress string `json:"backup_ip_address"`
}

type backupServersRoot struct {
  BackupServer  *BackupServer  `json:"backup_server"`
}

type backupServerCreateRequestRoot struct {
  BackupServerCreateRequest  *BackupServerCreateRequest  `json:"backup_server"`
}

func (d BackupServerCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all BackupServers.
func (s *BackupServersServiceOp) List(ctx context.Context, opt *ListOptions) ([]BackupServer, *Response, error) {
  path := backupServerBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]BackupServer
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]BackupServer, len(out))
  for i := range arr {
    arr[i] = out[i]["backup_server"]
  }

  return arr, resp, err
}

// Get individual BackupServer.
func (s *BackupServersServiceOp) Get(ctx context.Context, id int) (*BackupServer, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupServerBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(backupServersRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupServer, resp, err
}

// Delete BackupServer.
func (s *BackupServersServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupServerBasePath, id, apiFormat)
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
    AssociatedObjectID    int
    AssociatedObjectType  string
  }{
    AssociatedObjectID    : id,
    AssociatedObjectType  : "BackupServer",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "BackupServer")
}

// Debug - print formatted BackupServer structure
func (obj BackupServer) Debug() {
  fmt.Printf("             ID: %d\n", obj.ID)
  fmt.Printf("          Label: %s\n", obj.Label)
  fmt.Printf("       Capacity: %d\n", obj.Capacity)
  fmt.Printf("      IPAddress: %s\n", obj.IPAddress)
  fmt.Printf("BackupIPAddress: %s\n", obj.BackupIPAddress)
  fmt.Printf("        Enabled: %t\n", obj.Enabled)
  fmt.Printf("         Uptime: %s\n", obj.Uptime)
  fmt.Printf("      UpdatedAt: %s\n", obj.UpdatedAt)
}
