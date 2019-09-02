package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const imageTemplateBasePath = "templates"

// ImageTemplatesService is an interface for interfacing with the ImageTemplate
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/templates
type ImageTemplatesService interface {
  List(context.Context, *ListOptions) ([]ImageTemplate, *Response, error)
  Get(context.Context, int) (*ImageTemplate, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]ImageTemplate, *Response, error)
}

// ImageTemplatesServiceOp handles communication with the ImageTemplate related methods of the
// OnApp API.
type ImageTemplatesServiceOp struct {
  client *Client
}

var _ ImageTemplatesService = &ImageTemplatesServiceOp{}

// ImageTemplate - represent a template of OnApp API
type ImageTemplate struct {
  ID                        int         `json:"id,omitempty"`
  Label                     string      `json:"label,omitempty"`
  CreatedAt                 string      `json:"created_at,omitempty"`
  UpdatedAt                 string      `json:"updated_at,omitempty"`
  Version                   string      `json:"version,omitempty"`
  FileName                  string      `json:"file_name,omitempty"`
  OperatingSystem           string      `json:"operating_system,omitempty"`
  OperatingSystemDistro     string      `json:"operating_system_distro,omitempty"`
  AllowedSwap               bool        `json:"allowed_swap,bool"`
  State                     string      `json:"state,omitempty"`
  Checksum                  string      `json:"checksum,omitempty"`
  AllowResizeWithoutReboot  bool        `json:"allow_resize_without_reboot,bool"`
  MinDiskSize               int         `json:"min_disk_size,omitempty"`
  UserID                    int         `json:"user_id,omitempty"`
  TemplateSize              int         `json:"template_size,omitempty"`
  AllowedHotMigrate         bool        `json:"allowed_hot_migrate,bool"`
  OperatingSystemArch       string      `json:"operating_system_arch,omitempty"`
  OperatingSystemEdition    string      `json:"operating_system_edition,omitempty"`
  OperatingSystemTail       string      `json:"operating_system_tail,omitempty"`
  ParentTemplateID          int         `json:"parent_template_id,omitempty"`
  Virtualization            []string    `json:"virtualization,omitempty"`
  MinMemorySize             int         `json:"min_memory_size,omitempty"`
  DiskTargetDevice          string      `json:"disk_target_device,omitempty"`
  Cdn                       bool        `json:"cdn,bool"`
  BackupServerID            int         `json:"backup_server_id,omitempty"`
  Ext4                      bool        `json:"ext4,bool"`
  SmartServer               bool        `json:"smart_server,bool"`
  BaremetalServer           bool        `json:"baremetal_server,bool"`
  InitialPassword           string      `json:"initial_password,omitempty"`
  InitialUsername           string      `json:"initial_username,omitempty"`
  RemoteID                  int         `json:"remote_id,omitempty"`
  ManagerID                 string      `json:"manager_id,omitempty"`
  // ResizeWithoutRebootPolicy interface{} `json:"resize_without_reboot_policy,omitempty"`
  ApplicationServer         bool        `json:"application_server,bool"`
  Draas                     bool        `json:"draas,bool"`
  Properties                string      `json:"properties,omitempty"`
  Locked                    bool        `json:"locked,bool"`
  OpenstackID               int         `json:"openstack_id,omitempty"`
  Type                      string      `json:"type,omitempty"`
}

type imageTemplatesRoot struct {
  ImageTemplate  *ImageTemplate  `json:"image_template"`
}

// List all ImageTemplates.
func (s *ImageTemplatesServiceOp) List(ctx context.Context, opt *ListOptions) ([]ImageTemplate, *Response, error) {
  path := imageTemplateBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]ImageTemplate
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]ImageTemplate, len(out))
  for i := range arr {
    arr[i] = out[i]["image_template"]
  }

  return arr, resp, err
}

// Get individual ImageTemplate.
func (s *ImageTemplatesServiceOp) Get(ctx context.Context, id int) (*ImageTemplate, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", imageTemplateBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(imageTemplatesRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.ImageTemplate, resp, err
}

// Delete ImageTemplate.
func (s *ImageTemplatesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", imageTemplateBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }

  return s.client.Do(ctx, req, nil)
}
