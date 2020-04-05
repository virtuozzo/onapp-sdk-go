package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const networksBasePath string = "settings/networks"

// NetworksService is an interface for interfacing with the Network
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/networks
type NetworksService interface {
	List(context.Context, *ListOptions) ([]Network, *Response, error)
	Get(context.Context, int) (*Network, *Response, error)
	Create(context.Context, *NetworkCreateRequest) (*Network, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *NetworkEditRequest) (*Response, error)
}

// NetworksServiceOp handles communication with the Networks related methods of the
// OnApp API.
type NetworksServiceOp struct {
	client *Client
}

var _ NetworksService = &NetworksServiceOp{}

// Network represents a Network
type Network struct {
	ID                        int    `json:"id,omitempty"`
	Label                     string `json:"label,omitempty"`
	Identifier                string `json:"identifier,omitempty"`
	CreatedAt                 string `json:"created_at,omitempty"`
	UpdatedAt                 string `json:"updated_at,omitempty"`
	Vlan                      int    `json:"vlan,omitempty"`
	NetworkGroupID            int    `json:"network_group_id,omitempty"`
	Type                      string `json:"type,omitempty"`
	UserID                    int    `json:"user_id,omitempty"`
	IPAddressPoolID           int    `json:"ip_address_pool_id,omitempty"`
	DefaultOutsideIPAddressID int    `json:"default_outside_ip_address_id,omitempty"`
	DefaultNatRuleNumber      int    `json:"default_nat_rule_number,omitempty"`
	PrefixSize                int    `json:"prefix_size,omitempty"`
	IsNated                   bool   `json:"is_nated,bool"`
	VappID                    int    `json:"vapp_id,omitempty"`
	VdcID                     int    `json:"vdc_id,omitempty"`
	Enabled                   bool   `json:"enabled,bool"`
	Gateway                   string `json:"gateway,omitempty"`
	Netmask                   string `json:"netmask,omitempty"`
	PrimaryDNS                string `json:"primary_dns,omitempty"`
	SecondaryDNS              string `json:"secondary_dns,omitempty"`
	DNSSuffix                 string `json:"dns_suffix,omitempty"`
	Shared                    bool   `json:"shared,bool"`
	FenceMode                 string `json:"fence_mode,omitempty"`
	VcenterIdentifier         string `json:"vcenter_identifier,omitempty"`
	ParentNetworkID           int    `json:"parent_network_id,omitempty"`
	OpenstackID               int    `json:"openstack_id,omitempty"`
	DvSwitchID                int    `json:"dv_switch_id,omitempty"`

	// OnApp 6.2
	VdcGroupID        int `json:"vdc_group_id,omitempty"`
	UniversalRouterID int `json:"universal_router_id,omitempty"`
	AssignedUser      int `json:"assigned_user,omitempty"`
}

// NetworkCreateRequest represents a request to create a Network
type NetworkCreateRequest struct {
	Label          string `json:"label,omitempty"`
	NetworkGroupID int    `json:"network_group_id,omitempty"`
	Vlan           int    `json:"vlan,omitempty"`

	// Must be set as default value: "Networking::Network"
	Type string `json:"type,omitempty"`
}

// NetworkEditRequest represents a request to edit a Network
type NetworkEditRequest struct {
	Label          string `json:"label,omitempty"`
	NetworkGroupID int    `json:"network_group_id,omitempty"`
	Vlan           int    `json:"vlan,omitempty"`
}

type networkCreateRequestRoot struct {
	NetworkCreateRequest *NetworkCreateRequest `json:"network"`
}

type networkRoot struct {
	Network *Network `json:"network"`
}

func (d NetworkCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Networks.
func (s *NetworksServiceOp) List(ctx context.Context, opt *ListOptions) ([]Network, *Response, error) {
	path := networksBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]Network
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Network, len(out))
	for i := range arr {
		arr[i] = out[i]["network"]
	}

	return arr, resp, err
}

// Get individual Network.
func (s *NetworksServiceOp) Get(ctx context.Context, id int) (*Network, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", networksBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(networkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Network, resp, err
}

// Create Network.
func (s *NetworksServiceOp) Create(ctx context.Context, createRequest *NetworkCreateRequest) (*Network, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("Network createRequest", "cannot be nil")
	}

	path := networksBasePath + apiFormat
	rootRequest := &networkCreateRequestRoot{
		NetworkCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("Network [Create] req: ", req)

	root := new(networkRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Network, resp, err
}

// Delete Network.
func (s *NetworksServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", networksBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Network [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit Network.
func (s *NetworksServiceOp) Edit(ctx context.Context, id int, editRequest *NetworkEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("Network [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", networksBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Network [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
