package onappgo

// Links manages links that are returned along with a List
type Links struct {
  PerPage   int
  CurPage   int
  Total     int
  NumPages  int
}

// CurrentPage is current page of the list
func (l *Links) CurrentPage() (int, error) {
  return l.CurPage, nil
}

// IsLastPage returns true if the current page is the last
func (l *Links) IsLastPage() bool {
  return l.CurPage == l.NumPages
}
