package ctrProcDef

import (
	"go-workflow/internal/app/dto/dtoProcDef"
	"go-workflow/internal/app/service/svcProcDef"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gcontext/ginRender"
)

type ProcDefCtr interface {
	Save(c *gin.Context)
	Delete(c *gin.Context)
	Detail(c *gin.Context)
	PageList(c *gin.Context)
}

type procDefCtr struct {
	procDefSvc svcProcDef.ProcDefSvc
}

var _ ProcDefCtr = (*procDefCtr)(nil)

func NewProcDefCtr() ProcDefCtr {
	return &procDefCtr{
		procDefSvc: svcProcDef.NewProcDefSvc(),
	}
}

// Save 创建审批流程定义
// @Tags 审批流程定义
// @Summary 创建审批流程定义
// @accept application/json
// @Produce application/json
// @Param req body dtoProcDef.ProcDefSaveReq true "创建审批流程定义"
// @Success 200 {object} dto.DefaultRender{data=dtoProcDef.ProcDefSaveResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procDef/save [post]
func (ctr *procDefCtr) Save(c *gin.Context) {
	var req dtoProcDef.ProcDefSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procDefSvc.Save(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}

// Delete 删除审批流程定义
// @Tags 审批流程定义
// @Summary 删除审批流程定义
// @accept application/json
// @Produce application/json
// @Param req body dtoProcDef.ProcDefDeleteReq true "删除审批流程定义"
// @Success 200 {object} dto.DefaultRender{data=string} "{"code": 0,"data": "ok","msg": "删除成功"}"
// @Router /workflow/procDef/delete [post]
func (ctr *procDefCtr) Delete(c *gin.Context) {
	var req dtoProcDef.ProcDefDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}

	if err := ctr.procDefSvc.Delete(c, &req); err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, "删除成功")
	}
}

// Detail 审批流程定义详情
// @Tags 审批流程定义
// @Summary 审批流程定义详情
// @accept application/json
// @Produce application/json
// @Param req query dtoProcDef.ProcDefDetailReq true "审批流程定义详情"
// @Success 200 {object} dto.DefaultRender{data=dtoProcDef.ProcDefDetailResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procDef/detail [get]
func (ctr *procDefCtr) Detail(c *gin.Context) {
	var req dtoProcDef.ProcDefDetailReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procDefSvc.Detail(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.SuccessWithFormat(c, res)
	}
}

// PageList 审批流程定义列表
// @Tags 审批流程定义
// @Summary 审批流程定义列表分页
// @accept application/json
// @Produce application/json
// @Param req query dtoProcDef.ProcDefPageListReq true "审批流程定义列表"
// @Success 200 {object} dto.DefaultRender{data=dtoProcDef.ProcDefPageListResp} "{"code": 0,"data": "ok","msg": "success"}"
// @Router /workflow/procDef/pageList [get]
func (ctr *procDefCtr) PageList(c *gin.Context) {
	var req dtoProcDef.ProcDefPageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		ginRender.Fail(c, err)
		return
	}
	res, err := ctr.procDefSvc.PageList(c, &req)
	if err != nil {
		ginRender.Fail(c, err)
		return
	} else {
		ginRender.Success(c, res)
	}
}
