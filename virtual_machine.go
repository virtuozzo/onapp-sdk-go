package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const virtualMachineBasePath = "virtual_machines"

// VirtualMachinesService is an interface for interfacing with the VirtualMachine
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/virtual-servers
type VirtualMachinesService interface {
  List(context.Context, *ListOptions) ([]VirtualMachine, *Response, error)
  Get(context.Context, int) (*VirtualMachine, *Response, error)
  Create(context.Context, *VirtualMachineCreateRequest) (*VirtualMachine, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]VirtualMachine, *Response, error)

  Backups(context.Context, int, *ListOptions) ([]Backup, *Response, error)
  Transactions(context.Context, int, *ListOptions) ([]Transaction, *Response, error)
  Disks(context.Context, int, *ListOptions) ([]Disk, *Response, error)

  ListNetworkInterfaces(context.Context, int, *ListOptions) ([]NetworkInterface, *Response, error)
  ListFirewallRules(context.Context, int, *ListOptions) ([]FirewallRule, *Response, error)
}

// VirtualMachinesServiceOp handles communication with the VirtualMachine related methods of the
// OnApp API.
type VirtualMachinesServiceOp struct {
  client *Client
}

var _ VirtualMachinesService = &VirtualMachinesServiceOp{}

// VirtualMachine represent VirtualServer from OnApp API
type VirtualMachine struct {
  Acceleration                  bool                    `json:"acceleration,bool"`
  AccelerationAllowed           bool                    `json:"acceleration_allowed,bool"`
  AccelerationStatus            string                  `json:"acceleration_status,omitempty"`
  AddToMarketplace              string                  `json:"add_to_marketplace,omitempty"`
  AdminNote                     string                  `json:"admin_note,omitempty"`
  AllowedHotMigrate             bool                    `json:"allowed_hot_migrate,bool"`
  AllowedSwap                   bool                    `json:"allowed_swap,bool"`
  AutoscaleService              string                  `json:"autoscale_service,omitempty"`
  Booted                        bool                    `json:"booted,bool"`
  Built                         bool                    `json:"built,bool"`
  BuiltFromIso                  bool                    `json:"built_from_iso,bool"`
  BuiltFromOva                  bool                    `json:"built_from_ova,bool"`
  CDboot                        bool                    `json:"cdboot,bool"`
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
  EnableAutoscale               bool                    `json:"enable_autoscale,bool"`
  FirewallNotrack               bool                    `json:"firewall_notrack,bool"`
  Hostname                      string                  `json:"hostname,omitempty"`
  HotAddCPU                     string                  `json:"hot_add_cpu,omitempty"`
  HotAddMemory                  string                  `json:"hot_add_memory,omitempty"`
  HypervisorID                  int                     `json:"hypervisor_id,omitempty"`
  HypervisorType                string                  `json:"hypervisor_type,omitempty"`
  ID                            int                     `json:"id,omitempty"`
  Identifier                    string                  `json:"identifier,omitempty"`
  InitialRootPassword           string                  `json:"initial_root_password,omitempty"`
  InitialRootPasswordEncrypted  bool                    `json:"initial_root_password_encrypted,bool"`
  InstancePackageID             int                     `json:"instance_package_id,omitempty"`
  IPAddresses                   []map[string]IPAddress  `json:"ip_addresses,omitempty"`
  IsoID                         int                     `json:"iso_id,omitempty"`
  Label                         string                  `json:"label,omitempty"`
  LocalRemoteAccessIPAddress    string                  `json:"local_remote_access_ip_address,omitempty"`
  LocalRemoteAccessPort         int                     `json:"local_remote_access_port,omitempty"`
  Locked                        bool                    `json:"locked,bool"`
  Memory                        int                     `json:"memory,omitempty"`
  MinDiskSize                   int                     `json:"min_disk_size,omitempty"`
  MonthlyBandwidthUsed          float32                 `json:"monthly_bandwidth_used,omitempty"`
  Note                          string                  `json:"note,omitempty"`
  OpenstackID                   int                     `json:"openstack_id,omitempty"`
  OperatingSystem               string                  `json:"operating_system,omitempty"`
  OperatingSystemDistro         string                  `json:"operating_system_distro,omitempty"`
  // PreferredHVS                  []HVS                   `json:"preferred_hvs"`
  PricePerHour                  float32                 `json:"price_per_hour,omitempty"`
  PricePerHourPoweredOff        float32                 `json:"price_per_hour_powered_off,omitempty"`
  // Properties                    Propertie               `json:"properties"`
  RecoveryMode                  bool                    `json:"recovery_mode,bool"`
  RemoteAccessPassword          string                  `json:"remote_access_password,omitempty"`
  ServicePassword               string                  `json:"service_password,omitempty"`
  State                         string                  `json:"state,omitempty"`
  StorageServerType             string                  `json:"storage_server_type,omitempty"`
  StrictVirtualMachineID        int                     `json:"strict_virtual_machine_id,omitempty"`
  SupportIncrementalBackups     bool                    `json:"support_incremental_backups,bool"`
  Suspended                     bool                    `json:"suspended,bool"`
  TemplateID                    int                     `json:"template_id,omitempty"`
  TemplateLabel                 string                  `json:"template_label,omitempty"`
  TemplateVersion               string                  `json:"template_version,omitempty"`
  TimeZone                      string                  `json:"time_zone,omitempty"`
  TotalDiskSize                 int                     `json:"total_disk_size,omitempty"`
  UpdatedAt                     string                  `json:"updated_at,omitempty"`
  UserID                        int                     `json:"user_id,omitempty"`
  VappID                        int                     `json:"vapp_id,omitempty"`
  VcenterClusterID              int                     `json:"vcenter_cluster_id,omitempty"`
  VcenterMoref                  string                  `json:"vcenter_moref,omitempty"`
  VcenterReservedMemory         int                     `json:"vcenter_reserved_memory,omitempty"`
  Vip                           string                  `json:"vip,omitempty"`
  VmwareTools                   string                  `json:"vmware_tools,omitempty"`
  XenID                         int                     `json:"xen_id,omitempty"`

  // OnApp 6.1
  VirshConsole                  bool                    `json:"virsh_console,omitempty"`
}

