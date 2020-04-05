package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const imageTemplatesBasePath string = "templates"

// ImageTemplatesService is an interface for interfacing with the ImageTemplate
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/templates
//
// Describe templates *installed* on the OnApp cloud
type ImageTemplatesService interface {
	List(context.Context, *ListOptions) ([]ImageTemplate, *Response, error)
	Get(context.Context, int) (*ImageTemplate, *Response, error)
	Create(context.Context, *ImageTemplateCreateRequest) (*ImageTemplate, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *ImageTemplateEditRequest) (*Response, error)
}

// ImageTemplatesServiceOp handles communication with the ImageTemplate related methods of the
// OnApp API.
type ImageTemplatesServiceOp struct {
	client *Client
}

var _ ImageTemplatesService = &ImageTemplatesServiceOp{}

// ImageTemplate - represent a template of OnApp API from cloud
type ImageTemplate struct {
	AllowResizeWithoutReboot  bool                              `json:"allow_resize_without_reboot,bool"`
	AllowedHotMigrate         bool                              `json:"allowed_hot_migrate,bool"`
	AllowedSwap               bool                              `json:"allowed_swap,bool"`
	ApplicationServer         bool                              `json:"application_server,bool"`
	BackupServerID            string                            `json:"backup_server_id"`
	BaremetalServer           bool                              `json:"baremetal_server,bool"`
	Cdn                       bool                              `json:"cdn,bool"`
	Checksum                  string                            `json:"checksum,omitempty"`
	CreatedAt                 string                            `json:"created_at,omitempty"`
	DatacenterID              int                               `json:"datacenter_id,omitempty"`
	DiskTargetDevice          string                            `json:"disk_target_device,omitempty"`
	Draas                     bool                              `json:"draas,bool"`
	Ext4                      bool                              `json:"ext4,bool"`
	FileName                  string                            `json:"file_name,omitempty"`
	ID                        int                               `json:"id,omitempty"`
	Identifier                string                            `json:"identifier,omitempty"`
	InitialPassword           string                            `json:"initial_password,omitempty"`
	InitialUsername           string                            `json:"initial_username,omitempty"`
	Label                     string                            `json:"label,omitempty"`
	Locked                    bool                              `json:"locked,bool"`
	ManagerID                 string                            `json:"manager_id,omitempty"`
	MinDiskSize               int                               `json:"min_disk_size,omitempty"`
	MinMemorySize             int                               `json:"min_memory_size,omitempty"`
	OpenstackID               int                               `json:"openstack_id,omitempty"`
	OperatingSystem           string                            `json:"operating_system,omitempty"`
	OperatingSystemArch       string                            `json:"operating_system_arch,omitempty"`
	OperatingSystemDistro     string                            `json:"operating_system_distro,omitempty"`
	OperatingSystemEdition    string                            `json:"operating_system_edition,omitempty"`
	OperatingSystemTail       string                            `json:"operating_system_tail,omitempty"`
	ParentTemplateID          int                               `json:"parent_template_id,omitempty"`
	Properties                map[string]interface{}            `json:"properties,omitempty"`
	RemoteID                  string                            `json:"remote_id,omitempty"`
	ResizeWithoutRebootPolicy map[string]map[string]interface{} `json:"resize_without_reboot_policy,omitempty"`
	SmartServer               bool                              `json:"smart_server,bool"`
	State                     string                            `json:"state,omitempty"`
	TemplateSize              int                               `json:"template_size,omitempty"`
	Type                      string                            `json:"type,omitempty"`
	UpdatedAt                 string                            `json:"updated_at,omitempty"`
	UserID                    int                               `json:"user_id,omitempty"`
	Version                   string                            `json:"version,omitempty"`
	Virtualization            []string                          `json:"virtualization,omitempty"`
}

// ImageTemplateCreateRequest represents a request to install template
type ImageTemplateCreateRequest struct {
	ManagerID string `json:"manager_id,omitempty"`
	// don't use omitempty because BackupServerID can be empty
	BackupServerID string `json:"backup_server_id"`
}

type imageTemplateCreateRequestRoot struct {
	ImageTemplateCreateRequest *ImageTemplateCreateRequest `json:"image_template"`
}

// ImageTemplateEditRequest represents a request to edit template
type ImageTemplateEditRequest struct {
	Label             string `json:"label,omitempty"`
	FileName          string `json:"file_name,omitempty"`
	Version           string `json:"version,omitempty"`
	MinDiskSize       int    `json:"min_disk_size,omitempty"`
	MinMemorySize     int    `json:"min_memory_size,omitempty"`
	AllowedHotMigrate bool   `json:"allowed_hot_migrate,omitempty"`
}

type imageTemplateEditRequestRoot struct {
	ImageTemplateEditRequest *ImageTemplateEditRequest `json:"image_template"`
}

type imageTemplatesRoot struct {
	ImageTemplate *ImageTemplate `json:"image_template"`
}

// List all ImageTemplates.
func (s *ImageTemplatesServiceOp) List(ctx context.Context, opt *ListOptions) ([]ImageTemplate, *Response, error) {
	path := imageTemplatesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("ImageTemplate [List] req: ", req)

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

	path := fmt.Sprintf("%s/%d%s", imageTemplatesBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("ImageTemplate [Get] req: ", req)

	root := new(imageTemplatesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.ImageTemplate, resp, err
}

// Create ImageTemplate.
func (s *ImageTemplatesServiceOp) Create(ctx context.Context, createRequest *ImageTemplateCreateRequest) (*ImageTemplate, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("ImageTemplate createRequest", "cannot be nil")
	}

	path := imageTemplatesBasePath + apiFormat
	rootRequest := &imageTemplateCreateRequestRoot{
		ImageTemplateCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("ImageTemplate [Create] req: ", req)

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

	path := fmt.Sprintf("%s/%d%s", imageTemplatesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("ImageTemplate [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit ImageTemplate
func (s *ImageTemplatesServiceOp) Edit(ctx context.Context, id int, editRequest *ImageTemplateEditRequest) (*Response, error) {
	if editRequest == nil {
		return nil, godo.NewArgError("ImageTemplate editRequest", "cannot be nil")
	}

	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", imageTemplatesBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("ImageTemplate [Edit] req: ", req)

	return s.client.Do(ctx, req, nil)
}
