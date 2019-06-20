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
}

// VirtualMachineActionsServiceOp handles communication with the VirtualMachine action related
// methods of the OnApp API.
type VirtualMachineActionsServiceOp struct {
  client *Client
}

var _ VirtualMachineActionsService = &VirtualMachineActionsServiceOp{}

// Shutdown a VirtualMachine gracefully
func (s *VirtualMachineActionsServiceOp) Shutdown(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "shutdown"}
  return s.doAction(ctx, id, request)
}

// Stop a VirtualMachine forcefully
func (s *VirtualMachineActionsServiceOp) Stop(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "stop"}
  return s.doAction(ctx, id, request)
}

// Startup a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Startup(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "startup"}
  return s.doAction(ctx, id, request)
}

// Unlock a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Unlock(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "unlock"}
  return s.doAction(ctx, id, request)
}

// Reboot a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Reboot(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "reboot"}
  return s.doAction(ctx, id, request)
}

// Suspend a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Suspend(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "suspend"}
  return s.doAction(ctx, id, request)
}

// Unsuspend a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Unsuspend(ctx context.Context, id int) (*Transaction, *Response, error) {
  request := &ActionRequest{"type": "unsuspend"}
  return s.doAction(ctx, id, request)
}

func (s *VirtualMachineActionsServiceOp) doAction(ctx context.Context, id int, request *ActionRequest) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  if request == nil {
    return nil, nil, godo.NewArgError("request", "request can't be nil")
  }

  path := virtualMachineActionPath(id, request)

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
  if err != nil {
    return nil, nil, err
  }

  // out := new(virtualMachineRoot)
  // resp, err := s.client.Do(ctx, req, out)
  resp, err := s.client.Do(ctx, req, nil)
  if err != nil {
    return nil, resp, err
  }

  transactionsVM, resp, err := s.client.VirtualMachines.Transactions(ctx, id, nil)
  if err != nil {
    fmt.Printf("doAction.Transactions: %s\n\n", err)
    return nil, resp, err
  }

  // return out.VirtualMachine, resp, err
  return &transactionsVM[0], resp, err
}

func virtualMachineActionPath(virtualMachineID int, request *ActionRequest) string {
  return fmt.Sprintf("virtual_machines/%d/%s%s", virtualMachineID, (*request)["type"].(string), apiFormat)
}
