package router

import "github.com/gin-gonic/gin"

func RegisterRouter(routerGroup *gin.RouterGroup) {
	procDefRouter(routerGroup)
	procInstRouter(routerGroup)
	taskRouter(routerGroup)
}
