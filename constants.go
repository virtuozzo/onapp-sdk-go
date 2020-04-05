package onappgo

// Features
const (
	STORING                         = 1 + iota // "storing"
	SOLIDFIRE_STORING                          // "solidfire_storing"
	VDC_STORING                                // "vdc_storing"
	HYPERVISOR_STORING                         // "hypervisor_storing"
	NETWORKING                                 // "networking"
	EXTERNAL_NETWORKING                        // "external_networking"
	REGULAR_SERVER                             // "regular_server"
	VIRTUAL_SERVER                             // "virtual_server"
	PRECONFIGURED_SERVER                       // "preconfigured_server"
	AUTOSCALED_SERVER                          // "autoscaled_server"
	APPLICATION_SERVER                         // "application_server"
	ACCELERATED_SERVER                         // "accelerated_server"
	BAREMETAL_SERVER                           // "baremetal_server"
	CONTAINER_SERVER                           // "container_server"
	SMART_SERVER                               // "smart_server"
	VPC_SERVER                                 // "vpc_server"
	DRAAS                                      // "draas"
	CDN                                        // "cdn"
	REGULAR_BACKUP_ON_BACKUP_SERVER            // "regular_backup_on_backup_server"

	REGULAR_TEMPLATE_ON_BACKUP_SERVER = 21 + iota // "regular_template_on_backup_server"
	ISO_TEMPLATE                                  // "iso_template"
	OVA_TEMPLATE                                  // "ova_template"
	ALLOCATION_POOL                               // "allocation_pool"
	RESERVATION_POOL                              // "reservation_pool"
	ALLOCATION_VAPP                               // "allocation_vapp"
	RECIPE                                        // "recipe"
	SERVICE_ADDON                                 // "service_addon"
	REGULAR_BACKUP_ON_HYPERVISOR                  // "regular_backup_on_hypervisor"

	REGULAR_TEMPLATE_ON_HYPERVISOR    = 31 + iota // "regular_template_on_hypervisor"
	TEMPLATE_USAGE                                // "template_usage"
	RECOVERY_POINT_ON_BACKUP_RESOURCE             // "recovery_point_on_backup_resource"
	BACKUP_RESOURCE_USAGE                         // "backup_resource_usage"

	// Monthly Statistics features
	SUM_OF_HOURLY_ADJUSTMENT        = 95 + iota // "sum_of_hourly_adjustment"
	SUM_OF_HOURLY_STAT_NON_ARCHIVED             // "sum_of_hourly_stat_non_archived"
	SUM_OF_HOURLY_STAT                          // "sum_of_hourly_stat"
	SUM_OF_HOURLY_STAT_ARCHIVED                 // "sum_of_hourly_stat_archived"
	SUBSCRIPTION                                // "subscription"
)

