package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const bucketAccessControlsBasePath  = "billing/buckets/%d/access_controls"

// AccessControlsService is an interface for interfacing with the AccessControls
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/buckets/access-control
type AccessControlsService interface {
  List(context.Context, int, *ListOptions) ([]AccessControl, *Response, error)
  // Get(context.Context, int, int) (*AccessControl, *Response, error)
  Create(context.Context, *AccessControlCreateRequest) (*AccessControl, *Response, error)
  Delete(context.Context, *AccessControlDeleteRequest, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]AccessControl, *Response, error)
}

// AccessControlsServiceOp handles communication with the AccessControl related methods of the
// OnApp API.
type AccessControlsServiceOp struct {
  client *Client
}

var _ AccessControlsService = &AccessControlsServiceOp{}

type AccessControl struct {
  BucketID        int           `json:"bucket_id,omitempty"`
  ServerType      string        `json:"server_type,omitempty"`
  TargetID        int           `json:"target_id,omitempty"`
  Type            string        `json:"type,omitempty"`
  TimingStrategy  string        `json:"timing_strategy,omitempty"`
  TargetName      string        `json:"target_name,omitempty"`
  Preferences     interface{}   `json:"preferences,omitempty"`
  Limits          *Limits       `json:"limits,omitempty"`
}

type AccessControlCreateRequest struct {
  BucketID        int          `json:"bucket_id,omitempty"`
  ServerType      string       `json:"server_type,omitempty"`
  TargetID        int          `json:"target_id,omitempty"`
  Type            string       `json:"type,omitempty"`
  Limits          *Limits      `json:"limits,omitempty"`
}

// type accessControlCreateRequestRoot struct {
//   AccessControlCreateRequest  *AccessControlCreateRequest  `json:"access_control"`
// }

type accessControlRoot struct {
  AccessControl  *AccessControl  `json:"access_control"`
}

type AccessControlDeleteRequest AccessControlCreateRequest

// type accessControlDeleteRequestRoot struct {
//   AccessControlDeleteRequest  *AccessControlDeleteRequest  `json:"access_control"`
// }

func (d AccessControlCreateRequest) String() string {
  return godo.Stringify(d)
}

// List return AccessControls for Bucket.
func (s *AccessControlsServiceOp) List(ctx context.Context, id int, opt *ListOptions) ([]AccessControl, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(bucketAccessControlsBasePath, id) + apiFormat

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]AccessControl
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]AccessControl, len(out))
  for i := range arr {
    arr[i] = out[i]["access_control"]
  }

  return arr, resp, err
}

// Create AccessControl.
func (s *AccessControlsServiceOp) Create(ctx context.Context, createRequest *AccessControlCreateRequest) (*AccessControl, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(bucketAccessControlsBasePath, createRequest.BucketID) + apiFormat
  // rootRequest := &accessControlCreateRequestRoot {
  //   AccessControlCreateRequest: createRequest,
  // }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("AccessControl [Create]  req: ", req)

  root := new(accessControlRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, nil, err
  }

  return root.AccessControl, resp, err
}

// Delete AccessControl.
func (s *AccessControlsServiceOp) Delete(ctx context.Context, deleteRequest *AccessControlDeleteRequest, meta interface{}) (*Response, error) {
  if deleteRequest.BucketID < 1 {
    return nil, godo.NewArgError("bucketID", "cannot be less than 1")
  }

  path := fmt.Sprintf(bucketAccessControlsBasePath, deleteRequest.BucketID) + apiFormat
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  // rootRequest := &accessControlDeleteRequestRoot {
  //   AccessControlDeleteRequest : deleteRequest,
  // }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, deleteRequest)
  if err != nil {
    return nil, err
  }
  log.Println("AccessControl [Delete] req: ", req)

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, err
  }

  return resp, err
}

type Limits map[string]interface{}

func LimitsRef(serverType string, resourceType string) *Limits {
  if val, ok := (*(*AccessControls)[serverType])[resourceType]; ok {
    return val
  }

  return nil
}

const (
  MAX_PER_TARGET = "max_per_target"
  MIN_PER_ORIGIN = "min_per_origin"
  MAX_PER_ORIGIN = "max_per_origin"

  DEFAULT        = "default"
  PRESENCE       = "presence"
)

var AccessControls *AccessControlLimits

func init() {
  AccessControls = initializeAccessControlLimits()
}

