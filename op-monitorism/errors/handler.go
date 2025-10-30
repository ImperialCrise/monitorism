package errors

import (
	"errors"
	
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
)

type ErrorHandler struct {
	log     log.Logger
	metrics *prometheus.CounterVec
}

func NewErrorHandler(log log.Logger, metrics *prometheus.CounterVec) *ErrorHandler {
	return &ErrorHandler{
		log:     log,
		metrics: metrics,
	}
}

func (h *ErrorHandler) HandleError(err error, labels ...string) {
	if err == nil {
		return
	}

	code := GetCode(err)
	h.logError(err, code)
	
	if h.metrics != nil && len(labels) > 0 {
		h.metrics.WithLabelValues(labels...).Inc()
	}
}

func (h *ErrorHandler) logError(err error, code ErrorCode) {
	var me *MonitorError
	if Is(err, code) {
		errors.As(err, &me)
	}

	switch code {
	case ErrCodeNetwork, ErrCodeTimeout:
		if me != nil && len(me.Details) > 0 {
			h.log.Warn("network error", "error", err.Error(), "details", me.Details)
		} else {
			h.log.Warn("network error", "error", err.Error())
		}
	case ErrCodeValidation:
		if me != nil && len(me.Details) > 0 {
			h.log.Error("validation error", "error", err.Error(), "details", me.Details)
		} else {
			h.log.Error("validation error", "error", err.Error())
		}
	default:
		if me != nil && len(me.Details) > 0 {
			h.log.Error("error occurred", "error", err.Error(), "code", code, "details", me.Details)
		} else {
			h.log.Error("error occurred", "error", err.Error(), "code", code)
		}
	}
}

func (h *ErrorHandler) HandleErrorWithContext(err error, context string, labels ...string) {
	if err == nil {
		return
	}

	code := GetCode(err)
	h.logErrorWithContext(err, code, context)

	if h.metrics != nil && len(labels) > 0 {
		h.metrics.WithLabelValues(labels...).Inc()
	}
}

func (h *ErrorHandler) logErrorWithContext(err error, code ErrorCode, context string) {
	var me *MonitorError
	if Is(err, code) {
		errors.As(err, &me)
	}

	logFields := []interface{}{"error", err.Error(), "context", context}
	if me != nil && len(me.Details) > 0 {
		logFields = append(logFields, "details", me.Details)
	}

	switch code {
	case ErrCodeNetwork, ErrCodeTimeout:
		h.log.Warn("network error", logFields...)
	case ErrCodeValidation:
		h.log.Error("validation error", logFields...)
	default:
		logFields = append(logFields, "code", code)
		h.log.Error("error occurred", logFields...)
	}
}
