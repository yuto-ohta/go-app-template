package appmodel

import (
	"errors"
	"go-app-template/src/apperror"
	"net/http"
)

type SearchCondition struct {
	orderBy string
	order   Order
	page    int
	limit   int
}

/**************************************
	Constructor
**************************************/
func NewSearchCondition(orderBy string, order Order, page int, limit int) (*SearchCondition, error) {
	condition := &SearchCondition{
		orderBy: orderBy,
		order:   order,
		page:    page,
		limit:   limit,
	}
	if err := validate(*condition); err != nil {
		return nil, apperror.NewAppErrorWithStatus(err, http.StatusBadRequest)
	}

	return &SearchCondition{
		orderBy: orderBy,
		order:   order,
		page:    page,
		limit:   limit,
	}, nil
}

/**************************************
	Getter & Setter
**************************************/
func (s SearchCondition) GetOrderBy() string {
	return s.orderBy
}

func (s SearchCondition) GetOrder() Order {
	return s.order
}

func (s SearchCondition) GetPage() int {
	return s.page
}

func (s SearchCondition) GetLimit() int {
	return s.limit
}

/**************************************
	Validation
**************************************/
func validate(condition SearchCondition) error {
	// limitがない && pageが指定されているとき
	if condition.page > 0 && condition.limit <= 0 {
		return errors.New("pageNumを指定する場合は、limitも指定してください")
	}
	return nil
}
