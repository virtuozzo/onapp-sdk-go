package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const remoteTemplatesBasePath string = "templates/available"

// RemoteTemplatesService is an interface for interfacing with the ImageTemplate
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/templates/get-list-of-available-for-installation-templates
// 
// Describe templates *available* for install on the OnApp repository
type RemoteTemplatesService interface {
  List(context.Context, *ListOptions) ([]RemoteTemplate, *Response, error)
  // Get(context.Context, int) (*RemoteTemplate, *Response, error)
  // Create(context.Context, *ImageTemplateCreateRequest) (*RemoteTemplate, *Response, error)
  // Delete(context.Context, int, interface{}) (*Response, error)
  // Edit(context.Context, int, *ImageTemplateEditRequest) (*Response, error)
}

// RemoteTemplatesServiceOp handles communication with the ImageTemplate related methods of the
// OnApp API.
type RemoteTemplatesServiceOp struct {
  client *Client
}

var _ RemoteTemplatesService = &RemoteTemplatesServiceOp{}

// RemoteTemplate - represent a template of OnApp API from repository
type RemoteTemplate struct {
  AllowResizeWithoutReboot  bool   `json:"allow_resize_without_reboot,bool"`
  AllowedHotMigrate         bool   `json:"allowed_hot_migrate,bool"`
  AllowedSwap               bool   `json:"allowed_swap,bool"`
  ApplicationServer         bool   `json:"application_server,bool"`
  BaremetalServer           bool   `json:"baremetal_server,bool"`
  Cdn                       bool   `json:"cdn,bool"`
  Checksum                  string `json:"checksum,omitempty"`
  DiskTargetDevice          string `json:"disk_target_device,omitempty"`
  Ext4                      bool   `json:"ext4,bool"`
  FileName                  string `json:"file_name,omitempty"`
  Label                     string `json:"label,omitempty"`
  ManagerID                 string `json:"manager_id,omitempty"`
  MinDiskSize               int    `json:"min_disk_size"`
  MinMemorySize             int    `json:"min_memory_size"`
  OperatingSystem           string `json:"operating_system,omitempty"`
  OperatingSystemArch       string `json:"operating_system_arch,omitempty"`
  OperatingSystemDistro     string `json:"operating_system_distro,omitempty"`
  OperatingSystemEdition    string `json:"operating_system_edition,omitempty"`
  OperatingSystemTail       string `json:"operating_system_tail,omitempty"`
  ResizeWithoutRebootPolicy string `json:"resize_without_reboot_policy,omitempty"`
  SmartServer               bool   `json:"smart_server,bool"`
  TemplateSize              int    `json:"template_size,omitempty"`
  Version                   string `json:"version,omitempty"`
  Virtualization            string `json:"virtualization,omitempty"`
}

// List all RemoteTemplates.
func (s *RemoteTemplatesServiceOp) List(ctx context.Context, opt *ListOptions) ([]RemoteTemplate, *Response, error) {
  path := remoteTemplatesBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]RemoteTemplate
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]RemoteTemplate, len(out))
  for i := range arr {
    arr[i] = out[i]["remote_template"]
  }

  return arr, resp, err
}
