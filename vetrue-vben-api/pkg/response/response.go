package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	})
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: 200,
		Msg:  message,
		Data: data,
	})
}

func Fail(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg:  "fail",
		Data: nil,
	})
}

func FailWithMessage(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code: 1,
		Msg:  message,
		Data: nil,
	})
}

func FailWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  message,
		Data: nil,
	})
}
