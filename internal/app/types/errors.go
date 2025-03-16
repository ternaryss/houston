package types

type FieldError struct {
	Field string
	Msg   string
}

func NewFieldError(field, msg string) FieldError {
	return FieldError{
		Field: field,
		Msg:   msg,
	}
}
