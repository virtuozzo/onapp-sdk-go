package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

// Xen/KVM, VMware - CRUD
// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Get, Delete
const hypervisorsBasePath string = "settings/hypervisors"

// TODO: maybe later remove this because will be DataStoreJoins, NetworkJoins
// BackupServerJoins objects
const hypervisorsDataStoreJoins string = hypervisorsBasePath + "/%d/data_store_joins"
// const hypervisorsNetworkJoins string = hypervisorsBasePath + "/%d/network_joins"
const hypervisorsBackupServerJoins string = hypervisorsBasePath + "/%d/backup_server_joins"

// Used to get data for integrated storeg
const hypervisorHardwareDevicesBasePath string = hypervisorsBasePath + "/%d/hardware_devices/refresh"

// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Create, Edit
const cloudBootHypervisorsBasePath string = "settings/assets/%s/hypervisors"

// HypervisorsService is an interface for interfacing with the Hypervisor
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/compute-resources
type HypervisorsService interface {
	List(context.Context, *ListOptions) ([]Hypervisor, *Response, error)
	Get(context.Context, int) (*Hypervisor, *Response, error)
	Create(context.Context, *HypervisorCreateRequest) (*Hypervisor, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *HypervisorEditRequest) (*Response, error)

	// TODO: maybe later remove this because will be DataSoreJoins, NetworkJoins
	// BackupServerJoins objects
	AddDataStoreJoins(context.Context, int, int) (*Response, error)
	// AddNetworkJoins(context.Context, int, *HypervisorNetworkJoinCreateRequest) (*Response, error)
	AddBackupServerJoins(context.Context, int, int) (*Response, error)

	DeleteDataStoreJoins(context.Context, int, int) (*Response, error)
	// DeleteNetworkJoins(context.Context, int, int) (*Response, error)
	DeleteBackupServerJoins(context.Context, int, int) (*Response, error)

	Refresh(context.Context, int) (*HardwareDevices, *Response, error)
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

// DataStoreJoinCreateRequest -
type DataStoreJoinCreateRequest struct {
	DataStoreID int `json:"data_store_id,omitempty"`
}

// BackupServerJoinCreateRequest -
type BackupServerJoinCreateRequest struct {
	BackupServerID int `json:"backup_server_id,omitempty"`
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

	// mac := "00:00:00:00:00:00"
	// path := hypervisorPath(mac, createRequest.ServerType)
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
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("Hypervisor [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Hypervisor [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// AddDataStoreJoins - add Data Store Joins to the Hypervisor
func (s *HypervisorsServiceOp) AddDataStoreJoins(ctx context.Context, hvID int, dsID int) (*Response, error) {
	if hvID < 1 || dsID < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorsDataStoreJoins, hvID) + apiFormat

	rootRequest := &DataStoreJoinCreateRequest{
		DataStoreID: dsID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, err
	}

	log.Println("DataStoreJoins [Create] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// DeleteDataStoreJoins - delete Data Store Joins from the Hypervisor
func (s *HypervisorsServiceOp) DeleteDataStoreJoins(ctx context.Context, hvID int, id int) (*Response, error) {
	if hvID < 1 || id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorsDataStoreJoins, hvID)
	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Delete DataStore Joins from Hypervisor [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// AddBackupServerJoins - add Backup Server Joins to the Hypervisor
func (s *HypervisorsServiceOp) AddBackupServerJoins(ctx context.Context, hvID int, bsID int) (*Response, error) {
	if hvID < 1 || bsID < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorsBackupServerJoins, hvID) + apiFormat

	rootRequest := &BackupServerJoinCreateRequest{
		BackupServerID: bsID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, err
	}

	log.Println("BackupServerJoins [Create] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// DeleteBackupServerJoins - delete Backup Server Joins from the Hypervisor
func (s *HypervisorsServiceOp) DeleteBackupServerJoins(ctx context.Context, hvID int, id int) (*Response, error) {
	if hvID < 1 || id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorsBackupServerJoins, hvID)
	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Delete Backup Server Joins from Hypervisor [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

func hypervisorPath(mac string, serverType string) string {
	if serverType == "virtual" {
		path := hypervisorsBasePath + apiFormat
		return path
	} else if serverType == "smart" || serverType == "baremetal" {
		path := cloudBootHypervisorsBasePath + apiFormat
		return fmt.Sprintf(path, mac)
	}

	return ""
}

type HardwareDevices struct {
	HardwareCustomDevice           []*HardwareCustomDevice
	HardwareDiskDevice             []*HardwareDiskDevice
	HardwareDiskPciDevice          []*HardwareDiskPciDevice
	HardwareNetworkInterfaceDevice []*HardwareNetworkInterfaceDevice
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
	if resID < 1  {
		return nil, nil, godo.NewArgError("HypervisorsServiceOp.Refresh", "cannot be less than 1")
	}

	path := fmt.Sprintf(hypervisorHardwareDevicesBasePath, resID) + apiFormat
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("HypervisorsServiceOp [Refresh] req: ", req)

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
