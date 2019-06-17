package onappgo

import (
  "fmt"
  "time"
)

// IntegratedStorageCacheSettings - 
type IntegratedStorageCacheSettings struct {
}

// IoLimits - 
type IoLimits struct {
}

// Disk - 
type Disk struct {
  AddToFreebsdFstab              interface{}                    `json:"add_to_freebsd_fstab,omitempty"`
  AddToLinuxFstab                interface{}                    `json:"add_to_linux_fstab,omitempty"`
  Built                          bool                           `json:"built,bool,omitempty"`
  BurstBw                        int                            `json:"burst_bw,omitempty"`
  BurstIops                      int                            `json:"burst_iops,omitempty"`
  CreatedAt                      time.Time                      `json:"created_at,omitempty"`
  DataStoreID                    int                            `json:"data_store_id,omitempty"`
  DiskSize                       int                            `json:"disk_size,omitempty"`
  DiskVMNumber                   int                            `json:"disk_vm_number,omitempty"`
  FileSystem                     string                         `json:"file_system,omitempty"`
  HasAutobackups                 bool                           `json:"has_autobackups"`
  ID                             int                            `json:"id,omitempty"`
  Identifier                     string                         `json:"identifier,omitempty"`
  IntegratedStorageCacheEnabled  bool                           `json:"integrated_storage_cache_enabled,bool,omitempty"`
  IntegratedStorageCacheOverride bool                           `json:"integrated_storage_cache_override,bool,omitempty"`
  IntegratedStorageCacheSettings IntegratedStorageCacheSettings `json:"integrated_storage_cache_settings,omitempty"`
  IoLimits                       IoLimits                       `json:"io_limits,omitempty"`
  IoLimitsOverride               bool                           `json:"io_limits_override"`
  Iqn                            string                         `json:"iqn,omitempty"`
  IsSwap                         bool                           `json:"is_swap,bool,omitempty"`
  Label                          string                         `json:"label,omitempty"`
  Locked                         bool                           `json:"locked,bool,omitempty"`
  MaxBw                          int                            `json:"max_bw,omitempty"`
  MaxIops                        int                            `json:"max_iops,omitempty"`
  MinIops                        int                            `json:"min_iops,omitempty"`
  MountPoint                     string                         `json:"mount_point,omitempty"`
  Mounted                        bool                           `json:"mounted,bool,omitempty"`
  OpenstackID                    int                            `json:"openstack_id,omitempty"`
  Primary                        bool                           `json:"primary,bool,omitempty"`
  TemporaryVirtualMachineID      int                            `json:"temporary_virtual_machine_id,omitempty"`
  UpdatedAt                      time.Time                      `json:"updated_at,omitempty"`
  VirtualMachineID               int                            `json:"virtual_machine_id,omitempty"`
  VolumeID                       int                            `json:"volume_id,omitempty"`
}

// Debug - print formatted Disk structure
func (d Disk) Debug() {
  fmt.Printf("                   ID: %d\n", d.ID)
  fmt.Printf("           Identifier: %s\n", d.Identifier)
  fmt.Printf("     VirtualMachineID: %d\n", d.VirtualMachineID)
  fmt.Printf("                Label: %s\n", d.Label)
  fmt.Printf("           FileSystem: %s\n", d.FileSystem)
  fmt.Printf("            CreatedAt: %s\n", d.CreatedAt)
  fmt.Printf("               Locked: %t\n", d.Locked)
  fmt.Printf("             DiskSize: %d\n", d.DiskSize)
  fmt.Printf("           MountPoint: %s\n", d.MountPoint)
  fmt.Printf("             VolumeID: %d\n", d.VolumeID)
}
