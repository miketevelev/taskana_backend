package core_http_request

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	core_errors "github.com/miketevelev/taskana_backend/internal/core/errors"
)

func GetUserIdLimitOffsetQueryParams(r *http.Request) (
	*int,
	*int,
	error,
) {
	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to parse 'limit' query parameter: %w",
			err,
		)
	}

	offset, err := GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to parse 'offset' query parameter: %w",
			err,
		)
	}

	return limit, offset, nil
}

func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := strconv.Atoi(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param '%s' by key '%s' not a valid integer: %v: %w",
			param, key, err, core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}

func GetStringQueryParam(r *http.Request, key string) (*string, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}
	return &param, nil
}

func GetDateQueryParam(r *http.Request, key string) (*time.Time, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	layout := "2006-01-02"
	date, err := time.Parse(layout, param)
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid date: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &date, nil
}

func GetUUIDQueryParam(r *http.Request, key string) (*uuid.UUID, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return nil, nil
	}

	val, err := uuid.Parse(param)
	if err != nil {
		return nil, fmt.Errorf(
			"param '%s' by key='%s' is not a valid UUID: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return &val, nil
}
