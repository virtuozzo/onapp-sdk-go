package onappgo

import (
  "time"
)

// Backup represent VirtualMachine backup
type Backup struct {
  AllowResizeWithoutReboot bool        `json:"allow_resize_without_reboot,bool,omitempty"`
  AllowedHotMigrate        bool        `json:"allowed_hot_migrate,bool,omitempty"`
  AllowedSwap              bool        `json:"allowed_swap,bool,omitempty"`
  BackupServerID           int         `json:"backup_server_id,omitempty"`
  BackupSize               int         `json:"backup_size,omitempty"`
  BackupType               string      `json:"backup_type,omitempty"`
  Built                    bool        `json:"built,bool,omitempty"`
  BuiltAt                  time.Time   `json:"built_at,omitempty"`
  CreatedAt                time.Time   `json:"created_at,omitempty"`
  DataStoreType            string      `json:"data_store_type,omitempty"`
  DiskID                   int         `json:"disk_id,omitempty"`
  ID                       int         `json:"id,omitempty"`
  Identifier               string      `json:"identifier,omitempty"`
  Initiated                string      `json:"initiated,omitempty"`
  Iqn                      string      `json:"iqn,omitempty"`
  Locked                   bool        `json:"locked,bool,omitempty"`
  MarkedForDelete          bool        `json:"marked_for_delete,bool,omitempty"`
  MinDiskSize              int         `json:"min_disk_size,omitempty"`
  MinMemorySize            int         `json:"min_memory_size,omitempty"`
  Note                     string      `json:"note,omitempty"`
  OperatingSystem          string      `json:"operating_system,omitempty"`
  OperatingSystemDistro    string      `json:"operating_system_distro,omitempty"`
  TargetID                 int         `json:"target_id,omitempty"`
  TargetType               string      `json:"target_type,omitempty"`
  TemplateID               int         `json:"template_id,omitempty"`
  UpdatedAt                time.Time   `json:"updated_at,omitempty"`
  UserID                   int         `json:"user_id,omitempty"`
  VolumeID                 int         `json:"volume_id,omitempty"`
}
