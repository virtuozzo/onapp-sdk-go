package onappgo

import (
	"context"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const licensesBasePath string = "settings/license"
const licensesEditBasePath string = "settings"

// LicensesService is an interface for interfacing with the Licenses
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/buckets/access-control
type LicensesService interface {
	Get(context.Context) (*License, *Response, error)
	Edit(context.Context, *LicenseEditRequest) (*Response, error)
}

// LicensesServiceOp handles communication with the License related methods of the
// OnApp API.
type LicensesServiceOp struct {
	client *Client
}

var _ LicensesService = &LicensesServiceOp{}

type License struct {
	IntegratedStorageLimit string `json:"integrated_storage_limit,omitempty"`
	Key                    string `json:"key,omitempty"`
	KvmXenCoreLimit        string `json:"kvm_xen_core_limit,omitempty"`
	KvmXenHvLimit          string `json:"kvm_xen_hv_limit,omitempty"`
	KvmXenVMLimit          string `json:"kvm_xen_vm_limit,omitempty"`
	Status                 string `json:"status,omitempty"`
	SupplierAllowed        bool   `json:"supplier_allowed,bool"`
	SupplierStatus         string `json:"supplier_status,omitempty"`
	TraderAllowed          bool   `json:"trader_allowed,bool"`
	TraderStatus           string `json:"trader_status,omitempty"`
	Type                   string `json:"type,omitempty"`
	Valid                  bool   `json:"valid,bool"`
	IsolatedLicense        bool   `json:"isolated_license,bool"`
	VcenterCoreLimit       string `json:"vcenter_core_limit,omitempty"`
	VcenterVMLimit         string `json:"vcenter_vm_limit,omitempty"`
	LicenseKey             string `json:"license_key,omitempty"`
}

type LicenseEditRequest struct {
	IsolatedLicense bool   `json:"isolated_license,bool"`
	LicenseKey      string `json:"license_key,omitempty"`
}

type licenseEditRequestRoot struct {
	LicenseEditRequest *LicenseEditRequest `json:"configuration"`
}

type licenseRoot struct {
	License *License `json:"license"`
}

// Get individual License.
func (s *LicensesServiceOp) Get(ctx context.Context) (*License, *Response, error) {
	path := licensesBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("License [Get] req: ", req)

	root := new(licenseRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.License, resp, err
}

// Edit individual License.
func (s *LicensesServiceOp) Edit(ctx context.Context, editRequest *LicenseEditRequest) (*Response, error) {
	if editRequest == nil {
		return nil, godo.NewArgError("License editRequest", "cannot be nil")
	}

	path := licensesEditBasePath + apiFormat

	rootRequest := &licenseEditRequestRoot{
		LicenseEditRequest: editRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, rootRequest)
	if err != nil {
		return nil, err
	}

	log.Println("License [Create/Edit] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// IsValid check if license is valid
func (res License) IsValid() bool {
	return res.Valid == true
}
