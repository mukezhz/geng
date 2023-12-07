package domain

import (
    "{{.ProjectModuleName}}/domain/features"
    "{{.ProjectModuleName}}/domain/middlewares"

    "go.uber.org/fx"
)

var Module = fx.Options(
	middlewares.Module,
	features.Module,
)
