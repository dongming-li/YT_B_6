package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServerError resolves the given context with 500 server error
func ServerError(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"status":  "error",
		"message": "Server Error",
	})
}

// NotFound resolves the given context with 404 not found
func NotFound(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
		"status":  "error",
		"message": "Resource Not Found",
	})
}

// Unauthorized resolves the given context with 401 unauthorized
func Unauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"status":  "error",
		"message": "Unauthorized",
	})
}

// BadRequest resolves the given context with 400 bad request
func BadRequest(c *gin.Context) {
	BadRequestMessage(c, "Bad Request")
}

// BadRequestMessage resolves the given context with 400 and the given message
func BadRequestMessage(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"status":  "error",
		"message": msg,
	})
}

// Ok resolves the given context with 200 status ok
func Ok(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// Created resolves the given context with 202 status ok
func Created(c *gin.Context, model interface{}) {
	c.JSON(http.StatusCreated, model)
}

// Data resolves the given contxt status ok with data
func Data(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}
