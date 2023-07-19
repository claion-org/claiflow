package echov4

import "fmt"

var (
	ErrorInvalidRequestParameter = fmt.Errorf("invalid request parameter")
	// ErrorBindRequestObject       = fmt.Errorf("could not bind request object")
	ErrorInvalidRequestPath = fmt.Errorf("invalid request path")
)
