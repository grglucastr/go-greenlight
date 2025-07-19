package data

import (
	"strings"

	"github.com/grglucastr/go-greenlight/internal/validator"
)

type Metadata struct {
	CurrentPage  int `json:"current_page,omitzero"`
	PageSize     int `json:"page_size,omitzero"`
	FirstPage    int `json:"first_page,omitzero"`
	LastPage     int `json:"last_page,omitzero"`
	TotalRecords int `json:"total_records,omitzero"`
}

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	if pageSize == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     (totalRecords + pageSize - 1) / pageSize,
		TotalRecords: totalRecords,
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}

// Check that the client-provided Sort field matches one of the entries in our safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).

func (f Filters) sortColumn() string {

	for _, saveValue := range f.SortSafelist {
		if f.Sort == saveValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) limit() int {
	return f.PageSize
}

func (f Filters) offset() int {
	return (f.Page - 1) * f.PageSize
}