// Allowed set of limits for Resource based on the ServerType value
func initializeAccessControlLimits() *AccessControlLimits {
  return &AccessControlLimits{
    VIRTUAL :   &LimitResourceRoots{
      COMPUTE_ZONE_RESOURCE              : &Limits{
        "limit_cpu"                         : 0.0,
        "limit_cpu_share"                   : 0.0,
        "limit_cpu_units"                   : 0.0,
        "limit_memory"                      : 0.0,
        "limit_default_cpu"                 : 0.0,
        "limit_min_cpu"                     : 0.0,
        "limit_min_memory"                  : 0.0,
        "limit_default_cpu_share"           : 0.0,
        "limit_min_cpu_priority"            : 0.0,
        "use_cpu_units"                     : false,
        "use_default_cpu"                   : false,
        "use_default_cpu_share"             : false,
      },
      DATA_STORE_ZONE_RESOURCE           : &Limits{
        "limit"                             : 0.0,
      },
      NETWORK_ZONE_RESOURCE              : &Limits{
        "limit_ip"                          : 0.0,
        "limit_rate"                        : 0.0,
      },
      BACKUP_SERVER_ZONE_RESOURCE        : &Limits{
        "limit_backup"                      : 0.0,
        "limit_backup_disk_size"            : 0.0,
        "limit_template"                    : 0.0,
        "limit_template_disk_size"          : 0.0,
        "limit_ova"                         : 0.0,
        "limit_ova_disk_size"               : 0.0,
      },
      VIRTUAL_SERVERS_RESOURCE           : &Limits{
        "limit"                             : 0.0,
      },
      AUTOSCALED_SERVERS_RESOURCE        : &Limits{
        "limit"                             : 0.0,
      },
      COMPUTE_RESOURCE_STORING_RESOURCE  : &Limits{
        "limit"                             : 0.0,
      },
      BACKUPS_RESOURCE                   : &Limits{
        "limit"                             : 0.0,
      },
      TEMPLATES_RESOURCE                 : &Limits{
        "limit"                             : 0.0,
      },
      ISO_TEMPLATES_RESOURCE             : &Limits{
        "limit"                             : 0.0,
      },
      APPLICATION_SERVERS_RESOURCE       : &Limits{
        "limit"                             : 0.0,
      },
      CONTAINER_SERVERS_RESOURCE         : &Limits{
        "limit"                             : 0.0,
      },
      SOLIDFIRE_DATA_STORE_ZONE_RESOURCE : &Limits{
        "limit"                             : 0.0,
      },
      PRECONFIGURED_SERVERS_RESOURCE     : &Limits{},
    },

    SMART:      &LimitResourceRoots{
      COMPUTE_ZONE_RESOURCE              : &Limits{
        "limit_cpu"                         : 0.0,
        "limit_cpu_share"                   : 0.0,
        "limit_cpu_units"                   : 0.0,
        "limit_memory"                      : 0.0,
        "use_cpu_units"                     : false,
      },
      DATA_STORE_ZONE_RESOURCE           : &Limits{
        "limit"                             : 0.0,
      },
      NETWORK_ZONE_RESOURCE              : &Limits{
        "limit_ip"                          : 0.0,
        "limit_rate"                        : 0.0,
      },
      BACKUP_SERVER_ZONE_RESOURCE        : &Limits{
        "limit_backup"                      : 0.0,
        "limit_backup_disk_size"            : 0.0,
        "limit_template"                    : 0.0,
        "limit_template_disk_size"          : 0.0,
      },
      SMART_SERVERS_RESOURCE             : &Limits{
        "limit"                             : 0.0,
      },
      COMPUTE_RESOURCE_STORING_RESOURCE  : &Limits{
        "limit"                             : 0.0,
      },
      BACKUPS_RESOURCE                   : &Limits{
        "limit"                             : 0.0,
      },
    },

    BARE_METAL: &LimitResourceRoots{
      BARE_METAL_SERVERS_RESOURCE        : &Limits{
        "limit"                             : 0.0,
      },
      COMPUTE_ZONE_RESOURCE              : &Limits{},
      NETWORK_ZONE_RESOURCE              : &Limits{
        "limit_ip"                          : 0.0,
        "limit_rate"                        : 0.0,
      },
    },

    VPC:        &LimitResourceRoots{
      VIRTUAL_SERVERS_RESOURCE           : &Limits{
        "limit"                             : 0.0,
      },
      APPLICATION_SERVERS_RESOURCE       : &Limits{
        "limit"                             : 0.0,
      },
      COMPUTE_ZONE_RESOURCE              : &Limits{
        "limit_min_allocation_cpu_allocation"               : 0.0,
        "limit_min_allocation_memory_allocation"            : 0.0,
        "limit_min_allocation_cpu_resources_guaranteed"     : 0.0,
        "limit_min_allocation_memory_resources_guaranteed"  : 0.0,
        "limit_min_allocation_vcpu_speed"                   : 0.0,
        "limit_allocation_cpu_allocation"                   : 0.0,
        "limit_allocation_memory_allocation"                : 0.0,
        "limit_allocation_cpu_resources_guaranteed"         : 0.0,
        "limit_allocation_memory_resources_guaranteed"      : 0.0,
        "limit_allocation_vcpu_speed"                       : 0.0,
        "limit_min_reservation_cpu_allocation"              : 0.0,
        "limit_min_reservation_memory_allocation"           : 0.0,
        "limit_reservation_cpu_allocation"                  : 0.0,
        "limit_reservation_memory_allocation"               : 0.0,
        "limit_min_pay_as_you_go_cpu_limit"                 : 0.0,
        "limit_min_pay_as_you_go_memory_limit"              : 0.0,
        "limit_pay_as_you_go_cpu_limit"                     : 0.0,
        "limit_pay_as_you_go_memory_limit"                  : 0.0,
        "limit_min_pay_as_you_go_vcpu_speed"                : 0.0,
        "limit_vs_cpu"                                      : 0.0,
        "limit_vs_memory"                                   : 0.0,
      },
      DATA_STORE_ZONE_RESOURCE           : &Limits{
        "limit_min_disk_size"                               : 0.0,
        "limit_disk_size"                                   : 0.0,
        "limit_vs_disk_size"                                : 0.0,
      },
      NETWORK_ZONE_RESOURCE              : &Limits{
        "limit_ip"                                          : 0.0,
        "limit_vs_ip"                                       : 0.0,
      },
    },
  }
}

