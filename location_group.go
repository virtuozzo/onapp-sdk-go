package onappgo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/godo"
)

const locationGroupsBasePath string = "settings/location_groups"
const locationGroupRefreshBasePath string = locationGroupsBasePath + "/refresh"

const locationGroupAttachHypervisorGroupsBasePath string = locationGroupsBasePath + "/%d/hypervisor_groups/attach_resource"
const locationGroupDetachHypervisorGroupsBasePath string = locationGroupsBasePath + "/%d/hypervisor_groups/%d/detach_resource"

const locationGroupAttachDataStoreGroupsBasePath string = locationGroupsBasePath + "/%d/data_store_groups/attach_resource"
const locationGroupDetachDataStoreGroupsBasePath string = locationGroupsBasePath + "/%d/data_store_groups/%d/detach_resource"

const locationGroupAttachNetworkGroupsBasePath string = locationGroupsBasePath + "/%d/network_groups/attach_resource"
const locationGroupDetachNetworkGroupsBasePath string = locationGroupsBasePath + "/%d/network_groups/%d/detach_resource"

const locationGroupAttachBackupServerGroupBasePath string = locationGroupsBasePath + "/%d/backup_server_groups/attach_resource"
const locationGroupDetachBackupServerGroupBasePath string = locationGroupsBasePath + "/%d/backup_server_groups/%d/detach_resource"

// LocationGroupsService is an interface for interfacing with the LocationGroup
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/location-groups
type LocationGroupsService interface {
	List(context.Context, *ListOptions) ([]LocationGroup, *Response, error)
	Get(context.Context, int) (*LocationGroup, *Response, error)

	Refresh(context.Context) (*Response, error)
}

// LocationGroupsServiceOp handles communication with the LocationGroup related methods of the
// OnApp API.
type LocationGroupsServiceOp struct {
	client *Client
}

var _ LocationGroupsService = &LocationGroupsServiceOp{}

// LocationGroup represent LocationGroup from OnApp API
type LocationGroup struct {
	ID           int     `json:"id,omitempty"`
	CreatedAt    string  `json:"created_at,omitempty"`
	UpdatedAt    string  `json:"updated_at,omitempty"`
	Country      string  `json:"country,omitempty"`
	City         string  `json:"city,omitempty"`
	FederationID int     `json:"federation_id,omitempty"`
	Lat          float64 `json:"lat,omitempty"`
	Lng          float64 `json:"lng,omitempty"`
	CdnEnabled   bool    `json:"cdn_enabled,bool"`
	Federated    bool    `json:"federated,bool"`
}

// LocationGroupCreateRequest represents a request to create a LocationGroup
type LocationGroupCreateRequest struct {
}

type locationGroupCreateRequestRoot struct {
	LocationGroupCreateRequest *LocationGroupCreateRequest `json:"location_group"`
}

type locationGroupRoot struct {
	LocationGroup *LocationGroup `json:"location_group"`
}

func (d LocationGroupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all LocationGroups.
func (s *LocationGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]LocationGroup, *Response, error) {
	path := locationGroupsBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]LocationGroup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]LocationGroup, len(out))
	for i := range arr {
		arr[i] = out[i]["location_group"]
	}

	return arr, resp, err
}

// Get individual LocationGroup.
func (s *LocationGroupsServiceOp) Get(ctx context.Context, id int) (*LocationGroup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", locationGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(locationGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.LocationGroup, resp, err
}

// Refresh LocationGroup.
func (s *LocationGroupsServiceOp) Refresh(ctx context.Context) (*Response, error) {
	path := locationGroupRefreshBasePath + apiFormat
	path, err := addOptions(path, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("LocationGroup [Refresh] req: ", req)

	return s.client.Do(ctx, req, nil)
}
