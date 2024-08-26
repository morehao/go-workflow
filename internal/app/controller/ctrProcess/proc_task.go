package ctrProcess

import (
	"go-workflow/internal/app/dto/dtoProcess"
	"go-workflow/internal/app/service/svcProcess"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type ProcTaskCtr interface {
	Complete(c *gin.Context)
	WithDraw(c *gin.Context)
}

type procTaskCtr struct {
	procTaskSvc svcProcess.ProcTaskSvc
}

var _ ProcTaskCtr = (*procTaskCtr)(nil)

func NewTaskCtr() ProcTaskCtr {
	return &procTaskCtr{
		procTaskSvc: svcProcess.NewProcTaskSvc(),
	}
}

// Complete 完成审批流程任务
// @Tags 审批流程任务
// @Summary 完成审批流程任务
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.TaskCompleteReq true "完成审批流程任务"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.TaskCompleteResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/task/complete [post]
func (ctr *procTaskCtr) Complete(c *gin.Context) {
	var req dtoProcess.TaskCompleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procTaskSvc.Complete(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// WithDraw 撤回审批流程任务
// @Tags 审批流程任务
// @Summary 撤回审批流程任务
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.TaskWithDrawReq true "撤回审批流程任务"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.TaskWithDrawResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/task/withDraw [post]
func (ctr *procTaskCtr) WithDraw(c *gin.Context) {
	var req dtoProcess.TaskWithDrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procTaskSvc.WithDraw(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
