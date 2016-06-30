package example

import (
	"github.com/iris-contrib/middleware/basicauth"
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

// IrisHandler creates fasthttp.RequestHandler using Iris web framework.
//
// Implemented API:
//  POST   /session  set session parameters
//  GET    /session  get session parameters
//  DELETE /session  delete session
func IrisHandler() fasthttp.RequestHandler {
	api := iris.New()

	api.Get("/things", func(c *iris.Context) {
		c.JSON(iris.StatusOK, []interface{}{
			iris.Map{
				"name":        "foo",
				"description": "foo thing",
			},
			iris.Map{
				"name":        "bar",
				"description": "bar thing",
			},
		})
	})

	api.Get("/params/:x/:y", func(c *iris.Context) {
		c.JSON(iris.StatusOK, iris.Map{
			"x":  c.Param("x"),
			"y":  c.Param("y"),
			"q":  c.URLParam("q"),
			"p1": c.PostFormValue("p1"),
			"p2": c.PostFormValue("p2"),
		})
	})

	auth := basicauth.Default(map[string]string{
		"ford": "betelgeuse7",
	})

	api.Get("/auth", auth, func(c *iris.Context) {
		c.Write("authenticated!")
	})

	api.Post("/session/set", func(c *iris.Context) {
		sess := iris.Map{}

		if err := c.ReadJSON(&sess); err != nil {
			panic(err.Error())
		}

		c.Session().Set("name", sess["name"])
	})

	api.Get("/session/get", func(c *iris.Context) {
		name := c.Session().GetString("name")

		c.JSON(iris.StatusOK, iris.Map{
			"name": name,
		})
	})

	return api.NoListen().Handler
}
