package svcProcess

import (
	"fmt"
	"go-workflow/internal/app/dto/dtoProcess"
	"go-workflow/internal/app/flow"
	"go-workflow/internal/app/helper"
	"go-workflow/internal/app/model/daoProcDef"
	"go-workflow/internal/app/model/daoProcess"
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objFlow"
	"go-workflow/internal/app/object/objProcess"
	"go-workflow/internal/pkg/constants"
	"go-workflow/internal/pkg/context"
	"go-workflow/internal/pkg/errorCode"
	"strings"

	"gorm.io/gorm"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type ProcInstSvc interface {
	Start(c *gin.Context, req *dtoProcess.ProcInstStartReq) (*dtoProcess.ProcInstStartResp, error)
	Delete(c *gin.Context, req *dtoProcess.ProcInstDeleteReq) error
	Update(c *gin.Context, req *dtoProcess.ProcInstUpdateReq) error
	Detail(c *gin.Context, req *dtoProcess.ProcInstDetailReq) (*dtoProcess.ProcInstDetailResp, error)
	PageList(c *gin.Context, req *dtoProcess.ProcInstPageListReq) (*dtoProcess.ProcInstPageListResp, error)
	CreatedPageList(c *gin.Context, req *dtoProcess.CreatedPageListReq) (*dtoProcess.CreatedPageListResp, error)
	TodoPageList(c *gin.Context, req *dtoProcess.TodoPageListReq) (*dtoProcess.TodoPageListResp, error)
	NotifyPageList(c *gin.Context, req *dtoProcess.NotifyPageListReq) (*dtoProcess.NotifyPageListResp, error)
}

type procInstSvc struct {
}

var _ ProcInstSvc = (*procInstSvc)(nil)

func NewProcInstSvc() ProcInstSvc {
	return &procInstSvc{}
}

// Start 创建审批流程实例
func (svc *procInstSvc) Start(c *gin.Context, req *dtoProcess.ProcInstStartReq) (*dtoProcess.ProcInstStartResp, error) {
	companyName, userID, userName, departmentName := context.GetCompanyName(c), context.GetUserID(c), context.GetUserName(c), context.GetDepartmentName(c)

	procDefEntity, getProcDefErr := daoProcDef.NewProcDefDao().GetByCond(c, &daoProcDef.ProcDefCond{
		Company: companyName,
		Name:    req.ProcDefName,
	})
	if getProcDefErr != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstCreate] daoProcdef GetByCond fail, err:%v, req:%s", getProcDefErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}
	if procDefEntity == nil || procDefEntity.ID == 0 {
		return nil, errorCode.ProcDefNotExistErr
	}
	now := time.Now()
	instEntity := &daoProcess.ProcInstEntity{
		ProcDefID:     procDefEntity.ID,
		ProcDefName:   procDefEntity.Name,
		NodeID:        procDefEntity.Resource.NodeID,
		Title:         req.Title,
		Company:       companyName,
		Department:    departmentName,
		StartUserID:   userID,
		StartUserName: userName,
		StartTime:     gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
		IsFinished:    constants.ProcessTaskStatusUnfinished,
	}
	execNodeLinkedList, parseErr := flow.ParseProcessConfig(procDefEntity.Resource, req.Var)
	if parseErr != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstCreate] parse process config fail, err:%v, req:%s", parseErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}
	execNodeLinkedList.PushBack(objFlow.ExecNode{
		NodeID: constants.NodeIDTextEnd,
	})
	execNodeLinkedList.PushFront(objFlow.ExecNode{
		NodeID:   constants.NodeIDTextStart,
		Type:     constants.NodeInfoTypeStarter,
		Approver: userID,
	})
	var execNodeList objFlow.ExecNodeList
	if err := gutils.LinkedListToArray(execNodeLinkedList, &execNodeList); err != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstCreate] parse process config fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}

	executionEntity := &daoProcess.ExecutionEntity{
		ProcDefID: procDefEntity.ID,
		NodeInfos: execNodeList,
	}

	taskEntity := &daoProcess.TaskEntity{
		NodeID:        constants.NodeIDTextStart,
		Assignee:      userID,
		IsFinished:    constants.ProcessTaskStatusFinished,
		ClaimTime:     gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
		Step:          0,
		MemberCount:   1,
		UnCompleteNum: 0,
		ActType:       constants.ActionTypeOr,
		AgreeNum:      1,
	}
	if len(execNodeList) > 0 {
		if execNodeList[0].ActType == constants.ActionTypeAnd {
			taskEntity.UnCompleteNum = execNodeList[0].MemberCount
			taskEntity.MemberCount = execNodeList[0].MemberCount
		}
	}

	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		if err := daoProcess.NewProcInstDao().WithTx(tx).Insert(c, instEntity); err != nil {
			return err
		}
		executionEntity.ProcInstID = instEntity.ID
		if err := daoProcess.NewExecutionDao().Insert(c, executionEntity); err != nil {
			return err
		}
		taskEntity.ProcInstID = instEntity.ID
		if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, taskEntity); err != nil {
			return err
		}
		moveStageParam := &flow.MoveStageParam{
			ExecNodeList: execNodeList,
			UserID:       userID,
			Company:      companyName,
			ProcInstID:   instEntity.ID,
			Comment:      "开始流程",
			TaskID:       taskEntity.ID,
			Step:         0,
			Pass:         true,
		}
		if err := flow.MoveStage(c, tx, moveStageParam); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstCreate] tx fail, err:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreateErr
	}

	return &dtoProcess.ProcInstStartResp{
		ID: instEntity.ID,
	}, nil
}

