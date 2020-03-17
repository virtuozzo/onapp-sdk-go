package onappgo

import (
  "context"
  "fmt"
  "log"
  "net/http"

  "github.com/digitalocean/godo"
)

const firewallRulesBasePath string = "virtual_machines/%d/firewall_rules"

// FirewallRulesService is an interface for interfacing with the FirewallRule
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/firewall-rules-for-vss
type FirewallRulesService interface {
  List(context.Context, int, *ListOptions) ([]FirewallRule, *Response, error)
  Get(context.Context, int, int) (*FirewallRule, *Response, error)
  Create(context.Context, int, *FirewallRuleCreateRequest) (*FirewallRule, *Response, error)
  Delete(context.Context, int, int, interface{}) (*Response, error)
  Edit(context.Context, int, int, *FirewallRuleCreateRequest) (*Response, error)
}

// FirewallRulesServiceOp handles communication with the FirewallRules related methods of the
// OnApp API.
type FirewallRulesServiceOp struct {
  client *Client
}

var _ FirewallRulesService = &FirewallRulesServiceOp{}


// FirewallRule - 
// https://docs.onapp.com/apim/latest/firewall-rules-for-vss
type FirewallRule struct {
  ID                  int     `json:"id,omitempty"`
  Position            int     `json:"position,omitempty"`
  Address             string  `json:"address,omitempty"`
  CreatedAt           string  `json:"created_at,omitempty"`
  UpdatedAt           string  `json:"updated_at,omitempty"`
  Command             string  `json:"command,omitempty"`
  Port                string  `json:"port,omitempty"`
  Protocol            string  `json:"protocol,omitempty"`
  NetworkInterfaceID  int     `json:"network_interface_id,omitempty"`
  Comment             string  `json:"comment,omitempty"`
}

// FirewallRuleCreateRequest represents a request to create a FirewallRule
type FirewallRuleCreateRequest struct {
  Address             string `json:"address,omitempty"`
  Command             string `json:"command,omitempty"`
  Protocol            string `json:"protocol,omitempty"`
  NetworkInterfaceID  int    `json:"network_interface_id,omitempty"`
  Comment             string `json:"comment,omitempty"`
  Port                string `json:"port,omitempty"`
}

type firewallRuleCreateRequestRoot struct {
  FirewallRuleCreateRequest *FirewallRuleCreateRequest `json:"firewall_rule"`
}

type firewallRuleRoot struct {
  FirewallRule *FirewallRule `json:"firewall_rule"`
}

func (d FirewallRuleCreateRequest) String() string {
  return godo.Stringify(d)
}

// List all FirewallRules
func (s *FirewallRulesServiceOp) List(ctx context.Context, vmID int, opt *ListOptions) ([]FirewallRule, *Response, error) {
  path := fmt.Sprintf(firewallRulesBasePath, vmID) + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]FirewallRule
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  arr := make([]FirewallRule, len(out))
  for i := range arr {
    arr[i] = out[i]["firewall_rule"]
  }

  return arr, resp, err
}

// Get individual FirewallRule
func (s *FirewallRulesServiceOp) Get(ctx context.Context, vmID int, id int) (*FirewallRule, *Response, error) {
  if vmID < 1 || id < 1 {
    return nil, nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  path := fmt.Sprintf(firewallRulesBasePath, vmID)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(firewallRuleRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.FirewallRule, resp, err
}

// Create FirewallRule
func (s *FirewallRulesServiceOp) Create(ctx context.Context, vmID int, createRequest *FirewallRuleCreateRequest) (*FirewallRule, *Response, error) {
  if vmID < 1 {
    return nil, nil, godo.NewArgError("vmID", "cannot be less than 1")
  }

  if createRequest == nil {
    return nil, nil, godo.NewArgError("FirewallRule createRequest", "cannot be nil")
  }

  path := fmt.Sprintf(firewallRulesBasePath, vmID) + apiFormat
  rootRequest := &firewallRuleCreateRequestRoot{
    FirewallRuleCreateRequest: createRequest,
  }

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
  if err != nil {
    return nil, nil, err
  }
  fmt.Println("FirewallRule [Create] req: ", req)

  root := new(firewallRuleRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.FirewallRule, resp, err
}

// Delete FirewallRule
func (s *FirewallRulesServiceOp) Delete(ctx context.Context, vmID int, id int, meta interface{}) (*Response, error) {
  if vmID < 1 || id < 1 {
    return nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  path := fmt.Sprintf(firewallRulesBasePath, vmID)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
  path, err := addOptions(path, meta)
  if err != nil {
    return nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
  if err != nil {
    return nil, err
  }
  fmt.Println("FirewallRule [Delete] req: ", req)

  return s.client.Do(ctx, req, nil)
}

// Edit FirewallRule
func (s *FirewallRulesServiceOp) Edit(ctx context.Context, vmID int, id int, editRequest *FirewallRuleCreateRequest) (*Response, error) {
  if vmID < 1 || id < 1 {
    return nil, godo.NewArgError("vmID || id", "cannot be less than 1")
  }

  if editRequest == nil {
    return nil, godo.NewArgError("FirewallRule [Edit] editRequest", "cannot be nil")
  }

  path := fmt.Sprintf(firewallRulesBasePath, vmID)
  path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)

  req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
  if err != nil {
    return nil, err
  }
  log.Println("FirewallRule [Edit]  req: ", req)

  return s.client.Do(ctx, req, nil)
}
