package endpoints

// APIRoute is a basic struct containing Method and URL
type APIRoute struct {
	Route
	method Method
}

// Method returns the request method used by the route
func (r APIRoute) Method() Method {
	return r.method
}

// NewAPIRoute generates a new discord api route struct
func NewAPIRoute(method Method, url string) APIRoute {
	return APIRoute{
		Route: Route{
			baseRoute:  API,
			route:      url,
			paramCount: countParams(url),
		},
		method: method,
	}
}