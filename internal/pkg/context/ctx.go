package context

import (
	"github.com/gin-gonic/gin"
)

const (
	HeaderToken = "token"

	CompanyID      = "companyId"
	CompanyName    = "companyName"
	DepartmentID   = "departmentId"
	DepartmentName = "departmentName"
	UserID         = "userId"
	UserName       = "userName"
	Token          = "token"
)

func GetClientIp(c *gin.Context) string {
	return c.ClientIP()
}

func GetUserID(c *gin.Context) string {
	return c.GetString(UserID)
}

func GetUserName(c *gin.Context) string {
	return c.GetString(UserName)
}

func GetToken(c *gin.Context) string {
	return c.GetHeader(HeaderToken)
}

func GetCompanyID(c *gin.Context) uint64 {
	return c.GetUint64(CompanyID)
}

func GetCompanyName(c *gin.Context) string {
	return c.GetString(CompanyName)
}

func GetDepartmentID(c *gin.Context) uint64 {
	return c.GetUint64(DepartmentID)
}

func GetDepartmentName(c *gin.Context) string {
	return c.GetString(DepartmentName)
}
