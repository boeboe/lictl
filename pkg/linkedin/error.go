package linkedin

import "net/http"

// HTTPError represents an HTTP error with a status code.
type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return http.StatusText(e.StatusCode)
}
