package {{.PackageName}}

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type {{.ModuleName}}Controller struct {
    service *{{.ModuleName}}Service
}

func New{{.ModuleName}}Controller(service *{{.ModuleName}}Service) *{{.ModuleName}}Controller {
    return &{{.ModuleName}}Controller{service: service}
}

func (ctrl *{{.ModuleName}}Controller) HandleRoot(c *gin.Context) {
    message := ctrl.service.GetMessage()
    c.JSON(http.StatusOK, gin.H{"message": message.Message})
}
