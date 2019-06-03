package onappgo

import (
  "context"
  "fmt"
  "net/http"
  "net/url"

  "github.com/digitalocean/godo"
)

// ActionRequest reprents OnApp Action Request
type ActionRequest map[string]interface{}

// VirtualMachineActionsService is an interface for interfacing with the VirtualMachine actions
// endpoints of the OnApp API
type VirtualMachineActionsService interface {
  Shutdown(context.Context, int) (*godo.Action, *Response, error)
  // ShutdownByTag(context.Context, string) ([]godo.Action, *Response, error)
  PowerOff(context.Context, int) (*godo.Action, *Response, error)
  // PowerOffByTag(context.Context, string) ([]godo.Action, *Response, error)
  PowerOn(context.Context, int) (*godo.Action, *Response, error)
  // PowerOnByTag(context.Context, string) ([]godo.Action, *Response, error)
  PowerCycle(context.Context, int) (*godo.Action, *Response, error)
  // PowerCycleByTag(context.Context, string) ([]godo.Action, *Response, error)
  Reboot(context.Context, int) (*godo.Action, *Response, error)
  Restore(context.Context, int, int) (*godo.Action, *Response, error)
  Resize(context.Context, int, string, bool) (*godo.Action, *Response, error)
  Rename(context.Context, int, string) (*godo.Action, *Response, error)
  Snapshot(context.Context, int, string) (*godo.Action, *Response, error)
  // SnapshotByTag(context.Context, string, string) ([]godo.Action, *Response, error)
  EnableBackups(context.Context, int) (*godo.Action, *Response, error)
  // EnableBackupsByTag(context.Context, string) ([]godo.Action, *Response, error)
  DisableBackups(context.Context, int) (*godo.Action, *Response, error)
  // DisableBackupsByTag(context.Context, string) ([]godo.Action, *Response, error)
  PasswordReset(context.Context, int) (*godo.Action, *Response, error)
  RebuildByImageID(context.Context, int, int) (*godo.Action, *Response, error)
  RebuildByImageSlug(context.Context, int, string) (*godo.Action, *Response, error)
  ChangeKernel(context.Context, int, int) (*godo.Action, *Response, error)
  EnableIPv6(context.Context, int) (*godo.Action, *Response, error)
  // EnableIPv6ByTag(context.Context, string) ([]godo.Action, *Response, error)
  EnablePrivateNetworking(context.Context, int) (*godo.Action, *Response, error)
  // EnablePrivateNetworkingByTag(context.Context, string) ([]godo.Action, *Response, error)
  Get(context.Context, int, int) (*godo.Action, *Response, error)
  GetByURI(context.Context, string) (*godo.Action, *Response, error)
}

// VirtualMachineActionsServiceOp handles communication with the VirtualMachine action related
// methods of the OnApp API.
type VirtualMachineActionsServiceOp struct {
  client *Client
}

var _ VirtualMachineActionsService = &VirtualMachineActionsServiceOp{}

// Shutdown a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Shutdown(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "shutdown"}
  return s.doAction(ctx, id, request)
}

// ShutdownByTag shuts down VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) ShutdownByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
// 	request := &ActionRequest{"type": "shutdown"}
// 	return s.doActionByTag(ctx, tag, request)
// }

// PowerOff a VirtualMachine
func (s *VirtualMachineActionsServiceOp) PowerOff(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "power_off"}
  return s.doAction(ctx, id, request)
}

// PowerOffByTag powers off VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) PowerOffByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "power_off"}
//   return s.doActionByTag(ctx, tag, request)
// }

// PowerOn a VirtualMachine
func (s *VirtualMachineActionsServiceOp) PowerOn(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "power_on"}
  return s.doAction(ctx, id, request)
}

// PowerOnByTag powers on VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) PowerOnByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "power_on"}
//   return s.doActionByTag(ctx, tag, request)
// }

// PowerCycle a VirtualMachine
func (s *VirtualMachineActionsServiceOp) PowerCycle(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "power_cycle"}
  return s.doAction(ctx, id, request)
}

// PowerCycleByTag power cycles VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) PowerCycleByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "power_cycle"}
//   return s.doActionByTag(ctx, tag, request)
// }

// Reboot a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Reboot(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "reboot"}
  return s.doAction(ctx, id, request)
}

// Restore an image to a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Restore(ctx context.Context, id, imageID int) (*godo.Action, *Response, error) {
  requestType := "restore"
  request := &ActionRequest{
    "type":  requestType,
    "image": float64(imageID),
  }
  return s.doAction(ctx, id, request)
}

// Resize a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Resize(ctx context.Context, id int, sizeSlug string, resizeDisk bool) (*godo.Action, *Response, error) {
  requestType := "resize"
  request := &ActionRequest{
    "type": requestType,
    "size": sizeSlug,
    "disk": resizeDisk,
  }
  return s.doAction(ctx, id, request)
}

// Rename a VirtualMachine
func (s *VirtualMachineActionsServiceOp) Rename(ctx context.Context, id int, name string) (*godo.Action, *Response, error) {
  requestType := "rename"
  request := &ActionRequest{
    "type": requestType,
    "name": name,
  }
  return s.doAction(ctx, id, request)
}

// Snapshot a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) Snapshot(ctx context.Context, id int, name string) (*godo.Action, *Response, error) {
  requestType := "snapshot"
  request := &ActionRequest{
    "type": requestType,
    "name": name,
  }
  return s.doAction(ctx, id, request)
}

