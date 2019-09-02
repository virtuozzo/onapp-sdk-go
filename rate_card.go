package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const bucketRateCardsBasePath = "billing/buckets/%d/rate_cards"

// RateCardsService is an interface for interfacing with the RateCard
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/buckets/rate-card
type RateCardsService interface {
  List(context.Context, int, *ListOptions) ([]RateCard, *Response, error)
  // Get(context.Context, int, int) (*RateCard, *Response, error)
  Create(context.Context, *RateCardCreateRequest) (*RateCard, *Response, error)
  Delete(context.Context, *RateCardDeleteRequest, interface{}) (*Response, error)
  // Edit(context.Context, int, *ListOptions) ([]RateCard, *Response, error)
}

// RateCardsServiceOp handles communication with the RateCard related methods of the
// OnApp API.
type RateCardsServiceOp struct {
  client *Client
}

var _ RateCardsService = &RateCardsServiceOp{}

type RateCard struct {
  BucketID       int      `json:"bucket_id,omitempty"`
  ServerType     string   `json:"server_type,omitempty"`
  TargetID       int      `json:"target_id,omitempty"`
  Type           string   `json:"type,omitempty"`
  TimingStrategy string   `json:"timing_strategy,omitempty"`
  TargetName     string   `json:"target_name,omitempty"`
  Prices         *Prices  `json:"prices,omitempty"`
}

type RateCardCreateRequest struct {
  BucketID                        int     `json:"bucket_id,omitempty"`
  ServerType                      string  `json:"server_type,omitempty"`
  TargetID                        int     `json:"target_id,omitempty"`
  Type                            string  `json:"type,omitempty"`
  TimingStrategy                  string  `json:"timing_strategy,omitempty"`
  ApplyToAllResourcesInTheBucket  bool    `json:"apply_to_all_resources_in_the_bucket,bool"`
  Prices                          *Prices `json:"prices,omitempty"`
}

// type rateCardCreateRequestRoot struct {
//   RateCardCreateRequest  *RateCardCreateRequest  `json:"rate_card"`
// }

type rateCardRoot struct {
  RateCard  *RateCard  `json:"rate_card"`
}

type RateCardDeleteRequest RateCardCreateRequest

// type rateCardDeleteRequestRoot struct {
//   RateCardDeleteRequest  *RateCardDeleteRequest  `json:"rate_card"`
// }

func (d RateCardCreateRequest) String() string {
  return godo.Stringify(d)
}

// List return RateCards for Bucket.
func (s *RateCardsServiceOp) List(ctx context.Context, id int, opt *ListOptions) ([]RateCard, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf(bucketRateCardsBasePath, id) + apiFormat

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]RateCard
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]RateCard, len(out))
  for i := range arr {
    arr[i] = out[i]["rate_card"]
  }

  return arr, resp, err
}

