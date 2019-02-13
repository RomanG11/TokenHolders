package main

import (
	"TokenHolders/cmd/initializer"
	"TokenHolders/internal/app"
)

func main ()  {

	appl := initializer.InitApplication()

	app.RunListener(appl)
}