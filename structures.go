package onappgo

import (
  // "fmt"
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
  ID             int         `json:"id,omitempty"`
  Address        string      `json:"address,omitempty"`
  Prefix         int         `json:"prefix,omitempty"`
  Broadcast      string      `json:"broadcast,omitempty"`
  NetworkAddress string      `json:"network_address,omitempty"`
  Gateway        string      `json:"gateway,omitempty"`
  CreatedAt      string      `json:"created_at,omitempty"`
  UpdatedAt      string      `json:"updated_at,omitempty"`
  Ipv4           bool        `json:"ipv4,bool"`
  UserID         int         `json:"user_id,omitempty"`
  Pxe            bool        `json:"pxe,bool"`
  HypervisorID   int         `json:"hypervisor_id,omitempty"`
  LockVersion    int         `json:"lock_version,omitempty"`
  IPRangeID      int         `json:"ip_range_id,omitempty"`
  NetworkID      int         `json:"network_id,omitempty"`
  IPNetID        int         `json:"ip_net_id,omitempty"`
}

// IPAddressJoin - 
type IPAddressJoin struct {
  ID                 int       `json:"id,omitempty"`
  IPAddressID        int       `json:"ip_address_id,omitempty"`
  NetworkInterfaceID int       `json:"network_interface_id,omitempty"`
  CreatedAt          string    `json:"created_at,omitempty"`
  UpdatedAt          string    `json:"updated_at,omitempty"`
  IPAddress          IPAddress `json:"ip_address,omitempty"`
}

// ConnectionOptions for VMware hypervisor
type ConnectionOptions struct {
  APIURL                       string `json:"api_url,omitempty"`
  Login                        string `json:"login,omitempty"`
  Password                     string `json:"password,omitempty"`
  ClusterName                  string `json:"cluster_name,omitempty"`
  DistributedVirtualSwitchName string `json:"distributed_virtual_switch_name,omitempty"`
}

// StorageDisk of storage
type StorageDisk struct {
  Scsi     string   `json:"scsi,omitempty"`
  Selected bool     `json:"selected,bool"`
}

// StorageNic of storage
type StorageNic struct {
  Mac  string  `json:"mac,omitempty"`
  Type int     `json:"type,omitempty"`
}

// StorageCustomPCI of storage
type StorageCustomPCI struct {
  Pci      string  `json:"pci,omitempty"`
  Selected bool    `json:"selected,bool"`
}

// Storage for hypervisor
type Storage struct {
  Disks      []StorageDisk        `json:"disks,omitempty"`
  Nics       []StorageNic         `json:"nics,omitempty"`
  CustomPCIs []StorageCustomPCI   `json:"custom_pcis,omitempty"`
}

// IntegratedStorageCacheSettings - 
type IntegratedStorageCacheSettings struct {
}

// IoLimits - 
type IoLimits struct {
  ReadIops          int `json:"read_iops,omitempty"`
  WriteIops         int `json:"write_iops,omitempty"`
  ReadThroughput    int `json:"read_throughput,omitempty"`
  WriteThroughput   int `json:"write_throughput,omitempty"`
}

// AdditionalField - 
type AdditionalField struct {
  Name  string `json:"name,omitempty"`
  Value string `json:"value,omitempty"`
}

// AdditionalFields - 
type AdditionalFields struct {
  AdditionalField AdditionalField `json:"additional_field,omitempty"`
}

// AdvancedOptions - 
type AdvancedOptions struct {
}

// NetworkInterface - 
// https://docs.onapp.com/apim/latest/network-interfaces
type NetworkInterface struct {
  ID                  int         `json:"id,omitempty"`
  VirtualMachineID    int         `json:"virtual_machine_id,omitempty"`
  Label               string      `json:"label,omitempty"`
  Identifier          string      `json:"identifier,omitempty"`
  CreatedAt           string      `json:"created_at,omitempty"`
  UpdatedAt           string      `json:"updated_at,omitempty"`
  Primary             bool        `json:"primary,bool"`
  MacAddress          string      `json:"mac_address,omitempty"`
  NetworkJoinID       int         `json:"network_join_id,omitempty"`
  Usage               float64     `json:"usage,omitempty"`
  UsageLastResetAt    string      `json:"usage_last_reset_at,omitempty"`
  UsageMonthRolledAt  string      `json:"usage_month_rolled_at,omitempty"`
  RateLimit           int         `json:"rate_limit,omitempty"`
  DefaultFirewallRule string      `json:"default_firewall_rule,omitempty"`
  Connected           bool        `json:"connected,bool"`
  EdgeGatewayID       int         `json:"edge_gateway_id,omitempty"`
  UseAsGateway        bool        `json:"use_as_gateway,bool"`
  OpenstackID         int         `json:"openstack_id,omitempty"`
  AdapterType         string      `json:"adapter_type,omitempty"`
}

// FirewallRule - 
// https://docs.onapp.com/apim/latest/firewall-rules-for-vss
type FirewallRule struct {
  ID                  int         `json:"id,omitempty"`
  Position            int         `json:"position,omitempty"`
  Address             string      `json:"address,omitempty"`
  CreatedAt           string      `json:"created_at,omitempty"`
  UpdatedAt           string      `json:"updated_at,omitempty"`
  Command             string      `json:"adapter_type,omitempty"`
  Port                int         `json:"port,omitempty"`
  Protocol            string      `json:"protocol,omitempty"`
  NetworkInterfaceID  int         `json:"network_interface_id,omitempty"`
}

// AssignIPAddress - used for assign IPAddress to the VirtualMachine or User
type AssignIPAddress struct {
  Address             string  `json:"address,omitempty"`
  NetworkInterfaceID  int     `json:"network_interface_id,omitempty"`
  IPNetID             int     `json:"ip_net_id,omitempty"`
  IPRangeID           int     `json:"ip_range_id,omitempty"`
  UsedIP              int     `json:"used_ip,omitempty"`
  OwnIP               int     `json:"own_ip,omitempty"`

  // 6 or 4
  IPVersion           int     `json:"ip_version,omitempty"`
}

type LimitResourceRoots     map[string]*Limits

type PriceResourceRoots     map[string]*Prices

type AccessControlLimits    map[string]*LimitResourceRoots

type RateCardLimits         map[string]*PriceResourceRoots