// CustomRecipeVariableAttributes - 
type CustomRecipeVariableAttributes struct {
  Enabled         int    `json:"enabled,omitempty"`
  Name            string `json:"name,omitempty"`
  Value           string `json:"value,omitempty"`
}

// VirtualMachineCreateRequest represents a request to create a VirtualMachine
type VirtualMachineCreateRequest struct {
  CustomRecipeVariablesAttributes   []CustomRecipeVariableAttributes `json:"custom_recipe_variables_attributes,omitempty"`

  AccelerationAllowed               bool      `json:"acceleration_allowed,bool"`
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
  InstancePackageID                 int       `json:"instance_package_id,omitempty"`
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
  RequiredIPAddressAssignment       bool      `json:"required_ip_address_assignment,bool"`
  // *
  RequiredVirtualMachineBuild       bool      `json:"required_virtual_machine_build,bool"`
  RequiredVirtualMachineStartup     bool      `json:"required_virtual_machine_startup,bool"`
  SelectedIPAddress                 string    `json:"selected_ip_address,omitempty"`
  ServiceAddonIds                   []int     `json:"service_addon_ids,omitempty"`
  SwapDiskMinIops                   int       `json:"swap_disk_min_iops,omitempty"`
  // * in gigabytes 5, 10
  SwapDiskSize                      int       `json:"swap_disk_size,omitempty"`
  // *
  TemplateID                        int       `json:"template_id,omitempty"`
  TimeZone                          string    `json:"time_zone,omitempty"`
  TypeOfFormat                      string    `json:"type_of_format,omitempty"`

  // OnApp 6.1
  VirshConsole                      bool      `json:"virsh_console,omitempty"`
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

// List all VirtualMachines.
func (s *VirtualMachinesServiceOp) List(ctx context.Context, opt *ListOptions) ([]VirtualMachine, *Response, error) {
  path := virtualMachineBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]VirtualMachine
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  arr := make([]VirtualMachine, len(out))
  for i := range arr {
    arr[i] = out[i]["virtual_machine"]
  }

  return arr, resp, err
}

