package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

// Xen/KVM, VMware - CRUD
// CloudBoot, Smart CloudBoot, Baremetal CloudBoot - Get, Delete
const hypervisorGroupsBasePath string = "settings/hypervisor_zones"
const listOfAttachedComputeResources string = "settings/hypervisor_zones/%d/hypervisors"

// HypervisorGroupsService is an interface for interfacing with the Compute Zone
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/compute-zones
type HypervisorGroupsService interface {
	List(context.Context, *ListOptions) ([]HypervisorGroup, *Response, error)
	Get(context.Context, int) (*HypervisorGroup, *Response, error)
	Create(context.Context, *HypervisorGroupCreateRequest) (*HypervisorGroup, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *HypervisorGroupEditRequest) (*Response, error)

	ListOfAttachedComputeResources(context.Context, int) ([]Hypervisor, *Response, error)
}

// HypervisorGroupsServiceOp handles communication with the Compute Zone
// related methods of the OnApp API.
type HypervisorGroupsServiceOp struct {
	client *Client
}

var _ HypervisorGroupsService = &HypervisorGroupsServiceOp{}

// HypervisorGroup represent Compute Zone of the OnApp API
type HypervisorGroup struct {
	AdditionalFields            []AdditionalFields `json:"additional_fields,omitempty"`
	Closed                      bool               `json:"closed,bool"`
	CPUFlags                    []string           `json:"cpu_flags,omitempty"`
	CPUFlagsEnabled             bool               `json:"cpu_flags_enabled,bool"`
	CPUModelConfiguration       string             `json:"cpu_model_configuration"`
	CPUUnits                    int                `json:"cpu_units,omitempty"`
	CreatedAt                   string             `json:"created_at,omitempty"`
	CustomConfig                string             `json:"custom_config,omitempty"`
	DatacenterID                int                `json:"datacenter_id,omitempty"`
	DraasID                     int                `json:"draas_id,omitempty"`
	FailoverTimeout             int                `json:"failover_timeout,omitempty"`
	FederationEnabled           bool               `json:"federation_enabled,bool"`
	FederationID                string             `json:"federation_id,omitempty"`
	HypervisorID                int                `json:"hypervisor_id,omitempty"`
	ID                          int                `json:"id,omitempty"`
	Identifier                  string             `json:"identifier,omitempty"`
	Label                       string             `json:"label,omitempty"`
	LocationGroupID             int                `json:"location_group_id,omitempty"`
	MaxHostCPU                  int                `json:"max_host_cpu,omitempty"`
	MaxHostFreeMemory           int                `json:"max_host_free_memory,omitempty"`
	NetworkFailure              bool               `json:"network_failure,bool"`
	PreconfiguredOnly           bool               `json:"preconfigured_only,bool"`
	PreferLocalReads            bool               `json:"prefer_local_reads,bool"`
	ProviderName                string             `json:"provider_name,omitempty"`
	ProviderVdcID               int                `json:"provider_vdc_id,omitempty"`
	RecoveryType                string             `json:"recovery_type,omitempty"`
	ReleaseResourceType         string             `json:"release_resource_type,omitempty"`
	RunSysprep                  bool               `json:"run_sysprep,bool"`
	ScheduledForDeletion        string             `json:"scheduled_for_deletion,omitempty"`
	ServerType                  string             `json:"server_type,omitempty"`
	StorageChannel              int                `json:"storage_channel,omitempty"`
	SupplierProvider            string             `json:"supplier_provider,omitempty"`
	SupplierVersion             string             `json:"supplier_version,omitempty"`
	SupportsVirtualServerMotion bool               `json:"supports_virtual_server_motion,bool"`
	Tier                        string             `json:"tier,omitempty"`
	Traded                      bool               `json:"traded,bool"`
	UpdatedAt                   string             `json:"updated_at,omitempty"`
}

