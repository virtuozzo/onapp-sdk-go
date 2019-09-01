package onappgo

import (
  "context"
  "fmt"
  "net/http"
  "reflect"
  "container/list"

  "github.com/digitalocean/godo"
)

const (
  transactionsBasePath = "transactions"

  // TransactionRunning is a running transaction status
  TransactionRunning = "running"

  // TransactionComplete is a completed transaction status
  TransactionComplete = "complete"

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
  // ListByGroup(context.Context, int, string, *ListOptions) (*list.List, *Response, error)
  ListByGroup(context.Context, interface{}, *ListOptions) (*list.List, *Response, error)
}

// TransactionsServiceOp handles communition with the image action related methods of the
// OnApp API.
type TransactionsServiceOp struct {
  client *Client
}

var _ TransactionsService = &TransactionsServiceOp{}

// Transaction represents a OnApp Transaction
type Transaction struct {
  Action                  string        `json:"action,omitempty"`
  Actor                   string        `json:"actor,omitempty"`
  AllowedCancel           bool          `json:"allowed_cancel,bool"`
  AssociatedObjectID      int           `json:"associated_object_id,omitempty"`
  AssociatedObjectType    string        `json:"associated_object_type,omitempty"`
  ChainID                 int           `json:"chain_id,omitempty"`
  CreatedAt               string        `json:"created_at,omitempty"`
  DependentTransactionID  int           `json:"dependent_transaction_id,omitempty"`
  ID                      int           `json:"id,omitempty"`
  Identifier              string        `json:"identifier,omitempty"`
  LockVersion             int           `json:"lock_version,omitempty"`
  ParentID                int           `json:"parent_id,omitempty"`
  ParentType              string        `json:"parent_type,omitempty"`
  Pid                     int           `json:"pid,omitempty"`
  Priority                int           `json:"priority,omitempty"`
  Scheduled               bool          `json:"scheduled,bool"`
  StartAfter              string        `json:"start_after,omitempty"`
  StartedAt               string        `json:"started_at,omitempty"`
  Status                  string        `json:"status,omitempty"`
  UpdatedAt               string        `json:"updated_at,omitempty"`
  UserID                  int           `json:"user_id,omitempty"`

  Params                  map[string]interface{}  `json:"params,omitempty"`
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

  trx := make([]Transaction, len(out))
  for i := range trx {
    trx[i] = out[i]["transaction"]
  }

  return trx, resp, err
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

// ListByGroup return group of transcations depended by action
// func (s *TransactionsServiceOp) ListByGroup(ctx context.Context, id int, objectType string, opt *ListOptions) (*list.List, *Response, error) {
func (s *TransactionsServiceOp) ListByGroup(ctx context.Context, meta interface{}, opt *ListOptions) (*list.List, *Response, error) {
  var associatedObjectID, parentID int
  var associatedObjectType, parentType string

  val := reflect.ValueOf(meta)

  v1 := val.FieldByName("AssociatedObjectID")
  if v1.IsValid() {
    associatedObjectID = v1.Interface().(int)
    // fmt.Printf("associatedObjectID: <%d>\n", associatedObjectID)

    if associatedObjectID < 1 {
      return nil, nil, godo.NewArgError("id", "cannot be less than 1")
    }
  }

  v2 := val.FieldByName("ParentID")
  if v2.IsValid() {
    parentID = v2.Interface().(int)
    // fmt.Printf("parentID: <%d>\n", parentID)

    if parentID < 1 {
      return nil, nil, godo.NewArgError("id", "cannot be less than 1")
    }
  }
  
  v1 = val.FieldByName("AssociatedObjectType")
  if v1.IsValid() {
    associatedObjectType = v1.String()
  }

  v2 = val.FieldByName("ParentType")
  if v2.IsValid() {
    parentType = v2.String()
  }

  trx, resp, err := s.client.Transactions.List(ctx, opt)
  if err != nil {
    return nil, resp, fmt.Errorf("ListByGroup.trx: %s", err)
  }

  var next *Transaction

  len := len(trx)
  groupList := list.New()

  for i := range trx {
    cur := trx[i]
    if associatedObjectType != "" {
      // fmt.Printf("cur.AssociatedObjectID: <%d> -> associatedObjectID: <%d>\n", cur.AssociatedObjectID, associatedObjectID)
      // fmt.Printf("cur.AssociatedObjectType: <%s> -> associatedObjectType: <%s>\n", cur.AssociatedObjectType, associatedObjectType)
      if cur.AssociatedObjectID != associatedObjectID || cur.AssociatedObjectType != associatedObjectType {
        continue
      }
    }
    
    if parentType != "" {
      // fmt.Printf("cur.ParentID: <%d> -> parentID: <%d>\n", cur.ParentID, parentID)
      // fmt.Printf("cur.ParentType: <%s> -> parentType: <%s>\n", cur.ParentType, parentType)
      if cur.ParentID != parentID || cur.ParentType != parentType {
        continue
      }
    }

    if cur.DependentTransactionID == 0 {
      groupList.PushBack(cur)
      break
    }

    if i+1 < len {
      next = &trx[i+1]
    }

    if next != nil {
      if associatedObjectType != "" {
        if cur.AssociatedObjectID == next.AssociatedObjectID &&
          cur.AssociatedObjectType == next.AssociatedObjectType &&
          cur.ChainID == next.ChainID {
          groupList.PushBack(cur)
        }
      }

      if parentType != "" {
        if cur.ParentID == next.ParentID &&
          cur.ParentType == next.ParentType &&
          cur.ChainID == next.ChainID {
          groupList.PushBack(cur)
        }
      }
    }
  }

  return groupList, resp, err
}

// GetByFilter find transaction with specified fields for virtual machine by ID.
func (s *TransactionsServiceOp) GetByFilter(ctx context.Context, id int, filter interface{}, opt *ListOptions) (*Transaction, *Response, error) {
  if id < 1 {
    return nil, nil, godo.NewArgError("id", "cannot be less than 1")
  }

  trx, resp, err := s.client.Transactions.ListByGroup(ctx, filter, opt)
  if err != nil {
    return nil, resp, fmt.Errorf("GetByFilter.trx: %s", err)
  }

  var root *Transaction
  for e := trx.Front(); e != nil; e = e.Next() {
    val := e.Value.(Transaction)
    root = &val
    if root.equal(filter) {
      break
    } else {
      root = nil
    }
  }

  if root != nil {
    return root, resp, err
  }

  return nil, nil, fmt.Errorf("Transaction not found or wrong filter %+v", filter)
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

// func lastTransaction(ctx context.Context, client *Client, id int, objectType string) (*Transaction, *Response, error) {
func lastTransaction(ctx context.Context, client *Client, filter interface{}) (*Transaction, *Response, error) {
  opt := &ListOptions{
    PerPage : searchTransactions,
  }

  // trx, resp, err := client.Transactions.ListByGroup(ctx, id, aot, opt)
  trx, resp, err := client.Transactions.ListByGroup(ctx, filter, opt)

  var root *Transaction
  e := trx.Front()
  if e != nil {
    val := e.Value.(Transaction)
    root = &val
    return root, resp, err
  }

  return nil, nil, err
}

func (trx Transaction) String() string {
  return godo.Stringify(trx)
}

// Running check if transaction state is 'runing'
func (trx Transaction) Running() bool {
  return trx.Status == TransactionRunning
}

// Pending check if transaction state is 'pending'
func (trx Transaction) Pending() bool {
  return trx.Status == TransactionPending
}

// Incomplete check if transaction state is 'running' or 'pending'
func (trx Transaction) Incomplete() bool {
  return trx.Running() || trx.Pending()
}

// Complete check if transaction state is 'complete'
func (trx Transaction) Complete() bool {
  return trx.Status == TransactionComplete
}

// Failed check if transaction state is 'failed'
func (trx Transaction) Failed() bool {
  return trx.Status == TransactionFailed
}

// Cancelled check if transaction state is 'cancelled'
func (trx Transaction) Cancelled() bool {
  return trx.Status == TransactionCancelled
}

// Unlucky check if transaction state is 'failed' or 'cancelled'
func (trx Transaction) Unlucky() bool {
  return trx.Failed() || trx.Cancelled()
}

// Finished check if transaction state is
// 'complete' or 'failed' or 'cancelled'
func (trx Transaction) Finished() bool {
  return trx.Complete() || trx.Unlucky()
}
