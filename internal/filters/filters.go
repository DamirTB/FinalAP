package filters

import (
	"strings" // New import
	"damir/internal/validator" // New import
)

type Filters struct {
	Page int
	PageSize int
	Sort string
	SortSafelist []string
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
	if f.Sort == safeValue {
	return strings.TrimPrefix(f.Sort, "-")
	}
	}
	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
	return "DESC"
	}
	return "ASC"
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// Check that the page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}