package main

import "github.com/labstack/echo/v4"

type CustomContext struct {
	echo.Context
	RepositoryOwner string
	Repository      string
	APIToken        string
}
