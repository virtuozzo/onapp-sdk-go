package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const userGroupsBasePath string = "user_groups"

// UserGroupsService is an interface for interfacing with the UserGroup
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/user-groups
type UserGroupsService interface {
	List(context.Context, *ListOptions) ([]UserGroup, *Response, error)
	Get(context.Context, int) (*UserGroup, *Response, error)
	Create(context.Context, *UserGroupCreateRequest) (*UserGroup, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *UserGroupEditRequest) (*Response, error)
}

// UserGroupsServiceOp handles communication with the UserGroup related methods of the
// OnApp API.
type UserGroupsServiceOp struct {
	client *Client
}

var _ UserGroupsService = &UserGroupsServiceOp{}

// UserBucket -
type UserBucket struct {
	AllowsKms    bool   `json:"allows_kms,bool"`
	AllowsMak    bool   `json:"allows_mak,bool"`
	AllowsOwn    bool   `json:"allows_own,bool"`
	CreatedAt    string `json:"created_at,omitempty"`
	CurrencyCode string `json:"currency_code,omitempty"`
	ID           int    `json:"id,omitempty"`
	Label        string `json:"label,omitempty"`
	ShowPrice    bool   `json:"show_price,bool"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

// UserBuckets -
type UserBuckets struct {
	UserBucket UserBucket `json:"user_bucket"`
}

// UserGroup -
type UserGroup struct {
	// Must be as map[string]interface{}
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`

	BucketID          int           `json:"bucket_id,omitempty"`
	CreatedAt         string        `json:"created_at,omitempty"`
	DatacenterID      int           `json:"datacenter_id,omitempty"`
	DraasID           int           `json:"draas_id,omitempty"`
	HypervisorID      int           `json:"hypervisor_id,omitempty"`
	ID                int           `json:"id,omitempty"`
	Identifier        string        `json:"identifier,omitempty"`
	Label             string        `json:"label,omitempty"`
	PreconfiguredOnly bool          `json:"preconfigured_only,bool"`
	ProviderVdcID     int           `json:"provider_vdc_id,omitempty"`
	Roles             []Roles       `json:"roles,omitempty"`
	UpdatedAt         string        `json:"updated_at,omitempty"`
	UserBuckets       []UserBuckets `json:"user_buckets,omitempty"`
}

// UserGroupCreateRequest -
type UserGroupCreateRequest struct {
	Label     string `json:"label,omitempty"`
	BucketIDs []int  `json:"bucket_ids,omitempty"`
}

// UserGroupEditRequest -
type UserGroupEditRequest struct {
	Label string `json:"label,omitempty"`

	RoleIDs   []string `json:"role_ids,omitempty"`
	BucketIDs []int    `json:"bucket_ids,omitempty"`
}

type userGroupCreateRequestRoot struct {
	UserGroupCreateRequest *UserGroupCreateRequest `json:"user_group"`
}

type userGroupRoot struct {
	UserGroup *UserGroup `json:"user_group"`
}

func (d UserGroupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Users.
func (s *UserGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]UserGroup, *Response, error) {
	path := userGroupsBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]UserGroup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]UserGroup, len(out))
	for i := range arr {
		arr[i] = out[i]["user_group"]
	}

	return arr, resp, err
}

// Get individual UserGroup.
func (s *UserGroupsServiceOp) Get(ctx context.Context, id int) (*UserGroup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", userGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserGroup, resp, err
}

// Create UserGroup.
func (s *UserGroupsServiceOp) Create(ctx context.Context, createRequest *UserGroupCreateRequest) (*UserGroup, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := userGroupsBasePath + apiFormat
	rootRequest := &userGroupCreateRequestRoot{
		UserGroupCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("UserGroup [Create]  req: ", req)

	root := new(userGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserGroup, resp, err
}

// Delete UserGroup.
func (s *UserGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", userGroupsBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("UserGroup [Delete]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit UserGroup.
func (s *UserGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *UserGroupEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("UserGroup [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", userGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("UserGroup [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
