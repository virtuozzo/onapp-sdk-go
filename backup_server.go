package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const backupServersBasePath string = "settings/backup_servers"

// BackupServersService is an interface for interfacing with the Backup Server
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/backup-servers
type BackupServersService interface {
  List(context.Context, *ListOptions) ([]BackupServer, *Response, error)
  Get(context.Context, int) (*BackupServer, *Response, error)
  Create(context.Context, *BackupServerCreateRequest) (*BackupServer, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  Edit(context.Context, int, *BackupServerEditRequest) (*Response, error)
}

// BackupServersServiceOp handles communication with the Backup Server related methods of the
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

  // OnApp 6.1
  IntegratedStorage   bool        `json:"integrated_storage,bool"`
}

// BackupServerCreateRequest represents a request to create a BackupServer
type BackupServerCreateRequest struct {
  Label               string `json:"label,omitempty"`
  Enabled             bool   `json:"enabled,bool"`
  Capacity            int    `json:"capacity,omitempty"`
  IPAddress           string `json:"ip_address,omitempty"`
  BackupIPAddress     string `json:"backup_ip_address,omitempty"`
  BackupServerGroupID int    `json:"backup_server_group_id,omitempty"`

  // OnApp 6.1
  // default - false, enable - true
  IntegratedStorage   bool   `json:"integrated_storage,bool"`
}

// BackupServerEditRequest represents a request to edit a BackupServer
type BackupServerEditRequest struct {
  Label               string `json:"label,omitempty"`
  Enabled             bool   `json:"enabled,bool"`
  Capacity            int    `json:"capacity,omitempty"`
  IPAddress           string `json:"ip_address,omitempty"`
  BackupIPAddress     string `json:"backup_ip_address,omitempty"`
  IntegratedStorage   bool   `json:"integrated_storage,bool"`
}

type backupServerRoot struct {
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
  path := backupServersBasePath + apiFormat
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

  path := fmt.Sprintf("%s/%d%s", backupServersBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(backupServerRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupServer, resp, err
}

// Create BackupServer.
func (s *BackupServersServiceOp) Create(ctx context.Context, createRequest *BackupServerCreateRequest) (*BackupServer, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("BackupServer createRequest", "cannot be nil")
  }

  path := backupServersBasePath + apiFormat
  rootRequest := &backupServerCreateRequestRoot{
    BackupServerCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("BackupServer [Create] req: ", req)

  root := new(backupServerRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.BackupServer, resp, err
}

// Delete BackupServer.
func (s *BackupServersServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", backupServersBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("BackupServer [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit BackupServer.
func (s *BackupServersServiceOp) Edit(ctx context.Context, id int, editRequest *BackupServerEditRequest) (*Response, error) {
  path := fmt.Sprintf("%s/%d%s", backupServersBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("BackupServer [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}
