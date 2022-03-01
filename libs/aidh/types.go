package aidh

type FormFieldType int

const (
	FormFieldText FormFieldType = 1
	FormFieldFile FormFieldType = 2
)

//Error :
type Error struct {
	Message string `json:"msg"`
}

//Error :
func (e *Error) Error() string {
	return e.Message
}

//NewError :
func NewError(msg string) error {
	return &Error{Message: msg}
}

type FormField struct {
	Type  FormFieldType
	Key   string
	Value string
}

type FormFields []FormField
