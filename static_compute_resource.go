package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/google/uuid"
)

// Xen/KVM, VMware - CRUD
// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Get, Delete
const hypervisorsBasePath string = "settings/hypervisors"

// Used to get data for integrated storage
const hypervisorHardwareDeviceBasePath string = hypervisorsBasePath + "/%d/hardware_devices"
const hypervisorHardwareDeviceRefreshBasePath string = hypervisorHardwareDeviceBasePath + "/refresh"
const hypervisorIntegratedStorageSettingBasePath string = hypervisorsBasePath + "/%d/integrated_storage_settings"
const hypervisorRebootBasePath string = hypervisorsBasePath + "/%d/reboot"

// HypervisorsService is an interface for interfacing with the Hypervisor
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/compute-resources
type HypervisorsService interface {
	List(context.Context, *ListOptions) ([]Hypervisor, *Response, error)
	Get(context.Context, int) (*Hypervisor, *Response, error)
	Create(context.Context, *HypervisorCreateRequest) (*Hypervisor, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *HypervisorEditRequest) (*Response, error)

	Reboot(context.Context, int, *HypervisorRebootRequest) (*Response, error)

	Refresh(context.Context, int) (*HardwareDevices, *Response, error)
	Attach(context.Context, int, map[string]interface{}) (*Response, error)

	EditIntegratedStorageSettings(context.Context, int, *IntegratedStorageSettings) (*Response, error)
}

// HypervisorsServiceOp handles communication with the Hypervisor related methods of the
// OnApp API.
type HypervisorsServiceOp struct {
	client *Client
}

var _ HypervisorsService = &HypervisorsServiceOp{}

// Hypervisor represent Hypervisor of the OnApp API
type Hypervisor struct {
	AllowUnsafeAssignedInterrupts    bool              `json:"allow_unsafe_assigned_interrupts,bool"`
	AmqpExchangeName                 string            `json:"amqp_exchange_name,omitempty"`
	ApplyHypervisorGroupCustomConfig bool              `json:"apply_hypervisor_group_custom_config,bool"`
	Backup                           bool              `json:"backup,bool"`
	BackupIPAddress                  string            `json:"backup_ip_address,omitempty"`
	Blocked                          bool              `json:"blocked,bool"`
	Built                            bool              `json:"built,bool"`
	CalledInAt                       string            `json:"called_in_at,omitempty"`
	CloudBootOs                      string            `json:"cloud_boot_os,omitempty"`
	ConnectionOptions                ConnectionOptions `json:"connection_options,omitempty"`
	CPUCores                         int               `json:"cpu_cores,omitempty"`
	CPUFlags                         []string          `json:"cpu_flags,omitempty"`
	CPUIdle                          int               `json:"cpu_idle,omitempty"`
	CPUMhz                           string            `json:"cpu_mhz,omitempty"`
	CPUModel                         string            `json:"cpu_model,omitempty"`
	Cpus                             int               `json:"cpus,omitempty"`
	CPUUnits                         int               `json:"cpu_units,omitempty"`
	CrashDebug                       bool              `json:"crash_debug,bool"`
	CreatedAt                        string            `json:"created_at,omitempty"`
	CustomConfig                     string            `json:"custom_config,omitempty"`
	DisableFailover                  bool              `json:"disable_failover,bool"`
	Distro                           string            `json:"distro,omitempty"`
	Dom0MemorySize                   int               `json:"dom0_memory_size,omitempty"`
	Enabled                          bool              `json:"enabled,bool"`
	FailoverRecipeID                 int               `json:"failover_recipe_id,omitempty"`
	FailureCount                     int               `json:"failure_count,omitempty"`
	FormatDisks                      bool              `json:"format_disks,bool"`
	FreeDiskSpace                    map[string]int    `json:"free_disk_space,omitempty"`
	FreeMem                          int               `json:"free_mem,omitempty"`
	FreeMemory                       int               `json:"free_memory,omitempty"`
	Host                             string            `json:"host,omitempty"`
	HypervisorGroupID                int               `json:"hypervisor_group_id,omitempty"`
	HypervisorType                   string            `json:"hypervisor_type,omitempty"`
	ID                               int               `json:"id,omitempty"`
	InstanceUUID                     string            `json:"instance_uuid,omitempty"`
	IntegratedStorageDisabled        bool              `json:"integrated_storage_disabled,bool"`
	IPAddress                        string            `json:"ip_address,omitempty"`
	Label                            string            `json:"label,omitempty"`
	ListOfLogicalVolumes             string            `json:"list_of_logical_volumes,omitempty"`
	ListOfVolumeGroups               string            `json:"list_of_volume_groups,omitempty"`
	ListOfZombieDomains              string            `json:"list_of_zombie_domains,omitempty"`
	Locked                           bool              `json:"locked,bool"`
	Mac                              string            `json:"mac,omitempty"`
	Machine                          string            `json:"machine,omitempty"`
	MaintenanceMode                  bool              `json:"maintenance_mode,bool"`
	MemInfo                          int               `json:"mem_info,omitempty"`
	MemoryAllocatedByRunningVms      int               `json:"memory_allocated_by_running_vms,omitempty"`
	Online                           bool              `json:"online,bool"`
	OsVersion                        int               `json:"os_version,omitempty"`
	OsVersionMinor                   int               `json:"os_version_minor,omitempty"`
	PassthroughDisks                 bool              `json:"passthrough_disks,bool"`
	PowerCycleCommand                string            `json:"power_cycle_command,omitempty"`
	Rebooting                        bool              `json:"rebooting,bool"`
	Release                          string            `json:"release,omitempty"`
	SegregationOsType                string            `json:"segregation_os_type,omitempty"`
	ServerType                       string            `json:"server_type,omitempty"`
	Spare                            bool              `json:"spare,bool"`
	StaticIntegratedStorage          bool              `json:"static_integrated_storage,bool"`
	ThreadsPerCore                   int               `json:"threads_per_core,omitempty"`
	TotalCpus                        int               `json:"total_cpus,omitempty"`
	TotalMem                         int               `json:"total_mem,omitempty"`
	TotalMemory                      int               `json:"total_memory,omitempty"`
	TotalMemoryAllocatedByVms        int               `json:"total_memory_allocated_by_vms,omitempty"`
	TotalZombieMem                   int               `json:"total_zombie_mem,omitempty"`
	UpdatedAt                        string            `json:"updated_at,omitempty"`
	Uptime                           string            `json:"uptime,omitempty"`
	UsedCPUResources                 int               `json:"used_cpu_resources,omitempty"`
}

