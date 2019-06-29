package onappgo

import (
  "time"
  "fmt"
)

// Represent common structures for onappgo package

// Params represents a OnApp Transaction params
// type Params struct {
//   DestroyMsg       string   `json:"destroy_msg,omitempty"`
//   InitiatorID      int      `json:"initiator_id,omitempty"`
//   RealUserID       int      `json:"real_user_id,omitempty"`
//   RemoteIP         string   `json:"remote_ip,omitempty"`
//   SkipNotification bool     `json:"skip_notification,bool"`
//   ShutdownType     string   `json:"shutdown_type,omitempty"`
// }

// IPAddress - represents an IP address
type IPAddress struct {
  Address         string      `json:"address,omitempty"`
  Broadcast       string      `json:"broadcast,omitempty"`
  CreatedAt       time.Time   `json:"created_at,omitempty"`
  Free            bool        `json:"free,bool"`
  Gateway         string      `json:"gateway,omitempty"`
  HypervisorID    string      `json:"hypervisor_id,omitempty"`
  ID              int         `json:"id,omitempty"`
  IPRangeID       int         `json:"ip_range_id,omitempty"`
  IPv4            bool        `json:"ipv4,bool"`
  LockVersion     int         `json:"lock_version,omitempty"`
  Netmask         string      `json:"netmask,omitempty"`
  NetworkAddress  string      `json:"network_address,omitempty"`
  Pxe             bool        `json:"pxe,bool"`
  Prefix          int         `json:"prefix,omitempty"`
  UpdatedAt       time.Time   `json:"updated_at,omitempty"`
  UserID          string      `json:"user_id,omitempty"`
}

// Debug - print formatted IPAddress structure
func (obj IPAddress) Debug() {
  fmt.Printf("\t            ID: %d\n",  obj.ID)
  fmt.Printf("\t       Address: %s\n",  obj.Address)
  fmt.Printf("\t     Broadcast: %s\n",  obj.Broadcast)
  fmt.Printf("\t       Gateway: %s\n",  obj.Gateway)
  fmt.Printf("\t       Netmask: %s\n",  obj.Netmask)
  fmt.Printf("\tNetworkAddress: %s\n",  obj.NetworkAddress)
  fmt.Printf("\t        UserID: %s\n",  obj.UserID)
  fmt.Printf("\t     IPRangeID: %d\n",  obj.IPRangeID)
  fmt.Printf("\t          Free: %t\n",  obj.Free)
  fmt.Printf("\t  HypervisorID: %s\n",  obj.HypervisorID)
  fmt.Printf("\t   LockVersion: %d\n",  obj.LockVersion)
}

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
  Selected bool     `json:"selected,bool"`
}

// Nics of hypervisor
type Nics struct {
  Mac  string  `json:"mac,omitempty"`
  Type int     `json:"type,omitempty"`
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

// IntegratedStorageCacheSettings - 
type IntegratedStorageCacheSettings struct {
}

// IoLimits - 
type IoLimits struct {
}

type AdditionalField struct {
  Name  string `json:"name,omitempty"`
  Value string `json:"value,omitempty"`
}

type AdditionalFields struct {
  AdditionalField AdditionalField `json:"additional_field,omitempty"`
}
