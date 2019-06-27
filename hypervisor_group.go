package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "time"

  "github.com/digitalocean/godo"
)

const hypervisorGroupBasePath = "/settings/hypervisor_zones"

// HypervisorGroupsService is an interface for interfacing with the HypervisorGroup
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/compute-zones
type HypervisorGroupsService interface {
  List(context.Context, *ListOptions) ([]HypervisorGroup, *Response, error)
  Get(context.Context, int) (*HypervisorGroup, *Response, error)
  Create(context.Context, *HypervisorGroupCreateRequest) (*HypervisorGroup, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]HypervisorGroup, *Response, error)
}

// HypervisorZonesServiceOp handles communication with the HypervisorGroup related methods of the
// OnApp API.
type HypervisorZonesServiceOp struct {
  client *Client
}

var _ HypervisorGroupsService = &HypervisorZonesServiceOp{}

// HypervisorGroup represent VirtualServer from OnApp API
type HypervisorGroup struct {
  ID                          int           `json:"id,omitempty"`
  Label                       string        `json:"label,omitempty"`
  CreatedAt                   time.Time     `json:"created_at,omitempty"`
  UpdatedAt                   time.Time     `json:"updated_at,omitempty"`
  ServerType                  string        `json:"server_type,omitempty"`
  LocationGroupID             int           `json:"location_group_id,omitempty"`
  FederationEnabled           bool          `json:"federation_enabled,bool"`
  FederationID                int           `json:"federation_id,omitempty"`
  Traded                      bool          `json:"traded,bool"`
  Closed                      bool          `json:"closed,bool"`
  HypervisorID                int           `json:"hypervisor_id,omitempty"`
  Identifier                  string        `json:"identifier,omitempty"`
  DraasID                     int           `json:"draas_id,omitempty"`
  PreconfiguredOnly           bool          `json:"preconfigured_only,bool"`
  ProviderVdcID               int           `json:"provider_vdc_id,omitempty"`
  AdditionalFields            string        `json:"additional_fields,omitempty"`
  DatacenterID                int           `json:"datacenter_id,omitempty"`
  MaxHostFreeMemory           int           `json:"max_host_free_memory,omitempty"`
  MaxHostCPU                  int           `json:"max_host_cpu,omitempty"`
  PreferLocalReads            bool          `json:"prefer_local_reads,bool"`
  ReleaseResourceType         string        `json:"release_resource_type,omitempty"`
  NetworkFailure              bool          `json:"network_failure"`
  StorageChannel              int           `json:"storage_channel,omitempty"`
  RunSysprep                  bool          `json:"run_sysprep,bool"`
  RecoveryType                string        `json:"recovery_type,omitempty"`
  FailoverTimeout             int           `json:"failover_timeout,omitempty"`
  CPUUnits                    int           `json:"cpu_units,omitempty"`
  SupplierVersion             string        `json:"supplier_version,omitempty"`
  SupplierProvider            string        `json:"supplier_provider,omitempty"`
  ProviderName                string        `json:"provider_name,omitempty"`
  ScheduledForDeletion        string        `json:"scheduled_for_deletion,omitempty"`
  CPUFlagsEnabled             bool          `json:"cpu_flags_enabled,bool"`
  CPUFlags                    []string      `json:"cpu_flags,omitempty"`
  Tier                        string        `json:"tier,omitempty"`
  SupportsVirtualServerMotion string        `json:"supports_virtual_server_motion,omitempty"`
  CustomConfig                string        `json:"custom_config,omitempty"`
}

// HypervisorGroupCreateRequest represents a request to create a HypervisorGroup
type HypervisorGroupCreateRequest struct {
  CPUFlagsEnabled     bool   `json:"cpu_flags_enabled,bool"`
  CPUUnits            string `json:"cpu_units,omitempty"`
  CustomConfig        string `json:"custom_config,omitempty"`
  FailoverTimeout     int    `json:"failover_timeout,omitempty"`
  Label               string `json:"label,omitempty"`
  LocationGroupID     int    `json:"location_group_id,omitempty"`
  MaxVmsStartAtOnce   int    `json:"max_vms_start_at_once,omitempty"`
  PreconfiguredOnly   bool   `json:"preconfigured_only,bool"`
  RecoveryType        string `json:"recovery_type,omitempty"`
  ReleaseResourceType string `json:"release_resource_type,omitempty"`
  RunSysprep          bool   `json:"run_sysprep,bool"`
  ServerType          string `json:"server_type,omitempty"`

  // VMware parameters:
  DefaultGateway      string `json:"default_gateway,omitempty"`
  Vlan                string `json:"vlan,omitempty"`
}

type hypervisorGroupCreateRequestRoot struct {
  HypervisorGroupCreateRequest  *HypervisorGroupCreateRequest  `json:"hypervisor_group"`
}

type hypervisorGroupRoot struct {
  HypervisorGroup  *HypervisorGroup  `json:"hypervisor_group"`
}

func (d HypervisorGroupCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all HypervisorGroup.
func (s *HypervisorZonesServiceOp) List(ctx context.Context, opt *ListOptions) ([]HypervisorGroup, *Response, error) {
  path := hypervisorGroupBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]HypervisorGroup
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  vms := make([]HypervisorGroup, len(out))
  for i := range vms {
    vms[i] = out[i]["hypervisor_group"]
  }

  return vms, resp, err
}

// Get individual HypervisorGroup.
func (s *HypervisorZonesServiceOp) Get(ctx context.Context, id int) (*HypervisorGroup, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", hypervisorGroupBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(hypervisorGroupRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.HypervisorGroup, resp, err
}

// Create HypervisorGroup.
func (s *HypervisorZonesServiceOp) Create(ctx context.Context, createRequest *HypervisorGroupCreateRequest) (*HypervisorGroup, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("HypervisorZone createRequest", "cannot be nil")
  }

  path := hypervisorGroupBasePath + apiFormat

  rootRequest := &hypervisorGroupCreateRequestRoot{
    HypervisorGroupCreateRequest : createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nHypervisorZone [Create] req: ", req)

  root := new(hypervisorGroupRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.HypervisorGroup, resp, err
}

// Delete HypervisorGroup.
func (s *HypervisorZonesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", hypervisorGroupBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }

  resp, err := s.client.Do(ctx, req, nil)

  opt := &ListOptions{
    PerPage : searchTransactions,
  }

  trxVM, resp, err := s.client.Transactions.ListByGroup(ctx, id, "HypervisorGroup", opt)

  var root *Transaction
  e := trxVM.Front()
  if e != nil {
    val := e.Value.(Transaction)
    root = &val
    return root, resp, err
  }

  return nil, nil, err
}

// Debug - print formatted HypervisorGroup structure
func (h HypervisorGroup) Debug() {
  fmt.Println("[                  ID]: ", h.ID)
  fmt.Println("[          Identifier]: ", h.Identifier)
}
