package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"sort" validate:"gte=1,lte=20"`
	Offset int      `json:"sort" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	query := r.URL.Query()
	limit := query.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, err
		}
		fq.Limit = l
	}
	offset := query.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return fq, err
		}
		fq.Offset = o
	}
	sort := query.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}
	tags := query.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}
	search := query.Get("search")
	if search != "" {
		fq.Search = search
	}
	return fq, nil
}
