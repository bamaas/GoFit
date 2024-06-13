package data

import (
	"time"
)

type Filters struct {
	Page int
	PageSize int
	StartTime time.Time
	EndTime time.Time
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	TotalRecords int `json:"total_records,omitempty"`
	CurrentPage int `json:"current_page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
	FirstPage int `json:"first_page,omitempty"`
	LastPage int `json:"last_page,omitempty"`
}

func calculateMetadata(totalRecords, currentPage, pageSize int) Metadata {

	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		TotalRecords: totalRecords,
		CurrentPage: currentPage,
		PageSize: pageSize,
		FirstPage: 1,
		LastPage: (totalRecords + pageSize - 1) / pageSize,
	}
}