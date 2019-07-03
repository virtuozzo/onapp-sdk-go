package onappgo

import (
  "net/url"
  "strconv"
)

// Links manages links that are returned along with a List
type Links struct {
  Pages   *Pages       `json:"pages,omitempty"`
}

// Pages are pages specified in Links
type Pages struct {
  First string `json:"first,omitempty"`
  Prev  string `json:"prev,omitempty"`
  Last  string `json:"last,omitempty"`
  Next  string `json:"next,omitempty"`
}

// CurrentPage is current page of the list
func (l *Links) CurrentPage() (int, error) {
  return l.Pages.current()
}

func (p *Pages) current() (int, error) {
  switch {
  case p == nil:
    return 1, nil
  case p.Prev == "" && p.Next != "":
    return 1, nil
  case p.Prev != "":
    prevPage, err := pageForURL(p.Prev)
    if err != nil {
      return 0, err
    }

    return prevPage + 1, nil
  }

  return 0, nil
}

// IsLastPage returns true if the current page is the last
func (l *Links) IsLastPage() bool {
  if l.Pages == nil {
    return true
  }
  return l.Pages.isLast()
}

func (p *Pages) isLast() bool {
  return p.Last == ""
}

func pageForURL(urlText string) (int, error) {
  u, err := url.ParseRequestURI(urlText)
  if err != nil {
    return 0, err
  }

  pageStr := u.Query().Get("page")
  page, err := strconv.Atoi(pageStr)
  if err != nil {
    return 0, err
  }

  return page, nil
}
