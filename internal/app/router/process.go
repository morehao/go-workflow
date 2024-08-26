package router

import (
	"go-workflow/internal/app/controller/ctrProcDef"
	"go-workflow/internal/app/controller/ctrProcess"

	"github.com/gin-gonic/gin"
)

// procDefRouter 初始化审批流程定义路由信息
func procDefRouter(routerGroup *gin.RouterGroup) {
	procDefCtr := ctrProcDef.NewProcDefCtr()
	procDefGroup := routerGroup.Group("procDef")
	{
		procDefGroup.POST("create", procDefCtr.Save)      // 新建审批流程定义
		procDefGroup.POST("delete", procDefCtr.Delete)    // 删除审批流程定义
		procDefGroup.GET("detail", procDefCtr.Detail)     // 根据ID获取审批流程定义
		procDefGroup.GET("pageList", procDefCtr.PageList) // 获取审批流程定义列表
	}
}

// procInstRouter 初始化审批流程实例路由信息
func procInstRouter(routerGroup *gin.RouterGroup) {
	procInstCtr := ctrProcess.NewProcInstCtr()
	procInstGroup := routerGroup.Group("procInst")
	{
		procInstGroup.POST("start", procInstCtr.Start)                    // 新建审批流程实例
		procInstGroup.POST("delete", procInstCtr.Delete)                  // 删除审批流程实例
		procInstGroup.POST("update", procInstCtr.Update)                  // 更新审批流程实例
		procInstGroup.GET("detail", procInstCtr.Detail)                   // 根据ID获取审批流程实例
		procInstGroup.GET("pageList", procInstCtr.PageList)               // 获取审批流程实例列表
		procInstGroup.GET("createdPageList", procInstCtr.CreatedPageList) // 我创建的流程实例分页列表
		procInstGroup.GET("notifyPageList", procInstCtr.NotifyPageList)   // 抄送我的流程实例分页列表
		procInstGroup.GET("todoPageList", procInstCtr.TodoPageList)       // 待我审批的流程实例分页列表
	}
}

// taskRouter 初始化审批流程任务路由信息
func taskRouter(routerGroup *gin.RouterGroup) {
	taskCtr := ctrProcess.NewTaskCtr()
	taskGroup := routerGroup.Group("task")
	{
		taskGroup.POST("complete", taskCtr.Complete) // 完成审批流程任务
		taskGroup.POST("withDraw", taskCtr.WithDraw) // 撤回审批流程任务
	}
}