// HypervisorCreateRequest represents a request to create a Hypervisor
type HypervisorCreateRequest struct {
	Label                   string `json:"label,omitempty"`
	IPAddress               string `json:"ip_address,omitempty"`
	BackupIPAddress         string `json:"backup_ip_address,omitempty"`
	DisableFailover         bool   `json:"disable_failover,bool"`
	HypervisorType          string `json:"hypervisor_type,omitempty"`
	SegregationOsType       string `json:"segregation_os_type,omitempty"`
	Enabled                 bool   `json:"enabled,bool"`
	FailoverRecipeID        int    `json:"failover_recipe_id,omitempty"`
	HypervisorGroupID       int    `json:"hypervisor_group_id,omitempty"`
	CPUUnits                int    `json:"cpu_units,omitempty"`
	StaticIntegratedStorage bool   `json:"static_integrated_storage,bool"`
	PowerCycleCommand       string `json:"power_cycle_command,omitempty"`
}

// HypervisorEditRequest represents a request to edit a Hypervisor
type HypervisorEditRequest struct {
	Label             string `json:"label,omitempty"`
	IPAddress         string `json:"ip_address,omitempty"`
	BackupIPAddress   string `json:"backup_ip_address,omitempty"`
	SegregationOsType string `json:"segregation_os_type,omitempty"`
	Enabled           bool   `json:"enabled,bool"`
	FailoverRecipeID  int    `json:"failover_recipe_id,omitempty"`
	HypervisorGroupID int    `json:"hypervisor_group_id,omitempty"`
	CPUUnits          int    `json:"cpu_units,omitempty"`
}

type hypervisorCreateRequestRoot struct {
	HypervisorCreateRequest *HypervisorCreateRequest `json:"hypervisor"`
}

type hypervisorRoot struct {
	Hypervisor *Hypervisor `json:"hypervisor"`
}

