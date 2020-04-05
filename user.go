package onappgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/digitalocean/godo"
)

const usersBasePath string = "users"

// UsersService is an interface for interfacing with the User
// endpoints of the OnApp API
// See: https://docs.onapp.com/apim/latest/users
type UsersService interface {
	List(context.Context, *ListOptions) ([]User, *Response, error)
	Get(context.Context, int) (*User, *Response, error)
	Create(context.Context, *UserCreateRequest) (*User, *Response, error)
	Delete(context.Context, int, interface{}) (*Response, error)
	Edit(context.Context, int, *UserEditRequest) (*Response, error)
}

// UsersServiceOp handles communication with the User related methods of the
// OnApp API.
type UsersServiceOp struct {
	client *Client
}

var _ UsersService = &UsersServiceOp{}

// Infoboxes -
type Infoboxes struct {
	DisplayInfoboxes bool     `json:"display_infoboxes,bool"`
	HiddenInfoboxes  []string `json:"hidden_infoboxes,omitempty"`
}

// Roles -
type Roles struct {
	Role *Role `json:"role,omitempty"`
}

// User -
type User struct {
	ActivatedAt             string             `json:"activated_at,omitempty"`
	Avatar                  interface{}        `json:"avatar,omitempty"`
	BillingPlanID           int                `json:"billing_plan_id,omitempty"`
	CdnAccountStatus        string             `json:"cdn_account_status,omitempty"`
	CdnStatus               string             `json:"cdn_status,omitempty"`
	CreatedAt               string             `json:"created_at,omitempty"`
	DeletedAt               string             `json:"deleted_at,omitempty"`
	Email                   string             `json:"email,omitempty"`
	FirewallID              int                `json:"firewall_id,omitempty"`
	FirstName               string             `json:"first_name,omitempty"`
	GroupID                 int                `json:"group_id,omitempty"`
	ID                      int                `json:"id,omitempty"`
	Identifier              string             `json:"identifier,omitempty"`
	ImageTemplateGroupID    int                `json:"image_template_group_id,omitempty"`
	Infoboxes               Infoboxes          `json:"infoboxes,omitempty"`
	LastName                string             `json:"last_name,omitempty"`
	Locale                  string             `json:"locale,omitempty"`
	Login                   string             `json:"login,omitempty"`
	PasswordChangedAt       string             `json:"password_changed_at,omitempty"`
	RegisteredYubikey       bool               `json:"registered_yubikey,bool"`
	Status                  string             `json:"status,omitempty"`
	Supplied                bool               `json:"supplied,bool"`
	SuspendAt               string             `json:"suspend_at,omitempty"`
	SystemTheme             string             `json:"system_theme,omitempty"`
	TimeZone                string             `json:"time_zone,omitempty"`
	UpdatedAt               string             `json:"updated_at,omitempty"`
	UseGravatar             bool               `json:"use_gravatar,bool"`
	UserGroupID             int                `json:"user_group_id,omitempty"`
	BucketID                int                `json:"bucket_id,omitempty"`
	UsedCpus                int                `json:"used_cpus,omitempty"`
	UsedMemory              int                `json:"used_memory,omitempty"`
	UsedCPUShares           int                `json:"used_cpu_shares,omitempty"`
	UsedDiskSize            int                `json:"used_disk_size,omitempty"`
	MemoryAvailable         float64            `json:"memory_available,omitempty"`
	DiskSpaceAvailable      float64            `json:"disk_space_available,omitempty"`
	Roles                   []Roles            `json:"roles,omitempty"`
	MonthlyPrice            float64            `json:"monthly_price,omitempty"`
	PaymentAmount           float64            `json:"payment_amount,omitempty"`
	OutstandingAmount       float64            `json:"outstanding_amount,omitempty"`
	TotalAmount             float64            `json:"total_amount,omitempty"`
	DiscountDueToFree       float64            `json:"discount_due_to_free,omitempty"`
	TotalAmountWithDiscount float64            `json:"total_amount_with_discount,omitempty"`
	AdditionalFields        []AdditionalFields `json:"additional_fields,omitempty"`
	UsedIPAddresses         []IPAddresses      `json:"used_ip_addresses,omitempty"`
}

