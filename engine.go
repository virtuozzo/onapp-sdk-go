package onappgo

import (
	"context"
	"log"
	"net/http"
)

const engineBasePath string       = "sysadmin_tools/daemon/status"
const engineStartBasePath string  = "sysadmin_tools/daemon/start"
const engineStopBasePath string   = "sysadmin_tools/daemon/stop"
const engineReloadBasePath string = "sysadmin_tools/daemon/reload"

// EnginesService is an interface for interfacing with the Engines
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/onapp-engine
type EnginesService interface {
	Status(context.Context) (*Engine, *Response, error)
	Start(context.Context) (*Engine, *Response, error)
	Stop(context.Context) (*Engine, *Response, error)
	Reload(context.Context) (*Engine, *Response, error)
}

// EnginesServiceOp handles communication with the Engine related methods of the
// OnApp API.
type EnginesServiceOp struct {
	client *Client
}

var _ EnginesService = &EnginesServiceOp{}

type Engine struct {
	Status string `json:"status,omitempty"`
	IP     string `json:"ip,omitempty"`
}

// Status of Engine.
func (s *EnginesServiceOp) Status(ctx context.Context) (*Engine, *Response, error) {
	path := engineBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Engine [Status] req: ", req)

	root := &Engine{}
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Start Engine.
func (s *EnginesServiceOp) Start(ctx context.Context) (*Engine, *Response, error) {
	path := engineStartBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Engine [Start] req: ", req)

	root := &Engine{}
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Stop Engine.
func (s *EnginesServiceOp) Stop(ctx context.Context) (*Engine, *Response, error) {
	path := engineStopBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Engine [Stop] req: ", req)

	root := &Engine{}
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}

// Reload Engine.
func (s *EnginesServiceOp) Reload(ctx context.Context) (*Engine, *Response, error) {
	path := engineStopBasePath + apiFormat

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Engine [Reload] req: ", req)

	root := &Engine{}
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root, resp, err
}
