package httputils

import "net/http"

type Controller interface {

	// =======================
	// CHAIN OF RESPONSIBILITY
	// =======================
	http.Handler
	SetNext(next Controller) Controller

	// =======================
	// AUXILIARIES
	// =======================
	// Global variables of the chain
	SetErrorPrefix(errorPrefix string) Controller
}

type BaseController struct {
	ErrPrefix string
	Next      Controller
}

func (c *BaseController) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	// call next
	if c.Next != nil {
		c.Next.ServeHTTP(w, r)
	}
}

func (c *BaseController) SetErrorPrefix(
	errPrefix string,
) Controller {
	c.ErrPrefix = errPrefix
	if c.Next != nil {
		c.Next.SetErrorPrefix(errPrefix)
	}
	return c
}

func (c *BaseController) SetNext(
	next Controller,
) Controller {
	c.Next = next
	return c
}
