package apispec

// Service stores information about an API.
type Service struct {
	Domain string
	Host   string
	Routes []Route
}

// Route stores information about an API Route.
type Route struct {
	Method      string        // HTTP Method [GET, POST, PATCH, DELETE]
	Path        string        // route path
	Headers     []*Header     // HTTP headers
	URLParams   []*URLParam   // URL parameters withing the route path
	QueryParams []*QueryParam // HTTP query parameters
	Payload     []*Field      // HTTP request body fields
}

func (r Route) String() string {
	return r.Method + " " + r.Path
}

// Header stores information about HTTP Headers.
type Header struct {
	Name string
}

// URLParam stores information about URL params.
type URLParam struct {
	Name      string
	Type      BasicType
	MaxLength int
}

// QueryParam stores information about query parameters.
type QueryParam struct {
	Name     string
	Type     BasicType
	Optional bool
}
