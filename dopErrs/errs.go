package dopErrs

// Err

type Err string

func (e Err) Error() string {
	return string(e)
}

// ErrWithDesc

type ErrWithDesc struct {
	Err  error
	Desc string
}

func (e ErrWithDesc) Error() string {
	return e.Err.Error() + ", desc:" + e.Desc
}

// FormErr

type FormErr struct {
	Fields map[string]error
}

func (e FormErr) Error() string {
	return FormValidate.Error()
}

// errors

const (
	NoRows            = Err("err_no_rows")
	BadColumnName     = Err("bad_column_name")
	BadJson           = Err("bad_json")
	BadJwt            = Err("bad_jwt")
	BadQueryParams    = Err("bad_query_params")
	ServiceNA         = Err("service_not_available")
	NotImplemented    = Err("not_implemented")
	NotAuthorized     = Err("not_authorized")
	PermissionDenied  = Err("permission_denied")
	ObjectNotFound    = Err("object_not_found")
	IncorrectPageSize = Err("incorrect_page_size")
	BadStatusCode     = Err("bad_status_code")
	FormValidate      = Err("form_validate")
)
