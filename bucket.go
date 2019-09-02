package onappgo

import (
  "context"
  "net/http"
  "fmt"
  "log"

  "github.com/digitalocean/godo"
)

const bucketBasePath  = "billing/buckets"

// BucketsService is an interface for interfacing with the Bucket
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/buckets
type BucketsService interface {
  List(context.Context, *ListOptions) ([]Bucket, *Response, error)
  Get(context.Context, int) (*Bucket, *Response, error)
  Create(context.Context, *BucketCreateRequest) (*Bucket, *Response, error)
  Delete(context.Context, int, interface{}) (*Response, error)
  Edit(context.Context, int, *BucketEditRequest) (*Response, error)
}

// BucketsServiceOp handles communication with the Bucket related methods of the
// OnApp API.
type BucketsServiceOp struct {
  client *Client
}

var _ BucketsService = &BucketsServiceOp{}

// Bucket - 
type Bucket struct {
  ID            int         `json:"id,omitempty"`
  Label         string      `json:"label,omitempty"`
  CreatedAt     string      `json:"created_at,omitempty"`
  UpdatedAt     string      `json:"updated_at,omitempty"`
  CurrencyCode  string      `json:"currency_code,omitempty"`
  ShowPrice     bool        `json:"show_price,bool"`
  AllowsMak     bool        `json:"allows_mak,bool"`
  AllowsKms     bool        `json:"allows_kms,bool"`
  AllowsOwn     bool        `json:"allows_own,bool"`
  MonthlyPrice  float64     `json:"monthly_price,omitempty"`
}

// BucketCreateRequest - 
type BucketCreateRequest struct {
  Label         string  `json:"label,omitempty"`
  CurrencyCode  string  `json:"currency_code,omitempty"`
  MonthlyPrice  float64 `json:"monthly_price,omitempty"`
  AllowsKms     bool    `json:"allows_kms,bool"`
  AllowsMak     bool    `json:"allows_mak,bool"`
  AllowsOwn     bool    `json:"allows_own,bool"`
}

type bucketCreateRequestRoot struct {
  BucketCreateRequest  *BucketCreateRequest  `json:"bucket"`
}

type BucketEditRequest BucketCreateRequest

type bucketRoot struct {
  Bucket  *Bucket  `json:"bucket"`
}

func (d BucketCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all Buckets.
func (s *BucketsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Bucket, *Response, error) {
  path := bucketBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Bucket
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]Bucket, len(out))
  for i := range arr {
    arr[i] = out[i]["bucket"]
  }

  return arr, resp, err
}

// Get individual Bucket.
func (s *BucketsServiceOp) Get(ctx context.Context, id int) (*Bucket, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", bucketBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(bucketRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Bucket, resp, err
}

// Create Bucket.
func (s *BucketsServiceOp) Create(ctx context.Context, createRequest *BucketCreateRequest) (*Bucket, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
  }

  path := bucketBasePath + apiFormat

  rootRequest := &bucketCreateRequestRoot{
    BucketCreateRequest : createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("Bucket [Create]  req: ", req)

  root := new(bucketRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Bucket, resp, err
}

// Delete Bucket.
func (s *BucketsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
  if id < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", bucketBasePath, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  log.Println("Bucket [Delete]  req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit Bucket.
func (s *BucketsServiceOp) Edit(ctx context.Context, id int, editRequest *BucketEditRequest) (*Response, error) {
  path := fmt.Sprintf("%s/%d%s", bucketBasePath, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("Bucket [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}
