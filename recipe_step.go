package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const recipeStepsBasePath string = recipesBasePath + "/%d/recipe_steps"
const recipeStepsSwapPath string = "%d/move_to/%d"

// RecipeStepsService is an interface for interfacing with the RecipeStep
// endpoints of the OnApp API
// https://docs.onapp.com/apim/latest/recipes/manage-recipe-steps
type RecipeStepsService interface {
	List(context.Context, int, *ListOptions) ([]RecipeStep, *Response, error)
	Get(context.Context, int, int) (*RecipeStep, *Response, error)
	Create(context.Context, int, *RecipeStepCreateRequest) (*RecipeStep, *Response, error)
	Delete(context.Context, int, int, interface{}) (*Response, error)
	Swap(context.Context, int, int, int, interface{}) (*Response, error)
}

// RecipeStepsServiceOp handles communication with the Data Store related methods of the
// OnApp API.
type RecipeStepsServiceOp struct {
	client *Client
}

var _ RecipeStepsService = &RecipeStepsServiceOp{}

type RecipeStep struct {
	ID               int    `json:"id"`
	RecipeID         int    `json:"recipe_id"`
	Number           int    `json:"number,omitempty"`
	Script           string `json:"script,omitempty"`
	OnSuccess        string `json:"on_success,omitempty"`
	OnFailure        string `json:"on_failure,omitempty"`
	SuccessGotoStep  int    `json:"success_goto_step"`
	CreatedAt        string `json:"created_at,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	ResultSource     string `json:"result_source,omitempty"`
	PassValues       string `json:"pass_values,omitempty"`
	PassAnythingElse bool   `json:"pass_anything_else"`
	FailValues       string `json:"fail_values,omitempty"`
	FailAnythingElse bool   `json:"fail_anything_else"`
	FailureGotoStep  int    `json:"failure_goto_step"`
}

// RecipeStepCreateRequest represents a request to create a RecipeStep
type RecipeStepCreateRequest struct {
	Script           string `json:"script"`
	ResultSource     string `json:"result_source"`
	PassAnythingElse bool   `json:"pass_anything_else"`
	PassValues       string `json:"pass_values"`
	OnSuccess        string `json:"on_success"`
	SuccessGotoStep  int    `json:"success_goto_step"`
	FailAnythingElse bool   `json:"fail_anything_else"`
	FailValues       string `json:"fail_values"`
	OnFailure        string `json:"on_failure"`
	FailureGotoStep  int    `json:"failure_goto_step"`
}

type RecipeSteps struct {
	RecipeStep RecipeStep `json:"recipe_step,omitempty"`
}

type recipeStepCreateRequestRoot struct {
	RecipeStepCreateRequest *RecipeStepCreateRequest `json:"recipe_step"`
}

type recipeStepRoot struct {
	RecipeStep *RecipeStep `json:"recipe_step"`
}

func (d RecipeStepCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all RecipeSteps.
func (s *RecipeStepsServiceOp) List(ctx context.Context, recipeID int, opt *ListOptions) ([]RecipeStep, *Response, error) {
	path := fmt.Sprintf(recipeStepsBasePath, recipeID) + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]RecipeStep
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]RecipeStep, len(out))
	for i := range arr {
		arr[i] = out[i]["recipe_step"]
	}

	return arr, resp, err
}

// Get individual RecipeStep.
func (s *RecipeStepsServiceOp) Get(ctx context.Context, recipeID int, recipeStepID int) (*RecipeStep, *Response, error) {
	if recipeStepID < 1 || recipeID < 1 {
		return nil, nil, godo.NewArgError("recipeID or recipeStepID", "cannot be less than 1")
	}

	path := fmt.Sprintf(recipeStepsBasePath, recipeID)
	path = fmt.Sprintf("%s/%d%s", path, recipeStepID, apiFormat)
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(recipeStepRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RecipeStep, resp, err
}

// Create RecipeStep.
func (s *RecipeStepsServiceOp) Create(ctx context.Context, recipeID int, createRequest *RecipeStepCreateRequest) (*RecipeStep, *Response, error) {
	if recipeID < 1 {
		return nil, nil, godo.NewArgError("recipeID", "cannot be less than 1")
	}

	if createRequest == nil {
		return nil, nil, godo.NewArgError("RecipeStep createRequest", "cannot be nil")
	}

	path := fmt.Sprintf(recipeStepsBasePath, recipeID) + apiFormat
	rootRequest := &recipeStepCreateRequestRoot{
		RecipeStepCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("RecipeStep [Create] req: ", req)

	root := new(recipeStepRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.RecipeStep, resp, err
}

// Delete RecipeStep.
func (s *RecipeStepsServiceOp) Delete(ctx context.Context, recipeID int, recipeStepID int, meta interface{}) (*Response, error) {
	if recipeStepID < 1 || recipeID < 1 {
		return nil, godo.NewArgError("recipeID or recipeStepID", "cannot be less than 1")
	}

	path := fmt.Sprintf(recipeStepsBasePath, recipeID)
	path = fmt.Sprintf("%s/%d%s", path, recipeStepID, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("RecipeStep [Delete] req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Swap RecipeStep.
func (s *RecipeStepsServiceOp) Swap(ctx context.Context, recipeID int, recipeStepIDFrom int, recipeStepIDTo int, meta interface{}) (*Response, error) {
	if recipeStepIDFrom < 1 || recipeStepIDTo < 1 || recipeID < 1 {
		return nil, godo.NewArgError("recipeID or recipeStepIDFrom or recipeStepIDTo ", "cannot be less than 1")
	}

	recipeStepPath := fmt.Sprintf(recipeStepsBasePath, recipeID)
	swapPath := fmt.Sprintf(recipeStepsSwapPath, recipeStepIDFrom, recipeStepIDTo)
	path := fmt.Sprintf("%s/%s%s", recipeStepPath, swapPath, apiFormat)
	path, err := addOptions(path, meta)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, nil)
	if err != nil {
		return nil, err
	}
	log.Println("RecipeStep [Swap] req: ", req)

	return s.client.Do(ctx, req, nil)
}