// HypervisorGroupCreateRequest represents a request to create a Compute Zone
type HypervisorGroupCreateRequest struct {
	CPUFlagsEnabled     bool   `json:"cpu_flags_enabled,bool"`
	CPUUnits            int    `json:"cpu_units,omitempty"`
	CustomConfig        string `json:"custom_config,omitempty"`
	DefaultGateway      string `json:"default_gateway,omitempty"`
	FailoverTimeout     int    `json:"failover_timeout,omitempty"`
	Label               string `json:"label,omitempty"`
	LocationGroupID     int    `json:"location_group_id,omitempty"`
	MaxVmsStartAtOnce   int    `json:"max_vms_start_at_once,omitempty"`
	PreconfiguredOnly   bool   `json:"preconfigured_only,bool"`
	RecoveryType        string `json:"recovery_type,omitempty"`
	ReleaseResourceType string `json:"release_resource_type,omitempty"`
	RunSysprep          bool   `json:"run_sysprep,bool"`
	ServerType          string `json:"server_type,omitempty"`
	Vlan                string `json:"vlan,omitempty"`
}

// HypervisorGroupEditRequest represents a request to edit a Compute Zone
type HypervisorGroupEditRequest struct {
	CPUFlagsEnabled     bool   `json:"cpu_flags_enabled,bool"`
	CPUUnits            int    `json:"cpu_units,omitempty"`
	CustomConfig        string `json:"custom_config,omitempty"`
	FailoverTimeout     int    `json:"failover_timeout,omitempty"`
	Label               string `json:"label,omitempty"`
	LocationGroupID     int    `json:"location_group_id,omitempty"`
	MaxVmsStartAtOnce   int    `json:"max_vms_start_at_once,omitempty"`
	OnlyStartedVms      int    `json:"only_started_vms,omitempty"`
	PreconfiguredOnly   bool   `json:"preconfigured_only,bool"`
	PreferLocalReads    int    `json:"prefer_local_reads,omitempty"`
	RecoveryType        string `json:"recovery_type,omitempty"`
	ReleaseResourceType string `json:"release_resource_type,omitempty"`
	RunSysprep          bool   `json:"run_sysprep,bool"`
	ServerType          string `json:"server_type,omitempty"`
	UpdateCPUUnits      int    `json:"update_cpu_units,omitempty"`
	СPUGuarantee        int    `json:"cpu_guarantee,omitempty"`
}

type hypervisorGroupCreateRequestRoot struct {
	HypervisorGroupCreateRequest *HypervisorGroupCreateRequest `json:"hypervisor_group"`
}

type hypervisorGroupRoot struct {
	HypervisorGroup *HypervisorGroup `json:"hypervisor_group"`
}

func (d HypervisorGroupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all HypervisorGroup.
func (s *HypervisorGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]HypervisorGroup, *Response, error) {
	path := hypervisorGroupsBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]HypervisorGroup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]HypervisorGroup, len(out))
	for i := range arr {
		arr[i] = out[i]["hypervisor_group"]
	}

	return arr, resp, err
}

// Get individual HypervisorGroup.
func (s *HypervisorGroupsServiceOp) Get(ctx context.Context, id int) (*HypervisorGroup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(hypervisorGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.HypervisorGroup, resp, err
}

// Create HypervisorGroup.
func (s *HypervisorGroupsServiceOp) Create(ctx context.Context, createRequest *HypervisorGroupCreateRequest) (*HypervisorGroup, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("HypervisorZone createRequest", "cannot be nil")
	}

	path := hypervisorGroupsBasePath + apiFormat

	rootRequest := &hypervisorGroupCreateRequestRoot{
		HypervisorGroupCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("HypervisorGroup [Create] req: ", req)

	root := new(hypervisorGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.HypervisorGroup, resp, err
}

// Delete HypervisorGroup.
func (s *HypervisorGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorGroupsBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("HypervisorGroup [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit HypervisorGroup
func (s *HypervisorGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *HypervisorGroupEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("HypervisorZone [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", hypervisorGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("HypervisorGroup [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// ListOfAttachedComputeResources -
func (s *HypervisorGroupsServiceOp) ListOfAttachedComputeResources(ctx context.Context, hvgID int) ([]Hypervisor, *Response, error) {
	if hvgID < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(listOfAttachedComputeResources, hvgID) + apiFormat
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}
	log.Println("HypervisorGroup [ListOfAttachedComputeResources]  req: ", req)

	var out []map[string]Hypervisor
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Hypervisor, len(out))
	for i := range arr {
		arr[i] = out[i]["hypervisor"]
	}

	return arr, resp, err
}
