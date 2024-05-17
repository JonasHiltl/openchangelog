package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

// DecodeJSON will decode the response body to a JSON structure. This
// will consume the response body, but will not close it. Close must
// still be called.
func (r *Response) DecodeJSON(out interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	return dec.Decode(out)
}

// Error returns an error response if there is one. If there is an error,
// this will fully consume the response body, but will not close it. The
// body must still be closed manually.
func (r *Response) Error() error {
	// 200 to 399 are okay status codes.
	if r.StatusCode >= 200 && r.StatusCode < 400 {
		return nil
	}

	// We have an error. Let's copy the body into our own buffer first,
	// so that if we can't decode JSON, we can at least copy it raw.
	bodyBuf := &bytes.Buffer{}
	if _, err := io.Copy(bodyBuf, r.Body); err != nil {
		return err
	}

	r.Body.Close()
	r.Body = io.NopCloser(bodyBuf)

	// Build up the error object
	respErr := &ApiError{
		HTTPMethod: r.Request.Method,
		URL:        r.Request.URL.String(),
		StatusCode: r.StatusCode,
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	dec := json.NewDecoder(bytes.NewReader(bodyBuf.Bytes()))
	dec.UseNumber()
	if err := dec.Decode(&resp); err != nil {
		// failed to decode ApiError, just return body as string
		resp.Message = bodyBuf.String()
	} else {
		respErr.Message = resp.Message
	}

	return respErr
}

type ApiError struct {
	// HTTPMethod is the HTTP method for the request (PUT, GET, etc).
	HTTPMethod string

	// URL is the URL of the request.
	URL string

	// StatusCode is the HTTP status code.
	StatusCode int

	Message string
}

func (a ApiError) Error() string {
	return a.Message
}
