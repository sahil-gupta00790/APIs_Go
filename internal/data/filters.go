package data

import "Apis_go.sahil.net/internal/validator"

type Filters struct {
	Page     int
	PageSize int
	Sort     string
	SortList []string
}

func ValidateFilters(v *validator.Validator, f Filters) {
	//check pafe and page size parameters
	v.Check(f.Page > 0, "page", "must b greater than 0")
	v.Check(f.Page <= 1000, "page", "must be a maximum of 1000")
	v.Check(f.PageSize > 0, "page_size", "must b greater than 0")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")
	//checkl for sort parameter matches a value in acceptable sort list
	v.Check(validator.IN(f.Sort, f.SortList...), "sort", "invalid sort value")

}
