package controller

import (
	"encoding/base64"
	"x-ui/web/service"

	"github.com/gin-gonic/gin"
)

type SUBController struct {
	BaseController

	subService service.SubService
}

func NewSUBController(g *gin.RouterGroup) *SUBController {
	a := &SUBController{}
	a.initRouter(g)
	return a
}

func (a *SUBController) initRouter(g *gin.RouterGroup) {
	g = g.Group("/sub")

	g.GET("/:subid", a.subs)
}

func (a *SUBController) subs(c *gin.Context) {
	subId := c.Param("subid")
	host := c.Request.Host
	subs, err := a.subService.GetSubs(subId, host)
	result := ""
	for _, sub := range subs {
		result += sub + "\n"
	}
	if err != nil {
		c.String(400, "Error!")
	} else {
		c.String(200, base64.StdEncoding.EncodeToString([]byte(result)))
	}
}