// Get individual VirtualMachine.
func (s *VirtualMachinesServiceOp) Get(ctx context.Context, id int) (*VirtualMachine, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", virtualMachineBasePath, id, apiFormat)

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

// Delete VirtualMachine.
func (s *VirtualMachinesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", virtualMachineBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, nil, err
  }

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  filter := struct{
    AssociatedObjectID    int
    AssociatedObjectType  string
  }{
    AssociatedObjectID    : id,
    AssociatedObjectType  : "VirtualMachine",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "VirtualMachine")
}

// Backups lists the backups for a VirtualMachine
func (s *VirtualMachinesServiceOp) Backups(ctx context.Context, id int, opt *ListOptions) ([]Backup, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  resourceType := "backup"
  path := fmt.Sprintf("%s/%d/%s%s", virtualMachineBasePath, id, resourceType+"s", apiFormat)
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
    backups[i] = out[i][resourceType]
  }

  return backups, resp, err
}

// Transactions lists the transactions for a VirtualMachine.
func (s *VirtualMachinesServiceOp) Transactions(ctx context.Context, id int, opt *ListOptions) ([]Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  resourceType := "transaction"
  path := fmt.Sprintf("%s/%d/%s%s", virtualMachineBasePath, id, resourceType+"s", apiFormat)
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
    transactions[i] = out[i][resourceType]
  }

  return transactions, resp, err
}

// Disks lists the disk for a VirtualMachine.
func (s *VirtualMachinesServiceOp) Disks(ctx context.Context, id int, opt *ListOptions) ([]Disk, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  resourceType := "disk"
  path := fmt.Sprintf("%s/%d/%s%s", virtualMachineBasePath, id, resourceType+"s", apiFormat)
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Disk
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  disks := make([]Disk, len(out))
  for i := range disks {
    disks[i] = out[i][resourceType]
  }

  return disks, resp, err
}

// ListNetworkInterfaces a VirtualMachine
func (s *VirtualMachinesServiceOp) ListNetworkInterfaces(ctx context.Context, id int, opt *ListOptions) ([]NetworkInterface, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  resourceType := "network_interface"
  path := fmt.Sprintf("%s/%d/%s%s", virtualMachineBasePath, id, resourceType+"s", apiFormat)
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]NetworkInterface
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  nets := make([]NetworkInterface, len(out))
  for i := range nets {
    nets[i] = out[i][resourceType]
  }

  return nets, resp, err
}

// ListFirewallRules a VirtualMachine
func (s *VirtualMachinesServiceOp) ListFirewallRules(ctx context.Context, id int, opt *ListOptions) ([]FirewallRule, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  resourceType := "firewall_rule"
  path := fmt.Sprintf("%s/%d/%s%s", virtualMachineBasePath, id, resourceType+"s", apiFormat)
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]FirewallRule
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  fwr := make([]FirewallRule, len(out))
  for i := range fwr {
    fwr[i] = out[i][resourceType]
  }

  return fwr, resp, err
}

// Debug - print formatted VirtualMachine structure
func (obj VirtualMachine) Debug() {
  fmt.Printf("                 ID: %d\n", obj.ID)
  fmt.Printf("         Identifier: %s\n", obj.Identifier)
  fmt.Printf("              Label: %s\n", obj.Label)
  fmt.Printf("InitialRootPassword: %s\n", obj.InitialRootPassword)
  fmt.Printf("      TemplateLabel: %s\n", obj.TemplateLabel)
  fmt.Printf("          CreatedAt: %s\n", obj.CreatedAt)
  fmt.Printf("              State: %s\n", obj.State)
  fmt.Printf("              Built: %t\n", obj.Built)
  fmt.Printf("             Booted: %t\n", obj.Booted)

  for i := range obj.IPAddresses {
    ip := obj.IPAddresses[i]["ip_address"]
    fmt.Printf("\t   IPAddresses: [%d]\n", i)
    ip.Debug()
    fmt.Println("")
  }
}
