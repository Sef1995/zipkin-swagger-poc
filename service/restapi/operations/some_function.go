// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// SomeFunctionHandlerFunc turns a function with the right signature into a some function handler
type SomeFunctionHandlerFunc func(SomeFunctionParams) middleware.Responder

// Handle executing the request and returning a response
func (fn SomeFunctionHandlerFunc) Handle(params SomeFunctionParams) middleware.Responder {
	return fn(params)
}

// SomeFunctionHandler interface for that can handle valid some function params
type SomeFunctionHandler interface {
	Handle(SomeFunctionParams) middleware.Responder
}

// NewSomeFunction creates a new http.Handler for the some function operation
func NewSomeFunction(ctx *middleware.Context, handler SomeFunctionHandler) *SomeFunction {
	return &SomeFunction{Context: ctx, Handler: handler}
}

/*SomeFunction swagger:route GET /somefunction someFunction

someFunction

*/
type SomeFunction struct {
	Context *middleware.Context
	Handler SomeFunctionHandler
}

func (o *SomeFunction) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewSomeFunctionParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
