package data

import (
	"slices"
	"strings"

	"github.com/andreshungbz/cmps4191-lab-01-synchronous-api/internal/validator"
)

// Filters holds information on desired pagination, sort, and sort safe list settings.
type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string // the valid string values for a sort column
}

// sortColumn parses the column to be sorted.
func (f Filters) sortColumn() string {
	if slices.Contains(f.SortSafelist, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	panic("Unsafe sort parameter: " + f.Sort)
}

// sortDirection determines which SQL sorting keyword to use based on the dash character.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

// limit returns the Filters PageSize.
func (f Filters) limit() int {
	return f.PageSize
}

// offset returns the SQL offset value required to show the correct page.
func (f Filters) offset() int {
	// we decrement by 1 since the 1st page requires no offset value
	return (f.Page - 1) * f.PageSize
}

// ValidateFilters ensures filter values are reasonable.
func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "Must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "Must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "Must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "Must be a maximum of 100")
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "Invalid sort value")
}

// Metadata holds pagination information for a query response of multiple records.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitzero"`
	PageSize     int `json:"page_size,omitzero"`
	FirstPage    int `json:"first_page,omitzero"`
	LastPage     int `json:"last_page,omitzero"`
	TotalRecords int `json:"total_records,omitzero"`
}

// calculateMetadata calculates page information from the object and total records.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
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
