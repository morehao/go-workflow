package middleware

import (
	"go-workflow/internal/app/model/daoUser"
	"go-workflow/internal/pkg/context"
	"go-workflow/internal/pkg/errorCode"

	"github.com/morehao/go-tools/glog"

	"github.com/morehao/go-tools/gcontext/ginRender"

	"github.com/morehao/go-tools/gutils"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(context.HeaderToken)
		if token == "" {
			glog.Errorf(c, "[Middleware.Auth] token is empty")
			ginRender.Abort(c, errorCode.UserAuthErr)
			return
		}
		userId := token
		userEntity, getUserErr := daoUser.NewUserDao().GetById(c, gutils.VToUint64(userId))
		if getUserErr != nil {
			glog.Errorf(c, "[Middleware.Auth] get user fail, err: %v, userId:%s", getUserErr, userId)
			ginRender.Abort(c, errorCode.UserAuthErr)
			return
		}
		if userEntity == nil || userEntity.ID == 0 {
			ginRender.Abort(c, errorCode.UserAuthErr.ResetMsg("用户不存在"))
			return
		}
		c.Set(context.Token, token)
		c.Set(context.CompanyID, userEntity.CompanyID)
		c.Set(context.CompanyName, userEntity.CompanyName)
		c.Set(context.DepartmentID, userEntity.DepartmentID)
		c.Set(context.DepartmentName, userEntity.DepartmentName)
		c.Set(context.UserID, userId)
		c.Set(context.UserName, userEntity.Name)
		c.Next()
	}
}
