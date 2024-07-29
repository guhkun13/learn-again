package handler

import (
	"github.com/gin-gonic/gin"
)

type CompressionHandler interface {
	Create(ctx *gin.Context)
	Read(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	List(ctx *gin.Context)
}
