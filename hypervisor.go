package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "time"

  "github.com/digitalocean/godo"
)

// Xen/KVM, VMware
const hypervisorBasePath = "/settings/hypervisors"

// CloudBoot, Smart CloudBoot, Baremetal CloudBoot
const cloudBootHypervisorBasePath = "/settings/assets/%s/hypervisors"


// HypervisorsService is an interface for interfacing with the Hypervisor
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/compute-resources
type HypervisorsService interface {
  List(context.Context, *ListOptions) ([]Hypervisor, *Response, error)
  Get(context.Context, int) (*Hypervisor, *Response, error)
  Create(context.Context, *HypervisorCreateRequest) (*Hypervisor, *Response, error)
  // Delete(context.Context, int) (*Response, error)
  Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
  // Edit(context.Context, int, *ListOptions) ([]Hypervisor, *Response, error)
}

// HypervisorsServiceOp handles communication with the Hypervisor related methods of the
// OnApp API.
type HypervisorsServiceOp struct {
  client *Client
}

var _ HypervisorsService = &HypervisorsServiceOp{}

// ConnectionOptions for VMware hypervisor
type ConnectionOptions struct {
  APIURL                       string `json:"api_url,omitempty"`
  Login                        string `json:"login,omitempty"`
  Password                     string `json:"password,omitempty"`
  ClusterName                  string `json:"cluster_name,omitempty"`
  DistributedVirtualSwitchName string `json:"distributed_virtual_switch_name,omitempty"`
}

// Disks of hypervisor
type Disks struct {
  Scsi     string   `json:"scsi,omitempty"`
  Selected bool    `json:"selected,bool"`
}

// Nics of hypervisor
type Nics struct {
  Mac  string  `json:"mac,omitempty"`
  Type int     `json:"type"`
}

// CustomPCIs of hypervisor
type CustomPCIs struct {
  Pci      string  `json:"pci,omitempty"`
  Selected bool    `json:"selected,bool"`
}

// Storage for Hypervisor
type Storage struct {
  Disks      []Disks        `json:"disks,omitempty"`
  Nics       []Nics         `json:"nics,omitempty"`
  CustomPCIs []CustomPCIs   `json:"custom_pcis,omitempty"`
}

