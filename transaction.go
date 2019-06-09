package onappgo

import (
  "context"
  "fmt"
  "net/http"
  "time"

  "github.com/digitalocean/godo"
)

const (
  transActionsBasePath = "transactions"

  // TransActionInProgress is an in progress action status
  TransActionInProgress = "in-progress"

  // TransActionCompleted is a completed action status
  TransActionCompleted = "completed"
)

// TransactionsService handles communction with action related methods of the
// OnApp API: https://docs.onapp.com/apim/latest/transactions
type TransactionsService interface {
  List(context.Context, *ListOptions) ([]Transaction, *Response, error)
  // Get(context.Context, int) (*Transaction, *Response, error)
}

// TransactionsServiceOp handles communition with the image action related methods of the
// OnApp API.
type TransactionsServiceOp struct {
  client *Client
}

var _ TransactionsService = &TransactionsServiceOp{}

// Params represents a OnApp Transaction params
type Params struct {
  DestroyMsg       string   `json:"destroy_msg,omitempty"`
  InitiatorID      int      `json:"initiator_id,omitempty"`
  RealUserID       int      `json:"real_user_id,omitempty"`
  RemoteIP         string   `json:"remote_ip,omitempty"`
  SkipNotification bool     `json:"skip_notification,bool,omitempty"`
}

// Transaction represents a OnApp Transaction
type Transaction struct {
  Action                  string        `json:"action,omitempty"`
  Actor                   string        `json:"actor,omitempty"`
  AllowedCancel           bool          `json:"allowed_cancel,bool,omitempty"`
  AssociatedObjectID      int           `json:"associated_object_id,omitempty"`
  AssociatedObjectType    string        `json:"associated_object_type,omitempty"`
  ChainID                 int           `json:"chain_id,omitempty"`
  CreatedAt               time.Time     `json:"created_at,omitempty"`
  DependentTransactionID  int           `json:"dependent_transaction_id,omitempty"`
  ID                      int           `json:"id,omitempty"`
  Identifier              string        `json:"identifier,omitempty"`
  LockVersion             int           `json:"lock_version,omitempty"`
  Params                  Params        `json:"params,omitempty"`
  ParentID                int           `json:"parent_id,omitempty"`
  ParentType              string        `json:"parent_type,omitempty"`
  Pid                     int           `json:"pid,omitempty"`
  Priority                int           `json:"priority,omitempty"`
  Scheduled               bool          `json:"scheduled,bool,omitempty"`
  StartAfter              time.Time     `json:"start_after,omitempty"`
  StartedAt               time.Time     `json:"started_at,omitempty"`
  Status                  string        `json:"status,omitempty"`
  UpdatedAt               time.Time     `json:"updated_at,omitempty"`
  UserID                  int           `json:"user_id,omitempty"`
}

type transactionRoot struct {
  Event *Transaction `json:"transaction"`
}

type transactionsRoot struct {
  Transactions []Transaction `json:"transaction"`
}

// List all transactions
func (s *TransactionsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Transaction, *Response, error) {
  path := transActionsBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  // root := new(transactionsRoot)
  // resp, err := s.client.Do(ctx, req, root)
  // if err != nil {
  //   return nil, resp, err
  // }

  // return root.Transactions, resp, err

  var out []map[string]Transaction
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  fmt.Printf("out: %#v\n", out)
  txs := make([]Transaction, len(out))
  for i := range txs {
    txs[i] = out[i]["transaction"]
  }

  return txs, resp, err
}

// Get an transaction by ID.
// func (s *TransactionsServiceOp) Get(ctx context.Context, id int) (*Transaction, *Response, error) {
//   if id < 1 {
//     return nil, nil, NewArgError("id", "cannot be less than 1")
//   }

//   path := fmt.Sprintf("%s/%d%s", transActionsBasePath, id, apiFormat)
//   req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
//   if err != nil {
//     return nil, nil, err
//   }

//   root := new(transactionRoot)
//   resp, err := s.client.Do(ctx, req, root)
//   if err != nil {
//     return nil, resp, err
//   }

//   return root.Event, resp, err
// }

func (a Transaction) String() string {
  return godo.Stringify(a)
}
