package pagination

import (
	"errors"
	"go-app-template/src/apperror"
)

type Pageable []interface{}

/**************************************
	Getter & Setter
**************************************/

func (p Pageable) GetPage(pageNum int, limit int) (*Page, error) {
	var err error

	var pageInfo *PageInfo
	if pageInfo, err = newPageInfo(pageNum, limit, len(p)); err != nil {
		return nil, apperror.NewAppError(err)
	}

	var page *Page
	if page, err = newPage(*pageInfo, p); err != nil {
		return nil, apperror.NewAppError(err)
	}

	return page, nil
}

/**************************************
	private
**************************************/

func (p Pageable) skip(offset int) (Pageable, error) {
	if offset < 0 {
		return Pageable{}, errors.New("offsetは0以上にしてください")
	}
	if len(p) <= offset {
		return Pageable{}, nil
	}

	list := make(Pageable, len(p)-offset)
	for i, e := range p {
		if i < offset {
			continue
		}
		list[i-offset] = e
	}
	return list, nil
}

func (p Pageable) limit(limit int) (Pageable, error) {
	if limit < 0 {
		return Pageable{}, errors.New("limitは0以上にしてください")
	}
	var list Pageable
	if len(p) < limit {
		list = make(Pageable, len(p))
	} else {
		list = make(Pageable, limit)
	}
	for i, e := range p {
		if i == limit {
			break
		}
		list[i] = e
	}
	return list, nil
}
