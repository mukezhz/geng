package domain

import (
	"test/domain/hello"
	"test/domain/middlewares"

	"go.uber.org/fx"
)

var Module = fx.Options(
	middlewares.Module,
	hello.Module,
)
