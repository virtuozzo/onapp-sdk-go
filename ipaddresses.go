package onappgo

import (
	"context"
	"fmt"
	// "log"
	"net/http"

	"github.com/digitalocean/godo"
)

const ipAddressesVSBasePath string = "virtual_machines/%d/ip_addresses"

const ipAddressesAssignUserBasePath string = "settings/networks/%d/ip_addresses/assign"
const ipAddressesUnassignUserBasePath string = "settings/networks/%d/ip_addresses/unassign"

// IPAddressesService is an interface for interfacing with the IPAddress
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ip-addresses
type IPAddressesService interface {
	List(context.Context, int, *ListOptions) ([]IPAddressJoin, *Response, error)

	// AssignVS(context.Context, int, *AssignIPAddress) (*IPAddress, *Response, error)
	// UnassingVS(context.Context, int) (*Response, error)

	// AssignUser(context.Context, int, *AssignIPAddress) (*IPAddress, *Response, error)
	// UnassingUser(context.Context, int) (*Response, error)
}

// IPAddressesServiceOp handles communication with the IPAddresses related methods of the
// OnApp API.
type IPAddressesServiceOp struct {
	client *Client
}

var _ IPAddressesService = &IPAddressesServiceOp{}

// IPAddress -
type IPAddress struct {
	Address         string `json:"address,omitempty"`
	Broadcast       string `json:"broadcast,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	ExternalAddress string `json:"external_address,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	HypervisorID    int    `json:"hypervisor_id,omitempty"`
	ID              int    `json:"id,omitempty"`
	IPNetID         int    `json:"ip_net_id,omitempty"`
	IPRangeID       int    `json:"ip_range_id,omitempty"`
	Ipv4            bool   `json:"ipv4"`
	LockVersion     int    `json:"lock_version,omitempty"`
	NetworkAddress  string `json:"network_address,omitempty"`
	NetworkID       int    `json:"network_id,omitempty"`
	Prefix          int    `json:"prefix,omitempty"`
	Pxe             bool   `json:"pxe"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
}

// IPAddressJoin -
type IPAddressJoin struct {
	CreatedAt          string    `json:"created_at"`
	ID                 int       `json:"id,omitempty"`
	IPAddress          IPAddress `json:"ip_address,omitempty"`
	IPAddressID        int       `json:"ip_address_id,omitempty"`
	NetworkInterfaceID int       `json:"network_interface_id,omitempty"`
	UpdatedAt          string    `json:"updated_at,omitempty"`
}

// IPAddressesJoin -
type IPAddressesJoin struct {
	IPAddressJoin []IPAddressJoin `json:"ip_address_join,omitempty"`
}

// IPAddresses -
type IPAddresses struct {
	IPAddress IPAddress `json:"ip_address,omitempty"`
}

// AssignIPAddress - used for assign IPAddress to the VirtualMachine or User
type AssignIPAddress struct {
	Address            string `json:"address,omitempty"`
	IPNetID            int    `json:"ip_net_id,omitempty"`
	IPRangeID          int    `json:"ip_range_id,omitempty"`
	IPVersion          int    `json:"ip_version,omitempty"`
	NetworkInterfaceID int    `json:"network_interface_id,omitempty"`
	OwnIP              int    `json:"own_ip,omitempty"`
	UsedIP             int    `json:"used_ip,omitempty"`
}

// List all IPAddress
func (s *IPAddressesServiceOp) List(ctx context.Context, id int, opt *ListOptions) ([]IPAddressJoin, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(ipAddressesVSBasePath, id) + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]IPAddressJoin
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]IPAddressJoin, len(out))
	for i := range arr {
		arr[i] = out[i]["ip_address_join"]
	}

	return arr, resp, err
}
