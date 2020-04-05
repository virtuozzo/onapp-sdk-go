package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const instancePackagesBasePath string = "instance_packages"

// InstancePackagesService is an interface for interfacing with the Instance Packages
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/instance-packages
type InstancePackagesService interface {
	List(context.Context, *ListOptions) ([]InstancePackage, *Response, error)
	Get(context.Context, int) (*InstancePackage, *Response, error)
	Create(context.Context, *InstancePackageCreateRequest) (*InstancePackage, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *InstancePackageEditRequest) (*Response, error)
}

// InstancePackagesServiceOp handles communication with the Instance Package related methods of the
// OnApp API.
type InstancePackagesServiceOp struct {
	client *Client
}

var _ InstancePackagesService = &InstancePackagesServiceOp{}

// InstancePackage represents a InstancePackage
type InstancePackage struct {
	Bandwidth   int    `json:"bandwidth,omitempty"`
	BucketsIds  []int  `json:"buckets_ids,omitempty"`
	Cpus        int    `json:"cpus,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	DiskSize    int    `json:"disk_size,omitempty"`
	ID          int    `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Memory      int    `json:"memory,omitempty"`
	OpenstackID int    `json:"openstack_id,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
}

// InstancePackageCreateRequest represents a request to create a InstancePackage
type InstancePackageCreateRequest struct {
	Label     string `json:"label,omitempty"`
	Cpus      int    `json:"cpus,omitempty"`
	Memory    int    `json:"memory,omitempty"`
	DiskSize  int    `json:"disk_size,omitempty"`
	Bandwidth int    `json:"bandwidth,omitempty"`
}

type InstancePackageEditRequest InstancePackageCreateRequest

type instancePackageCreateRequestRoot struct {
	InstancePackageCreateRequest *InstancePackageCreateRequest `json:"instance_package"`
}

type instancePackageRoot struct {
	InstancePackage *InstancePackage `json:"instance_package"`
}

func (d InstancePackageCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *InstancePackagesServiceOp) List(ctx context.Context, opt *ListOptions) ([]InstancePackage, *Response, error) {
	path := instancePackagesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]InstancePackage
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]InstancePackage, len(out))
	for i := range arr {
		arr[i] = out[i]["instance_package"]
	}

	return arr, resp, err
}

// Get individual InstancePackage.
func (s *InstancePackagesServiceOp) Get(ctx context.Context, id int) (*InstancePackage, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", instancePackagesBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(instancePackageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePackage, resp, err
}

// Create InstancePackage.
func (s *InstancePackagesServiceOp) Create(ctx context.Context, createRequest *InstancePackageCreateRequest) (*InstancePackage, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("InstancePackage createRequest", "cannot be nil")
	}

	path := instancePackagesBasePath + apiFormat
	rootRequest := &instancePackageCreateRequestRoot{
		InstancePackageCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("InstancePackage [Create] req: ", req)

	root := new(instancePackageRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.InstancePackage, resp, err
}

// Delete InstancePackage.
func (s *InstancePackagesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", instancePackagesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("InstancePackage [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit InstancePackage.
func (s *InstancePackagesServiceOp) Edit(ctx context.Context, id int, editRequest *InstancePackageEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("InstancePackage editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", instancePackagesBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("InstancePackage [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