// UserCreateRequest -
type UserCreateRequest struct {
	Login            string              `json:"login,omitempty"`
	Email            string              `json:"email,omitempty"`
	FirstName        string              `json:"first_name,omitempty"`
	LastName         string              `json:"last_name,omitempty"`
	Password         string              `json:"password,omitempty"`
	UserGroupID      int                 `json:"user_group_id,omitempty"`
	BucketID         int                 `json:"bucket_id,omitempty"`
	RoleIDs          []string            `json:"role_ids,omitempty"`
	AdditionalFields []*AdditionalFields `json:"additional_fields,omitempty"`
}

// UserEditRequest -
type UserEditRequest struct {
	Email             string              `json:"email,omitempty"`
	FirstName         string              `json:"first_name,omitempty"`
	LastName          string              `json:"last_name,omitempty"`
	Password          string              `json:"password,omitempty"`
	UserGroupID       int                 `json:"user_group_id,omitempty"`
	BucketID          int                 `json:"bucket_id,omitempty"`
	RoleIDs           []string            `json:"role_ids,omitempty"`
	AdditionalFields  []*AdditionalFields `json:"additional_fields,omitempty"`
	SuspendAt         string              `json:"suspend_at,omitempty"`
	RegisteredYubikey bool                `json:"registered_yubikey,bool"`
}

type userCreateRequestRoot struct {
	UserCreateRequest *UserCreateRequest `json:"user"`
}

type userRoot struct {
	User *User `json:"user"`
}

func (d UserCreateRequest) String() string {
	return godo.Stringify(d)
}

// List all Users.
func (s *UsersServiceOp) List(ctx context.Context, opt *ListOptions) ([]User, *Response, error) {
	path := usersBasePath + apiFormat
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	var out []map[string]User
	resp, err := s.client.Do(ctx, req, &out)
	if err != nil {
		return nil, resp, err
	}

	arr := make([]User, len(out))
	for i := range arr {
		arr[i] = out[i]["user"]
	}

	return arr, resp, err
}

// Get individual User.
func (s *UsersServiceOp) Get(ctx context.Context, id int) (*User, *Response, error) {
	if id < 1 {
		return nil, nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", usersBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(userRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.User, resp, err
}

// Create User.
func (s *UsersServiceOp) Create(ctx context.Context, createRequest *UserCreateRequest) (*User, *Response, error) {
	if createRequest == nil {
		return nil, nil, godo.NewArgError("createRequest", "cannot be nil")
	}

	path := usersBasePath + apiFormat
	rootRequest := &userCreateRequestRoot{
		UserCreateRequest: createRequest,
	}

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, rootRequest)
	if err != nil {
		return nil, nil, err
	}
	log.Println("User [Create]  req: ", req)

	root := new(userRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.User, resp, err
}

// UserDeleteRequest -
type UserDeleteRequest struct {
	// Force int `url:"force"`
	Force int `json:"force,omitempty"`
}

// Delete User.
func (s *UsersServiceOp) Delete(ctx context.Context, id int, meta interface{}) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	path := fmt.Sprintf("%s/%d%s", usersBasePath, id, apiFormat)
	path, err := addOptions(path, nil)
	if err != nil {
		return nil, err
	}

	opts := &UserDeleteRequest{
		Force: 1,
	}

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, opts)
	if err != nil {
		return nil, err
	}
	log.Println("User [Delete]  req: ", req)

	return s.client.Do(ctx, req, nil)
}

// Edit User
func (s *UsersServiceOp) Edit(ctx context.Context, id int, editRequest *UserEditRequest) (*Response, error) {
	if id < 1 {
		return nil, godo.NewArgError("id", "cannot be less than 1")
	}

	if editRequest == nil {
		return nil, godo.NewArgError("User [Edit] editRequest", "cannot be nil")
	}

	path := fmt.Sprintf("%s/%d%s", usersBasePath, id, apiFormat)

	req, err := s.client.NewRequest(ctx, http.MethodPut, path, editRequest)
	if err != nil {
		return nil, err
	}
	log.Println("User [Edit]  req: ", req)

	return s.client.Do(ctx, req, nil)
}
