package onappgo

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/godo"
)

const hypervisorZonesBasePath string = "federation/hypervisor_zones/unsubscribed"

// HypervisorZonesService is an interface for interfacing with the HypervisorZone
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/federation/get-list-of-federated-resources
type HypervisorZonesService interface {
	List(context.Context, *ListOptions) ([]HypervisorZone, *Response, error)
	Get(context.Context, int) (*HypervisorZone, *Response, error)
	// Delete(context.Context, int) (*Response, error)
	Delete(context.Context, int, interface{}) (*Transaction, *Response, error)
	// Edit(context.Context, int, *ListOptions) ([]HypervisorZone, *Response, error)
}

// HypervisorZonesServiceOp handles communication with the HypervisorZone related methods of the
// OnApp API.
type HypervisorZonesServiceOp struct {
	client *Client
}

var _ HypervisorZonesService = &HypervisorZonesServiceOp{}

// Certificate -
type Certificate struct {
	ExpireAt string `json:"expire_at,omitempty"`
	Name     string `json:"name,omitempty"`
}

// Certificates -
type Certificates struct {
	Certificate Certificate `json:"certificate,omitempty"`
}

// DataStoreZonePricing -
type DataStoreZonePricing struct {
	DataRead       string `json:"data_read,omitempty"`
	DataWrite      string `json:"data_write,omitempty"`
	DiskSizeMax    string `json:"disk_size_max,omitempty"`
	DiskSizeOff    string `json:"disk_size_off,omitempty"`
	DiskSizeOn     string `json:"disk_size_on,omitempty"`
	InputRequests  string `json:"input_requests,omitempty"`
	OutputRequests string `json:"output_requests,omitempty"`
}

// HypervisorZonePricing -
type HypervisorZonePricing struct {
	CPUMax         string `json:"cpu_max,omitempty"`
	CPUOff         string `json:"cpu_off,omitempty"`
	CPUOn          string `json:"cpu_on,omitempty"`
	CPUPriorityMax string `json:"cpu_priority_max,omitempty"`
	CPUPriorityOff string `json:"cpu_priority_off,omitempty"`
	CPUPriorityOn  string `json:"cpu_priority_on,omitempty"`
	MemoryMax      string `json:"memory_max,omitempty"`
	MemoryOff      string `json:"memory_off,omitempty"`
	MemoryOn       string `json:"memory_on,omitempty"`
}

// NetworkZonePricing -
type NetworkZonePricing struct {
	DataRxed       string `json:"data_rxed,omitempty"`
	DataSent       string `json:"data_sent,omitempty"`
	IPAddressesMax string `json:"ip_addresses_max,omitempty"`
	IPAddressesOff string `json:"ip_addresses_off,omitempty"`
	IPAddressesOn  string `json:"ip_addresses_on,omitempty"`
	PortSpeed      string `json:"port_speed,omitempty"`
	PortSpeedMax   string `json:"port_speed_max,omitempty"`
}

// TierOptions -
type TierOptions struct {
	Backups            bool `json:"backups,bool"`
	DdosProtection     bool `json:"ddos_protection,bool"`
	DNS                bool `json:"dns,bool"`
	Ha                 bool `json:"ha,bool"`
	Ipv6               bool `json:"ipv6,bool"`
	Motion             bool `json:"motion,bool"`
	Replication        bool `json:"replication,bool"`
	SLA                bool `json:"sla,bool"`
	StoragePerformance bool `json:"storage_performance,bool"`
	Templates          bool `json:"templates,bool"`
	WindowsLicense     bool `json:"windows_license,bool"`
}

// UserVirtualServerPricing -
type UserVirtualServerPricing struct {
	AutoScaling            string `json:"auto_scaling,omitempty"`
	AutoScalingMax         string `json:"auto_scaling_max,omitempty"`
	Backup                 string `json:"backup,omitempty"`
	BackupMax              string `json:"backup_max,omitempty"`
	Template               string `json:"template,omitempty"`
	TemplateBackupStore    string `json:"template_backup_store,omitempty"`
	TemplateBackupStoreMax string `json:"template_backup_store_max,omitempty"`
	TemplateMax            string `json:"template_max,omitempty"`
}

// HypervisorZone -
type HypervisorZone struct {
	BandwidthIndex           int                      `json:"bandwidth_index,omitempty"`
	BandwidthScore           int                      `json:"bandwidth_score,omitempty"`
	Certificates             []Certificates           `json:"certificates,omitempty"`
	City                     string                   `json:"city,omitempty"`
	CloudIndex               int                      `json:"cloud_index,omitempty"`
	Country                  string                   `json:"country,omitempty"`
	CPUIndex                 int                      `json:"cpu_index,omitempty"`
	CPUScore                 int                      `json:"cpu_score,omitempty"`
	DataStoreZonePricing     DataStoreZonePricing     `json:"data_store_zone_pricing,omitempty"`
	Description              string                   `json:"description,omitempty"`
	DiskIndex                int                      `json:"disk_index,omitempty"`
	DiskScore                int                      `json:"disk_score,omitempty"`
	FederationID             string                   `json:"federation_id,omitempty"`
	HypervisorZonePricing    HypervisorZonePricing    `json:"hypervisor_zone_pricing,omitempty"`
	Label                    string                   `json:"label,omitempty"`
	Latitude                 float64                  `json:"latitude,omitempty"`
	Longitude                float64                  `json:"longitude,omitempty"`
	NetworkZonePricing       NetworkZonePricing       `json:"network_zone_pricing,omitempty"`
	ProviderName             string                   `json:"provider_name,omitempty"`
	SellerPageURL            string                   `json:"seller_page_url,omitempty"`
	Tier                     string                   `json:"tier,omitempty"`
	TierBandwidthIndex       int                      `json:"tier_bandwidth_index,omitempty"`
	TierCloudIndex           int                      `json:"tier_cloud_index,omitempty"`
	TierCPUIndex             int                      `json:"tier_cpu_index,omitempty"`
	TierDiskIndex            int                      `json:"tier_disk_index,omitempty"`
	TierOptions              TierOptions              `json:"tier_options,omitempty"`
	UptimePercentage         int                      `json:"uptime_percentage,omitempty"`
	UserVirtualServerPricing UserVirtualServerPricing `json:"user_virtual_server_pricing,omitempty"`
}

type hypervisorZonesRoot struct {
	HypervisorZone *HypervisorZone `json:"hypervisor_zone"`
}

// List all HypervisorZones.
func (s *HypervisorZonesServiceOp) List(ctx context.Context, opt *ListOptions) ([]HypervisorZone, *Response, error) {
	path := hypervisorZonesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]HypervisorZone
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]HypervisorZone, len(out))
	for i := range arr {
		arr[i] = out[i]["hypervisor_zone"]
	}

	return arr, resp, err
}

// Get individual HypervisorZone.
func (s *HypervisorZonesServiceOp) Get(ctx context.Context, id int) (*HypervisorZone, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorZonesBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(hypervisorZonesRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.HypervisorZone, resp, err
}

// Delete HypervisorZone.
func (s *HypervisorZonesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Transaction, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorZonesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return nil, resp, err
	}

	filter := struct {
		AssociatedObjectID   int
		AssociatedObjectType string
	}{
		AssociatedObjectID:   id,
		AssociatedObjectType: "HypervisorZone",
	}

	return lastTransaction(ctx, s.client, filter)
}
