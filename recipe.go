package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const recipesBasePath string = "recipes"

// RecipesService is an interface for interfacing with the Recipe
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/recipes
type RecipesService interface {
	List(context.Context, *ListOptions) ([]Recipe, *Response, error)
	Get(context.Context, int) (*Recipe, *Response, error)
	Create(context.Context, *RecipeCreateRequest) (*Recipe, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *RecipeCreateRequest) (*Response, error)
}

// RecipesServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type RecipesServiceOp struct {
	client *Client
}

var _ RecipesService = &RecipesServiceOp{}

// Recipe represents a Recipe
type Recipe struct {
	ID             int           `json:"id,omitempty"`
	UserID         int           `json:"user_id,omitempty"`
	CreatedAt      string        `json:"created_at,omitempty"`
	UpdatedAt      string        `json:"updated_at,omitempty"`
	Label          string        `json:"label,omitempty"`
	Description    string        `json:"description,omitempty"`
	ScriptType     string        `json:"script_type,omitempty"`
	CompatibleWith string        `json:"compatible_with,omitempty"`
	RecipeSteps    []RecipeSteps `json:"recipe_steps,omitempty"`
}

// RecipeCreateRequest represents a request to create a Recipe
type RecipeCreateRequest struct {
	Label          string `json:"label"`
	Description    string `json:"description"`
	CompatibleWith string `json:"compatible_with"`
	ScriptType     string `json:"script_type"`
}

type recipeCreateRequestRoot struct {
	RecipeCreateRequest *RecipeCreateRequest `json:"recipe"`
}

type recipeRoot struct {
	Recipe *Recipe `json:"recipe"`
}

func (d RecipeCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Recipes.
func (s *RecipesServiceOp) List(ctx context.Context, opt *ListOptions) ([]Recipe, *Response, error) {
	path := recipesBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]Recipe
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]Recipe, len(out))
	for i := range arr {
		arr[i] = out[i]["recipe"]
	}

	return arr, resp, err
}

// Get individual Recipe.
func (s *RecipesServiceOp) Get(ctx context.Context, id int) (*Recipe, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", recipesBasePath, id, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(recipeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Recipe, resp, err
}

// Create Recipe.
func (s *RecipesServiceOp) Create(ctx context.Context, createRequest *RecipeCreateRequest) (*Recipe, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("Recipe createRequest", "cannot be nil")
	}

	path := recipesBasePath + apiFormat
	rootRequest := &recipeCreateRequestRoot{
		RecipeCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("Recipe [Create] req: ", req)

	root := new(recipeRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Recipe, resp, err
}

// Delete Recipe.
func (s *RecipesServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", recipesBasePath, id, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("Recipe [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit Recipe.
func (s *RecipesServiceOp) Edit(ctx context.Context, id int, editRequest *RecipeCreateRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("Recipe [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", recipesBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("Recipe [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
