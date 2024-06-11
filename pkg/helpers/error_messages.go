package helpers

const (
	ErrorDefaultIncomplete = "Something gone wrong, please contact your system administrator"
	ErrorDefaultBadRequest = "Bad Request"
	ErrorDefaultNotFound   = "Not Found"
)

var ErrorDefaultMessage = map[int64]string{
	404: ErrorDefaultNotFound,
	400: ErrorDefaultBadRequest,
}
