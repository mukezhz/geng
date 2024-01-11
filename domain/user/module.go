package user

import (
    "go.uber.org/fx"
)

var Module = fx.Module("user",
	fx.Options(
		fx.Provide(
			NewUserService,
			NewUserController,
			NewUserRepository,
		),
		fx.Invoke(NewUserRoute),
	),
)
