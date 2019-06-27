package onappgo

import (
  "context"
  "fmt"
  "net/http"

  "github.com/digitalocean/godo"
)

// ActionRequest reprents OnApp Action Request
type ActionRequest map[string]interface{}

// VirtualMachineActionsService is an interface for interfacing with the VirtualMachine actions
// endpoints of the OnApp API
type VirtualMachineActionsService interface {
  Shutdown(context.Context, int) (*Transaction, *Response, error)
  Stop(context.Context, int) (*Transaction, *Response, error)
  Startup(context.Context, int) (*Transaction, *Response, error)
  Unlock(context.Context, int) (*Transaction, *Response, error)
  Reboot(context.Context, int) (*Transaction, *Response, error)
  Suspend(context.Context, int) (*Transaction, *Response, error)
  Unsuspend(context.Context, int) (*Transaction, *Response, error)

  ResetPassword(context.Context, int, string, string) (*Transaction, *Response, error)
  FQDN(context.Context, int, string, string) (*Transaction, *Response, error)
}

// VirtualMachineActionsServiceOp handles communication with the VirtualMachine action related
// methods of the OnApp API.
type VirtualMachineActionsServiceOp struct {
  client *Client
}

var _ VirtualMachineActionsService = &VirtualMachineActionsServiceOp{}

// Shutdown a VirtualMachine gracefully
func (s *VirtualMachineActionsServiceOp) Shutdown(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "shutdown", "action": "stop_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Stop a VirtualMachine forcefully
func (s *VirtualMachineActionsServiceOp) Stop(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "stop", "action": "stop_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Startup a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Startup(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "startup", "action": "startup_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Unlock a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Unlock(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "unlock", "action": "startup_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Reboot a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Reboot(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "reboot", "action": "reboot_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Suspend a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Suspend(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "suspend", "action": "stop_virtual_machine"}
  return s.doAction(ctx, id, request, nil)
}

// Unsuspend a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Unsuspend(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{
    "method" : http.MethodPost,
    "type"   : "unsuspend",
    "path"   : "suspend",
    "action" : "stop_virtual_machine",
  }
  return s.doAction(ctx, id, request, nil)
}

type resetPassword  struct {
  InitialRootPassword               string  `json:"initial_root_password,omitempty"`
  InitialRootPasswordEncryptionKey  string  `json:"initial_root_password_encryption_key,omitempty"`
}

type rootResetPassword struct {
  ResetPassword *resetPassword `json:"virtual_machine"`
}

// ResetPassword a VirtualMachine
func (s *VirtualMachineActionsServiceOp) ResetPassword(ctx context.Context, id int, password string, key string) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPost, "type": "reset_password", "action": "reset_root_password"}

  vmPassword := &resetPassword{
    InitialRootPassword : password,
    InitialRootPasswordEncryptionKey : key,
  }

  root := &rootResetPassword{
    ResetPassword : vmPassword,
  }

  return s.doAction(ctx, id, request, root)
}

// FQDN a VirtualMachine
func (s *VirtualMachineActionsServiceOp) FQDN(ctx context.Context, id int, hostname string, domain string) (*Transaction, *Response, error) {
  request := &ActionRequest{"method": http.MethodPatch, "type": "fqdn", "action": "update_fqdn"}

  vmFQDN := &VirtualMachine{
    Domain : domain,
    Hostname : hostname,
  }

  root := &virtualMachineRoot{
    VirtualMachine : vmFQDN,
  }

  return s.doAction(ctx, id, request, root)
}

func (s *VirtualMachineActionsServiceOp) doAction(ctx context.Context, id int, request *ActionRequest, jsonParams interface{}) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  if request == nil {
    return nil, nil, godo.NewArgError("request", "request can't be nil")
  }

  path, err := virtualMachineActionPath(id, request)
  if err != nil {
    return nil, nil, err
  }

  if (*request)["method"] == nil {
    return nil, nil, godo.NewArgError("method", "must be specified")
  }
  httpMethod := (*request)["method"].(string)

  req, err := s.client.NewRequest(ctx, httpMethod, path, jsonParams)
  if err != nil {
    return nil, nil, err
  }
  // fmt.Printf(" doAction.req: [%+v]\n", req)

  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  opt := &ListOptions{
    PerPage : searchTransactions,
  }

  trxVM, resp, err := s.client.Transactions.ListByGroup(ctx, id, "VirtualMachine", opt)

  // Return last transaction from list of transactions
  var root *Transaction
  e := trxVM.Front()
  if e != nil {
    val := e.Value.(Transaction)
    root = &val
    return root, resp, err
  }

  return nil, nil, err
}

func virtualMachineActionPath(id int, request *ActionRequest) (string, error) {
  if (*request)["type"] == nil {
    return "", godo.NewArgError("type", "must be specified")
  }

  path := (*request)["type"].(string)

  if (*request)["path"] != nil {
    path = (*request)["path"].(string)
  }

  return fmt.Sprintf("virtual_machines/%d/%s%s", id, path, apiFormat), nil
}