func (d HypervisorCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Hypervisors.
func (s *HypervisorsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Hypervisor, *Response, error) {
	path := hypervisorsBasePath + apiFormat
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

	path := fmt.Sprintf("%s/%d%s", hypervisorsBasePath, id, apiFormat)
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

	path := hypervisorsBasePath + apiFormat
	rootRequest := &hypervisorCreateRequestRoot{
		HypervisorCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Hypervisor [Create] req: ", req)

	root := new(hypervisorRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Hypervisor, resp, err
}

// Delete Hypervisor.
func (s *HypervisorsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorsBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit Hypervisor.
func (s *HypervisorsServiceOp) Edit(ctx context.Context, id int, editRequest *HypervisorEditRequest) (*Response, error) {
	if editRequest == nil || id < 1 {
		return nil, godo.NewArgError("editRequest || id", "cannot be nil or less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

type HypervisorRebootRequest struct {
	SkipPoweredOffVmsMigration int `json:"skip_powered_off_vms_migration,omitempty"`
	SheduleFailover            int `json:"schedule_failover,omitempty"`
	Force                      int `json:"force,omitempty"`
	Confirm                    int `json:"confirm,omitempty"`
}

// Reboot Hypervisor
func (s *HypervisorsServiceOp) Reboot(ctx context.Context, id int, rebootRequest *HypervisorRebootRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorRebootBasePath, id) + apiFormat
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, rebootRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [Reboot] req: ", req)

	return s.client.Do(ctx, req, nil)
}

type HardwareDevices struct {
	HardwareCustomDevice           []*HardwareCustomDevice           `json:"hardware_custom_device,omitempty"`
	HardwareDiskDevice             []*HardwareDiskDevice             `json:"hardware_disk_device,omitempty"`
	HardwareDiskPciDevice          []*HardwareDiskPciDevice          `json:"hardware_disk_pci_device,omitempty"`
	HardwareNetworkInterfaceDevice []*HardwareNetworkInterfaceDevice `json:"hardware_network_interface_device,omitempty"`
}

type rootHardware []struct {
	HardwareCustomDevice           *HardwareCustomDevice           `json:"hardware_custom_device,omitempty"`
	HardwareDiskDevice             *HardwareDiskDevice             `json:"hardware_disk_device,omitempty"`
	HardwareDiskPciDevice          *HardwareDiskPciDevice          `json:"hardware_disk_pci_device,omitempty"`
	HardwareNetworkInterfaceDevice *HardwareNetworkInterfaceDevice `json:"hardware_network_interface_device,omitempty"`
}

type HardwareDiskDevice struct {
	ID         int    `json:"id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"`
	Status     string `json:"status,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	ParentType string `json:"parent_type,omitempty"`
	Name       string `json:"name,omitempty"`
	Scsi       string `json:"scsi,omitempty"`
}

type HardwareNetworkInterfaceDevice struct {
	ID            int    `json:"id,omitempty"`
	ParentID      int    `json:"parent_id,omitempty"`
	Status        string `json:"status,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	ParentType    string `json:"parent_type,omitempty"`
	Name          string `json:"name,omitempty"`
	Pci           string `json:"pci,omitempty"`
	Mac           string `json:"mac,omitempty"`
	InterfaceType string `json:"interface_type,omitempty"`
}

type HardwareCustomDevice struct {
	ID         int    `json:"id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"`
	Status     string `json:"status,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	ParentType string `json:"parent_type,omitempty"`
	Name       string `json:"name,omitempty"`
	Pci        string `json:"pci,omitempty"`
	Code       string `json:"code,omitempty"`
}

type HardwareDiskPciDevice struct {
	ID         int    `json:"id,omitempty"`
	ParentID   int    `json:"parent_id,omitempty"`
	Status     string `json:"status,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	ParentType string `json:"parent_type,omitempty"`
	Pci        string `json:"pci,omitempty"`
}

// Refresh - get list of hardware devices (disks, network interfaces) from hypervisor with enabled integrated storage
func (s *HypervisorsServiceOp) Refresh(ctx context.Context, resID int) (*HardwareDevices, *Response, error) {
	if resID < 1 {
		return nil, nil, godo.NewArgError("Hypervisor.Refresh", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorHardwareDeviceRefreshBasePath, resID) + apiFormat
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	key, _ := uuid.NewRandom()
	req.Header.Add("X-Idempotency-Key", key.String())

	log.Println("Hypervisor [Refresh] req: ", req)

	out := &rootHardware{}
	resp, err := s.client.Do(ctx, req, out)
	if err != nil {
		return nil, resp, err
	}

	// log.Printf("IntegratedStorages out: %+v\n", out)

	res := &HardwareDevices{}
	res.initHardwareDevices(out)

	return res, resp, err
}

func (obj *HardwareDevices) initHardwareDevices(s *rootHardware) {
	for _, v := range *s {
		if v.HardwareCustomDevice != nil {
			obj.HardwareCustomDevice = append(obj.HardwareCustomDevice, v.HardwareCustomDevice)
		} else if v.HardwareDiskDevice != nil {
			obj.HardwareDiskDevice = append(obj.HardwareDiskDevice, v.HardwareDiskDevice)
		} else if v.HardwareDiskPciDevice != nil {
			obj.HardwareDiskPciDevice = append(obj.HardwareDiskPciDevice, v.HardwareDiskPciDevice)
		} else if v.HardwareNetworkInterfaceDevice != nil {
			obj.HardwareNetworkInterfaceDevice = append(obj.HardwareNetworkInterfaceDevice, v.HardwareNetworkInterfaceDevice)
		}
	}
}

// AttachDiskHardwareDevice represents a request to attach disks from hypervisor to the integrated data store
type AttachDiskHardwareDevice struct {
	Status string `json:"status,omitempty"`
	Format bool   `json:"format,bool"`
}

// AttachNetworkInterfaceHardwareDevice represents a request to attach network interfaces from hypervisor to the integrated data store
type AttachNetworkInterfaceHardwareDevice struct {
	Status string `json:"status,omitempty"`
}

type hardwareDevicesRoot struct {
	AttachHardwareDevices map[string]interface{} `json:"hardware_devices"`
}

const AssignedToCache string = "assigned_to_cache"
const AssignedToStorage string = "assigned_to_storage"
const AssignedToSAN string = "assigned_to_san"
const Unassigned string = "unassigned"

// Attach - attach disks, network interfaces of hypervisor to the integrated data store
func (s *HypervisorsServiceOp) Attach(ctx context.Context, resID int, attachRequest map[string]interface{}) (*Response, error) {
	if resID < 1 {
		return nil, godo.NewArgError("Hypervisor.Attach", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorHardwareDeviceBasePath, resID) + apiFormat

	rootRequest := &hardwareDevicesRoot{
		AttachHardwareDevices: attachRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, rootRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [Attach] req: ", req)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// IntegratedStorageSettings -
// TODO: we must change all fields to the type 'int' only skip 'bonding_mode' this field must be string
// because we got normal structure, but for editing we must send string so there is some data inconsistency
type IntegratedStorageSettings struct {
	BondingMode          string `json:"bonding_mode,omitempty"`
	CacheMirrors         string `json:"cache_mirrors,omitempty"`
	CacheStripes         string `json:"cache_stripes,omitempty"`
	ControllerDbSize     string `json:"controller_db_size,omitempty"`
	ControllerMemorySize string `json:"controller_memory_size,omitempty"`
	DisksPerController   string `json:"disks_per_controller,omitempty"`
	Mtu                  string `json:"mtu,omitempty"`
	Vlan                 string `json:"vlan,omitempty"`
}

// From cloudboot
// type IntegratedStorageSettings struct {
// 	BondingMode          string `json:"bonding_mode,omitempty"`
// 	BondName             string `json:"bond_name,omitempty"`
// 	BridgeName           string `json:"bridge_name,omitempty"`
// 	CacheMirrors         int    `json:"cache_mirrors,omitempty"`
// 	CacheStripes         int    `json:"cache_stripes,omitempty"`
// 	ControllerMemorySize int    `json:"controller_memory_size,omitempty"`
// 	DbSize               int    `json:"db_size,omitempty"`
// 	DisksPerController   int    `json:"disks_per_controller,omitempty"`
// 	Mtu                  int    `json:"mtu,omitempty"`
// 	Vlan                 string `json:"vlan,omitempty"`
// }

type integratedStorageSettingCreateRequestRoot struct {
	IntegratedStorageSettings *IntegratedStorageSettings `json:"integrated_storage_settings,omitempty"`
}

// EditIntegratedStorageSettings -
func (s *HypervisorsServiceOp) EditIntegratedStorageSettings(ctx context.Context, id int, editRequest *IntegratedStorageSettings) (*Response, error) {
	if editRequest == nil || id < 1 {
		return nil, godo.NewArgError("editRequest || id", "cannot be nill or less than 1")
	}

	path := fmt.Sprintf(hypervisorIntegratedStorageSettingBasePath, id) + apiFormat

	rootRequest := &integratedStorageSettingCreateRequestRoot{
		IntegratedStorageSettings: editRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, rootRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [EditIntegratedStorageSettings]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