// Create RateCard.
func (s *RateCardsServiceOp) Create(ctx context.Context, createRequest *RateCardCreateRequest) (*RateCard, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(bucketRateCardsBasePath, createRequest.BucketID) + apiFormat
  // rootRequest := &rateCardCreateRequestRoot {
  //   RateCardCreateRequest: createRequest,
  // }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("RateCard [Create]  req: ", req)

  root := new(rateCardRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.RateCard, resp, err
}

// Delete RateCard.
func (s *RateCardsServiceOp) Delete(ctx context.Context, deleteRequest *RateCardDeleteRequest, meta interface{}) (*Response, error) {
  if deleteRequest.BucketID < 1 {
    return nil, godo.NewArgError("bucket_id", "cannot be less than 1")
  }

  path := fmt.Sprintf(bucketRateCardsBasePath, deleteRequest.BucketID) + apiFormat
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  // rootRequest := &rateCardDeleteRequestRoot {
  //   RateCardDeleteRequest: deleteRequest,
  // }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, deleteRequest)
  if err != nil {
    return nil, err
  }
  log.Println("RateCard [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}

type Prices map[string]interface{}

func PricesRef(serverType string, resourceType string) *Prices {
  if val, ok := (*(*RateCards)[serverType])[resourceType]; ok {
    return val
  }

  return nil
}

const (
  FREE_PER_TARGET         = "free_per_target"
  FREE_PER_TARGET_MONTHLY = "free_per_target_monthly"
  FREE_PER_ORIGIN         = "free_per_origin"

  ON                      = "on"      // powered_on = true
  OFF                     = "off"     // powered_on = false
  ONOFF                   = "on/off"  // any state of powered_on
)

var RateCards *RateCardLimits

func init() {
  RateCards = initializeRateCardsLimits()
}

// Allowed set of limits for Resource based on the ServerType value
func initializeRateCardsLimits() *RateCardLimits {
  return &RateCardLimits {
    VIRTUAL : &PriceResourceRoots{
      COMPUTE_ZONE_RESOURCE : &Prices{
        "limit_free_cpu"                            :  0.0,
        "limit_free_cpu_share"                      :  0.0,
        "limit_free_cpu_units"                      :  0.0,
        "limit_free_memory"                         :  0.0,
        "price_on_cpu"                              :  0.0,
        "price_off_cpu"                             :  0.0,
        "price_on_cpu_share"                        :  0.0,
        "price_off_cpu_share"                       :  0.0,
        "price_on_cpu_units"                        :  0.0,
        "price_off_cpu_units"                       :  0.0,
        "price_on_memory"                           :  0.0,
        "price_off_memory"                          :  0.0,
      },
      DATA_STORE_ZONE_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "limit_data_read_free"                      :  0.0,
        "limit_data_written_free"                   :  0.0,
        "limit_reads_completed_free"                :  0.0,
        "limit_writes_completed_free"               :  0.0,
        "limit_free_monthly"                        :  0.0,
        "limit_data_read_free_monthly"              :  0.0,
        "limit_data_written_free_monthly"           :  0.0,
        "limit_reads_completed_free_monthly"        :  0.0,
        "limit_writes_completed_free_monthly"       :  0.0,
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
        "price_data_read"                           :  0.0,
        "price_data_written"                        :  0.0,
        "price_reads_completed"                     :  0.0,
        "price_writes_completed"                    :  0.0,
      },
      NETWORK_ZONE_RESOURCE : &Prices{
        "limit_rate_free"                           :  0.0,
        "limit_ip_free"                             :  0.0,
        "limit_data_sent_free"                      :  0.0,
        "limit_data_received_free"                  :  0.0,
        "limit_ip_free_monthly"                     :  0.0,
        "limit_data_sent_free_monthly"              :  0.0,
        "limit_data_received_free_monthly"          :  0.0,
        "price_rate_on"                             :  0.0,
        "price_rate_off"                            :  0.0,
        "price_ip_on"                               :  0.0,
        "price_ip_off"                              :  0.0,
        "price_data_sent"                           :  0.0,
        "price_data_received"                       :  0.0,
      },
      BACKUP_SERVER_ZONE_RESOURCE : &Prices{
        "limit_backup_free"                         :  0.0,
        "limit_backup_disk_size_free"               :  0.0,
        "limit_template_free"                       :  0.0,
        "limit_template_disk_size_free"             :  0.0,
        "limit_ova_free"                            :  0.0,
        "limit_ova_disk_size_free"                  :  0.0,
        "price_backup"                              :  0.0,
        "price_backup_disk_size"                    :  0.0,
        "price_template"                            :  0.0,
        "price_template_disk_size"                  :  0.0,
        "price_ova"                                 :  0.0,
        "price_ova_disk_size"                       :  0.0,
      },
      AUTOSCALED_SERVERS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      COMPUTE_RESOURCE_STORING_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      BACKUPS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      TEMPLATES_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      ISO_TEMPLATES_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      ACCELERATED_SERVERS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      DRAAS_RESOURCE : &Prices{
        "price_disk_size"                           :  0.0,
        "price_memory"                              :  0.0,
        "price_cpus"                                :  0.0,
        "price_cpu_shares"                          :  0.0,
        "price_cpu_units"                           :  0.0,
        "price_nodes"                               :  0.0,
      },
      SOLIDFIRE_DATA_STORE_ZONE_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
      },
      PRECONFIGURED_SERVERS_RESOURCE : &Prices{
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
        "price_overused_bandwidth"                  :  0.0,
      },
    },

    SMART       : &PriceResourceRoots{
      COMPUTE_ZONE_RESOURCE : &Prices{
        "limit_free_cpu"                            :  0.0,
        "limit_free_cpu_share"                      :  0.0,
        "limit_free_cpu_units"                      :  0.0,
        "limit_free_memory"                         :  0.0,
        "price_on_cpu"                              :  0.0,
        "price_off_cpu"                             :  0.0,
        "price_on_cpu_share"                        :  0.0,
        "price_off_cpu_share"                       :  0.0,
        "price_on_cpu_units"                        :  0.0,
        "price_off_cpu_units"                       :  0.0,
        "price_on_memory"                           :  0.0,
        "price_off_memory"                          :  0.0,
      },
      DATA_STORE_ZONE_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "limit_data_read_free"                      :  0.0,
        "limit_data_written_free"                   :  0.0,
        "limit_reads_completed_free"                :  0.0,
        "limit_writes_completed_free"               :  0.0,
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
        "price_data_read"                           :  0.0,
        "price_data_written"                        :  0.0,
        "price_reads_completed"                     :  0.0,
        "price_writes_completed"                    :  0.0,
      },
      NETWORK_ZONE_RESOURCE : &Prices{
        "limit_rate_free"                           :  0.0,
        "limit_ip_free"                             :  0.0,
        "limit_data_sent_free"                      :  0.0,
        "limit_data_received_free"                  :  0.0,
        "price_rate_on"                             :  0.0,
        "price_rate_off"                            :  0.0,
        "price_ip_on"                               :  0.0,
        "price_ip_off"                              :  0.0,
        "price_data_sent"                           :  0.0,
        "price_data_received"                       :  0.0,
      },
      BACKUP_SERVER_ZONE_RESOURCE : &Prices{
        "limit_backup_free"                         :  0.0,
        "limit_backup_disk_size_free"               :  0.0,
        "limit_template_free"                       :  0.0,
        "limit_template_disk_size_free"             :  0.0,
        "price_backup"                              :  0.0,
        "price_backup_disk_size"                    :  0.0,
        "price_template"                            :  0.0,
        "price_template_disk_size"                  :  0.0,
      },
      AUTOSCALED_SERVERS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      COMPUTE_RESOURCE_STORING_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      BACKUPS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      TEMPLATES_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      ISO_TEMPLATES_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      ACCELERATED_SERVERS_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price"                                     :  0.0,
      },
      DRAAS_RESOURCE : &Prices{
        "price_disk_size"                           :  0.0,
        "price_memory"                              :  0.0,
        "price_cpus"                                :  0.0,
        "price_cpu_shares"                          :  0.0,
        "price_cpu_units"                           :  0.0,
        "price_nodes"                               :  0.0,
      },
      SOLIDFIRE_DATA_STORE_ZONE_RESOURCE : &Prices{
        "limit_free"                                :  0.0,
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
      },
      PRECONFIGURED_SERVERS_RESOURCE : &Prices{
        "price_on"                                  :  0.0,
        "price_off"                                 :  0.0,
        "price_overused_bandwidth"                  :  0.0,
      },
    },

    BARE_METAL  : &PriceResourceRoots{},
    VPC         : &PriceResourceRoots{},
  }
}

