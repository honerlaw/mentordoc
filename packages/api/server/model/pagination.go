package model

import (
	"net/http"
	"strconv"
)

const defaultPageCount = "25"

type Pagination struct {
	Page  int
	Count int
}

func NewPagination(req *http.Request) *Pagination {
	page := req.URL.Query().Get("page")
	count := req.URL.Query().Get("count")

	// no page don't return anything
	if len(page) == 0 {
		return nil
	}

	// default to 25
	if len(count) == 0 {
		count = defaultPageCount
	}

	pageNum, err := strconv.Atoi(page)
	if err != nil {
		return nil
	}
	countNum, err := strconv.Atoi(count)
	if err != nil {
		return nil
	}

	return &Pagination{
		Page:  pageNum,
		Count: countNum,
	}
}
