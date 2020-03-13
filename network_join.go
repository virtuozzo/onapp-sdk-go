package onappgo

import (
  "context"
  "fmt"
  "log"
  "net/http"

  "github.com/digitalocean/godo"
)

var networkJoinPaths = map[string]string {
  "Hypervisor"      : "settings/hypervisors/%d/network_joins",
  "HypervisorGroup" : "settings/hypervisor_zones/%d/network_joins",
}

// NetworkJoinsService is an interface for interfacing with the NetworkJoin
type NetworkJoinsService interface {
  List(context.Context, *NetworkJoinCreateRequest, *ListOptions) ([]NetworkJoin, *Response, error)
  Get(context.Context, string, int, int) (*NetworkJoin, *Response, error)
  Create(context.Context, *NetworkJoinCreateRequest) (*NetworkJoin, *Response, error)
  Delete(context.Context, *NetworkJoinDeleteRequest, interface{}) (*Response, error)
}

// NetworkJoinsServiceOp - 
type NetworkJoinsServiceOp struct {
  client *Client
}

var _ NetworkJoinsService = &NetworkJoinsServiceOp{}

// NetworkJoin represents a NetworkJoin
type NetworkJoin struct {
  ID             int    `json:"id,omitempty"`
  NetworkID      int    `json:"network_id,omitempty"`
  Interface      string `json:"interface,omitempty"`
  CreatedAt      string `json:"created_at,omitempty"`
  UpdatedAt      string `json:"updated_at,omitempty"`
  TargetJoinID   int    `json:"target_join_id,omitempty"`
  TargetJoinType string `json:"target_join_type,omitempty"`
  Identifier     string `json:"identifier,omitempty"`
}

// NetworkJoinCreateRequest represents a request to create a NetworkJoin
type NetworkJoinCreateRequest struct {
  NetworkID      int    `json:"network_id,omitempty"`
  Interface      string `json:"interface,omitempty"`

  // helpers
  TargetJoinID   int    `json:"-"`
  TargetJoinType string `json:"-"`
}

// NetworkJoinDeleteRequest represents a request to delete a NetworkJoin
type NetworkJoinDeleteRequest struct {
  ID             int

  TargetJoinID   int
  TargetJoinType string
}

type networkJoinCreateRequestRoot struct {
  NetworkJoinCreateRequest *NetworkJoinCreateRequest `json:"network_join"`
}

type networkJoinRoot struct {
  NetworkJoin *NetworkJoin `json:"networking_network_join"`
}

func (d NetworkJoinCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all NetworkJoins.
func (s *NetworkJoinsServiceOp) List(ctx context.Context, createRequest *NetworkJoinCreateRequest, opt *ListOptions) ([]NetworkJoin, *Response, error) {
  path := ""
  if val, ok := networkJoinPaths[createRequest.TargetJoinType]; ok {
    path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
  } else {
    return nil, nil, godo.NewArgError("NetworkJoin List: map key not found", createRequest.TargetJoinType)
  }
  
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]NetworkJoin
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]NetworkJoin, len(out))
  for i := range arr {
    arr[i] = out[i]["networking_network_join"]
  }

  return arr, resp, err
}

// Get individual NetworkJoin.
func (s *NetworkJoinsServiceOp) Get(ctx context.Context, targetJoinType string,  targetJoinID int, id int ) (*NetworkJoin, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }
  
  path := ""
  if val, ok := networkJoinPaths[targetJoinType]; ok {
    path = fmt.Sprintf(val, targetJoinID)
  } else {
    return nil, nil, godo.NewArgError("NetworkJoin Get: map key not found", targetJoinType)
  }

  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(networkJoinRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkJoin, resp, err
}

// Create NetworkJoin.
func (s *NetworkJoinsServiceOp) Create(ctx context.Context, createRequest *NetworkJoinCreateRequest) (*NetworkJoin, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("NetworkJoin createRequest", "cannot be nil")
  }

  path := ""
  if val, ok := networkJoinPaths[createRequest.TargetJoinType]; ok {
    path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
  } else {
    return nil, nil, godo.NewArgError("NetworkJoin Create: map key not found", createRequest.TargetJoinType)
  }

  rootRequest := &networkJoinCreateRequestRoot{
    NetworkJoinCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  log.Println("NetworkJoin [Create] req: ", req)

  root := new(networkJoinRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.NetworkJoin, resp, err
}

// Delete NetworkJoin.
func (s *NetworkJoinsServiceOp) Delete(ctx context.Context, deleteRequest *NetworkJoinDeleteRequest, meta interface{}) (*Response, error) {
  if deleteRequest == nil {
    return nil, godo.NewArgError("NetworkJoin deleteRequest", "cannot be nil")
  }

  if deleteRequest.ID < 1 {
    return nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := ""
  if val, ok := networkJoinPaths[deleteRequest.TargetJoinType]; ok {
    path = fmt.Sprintf(val, deleteRequest.TargetJoinID)
  } else {
    return nil, godo.NewArgError("NetworkJoin Delete: map key not found", deleteRequest.TargetJoinType)
  }

  path = fmt.Sprintf("%s/%d%s", path, deleteRequest.ID, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  fmt.Println("NetworkJoin [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}
