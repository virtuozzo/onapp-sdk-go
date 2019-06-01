package onappgo

import (
  "context"
  "errors"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const virtualMachineBasePath = "virtual_machines"

var errNoNetworks = errors.New("no networks have been defined")

// VirtualMachinesService is an interface for interfacing with the VirtualMachine
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/virtual-servers/get-list-of-vss
type VirtualMachinesService interface {
  List(context.Context, *ListOptions) ([]VirtualMachine, *Response, error)
  // ListByTag(context.Context, string, *ListOptions) ([]VirtualMachine, *Response, error)
  Get(context.Context, string) (*VirtualMachine, *Response, error)
  Create(context.Context, *VirtualMachineCreateRequest) (*VirtualMachine, *Response, error)
  // CreateMultiple(context.Context, *VirtualMachineMultiCreateRequest) ([]VirtualMachine, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  // DeleteByTag(context.Context, string) (*Response, error)
  // Snapshots(context.Context, int, *ListOptions) ([]Image, *Response, error)
  // Backups(context.Context, int, *ListOptions) ([]Image, *Response, error)
  // Actions(context.Context, int, *ListOptions) ([]Action, *Response, error)
  // Neighbors(context.Context, int) ([]VirtualMachine, *Response, error)
}

// VirtualMachinesServiceOp handles communication with the VirtualMachine related methods of the
// OnApp API.
type VirtualMachinesServiceOp struct {
  client *Client
}

// var _ VirtualMachinesService = &VirtualMachinesServiceOp{}

// VirtualMachine represent VirtualServer from OnApp API
type VirtualMachine struct {
  Acceleration                  bool                    `json:"acceleration,bool,omitempty"`
  AccelerationAllowed           bool                    `json:"acceleration_allowed,bool,omitempty"`
  AccelerationStatus            string                  `json:"acceleration_status,omitempty"`
  AddToMarketplace              string                  `json:"add_to_marketplace,omitempty"`
  AdminNote                     string                  `json:"admin_note,omitempty"`
  AllowedHotMigrate             bool                    `json:"allowed_hot_migrate,bool,omitempty"`
  AllowedSwap                   bool                    `json:"allowed_swap,bool,omitempty,bool,omitempty"`
  AutoscaleService              string                  `json:"autoscale_service,omitempty"`
  Booted                        bool                    `json:"booted,bool,omitempty"`
  Built                         bool                    `json:"built,bool,omitempty"`
  BuiltFromIso                  bool                    `json:"built_from_iso,bool,omitempty"`
  BuiltFromOva                  bool                    `json:"built_from_ova,bool,omitempty"`
  CDboot                        bool                    `json:"cdboot,bool,omitempty"`
  CoresPerSocket                int                     `json:"cores_per_socket,omitempty"`
  CPUPriority                   int                     `json:"cpu_priority,omitempty"`
  CPUShares                     int                     `json:"cpu_shares,omitempty"`
  CPUSockets                    string                  `json:"cpu_sockets,omitempty"`
  CPUUnits                      int                     `json:"cpu_units,omitempty"`
  Cpus                          int                     `json:"cpus,omitempty"`
  CreatedAt                     string                  `json:"created_at,omitempty"`
  DeletedAt                     string                  `json:"deleted_at,omitempty"`
  Domain                        string                  `json:"domain,omitempty"`
  DraasKeys                     []string                `json:"draas_keys,omitempty"`
  DraasMode                     int                     `json:"draas_mode,omitempty"`
  EdgeServerType                string                  `json:"edge_server_type,omitempty"`
  EnableAutoscale               bool                    `json:"enable_autoscale,bool,omitempty"`
  FirewallNotrack               bool                    `json:"firewall_notrack,bool,omitempty"`
  Hostname                      string                  `json:"hostname,omitempty"`
  HotAddCPU                     string                  `json:"hot_add_cpu,omitempty"`
  HotAddMemory                  string                  `json:"hot_add_memory,omitempty"`
  HypervisorID                  int                     `json:"hypervisor_id,omitempty"`
  HypervisorType                string                  `json:"hypervisor_type,omitempty"`
  ID                            int                     `json:"id,omitempty"`
  Identifier                    string                  `json:"identifier,omitempty"`
  InitialRootPassword           string                  `json:"initial_root_password,omitempty"`
  InitialRootPasswordEncrypted  bool                    `json:"initial_root_password_encrypted,bool,omitempty"`
  InstancePackageID             string                  `json:"instance_package_id,omitempty"`
  IPAddressesRaw                []map[string]IPAddress  `json:"ip_addresses,omitempty"`
  IsoID                         string                  `json:"iso_id,omitempty"`
  Label                         string                  `json:"label,omitempty"`
  LocalRemoteAccessIPAddress    string                  `json:"local_remote_access_ip_address,omitempty"`
  LocalRemoteAccessPort         int                     `json:"local_remote_access_port,omitempty"`
  Locked                        bool                    `json:"locked,bool,omitempty"`
  Memory                        int                     `json:"memory,omitempty"`
  MinDiskSize                   int                     `json:"min_disk_size,omitempty"`
  MonthlyBandwidthUsed          float32                 `json:"monthly_bandwidth_used,omitempty"`
  Note                          string                  `json:"note,omitempty"`
  OpenstackID                   string                  `json:"openstack_id,omitempty"`
  OperatingSystem               string                  `json:"operating_system,omitempty"`
  OperatingSystemDistro         string                  `json:"operating_system_distro,omitempty"`
  // PreferredHVS                  []HVS                   `json:"preferred_hvs"`
  PricePerHour                  float32                 `json:"price_per_hour,omitempty"`
  PricePerHourPoweredOff        float32                 `json:"price_per_hour_powered_off,omitempty"`
  // Properties                    []string                `json:"properties"`
  RecoveryMode                  bool                    `json:"recovery_mode,bool,omitempty"`
  RemoteAccessPassword          string                  `json:"remote_access_password,omitempty"`
  ServicePassword               string                  `json:"service_password,omitempty"`
  State                         string                  `json:"state,omitempty"`
  StorageServerType             string                  `json:"storage_server_type,omitempty"`
  StrictVirtualMachineID        string                  `json:"strict_virtual_machine_id,omitempty"`
  SupportIncrementalBackups     bool                    `json:"support_incremental_backups,bool,omitempty"`
  Suspended                     bool                    `json:"suspended,bool,omitempty"`
  TemplateID                    int                     `json:"template_id,omitempty"`
  TemplateLabel                 string                  `json:"template_label,omitempty"`
  TemplateVersion               string                  `json:"template_version,omitempty"`
  TimeZone                      string                  `json:"time_zone,omitempty"`
  TotalDiskSize                 int                     `json:"total_disk_size,omitempty"`
  UpdatedAt                     string                  `json:"updated_at,omitempty"`
  UserID                        int                     `json:"user_id,omitempty"`
  VappID                        string                  `json:"vapp_id,omitempty"`
  VcenterClusterID              string                  `json:"vcenter_cluster_id,omitempty"`
  VcenterMoref                  string                  `json:"vcenter_moref,omitempty"`
  VcenterReservedMemory         int                     `json:"vcenter_reserved_memory,omitempty"`
  Vip                           string                  `json:"vip,omitempty"`
  VmwareTools                   string                  `json:"vmware_tools,omitempty"`
  XenID                         int                     `json:"xen_id,omitempty"`
}

// IPAddress - represents an ip address of VirtualMachine
type IPAddress struct {
  Address         string    `json:"address,omitempty"`
  Broadcast       string    `json:"broadcast,omitempty"`
  CreatedAt       string    `json:"created_at,omitempty"`
  Free            bool      `json:"free,bool,omitempty"`
  Gateway         string    `json:"gateway,omitempty"`
  HypervisorID    string    `json:"hypervisor_id,omitempty"`
  ID              int       `json:"id,omitempty"`
  IPRangeID       int       `json:"ip_range_id,omitempty"`
  Netmask         string    `json:"netmask,omitempty"`
  NetworkAddress  string    `json:"network_address,omitempty"`
  Pxe             bool      `json:"pxe,bool,omitempty"`
  UpdatedAt       string    `json:"updated_at,omitempty"`
  UserID          string    `json:"user_id,omitempty"`
}

// VirtualMachineCreateRequest represents a request to create a VirtualMachine.
type VirtualMachineCreateRequest struct {
  // location_group_id
  // type_of_format
  // network_id
  // initial_root_password_encryption_key
  // custom_variables_attributes
        // [
        //   enabled - true, if the variable is enabled, otherwise false
        //   id - variable ID
        //   name - variable name
        //   value - variable value script
        // ]
  // service_addon_ids

  AccelerationAllowed             bool      `json:"acceleration_allowed,bool,omitempty"`
  AdminNote                       string    `json:"admin_note,omitempty"`
  CPUShares                       int       `json:"cpu_shares,omitempty"`
  CPUSockets                      string    `json:"cpu_sockets,omitempty"`
  Cpus                            int       `json:"cpus,omitempty"`
  Domain                          string    `json:"domain,omitempty"`
  DataStoreGroupPrimaryID         int       `json:"data_store_group_primary_id,omitempty"`
  DataStoreGroupSwapID            int       `json:"data_store_group_swap_id,omitempty"`
  EnableAutoscale                 int       `json:"enable_autoscale,omitempty"`
  Hostname                        string    `json:"hostname,omitempty"`
  HypervisorGroupID               int       `json:"hypervisor_group_id,omitempty"`
  HypervisorID                    int       `json:"hypervisor_id,omitempty"`
  InitialRootPassword             string    `json:"initial_root_password,omitempty"`
  InstancePackageID               string    `json:"instance_package_id,omitempty"`
  Label                           string    `json:"label,omitempty"`
  LicensingKey                    string    `json:"licensing_key,omitempty"`
  LicensingServerID               int       `json:"licensing_server_id,omitempty"`
  LicensingType                   string    `json:"licensing_type,omitempty"`
  LocationGroupID                 int       `json:"location_group_id,omitempty"`
  Memory                          int       `json:"memory,omitempty"`
  PrimaryDiskMinIops              int       `json:"primary_disk_min_iops,omitempty"`
  PrimaryDiskSize                 int       `json:"primary_disk_size,omitempty"`
  PrimaryNetworkGroupID           int       `json:"primary_network_group_id,omitempty"`
  RateLimit                       int       `json:"rate_limit,omitempty"`
  RecipeJoinsAttributes           []string  `json:"recipe_joins_attributes,omitempty"`
  RequiredAutomaticBackup         int       `json:"required_automatic_backup,omitempty"`
  RequiredIPAddressAssignment     int       `json:"required_ip_address_assignment,omitempty"`
  RequiredVirtualMachineBuild     int       `json:"required_virtual_machine_build,omitempty"`
  RequiredVirtualMachineStartup   int       `json:"required_virtual_machine_startup,omitempty"`
  SelectedIPAddress               string    `json:"selected_ip_address,omitempty"`
  SwapDiskMinIops                 int       `json:"swap_disk_min_iops,omitempty"`
  SwapDiskSize                    int       `json:"swap_disk_size,omitempty"`
  TemplateID                      int       `json:"template_id,omitempty"`
  TimeZone                        string    `json:"time_zone,omitempty"`
}

// VirtualMachineMultiCreateRequest is a request to create multiple VirtualMachine.
// type VirtualMachineMultiCreateRequest struct {
// }

// vmRoot represents a VirtualMachine root
type vmRoot struct {
  VirtualMachine *VirtualMachine
}

func (d VirtualMachineCreateRequest) String() string {
  return godo.Stringify(d)
}

// func (d VirtualMachineMultiCreateRequest) String() string {
//   return godo.Stringify(d)
// }

// Performs a list request given a path.
func (s *VirtualMachinesServiceOp) list(ctx context.Context, path string) ([]VirtualMachine, *Response, error) {
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]VirtualMachine
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  vms := make([]VirtualMachine, len(out))
  for i := range vms {
    vms[i] = out[i]["virtual_machine"]
  }

  return vms, resp, err
}

