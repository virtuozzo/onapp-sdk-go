package onappgo

import (
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

// Debug - print formatted IPAddress structure
func (obj IPAddress) Debug() {
  fmt.Printf("\t            ID: %d\n", obj.ID)
  fmt.Printf("\t       Address: %s\n", obj.Address)
  fmt.Printf("\t     Broadcast: %s\n", obj.Broadcast)
  fmt.Printf("\t       Gateway: %s\n", obj.Gateway)
  // fmt.Printf("\t       Netmask: %s\n", obj.Netmask)
  fmt.Printf("\tNetworkAddress: %s\n", obj.NetworkAddress)
  fmt.Printf("\t        UserID: %d\n", obj.UserID)
  fmt.Printf("\t     IPRangeID: %d\n", obj.IPRangeID)
  // fmt.Printf("\t          Free: %t\n", obj.Free)
  fmt.Printf("\t  HypervisorID: %d\n", obj.HypervisorID)
  fmt.Printf("\t   LockVersion: %d\n", obj.LockVersion)
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

// Debug - print formatted NetworkInterface structure
func (obj NetworkInterface) Debug() {
  fmt.Printf("                 ID: %d\n", obj.ID)
  fmt.Printf("   VirtualMachineID: %d\n", obj.VirtualMachineID)
  fmt.Printf("              Label: %s\n", obj.Label)
  fmt.Printf("         Identifier: %s\n", obj.Identifier)
  fmt.Printf("         MacAddress: %s\n", obj.MacAddress)
  fmt.Printf("        AdapterType: %s\n", obj.AdapterType)
  fmt.Printf("            Primary: %t\n", obj.Primary)
  fmt.Printf("DefaultFirewallRule: %s\n", obj.DefaultFirewallRule)
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

// Debug - print formatted FirewallRule structure
func (obj FirewallRule) Debug() {
  fmt.Printf("      ID: %d\n", obj.ID)
  fmt.Printf("Position: %d\n", obj.Position)
  fmt.Printf(" Address: %s\n", obj.Address)
  fmt.Printf(" Command: %s\n", obj.Command)
  fmt.Printf("    Port: %d\n", obj.Port)
  fmt.Printf("Protocol: %s\n", obj.Protocol)
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

// Debug - print formatted AssignIPAddress structure
func (obj AssignIPAddress) Debug() {
  fmt.Printf("           Address: %s\n", obj.Address)
  fmt.Printf("NetworkInterfaceID: %d\n", obj.NetworkInterfaceID)
  fmt.Printf("           IPNetID: %d\n", obj.IPNetID)
  fmt.Printf("         IPRangeID: %d\n", obj.IPRangeID)
  fmt.Printf("            UsedIP: %d\n", obj.UsedIP)
  fmt.Printf("             OwnIP: %d\n", obj.OwnIP)
}