// Delete 删除审批流程实例
func (svc *procInstSvc) Delete(c *gin.Context, req *dtoProcess.ProcInstDeleteReq) error {
	if err := daoProcess.NewProcInstDao().Delete(c, req.ID, 0); err != nil {
		glog.Errorf(c, "[procInstSvc.Delete] daoProcInst Delete fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcInstDeleteErr
	}
	return nil
}

// Update 更新审批流程实例
func (svc *procInstSvc) Update(c *gin.Context, req *dtoProcess.ProcInstUpdateReq) error {
	updateEntity := &daoProcess.ProcInstEntity{
		ID: req.ID,
	}
	if err := daoProcess.NewProcInstDao().Update(c, updateEntity); err != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstUpdate] daoProcInst Update fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcInstUpdateErr
	}
	return nil
}

// Detail 根据id获取审批流程实例
func (svc *procInstSvc) Detail(c *gin.Context, req *dtoProcess.ProcInstDetailReq) (*dtoProcess.ProcInstDetailResp, error) {
	detailEntity, err := daoProcess.NewProcInstDao().GetById(c, req.ID)
	if err != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstDetail] daoProcInst GetById fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstGetDetailErr
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, errorCode.ProcInstNotExistErr
	}
	Resp := &dtoProcess.ProcInstDetailResp{
		ID: detailEntity.ID,
		ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
			Candidate:     detailEntity.Candidate,
			Company:       detailEntity.Company,
			Department:    detailEntity.Department,
			Duration:      detailEntity.Duration,
			EndTime:       detailEntity.EndTime,
			IsFinished:    detailEntity.IsFinished,
			NodeID:        detailEntity.NodeID,
			ProcDefID:     detailEntity.ProcDefID,
			ProcDefName:   detailEntity.ProcDefName,
			StartTime:     detailEntity.StartTime,
			StartUserID:   detailEntity.StartUserID,
			StartUserName: detailEntity.StartUserName,
			TaskID:        detailEntity.TaskID,
			Title:         detailEntity.Title,
		},
		OperatorBaseInfo: objCommon.OperatorBaseInfo{},
	}
	return Resp, nil
}

