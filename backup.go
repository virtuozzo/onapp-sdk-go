package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const listOfAllVSBackupsBasePath string = "virtual_machines/%d/backups"
const listOfDiskBackupsBasePath string = "virtual_machines/%d/disks/%d/backups"

const createDiskBackupsBasePath string = "settings/disks/%d/backups"
const convertBackupToTemplateBasePath string = "backups/%d/convert"
const deleteBackupsBasePath string = "backups"
const backupNoteBasePath string = "backups/%d/note"

// BackupsService is an interface for interfacing with the Backup
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/backups-snapshots
type BackupsService interface {
	List(context.Context, int, *ListOptions) ([]Backup, *Response, error)
	Get(context.Context, int) (*Backup, *Response, error)
	Create(context.Context, *BackupCreateRequest) (*Backup, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)

	// TODO !!!
	// Move next functions to the BackupActionsService
	AllComputeResourceBackups(context.Context, int) ([]Backup, *Response, error)
	ListOfDiskBackups(context.Context, int, int) ([]Backup, *Response, error)
	BackupNote(context.Context, int, *BackupNoteRequest) (*Response, error)
	ConvertBackupToTemplate(context.Context, int, *ConvertBackupToTemplateRequest) (*Response, error)
}

// BackupsServiceOp handles communication with the Backup related methods of the
// OnApp API.
type BackupsServiceOp struct {
	client *Client
}

var _ BackupsService = &BackupsServiceOp{}

// Backup represent VirtualMachine backup
type Backup struct {
	AllowedHotMigrate        bool   `json:"allowed_hot_migrate,bool"`
	AllowedSwap              bool   `json:"allowed_swap,bool"`
	AllowResizeWithoutReboot bool   `json:"allow_resize_without_reboot,bool"`
	BackupServerID           int    `json:"backup_server_id,omitempty"`
	BackupSize               int    `json:"backup_size,omitempty"`
	BackupType               string `json:"backup_type,omitempty"`
	Built                    bool   `json:"built,bool"`
	BuiltAt                  string `json:"built_at,omitempty"`
	CreatedAt                string `json:"created_at,omitempty"`
	DataStoreType            string `json:"data_store_type,omitempty"`
	DiskID                   int    `json:"disk_id,omitempty"`
	ID                       int    `json:"id,omitempty"`
	Identifier               string `json:"identifier,omitempty"`
	Initiated                string `json:"initiated,omitempty"`
	Iqn                      string `json:"iqn,omitempty"`
	Locked                   bool   `json:"locked,bool"`
	MarkedForDelete          bool   `json:"marked_for_delete,bool"`
	MinDiskSize              int    `json:"min_disk_size,omitempty"`
	MinMemorySize            int    `json:"min_memory_size,omitempty"`
	Note                     string `json:"note,omitempty"`
	OperatingSystem          string `json:"operating_system,omitempty"`
	OperatingSystemDistro    string `json:"operating_system_distro,omitempty"`
	TargetID                 int    `json:"target_id,omitempty"`
	TargetType               string `json:"target_type,omitempty"`
	TemplateID               int    `json:"template_id,omitempty"`
	UpdatedAt                string `json:"updated_at,omitempty"`
	UserID                   int    `json:"user_id,omitempty"`
	VolumeID                 int    `json:"volume_id,omitempty"`
}

// BackupCreateRequest - data for creating Backup
type BackupCreateRequest struct {
	DiskID             int    `json:"disk_id,omitempty"`
	Note               string `json:"note,omitempty"`
	ForceWindowsBackup int    `json:"force_windows_backup,omitempty"`
	VirtualMachineID   int    `json:"-"` // Additional fie42ld to determine Virtual Machine to create disk backup
}

// BackupNoteRequest - data for add/edit backup note
type BackupNoteRequest struct {
	Note string `json:"note,omitempty"`
}

// ConvertBackupToTemplateRequest - data for converting Backup to the Template
type ConvertBackupToTemplateRequest struct {
	Label         string `json:"label,omitempty"`
	MinDiskSize   int    `json:"min_disk_size,omitempty"`
	MinMemorySize int    `json:"min_memory_size,omitempty"`
}

type backupCreateRequestRoot struct {
	BackupCreateRequest *BackupCreateRequest `json:"backup"`
}

type backupRoot struct {
	Backup *Backup `json:"backup"`
}

func (d BackupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Backups in the cloud
func (s *BackupsServiceOp) List(ctx context.Context, vmID int, opt *ListOptions) ([]Backup, *Response, error) {
	path := fmt.Sprintf(listOfAllVSBackupsBasePath, vmID) + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]Backup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Backup, len(out))
	for i := range arr {
		arr[i] = out[i]["backup"]
	}

	return arr, resp, err
}

// Get individual Backup
func (s *BackupsServiceOp) Get(ctx context.Context, id int) (*Backup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	lst, resp, err := s.List(ctx, id, nil)

	return &lst[0], resp, err
}

// Create Backup
func (s *BackupsServiceOp) Create(ctx context.Context, createRequest *BackupCreateRequest) (*Backup, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(createDiskBackupsBasePath, createRequest.DiskID) + apiFormat

	rootRequest := &backupCreateRequestRoot{
		BackupCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Backup [Create]  req: ", req)

	root := new(backupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Backup, resp, err
}

// Delete Backup
func (s *BackupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", deleteBackupsBasePath, id, apiFormat)

	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Backup [Delete]  req: ", req)

	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// AllComputeResourceBackups - 
func (s *BackupsServiceOp) AllComputeResourceBackups(ctx context.Context, vmID int) ([]Backup, *Response, error) {
	return s.List(ctx, vmID, nil)
}

// ListOfDiskBackups -
func (s *BackupsServiceOp) ListOfDiskBackups(ctx context.Context, vmID int, diskID int) ([]Backup, *Response, error) {
	path := fmt.Sprintf(listOfDiskBackupsBasePath, vmID, diskID) + apiFormat
	path, err := addOptions(path, nil)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]Backup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Backup, len(out))
	for i := range arr {
		arr[i] = out[i]["backup"]
	}

	return arr, resp, err
}

// BackupNote - Add/Edit Note of Backup
func (s *BackupsServiceOp) BackupNote(ctx context.Context, id int, noteRequest *BackupNoteRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(backupNoteBasePath, id) + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, noteRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Backup [BackupNote]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// ConvertBackupToTemplate - Convert Backup to the Template
func (s *BackupsServiceOp) ConvertBackupToTemplate(ctx context.Context, id int, convertRequest *ConvertBackupToTemplateRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(convertBackupToTemplateBasePath, id) + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, convertRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Backup [ConvertBackupToTemplate]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
