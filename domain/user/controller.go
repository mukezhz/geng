package user

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type UserController struct {
    service *UserService
}

func NewUserController(service *UserService) *UserController {
    return &UserController{service: service}
}

func (ctrl *UserController) HandleRoot(c *gin.Context) {
    message := ctrl.service.GetMessage()
    c.JSON(http.StatusOK, gin.H{"message": message.Message})
}
