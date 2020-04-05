package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const backupResourcesBasePath string = "settings/backups/resources"

// BackupResourcesService is an interface for interfacing with the Backup Resources
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/backup-resources
type BackupResourcesService interface {
	List(context.Context, *ListOptions) ([]BackupResource, *Response, error)
	Get(context.Context, int) (*BackupResource, *Response, error)
	Create(context.Context, *BackupResourceCreateRequest) (*BackupResource, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	// Edit(context.Context, int, *ListOptions) ([]BackupResource, *Response, error)
}

// BackupResourcesServiceOp handles communication with the Backup Resource related methods of the
// OnApp API.
type BackupResourcesServiceOp struct {
	client *Client
}

var _ BackupResourcesService = &BackupResourcesServiceOp{}

// BackupResource represents a BackupResource
type BackupResource struct {
	AdvancedOptions []AdvancedOptions `json:"advanced_options"`
	ID              int               `json:"id,omitempty"`
	CreatedAt       string            `json:"created_at,omitempty"`
	Label           string            `json:"label,omitempty"`
	Enabled         bool              `json:"enabled,bool"`
	Plugin          string            `json:"plugin,omitempty"`
	PrimaryHost     string            `json:"primary_host,omitempty"`
	SecondaryHost   string            `json:"secondary_host,omitempty"`
	UpdatedAt       string            `json:"updated_at,omitempty,omitempty"`
	Username        string            `json:"username,omitempty"`
	Password        string            `json:"password,omitempty"`
	ResourceZoneID  int               `json:"resource_zone_id,omitempty"`

	// OnApp 6.1
	DayToRunOn int    `json:"day_to_run_on,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
}

// BackupResourceCreateRequest represents a request to create a BackupResource
type BackupResourceCreateRequest struct {
	Label          string `json:"label,omitempty"`
	Plugin         string `json:"plugin,omitempty"`
	PrimaryHost    string `json:"primary_host,omitempty"`
	Username       string `json:"username,omitempty"`
	Password       string `json:"password,omitempty"`
	ResourceZoneID int    `json:"resource_zone_id,omitempty"`

	// OnApp 6.1
	// 0 - Sunday, 1 - Monday, 2 - Tuesday, 3 - Wednesday, 4 - Thursday, 5 - Friday, 6 - Saturday
	DayToRunOn int    `json:"day_to_run_on,omitempty"`
	StartTime  string `json:"start_time,omitempty"`
}

type backupResourceCreateRequestRoot struct {
	BackupResourceCreateRequest *BackupResourceCreateRequest `json:"backup_resource"`
}

type backupResourceRoot struct {
	BackupResource *BackupResource `json:"backup_resource"`
}

func (d BackupResourceCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all DataStoreGroups.
func (s *BackupResourcesServiceOp) List(ctx context.Context, opt *ListOptions) ([]BackupResource, *Response, error) {
	path := backupResourcesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]BackupResource
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]BackupResource, len(out))
	for i := range arr {
		arr[i] = out[i]["backup_resource"]
	}

	return arr, resp, err
}

// Get individual BackupResource.
func (s *BackupResourcesServiceOp) Get(ctx context.Context, id int) (*BackupResource, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", backupResourcesBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(backupResourceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.BackupResource, resp, err
}

// Create BackupResource.
func (s *BackupResourcesServiceOp) Create(ctx context.Context, createRequest *BackupResourceCreateRequest) (*BackupResource, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("BackupResource createRequest", "cannot be nil")
	}

	path := backupResourcesBasePath + apiFormat
	rootRequest := &backupResourceCreateRequestRoot{
		BackupResourceCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}

	log.Println("BackupResource [Create] req: ", req)

	root := new(backupResourceRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.BackupResource, resp, err
}

// Delete BackupResource.
func (s *BackupResourcesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", backupResourcesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("BackupResource [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}
