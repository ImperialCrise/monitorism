package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

func WriteHTTPError(w http.ResponseWriter, err error) {
	var resp ErrorResponse
	
	code := GetCode(err)
	status := getHTTPStatus(code)
	
	var me *MonitorError
	if errors.As(err, &me) {
		resp = ErrorResponse{
			Code:    me.Code,
			Message: me.Message,
			Details: me.Details,
		}
	} else {
		resp = ErrorResponse{
			Code:    code,
			Message: err.Error(),
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func getHTTPStatus(code ErrorCode) int {
	switch code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodePermission:
		return http.StatusForbidden
	case ErrCodeTimeout:
		return http.StatusRequestTimeout
	case ErrCodeNetwork:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