// List all VirtualMachines.
func (s *VirtualMachinesServiceOp) List(ctx context.Context, opt *ListOptions) ([]VirtualMachine, *Response, error) {
  path := virtualMachineBasePath
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  return s.list(ctx, path)
}

// Get individual VirtualMachine.
func (s *VirtualMachinesServiceOp) Get(ctx context.Context, VirtualMachineID string) (*VirtualMachine, *Response, error) {
  if len(VirtualMachineID) < 1 {
    return nil, nil, godo.NewArgError("VirtualMachineID", "cannot be empty")
  }

  path := fmt.Sprintf("%s/%s", virtualMachineBasePath, VirtualMachineID)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nreq: ", req)

  var out map[string]VirtualMachine
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  vm := out["virtual_machine"]

  return &vm, resp, err
}

// Create VirtualMachine.
func (s *VirtualMachinesServiceOp) Create(ctx context.Context, createRequest *VirtualMachineCreateRequest) (*VirtualMachine, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := virtualMachineBasePath

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("[Create] req: ", req)

  // out := new(virtualMachineRoot)
  // resp, err := s.client.Do(ctx, req, out)
  // if err != nil {
  //  return nil, resp, err
  // }

  // return out, resp, err
  return nil, nil, err
}

// CreateMultiple - creates multiple VirtualMachines.
// func (s *VirtualMachinesServiceOp) CreateMultiple(ctx context.Context, createRequest *VirtualMachineMultiCreateRequest) ([]VirtualMachine, *Response, error) {
//   if createRequest == nil {
//     return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
//   }

//   path := virtualMachineBasePath

//   req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
//   if err != nil {
//     return nil, nil, err
//   }

//   fmt.Println("[CreateMultiple] req: ", req)

//   // out := new(virtualMachinesRoot)
//   // resp, err := s.client.Do(ctx, req, out)
//   // if err != nil {
//   //  return nil, resp, err

//   // return out, resp, err
//   return nil, nil, err
// }
