package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const resolverBasePath string = "settings/nameservers"

// ResolversService is an interface for interfacing with the Resolver
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/firewall-rules-for-vss
type ResolversService interface {
	List(context.Context, *ListOptions) ([]Resolver, *Response, error)
	Get(context.Context, int) (*Resolver, *Response, error)
	Create(context.Context, *ResolverCreateRequest) (*Resolver, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *ResolverCreateRequest) (*Response, error)
}

// ResolversServiceOp handles communication with the Resolvers related methods of the
// OnApp API.
type ResolversServiceOp struct {
	client *Client
}

var _ ResolversService = &ResolversServiceOp{}

// Resolver -
// https://docs.onapp.com/apim/latest/resolvers
type Resolver struct {
	Address   string `json:"address,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	ID        int    `json:"id,omitempty"`
	NetworkID int    `json:"network_id,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

// ResolverCreateRequest represents a request to create a Resolver
type ResolverCreateRequest struct {
	Address   string `json:"address,omitempty"`
	NetworkID int    `json:"network_id,omitempty"`
}

type resolverCreateRequestRoot struct {
	ResolverCreateRequest *ResolverCreateRequest `json:"nameserver"`
}

type resolverRoot struct {
	Resolver *Resolver `json:"nameserver"`
}

func (d ResolverCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Resolvers
func (s *ResolversServiceOp) List(ctx context.Context, opt *ListOptions) ([]Resolver, *Response, error) {
	path := resolverBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]Resolver
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Resolver, len(out))
	for i := range arr {
		arr[i] = out[i]["nameserver"]
	}

	return arr, resp, err
}

// Get individual Resolver
func (s *ResolversServiceOp) Get(ctx context.Context, id int) (*Resolver, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", resolverBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(resolverRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Resolver, resp, err
}

// Create Resolver
func (s *ResolversServiceOp) Create(ctx context.Context, createRequest *ResolverCreateRequest) (*Resolver, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("Resolver createRequest", "cannot be nil")
	}

	path := resolverBasePath + apiFormat
	rootRequest := &resolverCreateRequestRoot{
		ResolverCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("Resolver [Create] req: ", req)

	root := new(resolverRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Resolver, resp, err
}

// Delete Resolver
func (s *ResolversServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("vmID || id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", resolverBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("Resolver [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit Resolver
func (s *ResolversServiceOp) Edit(ctx context.Context, id int, editRequest *ResolverCreateRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("Resolver [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", resolverBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Resolver [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
