package cango

import "github.com/gin-gonic/gin"

type requestBind interface {
	bind() func(ctx *gin.Context, i any) error
}

type MustBind struct {
}

func (bind MustBind) bind() func(ctx *gin.Context, i any) error {
	return func(ctx *gin.Context, i any) error {
		return ctx.Bind(i)
	}
}

type ShouldBind struct {
}

func (bind ShouldBind) bind() func(ctx *gin.Context, i any) error {
	return func(ctx *gin.Context, i any) error {
		return ctx.ShouldBind(i)
	}
}

type NotBind struct {
}

func (bind NotBind) bind() func(ctx *gin.Context, i any) error {
	return func(ctx *gin.Context, i any) error {
		return nil
	}
}