// Hypervisor represent Hypervisor from OnApp API
type Hypervisor struct {
  ID                               int               `json:"id,omitempty"`
  Label                            string            `json:"label,omitempty"`
  IPAddress                        string            `json:"ip_address,omitempty"`
  CreatedAt                        time.Time         `json:"created_at,omitempty"`
  UpdatedAt                        time.Time         `json:"updated_at,omitempty"`
  Locked                           bool              `json:"locked"`
  CalledInAt                       string            `json:"called_in_at,omitempty"`
  Online                           bool              `json:"online"`
  Spare                            bool              `json:"spare"`
  FailureCount                     int               `json:"failure_count,omitempty"`
  HypervisorType                   string            `json:"hypervisor_type,omitempty"`
  HypervisorGroupID                int               `json:"hypervisor_group_id,omitempty"`
  Enabled                          bool              `json:"enabled"`
  Uptime                           string            `json:"uptime,omitempty"`
  ListOfZombieDomains              string            `json:"list_of_zombie_domains,omitempty"`
  ListOfVolumeGroups               string            `json:"list_of_volume_groups,omitempty"`
  Host                             string            `json:"host,omitempty"`
  Release                          string            `json:"release,omitempty"`
  Machine                          string            `json:"machine,omitempty"`
  CPUMhz                           string            `json:"cpu_mhz,omitempty"`
  Cpus                             int               `json:"cpus,omitempty"`
  ThreadsPerCore                   int               `json:"threads_per_core,omitempty"`
  TotalMem                         int               `json:"total_mem,omitempty"`
  TotalZombieMem                   int               `json:"total_zombie_mem,omitempty"`
  DisableFailover                  bool              `json:"disable_failover"`
  Mac                              string            `json:"mac,omitempty"`
  CustomConfig                     string            `json:"custom_config,omitempty"`
  FormatDisks                      bool              `json:"format_disks"`
  PassthroughDisks                 bool              `json:"passthrough_disks"`
  ConnectionOptions                ConnectionOptions `json:"connection_options,omitempty"`
  HostID                           string            `json:"host_id,omitempty"`
  FreeMem                          int               `json:"free_mem,omitempty"`
  BackupIPAddress                  string            `json:"backup_ip_address,omitempty"`
  Built                            bool              `json:"built"`
  Blocked                          bool              `json:"blocked"`
  ServerType                       string            `json:"server_type,omitempty"`
  Backup                           bool              `json:"backup"`
  CPUIdle                          int               `json:"cpu_idle,omitempty"`
  Mtu                              int               `json:"mtu,omitempty"`
  MemInfo                          int               `json:"mem_info,omitempty"`
  StorageControllerMemorySize      int               `json:"storage_controller_memory_size,omitempty"`
  DisksPerStorageController        int               `json:"disks_per_storage_controller,omitempty"`
  CloudBootOs                      string            `json:"cloud_boot_os,omitempty"`
  AllowUnsafeAssignedInterrupts    bool              `json:"allow_unsafe_assigned_interrupts"`
  Dom0MemorySize                   int               `json:"dom0_memory_size,omitempty"`
  CPUCores                         int               `json:"cpu_cores,omitempty"`
  CPUUnits                         int               `json:"cpu_units,omitempty"`
  PowerCycleCommand                string            `json:"power_cycle_command,omitempty"`
  Rebooting                        bool              `json:"rebooting"`
  MaintenanceMode                  bool              `json:"maintenance_mode"`
  CPUFlags                         []string          `json:"cpu_flags,omitempty"`
  AmqpExchangeName                 string            `json:"amqp_exchange_name,omitempty"`
  CacheMirrors                     int               `json:"cache_mirrors,omitempty"`
  CacheStripes                     int               `json:"cache_stripes,omitempty"`
  StorageControllerDbSize          int               `json:"storage_controller_db_size,omitempty"`
  StorageBondingMode               string            `json:"storage_bonding_mode,omitempty"`
  OsVersion                        int               `json:"os_version,omitempty"`
  OsVersionMinor                   int               `json:"os_version_minor,omitempty"`
  IntegratedStorageDisabled        bool              `json:"integrated_storage_disabled"`
  StorageVlan                      string            `json:"storage_vlan,omitempty"`
  ApplyHypervisorGroupCustomConfig bool              `json:"apply_hypervisor_group_custom_config"`
  CPUModel                         string            `json:"cpu_model,omitempty"`
  SegregationOsType                string            `json:"segregation_os_type,omitempty"`
  CrashDebug                       bool              `json:"crash_debug"`
  FailoverRecipeID                 string            `json:"failover_recipe_id,omitempty"`
  TotalCpus                        int               `json:"total_cpus,omitempty"`
  FreeMemory                       int               `json:"free_memory,omitempty"`
  UsedCPUResources                 int               `json:"used_cpu_resources,omitempty"`
  TotalMemory                      int               `json:"total_memory,omitempty"`
  FreeDiskSpace                    map[string]int    `json:"free_disk_space,omitempty"`
  MemoryAllocatedByRunningVms      int               `json:"memory_allocated_by_running_vms,omitempty"`
  TotalMemoryAllocatedByVms        int               `json:"total_memory_allocated_by_vms,omitempty"`
  Storage                          Storage           `json:"storage,omitempty"`
  ListOfLogicalVolumes             string            `json:"list_of_logical_volumes,omitempty"`
  Distro                           string            `json:"distro,omitempty"`
}

