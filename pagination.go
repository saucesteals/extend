package extend

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Paginator[T any, R PaginatedResponse[T]] struct {
	hasNext  bool
	nextPage int
	api      *Client
	path     string
	query    url.Values
}

func newPaginator[T any, R PaginatedResponse[T]](api *Client, options PaginationOptions, path string, query url.Values) *Paginator[T, R] {
	query.Set("count", strconv.Itoa(options.Count))
	query.Set("sortDirection", string(options.SortDirection))
	query.Set("sortField", options.SortField)
	return &Paginator[T, R]{true, options.Page, api, path, query}
}

func (p *Paginator[T, R]) Next() bool {
	return p.hasNext
}

func (p *Paginator[T, R]) Get(ctx context.Context) (*R, error) {
	if !p.hasNext {
		return nil, errors.New("no more items")
	}

	p.query.Set("page", strconv.Itoa(p.nextPage))
	p.nextPage++

	var response R
	err := p.api.jsonRequest(ctx, http.MethodGet, p.path+"?"+p.query.Encode(), nil, &response)
	if err != nil {
		return nil, err
	}

	pagination := response.Pagination()
	p.hasNext = pagination.Page < pagination.NumberOfPages

	return &response, nil
}

type SortDirection string

const (
	SortDirectionAsc  SortDirection = "ASC"
	SortDirectionDesc SortDirection = "DESC"
)

type PaginationOptions struct {
	// Page is zero indexed
	Page int

	// Count is the number of items per page
	Count int

	// SortDirection is the direction of the sort
	SortDirection SortDirection

	// SortField is the field to sort by
	SortField string
}

type Pagination struct {
	// Page is the current page
	Page int `json:"page"`

	// PageItemCount is the number of items on the current page
	PageItemCount int `json:"pageItemCount"`

	// TotalItems is the total number of items
	TotalItems int `json:"totalItems"`

	// NumberOfPages is the total number of pages
	NumberOfPages int `json:"numberOfPages"`
}

type PaginationResponse struct {
	PaginationData Pagination `json:"pagination"`
}

func (p PaginationResponse) Pagination() Pagination {
	return p.PaginationData
}

type PaginatedResponse[T any] interface {
	Items() []T
	Pagination() Pagination
}
