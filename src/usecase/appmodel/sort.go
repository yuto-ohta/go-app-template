package appmodel

type Order int

const (
	ASC Order = iota + 1
	DESC
)

func (o Order) String() string {
	var messages = map[Order]string{
		ASC:  "ASC",
		DESC: "DESC",
	}
	return messages[o]
}

/**************************************
	Getter & Setter
**************************************/

func GetOrder(param string) Order {
	switch param {
	case ASC.String():
		return ASC
	case DESC.String():
		return DESC
	default:
		return -1
	}
}
