package model

import (
	"go-workflow/internal/app/helper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Base struct {
	Tx *gorm.DB `json:"-"`
}

// Db 获取Db
func (base *Base) Db(ctx *gin.Context) (db *gorm.DB) {
	if base.Tx != nil {
		return base.Tx.WithContext(ctx)
	}

	db = helper.MysqlClient.WithContext(ctx)
	return
}
