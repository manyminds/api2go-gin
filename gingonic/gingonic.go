package gingonic

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/manyminds/api2go/routing"
)

type ginRouter struct {
	router    *gin.Engine
	whitelist []string
}

func (g ginRouter) Handler() http.Handler {
	return g.router
}

func (g ginRouter) Handle(protocol, route string, handler routing.HandlerFunc) {
	if g.whiteListed(route) {
		wrappedCallback := func(c *gin.Context) {
			params := map[string]string{}
			for _, p := range c.Params {
				params[p.Key] = p.Value
			}

			handler(c.Writer, c.Request, params)
		}

		g.router.Handle(protocol, route, wrappedCallback)
	}
}

func (g ginRouter) whiteListed(route string) bool {
	if len(g.whitelist) == 0 {
		// if the white list is omitted all routes are valid
		return true
	}

	// scan through the white list for allowed route prefixes
	for _, allowed := range g.whitelist {
		if strings.HasPrefix(route, allowed) {
			return true
		}
	}

	// not found in white list
	return false
}

//New creates a new api2go router to use with the gin framework
func New(g *gin.Engine, w ...string) routing.Routeable {
	return &ginRouter{router: g, whitelist: w}
}
