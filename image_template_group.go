package onappgo

import (
  "context"
  "fmt"
  "log"
  "net/http"

  "github.com/digitalocean/godo"
)

const imageTemplateGroupsBasePath string = "settings/image_template_groups"
const relationGroupTemplateBasePath string = "settings/image_template_groups/%d/relation_group_templates"

// ImageTemplateGroupsService is an interface for interfacing with the ImageTemplateGroup
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/template-store
type ImageTemplateGroupsService interface {
  List(context.Context, *ListOptions) ([]ImageTemplateGroup, *Response, error)
  Get(context.Context, int) (*ImageTemplateGroup, *Response, error)
  Create(context.Context, *ImageTemplateGroupCreateRequest) (*ImageTemplateGroup, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  Edit(context.Context, int, *ImageTemplateGroupEditRequest) (*Response, error)

  Attach(context.Context, int, *ImageTemplateGroupAttachRequest) (*ImageTemplateGroup, *Response, error)
  Detach(context.Context, int, int) (*Response, error)
}

// ImageTemplateGroupsServiceOp handles communication with the ImageTemplateGroup related methods of the
// OnApp API.
type ImageTemplateGroupsServiceOp struct {
  client *Client
}

var _ ImageTemplateGroupsService = &ImageTemplateGroupsServiceOp{}

// ImageTemplateGroup - represent a template of OnApp API
type ImageTemplateGroup struct {
  CreatedAt         string `json:"created_at,omitempty"`
  Depth             int    `json:"depth,omitempty"`
  HypervisorGroupID int    `json:"hypervisor_group_id,omitempty"`
  ID                int    `json:"id,omitempty"`
  Kms               bool   `json:"kms,bool"`
  KmsHost           string `json:"kms_host,omitempty"`
  KmsPort           string `json:"kms_port,omitempty"`
  KmsServerLabel    string `json:"kms_server_label,omitempty"`
  Label             string `json:"label,omitempty"`
  Lft               int    `json:"lft,omitempty"`
  Mak               bool   `json:"mak,bool"`
  Own               bool   `json:"own,bool"`
  ParentID          int    `json:"parent_id,omitempty"`
  Rgt               int    `json:"rgt,omitempty"`
  SystemGroup       bool   `json:"system_group,bool"`
  UpdatedAt         string `json:"updated_at,omitempty"`
  UserID            int    `json:"user_id,omitempty"`
}

// ImageTemplateGroupCreateRequest represents a request to create a ImageTemplateGroup
type ImageTemplateGroupCreateRequest struct {
  Kms            bool   `json:"kms,bool"`
  KmsHost        string `json:"kms_host,omitempty"`
  KmsPort        string `json:"kms_port,omitempty"`
  KmsServerLabel string `json:"kms_server_label,omitempty"`
  Label          string `json:"label,omitempty"`
  Mak            bool   `json:"mak,bool"`
  Own            bool   `json:"own,bool"`
  UserID         int    `json:"user_id,omitempty"`
}

// ImageTemplateGroupEditRequest represents a request to edit a ImageTemplateGroup
type ImageTemplateGroupEditRequest struct {
  KmsHost        string `json:"kms_host,omitempty"`
  KmsPort        string `json:"kms_port,omitempty"`
  KmsServerLabel string `json:"kms_server_label,omitempty"`
  Label          string `json:"label,omitempty"`
  Mak            bool   `json:"mak,bool"`
  Own            bool   `json:"own,bool"`
}

type imageTemplateGroupCreateRequestRoot struct {
  ImageTemplateGroupCreateRequest *ImageTemplateGroupCreateRequest `json:"image_template_group"`
}

type imageTemplateGroupsRoot struct {
  ImageTemplateGroup *ImageTemplateGroup `json:"image_template_group"`
}

// ImageTemplateGroupAttachRequest represents a request to attach template to the ImageTemplateGroup
type ImageTemplateGroupAttachRequest struct {
  TemplateID  int     `json:"template_id,omitempty"`
  // Price       float64 `json:"price,omitempty"`
}

// Attach Template to Template Group
type imageTemplateGroupAttachRequestRoot struct {
  ImageTemplateGroupAttachRequest *ImageTemplateGroupAttachRequest `json:"relation_group_template"`
}

// List all ImageTemplateGroups.
func (s *ImageTemplateGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]ImageTemplateGroup, *Response, error) {
  path := imageTemplateGroupsBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]ImageTemplateGroup
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]ImageTemplateGroup, len(out))
  for i := range arr {
    arr[i] = out[i]["image_template_group"]
  }

  return arr, resp, err
}

// Get individual ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Get(ctx context.Context, id int) (*ImageTemplateGroup, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", imageTemplateGroupsBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(imageTemplateGroupsRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.ImageTemplateGroup, resp, err
}

// Create ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Create(ctx context.Context, createRequest *ImageTemplateGroupCreateRequest) (*ImageTemplateGroup, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("ImageTemplateGroup createRequest", "cannot be nil")
  }

  path := imageTemplateGroupsBasePath + apiFormat
  rootRequest := &imageTemplateGroupCreateRequestRoot{
    ImageTemplateGroupCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("ImageTemplateGroup [Create] req: ", req)

  root := new(imageTemplateGroupsRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.ImageTemplateGroup, resp, err
}

// Delete ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", imageTemplateGroupsBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("ImageTemplateGroup [Delete]  req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *ImageTemplateGroupEditRequest) (*Response, error) {
  if editRequest == nil {
    return nil, godo.NewArgError("ImageTemplateGroup editRequest", "cannot be nil")
  }

  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", imageTemplateGroupsBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("ImageTemplateGroup [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Attach template to the ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Attach(ctx context.Context, groupID int, attachRequest *ImageTemplateGroupAttachRequest) (*ImageTemplateGroup, *Response, error) {
  if groupID < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  if attachRequest == nil {
    return nil, nil, godo.NewArgError("ImageTemplateGroup attachRequest", "cannot be nil")
  }

  path := fmt.Sprintf(relationGroupTemplateBasePath, groupID)
  rootRequest := &imageTemplateGroupAttachRequestRoot{
    ImageTemplateGroupAttachRequest: attachRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("ImageTemplateGroup [Attach] req: ", req)

  root := new(imageTemplateGroupsRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.ImageTemplateGroup, resp, err
}

// Detach template to the ImageTemplateGroup.
func (s *ImageTemplateGroupsServiceOp) Detach(ctx context.Context, groupID int, id int) (*Response, error) {
  if groupID < 1 || id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(relationGroupTemplateBasePath, groupID)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("ImageTemplateGroup [Detach] req: ", req)

  root := new(imageTemplateGroupsRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return resp, err
  }

  return resp, err
}
