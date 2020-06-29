package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const sshKeyBasePath string = "settings/ssh_keys"
const addSSHKeyBasePath string = "users/%d/ssh_keys"

// SSHKeysService is an interface for interfacing with the SSHKey
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/ssh-keys
type SSHKeysService interface {
	List(context.Context, *ListOptions) ([]SSHKey, *Response, error)
	Get(context.Context, int) (*SSHKey, *Response, error)
	Create(context.Context, *SSHKeyCreateRequest) (*SSHKey, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *SSHKeyEditRequest) (*Response, error)
}

// SSHKeysServiceOp handles communication with the SSHKey related methods of the
// OnApp API.
type SSHKeysServiceOp struct {
	client *Client
}

var _ SSHKeysService = &SSHKeysServiceOp{}

// SSHKey - represent disk from Virtual Machine
type SSHKey struct {
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	Key       string `json:"key,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// SSHKeyCreateRequest - data for creating SSHKey
type SSHKeyCreateRequest struct {
	UserID int    `json:"user_id,omitempty"`
	Key    string `json:"key,omitempty"`
}

// SSHKeyEditRequest - data for editing SSHKey
type SSHKeyEditRequest struct {
	UserID int    `json:"user_id,omitempty"`
	Key    string `json:"key,omitempty"`
}

type sshKeyCreateRequestRoot struct {
	SSHKeyCreateRequest *SSHKeyCreateRequest `json:"ssh_key"`
}

type sshKeyRoot struct {
	SSHKey *SSHKey `json:"ssh_key"`
}

func (d SSHKeyCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all SSH Keys in the cloud.
func (s *SSHKeysServiceOp) List(ctx context.Context, opt *ListOptions) ([]SSHKey, *Response, error) {
	path := sshKeyBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]SSHKey
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]SSHKey, len(out))
	for i := range arr {
		arr[i] = out[i]["ssh_key"]
	}

	return arr, resp, err
}

// Get individual SSH key.
func (s *SSHKeysServiceOp) Get(ctx context.Context, id int) (*SSHKey, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", sshKeyBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(sshKeyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// Create SSHKey.
func (s *SSHKeysServiceOp) Create(ctx context.Context, createRequest *SSHKeyCreateRequest) (*SSHKey, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(addSSHKeyBasePath, createRequest.UserID) + apiFormat

	rootRequest := &sshKeyCreateRequestRoot{
		SSHKeyCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("SSHKey [Create]  req: ", req)

	root := new(sshKeyRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.SSHKey, resp, err
}

// Delete SSHKey.
func (s *SSHKeysServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", sshKeyBasePath, id, apiFormat)

	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("SSHKey [Delete]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit SSHKey.
func (s *SSHKeysServiceOp) Edit(ctx context.Context, id int, editRequest *SSHKeyEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", sshKeyBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("SSHKey [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
