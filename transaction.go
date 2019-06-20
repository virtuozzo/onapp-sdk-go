package onappgo

import (
  "context"
  "fmt"
  "net/http"
  "time"
  "reflect"

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

  // TransactionCancelled is a cancelled transaction status
  TransactionCancelled = "cancelled"

  // TransactionFailed is a failed transaction status
  TransactionFailed = "failed"
)

// TransactionsService handles communction with action related methods of the
// OnApp API: https://docs.onapp.com/apim/latest/transactions
type TransactionsService interface {
  List(context.Context, *ListOptions) ([]Transaction, *Response, error)
  Get(context.Context, int) (*Transaction, *Response, error)

  GetByFilter(context.Context, int, interface{}, *ListOptions) (*Transaction, *Response, error)
}

// TransactionsServiceOp handles communition with the image action related methods of the
// OnApp API.
type TransactionsServiceOp struct {
  client *Client
}

var _ TransactionsService = &TransactionsServiceOp{}

// Params represents a OnApp Transaction params
// type Params struct {
//   DestroyMsg       string   `json:"destroy_msg,omitempty"`
//   InitiatorID      int      `json:"initiator_id,omitempty"`
//   RealUserID       int      `json:"real_user_id,omitempty"`
//   RemoteIP         string   `json:"remote_ip,omitempty"`
//   SkipNotification bool     `json:"skip_notification,bool,omitempty"`
//   ShutdownType     string   `json:"shutdown_type,omitempty"`
// }

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
  // Params                  Params        `json:"params,omitempty"`
  Params                  interface{}   `json:"params,omitempty"`
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

// GetByFilter find transaction with specified fields for virtual machine by ID.
func (s *TransactionsServiceOp) GetByFilter(ctx context.Context, vmID int, filter interface{}, opt *ListOptions) (*Transaction, *Response, error) {
  if vmID < 1 {
    return nil, nil, godo.NewArgError("vmID", "cannot be less than 1")
  }

  path := fmt.Sprintf("%s%s", transactionsBasePath, apiFormat)
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

  var root *Transaction
  for i := range out {
    trx := out[i]["transaction"]

    root = &trx
    if root.equal(filter) {
      break
    } else {
      root = nil
    }
  }

  if root != nil {
    return root, resp, err
  }

  return nil, nil, fmt.Errorf("Transaction not found or wrong filter [%+v].\n", filter)
}

func (trx *Transaction) equal(filter interface{}) bool {
  val := reflect.ValueOf(filter)
  filterFields := reflect.Indirect(reflect.ValueOf(trx))

  // fmt.Printf("        equal.filterFields: %#v\n", filterFields)
  for i := 0; i < val.NumField(); i++ {
    typeField := val.Type().Field(i)
    value := val.Field(i)
    filterValue := filterFields.FieldByName(typeField.Name)

    // fmt.Printf("%s: %s[%#v] -> %s[%#v]\n", typeField.Name, value.Type(), value, filterValue.Type(), filterValue)

    if value.Interface() != filterValue.Interface() {
      // fmt.Printf("[false] return on filed [%s]\n\n", typeField.Name)
      return false
    }
  }

  // fmt.Printf("[true] transaction found with id[%d]\n\n", trx.ID)
  return true
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

// NotCompleted check if transaction state is
// 'running' or 'pending' or 'failed' or 'cancelled'
func (trx Transaction) NotCompleted() bool {
  return trx.Running() || trx.Pending() || trx.Failed() || trx.Cancelled()
}

// Pending check if transaction state is 'pending'
func (trx Transaction) Pending() bool {
  return trx.Status == TransactionPending
}

// Failed check if transaction state is 'failed'
func (trx Transaction) Failed() bool {
  return trx.Status == TransactionFailed
}

// Cancelled check if transaction state is 'cancelled'
func (trx Transaction) Cancelled() bool {
  return trx.Status == TransactionCancelled
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
  fmt.Printf("                  Params: %+v\n", trx.Params)
  // fmt.Printf("     Params.ShutdownType: %s\n",  trx.Params.ShutdownType)
}
