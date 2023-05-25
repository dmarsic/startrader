package response

import (
	"github.com/gin-gonic/gin"
)

type StatusType string

const (
	Ok      StatusType = "ok"
	Error              = "error"
	Warning            = "warning"
)

type Response struct {
	Status  StatusType `json:"status"`
	Message string     `json:"message"`
	Data    any        `json:"data"`
}

func WriteResponse(c *gin.Context, e Response, statusCode int) {
	c.Status(statusCode)
	c.Header("Content-Type", "application/json")
	c.JSON(statusCode, e)
}
