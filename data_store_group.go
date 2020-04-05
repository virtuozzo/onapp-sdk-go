package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const dataStoreGroupsBasePath string = "settings/data_store_zones"

// DataStoreGroupsService is an interface for interfacing with the Data Store Zones
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/data-store-zones
type DataStoreGroupsService interface {
	List(context.Context, *ListOptions) ([]DataStoreGroup, *Response, error)
	Get(context.Context, int) (*DataStoreGroup, *Response, error)
	Create(context.Context, *DataStoreGroupCreateRequest) (*DataStoreGroup, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *DataStoreGroupEditRequest) (*Response, error)
}

// DataStoreGroupsServiceOp handles communication with the Data Store Groups related methods of the
// OnApp API.
type DataStoreGroupsServiceOp struct {
	client *Client
}

var _ DataStoreGroupsService = &DataStoreGroupsServiceOp{}

// DataStoreGroup represents a DataStoreGroup
type DataStoreGroup struct {
	ID                int                `json:"id,omitempty"`
	Label             string             `json:"label,omitempty"`
	CreatedAt         string             `json:"created_at,omitempty"`
	UpdatedAt         string             `json:"updated_at,omitempty"`
	ServerType        string             `json:"server_type,omitempty"`
	LocationGroupID   int                `json:"location_group_id,omitempty"`
	FederationEnabled bool               `json:"federation_enabled,bool"`
	FederationID      int                `json:"federation_id,omitempty"`
	Traded            bool               `json:"traded,bool"`
	Closed            bool               `json:"closed,bool"`
	HypervisorID      int                `json:"hypervisor_id,omitempty"`
	Identifier        string             `json:"identifier,omitempty"`
	PreconfiguredOnly bool               `json:"preconfigured_only,bool"`
	DraasID           int                `json:"draas_id,omitempty"`
	ProviderVdcID     int                `json:"provider_vdc_id,omitempty"`
	AdditionalFields  []AdditionalFields `json:"additional_fields,omitempty"`
	DatacenterID      int                `json:"datacenter_id,omitempty"`
	DefaultMaxIops    int                `json:"default_max_iops,omitempty"`
	DefaultBurstIops  int                `json:"default_burst_iops,omitempty"`
	MinDiskSize       int                `json:"min_disk_size,omitempty"`
}

// DataStoreGroupCreateRequest represents a request to create a DataStoreGroup
type DataStoreGroupCreateRequest struct {
	Label             string `json:"label,omitempty"`
	LocationGroupID   int    `json:"location_group_id,omitempty"`
	PreconfiguredOnly bool   `json:"preconfigured_only,bool"`
	ServerType        string `json:"server_type,omitempty"`
}

// DataStoreGroupEditRequest represents a request to edit a DataStoreGroup
type DataStoreGroupEditRequest struct {
	Label             string `json:"label,omitempty"`
	LocationGroupID   int    `json:"location_group_id,omitempty"`
	PreconfiguredOnly bool   `json:"preconfigured_only,bool"`
}

type dataStoreGroupCreateRequestRoot struct {
	DataStoreGroupCreateRequest *DataStoreGroupCreateRequest `json:"data_store_group"`
}

type dataStoreGroupRoot struct {
	DataStoreGroup *DataStoreGroup `json:"data_store_group"`
}

func (d DataStoreGroupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *DataStoreGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]DataStoreGroup, *Response, error) {
	path := dataStoreGroupsBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]DataStoreGroup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]DataStoreGroup, len(out))
	for i := range arr {
		arr[i] = out[i]["data_store_group"]
	}

	return arr, resp, err
}

// Get individual DataStoreGroup.
func (s *DataStoreGroupsServiceOp) Get(ctx context.Context, id int) (*DataStoreGroup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", dataStoreGroupsBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dataStoreGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DataStoreGroup, resp, err
}

// Create DataStoreGroup.
func (s *DataStoreGroupsServiceOp) Create(ctx context.Context, createRequest *DataStoreGroupCreateRequest) (*DataStoreGroup, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("DataStoreGroup createRequest", "cannot be nil")
	}

	path := dataStoreGroupsBasePath + apiFormat
	rootRequest := &dataStoreGroupCreateRequestRoot{
		DataStoreGroupCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("DataStoreGroup [Create] req: ", req)

	root := new(dataStoreGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DataStoreGroup, resp, err
}

// Delete DataStoreGroup.
func (s *DataStoreGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", dataStoreGroupsBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("DataStoreGroup [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit DataStoreGroup.
func (s *DataStoreGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *DataStoreGroupEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("DataStoreGroup [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", dataStoreGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("DataStoreGroup [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
