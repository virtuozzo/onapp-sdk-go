package onappgo

import (
  "context"
  "fmt"
  "log"
  "net/http"

  "github.com/digitalocean/godo"
)

var pathNetworkJoin = map[string]string {
  "Hypervisor"        : "settings/hypervisors/%d/network_joins",
  "HypervisorZone"    : "settings/hypervisor_zones/%d/network_joins",
}

// NetworkJoinsService is an interface for interfacing with the NetworkJoin
type NetworkJoinsService interface {
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
  Interface      string `json:"interface,omitempty"`
  NetworkID      int    `json:"network_id,omitempty"`

  // helpers
  TargetJoinID   int    `json:"target_join_id,omitempty"`
  TargetJoinType string `json:"target_join_type,omitempty"`
}

// NetworkJoinDeleteRequest represents a request to delete a NetworkJoin
type NetworkJoinDeleteRequest struct {
  ID             int    `json:"id,omitempty"`
  TargetJoinType string `json:"target_join_type,omitempty"`
}

type networkJoinCreateRequestRoot struct {
  NetworkJoinCreateRequest *NetworkJoinCreateRequest `json:"network_join"`
}

type networkJoinRoot struct {
  NetworkJoin *NetworkJoin `json:"network_join"`
}

func (d NetworkJoinCreateRequest) String() string {
  return godo.Stringify(d)
}

// Create NetworkJoin.
func (s *NetworkJoinsServiceOp) Create(ctx context.Context, createRequest *NetworkJoinCreateRequest) (*NetworkJoin, *Response, error) {
  if createRequest == nil {
    return nil, nil, godo.NewArgError("NetworkJoin createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(pathNetworkJoin[createRequest.TargetJoinType], createRequest.TargetJoinID) + apiFormat
  rootRequest := &networkJoinCreateRequestRoot{
    NetworkJoinCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  fmt.Println("NetworkJoin [Create] req: ", req)

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

  path := fmt.Sprintf(pathNetworkJoin[deleteRequest.TargetJoinType], deleteRequest.ID) + apiFormat
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
