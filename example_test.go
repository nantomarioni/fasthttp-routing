package routing_test

import (
	"log"

	"github.com/nantomarioni/fasthttp-routing"
	"github.com/nantomarioni/fasthttp-routing/access"
	"github.com/nantomarioni/fasthttp-routing/content"
	"github.com/nantomarioni/fasthttp-routing/fault"
	"github.com/nantomarioni/fasthttp-routing/file"
	"github.com/nantomarioni/fasthttp-routing/slash"
	"github.com/valyala/fasthttp"
)

func Example() {
	router := routing.New()

	router.Use(
		// all these handlers are shared by every route
		access.Logger(log.Printf),
		slash.Remover(fasthttp.StatusMovedPermanently),
		fault.Recovery(log.Printf),
	)

	// serve RESTful APIs
	api := router.Group("/api")
	api.Use(
		// these handlers are shared by the routes in the api group only
		content.TypeNegotiator(content.JSON, content.XML),
	)
	api.Get("/users", func(c *routing.Context) error {
		return c.Write("user list")
	})
	api.Post("/users", func(c *routing.Context) error {
		return c.Write("create a new user")
	})
	api.Put(`/users/<id:\d+>`, func(c *routing.Context) error {
		return c.Write("update user " + c.Param("id"))
	})

	// serve index file
	router.Get("/", file.Content("ui/index.html"))
	// serve files under the "ui" subdirectory
	router.Get("/*", file.Server(file.PathMap{
		"/": "/ui/",
	}))

	fasthttp.ListenAndServe(":8080", router.HandleRequest)
}