// HypervisorCreateRequest represents a request to create a default XEN/KVM Hypervisor
type HypervisorCreateRequest struct {
  Label                         string `json:"label,omitempty"`

  // VMware
  IPAddress                     string `json:"ip_address,omitempty"`

  // CloudBoot, SmartCloudBoot, VMware
  BackupIPAddress               string `json:"backup_ip_address,omitempty"`
  CollectStats                  bool   `json:"collect_stats,bool"`
  DisableFailover               bool   `json:"disable_failover,bool"`

  HypervisorType                string `json:"hypervisor_type,omitempty"`
  SegregationOsType             string `json:"segregation_os_type,omitempty"`
  Enabled                       bool   `json:"enabled,bool"`

  // BaremetalCloudBoot
  FailoverRecipeID              int    `json:"failover_recipe_id,omitempty"`

  HypervisorGroupID             int    `json:"hypervisor_group_id,omitempty"`
  CPUUnits                      int    `json:"cpu_units,omitempty"`

  // SmartCloudBoot, BaremetalCloudBoot
  PxeIPAddressID                int     `json:"pxe_ip_address_id,omitempty"`

  // CloudBoot, SmartCloudBoot, BaremetalCloudBoot
  ServerType                    string  `json:"server_type,omitempty"`
  Backup                        bool    `json:"backup,bool"`

  // only for VMware
  ConnectionOptions             ConnectionOptions `json:"connection_options,omitempty"`

  // CloudBoot, SmartCloudBoot
  FormatDisks                   bool    `json:"format_disks,bool"`
  PassthroughDisks              bool    `json:"passthrough_disks,bool"`
  Storage                       Storage `json:"storage,omitempty"`
  Mtu                           int     `json:"mtu,omitempty"`
  StorageControllerMemorySize   int     `json:"storage_controller_memory_size,omitempty"`
  DisksPerStorageController     int     `json:"disks_per_storage_controller,omitempty"`
  CustomConfig                  string  `json:"custom_config,omitempty"`

  // only for CloudBoot
  CloudBootOs                   string  `json:"cloud_boot_os,omitempty"`

  // only for SmartCloudBoot
  PassthroughCustomPcis         string  `json:"passthrough_custom_pcis,omitempty"`
  AllowUnsafeAssignedInterrupts bool    `json:"allow_unsafe_assigned_interrupts,bool"`
}

// CloudBootCreateRequest represents a request to create a CloudBoot hypervisor
type CloudBootCreateRequest struct {
  // Label                       string  `json:"label,omitempty"`
  // HypervisorType              string  `json:"hypervisor_type,omitempty"`
  // SegregationOsType           string  `json:"segregation_os_type,omitempty"`
  // ServerType                  string  `json:"server_type,omitempty"`
  // Backup                      bool    `json:"backup"`
  // BackupIPAddress             string  `json:"backup_ip_address,omitempty"`
  // Enabled                     bool    `json:"enabled"`
  // CollectStats                bool    `json:"collect_stats"`
  // DisableFailover             bool    `json:"disable_failover"`
  // FormatDisks                 bool    `json:"format_disks"`
  // PassthroughDisks            bool    `json:"passthrough_disks"`
  // Storage                     Storage `json:"storage,omitempty"`
  // Mtu                         int     `json:"mtu,omitempty"`
  // StorageControllerMemorySize int     `json:"storage_controller_memory_size,omitempty"`
  // DisksPerStorageController   int     `json:"disks_per_storage_controller,omitempty"`
  // CloudBootOs                 string  `json:"cloud_boot_os,omitempty"`
  // CustomConfig                string  `json:"custom_config,omitempty"`
}

// SmartCloudBootCreateRequest represents a request to create a Smart CloudBoot hypervisor
type SmartCloudBootCreateRequest struct {
  // Label                         string  `json:"label,omitempty"`
  // PxeIPAddressID                int     `json:"pxe_ip_address_id,omitempty"`
  // HypervisorType                string  `json:"hypervisor_type,omitempty"`
  // SegregationOsType             string  `json:"segregation_os_type,omitempty"`
  // ServerType                    string  `json:"server_type,omitempty"`
  // BackupIPAddress               string  `json:"backup_ip_address,omitempty"`
  // Enabled                       bool    `json:"enabled"`
  // CollectStats                  bool    `json:"collect_stats"`
  // DisableFailover               bool    `json:"disable_failover"`
  // FormatDisks                   bool    `json:"format_disks"`
  // PassthroughDisks              bool    `json:"passthrough_disks"`
  // Storage                       Storage `json:"storage,omitempty"`
  // PassthroughCustomPcis         string  `json:"passthrough_custom_pcis,omitempty"`
  // Mtu                           int     `json:"mtu,omitempty"`
  // StorageControllerMemorySize   int     `json:"storage_controller_memory_size,omitempty"`
  // DisksPerStorageController     int     `json:"disks_per_storage_controller,omitempty"`
  // AllowUnsafeAssignedInterrupts bool    `json:"allow_unsafe_assigned_interrupts"`
  // CustomConfig                  string  `json:"custom_config,omitempty"`
}

