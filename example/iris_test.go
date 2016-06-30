package example

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
)

func irisTester(t *testing.T) *httpexpect.Expect {
	handler := IrisHandler()

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL: "http://example.com",
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func TestIrisThings(t *testing.T) {
	schema := `{
		"type": "array",
		"items": {
			"type": "object",
			"properties": {
				"name":        {"type": "string"},
				"description": {"type": "string"}
			},
			"required": ["name", "description"]
		}
	}`

	e := irisTester(t)

	things := e.GET("/things").
		Expect().
		Status(http.StatusOK).JSON()

	things.Schema(schema)

	names := things.Path("$[*].name").Array()

	names.Elements("foo", "bar")

	for n, desc := range things.Path("$..description").Array().Iter() {
		m := desc.String().Match("(.+) (.+)")

		m.Index(1).Equal(names.Element(n).String().Raw())
		m.Index(2).Equal("thing")
	}
}

func TestIrisParams(t *testing.T) {
	e := irisTester(t)

	type Form struct {
		P1 string `form:"p1"`
		P2 string `form:"p2"`
	}

	// GET /params/xxx/yyy?q=qqq
	//  p1=P1&p2=P2

	r := e.GET("/params/{x}/{y}", "xxx", "yyy").
		WithQuery("q", "qqq").WithForm(Form{P1: "P1", P2: "P2"}).
		Expect().
		Status(http.StatusOK).JSON().Object()

	r.Value("x").Equal("xxx")
	r.Value("y").Equal("yyy")
	r.Value("q").Equal("qqq")

	r.ValueEqual("p1", "P1")
	r.ValueEqual("p2", "P2")
}

func TestIrisAuth(t *testing.T) {
	e := irisTester(t)

	e.GET("/auth").
		Expect().
		Status(http.StatusUnauthorized)

	e.GET("/auth").WithBasicAuth("ford", "<bad password>").
		Expect().
		Status(http.StatusUnauthorized)

	e.GET("/auth").WithBasicAuth("ford", "betelgeuse7").
		Expect().
		Status(http.StatusOK).Body().Equal("authenticated!")
}

func TestIrisSession(t *testing.T) {
	e := irisTester(t)

	e.POST("/session/set").WithJSON(map[string]string{"name": "test"}).
		Expect().
		Status(http.StatusOK).Cookies().NotEmpty()

	r := e.GET("/session/get").
		Expect().
		Status(http.StatusOK).JSON().Object()

	r.Equal(map[string]string{
		"name": "test",
	})
}