// Metrics
// const (
//   COUNT                   = 1 + iota  // "count"
//   CPUS                                // "cpus"
//   CPU_SHARES                          // "cpu_shares"
//   CPU_PRIORITY                        // "cpu_priority"
//   CPU_UNITS                           // "cpu_units"
//   MEMORY                              // "memory"
//   MEMORY_USED                         // "memory_used"
//   DISK_SIZE                           // "disk_size"
//   DISK_SIZE_USED                      // "disk_size_used"
//   UNLIMITED_DISK_SIZE                 // "disk_size_unlimited"
//   DATA_READ                           // "data_read"
//   DATA_WRITTEN                        // "data_written"
//   INPUT_REQUESTS                      // "input_requests"
//   OUTPUT_REQUESTS                     // "output_requests"
//   IP_ADDRESSES                        // "ip_addresses"
//   PORT_SPEED                          // "port_speed"
//   DATA_SENT                           // "data_sent"
//   DATA_RECEIVED                       // "data_received"
//   OVERUSED_BANDWIDTH                  // "overused_bandwidth"
//   IOPS                                // "iops"
//   NODES                               // "nodes"
//   CPU_USED                            // "cpu_used"
//   VCPU_SPEED                          // "vcpu_speed"
//   CPU_GUARANTEED                      // "cpu_guaranteed"
//   MEMORY_GUARANTEED                   // "memory_guaranteed"
//   CPU_ALLOCATION                      // "cpu_allocation"
//   MEMORY_ALLOCATION                   // "memory_allocation"
//   REGULAR_CPU_QUOTA                   // "regular_cpu_quota"
//   REGULAR_MEMORY_QUOTA                // "regular_memory_quota"
//   UNLIMITED_CPU_QUOTA                 // "unlimited_cpu_quota"
//   UNLIMITED_MEMORY_QUOTA              // "unlimited_memory_quota"
//   DEPLOYED_EDGE_GATEWAYS              // "deployed_edge_gateways"
//   DEPLOYED_ORG_NETWORKS               // "deployed_org_networks"
//   FAST_PROVISIONING_SET               // "fast_provisioning_set"
//   THIN_PROVISIONING_SET               // "thin_provisioning_set"
//   VS_COUNT                            // "vs_count"
//   VS_LIMIT                            // "vs_limit"
//   CPU_LIMIT                           // "cpu_limit"
//   MEMORY_LIMIT                        // "memory_limit"
//   VRAM                                // "vram"
//   CPU_TIME                            // "cpu_time"
//   UNLIMITED_PORT_SPEED                // "unlimited_port_speed"
//   DATA_CACHED                         // "data_cached"
//   DATA_NON_CACHED                     // "data_non_cached"
//   CDN_BANDWIDTH                       // "cdn_bandwidth"
// )

// Targets
const (
	COMPUTE_ZONE        = 1 + iota // "compute_zone"
	DATA_STORE_ZONE                // "data_store_zone"
	NETWORK_ZONE                   // "network_zone"
	BACKUP_SERVER_ZONE             // "backup_server_zone"
	INSTANCE_PACKAGE               // "instance_package"
	TEMPLATE_GROUP                 // "template_group"
	EDGE_GROUP                     // "edge_group"
	RECIPE_GROUP                   // "recipe_group"
	SERVICE_ADDON_GROUP            // "service_addon_group"
)

// ResourceRoots
const (
	COMPUTE_ZONE_RESOURCE              = "compute_zone_resource"
	DATA_STORE_ZONE_RESOURCE           = "data_store_zone_resource"
	NETWORK_ZONE_RESOURCE              = "network_zone_resource"
	BACKUP_SERVER_ZONE_RESOURCE        = "backup_server_zone_resource"
	VIRTUAL_SERVERS_RESOURCE           = "virtual_servers_resource"
	AUTOSCALED_SERVERS_RESOURCE        = "autoscaled_servers_resource"
	COMPUTE_RESOURCE_STORING_RESOURCE  = "compute_resource_storing_resource"
	BACKUPS_RESOURCE                   = "backups_resource"
	TEMPLATES_RESOURCE                 = "templates_resource"
	ISO_TEMPLATES_RESOURCE             = "iso_templates_resource"
	ACCELERATED_SERVERS_RESOURCE       = "accelerated_servers_resource"
	APPLICATION_SERVERS_RESOURCE       = "application_servers_resource"
	CONTAINER_SERVERS_RESOURCE         = "container_servers_resource"
	DRAAS_RESOURCE                     = "draas_resource"
	SOLIDFIRE_DATA_STORE_ZONE_RESOURCE = "solidfire_data_store_zone_resource"
	PRECONFIGURED_SERVERS_RESOURCE     = "preconfigured_servers_resource"
	EDGE_GROUPS_RESOURCE               = "edge_groups_resource"
	RECIPE_GROUPS_RESOURCE             = "recipe_groups_resource"
	TEMPLATE_GROUPS_RESOURCE           = "template_groups_resource"
	SERVICE_ADDON_GROUPS_RESOURCE      = "service_addon_groups_resource"
	TEMPLATE_RESOURCE                  = "template_resource"
	SERVICE_ADDON_RESOURCE             = "service_addon_resource"
	BARE_METAL_SERVERS_RESOURCE        = "bare_metal_servers_resource"
	SMART_SERVERS_RESOURCE             = "smart_servers_resource"
	ORCHESTRATION_MODEL_RESOURCE       = "orchestration_model_resource"
	BACKUP_RESOURCE_ZONE_RESOURCE      = "backup_resource_zone_resource"
	CDN_BANDWIDTH_RESOURCE             = "cdn_bandwidth_resource"
)