// BaremetalCloudBootCreateRequest represents a request to create a Baremetal CloudBoot hypervisor
type BaremetalCloudBootCreateRequest struct {
  // Label            string   `json:"label,omitempty"`
  // PxeIPAddressID   int      `json:"pxe_ip_address_id,omitempty"`
  // HypervisorType   string   `json:"hypervisor_type,omitempty"`
  // ServerType       string   `json:"server_type,omitempty"`
  // Enabled          bool     `json:"enabled"`
  // FailoverRecipeID int      `json:"failover_recipe_id,omitempty"`
}

// VMwareCreateRequest represents a request to create a VMware hypervisor
type VMwareCreateRequest struct {
  // Label             string            `json:"label,omitempty"`
  // IPAddress         string            `json:"ip_address,omitempty"`
  // BackupIPAddress   string            `json:"backup_ip_address,omitempty"`
  // HypervisorType    string            `json:"hypervisor_type,omitempty"`
  // SegregationOsType string            `json:"segregation_os_type,omitempty"`
  // Enabled           bool              `json:"enabled"`
  // CollectStats      bool              `json:"collect_stats"`
  // DisableFailover   bool              `json:"disable_failover"`
  // ConnectionOptions ConnectionOptions `json:"connection_options,omitempty"`
}

type hypervisorCreateRequestRoot struct {
  HypervisorCreateRequest  *HypervisorCreateRequest  `json:"hypervisor"`
}

type hypervisorRoot struct {
  Hypervisor  *Hypervisor  `json:"hypervisor"`
}

func (d HypervisorCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all Hypervisors.
func (s *HypervisorsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Hypervisor, *Response, error) {
  path := hypervisorBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Hypervisor
  resp, err := s.client.Do(ctx, req, &out)

  if err != nil {
    return nil, resp, err
  }

  vms := make([]Hypervisor, len(out))
  for i := range vms {
    vms[i] = out[i]["hypervisor"]
  }

  return vms, resp, err
}

// Get individual Hypervisor.
func (s *HypervisorsServiceOp) Get(ctx context.Context, id int) (*Hypervisor, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", hypervisorBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(hypervisorRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Hypervisor, resp, err
}

// Create Hypervisor.
func (s *HypervisorsServiceOp) Create(ctx context.Context, createRequest *HypervisorCreateRequest) (*Hypervisor, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("Hypervisor createRequest", "cannot be nil")
  }

  // if hypervisorType == "xen" || hypervisorType == "kvm" || hypervisorType == "vmware" {
  //   path := hypervisorBasePath + apiFormat
  // } else {
  //   path := cloudBootHypervisorBasePath + apiFormat
  // }
  path := hypervisorBasePath + apiFormat

  rootRequest := &hypervisorCreateRequestRoot{
    HypervisorCreateRequest : createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }

  fmt.Println("\nHypervisor [Create] req: ", req)

  root := new(hypervisorRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.Hypervisor, resp, err
}

// Delete Hypervisor.
func (s *HypervisorsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", hypervisorBasePath, id, apiFormat)
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

  trxVM, resp, err := s.client.Transactions.ListByGroup(ctx, id, "Hypervisor", opt)

  var root *Transaction
  e := trxVM.Front()
  if e != nil {
    val := e.Value.(Transaction)
    root = &val
    return root, resp, err
  }

  return nil, nil, err
}

// Debug - print formatted Hypervisor structure
func (h Hypervisor) Debug() {
  fmt.Println("[        ID]: ", h.ID)
  fmt.Println("[     Label]: ", h.Label)
  fmt.Println("[ IPAddress]: ", h.IPAddress)
  fmt.Println("[      Host]: ", h.Host)
  fmt.Println("[       Mac]: ", h.Mac)
  fmt.Println("[    Online]: ", h.Online)
  fmt.Println("[ TotalCpus]: ", h.TotalCpus)
  fmt.Println("[    Uptime]: ", h.Uptime)
}
