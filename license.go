package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const licensesBasePath      string = "settings/license"
const licensesEditBasePath  string = "settings"

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
  Type                   string `json:"type,omitempty"`
  Key                    string `json:"key,omitempty"`
  Valid                  bool   `json:"valid,bool"`
  Status                 string `json:"status,omitempty"`
  KvmXenHvLimit          string `json:"kvm_xen_hv_limit,omitempty"`
  KvmXenVMLimit          string `json:"kvm_xen_vm_limit,omitempty"`
  VcenterVMLimit         string `json:"vcenter_vm_limit,omitempty"`
  KvmXenCoreLimit        string `json:"kvm_xen_core_limit,omitempty"`
  VcenterCoreLimit       string `json:"vcenter_core_limit,omitempty"`
  IntegratedStorageLimit string `json:"integrated_storage_limit,omitempty"`
  TraderStatus           string `json:"trader_status,omitempty"`
  TraderAllowed          bool   `json:"trader_allowed,bool"`
  SupplierStatus         string `json:"supplier_status,omitempty"`
  SupplierAllowed        bool   `json:"supplier_allowed,bool"`
}

type LicenseEditRequest struct {
  IsolatedLicense bool    `json:"isolated_license,bool"`
  LicenseKey      string  `json:"license_key,omitempty"`
}

type licenseEditRequestRoot struct {
  LicenseEditRequest  *LicenseEditRequest  `json:"configuration"`
}

type licenseRoot struct {
  License  *License  `json:"license"`
}

// Get individual License.
func (s *LicensesServiceOp) Get(ctx context.Context) (*License, *Response, error) {
  path := fmt.Sprintf("%s%s", licensesBasePath, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

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

  path := fmt.Sprintf("%s%s", licensesEditBasePath, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
  if err != nil {
    return nil, err
  }

  rootRequest := &licenseEditRequestRoot{
    LicenseEditRequest : editRequest,
  }

  return s.client.Do(ctx, req, rootRequest)
}
