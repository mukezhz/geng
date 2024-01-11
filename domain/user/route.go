package user

import (
  "github.com/gin-gonic/gin"
  "github.com/mukezhz/geng/pkg/infrastructure"
)

type UserRoute struct {
    router *infrastructure.Router
    controller *UserController
    groupRouter *gin.RouterGroup
}

func NewUserRoute(router *infrastructure.Router, controller *UserController) *UserRoute {
	route := UserRoute{router: router, controller: controller}
  route.groupRouter = router.Group("api/user")
	route.RegisterHelloRoutes()
	return &route
}

func (r *UserRoute) RegisterHelloRoutes() {
	r.groupRouter.GET("", r.controller.HandleRoot)
}
