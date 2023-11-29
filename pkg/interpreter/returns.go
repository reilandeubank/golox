package interpreter

type ReturnError struct {
	Value interface{}
}

func (r ReturnError) Error() interface{} {
	return r.Value
}