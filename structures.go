package onappgo

// Represent common structures for onappgo package

// // IPAddress - represents an IP address
// type IPAddress struct {
// 	Address         string      `json:"address,omitempty"`
// 	Broadcast       string      `json:"broadcast,omitempty"`
// 	CreatedAt       string      `json:"created_at,omitempty"`
// 	ExternalAddress interface{} `json:"external_address,omitempty"`
// 	Free            bool        `json:"free,bool"`
// 	Gateway         string      `json:"gateway,omitempty"`
// 	HypervisorID    int         `json:"hypervisor_id,omitempty"`
// 	ID              int         `json:"id,omitempty"`
// 	IPNetID         int         `json:"ip_net_id,omitempty"`
// 	IPRangeID       int         `json:"ip_range_id,omitempty"`
// 	Ipv4            bool        `json:"ipv4,bool"`
// 	LockVersion     int         `json:"lock_version,omitempty"`
// 	NetworkAddress  string      `json:"network_address,omitempty"`
// 	NetworkID       int         `json:"network_id,omitempty"`
// 	Netmask         string      `json:"netmask,omitempty"`
// 	Prefix          int         `json:"prefix,omitempty"`
// 	Pxe             bool        `json:"pxe,bool"`
// 	UpdatedAt       string      `json:"updated_at,omitempty"`
// 	UserID          int         `json:"user_id,omitempty"`
// }

type IPAddress struct {
	ID              int    `json:"id,omitempty"`
	Address         string `json:"address,omitempty"`
	Broadcast       string `json:"broadcast,omitempty"`
	NetworkAddress  string `json:"network_address,omitempty"`
	Gateway         string `json:"gateway,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	UpdatedAt       string `json:"updated_at,omitempty"`
	UserID          int    `json:"user_id,omitempty"`
	Pxe             bool   `json:"pxe,bool"`
	HypervisorID    int    `json:"hypervisor_id,omitempty"`
	IPRangeID       int    `json:"ip_range_id,omitempty"`
	ExternalAddress string `json:"external_address,omitempty"`
	Free            bool   `json:"free,bool"`
	Netmask         string `json:"netmask,omitempty"`
}

// IPAddresses -
type IPAddresses struct {
	IPAddress IPAddress `json:"ip_address,omitempty"`
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
	APIURL        string `json:"api_url,omitempty"`
	Login         string `json:"login,omitempty"`
	OperationMode string `json:"operation_mode,omitempty"`
	Password      string `json:"password,omitempty"`
}

// IntegratedStorageCacheSettings -
type IntegratedStorageCacheSettings struct {
}

// IoLimits -
type IoLimits struct {
	ReadIops        int `json:"read_iops,omitempty"`
	WriteIops       int `json:"write_iops,omitempty"`
	ReadThroughput  int `json:"read_throughput,omitempty"`
	WriteThroughput int `json:"write_throughput,omitempty"`
}

// AdditionalFields -
type AdditionalFields struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// AdvancedOptions -
type AdvancedOptions struct {
}

// AssignIPAddress - used for assign IPAddress to the VirtualMachine or User
type AssignIPAddress struct {
	Address            string `json:"address,omitempty"`
	NetworkInterfaceID int    `json:"network_interface_id,omitempty"`
	IPNetID            int    `json:"ip_net_id,omitempty"`
	IPRangeID          int    `json:"ip_range_id,omitempty"`
	UsedIP             int    `json:"used_ip,omitempty"`
	OwnIP              int    `json:"own_ip,omitempty"`
	IPVersion          int    `json:"ip_version,omitempty"`
}

type LimitResourceRoots map[string]*Limits

type PriceResourceRoots map[string]*Prices

type AccessControlLimits map[string]*LimitResourceRoots

type RateCardLimits map[string]*PriceResourceRoots