// PageList 分页获取审批流程实例列表
func (svc *procInstSvc) PageList(c *gin.Context, req *dtoProcess.ProcInstPageListReq) (*dtoProcess.ProcInstPageListResp, error) {
	cond := &daoProcess.ProcInstCond{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	dataList, total, err := daoProcess.NewProcInstDao().GetPageListByCond(c, cond)
	if err != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstPageList] daoProcInst GetPageListByCond fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstGetPageListErr
	}
	list := make([]dtoProcess.ProcInstPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtoProcess.ProcInstPageListItem{
			ID: v.ID,
			ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
				Candidate:     v.Candidate,
				Company:       v.Company,
				Department:    v.Department,
				Duration:      v.Duration,
				EndTime:       v.EndTime,
				IsFinished:    v.IsFinished,
				NodeID:        v.NodeID,
				ProcDefID:     v.ProcDefID,
				ProcDefName:   v.ProcDefName,
				StartTime:     v.StartTime,
				StartUserID:   v.StartUserID,
				StartUserName: v.StartUserName,
				TaskID:        v.TaskID,
				Title:         v.Title,
			},
		})
	}
	return &dtoProcess.ProcInstPageListResp{
		List:  list,
		Total: total,
	}, nil
}
func (svc *procInstSvc) CreatedPageList(c *gin.Context, req *dtoProcess.CreatedPageListReq) (*dtoProcess.CreatedPageListResp, error) {
	companyName, userID := context.GetCompanyName(c), context.GetUserID(c)
	cond := &daoProcess.ProcInstCond{
		Company:     companyName,
		StartUserID: userID,
		Page:        req.Page,
		PageSize:    req.PageSize,
		OrderField:  "start_time desc",
	}
	entityList, total, getListErr := daoProcess.NewProcInstDao().GetPageListByCond(c, cond)
	if getListErr != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstPageList] daoProcInst GetPageListByCond fail, err:%v, req:%s", getListErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreatedPageListErr
	}
	list := make([]dtoProcess.ProcInstPageListItem, 0, len(entityList))
	for _, v := range entityList {
		list = append(list, dtoProcess.ProcInstPageListItem{
			ID: v.ID,
			ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
				Candidate:     v.Candidate,
				Company:       v.Company,
				Department:    v.Department,
				Duration:      v.Duration,
				EndTime:       v.EndTime,
				IsFinished:    v.IsFinished,
				NodeID:        v.NodeID,
				ProcDefID:     v.ProcDefID,
				ProcDefName:   v.ProcDefName,
				StartTime:     v.StartTime,
				StartUserID:   v.StartUserID,
				StartUserName: v.StartUserName,
				TaskID:        v.TaskID,
				Title:         v.Title,
			},
		})
	}
	return &dtoProcess.CreatedPageListResp{
		List:  list,
		Total: total,
	}, nil
}
func (svc *procInstSvc) TodoPageList(c *gin.Context, req *dtoProcess.TodoPageListReq) (*dtoProcess.TodoPageListResp, error) {
	companyName, userID := context.GetCompanyName(c), context.GetUserID(c)
	db := helper.MysqlClient
	var orCondList []string
	if len(req.DepartmentList) > 0 {
		departmentCond := fmt.Sprintf("%s.department in (%s)", daoProcess.TblNameProcInst, strings.Join(req.DepartmentList, ","))
		candidateCond := fmt.Sprintf("%s.candidate = %s", daoProcess.TblNameProcInst, "'主管'")
		orCondList = append(orCondList, fmt.Sprintf("(%s and %s)", departmentCond, candidateCond))
	}
	if len(req.GroupList) > 0 {
		orCondList = append(orCondList, fmt.Sprintf("%s.candidate in (%s)", daoProcess.TblNameProcInst, strings.Join(req.GroupList, ",")))
	}
	orCondList = append(orCondList, fmt.Sprintf("%s.candidate = %s", daoProcess.TblNameProcInst, userID))
	db = db.Where(strings.Join(orCondList, " or "))
	cond := &daoProcess.ProcInstCond{
		Company:    companyName,
		IsFinished: constants.ProcessTaskStatusUnfinished,
		Page:       req.Page,
		PageSize:   req.PageSize,
		OrderField: "start_time desc",
	}
	entityList, total, getListErr := daoProcess.NewProcInstDao().WithTx(db).GetPageListByCond(c, cond)
	if getListErr != nil {
		glog.Errorf(c, "[procInstSvc.TodoPageList] daoProcInst GetPageListByCond fail, err:%v, req:%s", getListErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstCreatedPageListErr
	}
	list := make([]dtoProcess.ProcInstPageListItem, 0, len(entityList))
	for _, v := range entityList {
		list = append(list, dtoProcess.ProcInstPageListItem{
			ID: v.ID,
			ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
				Candidate:     v.Candidate,
				Company:       v.Company,
				Department:    v.Department,
				Duration:      v.Duration,
				EndTime:       v.EndTime,
				IsFinished:    v.IsFinished,
				NodeID:        v.NodeID,
				ProcDefID:     v.ProcDefID,
				ProcDefName:   v.ProcDefName,
				StartTime:     v.StartTime,
				StartUserID:   v.StartUserID,
				StartUserName: v.StartUserName,
				TaskID:        v.TaskID,
				Title:         v.Title,
			},
		})
	}
	return &dtoProcess.TodoPageListResp{
		List:  list,
		Total: total,
	}, nil
}
func (svc *procInstSvc) NotifyPageList(c *gin.Context, req *dtoProcess.NotifyPageListReq) (*dtoProcess.NotifyPageListResp, error) {
	companyName, userID := context.GetCompanyName(c), context.GetUserID(c)
	// TODO:group条件处理
	db := helper.MysqlClient
	if len(req.GroupList) > 0 {
		db = db.Joins("join identitylink on proc_inst.id = identitylink.proc_inst_id and (identitylink.user_id = ? or identitylink.group in (?))", userID, req.GroupList)
	} else {
		db = db.Joins("join identitylink on proc_inst.id = identitylink.proc_inst_id and identitylink.user_id = ?", userID)
	}
	cond := &daoProcess.ProcInstCond{
		Company:  companyName,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	entityList, total, getListErr := daoProcess.NewProcInstDao().WithTx(db).GetPageListByCond(c, cond)
	if getListErr != nil {
		glog.Errorf(c, "[procInstSvc.ProcInstPageList] daoProcInst GetPageListByCond fail, err:%v, req:%s", getListErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcInstNotifyPageListErr
	}
	list := make([]dtoProcess.ProcInstPageListItem, 0, len(entityList))
	for _, v := range entityList {
		list = append(list, dtoProcess.ProcInstPageListItem{
			ID: v.ID,
			ProcInstBaseInfo: objProcess.ProcInstBaseInfo{
				Candidate:     v.Candidate,
				Company:       v.Company,
				Department:    v.Department,
				Duration:      v.Duration,
				EndTime:       v.EndTime,
				IsFinished:    v.IsFinished,
				NodeID:        v.NodeID,
				ProcDefID:     v.ProcDefID,
				ProcDefName:   v.ProcDefName,
				StartTime:     v.StartTime,
				StartUserID:   v.StartUserID,
				StartUserName: v.StartUserName,
				TaskID:        v.TaskID,
				Title:         v.Title,
			},
		})
	}

	return &dtoProcess.NotifyPageListResp{
		List:  list,
		Total: total,
	}, nil
}
