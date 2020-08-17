package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const softwareLicenseBasePath string = "/software_licenses"

// SoftwareLicensesService is an interface for interfacing with the SoftwareLicense
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/software-licenses
type SoftwareLicensesService interface {
	List(context.Context, *ListOptions) ([]SoftwareLicense, *Response, error)
	Get(context.Context, int) (*SoftwareLicense, *Response, error)
	Create(context.Context, *SoftwareLicenseCreateRequest) (*SoftwareLicense, *Response, error)
	Delete(context.Context, int) (*Response, error)
	Edit(context.Context, int, *SoftwareLicenseEditRequest) (*Response, error)
}

// SoftwareLicensesServiceOp handles communication with the SoftwareLicense related methods of the
// OnApp API.
type SoftwareLicensesServiceOp struct {
	client *Client
}

var _ SoftwareLicensesService = &SoftwareLicensesServiceOp{}

// SoftwareLicense - represent disk from Virtual Machine
type SoftwareLicense struct {
	ID        int      `json:"id,omitempty"`
	Arch      string   `json:"arch,omitempty"`
	Edition   []string `json:"edition,omitempty"`
	Tail      string   `json:"tail,omitempty"`
	Distro    string   `json:"distro,omitempty"`
	License   string   `json:"license,omitempty"`
	Total     int      `json:"total,omitempty"`
	Count     int      `json:"count,omitempty"`
	CreatedAt string   `json:"created_at,omitempty"`
	UpdatedAt string   `json:"updated_at,omitempty"`
}

// SoftwareLicenseCreateRequest - data for creating SoftwareLicense
type SoftwareLicenseCreateRequest struct {
	Arch      string   `json:"arch,omitempty"`
	Edition   []string `json:"edition,omitempty"`
	Tail      string   `json:"tail,omitempty"`
	Distro    string   `json:"distro,omitempty"`
	License   string   `json:"license,omitempty"`
	Total     int      `json:"total,omitempty"`
	Count     int      `json:"count,omitempty"`
}

// SoftwareLicenseEditRequest - data for editing SoftwareLicense
type SoftwareLicenseEditRequest SoftwareLicenseCreateRequest

type softwareLicenseCreateRequestRoot struct {
	SoftwareLicenseCreateRequest *SoftwareLicenseCreateRequest `json:"software_license"`
}

type softwareLicenseRoot struct {
	SoftwareLicense *SoftwareLicense `json:"software_license"`
}

func (d SoftwareLicenseCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Software License
func (s *SoftwareLicensesServiceOp) List(ctx context.Context, opt *ListOptions) ([]SoftwareLicense, *Response, error) {
	path := softwareLicenseBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]SoftwareLicense
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]SoftwareLicense, len(out))
	for i := range arr {
		arr[i] = out[i]["software_license"]
	}

	return arr, resp, err
}

// Get individual Software License
func (s *SoftwareLicensesServiceOp) Get(ctx context.Context, id int) (*SoftwareLicense, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", softwareLicenseBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(softwareLicenseRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SoftwareLicense, resp, err
}

// Create Software License
func (s *SoftwareLicensesServiceOp) Create(ctx context.Context, createRequest *SoftwareLicenseCreateRequest) (*SoftwareLicense, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := softwareLicenseBasePath + apiFormat

	rootRequest := &softwareLicenseCreateRequestRoot{
		SoftwareLicenseCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("SoftwareLicense [Create]  req: ", req)

	root := new(softwareLicenseRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SoftwareLicense, resp, err
}

// Delete Software License
func (s *SoftwareLicensesServiceOp) Delete(ctx context.Context, id int) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", softwareLicenseBasePath, id, apiFormat)

	path, err := addOptions(path, nil)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("SoftwareLicense [Delete]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit Software License
func (s *SoftwareLicensesServiceOp) Edit(ctx context.Context, id int, editRequest *SoftwareLicenseEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", softwareLicenseBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("SoftwareLicense [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
