package dto

type PageInfo struct {
	PageNum     int `json:"pageNumber"`
	LastPageNum int `json:"lastPageNumber"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
}
