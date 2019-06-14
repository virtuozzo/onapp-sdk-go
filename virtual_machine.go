package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "time"

  "github.com/digitalocean/godo"
)

const virtualMachineBasePath = "virtual_machines"

// VirtualMachinesService is an interface for interfacing with the VirtualMachine
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/virtual-servers/get-list-of-vss
type VirtualMachinesService interface {
  List(context.Context, *ListOptions) ([]VirtualMachine, *Response, error)
  Get(context.Context, int) (*VirtualMachine, *Response, error)
  Create(context.Context, *VirtualMachineCreateRequest) (*VirtualMachine, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  // Snapshots(context.Context, int, *ListOptions) ([]Image, *Response, error)
  Backups(context.Context, int, *ListOptions) ([]Backup, *Response, error)
  Transactions(context.Context, int, *ListOptions) ([]Transaction, *Response, error)
}

// VirtualMachinesServiceOp handles communication with the VirtualMachine related methods of the
// OnApp API.
type VirtualMachinesServiceOp struct {
  client *Client
}

var _ VirtualMachinesService = &VirtualMachinesServiceOp{}

// VirtualMachine represent VirtualServer from OnApp API
type VirtualMachine struct {
  Acceleration                  bool                    `json:"acceleration,bool,omitempty"`
  AccelerationAllowed           bool                    `json:"acceleration_allowed,bool,omitempty"`
  AccelerationStatus            string                  `json:"acceleration_status,omitempty"`
  AddToMarketplace              string                  `json:"add_to_marketplace,omitempty"`
  AdminNote                     string                  `json:"admin_note,omitempty"`
  AllowedHotMigrate             bool                    `json:"allowed_hot_migrate,bool,omitempty"`
  AllowedSwap                   bool                    `json:"allowed_swap,bool,omitempty"`
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
  CreatedAt                     time.Time               `json:"created_at,omitempty"`
  DeletedAt                     time.Time               `json:"deleted_at,omitempty"`
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
  IPAddresses                   []IPAddress             `json:"ip_addresses,omitempty"`
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
  // Properties                    Propertie               `json:"properties"`
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
  UpdatedAt                     time.Time               `json:"updated_at,omitempty"`
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
  Address         string      `json:"address,omitempty"`
  Broadcast       string      `json:"broadcast,omitempty"`
  CreatedAt       time.Time   `json:"created_at,omitempty"`
  Free            bool        `json:"free,bool,omitempty"`
  Gateway         string      `json:"gateway,omitempty"`
  HypervisorID    string      `json:"hypervisor_id,omitempty"`
  ID              int         `json:"id,omitempty"`
  IPRangeID       int         `json:"ip_range_id,omitempty"`
  Netmask         string      `json:"netmask,omitempty"`
  NetworkAddress  string      `json:"network_address,omitempty"`
  Pxe             bool        `json:"pxe,bool,omitempty"`
  UpdatedAt       time.Time   `json:"updated_at,omitempty"`
  UserID          string      `json:"user_id,omitempty"`
}

// VirtualMachineCreateRequest represents a request to create a VirtualMachine
type VirtualMachineCreateRequest struct {
  // custom_variables_attributes
      // [
      //   enabled - true, if the variable is enabled, otherwise false
      //   id - variable ID
      //   name - variable name
      //   value - variable value script
      // ]
  // service_addon_ids

  AccelerationAllowed               bool      `json:"acceleration_allowed,bool,omitempty"`
  AdminNote                         string    `json:"admin_note,omitempty"`
  // *
  CPUShares                         int       `json:"cpu_shares,omitempty"`
  CPUSockets                        string    `json:"cpu_sockets,omitempty"`
  // *
  Cpus                              int       `json:"cpus,omitempty"`
  Domain                            string    `json:"domain,omitempty"`
  DataStoreGroupPrimaryID           int       `json:"data_store_group_primary_id,omitempty"`
  DataStoreGroupSwapID              int       `json:"data_store_group_swap_id,omitempty"`
  EnableAutoscale                   int       `json:"enable_autoscale,omitempty"`
  // *
  Hostname                          string    `json:"hostname,omitempty"`
  HypervisorGroupID                 int       `json:"hypervisor_group_id,omitempty"`
  HypervisorID                      int       `json:"hypervisor_id,omitempty"`
  InitialRootPassword               string    `json:"initial_root_password,omitempty"`
  InitialRootPasswordEncryptionKey  string    `json:"initial_root_password_encryption_key,omitempty"`
  InstancePackageID                 string    `json:"instance_package_id,omitempty"`
  // *
  Label                             string    `json:"label,omitempty"`
  // *
  LicensingKey                      string    `json:"licensing_key,omitempty"`
  LicensingServerID                 int       `json:"licensing_server_id,omitempty"`
  // *
  LicensingType                     string    `json:"licensing_type,omitempty"`
  LocationGroupID                   int       `json:"location_group_id,omitempty"`
  // * in megabytes 1024, 2048
  Memory                            int       `json:"memory,omitempty"`
  NetworkID                         int       `json:"network_id,omitempty"`
  PrimaryDiskMinIops                int       `json:"primary_disk_min_iops,omitempty"`
  // * in gigabytes 5, 10, 20
  PrimaryDiskSize                   int       `json:"primary_disk_size,omitempty"`
  PrimaryNetworkGroupID             int       `json:"primary_network_group_id,omitempty"`
  RateLimit                         int       `json:"rate_limit,omitempty"`
  RecipeJoinsAttributes             []string  `json:"recipe_joins_attributes,omitempty"`
  RequiredAutomaticBackup           int       `json:"required_automatic_backup,omitempty"`
  // *
  RequiredIPAddressAssignment       bool      `json:"required_ip_address_assignment,bool,omitempty"`
  // *
  RequiredVirtualMachineBuild       bool      `json:"required_virtual_machine_build,bool,omitempty"`
  RequiredVirtualMachineStartup     bool      `json:"required_virtual_machine_startup,bool,omitempty"`
  SelectedIPAddress                 string    `json:"selected_ip_address,omitempty"`
  SwapDiskMinIops                   int       `json:"swap_disk_min_iops,omitempty"`
  // * in gigabytes 5, 10
  SwapDiskSize                      int       `json:"swap_disk_size,omitempty"`
  // *
  TemplateID                        int       `json:"template_id,omitempty"`
  TimeZone                          string    `json:"time_zone,omitempty"`
  TypeOfFormat                      string    `json:"type_of_format,omitempty"`
}

type virtualMachineCreateRequestRoot struct {
  VirtualMachineCreateRequest  *VirtualMachineCreateRequest  `json:"virtual_machine"`
}

type virtualMachineRoot struct {
  VirtualMachine  *VirtualMachine  `json:"virtual_machine"`
}

func (d VirtualMachineCreateRequest) String() string {
  return godo.Stringify(d)
}

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
func (s *VirtualMachinesServiceOp) Get(ctx context.Context, vmID int) (*VirtualMachine, *Response, error) {
  if vmID < 1 {
    return nil, nil, godo.NewArgError("vmID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", virtualMachineBasePath, vmID, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(virtualMachineRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.VirtualMachine, resp, err
}

// Create VirtualMachine.
func (s *VirtualMachinesServiceOp) Create(ctx context.Context, createRequest *VirtualMachineCreateRequest) (*VirtualMachine, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := virtualMachineBasePath + apiFormat

  rootRequest := &virtualMachineCreateRequestRoot{
    VirtualMachineCreateRequest : createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\n[Create]  req: ", req)

  root := new(virtualMachineRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.VirtualMachine, resp, err
}

// Performs a delete request given a path
// func (s *VirtualMachinesServiceOp) delete(ctx context.Context, path string) (*Response, error) {
//   req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
//   if err != nil {
//     return nil, err
//   }

//   resp, err := s.client.Do(ctx, req, nil)

//   return resp, err
// }

// Delete VirtualMachine.
// func (s *VirtualMachinesServiceOp) Delete(ctx context.Context, virtualMachineID int) (*Response, error) {
//   if virtualMachineID < 1 {
//     return nil, godo.NewArgError("virtualMachineID", "cannot be less than 1")
//   }

//   path := fmt.Sprintf("%s/%d%s", virtualMachineBasePath, virtualMachineID, apiFormat)

//   return s.delete(ctx, path)
// }

// Performs a delete request given a path
func (s *VirtualMachinesServiceOp) delete(ctx context.Context, path string) (*Response, error) {
  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }

  fmt.Printf("delete.req: %s\n",  req.URL)
  resp, err := s.client.Do(ctx, req, nil)

  return resp, err
}

// Delete VirtualMachine.
func (s *VirtualMachinesServiceOp) Delete(ctx context.Context, virtualMachineID int, meta interface{}) (*Response, error) {
  if virtualMachineID < 1 {
    return nil, godo.NewArgError("virtualMachineID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", virtualMachineBasePath, virtualMachineID, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  return s.delete(ctx, path)
}

// Backups lists the backups for a VirtualMachine
func (s *VirtualMachinesServiceOp) Backups(ctx context.Context, virtualMachineID int, opt *ListOptions) ([]Backup, *Response, error) {
  if virtualMachineID < 1 {
    return nil, nil, godo.NewArgError("virtualMachineID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d/backups%s", virtualMachineBasePath, virtualMachineID, apiFormat)
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Backup
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  backups := make([]Backup, len(out))
  for i := range backups {
    backups[i] = out[i]["backup"]
  }

  return backups, resp, err
}

// Transactions lists the transactions for a VirtualMachine.
func (s *VirtualMachinesServiceOp) Transactions(ctx context.Context, virtualMachineID int, opt *ListOptions) ([]Transaction, *Response, error) {
  if virtualMachineID < 1 {
    return nil, nil, godo.NewArgError("virtualMachineID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d/transactions%s", virtualMachineBasePath, virtualMachineID, apiFormat)
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Transaction
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  transactions := make([]Transaction, len(out))
  for i := range transactions {
    transactions[i] = out[i]["transaction"]
  }

  return transactions, resp, err
}
