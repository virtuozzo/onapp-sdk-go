package onappgo

import (
  "context"
  "net/http"
  "fmt"

  "github.com/digitalocean/godo"
)

const hypervisorZoneBasePath = "federation/hypervisor_zones/unsubscribed"

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

type HypervisorZonePricing struct {
  CPUOn          string `json:"cpu_on,omitempty"`
  CPUOff         string `json:"cpu_off,omitempty"`
  CPUPriorityOn  string `json:"cpu_priority_on,omitempty"`
  CPUPriorityOff string `json:"cpu_priority_off,omitempty"`
  MemoryOn       string `json:"memory_on,omitempty"`
  MemoryOff      string `json:"memory_off,omitempty"`
  CPUMax         string `json:"cpu_max,omitempty"`
  CPUPriorityMax string `json:"cpu_priority_max,omitempty"`
  MemoryMax      string `json:"memory_max,omitempty"`
}

type DataStoreZonePricing struct {
  DiskSizeOn     string `json:"disk_size_on,omitempty"`
  DiskSizeOff    string `json:"disk_size_off,omitempty"`
  DataRead       string `json:"data_read,omitempty"`
  DataWrite      string `json:"data_write,omitempty"`
  InputRequests  string `json:"input_requests,omitempty"`
  OutputRequests string `json:"output_requests,omitempty"`
  DiskSizeMax    string `json:"disk_size_max,omitempty"`
}

type NetworkZonePricing struct {
  IPAddressesOn  string `json:"ip_addresses_on,omitempty"`
  IPAddressesOff string `json:"ip_addresses_off,omitempty"`
  PortSpeed      string `json:"port_speed,omitempty"`
  DataRxed       string `json:"data_rxed,omitempty"`
  DataSent       string `json:"data_sent,omitempty"`
  IPAddressesMax string `json:"ip_addresses_max,omitempty"`
  PortSpeedMax   string `json:"port_speed_max,omitempty"`
}

type UserVirtualServerPricing struct {
  AutoScaling            string `json:"auto_scaling,omitempty"`
  TemplateBackupStore    string `json:"template_backup_store,omitempty"`
  Backup                 string `json:"backup,omitempty"`
  Template               string `json:"template,omitempty"`
  AutoScalingMax         string `json:"auto_scaling_max,omitempty"`
  TemplateBackupStoreMax string `json:"template_backup_store_max,omitempty"`
  BackupMax              string `json:"backup_max,omitempty"`
  TemplateMax            string `json:"template_max,omitempty"`
}

type TierOptions struct {
  Ha                 bool `json:"ha,bool"`
  SLA                bool `json:"sla,bool"`
  StoragePerformance bool `json:"storage_performance,bool"`
  Backups            bool `json:"backups,bool"`
  Templates          bool `json:"templates,bool"`
  WindowsLicense     bool `json:"windows_license,bool"`
  DdosProtection     bool `json:"ddos_protection,bool"`
  Ipv6               bool `json:"ipv6,bool"`
  DNS                bool `json:"dns,bool"`
  Motion             bool `json:"motion,bool"`
  Replication        bool `json:"replication,bool"`
}

type HypervisorZone struct {
  Label                    string                   `json:"label,omitempty"`
  ProviderName             string                   `json:"provider_name,omitempty"`
  SellerPageURL            string                   `json:"seller_page_url,omitempty"`
  Description              string                   `json:"description,omitempty"`
  FederationID             string                   `json:"federation_id,omitempty"`
  Country                  string                   `json:"country,omitempty"`
  City                     string                   `json:"city,omitempty"`
  UptimePercentage         int                      `json:"uptime_percentage,omitempty"`
  Tier                     string                   `json:"tier,omitempty"`
  Latitude                 float64                  `json:"latitude,omitempty"`
  Longitude                float64                  `json:"longitude,omitempty"`
  CPUScore                 int                      `json:"cpu_score,omitempty"`
  CPUIndex                 int                      `json:"cpu_index,omitempty"`
  BandwidthScore           int                      `json:"bandwidth_score,omitempty"`
  BandwidthIndex           int                      `json:"bandwidth_index,omitempty"`
  DiskScore                int                      `json:"disk_score,omitempty"`
  DiskIndex                int                      `json:"disk_index,omitempty"`
  CloudIndex               int                      `json:"cloud_index,omitempty"`
  TierCPUIndex             int                      `json:"tier_cpu_index,omitempty"`
  TierDiskIndex            int                      `json:"tier_disk_index,omitempty"`
  TierBandwidthIndex       int                      `json:"tier_bandwidth_index,omitempty"`
  TierCloudIndex           int                      `json:"tier_cloud_index,omitempty"`
  Certificates             []interface{}            `json:"certificates,omitempty"`
  HypervisorZonePricing    HypervisorZonePricing    `json:"hypervisor_zone_pricing,omitempty"`
  DataStoreZonePricing     DataStoreZonePricing     `json:"data_store_zone_pricing,omitempty"`
  NetworkZonePricing       NetworkZonePricing       `json:"network_zone_pricing,omitempty"`
  UserVirtualServerPricing UserVirtualServerPricing `json:"user_virtual_server_pricing,omitempty"`
  TierOptions              TierOptions              `json:"tier_options,omitempty"`
}

type hypervisorZonesRoot struct {
  HypervisorZone  *HypervisorZone  `json:"hypervisor_zone"`
}

// List all HypervisorZones.
func (s *HypervisorZonesServiceOp) List(ctx context.Context, opt *ListOptions) ([]HypervisorZone, *Response, error) {
  path := hypervisorZoneBasePath + apiFormat
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

  path := fmt.Sprintf("%s/%d%s", hypervisorZoneBasePath, id, apiFormat)

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

  path := fmt.Sprintf("%s/%d%s", hypervisorZoneBasePath, id, apiFormat)
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

  filter := struct{
    AssociatedObjectID    int
    AssociatedObjectType  string
  }{
    AssociatedObjectID    : id,
    AssociatedObjectType  : "HypervisorZone",
  }

  return lastTransaction(ctx, s.client, filter)
  // return lastTransaction(ctx, s.client, id, "HypervisorZone")
}

// Debug - print formatted HypervisorZone structure
func (obj HypervisorZone) Debug() {
  fmt.Printf("           Label: %s\n", obj.Label)
  fmt.Printf("    ProviderName: %s\n", obj.ProviderName)
  fmt.Printf("         Country: %s\n", obj.Country)
  fmt.Printf("            City: %s\n", obj.City)
  fmt.Printf("UptimePercentage: %d\n", obj.UptimePercentage)
  fmt.Printf("    FederationID: %s\n", obj.FederationID)
}
