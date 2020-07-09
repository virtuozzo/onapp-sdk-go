package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const cloudBootIPAddressesBasePath string = "cloud_boot_ip_addresses"

// CloudbootIPAddressesService is an interface for interfacing with the CloudbootIPAddress
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/cloudboot-ip-addresses
type CloudbootIPAddressesService interface {
	List(context.Context, *ListOptions) ([]CloudbootIPAddress, *Response, error)
	// Get(context.Context, int) (*CloudbootIPAddress, *Response, error)
	Create(context.Context, *CloudbootIPAddressCreateRequest) (*CloudbootIPAddress, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
}

// CloudbootIPAddressesServiceOp handles communication with the CloudbootIPAddress related methods of the
// OnApp API.
type CloudbootIPAddressesServiceOp struct {
	client *Client
}

var _ CloudbootIPAddressesService = &CloudbootIPAddressesServiceOp{}

type CloudbootIPAddress struct {
	ID              int    `json:"id,omitempty"`
	Address         string `json:"address,omitempty"`
	Broadcast       string `json:"broadcast,omitempty"`
	NetworkAddress  string `json:"network_address,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	Pxe             bool   `json:"pxe,bool"`
	HypervisorID    int    `json:"hypervisor_id,omitempty"`
	IPRangeID       int    `json:"ip_range_id,omitempty"`
	ExternalAddress string `json:"external_address,omitempty"`
	Free            bool   `json:"free,bool"`
	Netmask         string `json:"netmask,omitempty"`
}

type CloudbootIPAddressCreateRequest struct {
	Address string `json:"address,omitempty"`
}

type cloudbootIPAddressRoot struct {
	CloudbootIPAddress *CloudbootIPAddress `json:"ip_address"`
}

type cloudbootIPAddressCreateRequestRoot struct {
	CloudbootIPAddressCreateRequest *CloudbootIPAddressCreateRequest `json:"ip_address"`
}

func (d CloudbootIPAddressCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Cloudboot CloudbootIPAddresss
func (s *CloudbootIPAddressesServiceOp) List(ctx context.Context, opt *ListOptions) ([]CloudbootIPAddress, *Response, error) {
	path := cloudBootIPAddressesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]CloudbootIPAddress
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]CloudbootIPAddress, len(out))
	for i := range arr {
		arr[i] = out[i]["ip_address"]
	}

	return arr, resp, err
}

// // Get individual Cloudboot CloudbootIPAddress
// func (s *CloudbootIPAddressesServiceOp) Get(ctx context.Context, id int) (*CloudbootIPAddress, *Response, error) {
// 	if id < 1 {
// 		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
// 	}

// 	path := fmt.Sprintf("%s/%d%s", cloudBootIPAddressesBasePath, id, apiFormat)
// 	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	root := new(cloudbootIPAddressRoot)
// 	resp, err := s.client.Do(ctx, req, root)
// 	if err != nil {
// 		return nil, resp, err
// 	}

// 	return root.CloudbootIPAddress, resp, err
// }

// Create Cloudboot CloudbootIPAddress
func (s *CloudbootIPAddressesServiceOp) Create(ctx context.Context, createRequest *CloudbootIPAddressCreateRequest) (*CloudbootIPAddress, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("CloudbootIPAddress createRequest", "cannot be nil")
	}

	path := cloudBootIPAddressesBasePath + apiFormat
	rootRequest := &cloudbootIPAddressCreateRequestRoot{
		CloudbootIPAddressCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("CloudbootIPAddress [Create] req: ", req)

	root := new(cloudbootIPAddressRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.CloudbootIPAddress, resp, err
}

// Delete Cloudboot CloudbootIPAddress
func (s *CloudbootIPAddressesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", cloudBootIPAddressesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("CloudbootIPAddress [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}
