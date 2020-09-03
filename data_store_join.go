package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

var dataStoreJoinPaths = map[string]string{
	"Hypervisor":      "settings/hypervisors/%d/data_store_joins",
	"HypervisorGroup": "settings/hypervisor_zones/%d/data_store_joins",
}

// DataStoreJoinsService is an interface for interfacing with the DataStoreJoin
type DataStoreJoinsService interface {
	List(context.Context, *DataStoreJoinCreateRequest, *ListOptions) ([]DataStoreJoin, *Response, error)
	Get(context.Context, string, int, int) (*DataStoreJoin, *Response, error)
	Create(context.Context, *DataStoreJoinCreateRequest) (*DataStoreJoin, *Response, error)
	Delete(context.Context, *DataStoreJoinDeleteRequest, interface{}) (*Response, error)
}

// DataStoreJoinsServiceOp -
type DataStoreJoinsServiceOp struct {
	client *Client
}

var _ DataStoreJoinsService = &DataStoreJoinsServiceOp{}

// DataStoreJoin represents a DataStoreJoin
type DataStoreJoin struct {
	ID             int    `json:"id,omitempty"`
	DataStoreID    int    `json:"data_store_id,omitempty"`
	CreatedAt      string `json:"created_at,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
	TargetJoinID   int    `json:"target_join_id,omitempty"`
	TargetJoinType string `json:"target_join_type,omitempty"`
	Identifier     string `json:"identifier,omitempty"`
}

// DataStoreJoinCreateRequest represents a request to create a DataStoreJoin
type DataStoreJoinCreateRequest struct {
	DataStoreID    int    `json:"data_store_id,omitempty"`
	TargetJoinID   int    `json:"-"`
	TargetJoinType string `json:"-"`
}

// DataStoreJoinDeleteRequest represents a request to delete a DataStoreJoin
type DataStoreJoinDeleteRequest struct {
	ID             int
	TargetJoinID   int
	TargetJoinType string
}

type dataStoreJoinCreateRequestRoot struct {
	DataStoreID int `json:"data_store_id,omitempty"`
}

type dataStoreJoinRoot struct {
	DataStoreJoin *DataStoreJoin `json:"data_store_join"`
}

func (d DataStoreJoinCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all DataStoreJoins.
func (s *DataStoreJoinsServiceOp) List(ctx context.Context, createRequest *DataStoreJoinCreateRequest, opt *ListOptions) ([]DataStoreJoin, *Response, error) {
	path := ""
	if val, ok := dataStoreJoinPaths[createRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
	} else {
		return nil, nil, godo.NewArgError("DataStoreJoin List: map key not found", createRequest.TargetJoinType)
	}

	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]DataStoreJoin
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]DataStoreJoin, len(out))
	for i := range arr {
		arr[i] = out[i]["data_store_join"]
	}

	return arr, resp, err
}

// Get individual DataStoreJoin.
func (s *DataStoreJoinsServiceOp) Get(ctx context.Context, targetJoinType string, targetJoinID int, id int) (*DataStoreJoin, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := ""
	if val, ok := dataStoreJoinPaths[targetJoinType]; ok {
		path = fmt.Sprintf(val, targetJoinID)
	} else {
		return nil, nil, godo.NewArgError("DataStoreJoin Get: map key not found", targetJoinType)
	}

	path = fmt.Sprintf("%s/%d%s", path, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(dataStoreJoinRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DataStoreJoin, resp, err
}

// Create DataStoreJoin.
func (s *DataStoreJoinsServiceOp) Create(ctx context.Context, createRequest *DataStoreJoinCreateRequest) (*DataStoreJoin, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("DataStoreJoin createRequest", "cannot be nil")
	}

	path := ""
	if val, ok := dataStoreJoinPaths[createRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
	} else {
		return nil, nil, godo.NewArgError("DataStoreJoin Create: map key not found", createRequest.TargetJoinType)
	}

	rootRequest := &dataStoreJoinCreateRequestRoot{
		DataStoreID: createRequest.DataStoreID,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("DataStoreJoin [Create] req: ", req)

	root := new(dataStoreJoinRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.DataStoreJoin, resp, err
}

// Delete DataStoreJoin.
func (s *DataStoreJoinsServiceOp) Delete(ctx context.Context, deleteRequest *DataStoreJoinDeleteRequest, meta interface{}) (*Response, error) {
	if deleteRequest == nil {
		return nil, godo.NewArgError("DataStoreJoin deleteRequest", "cannot be nil")
	}

	if deleteRequest.ID < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := ""
	if val, ok := dataStoreJoinPaths[deleteRequest.TargetJoinType]; ok {
		path = fmt.Sprintf(val, deleteRequest.TargetJoinID)
	} else {
		return nil, godo.NewArgError("DataStoreJoin Delete: wrong TargetJoinType", deleteRequest.TargetJoinType)
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
	fmt.Println("DataStoreJoin [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}
