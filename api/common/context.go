package common

import (
	"errors"
	"github.com/labstack/echo"
	"net/http"
	"s7ab-platform-hyperledger/platform/core"
	"s7ab-platform-hyperledger/platform/core/logger"
	"s7ab-platform-hyperledger/platform/core/observer"
	"strconv"
)

const (
	defaultPaginationLimit = 20
	limitQueryParam        = `limit`
	offsetQueryParam       = `offset`
)

var (
	ErrInvalidContext = errors.New(`invalid context`)
)

type Context struct {
	echo.Context
	Log      logger.Logger
	SDK      *core.SDKCore
	Observer *observer.Observer
}

// WriteError
// Format error and logging it
func (c *Context) WriteError(err error) error {
	c.Log.Debug(c.Path(), logger.KV(`error`, err.Error()))
	return c.JSON(http.StatusInternalServerError, map[string]interface{}{`success`: false, `error`: err.Error()})
}

// WriteSuccess
// Format response and logging it
func (c *Context) WriteSuccess(val interface{}) error {
	c.Log.Debug(c.Path(), logger.KV(`data`, val))
	return c.JSON(http.StatusOK, map[string]interface{}{`result`: val, `success`: true, `error`: ``})
}

func (c *Context) WriteClearSuccess(val interface{}) error {
	c.Log.Debug(c.Path(), logger.KV(`data`, val))
	return c.JSON(http.StatusOK, val)
}

// GetPagination
// Get pagination params from query string or use default
func (c *Context) GetPagination() (int, int) {
	var (
		limit  int
		offset int
		err    error
	)
	qp := c.QueryParams()
	if limitStr := qp.Get(limitQueryParam); limitStr != `` {
		if limit, err = strconv.Atoi(limitStr); err != nil {
			limit = defaultPaginationLimit
		}
	} else {
		limit = defaultPaginationLimit
	}

	if offsetStr := qp.Get(offsetQueryParam); offsetStr != `` {
		offset, _ = strconv.Atoi(offsetStr)
	}
	return limit, offset
}

func NewContext(e echo.Context, s *core.SDKCore, l logger.Logger) Context {
	c := Context{}
	c.SDK = s
	c.Context = e
	c.Observer = observer.NewObserver(s)
	c.Log = l
	return c
}
