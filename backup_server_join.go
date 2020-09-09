package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

var backupServerJoinPaths = map[string]string{
	"Hypervisor":      "settings/hypervisors/%d/backup_server_joins",
	"HypervisorGroup": "settings/hypervisor_zones/%d/backup_server_joins",
}

// BackupServerJoinsService is an interface for interfacing with the BackupServerJoin
type BackupServerJoinsService interface {
	List(context.Context, *BackupServerJoinCreateRequest, *ListOptions) ([]BackupServerJoin, *Response, error)
	Get(context.Context, string, int, int) (*BackupServerJoin, *Response, error)
	Create(context.Context, *BackupServerJoinCreateRequest) (*BackupServerJoin, *Response, error)
	Delete(context.Context, *BackupServerJoinDeleteRequest, interface{}) (*Response, error)
}

// BackupServerJoinsServiceOp -
type BackupServerJoinsServiceOp struct {
	client *Client
}

var _ BackupServerJoinsService = &BackupServerJoinsServiceOp{}

// BackupServerJoin represents a BackupServerJoin
type BackupServerJoin struct {
	ID             int    `json:"id,omitempty"`
	BackupServerID int    `json:"backup_server_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	TargetJoinID   int    `json:"target_join_id,omitempty"`
	TargetJoinType string `json:"target_join_type,omitempty"`
}

// BackupServerJoinCreateRequest represents a request to create a BackupServerJoin
type BackupServerJoinCreateRequest struct {
	BackupServerID int    `json:"backup_server_id,omitempty"`
	TargetJoinID   int    `json:"-"`
	TargetJoinType string `json:"-"`
}

// BackupServerJoinDeleteRequest represents a request to delete a BackupServerJoin
type BackupServerJoinDeleteRequest struct {
	ID             int
	TargetJoinID   int
	TargetJoinType string
}

type backupServerJoinCreateRequestRoot struct {
	BackupServerID int `json:"backup_server_id,omitempty"`
}

type backupServerJoinRoot struct {
	BackupServerJoin *BackupServerJoin `json:"backup_server_join"`
}

func (d BackupServerJoinCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all BackupServerJoins.
func (s *BackupServerJoinsServiceOp) List(ctx context.Context, createRequest *BackupServerJoinCreateRequest, opt *ListOptions) ([]BackupServerJoin, *Response, error) {
	path := ""
	if val, ok := backupServerJoinPaths[createRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
	} else {
		return nil, nil, godo.NewArgError("BackupServerJoin List: wrong TargetJoinType", createRequest.TargetJoinType)
	}

	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]BackupServerJoin
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]BackupServerJoin, len(out))
	for i := range arr {
		arr[i] = out[i]["backup_server_join"]
	}

	return arr, resp, err
}

// Get individual BackupServerJoin.
func (s *BackupServerJoinsServiceOp) Get(ctx context.Context, targetJoinType string, targetJoinID int, id int) (*BackupServerJoin, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := ""
	if val, ok := backupServerJoinPaths[targetJoinType]; ok {
		path = fmt.Sprintf(val, targetJoinID)
	} else {
		return nil, nil, godo.NewArgError("BackupServerJoin Get: wrong TargetJoinType", targetJoinType)
	}

	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupServerJoinRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.BackupServerJoin, resp, err
}

// Create BackupServerJoin.
func (s *BackupServerJoinsServiceOp) Create(ctx context.Context, createRequest *BackupServerJoinCreateRequest) (*BackupServerJoin, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("BackupServerJoin createRequest", "cannot be nil")
	}

	path := ""
	if val, ok := backupServerJoinPaths[createRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
	} else {
		return nil, nil, godo.NewArgError("BackupServerJoin Create: wrong TargetJoinType", createRequest.TargetJoinType)
	}

	rootRequest := &backupServerJoinCreateRequestRoot{
		BackupServerID: createRequest.BackupServerID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("BackupServerJoin [Create] req: ", req)

	root := new(backupServerJoinRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.BackupServerJoin, resp, err
}

// Delete BackupServerJoin.
func (s *BackupServerJoinsServiceOp) Delete(ctx context.Context, deleteRequest *BackupServerJoinDeleteRequest, meta interface{}) (*Response, error) {
	if deleteRequest == nil {
		return nil, godo.NewArgError("BackupServerJoin deleteRequest", "cannot be nil")
	}

	if deleteRequest.ID < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := ""
	if val, ok := backupServerJoinPaths[deleteRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, deleteRequest.TargetJoinID)
	} else {
		return nil, godo.NewArgError("BackupServerJoin Delete: wrong TargetJoinType", deleteRequest.TargetJoinType)
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
	fmt.Println("BackupServerJoin [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}
