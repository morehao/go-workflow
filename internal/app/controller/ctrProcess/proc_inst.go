package ctrProcess

import (
	"go-workflow/internal/app/dto/dtoProcess"
	"go-workflow/internal/app/service/svcProcess"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type ProcInstCtr interface {
	Start(c *gin.Context)
	Delete(c *gin.Context)
	Update(c *gin.Context)
	Detail(c *gin.Context)
	PageList(c *gin.Context)
	CreatedPageList(c *gin.Context)
	TodoPageList(c *gin.Context)
	NotifyPageList(c *gin.Context)
}

type procInstCtr struct {
	procInstSvc svcProcess.ProcInstSvc
}

var _ ProcInstCtr = (*procInstCtr)(nil)

func NewProcInstCtr() ProcInstCtr {
	return &procInstCtr{
		procInstSvc: svcProcess.NewProcInstSvc(),
	}
}

// Start 启动审批流程实例
// @Tags 审批流程实例
// @Summary 启动审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstStartReq true "启动审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstStartResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/start [post]
func (ctr *procInstCtr) Start(c *gin.Context) {
	var req dtoProcess.ProcInstStartReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.Start(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// Delete 删除审批流程实例
// @Tags 审批流程实例
// @Summary 删除审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstDeleteReq true "删除审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /workflow/procInst/delete [post]
func (ctr *procInstCtr) Delete(c *gin.Context) {
	var req dtoProcess.ProcInstDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}

	if err := ctr.procInstSvc.Delete(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "删除成功")
	}
}

// Update 修改审批流程实例
// @Tags 审批流程实例
// @Summary 修改审批流程实例
// @accept application/json
// @Produce application/json
// @Param req body dtoProcess.ProcInstUpdateReq true "修改审批流程实例"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "修改成功"}"
// @Router /workflow/procInst/update [post]
func (ctr *procInstCtr) Update(c *gin.Context) {
	var req dtoProcess.ProcInstUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	if err := ctr.procInstSvc.Update(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "修改成功")
	}
}

// Detail 审批流程实例详情
// @Tags 审批流程实例
// @Summary 审批流程实例详情
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.ProcInstDetailReq true "审批流程实例详情"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/detail [get]
func (ctr *procInstCtr) Detail(c *gin.Context) {
	var req dtoProcess.ProcInstDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.Detail(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// PageList 审批流程实例列表
// @Tags 审批流程实例
// @Summary 审批流程实例列表分页
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.ProcInstPageListReq true "审批流程实例列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.ProcInstPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/pageList [get]
func (ctr *procInstCtr) PageList(c *gin.Context) {
	var req dtoProcess.ProcInstPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.PageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// CreatedPageList 我创建的流程实例分页列表
// @Tags 审批流程实例
// @Summary 我创建的流程实例分页列表
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.CreatedPageListReq true "我创建的流程实例分页列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.CreatedPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/createdPageList [get]
func (ctr *procInstCtr) CreatedPageList(c *gin.Context) {
	var req dtoProcess.CreatedPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.CreatedPageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// TodoPageList 待我审批的流程实例分页列表
// @Tags 审批流程实例
// @Summary 待我审批的流程实例分页列表
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.TodoPageListReq true "待我审批的流程实例分页列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.TodoPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/todoPageList [get]
func (ctr *procInstCtr) TodoPageList(c *gin.Context) {
	var req dtoProcess.TodoPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.TodoPageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// NotifyPageList 抄送我的的流程实例分页列表
// @Tags 审批流程实例
// @Summary 抄送我的的流程实例分页列表
// @accept application/json
// @Produce application/json
// @Param req query dtoProcess.NotifyPageListReq true "抄送我的的流程实例分页列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcess.NotifyPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procInst/notifyPageList [get]
func (ctr *procInstCtr) NotifyPageList(c *gin.Context) {
	var req dtoProcess.NotifyPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procInstSvc.NotifyPageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
