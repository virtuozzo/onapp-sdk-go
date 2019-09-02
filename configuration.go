package onappgo

import (
  "context"
  "net/http"
  "fmt"
)

const configurationsBasePath      = "settings/configuration"
const configurationsEditBasePath  = "settings"

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
  ForceSamlLoginOnly                    bool        `json:"force_saml_login_only,bool"`
  SystemEmail                           string      `json:"system_email,omitempty"`
  SystemHost                            string      `json:"system_host,omitempty"`
  SystemNotification                    bool        `json:"system_notification,bool"`
  EnableNotifications                   bool        `json:"enable_notifications,bool"`
  SystemSupportEmail                    string      `json:"system_support_email,omitempty"`
  SystemTheme                           string      `json:"system_theme,omitempty"`
  PaginationMaxItemsLimit               int         `json:"pagination_max_items_limit,omitempty"`
  AppName                               string      `json:"app_name,omitempty"`
  StorageUnicast                        bool        `json:"storage_unicast,bool"`
  EnableDailyStorageReport              bool        `json:"enable_daily_storage_report,bool"`
  EnableHourlyStorageReport             bool        `json:"enable_hourly_storage_report,bool"`
  DefaultCustomTheme                    string      `json:"default_custom_theme,omitempty"`
  SessionTimeout                        int         `json:"session_timeout,omitempty"`
  MaxNetworkInterfacePortSpeed          int         `json:"max_network_interface_port_speed,omitempty"`
  SslPemPath                            string      `json:"ssl_pem_path,omitempty"`
  UseYubikeyLogin                       bool        `json:"use_yubikey_login,bool"`
  YubikeyAPIKey                         string      `json:"yubikey_api_key,omitempty"`
  YubikeyAPIID                          string      `json:"yubikey_api_id,omitempty"`
  Localdomain                           string      `json:"localdomain,omitempty"`
  RabbitmqHost                          string      `json:"rabbitmq_host,omitempty"`
  RabbitmqPort                          int         `json:"rabbitmq_port,omitempty"`
  RabbitmqVhost                         string      `json:"rabbitmq_vhost,omitempty"`
  RabbitmqLogin                         string      `json:"rabbitmq_login,omitempty"`
  RabbitmqPassword                      string      `json:"rabbitmq_password,omitempty"`
  AllowIncrementalBackups               bool        `json:"allow_incremental_backups,bool"`
  UseSSHFileTransfer                    bool        `json:"use_ssh_file_transfer,bool"`
  SSHFileTransferServer                 string      `json:"ssh_file_transfer_server,omitempty"`
  SSHFileTransferUser                   string      `json:"ssh_file_transfer_user,omitempty"`
  SSHFileTransferOptions                string      `json:"ssh_file_transfer_options,omitempty"`
  SSHPort                               int         `json:"ssh_port,omitempty"`
  SSHTimeout                            int         `json:"ssh_timeout,omitempty"`
  TemplatePath                          string      `json:"template_path,omitempty"`
  BackupsPath                           string      `json:"backups_path,omitempty"`
  DataPath                              string      `json:"data_path,omitempty"`
  UpdateServerURL                       string      `json:"update_server_url,omitempty"`
  DeleteTemplateSourceAfterInstall      bool        `json:"delete_template_source_after_install,bool"`
  LicenseKey                            string      `json:"license_key,omitempty"`
  GenerateComment                       string      `json:"generate_comment,omitempty"`
  SimultaneousBackups                   int         `json:"simultaneous_backups,omitempty"`
  SimultaneousBackupsPerDatastore       int         `json:"simultaneous_backups_per_datastore,omitempty"`
  SimultaneousBackupsPerHypervisor      int         `json:"simultaneous_backups_per_hypervisor,omitempty"`
  SimultaneousTransactions              int         `json:"simultaneous_transactions,omitempty"`
  SimultaneousStorageResyncTransactions int         `json:"simultaneous_storage_resync_transactions,omitempty"`
  SimultaneousPersonalDeliviries        int         `json:"simultaneous_personal_deliviries,omitempty"`
  GuestWaitTimeBeforeDestroy            int         `json:"guest_wait_time_before_destroy,omitempty"`
  RemoteAccessSessionStartPort          int         `json:"remote_access_session_start_port,omitempty"`
  RemoteAccessSessionLastPort           int         `json:"remote_access_session_last_port,omitempty"`
  AjaxPowerUpdateTime                   int         `json:"ajax_power_update_time,omitempty"`
  AjaxPaginationUpdateTime              int         `json:"ajax_pagination_update_time,omitempty"`
  AjaxLogUpdateInterval                 int         `json:"ajax_log_update_interval,omitempty"`
  HypervisorLiveTimes                   int         `json:"hypervisor_live_times,omitempty"`
  IsoPathOnCp                           string      `json:"iso_path_on_cp,omitempty"`
  OvaPath                               string      `json:"ova_path,omitempty"`
  IsoPathOnHv                           string      `json:"iso_path_on_hv,omitempty"`
  RemoveBackupsOnDestroyVM              bool        `json:"remove_backups_on_destroy_vm,bool"`
  DisableHypervisorFailover             bool        `json:"disable_hypervisor_failover,bool"`
  IpsAllowedForLogin                    string      `json:"ips_allowed_for_login,omitempty"`
  MonitisPath                           string      `json:"monitis_path,omitempty"`
  MonitisAccount                        string      `json:"monitis_account,omitempty"`
  MonitisApikey                         string      `json:"monitis_apikey,omitempty"`
  Locales                               []string    `json:"locales,omitempty"`
  MaxMemoryRatio                        int         `json:"max_memory_ratio,omitempty"`
  KvmMaxMemoryRate                      int         `json:"kvm_max_memory_rate,omitempty"`
  KvmAvailableFreeMemoryPercentage      int         `json:"kvm_available_free_memory_percentage,omitempty"`
  DefaultImageTemplate                  int         `json:"default_image_template,omitempty"`
  ServiceAccountName                    string      `json:"service_account_name,omitempty"`
  DefaultFirewallPolicy                 string      `json:"default_firewall_policy,omitempty"`
  DropFirewallPolicyAllowedIps          string      `json:"drop_firewall_policy_allowed_ips,omitempty"`
  ShowIPAddressSelectionForNewVM        bool        `json:"show_ip_address_selection_for_new_vm,bool"`
  BackupTakerDelay                      int         `json:"backup_taker_delay,omitempty"`
  ClusterMonitorDelay                   int         `json:"cluster_monitor_delay,omitempty"`
  HypervisorMonitorDelay                int         `json:"hypervisor_monitor_delay,omitempty"`
  CdnSyncDelay                          int         `json:"cdn_sync_delay,omitempty"`
  ScheduleRunnerDelay                   int         `json:"schedule_runner_delay,omitempty"`
  TransactionRunnerDelay                int         `json:"transaction_runner_delay,omitempty"`
  TransactionApprovals                  bool        `json:"transaction_approvals,bool"`
  BillingTransactionRunnerDelay         int         `json:"billing_transaction_runner_delay,omitempty"`
  ZombieTransactionTime                 int         `json:"zombie_transaction_time,omitempty"`
  ZombieDiskSpaceUpdaterDelay           int         `json:"zombie_disk_space_updater_delay,omitempty"`
  RunRecipeOnVsSleepSeconds             int         `json:"run_recipe_on_vs_sleep_seconds,omitempty"`
  DNSEnabled                            bool        `json:"dns_enabled,bool"`
  AllowStartVmsWithOneIP                bool        `json:"allow_start_vms_with_one_ip,bool"`
  AllowInitialRootPasswordEncryption    bool        `json:"allow_initial_root_password_encryption,bool"`
  WipeOutDiskOnDestroy                  bool        `json:"wipe_out_disk_on_destroy,bool"`
  PartitionAlignOffset                  int         `json:"partition_align_offset,omitempty"`
  PasswordEnforceComplexity             bool        `json:"password_enforce_complexity,bool"`
  PasswordMinimumLength                 int         `json:"password_minimum_length,omitempty"`
  PasswordUpperLowercase                bool        `json:"password_upper_lowercase,bool"`
  PasswordLettersNumbers                bool        `json:"password_letters_numbers,bool"`
  PasswordSymbols                       bool        `json:"password_symbols,bool"`
  PasswordForceUnique                   bool        `json:"password_force_unique,bool"`
  PasswordLockoutAttempts               int         `json:"password_lockout_attempts,omitempty"`
  PasswordExpiry                        int         `json:"password_expiry,omitempty"`
  PasswordHistoryLength                 int         `json:"password_history_length,omitempty"`
  CloudBootEnabled                      bool        `json:"cloud_boot_enabled,bool"`
  IoLimitingEnabled                     bool        `json:"io_limiting_enabled,bool"`
  NfsRootIP                             string      `json:"nfs_root_ip,omitempty"`
  CloudBootTarget                       string      `json:"cloud_boot_target,omitempty"`
  DefaultAccelerationPolicy             bool        `json:"default_acceleration_policy,bool"`
  NumberOfNotificationsToShow           int         `json:"number_of_notifications_to_show,omitempty"`
  NotificationSubjectPrefix             string      `json:"notification_subject_prefix,omitempty"`
  MaxIPAddressesToAssignSimultaneously  int         `json:"max_ip_addresses_to_assign_simultaneously,omitempty"`
  WizardResourceReservationTTL          int         `json:"wizard_resource_reservation_ttl,omitempty"`
  StorageEnabled                        bool        `json:"storage_enabled,bool"`
  PreferLocalReads                      bool        `json:"prefer_local_reads,bool"`
  IntraHypervisorBalanceThresholdRatio  int         `json:"intra_hypervisor_balance_threshold_ratio,omitempty"`
  InterHypervisorBalanceThresholdRatio  int         `json:"inter_hypervisor_balance_threshold_ratio,omitempty"`
  UniformNodeCapacityThresholdRatio     int         `json:"uniform_node_capacity_threshold_ratio,omitempty"`
  AllowHypervisorPasswordEncryption     bool        `json:"allow_hypervisor_password_encryption,bool"`
  ArchiveStatsPeriod                    int         `json:"archive_stats_period,omitempty"`
  InstantStatsPeriod                    int         `json:"instant_stats_period,omitempty"`
  IsArchiveStatsEnabled                 bool        `json:"is_archive_stats_enabled,bool"`
  DisableBilling                        bool        `json:"disable_billing,bool"`
  SystemAlertReminderPeriod             int         `json:"system_alert_reminder_period,omitempty"`
  WrongActivatedLogicalVolumeAlerts     bool        `json:"wrong_activated_logical_volume_alerts,bool"`
  WrongActivatedLogicalVolumeMinutes    int         `json:"wrong_activated_logical_volume_minutes,omitempty"`
  UseHTML5VncConsole                    bool        `json:"use_html5_vnc_console,bool"`
  StorageEndpointOverride               string      `json:"storage_endpoint_override,omitempty"`
  URLForCustomTools                     string      `json:"url_for_custom_tools,omitempty"`
  BackupConvertCoefficient              float64     `json:"backup_convert_coefficient,omitempty"`
  RsyncOptionXattrs                     bool        `json:"rsync_option_xattrs,bool"`
  RsyncOptionAcls                       bool        `json:"rsync_option_acls,bool"`
  SimultaneousBackupsPerBackupServer    int         `json:"simultaneous_backups_per_backup_server,omitempty"`
  EmailDeliveryMethod                   string      `json:"email_delivery_method,omitempty"`
  SMTPAddress                           string      `json:"smtp_address,omitempty"`
  SMTPPort                              int         `json:"smtp_port,omitempty"`
  SMTPDomain                            string      `json:"smtp_domain,omitempty"`
  SMTPUsername                          string      `json:"smtp_username,omitempty"`
  SMTPPassword                          string      `json:"smtp_password,omitempty"`
  SMTPAuthentication                    string      `json:"smtp_authentication,omitempty"`
  SMTPEnableStarttlsAuto                bool        `json:"smtp_enable_starttls_auto,bool"`
  SnmptrapAddresses                     string      `json:"snmptrap_addresses,omitempty"`
  SnmptrapPort                          int         `json:"snmptrap_port,omitempty"`
  InfinibandCloudBootEnabled            bool        `json:"infiniband_cloud_boot_enabled,bool"`
  QemuDetachDeviceDelay                 int         `json:"qemu_detach_device_delay,omitempty"`
  QemuAttachDeviceDelay                 int         `json:"qemu_attach_device_delay,omitempty"`
  TcLatency                             int         `json:"tc_latency,omitempty"`
  TcBurst                               int         `json:"tc_burst,omitempty"`
  TcMtu                                 int         `json:"tc_mtu,omitempty"`
  HaEnabled                             bool        `json:"ha_enabled,bool"`
  DashboardAPIAccessToken               string      `json:"dashboard_api_access_token,omitempty"`
  AllowConnectAws                       bool        `json:"allow_connect_aws,bool"`
  FederationTrustsOnlyPrivate           bool        `json:"federation_trusts_only_private,bool"`
  MaximumPendingTasks                   int         `json:"maximum_pending_tasks,omitempty"`
  MaxUploadSize                         int         `json:"max_upload_size,omitempty"`
  TransactionStandbyPeriod              int         `json:"transaction_standby_period,omitempty"`
  LogCleanupPeriod                      int         `json:"log_cleanup_period,omitempty"`
  LogCleanupEnabled                     bool        `json:"log_cleanup_enabled,bool"`
  LogLevel                              string      `json:"log_level,omitempty"`
  CdnMaxResultsPerGetPage               int         `json:"cdn_max_results_per_get_page,omitempty"`
  InstancePackagesThresholdNum          int         `json:"instance_packages_threshold_num,omitempty"`
  AmountOfServiceInstances              int         `json:"amount_of_service_instances,omitempty"`
  GracefulStopTimeout                   int         `json:"graceful_stop_timeout,omitempty"`
  AllowToCollectErrors                  bool        `json:"allow_to_collect_errors,bool"`
  PasswordProtectionForDeleting         bool        `json:"password_protection_for_deleting,bool"`
  DraasEnabled                          bool        `json:"draas_enabled,bool"`
  ZabbixHost                            string      `json:"zabbix_host,omitempty"`
  ZabbixURL                             string      `json:"zabbix_url,omitempty"`
  ZabbixUser                            string      `json:"zabbix_user,omitempty"`
  ZabbixPassword                        string      `json:"zabbix_password,omitempty"`
  PingVmsBeforeInitFailover             bool        `json:"ping_vms_before_init_failover,bool"`
  VcloudStatsLevel1Period               int         `json:"vcloud_stats_level1_period,omitempty"`
  VcloudStatsLevel2Period               int         `json:"vcloud_stats_level2_period,omitempty"`
  VcloudStatsUsageInterval              int         `json:"vcloud_stats_usage_interval,omitempty"`
  VcloudPreventIdleSessionTimeout       int         `json:"vcloud_prevent_idle_session_timeout,omitempty"`
  BlockSize                             int         `json:"block_size,omitempty"`
  CloudBootDomainNameServers            string      `json:"cloud_boot_domain_name_servers,omitempty"`
  EnforceRedundancy                     bool        `json:"enforce_redundancy,bool"`
  DashboardStats                        []string    `json:"dashboard_stats,omitempty"`
  MigrationRateLimit                    int         `json:"migration_rate_limit,omitempty"`
  SimultaneousMigrationsPerHypervisor   int         `json:"simultaneous_migrations_per_hypervisor,omitempty"`
  AllowAdvancedVsManagement             bool        `json:"allow_advanced_vs_management,bool"`
  SupportHelpEmail                      string      `json:"support_help_email,omitempty"`
  IPAddressReservationTime              int         `json:"ip_address_reservation_time,omitempty"`
  RecipeTmpDir                          string      `json:"recipe_tmp_dir,omitempty"`
  SnmpStatsLevel1Period                 int         `json:"snmp_stats_level1_period,omitempty"`
  SnmpStatsLevel2Period                 int         `json:"snmp_stats_level2_period,omitempty"`
  SnmpStatsLevel3Period                 int         `json:"snmp_stats_level3_period,omitempty"`
  ActionGlobalLockExpirationTimeout     int         `json:"action_global_lock_expiration_timeout,omitempty"`
  ActionGlobalLockRetryDelay            int         `json:"action_global_lock_retry_delay,omitempty"`
  IsolatedLicense                       bool        `json:"isolated_license,bool"`
  PaginationDashboardPagesLimit         int         `json:"pagination_dashboard_pages_limit,omitempty"`
  TrustedProxies                        []string    `json:"trusted_proxies,omitempty"`
  DefaultTimeout                        int         `json:"default_timeout,omitempty"`
  DeleteVappTemplateTimeout             int         `json:"delete_vapp_template_timeout,omitempty"`
  DeleteVappTimeout                     int         `json:"delete_vapp_timeout,omitempty"`
  DeleteMediaTimeout                    int         `json:"delete_media_timeout,omitempty"`
  InstantiateVappTemplateTimeout        int         `json:"instantiate_vapp_template_timeout,omitempty"`
  PowerOnTimeout                        int         `json:"power_on_timeout,omitempty"`
  PowerOffTimeout                       int         `json:"power_off_timeout,omitempty"`
  SuspendTimeout                        int         `json:"suspend_timeout,omitempty"`
  DiscardSuspendTimeout                 int         `json:"discard_suspend_timeout,omitempty"`
  RebootTimeout                         int         `json:"reboot_timeout,omitempty"`
  UndeployTimeout                       int         `json:"undeploy_timeout,omitempty"`
  ProcessDescriptorVappTemplateTimeout  int         `json:"process_descriptor_vapp_template_timeout,omitempty"`
  HTTPRequestTimeout                    int         `json:"http_request_timeout,omitempty"`
  RecomposeVappTimeout                  int         `json:"recompose_vapp_timeout,omitempty"`
  CreateEdgeGatewayTimeout              int         `json:"create_edge_gateway_timeout,omitempty"`
  ComposeVappTemplateTimeout            int         `json:"compose_vapp_template_timeout,omitempty"`
  CreateSnapshotTimeout                 int         `json:"create_snapshot_timeout,omitempty"`
  CreateVdcTimeout                      int         `json:"create_vdc_timeout,omitempty"`
}

type configurationEditRequestRoot struct {
  Configuration  *Configuration  `json:"configuration"`
}

type configurationRoot struct {
  Configuration  *Configuration  `json:"settings"`
}

// Get individual Configuration.
func (s *ConfigurationsServiceOp) Get(ctx context.Context) (*Configuration, *Response, error) {
  path := fmt.Sprintf("%s%s", configurationsBasePath, apiFormat)

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
  path := fmt.Sprintf("%s%s", configurationsEditBasePath, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
  if err != nil {
    return nil, err
  }

  rootRequest := &configurationEditRequestRoot{
    Configuration : editRequest,
  }

  return s.client.Do(ctx, req, rootRequest)
}
