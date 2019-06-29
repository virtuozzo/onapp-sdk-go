package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "time"

  "github.com/digitalocean/godo"
)

// Xen/KVM, VMware - CRUD
// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Get, Delete
const hypervisorBasePath = "settings/hypervisors"

// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Create, Edit
const cloudBootHypervisorBasePath = "settings/assets/%s/hypervisors"

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

// Hypervisor represent Hypervisor of the OnApp API
type Hypervisor struct {
  ID                               int               `json:"id,omitempty"`
  Label                            string            `json:"label,omitempty"`
  IPAddress                        string            `json:"ip_address,omitempty"`
  CreatedAt                        time.Time         `json:"created_at,omitempty"`
  UpdatedAt                        time.Time         `json:"updated_at,omitempty"`
  Locked                           bool              `json:"locked,bool"`
  CalledInAt                       string            `json:"called_in_at,omitempty"`
  Online                           bool              `json:"online,bool"`
  Spare                            bool              `json:"spare,bool"`
  FailureCount                     int               `json:"failure_count,omitempty"`
  HypervisorType                   string            `json:"hypervisor_type,omitempty"`
  HypervisorGroupID                int               `json:"hypervisor_group_id,omitempty"`
  Enabled                          bool              `json:"enabled,bool"`
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
  DisableFailover                  bool              `json:"disable_failover,bool"`
  Mac                              string            `json:"mac,omitempty"`
  CustomConfig                     string            `json:"custom_config,omitempty"`
  FormatDisks                      bool              `json:"format_disks,bool"`
  PassthroughDisks                 bool              `json:"passthrough_disks,bool"`
  ConnectionOptions                ConnectionOptions `json:"connection_options,omitempty"`
  HostID                           string            `json:"host_id,omitempty"`
  FreeMem                          int               `json:"free_mem,omitempty"`
  BackupIPAddress                  string            `json:"backup_ip_address,omitempty"`
  Built                            bool              `json:"built,bool"`
  Blocked                          bool              `json:"blocked,bool"`
  ServerType                       string            `json:"server_type,omitempty"`
  Backup                           bool              `json:"backup,bool"`
  CPUIdle                          int               `json:"cpu_idle,omitempty"`
  Mtu                              int               `json:"mtu,omitempty"`
  MemInfo                          int               `json:"mem_info,omitempty"`
  StorageControllerMemorySize      int               `json:"storage_controller_memory_size,omitempty"`
  DisksPerStorageController        int               `json:"disks_per_storage_controller,omitempty"`
  CloudBootOs                      string            `json:"cloud_boot_os,omitempty"`
  AllowUnsafeAssignedInterrupts    bool              `json:"allow_unsafe_assigned_interrupts,bool"`
  Dom0MemorySize                   int               `json:"dom0_memory_size,omitempty"`
  CPUCores                         int               `json:"cpu_cores,omitempty"`
  CPUUnits                         int               `json:"cpu_units,omitempty"`
  PowerCycleCommand                string            `json:"power_cycle_command,omitempty"`
  Rebooting                        bool              `json:"rebooting,bool"`
  MaintenanceMode                  bool              `json:"maintenance_mode,bool"`
  CPUFlags                         []string          `json:"cpu_flags,omitempty"`
  AmqpExchangeName                 string            `json:"amqp_exchange_name,omitempty"`
  CacheMirrors                     int               `json:"cache_mirrors,omitempty"`
  CacheStripes                     int               `json:"cache_stripes,omitempty"`
  StorageControllerDbSize          int               `json:"storage_controller_db_size,omitempty"`
  StorageBondingMode               string            `json:"storage_bonding_mode,omitempty"`
  OsVersion                        int               `json:"os_version,omitempty"`
  OsVersionMinor                   int               `json:"os_version_minor,omitempty"`
  IntegratedStorageDisabled        bool              `json:"integrated_storage_disabled,bool"`
  StorageVlan                      string            `json:"storage_vlan,omitempty"`
  ApplyHypervisorGroupCustomConfig bool              `json:"apply_hypervisor_group_custom_config,bool"`
  CPUModel                         string            `json:"cpu_model,omitempty"`
  SegregationOsType                string            `json:"segregation_os_type,omitempty"`
  CrashDebug                       bool              `json:"crash_debug,bool"`
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

// HypervisorCreateRequest represents a request to create a Hypervisor
type HypervisorCreateRequest struct {
  Label                         string  `json:"label,omitempty"`

  // VMware
  IPAddress                     string  `json:"ip_address,omitempty"`

  // CloudBoot, SmartCloudBoot, VMware
  BackupIPAddress               string  `json:"backup_ip_address,omitempty"`
  CollectStats                  bool    `json:"collect_stats,bool"`
  DisableFailover               bool    `json:"disable_failover,bool"`

  // SmartCloudBoot only can be: kvm
  // VMware, CloudBoot: xen, kvm
  // BaremetalCloudBoot only can be: xen
  HypervisorType                string  `json:"hypervisor_type,omitempty"`
  SegregationOsType             string  `json:"segregation_os_type,omitempty"`
  Enabled                       bool    `json:"enabled,bool"`

  // BaremetalCloudBoot
  FailoverRecipeID              int     `json:"failover_recipe_id,omitempty"`

  HypervisorGroupID             int     `json:"hypervisor_group_id,omitempty"`
  CPUUnits                      int     `json:"cpu_units,omitempty"`

  // SmartCloudBoot, BaremetalCloudBoot
  PxeIPAddressID                int     `json:"pxe_ip_address_id,omitempty"`

  // CloudBoot, SmartCloudBoot, BaremetalCloudBoot
  // by default: virtual
  // SmartCloudBoot: smart
  // BaremetalCloudBoot: baremetal
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

  arr := make([]Hypervisor, len(out))
  for i := range arr {
    arr[i] = out[i]["hypervisor"]
  }

  return arr, resp, err
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

  mac := "00:00:00:00:00:00"
  path := hypervisorPath(mac, createRequest.ServerType)
  rootRequest := &hypervisorCreateRequestRoot{
    HypervisorCreateRequest: createRequest,
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
  if err != nil {
    return nil, resp, err
  }

  return lastTransaction(ctx, s.client, id, "Hypervisor")
}

func hypervisorPath(mac string, serverType string) string {
  if serverType == "virtual" {
    path := hypervisorBasePath + apiFormat
    return path
  } else if serverType == "smart" || serverType == "baremetal" {
    path := cloudBootHypervisorBasePath + apiFormat
    return fmt.Sprintf(path, mac)
  }

  return ""
}

// Debug - print formatted Hypervisor structure
func (obj Hypervisor) Debug() {
  fmt.Printf("            ID: %d\n", obj.ID)
  fmt.Printf("         Label: %s\n", obj.Label)
  fmt.Printf("     IPAddress: %s\n", obj.IPAddress)
  fmt.Printf("    ServerType: %s\n", obj.ServerType)
  fmt.Printf("HypervisorType: %s\n", obj.HypervisorType)
  fmt.Printf("          Host: %s\n", obj.Host)
  fmt.Printf("           Mac: %s\n", obj.Mac)
  fmt.Printf("        Online: %t\n", obj.Online)
  fmt.Printf("     TotalCpus: %d\n", obj.TotalCpus)
  fmt.Printf("      TotalMem: %d\n", obj.TotalMem)
  fmt.Printf("        Uptime: %s\n", obj.Uptime)
}
