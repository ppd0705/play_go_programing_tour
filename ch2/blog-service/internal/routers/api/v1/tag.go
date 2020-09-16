package v1

import "github.com/gin-gonic/gin"

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
func (t Tag) List(c *gin.Context)   {}

func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}
