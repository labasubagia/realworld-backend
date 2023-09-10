package exception

const (
	TypeInternal         = "ErrInternal"
	TypeValidation       = "ErrValidation"
	TypeNotFound         = "ErrNotFound"
	TypePermissionDenied = "ErrPermissionDenied"
	TypeTokenExpired     = "TokenExpired"
	TypeTokenInvalid     = "TokenInvalid"
)

type Err = map[string][]string

type Exception struct {
	Type    string `json:"-"`
	Message string `json:"-"`
	Cause   error  `json:"-"`
	Errors  Err    `json:"errors"`
}

func New(kind, message string, err error) *Exception {
	return &Exception{
		Type:    kind,
		Cause:   err,
		Message: message,
		Errors:  make(map[string][]string),
	}
}

func Validation() *Exception {
	return New(TypeValidation, "validation error", nil)
}

func (e *Exception) HasError() bool {
	return len(e.Errors) > 0
}

func (e *Exception) AddError(key, msg string) *Exception {
	e.Errors[key] = append(e.Errors[key], msg)
	return e
}

func Into(err error) *Exception {
	if err == nil {
		return nil
	}
	fail, ok := err.(*Exception)
	if ok {
		return fail
	}
	return New(TypeInternal, err.Error(), err)
}

func (fail *Exception) Error() string {
	return fail.Cause.Error()
}
