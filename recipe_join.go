package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const controlPanel string = "ControlPanel"

var recipeJoinPaths = map[string]string{
	controlPanel:        "settings/control_panel/recipe_joins",
	"HypervisorGroup":   "settings/hypervisor_zones/%d/recipe_joins",
	"VirtualServer":     "virtual_machines/%d/recipe_joins",
	"ImageTemplateBase": "templates/%d/recipe_joins",
	"SmartServer":       "smart_servers/%d/recipe_joins",
}

// RecipeJoinsService is an interface for interfacing with the RecipeJoin
type RecipeJoinsService interface {
	List(context.Context, *RecipeJoinCreateRequest, *ListOptions) (map[string]interface{}, *Response, error)
	Create(context.Context, *RecipeJoinCreateRequest) (*RecipeJoin, *Response, error)
	Delete(context.Context, *RecipeJoinDeleteRequest, interface{}) (*Response, error)
}

// RecipeJoinsServiceOp -
type RecipeJoinsServiceOp struct {
	client *Client
}

var _ RecipeJoinsService = &RecipeJoinsServiceOp{}

// RecipeJoin represents a RecipeJoin
type RecipeJoin struct {
	CreatedAt      string `json:"created_at,omitempty"`
	EventType      string `json:"event_type"`
	ID             int    `json:"id"`
	RecipeID       int    `json:"recipe_id"`
	TargetJoinID   int    `json:"target_join_id,omitempty"`
	TargetJoinType string `json:"target_join_type,omitempty"`
	UpdatedAt      string `json:"updated_at,omitempty"`
}

// RecipeJoinCreateRequest represents a request to create a ControlPanel RecipeJoin
type RecipeJoinCreateRequest struct {
	EventType      string `json:"event_type"`
	RecipeID       int    `json:"recipe_id"`
	TargetJoinID   int    `json:"-"`
	TargetJoinType string `json:"-"`
}

// RecipeJoinDeleteRequest represents a request to delete a RecipeJoin
type RecipeJoinDeleteRequest struct {
	ID             int
	TargetJoinID   int
	TargetJoinType string
}

type recipeJoinCreateRequestRoot struct {
	RecipeJoinCreateRequest *RecipeJoinCreateRequest `json:"recipe_join"`
}

type recipeJoinRoot struct {
	RecipeJoin *RecipeJoin `json:"recipe_join"`
}

func (d RecipeJoinCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all ControlPanelRecipeJoins.
func (s *RecipeJoinsServiceOp) List(ctx context.Context, createRequest *RecipeJoinCreateRequest, opt *ListOptions) (map[string]interface{}, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("RecipeJoin createRequest [List]", "cannot be nil")
	}

	path := ""
	if val, ok := recipeJoinPaths[createRequest.TargetJoinType]; ok {
		if createRequest.TargetJoinType == controlPanel {
			path = val + apiFormat
		} else {
			if createRequest.TargetJoinID < 1 {
				return nil, nil, godo.NewArgError("id", "cannot be less than 1")
			}
			path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
		}
	} else {
		return nil, nil, godo.NewArgError("RecipeJoin List: wrong TargetJoinType", createRequest.TargetJoinType)
	}

	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out map[string]interface{}
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	// log.Println("RecipeJoin [List] response: ", out)

	return out, resp, err
}

// Create ComputeZoneRecipeJoin.
func (s *RecipeJoinsServiceOp) Create(ctx context.Context, createRequest *RecipeJoinCreateRequest) (*RecipeJoin, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("RecipeJoin createRequest [Create]", "cannot be nil")
	}

	path := ""
	if val, ok := recipeJoinPaths[createRequest.TargetJoinType]; ok {
		if createRequest.TargetJoinType == controlPanel {
			path = val + apiFormat
		} else {
			if createRequest.TargetJoinID < 1 {
				return nil, nil, godo.NewArgError("id", "cannot be less than 1")
			}
			path = fmt.Sprintf(val, createRequest.TargetJoinID) + apiFormat
		}
	} else {
		return nil, nil, godo.NewArgError("RecipeJoin Create: wrong TargetJoinType", createRequest.TargetJoinType)
	}

	rootRequest := &recipeJoinCreateRequestRoot{
		RecipeJoinCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("RecipeJoin [Create] req: ", req)

	root := new(recipeJoinRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RecipeJoin, resp, err
}

// Delete ComputeZoneRecipeJoin.
func (s *RecipeJoinsServiceOp) Delete(ctx context.Context, deleteRequest *RecipeJoinDeleteRequest, meta interface{}) (*Response, error) {
	if deleteRequest == nil {
		return nil, godo.NewArgError("RecipeJoin deleteRequest", "cannot be nil")
	}

	if deleteRequest.ID < 1 || deleteRequest.TargetJoinID < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := ""
	if val, ok := recipeJoinPaths[deleteRequest.TargetJoinType]; ok {
		if deleteRequest.TargetJoinType == controlPanel {
			path = val
		} else {
			path = fmt.Sprintf(val, deleteRequest.TargetJoinID)
		}
	} else {
		return nil, godo.NewArgError("RecipeJoin Delete: wrong TargetJoinType", deleteRequest.TargetJoinType)
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
	fmt.Println("RecipeJoin [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}
