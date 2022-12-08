package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const recipeGroupsBasePath string = "recipe_groups"

// RecipeGroupsService is an interface for interfacing with the RecipeGroup
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/recipe_group-groups
type RecipeGroupsService interface {
	List(context.Context, *ListOptions) ([]RecipeGroup, *Response, error)
	Get(context.Context, int) (*RecipeGroup, *Response, error)
	Create(context.Context, *RecipeGroupCreateRequest) (*RecipeGroup, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *RecipeGroupEditRequest) (*Response, error)
}

// RecipeGroupsServiceOp handles communication with the RecipeGroups related methods of the
// OnApp API.
type RecipeGroupsServiceOp struct {
	client *Client
}

var _ RecipeGroupsService = &RecipeGroupsServiceOp{}

type Children struct {
	ID        int        `json:"id,omitempty"`
	Label     string     `json:"label,omitempty"`
	ParentID  int        `json:"parent_id,omitempty"`
	Lft       int        `json:"lft,omitempty"`
	Rgt       int        `json:"rgt,omitempty"`
	Depth     int        `json:"depth,omitempty"`
	CreatedAt string     `json:"created_at,omitempty"`
	UpdatedAt string     `json:"updated_at,omitempty"`
	Childrens []Children `json:"children,omitempty"` // rename filed but JSON leave as is
	Relations []Relation `json:"relations,omitempty"`
}

type Relation struct {
	ID            int    `json:"id,omitempty"`
	RecipeID      int    `json:"recipe_id,omitempty"`
	RecipeGroupID int    `json:"recipe_group_id,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
	Recipe        Recipe `json:"recipe,omitempty"`
}

// RecipeGroup represents a RecipeGroup
type RecipeGroup struct {
	ID        int        `json:"id,omitempty"`
	Label     string     `json:"label,omitempty"`
	ParentID  int        `json:"parent_id,omitempty"`
	Lft       int        `json:"lft,omitempty"`
	Rgt       int        `json:"rgt,omitempty"`
	Depth     int        `json:"depth,omitempty"`
	CreatedAt string     `json:"created_at,omitempty"`
	UpdatedAt string     `json:"updated_at,omitempty"`
	Childrens []Children `json:"children,omitempty"` // rename filed but JSON leave as is
	Relations []Relation `json:"relations,omitempty"`
}

// RecipeGroupCreateRequest represents a request to create a RecipeGroup
type RecipeGroupCreateRequest struct {
	Label string `json:"label,omitempty"`
}

// RecipeGroupEditRequest represents a request to edit a RecipeGroup
type RecipeGroupEditRequest struct {
	Label string `json:"label,omitempty"`
}

type recipeGroupCreateRequestRoot struct {
	RecipeGroupCreateRequest *RecipeGroupCreateRequest `json:"recipe_group"`
}

// recipeGroupRoot - used to get one RecipeGroup
type recipeGroupRoot struct {
	RecipeGroup *RecipeGroup `json:"recipe_group"`
}

func (d RecipeGroupCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Recipe Groups.
func (s *RecipeGroupsServiceOp) List(ctx context.Context, opt *ListOptions) ([]RecipeGroup, *Response, error) {
	path := recipeGroupsBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []RecipeGroup
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	return out, resp, err
}

// Get individual Recipe Group.
func (s *RecipeGroupsServiceOp) Get(ctx context.Context, id int) (*RecipeGroup, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", recipeGroupsBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	// OnApp 6.8 - /recipe_groups/1.json
	// {
	//   "recipe_group":{
	//     "id":1,
	//     "label":"recipe_group1",
	//     "parent_id":null,
	//     "lft":1,
	//     "rgt":10,
	//     "depth":0,
	//     "created_at":"2022-12-08T12:11:40.000Z",
	//     "updated_at":"2022-12-08T12:35:02.000Z"
	//   }
	// }
	// BUG of OnApp the JSON for RecipeGroup didn't return fields 'children' and 'relations'
	root := new(recipeGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RecipeGroup, resp, err
}

// Create Recipe Group.
func (s *RecipeGroupsServiceOp) Create(ctx context.Context, createRequest *RecipeGroupCreateRequest) (*RecipeGroup, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("RecipeGroup createRequest", "cannot be nil")
	}

	path := recipeGroupsBasePath + apiFormat
	rootRequest := &recipeGroupCreateRequestRoot{
		RecipeGroupCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("RecipeGroup [Create] req: ", req)

	root := new(recipeGroupRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RecipeGroup, resp, err
}

// Delete RecipeGroup.
func (s *RecipeGroupsServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", recipeGroupsBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("RecipeGroup [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit RecipeGroup.
func (s *RecipeGroupsServiceOp) Edit(ctx context.Context, id int, editRequest *RecipeGroupEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("RecipeGroup [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", recipeGroupsBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("RecipeGroup [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