// SnapshotByTag snapshots VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) SnapshotByTag(ctx context.Context, tag string, name string) ([]godo.Action, *Response, error) {
//   requestType := "snapshot"
//   request := &ActionRequest{
//     "type": requestType,
//     "name": name,
//   }
//   return s.doActionByTag(ctx, tag, request)
// }

// EnableBackups enables backups for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) EnableBackups(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "enable_backups"}
  return s.doAction(ctx, id, request)
}

// EnableBackupsByTag enables backups for VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) EnableBackupsByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "enable_backups"}
//   return s.doActionByTag(ctx, tag, request)
// }

// DisableBackups disables backups for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) DisableBackups(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "disable_backups"}
  return s.doAction(ctx, id, request)
}

// DisableBackupsByTag disables backups for VirtualMachine matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) DisableBackupsByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "disable_backups"}
//   return s.doActionByTag(ctx, tag, request)
// }

// PasswordReset resets the password for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) PasswordReset(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "password_reset"}
  return s.doAction(ctx, id, request)
}

// RebuildByImageID rebuilds a VirtualMachine from an image with a given id.
func (s *VirtualMachineActionsServiceOp) RebuildByImageID(ctx context.Context, id, imageID int) (*godo.Action, *Response, error) {
  request := &ActionRequest{
    "type": "rebuild",
    "image": imageID,
  }
  return s.doAction(ctx, id, request)
}

// RebuildByImageSlug rebuilds a VirtualMachine from an Image matched by a given Slug.
func (s *VirtualMachineActionsServiceOp) RebuildByImageSlug(ctx context.Context, id int, slug string) (*godo.Action, *Response, error) {
  request := &ActionRequest{
    "type": "rebuild",
    "image": slug,
  }
  return s.doAction(ctx, id, request)
}

// ChangeKernel changes the kernel for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) ChangeKernel(ctx context.Context, id, kernelID int) (*godo.Action, *Response, error) {
  request := &ActionRequest{
    "type": "change_kernel",
    "kernel": kernelID,
  }
  return s.doAction(ctx, id, request)
}

// EnableIPv6 enables IPv6 for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) EnableIPv6(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "enable_ipv6"}
  return s.doAction(ctx, id, request)
}

// EnableIPv6ByTag enables IPv6 for VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) EnableIPv6ByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "enable_ipv6"}
//   return s.doActionByTag(ctx, tag, request)
// }

// EnablePrivateNetworking enables private networking for a VirtualMachine.
func (s *VirtualMachineActionsServiceOp) EnablePrivateNetworking(ctx context.Context, id int) (*godo.Action, *Response, error) {
  request := &ActionRequest{"type": "enable_private_networking"}
  return s.doAction(ctx, id, request)
}

// EnablePrivateNetworkingByTag enables private networking for VirtualMachines matched by a Tag.
// func (s *VirtualMachineActionsServiceOp) EnablePrivateNetworkingByTag(ctx context.Context, tag string) ([]godo.Action, *Response, error) {
//   request := &ActionRequest{"type": "enable_private_networking"}
//   return s.doActionByTag(ctx, tag, request)
// }

func (s *VirtualMachineActionsServiceOp) doAction(ctx context.Context, id int, request *ActionRequest) (*godo.Action, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  if request == nil {
    return nil, nil, godo.NewArgError("request", "request can't be nil")
  }

  path := virtualMachineActionPath(id)

  req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
  if err != nil {
    return nil, nil, err
  }

  var out map[string]godo.Action
  resp, err := s.client.Do(ctx, req, out)
  if err != nil {
    return nil, resp, err
  }

  vm := out["action"]

  return &vm, resp, err
}

// func (s *VirtualMachineActionsServiceOp) doActionByTag(ctx context.Context, tag string, request *ActionRequest) ([]godo.Action, *Response, error) {
//   if tag == "" {
//     return nil, nil, godo.NewArgError("tag", "cannot be empty")
//   }

//   if request == nil {
//     return nil, nil, godo.NewArgError("request", "request can't be nil")
//   }

//   path := virtualMachineActionPathByTag(tag)

//   req, err := s.client.NewRequest(ctx, http.MethodPost, path, request)
//   if err != nil {
//     return nil, nil, err
//   }

//   root := new(actionsRoot)
//   resp, err := s.client.Do(ctx, req, root)
//   if err != nil {
//     return nil, resp, err
//   }

//   return root.Actions, resp, err
// }

// Get an action for a particular VirtualMachine by id.
func (s *VirtualMachineActionsServiceOp) Get(ctx context.Context, virtualMachineID, actionID int) (*godo.Action, *Response, error) {
  if virtualMachineID < 1 {
    return nil, nil, godo.NewArgError("virtualMachineID", "cannot be less than 1")
  }

  if actionID < 1 {
    return nil, nil, godo.NewArgError("actionID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d", virtualMachineActionPath(virtualMachineID), actionID)
  return s.get(ctx, path)
}

// GetByURI gets an action for a particular VirtualMachine by id.
func (s *VirtualMachineActionsServiceOp) GetByURI(ctx context.Context, rawurl string) (*godo.Action, *Response, error) {
  u, err := url.Parse(rawurl)
  if err != nil {
    return nil, nil, err
  }

  return s.get(ctx, u.Path)

}

func (s *VirtualMachineActionsServiceOp) get(ctx context.Context, path string) (*godo.Action, *Response, error) {
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out map[string]godo.Action
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  vm := out["action"]

  return &vm, resp, err
}

func virtualMachineActionPath(virtualMachineID int) string {
  return fmt.Sprintf("virtual_machines/%d/actions", virtualMachineID)
}

// func virtualMachineActionPathByTag(tag string) string {
//   return fmt.Sprintf("virtual_machines/actions?tag_name=%s", tag)
// }
