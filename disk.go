package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const disksBasePath string = "settings/disks"

// DisksService is an interface for interfacing with the Disk
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/disks
type DisksService interface {
	List(context.Context, *ListOptions) ([]Disk, *Response, error)
	Get(context.Context, int) (*Disk, *Response, error)
	Create(context.Context, *DiskCreateRequest) (*Disk, *Response, error)
	Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
	Edit(context.Context, int, *DiskEditRequest) (*Response, error)
}

// DisksServiceOp handles communication with the Disk related methods of the
// OnApp API.
type DisksServiceOp struct {
	client *Client
}

var _ DisksService = &DisksServiceOp{}

// Disk - represent disk from Virtual Machine
type Disk struct {
	AddToFreebsdFstab              bool                           `json:"add_to_freebsd_fstab,bool"`
	AddToLinuxFstab                bool                           `json:"add_to_linux_fstab,bool"`
	Built                          bool                           `json:"built,bool"`
	BurstBw                        int                            `json:"burst_bw,omitempty"`
	BurstIops                      int                            `json:"burst_iops,omitempty"`
	CreatedAt                      string                         `json:"created_at,omitempty"`
	DataStoreID                    int                            `json:"data_store_id,omitempty"`
	DiskSize                       int                            `json:"disk_size,omitempty"`
	DiskVMNumber                   int                            `json:"disk_vm_number,omitempty"`
	FileSystem                     string                         `json:"file_system,omitempty"`
	HasAutobackups                 bool                           `json:"has_autobackups"`
	ID                             int                            `json:"id,omitempty"`
	Identifier                     string                         `json:"identifier,omitempty"`
	IntegratedStorageCacheEnabled  bool                           `json:"integrated_storage_cache_enabled,bool"`
	IntegratedStorageCacheOverride bool                           `json:"integrated_storage_cache_override,bool"`
	IntegratedStorageCacheSettings IntegratedStorageCacheSettings `json:"integrated_storage_cache_settings,omitempty"`
	IoLimits                       IoLimits                       `json:"io_limits,omitempty"`
	IoLimitsOverride               bool                           `json:"io_limits_override"`
	Iqn                            string                         `json:"iqn,omitempty"`
	IsSwap                         bool                           `json:"is_swap,bool"`
	Label                          string                         `json:"label,omitempty"`
	Locked                         bool                           `json:"locked,bool"`
	MaxBw                          int                            `json:"max_bw,omitempty"`
	MaxIops                        int                            `json:"max_iops,omitempty"`
	MinIops                        int                            `json:"min_iops,omitempty"`
	MountPoint                     string                         `json:"mount_point,omitempty"`
	Mounted                        bool                           `json:"mounted,bool"`
	OpenstackID                    int                            `json:"openstack_id,omitempty"`
	Primary                        bool                           `json:"primary,bool"`
	TemporaryVirtualMachineID      int                            `json:"temporary_virtual_machine_id,omitempty"`
	UpdatedAt                      string                         `json:"updated_at,omitempty"`
	VirtualMachineID               int                            `json:"virtual_machine_id,omitempty"`
	VolumeID                       int                            `json:"volume_id,omitempty"`
}

// DiskCreateRequest - data for creating Disk
type DiskCreateRequest struct {
	AddToLinuxFstab   bool   `json:"add_to_linux_fstab,bool"`
	AddToFreebsdFstab bool   `json:"add_to_freebsd_fstab,bool"`
	Primary           bool   `json:"primary,bool"`
	DiskSize          int    `json:"disk_size,omitempty"`
	FileSystem        string `json:"file_system,omitempty"`
	DataStoreID       int    `json:"data_store_id,omitempty"`
	Label             string `json:"label,omitempty"`
	RequireFormatDisk bool   `json:"require_format_disk,bool"`
	MountPoint        string `json:"mount_point,omitempty"`
	HotAttach         bool   `json:"hot_attach,bool"`
	MinIops           int    `json:"min_iops,omitempty"`
	Mounted           bool   `json:"mounted,bool"`

	// Additional field to determine Virtual Machine to create disk
	VirtualMachineID  int    `json:"-"`
}

// DiskEditRequest - data for editing Disk
type DiskEditRequest struct {
	AddToLinuxFstab   bool   `json:"add_to_linux_fstab,bool"`
	AddToFreebsdFstab bool   `json:"add_to_freebsd_fstab,bool"`
	DiskSize          int    `json:"disk_size,omitempty"`
	FileSystem        string `json:"file_system,omitempty"`
	DataStoreID       int    `json:"data_store_id,omitempty"`
	Label             string `json:"label,omitempty"`
	RequireFormatDisk bool   `json:"require_format_disk,bool"`
	MountPoint        string `json:"mount_point,omitempty"`
	MinIops           int    `json:"min_iops,omitempty"`
	Mounted           string `json:"mounted,omitempty"`
}

type diskCreateRequestRoot struct {
	DiskCreateRequest *DiskCreateRequest `json:"disk"`
}

type diskRoot struct {
	Disk *Disk `json:"disk"`
}

func (d DiskCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Disks in the cloud.
func (s *DisksServiceOp) List(ctx context.Context, opt *ListOptions) ([]Disk, *Response, error) {
	path := disksBasePath + apiFormat
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

	arr := make([]Disk, len(out))
	for i := range arr {
		arr[i] = out[i]["disk"]
	}

	return arr, resp, err
}

// Get individual Disk.
func (s *DisksServiceOp) Get(ctx context.Context, id int) (*Disk, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", disksBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(diskRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Disk, resp, err
}

// Create Disk.
func (s *DisksServiceOp) Create(ctx context.Context, createRequest *DiskCreateRequest) (*Disk, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d/disks%s", virtualMachineBasePath, createRequest.VirtualMachineID, apiFormat)

	rootRequest := &diskCreateRequestRoot{
		DiskCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Disk [Create]  req: ", req)

	root := new(diskRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Disk, resp, err
}

// Delete Disk.
func (s *DisksServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", disksBasePath, id, apiFormat)

	path, err := addOptions(path, meta)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Disk [Delete]  req: ", req)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, resp, err
	}

	filter := struct {
		ParentID   int
		ParentType string
	}{
		ParentID:   id,
		ParentType: "Disk",
	}

	return lastTransaction(ctx, s.client, filter)
}

// Edit Disk.
func (s *DisksServiceOp) Edit(ctx context.Context, id int, editRequest *DiskEditRequest) (*Response, error) {
	path := fmt.Sprintf("%s/%d%s", disksBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Disk [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
