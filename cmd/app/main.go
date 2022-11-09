package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Init echo server
	e := echo.New()

	// Enable tracing middleware
	c := jaegertracing.New(e, nil)
	defer c.Close()

	// Use middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		fmt.Printf("Request Body: \t%s\n", string(reqBody))
		fmt.Printf("Response Body: \t%s\n", string(resBody))
	}))

	// Serve routers
	e.GET("/", func(c echo.Context) error {
		// Wrap slowFunc on a new span to trace it's execution passing the function arguments
		jaegertracing.TraceFunction(c, slowFunc, "Test String")
		return c.String(http.StatusOK, fmt.Sprintf("%s\n", "Hello, World!"))
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// A function to be wrapped. No need to change it's arguments due to tracing
func slowFunc(s string) {
	time.Sleep(2000 * time.Millisecond)
}
