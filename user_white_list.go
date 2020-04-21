package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const userWhiteListsBasePath string = "/users/%d/user_white_lists"

// UserWhiteListsService is an interface for interfacing with the UserWhiteList
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/whitelist-ips
type UserWhiteListsService interface {
	List(context.Context, int, *ListOptions) ([]UserWhiteList, *Response, error)
	Get(context.Context, int, int) (*UserWhiteList, *Response, error)
	Create(context.Context, int, *UserWhiteListCreateRequest) (*UserWhiteList, *Response, error)
	Delete(context.Context, int, int, interface{}) (*Response, error)
	Edit(context.Context, int, int, *UserWhiteListEditRequest) (*Response, error)
}

// UserWhiteListsServiceOp handles communication with the UserWhiteLists related methods of the
// OnApp API.
type UserWhiteListsServiceOp struct {
	client *Client
}

var _ UserWhiteListsService = &UserWhiteListsServiceOp{}

// UserWhiteList represents a UserWhiteList
type UserWhiteList struct {
	CreatedAt   string `json:"created_at,omitempty"`
	Description string `json:"description"` // can be empty
	ID          int    `json:"id,omitempty"`
	IP          string `json:"ip,omitempty"`
	UpdatedAt   string `json:"updated_at,omitempty"`
	UserID      int    `json:"user_id,omitempty"`
}

// UserWhiteListCreateRequest represents a request to create a UserWhiteList
type UserWhiteListCreateRequest struct {
	Description string `json:"description"` // can be empty
	IP          string `json:"ip,omitempty"`
}

// UserWhiteListEditRequest represents a request to edit a UserWhiteList
type UserWhiteListEditRequest UserWhiteListCreateRequest

type userWhiteListCreateRequestRoot struct {
	UserWhiteListCreateRequest *UserWhiteListCreateRequest `json:"user_white_list"`
}

type userWhiteListRoot struct {
	UserWhiteList *UserWhiteList `json:"user_white_list"`
}

func (d UserWhiteListCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all UserWhiteLists.
func (s *UserWhiteListsServiceOp) List(ctx context.Context, userID int, opt *ListOptions) ([]UserWhiteList, *Response, error) {
	path := fmt.Sprintf(userWhiteListsBasePath, userID) + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]UserWhiteList
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]UserWhiteList, len(out))
	for i := range arr {
		arr[i] = out[i]["user_white_list"]
	}

	return arr, resp, err
}

// Get individual UserWhiteList.
func (s *UserWhiteListsServiceOp) Get(ctx context.Context, userID int, id int) (*UserWhiteList, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(userWhiteListsBasePath, userID)
	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userWhiteListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserWhiteList, resp, err
}

// Create UserWhiteList.
func (s *UserWhiteListsServiceOp) Create(ctx context.Context, userID int, createRequest *UserWhiteListCreateRequest) (*UserWhiteList, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("UserWhiteList createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(userWhiteListsBasePath, userID) + apiFormat
	rootRequest := &userWhiteListCreateRequestRoot{
		UserWhiteListCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("UserWhiteList [Create] req: ", req)

	root := new(userWhiteListRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.UserWhiteList, resp, err
}

// Delete UserWhiteList.
func (s *UserWhiteListsServiceOp) Delete(ctx context.Context, userID int, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf(userWhiteListsBasePath, userID)
	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("UserWhiteList [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit UserWhiteList.
func (s *UserWhiteListsServiceOp) Edit(ctx context.Context, userID int, id int, editRequest *UserWhiteListEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("UserWhiteList [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf(userWhiteListsBasePath, userID)
	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("UserWhiteList [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
