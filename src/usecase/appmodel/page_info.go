package appmodel

import (
	"errors"
	"go-app-template/src/apperror"
)

type PageInfo struct {
	pageNum     int
	lastPageNum int
	limit       int
	offset      int
}

/**************************************
	Getter & Setter
**************************************/
func (p PageInfo) GetLastPageNum() int {
	return p.lastPageNum
}

func (p PageInfo) GetPageNum() int {
	return p.pageNum
}

func (p PageInfo) GetLimit() int {
	return p.limit
}

func (p PageInfo) GetOffset() int {
	return p.offset
}

/**************************************
	Constructor
**************************************/
func newPageInfo(pageNum int, limit int, maxCount int) (*PageInfo, error) {
	if maxCount <= 0 {
		return &PageInfo{}, apperror.NewAppError(errors.New("maxCountは1以上にしてください"))
	}
	// pageNumのみ指定の場合は、400とする
	if pageNum > 0 && limit <= 0 {
		return &PageInfo{}, apperror.NewAppError(errors.New("pageNumを指定する場合は、limitも指定してください"))
	}
	// limitのみ指定の場合は、page=1として扱う
	if pageNum <= 0 && limit > 0 {
		pageNum = 1
	}
	// pageNum, limitが未指定の場合は、全件ページとして扱う
	if pageNum <= 0 && limit <= 0 {
		return &PageInfo{
			pageNum:     1,
			lastPageNum: 1,
			limit:       maxCount,
			offset:      0,
		}, nil
	}

	// 各値の算出
	lastPageNum := maxCount / limit
	if maxCount%limit != 0 {
		lastPageNum++
	}
	if pageNum > lastPageNum {
		// 指定のページが最終ページを超過している場合は、最終ページとして扱う
		pageNum = lastPageNum
	}
	offset := (pageNum - 1) * limit

	return &PageInfo{
		pageNum:     pageNum,
		lastPageNum: lastPageNum,
		limit:       limit,
		offset:      offset,
	}, nil
}
