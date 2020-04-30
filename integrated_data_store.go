package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const integratedDataStoresBasePath string = "storage/%d/data_stores"
const backendNodesBasePath string = "storage/%d/hypervisors"

// IntegratedDataStoresService is an interface for interfacing with the IntegrateDataStores
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/integrated-storage
type IntegratedDataStoresService interface {
	List(context.Context, int, *ListOptions) ([]IntegratedDataStores, *Response, error)
	Get(context.Context, int, string) (*IntegratedDataStores, *Response, error)
	Create(context.Context, int, *IntegratedDataStoreCreateRequest) (*IntegratedDataStores, *Response, error)
	Delete(context.Context, int, string, interface{}) (*Response, error)
	Edit(context.Context, int, string, *IntegratedDataStoresEditRequest) (*Response, error)

	BackendNodes(context.Context, int) (*BackendNodes, *Response, error)
}

// IntegratedDataStoresServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type IntegratedDataStoresServiceOp struct {
	client *Client
}

var _ IntegratedDataStoresService = &IntegratedDataStoresServiceOp{}

type Node struct {
	ID string `json:"id,omitempty"`
}

type Nodes struct {
	Node Node `json:"node,omitempty"`
}

type IntegratedDataStores struct {
	ID              string  `json:"id,omitempty"`
	Name            string  `json:"name,omitempty"`
	Replicas        int     `json:"replicas,omitempty"`
	Stripes         int     `json:"stripes,omitempty"`
	Overcommit      int     `json:"overcommit,omitempty"`
	TotalSize       int64   `json:"total_size,omitempty"`
	FreeSize        int64   `json:"free_size,omitempty"`
	MaximumDiskSize int64   `json:"maximum_disk_size,omitempty"`
	Performance     int     `json:"performance,omitempty"`
	DiskCount       int     `json:"disk_count,omitempty"`
	Nodes           []Nodes `json:"nodes,omitempty"`
}

type HypervisorNode struct {
	ID    string  `json:"id,omitempty"`
	Nodes []Nodes `json:"nodes,omitempty"`
}

type BackendNodes []struct {
	HypervisorNode *HypervisorNode `json:"hypervisor,omitempty"`
}

// IntegratedDataStoreCreateRequest represents a request to create a IntegrateDataStores
type IntegratedDataStoreCreateRequest struct {
	Name       string   `json:"name,omitempty"`
	Replicas   int      `json:"replicas,omitempty"`
	Stripes    int      `json:"stripes,omitempty"`
	Overcommit int      `json:"overcommit,omitempty"`
	NodeIDs    []string `json:"node_ids,omitempty"`
}

// IntegratedDataStoresEditRequest represents a request to edit a IntegrateDataStores
type IntegratedDataStoresEditRequest struct {
	OwnerIDs []string `json:"owner_ids,omitempty"`
}

type integratedDataStoreCreateRequestRoot struct {
	IntegratedDataStoreCreateRequest *IntegratedDataStoreCreateRequest `json:"storage_data_store"`
}

type integratedDataStoreRoot struct {
	IntegratedDataStores *IntegratedDataStores `json:"data_store"`
}

func (d IntegratedDataStoreCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all
func (s *IntegratedDataStoresServiceOp) List(ctx context.Context, resID int, opt *ListOptions) ([]IntegratedDataStores, *Response, error) {
	if resID < 1 {
		return nil, nil, godo.NewArgError("resID", "cannot be less than 1")
	}

	path := fmt.Sprintf(integratedDataStoresBasePath, resID) + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("IntegratedDataStores [List]  req: ", req)

	var out []map[string]IntegratedDataStores
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]IntegratedDataStores, len(out))
	for i := range arr {
		arr[i] = out[i]["data_store"]
	}

	return arr, resp, err
}

// Get individual
func (s *IntegratedDataStoresServiceOp) Get(ctx context.Context, resID int, id string) (*IntegratedDataStores, *Response, error) {
	if resID < 1 || id == "" {
		return nil, nil, godo.NewArgError("resID or id", "cannot be empty or less than 1")
	}

	path := fmt.Sprintf(integratedDataStoresBasePath, resID)
	path = fmt.Sprintf("%s/%s%s", path, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("IntegratedDataStores [Get]  req: ", req)

	root := new(integratedDataStoreRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.IntegratedDataStores, resp, err
}

// Create -
func (s *IntegratedDataStoresServiceOp) Create(ctx context.Context, resID int, createRequest *IntegratedDataStoreCreateRequest) (*IntegratedDataStores, *Response, error) {
	if resID < 1 || createRequest == nil {
		return nil, nil, godo.NewArgError("IntegratedDataStores createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(integratedDataStoresBasePath, resID) + apiFormat
	rootRequest := &integratedDataStoreCreateRequestRoot{
		IntegratedDataStoreCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("IntegratedDataStores [Create] req: ", req)

	root := new(integratedDataStoreRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	// action	"plug_hardware_disk_device"

	return root.IntegratedDataStores, resp, err
}

// Delete -
func (s *IntegratedDataStoresServiceOp) Delete(ctx context.Context, resID int, id string, meta interface{}) (*Response, error) {
	if resID < 1 || id == "" {
		return nil, godo.NewArgError("resID or id", "cannot be empty or less than 1")
	}
	path := fmt.Sprintf(integratedDataStoresBasePath, resID)
	path = fmt.Sprintf("%s/%s%s", path, id, apiFormat)

	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("IntegratedDataStores [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit -
func (s *IntegratedDataStoresServiceOp) Edit(ctx context.Context, resID int, id string, editRequest *IntegratedDataStoresEditRequest) (*Response, error) {
	if resID < 1 || id == "" {
		return nil, godo.NewArgError("resID or id", "cannot be empty or less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("IntegratedDataStores [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf(integratedDataStoresBasePath, resID)
	path = fmt.Sprintf("%s/%s%s", path, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("IntegratedDataStores [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// BackendNodes get list of nodes from hypervisor
func (s *IntegratedDataStoresServiceOp) BackendNodes(ctx context.Context, id int) (*BackendNodes, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(backendNodesBasePath, id) + apiFormat
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("IntegratedDataStores [BackendNodes]  req: ", req)

	root := &BackendNodes{}
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
