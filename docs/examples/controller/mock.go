package controller

import "github.com/gin-gonic/gin"

type MockController struct{}

func (mc *MockController) Mock(c *gin.Context) {
	c.JSON()
}
