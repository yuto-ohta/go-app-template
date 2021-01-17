package appmodel

import "go-app-template/src/apperror"

type Page struct {
	info PageInfo
	list []interface{}
}

/**************************************
	Getter & Setter
**************************************/
func (p Page) GetInfo() PageInfo {
	return p.info
}

func (p Page) GetList() []interface{} {
	return p.list
}

/**************************************
	Constructor
**************************************/
func newPage(info PageInfo, target Pageable) (*Page, error) {
	var err error

	var skipped Pageable
	if skipped, err = target.skip(info.offset); err != nil {
		return nil, apperror.NewAppError(err)
	}
	var limited Pageable
	if limited, err = skipped.limit(info.limit); err != nil {
		return nil, apperror.NewAppError(err)
	}

	return &Page{
		info: info,
		list: limited,
	}, nil
}
