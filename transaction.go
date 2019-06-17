package onappgo

import (
  "context"
  "fmt"
  "net/http"
  "time"

  "github.com/digitalocean/godo"
)

const (
  transactionsBasePath = "transactions"

  // TransactionRunning is a running transaction status
  TransactionRunning = "running"

  // TransactionCompleted is a completed transaction status
  TransactionCompleted = "complete"

  // TransactionPending is a pending transaction status
  TransactionPending = "pending"
)

// TransactionsService handles communction with action related methods of the
// OnApp API: https://docs.onapp.com/apim/latest/transactions
type TransactionsService interface {
  List(context.Context, *ListOptions) ([]Transaction, *Response, error)
  Get(context.Context, int) (*Transaction, *Response, error)
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
  ShutdownType     string   `json:"shutdown_type,omitempty"`
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
  Transaction *Transaction `json:"transaction"`
}

// List all transactions
func (s *TransactionsServiceOp) List(ctx context.Context, opt *ListOptions) ([]Transaction, *Response, error) {
  path := transactionsBasePath + apiFormat
  path, err := addOptions(path, opt)
  if err != nil {
    return nil, nil, err
  }

  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  var out []map[string]Transaction
  resp, err := s.client.Do(ctx, req, &out)
  if err != nil {
    return nil, resp, err
  }

  txs := make([]Transaction, len(out))
  for i := range txs {
    txs[i] = out[i]["transaction"]
  }

  return txs, resp, err
}

// Get an transaction by ID.
func (s *TransactionsServiceOp) Get(ctx context.Context, id int) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s/%d%s", transactionsBasePath, id, apiFormat)
  req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
  if err != nil {
    return nil, nil, err
  }

  root := new(transactionRoot)
  resp, err := s.client.Do(ctx, req, root)
  if err != nil {
    return nil, resp, err
  }

  return root.Transaction, resp, err
}

func (trx Transaction) String() string {
  return godo.Stringify(trx)
}

// Running check if transaction state is 'runing'
func (trx Transaction) Running() bool {
  return trx.Status == TransactionRunning
}

// Completed check if transaction state is 'complete'
func (trx Transaction) Completed() bool {
  return trx.Status == TransactionCompleted
}

// Pending check if transaction state is 'pending'
func (trx Transaction) Pending() bool {
  return trx.Status == TransactionPending
}

// Debug - print formatted Transaction structure
func (trx Transaction) Debug() {
  fmt.Printf("                      ID: %d\n",  trx.ID)
  fmt.Printf("              Identifier: %s\n",  trx.Identifier)
  fmt.Printf("                  Action: %s\n",  trx.Action)
  fmt.Printf("      AssociatedObjectID: %d\n",  trx.AssociatedObjectID)
  fmt.Printf("    AssociatedObjectType: %s\n",  trx.AssociatedObjectType)
  fmt.Printf("                  Status: %s\n",  trx.Status)
  fmt.Printf("               CreatedAt: %s\n",  trx.CreatedAt)
  fmt.Printf("              ParentType: %s\n",  trx.ParentType)
  fmt.Printf("                 ChainID: %d\n",  trx.ChainID)
  fmt.Printf("  DependentTransactionID: %d\n",  trx.DependentTransactionID)
  fmt.Printf("                  Params: %#v\n", trx.Params)
  fmt.Printf("     Params.ShutdownType: %s\n",  trx.Params.ShutdownType)
}
