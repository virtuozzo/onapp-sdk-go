package onappgo

// Represent common structures for onappgo package

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

type LimitResourceRoots map[string]*Limits

type PriceResourceRoots map[string]*Prices

type AccessControlLimits map[string]*LimitResourceRoots

type RateCardLimits map[string]*PriceResourceRoots
