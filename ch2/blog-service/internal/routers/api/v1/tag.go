package v1

import (
	"block-service/global"
	"block-service/pkg/app"
	"block-service/pkg/errcode"
	"fmt"
	"github.com/gin-gonic/gin"
)

type Tag struct {
}

func NewTag() Tag {
	return Tag{}
}

//@Summary get multi tag
//@Produce json
//@Param name query string false "tag name" maxlength(100)
//@Param state query int false "tag state" Enums(0, 1) default(1)
//@Param page query int false "page number"
//@Param page_size query int false "number of content per page"
//@Success 200 {object} model.TagSwagger "200 OK"
//@Failure 400 {object} errcode.Error "400 Bad request"
//@Failure 500 {object} errcode.Error "500 Internal error"
//@Router /api/v1/tags [get]
func (t Tag) List(c *gin.Context)   {
	param := struct {
		Name string `form:"name" binding:"max=100"`
		State uint8 `form:"state,default=1" binding:"oneof=0 1"`
	}{}
	response := app.NewResponse(c)
	fmt.Printf("req: %v\n", c.Request)
	valid, errs := app.BindAndValid(c, &param)
	fmt.Printf("valid: %v, errs: %v\n", valid, errs)
	if valid == true {
		global.Logger.Errorf("app.BindAndValid errs: %v", errs)
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(errs.Errors()...))
		return
	}
	response.ToResponse(gin.H{})
}

func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}
