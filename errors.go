package saasquatch

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ErrNoSuchCode        = errors.New("referral: no such code")
	ErrInvalidRedemption = errors.New("reward: invalid redemption due to zero amount or exceed redeemable credit")
	ErrBadRequest        = errors.New("saasquatch: bad request")
)

type ErrorResponse struct {
	Response     *http.Response // HTTP response that caused this error
	StatusCode   int            `json:"statusCode"`
	Message      string         `json:"message"`
	ApiErrorCode string         `json:"apiErrorCode"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %+v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Message, r.ApiErrorCode)
}

func (r ErrorResponse) transalateError() error {
	switch r.ApiErrorCode {
	case "BAD_REQUEST":
		return ErrBadRequest
	case "REFERRAL_CODE_NOT_FOUND":
		return ErrNoSuchCode
	case "INVALD_REWARD_REDEMPTION":
		return ErrInvalidRedemption
	default:
		return nil
	}

	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)

		if err := errorResponse.transalateError(); err != nil {
			return err
		}
	}
	return errorResponse
}
