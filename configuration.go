package onappgo

import (
	"context"
	"log"
	"net/http"
)

const configurationsBasePath string = "settings/configuration"
const configurationsEditBasePath string = "settings"

// ConfigurationsService is an interface for interfacing with the Configurations
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/system-configuration
type ConfigurationsService interface {
	Get(context.Context) (*Configuration, *Response, error)
	Edit(context.Context, *Configuration) (*Response, error)
}

// ConfigurationsServiceOp handles communication with the Configuration related methods of the
// OnApp API.
type ConfigurationsServiceOp struct {
	client *Client
}

var _ ConfigurationsService = &ConfigurationsServiceOp{}

// Configuration - represent configuration settings of OnApp installation.
type Configuration struct {
	ActionGlobalLockExpirationTimeout     int      `json:"action_global_lock_expiration_timeout,omitempty"`
	ActionGlobalLockRetryDelay            int      `json:"action_global_lock_retry_delay,omitempty"`
	AdapterOpenConnectionTimeout          int      `json:"adapter_open_connection_timeout,omitempty"`
	AdapterResponseTimeout                int      `json:"adapter_response_timeout,omitempty"`
	AjaxLogUpdateInterval                 int      `json:"ajax_log_update_interval,omitempty"`
	AjaxPaginationUpdateTime              int      `json:"ajax_pagination_update_time,omitempty"`
	AjaxPowerUpdateTime                   int      `json:"ajax_power_update_time,omitempty"`
	AllowAdvancedVsManagement             bool     `json:"allow_advanced_vs_management,bool"`
	AllowConnectAws                       bool     `json:"allow_connect_aws,bool"`
	AllowHypervisorPasswordEncryption     bool     `json:"allow_hypervisor_password_encryption,bool"`
	AllowIncrementalBackups               bool     `json:"allow_incremental_backups,bool"`
	AllowInitialRootPasswordEncryption    bool     `json:"allow_initial_root_password_encryption,bool"`
	AllowStartVmsWithOneIP                bool     `json:"allow_start_vms_with_one_ip,bool"`
	AllowToCollectErrors                  bool     `json:"allow_to_collect_errors,bool"`
	AmountOfServiceInstances              int      `json:"amount_of_service_instances,omitempty"`
	AppName                               string   `json:"app_name,omitempty"`
	ArchiveStatsPeriod                    int      `json:"archive_stats_period,omitempty"`
	BackupConvertCoefficient              float64  `json:"backup_convert_coefficient,omitempty"`
	BackupsPath                           string   `json:"backups_path,omitempty"`
	BackupTakerDelay                      int      `json:"backup_taker_delay,omitempty"`
	BillingTransactionRunnerDelay         int      `json:"billing_transaction_runner_delay,omitempty"`
	BlockSize                             int      `json:"block_size,omitempty"`
	CaptureVappTimeout                    int      `json:"capture_vapp_timeout,omitempty"`
	CdnMaxResultsPerGetPage               int      `json:"cdn_max_results_per_get_page,omitempty"`
	CdnSyncDelay                          int      `json:"cdn_sync_delay,omitempty"`
	CloudBootDomainNameServers            string   `json:"cloud_boot_domain_name_servers,omitempty"`
	CloudBootEnabled                      bool     `json:"cloud_boot_enabled,bool"`
	CloudBootTarget                       string   `json:"cloud_boot_target,omitempty"`
	ClusterMonitorDelay                   int      `json:"cluster_monitor_delay,omitempty"`
	ComposeVappTimeout                    int      `json:"compose_vapp_timeout,omitempty"`
	CreateEdgeGatewayTimeout              int      `json:"create_edge_gateway_timeout,omitempty"`
	CreateSnapshotTimeout                 int      `json:"create_snapshot_timeout,omitempty"`
	CreateVappTemplateTimeout             int      `json:"create_vapp_template_timeout,omitempty"`
	CreateVdcTimeout                      int      `json:"create_vdc_timeout,omitempty"`
	DashboardAPIAccessToken               string   `json:"dashboard_api_access_token,omitempty"`
	DashboardStats                        []string `json:"dashboard_stats,omitempty"`
	DataPath                              string   `json:"data_path,omitempty"`
	DefaultAccelerationPolicy             bool     `json:"default_acceleration_policy,bool"`
	DefaultCustomTheme                    string   `json:"default_custom_theme,omitempty"`
	DefaultFirewallPolicy                 string   `json:"default_firewall_policy,omitempty"`
	DefaultImageTemplate                  int      `json:"default_image_template,omitempty"`
	DefaultTimeout                        int      `json:"default_timeout,omitempty"`
	DefaultVirshConsolePolicy             bool     `json:"default_virsh_console_policy,bool"`
	DeleteTemplateSourceAfterInstall      bool     `json:"delete_template_source_after_install,bool"`
	DisableBilling                        bool     `json:"disable_billing,bool"`
	DisableHypervisorFailover             bool     `json:"disable_hypervisor_failover,bool"`
	DNSEnabled                            bool     `json:"dns_enabled,bool"`
	DraasEnabled                          bool     `json:"draas_enabled,bool"`
	DraasShadowSSHPort                    int      `json:"draas_shadow_ssh_port,omitempty"`
	DraasShadowVpnPort                    int      `json:"draas_shadow_vpn_port,omitempty"`
	DraasVpnCidrBlock                     string   `json:"draas_vpn_cidr_block,omitempty"`
	DropFirewallPolicyAllowedIps          string   `json:"drop_firewall_policy_allowed_ips,omitempty"`
	EmailDeliveryMethod                   string   `json:"email_delivery_method,omitempty"`
	EnableDailyStorageReport              bool     `json:"enable_daily_storage_report,bool"`
	EnableDownloadTimeout                 int      `json:"enable_download_timeout"`
	EnableHourlyStorageReport             bool     `json:"enable_hourly_storage_report,bool"`
	EnableNotifications                   bool     `json:"enable_notifications,bool"`
	EnableSuperAdminPermissions           bool     `json:"enable_super_admin_permissions,bool"`
	EnforceRedundancy                     bool     `json:"enforce_redundancy,bool"`
	FederationTrustsOnlyPrivate           bool     `json:"federation_trusts_only_private,bool"`
	ForceSamlLoginOnly                    bool     `json:"force_saml_login_only,bool"`
	GenerateComment                       string   `json:"generate_comment,omitempty"`
	GracefulStopTimeout                   int      `json:"graceful_stop_timeout,omitempty"`
	GuestWaitTimeBeforeDestroy            int      `json:"guest_wait_time_before_destroy,omitempty"`
	HaEnabled                             bool     `json:"ha_enabled,bool"`
	HypervisorLiveTimes                   int      `json:"hypervisor_live_times,omitempty"`
	HypervisorMonitorDelay                int      `json:"hypervisor_monitor_delay,omitempty"`
	InfinibandCloudBootEnabled            bool     `json:"infiniband_cloud_boot_enabled,bool"`
	InstancePackagesThresholdNum          int      `json:"instance_packages_threshold_num,omitempty"`
	InstantiateVappTemplateTimeout        int      `json:"instantiate_vapp_template_timeout,omitempty"`
	InstantStatsPeriod                    int      `json:"instant_stats_period,omitempty"`
	InterHypervisorBalanceThresholdRatio  int      `json:"inter_hypervisor_balance_threshold_ratio,omitempty"`
	IntraHypervisorBalanceThresholdRatio  int      `json:"intra_hypervisor_balance_threshold_ratio,omitempty"`
	IoLimitingEnabled                     bool     `json:"io_limiting_enabled,bool"`
	IPAddressReservationTime              int      `json:"ip_address_reservation_time,omitempty"`
	IPHistoryKeepPeriod                   int      `json:"ip_history_keep_period,omitempty"`
	IpsAllowedForLogin                    string   `json:"ips_allowed_for_login,omitempty"`
	IsArchiveStatsEnabled                 bool     `json:"is_archive_stats_enabled,bool"`
	IscsiPortAvailabilityCheckTimeout     int      `json:"iscsi_port_availability_check_timeout,omitempty"`
	IsolatedLicense                       bool     `json:"isolated_license,bool"`
	IsoPathOnCp                           string   `json:"iso_path_on_cp,omitempty"`
	IsoPathOnHv                           string   `json:"iso_path_on_hv,omitempty"`
	KvmAvailableFreeMemoryPercentage      int      `json:"kvm_available_free_memory_percentage,omitempty"`
	KvmMaxMemoryRate                      int      `json:"kvm_max_memory_rate,omitempty"`
	LicenseKey                            string   `json:"license_key,omitempty"`
	Localdomain                           string   `json:"localdomain,omitempty"`
	Locales                               []string `json:"locales,omitempty"`
	LogCleanupEnabled                     bool     `json:"log_cleanup_enabled,bool"`
	LogCleanupPeriod                      int      `json:"log_cleanup_period,omitempty"`
	LogLevel                              string   `json:"log_level,omitempty"`
	MaxCPUQuota                           int      `json:"max_cpu_quota,omitempty"`
	MaximumPendingTasks                   int      `json:"maximum_pending_tasks,omitempty"`
	MaxIPAddressesToAssignSimultaneously  int      `json:"max_ip_addresses_to_assign_simultaneously,omitempty"`
	MaxMemoryQuota                        int      `json:"max_memory_quota"`
	MaxMemoryRatio                        int      `json:"max_memory_ratio,omitempty"`
	MaxNetworkInterfacePortSpeed          int      `json:"max_network_interface_port_speed,omitempty"`
	MaxUploadSize                         int      `json:"max_upload_size,omitempty"`
	MigrationRateLimit                    int      `json:"migration_rate_limit,omitempty"`
	MonitisAccount                        string   `json:"monitis_account,omitempty"`
	MonitisApikey                         string   `json:"monitis_apikey,omitempty"`
	MonitisPath                           string   `json:"monitis_path,omitempty"`
	MysqlBillingTransactionRetries        string   `json:"mysql_billing_transaction_retries,omitempty"`
	NfsRootIP                             string   `json:"nfs_root_ip,omitempty"`
	NotificationSubjectPrefix             string   `json:"notification_subject_prefix,omitempty"`
	NsxPollingInterval                    int      `json:"nsx_polling_interval,omitempty"`
	NumberOfNotificationsToShow           int      `json:"number_of_notifications_to_show,omitempty"`
	OvaPath                               string   `json:"ova_path,omitempty"`
	PaginationDashboardPagesLimit         int      `json:"pagination_dashboard_pages_limit,omitempty"`
	PaginationMaxItemsLimit               int      `json:"pagination_max_items_limit,omitempty"`
	PartitionAlignOffset                  int      `json:"partition_align_offset,omitempty"`
	PasswordEnforceComplexity             bool     `json:"password_enforce_complexity,bool"`
	PasswordExpiry                        int      `json:"password_expiry,omitempty"`
	PasswordForceUnique                   bool     `json:"password_force_unique,bool"`
	PasswordHistoryLength                 int      `json:"password_history_length,omitempty"`
	PasswordLettersNumbers                bool     `json:"password_letters_numbers,bool"`
	PasswordLockoutAttempts               int      `json:"password_lockout_attempts,omitempty"`
	PasswordMinimumLength                 int      `json:"password_minimum_length,omitempty"`
	PasswordProtectionForDeleting         bool     `json:"password_protection_for_deleting,bool"`
	PasswordSymbols                       bool     `json:"password_symbols,bool"`
	PasswordUpperLowercase                bool     `json:"password_upper_lowercase,bool"`
	PingVmsBeforeInitFailover             bool     `json:"ping_vms_before_init_failover,bool"`
	PowerOffTimeout                       int      `json:"power_off_timeout,omitempty"`
	PowerOnTimeout                        int      `json:"power_on_timeout,omitempty"`
	PreferLocalReads                      bool     `json:"prefer_local_reads,bool"`
	QemuAttachDeviceDelay                 int      `json:"qemu_attach_device_delay,omitempty"`
	QemuDetachDeviceDelay                 int      `json:"qemu_detach_device_delay,omitempty"`
	RabbitmqHost                          string   `json:"rabbitmq_host,omitempty"`
	RabbitmqLogin                         string   `json:"rabbitmq_login,omitempty"`
	RabbitmqPassword                      string   `json:"rabbitmq_password,omitempty"`
	RabbitmqPort                          int      `json:"rabbitmq_port,omitempty"`
	RabbitmqVhost                         string   `json:"rabbitmq_vhost,omitempty"`
	RebootTimeout                         int      `json:"reboot_timeout,omitempty"`
	RecipeTmpDir                          string   `json:"recipe_tmp_dir,omitempty"`
	RecomposeVappTimeout                  int      `json:"recompose_vapp_timeout,omitempty"`
	RemoteAccessSessionLastPort           int      `json:"remote_access_session_last_port,omitempty"`
	RemoteAccessSessionStartPort          int      `json:"remote_access_session_start_port,omitempty"`
	RemoveBackupsOnDestroyVM              bool     `json:"remove_backups_on_destroy_vm,bool"`
	ResetTimeout                          int      `json:"reset_timeout,omitempty"`
	RsyncOptionAcls                       bool     `json:"rsync_option_acls,bool"`
	RsyncOptionXattrs                     bool     `json:"rsync_option_xattrs,bool"`
	RunRecipeOnVsSleepSeconds             int      `json:"run_recipe_on_vs_sleep_seconds,omitempty"`
	ScheduleRunnerDelay                   int      `json:"schedule_runner_delay,omitempty"`
	ServiceAccountName                    string   `json:"service_account_name,omitempty"`
	SessionTimeout                        int      `json:"session_timeout,omitempty"`
	ShowIPAddressSelectionForNewVM        bool     `json:"show_ip_address_selection_for_new_vm,bool"`
	ShowNewWizard                         bool     `json:"show_new_wizard,bool"`
	ShutdownTimeout                       int      `json:"shutdown_timeout,omitempty"`
	SimultaneousBackups                   int      `json:"simultaneous_backups,omitempty"`
	SimultaneousBackupsPerBackupServer    int      `json:"simultaneous_backups_per_backup_server,omitempty"`
	SimultaneousBackupsPerDatastore       int      `json:"simultaneous_backups_per_datastore,omitempty"`
	SimultaneousBackupsPerHypervisor      int      `json:"simultaneous_backups_per_hypervisor,omitempty"`
	SimultaneousMigrationsPerHypervisor   int      `json:"simultaneous_migrations_per_hypervisor,omitempty"`
	SimultaneousPersonalDeliviries        int      `json:"simultaneous_personal_deliviries,omitempty"`
	SimultaneousStorageResyncTransactions int      `json:"simultaneous_storage_resync_transactions,omitempty"`
	SimultaneousTransactions              int      `json:"simultaneous_transactions,omitempty"`
	SMTPAddress                           string   `json:"smtp_address,omitempty"`
	SMTPAuthentication                    string   `json:"smtp_authentication,omitempty"`
	SMTPDomain                            string   `json:"smtp_domain,omitempty"`
	SMTPEnableStarttlsAuto                bool     `json:"smtp_enable_starttls_auto,bool"`
	SMTPPassword                          string   `json:"smtp_password,omitempty"`
	SMTPPort                              int      `json:"smtp_port,omitempty"`
	SMTPUsername                          string   `json:"smtp_username,omitempty"`
	SnmpStatsLevel1Period                 int      `json:"snmp_stats_level1_period,omitempty"`
	SnmpStatsLevel2Period                 int      `json:"snmp_stats_level2_period,omitempty"`
	SnmpStatsLevel3Period                 int      `json:"snmp_stats_level3_period,omitempty"`
	SnmptrapAddresses                     string   `json:"snmptrap_addresses,omitempty"`
	SnmptrapPort                          int      `json:"snmptrap_port,omitempty"`
	SSHFileTransferOptions                string   `json:"ssh_file_transfer_options,omitempty"`
	SSHFileTransferServer                 string   `json:"ssh_file_transfer_server,omitempty"`
	SSHFileTransferUser                   string   `json:"ssh_file_transfer_user,omitempty"`
	SSHPort                               int      `json:"ssh_port,omitempty"`
	SSHTimeout                            int      `json:"ssh_timeout,omitempty"`
	SslPemPath                            string   `json:"ssl_pem_path,omitempty"`
	StorageEnabled                        bool     `json:"storage_enabled,bool"`
	StorageEndpointOverride               string   `json:"storage_endpoint_override,omitempty"`
	StorageUnicast                        bool     `json:"storage_unicast,bool"`
	SupportHelpEmail                      string   `json:"support_help_email,omitempty"`
	SuspendTimeout                        int      `json:"suspend_timeout,omitempty"`
	SystemAlertReminderPeriod             int      `json:"system_alert_reminder_period,omitempty"`
	SystemEmail                           string   `json:"system_email,omitempty"`
	SystemHost                            string   `json:"system_host,omitempty"`
	SystemNotification                    bool     `json:"system_notification,bool"`
	SystemSupportEmail                    string   `json:"system_support_email,omitempty"`
	SystemTheme                           string   `json:"system_theme,omitempty"`
	TcBurst                               int      `json:"tc_burst,omitempty"`
	TcLatency                             int      `json:"tc_latency,omitempty"`
	TcMtu                                 int      `json:"tc_mtu,omitempty"`
	TemplatePath                          string   `json:"template_path,omitempty"`
	TransactionApprovals                  bool     `json:"transaction_approvals,bool"`
	TransactionRunnerDelay                int      `json:"transaction_runner_delay,omitempty"`
	TransactionStandbyPeriod              int      `json:"transaction_standby_period,omitempty"`
	TrustedProxies                        []string `json:"trusted_proxies,omitempty"`
	UndeployTimeout                       int      `json:"undeploy_timeout,omitempty"`
	UniformNodeCapacityThresholdRatio     int      `json:"uniform_node_capacity_threshold_ratio,omitempty"`
	UnsuspendTimeout                      int      `json:"unsuspend_timeout,omitempty"`
	UpdateServerURL                       string   `json:"update_server_url,omitempty"`
	UploadMediaTimeout                    int      `json:"upload_media_timeout,omitempty"`
	UploadVappTemplateTimeout             int      `json:"upload_vapp_template_timeout,omitempty"`
	URLForCustomTools                     string   `json:"url_for_custom_tools,omitempty"`
	UseHTML5VncConsole                    bool     `json:"use_html5_vnc_console,bool"`
	UseSSHFileTransfer                    bool     `json:"use_ssh_file_transfer,bool"`
	UseYubikeyLogin                       bool     `json:"use_yubikey_login,bool"`
	VcenterL1StatsTimeout                 int      `json:"vcenter_l1_stats_timeout,omitempty"`
	VcenterL2StatsTimeout                 int      `json:"vcenter_l2_stats_timeout,omitempty"`
	VcloudPreventIdleSessionTimeout       int      `json:"vcloud_prevent_idle_session_timeout,omitempty"`
	VcloudStatsLevel1Period               int      `json:"vcloud_stats_level1_period,omitempty"`
	VcloudStatsLevel2Period               int      `json:"vcloud_stats_level2_period,omitempty"`
	VcloudStatsUsageInterval              int      `json:"vcloud_stats_usage_interval,omitempty"`
	WipeOutDiskOnDestroy                  bool     `json:"wipe_out_disk_on_destroy,bool"`
	WizardResourceReservationTTL          int      `json:"wizard_resource_reservation_ttl,omitempty"`
	WrongActivatedLogicalVolumeAlerts     bool     `json:"wrong_activated_logical_volume_alerts,bool"`
	WrongActivatedLogicalVolumeMinutes    int      `json:"wrong_activated_logical_volume_minutes,omitempty"`
	YubikeyAPIID                          string   `json:"yubikey_api_id,omitempty"`
	YubikeyAPIKey                         string   `json:"yubikey_api_key,omitempty"`
	ZabbixHost                            string   `json:"zabbix_host,omitempty"`
	ZabbixPassword                        string   `json:"zabbix_password,omitempty"`
	ZabbixURL                             string   `json:"zabbix_url,omitempty"`
	ZabbixUser                            string   `json:"zabbix_user,omitempty"`
	ZombieDiskSpaceUpdaterDelay           int      `json:"zombie_disk_space_updater_delay,omitempty"`
	ZombieTransactionTime                 int      `json:"zombie_transaction_time,omitempty"`
}

type configurationEditRequestRoot struct {
	Configuration *Configuration `json:"configuration"`
}

type configurationRoot struct {
	Configuration *Configuration `json:"settings"`
}

// Get individual Configuration.
func (s *ConfigurationsServiceOp) Get(ctx context.Context) (*Configuration, *Response, error) {
	path := configurationsBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(configurationRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Configuration, resp, err
}

// Edit individual Configuration.
func (s *ConfigurationsServiceOp) Edit(ctx context.Context, editRequest *Configuration) (*Response, error) {
	path := configurationsEditBasePath + apiFormat

	rootRequest := &configurationEditRequestRoot{
		Configuration: editRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, rootRequest)
	if err != nil {
		return nil, err
	}

	log.Println("[DEBUG] Configuration [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