// ServerTypes
const (
	VIRTUAL        = "virtual"
	VPC            = "vpc"
	SMART          = "smart"
	BARE_METAL     = "baremetal"
	OCM            = "ocm"
	OTHER          = "other"
	INFRASTRUCTURE = "infrastructure"
	SUNLIGHT       = "sunlight"
)

var ServerTypes = []string{
	VIRTUAL,
	VPC,
	SMART,
	BARE_METAL,
	OCM,
	OTHER,
	INFRASTRUCTURE,
	SUNLIGHT,
}

var TimingStrategy = []string{
	"hourly",
	"monthly",
}

// type ServerTypesRestrictions map[string][]string

var SERVER_TYPES_RESTRICTIONS = map[string][]string{
	VIRTUAL: []string{
		COMPUTE_ZONE_RESOURCE,
		DATA_STORE_ZONE_RESOURCE,
		NETWORK_ZONE_RESOURCE,
		BACKUP_SERVER_ZONE_RESOURCE,
		VIRTUAL_SERVERS_RESOURCE,
		AUTOSCALED_SERVERS_RESOURCE,
		COMPUTE_RESOURCE_STORING_RESOURCE,
		BACKUPS_RESOURCE,
		TEMPLATES_RESOURCE,
		ISO_TEMPLATES_RESOURCE,
		ACCELERATED_SERVERS_RESOURCE,
		APPLICATION_SERVERS_RESOURCE,
		CONTAINER_SERVERS_RESOURCE,
		DRAAS_RESOURCE,
		SOLIDFIRE_DATA_STORE_ZONE_RESOURCE,
		PRECONFIGURED_SERVERS_RESOURCE,
	},

	VPC: []string{
		VIRTUAL_SERVERS_RESOURCE,
		APPLICATION_SERVERS_RESOURCE,
		COMPUTE_ZONE_RESOURCE,
		DATA_STORE_ZONE_RESOURCE,
		NETWORK_ZONE_RESOURCE,
	},

	SMART: []string{
		SMART_SERVERS_RESOURCE,
		COMPUTE_ZONE_RESOURCE,
		DATA_STORE_ZONE_RESOURCE,
		NETWORK_ZONE_RESOURCE,
		BACKUP_SERVER_ZONE_RESOURCE,
		COMPUTE_RESOURCE_STORING_RESOURCE,
		BACKUPS_RESOURCE,
	},

	BARE_METAL: []string{
		NETWORK_ZONE_RESOURCE,
		BARE_METAL_SERVERS_RESOURCE,
		COMPUTE_ZONE_RESOURCE,
	},

	OTHER: []string{
		EDGE_GROUPS_RESOURCE,
		RECIPE_GROUPS_RESOURCE,
		TEMPLATE_GROUPS_RESOURCE,
		SERVICE_ADDON_GROUPS_RESOURCE,
		TEMPLATE_RESOURCE,
		SERVICE_ADDON_RESOURCE,
		ORCHESTRATION_MODEL_RESOURCE,
		BACKUP_RESOURCE_ZONE_RESOURCE,
	},

	SUNLIGHT: []string{
		COMPUTE_ZONE_RESOURCE,
		DATA_STORE_ZONE_RESOURCE,
		NETWORK_ZONE_RESOURCE,
		VIRTUAL_SERVERS_RESOURCE,
		PRECONFIGURED_SERVERS_RESOURCE,
	},
}
