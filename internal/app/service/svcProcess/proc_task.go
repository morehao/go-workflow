package svcProcess

import (
	"go-workflow/internal/app/dto/dtoProcess"
	"go-workflow/internal/app/flow"
	"go-workflow/internal/app/helper"
	"go-workflow/internal/app/model/daoProcess"
	"go-workflow/internal/pkg/constants"
	"go-workflow/internal/pkg/context"
	"go-workflow/internal/pkg/errorCode"
	"math"

	"gorm.io/gorm"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type ProcTaskSvc interface {
	Complete(c *gin.Context, req *dtoProcess.TaskCompleteReq) (*dtoProcess.TaskCompleteResp, error)
	WithDraw(c *gin.Context, req *dtoProcess.TaskWithDrawReq) (*dtoProcess.TaskWithDrawResp, error)
}

type procTaskSvc struct {
}

var _ ProcTaskSvc = (*procTaskSvc)(nil)

func NewProcTaskSvc() ProcTaskSvc {
	return &procTaskSvc{}
}

// Complete 完成审批流程任务
func (svc *procTaskSvc) Complete(c *gin.Context, req *dtoProcess.TaskCompleteReq) (*dtoProcess.TaskCompleteResp, error) {
	taskEntity, getTaskErr := daoProcess.NewTaskDao().GetById(c, req.TaskID)
	if getTaskErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskComplete] daoTask Get fail, err:%v, req:%s", getTaskErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskCompleteErr
	}
	if taskEntity == nil || taskEntity.ID == 0 {
		glog.Errorf(c, "[procTaskSvc.TaskComplete] daoTask Get fail, taskEntity is nil, req:%s", gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskNotExistErr
	}
	if taskEntity.IsFinished == constants.ProcessTaskStatusFinished {
		glog.Errorf(c, "[procTaskSvc.TaskComplete] daoTask Get fail, taskEntity is finished, req:%s", gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskCompleteErr.ResetMsg("任务已完成")
	}

	companyName, userID, userName := context.GetCompanyName(c), context.GetUserID(c), context.GetUserName(c)
	if taskEntity.ActType == constants.ActionTypeAnd {
		count, countErr := daoProcess.NewIdentitylinkDao().CountByCond(c, &daoProcess.IdentitylinkCond{
			TaskID:  req.TaskID,
			Company: companyName,
			UserID:  userID,
		})
		if countErr != nil {
			glog.Errorf(c, "[procTaskSvc.TaskComplete] daoIdentitylink CountByCond fail, err:%v, req:%s", countErr, gutils.ToJsonString(req))
			return nil, errorCode.ProcTaskCompleteErr
		}
		if count > 0 {
			return nil, errorCode.ProcTaskCompleteErr.ResetMsg("您已经审批过了，请等待他人审批！")
		}
	}

	execEntity, getExecErr := daoProcess.NewExecutionDao().GetByCond(c, &daoProcess.ExecutionCond{
		ProcInstID: taskEntity.ProcInstID,
	})
	if getExecErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskComplete] daoExecution GetByCond fail, err:%v, req:%s", getExecErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskCompleteErr
	}

	now := time.Now()
	taskEntity.Assignee = userID
	taskEntity.ClaimTime = gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS)
	if req.OperateType == constants.TaskOperateTypeApprove {
		taskEntity.AgreeNum++
	} else {
		taskEntity.IsFinished = constants.ProcessTaskStatusFinished
	}
	taskEntity.UnCompleteNum--
	if taskEntity.UnCompleteNum == 0 {
		taskEntity.IsFinished = constants.ProcessTaskStatusFinished
	}
	// 启用事务
	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		if err := daoProcess.NewTaskDao().WithTx(tx).Update(c, taskEntity); err != nil {
			return err
		}
		if taskEntity.UnCompleteNum > 0 && req.OperateType == constants.TaskOperateTypeApprove {
			identityLinkEntity := &daoProcess.IdentitylinkEntity{
				Type:       constants.IdentityTypeParticipant,
				UserID:     userID,
				UserName:   userName,
				ProcInstID: taskEntity.ProcInstID,
				Step:       taskEntity.Step,
				Company:    companyName,
				TaskID:     req.TaskID,
				Comment:    req.Comment,
			}
			if err := daoProcess.NewIdentitylinkDao().WithTx(tx).Insert(c, identityLinkEntity); err != nil {
				return err
			}
		}
		moveStageParam := &flow.MoveStageParam{
			ExecNodeList: execEntity.NodeInfos, // 执行节点列表
			UserID:       userID,
			Username:     userName,
			Company:      companyName,
			Comment:      req.Comment,
			Candidate:    req.Candidate,
			TaskID:       req.TaskID,
			ProcInstID:   taskEntity.ProcInstID,
			Step:         taskEntity.Step,
			Pass:         req.OperateType == constants.TaskOperateTypeApprove,
			FinishStatus: taskEntity.IsFinished,
			// NotifyStatus: constants.IdentityTypeNotifier, TODO:待确认如何取值
		}
		if err := flow.MoveStage(c, tx, moveStageParam); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskCreate] daoTask Complete fail, err:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskCompleteErr
	}
	return &dtoProcess.TaskCompleteResp{
		ID: taskEntity.ID,
	}, nil
}
func (svc *procTaskSvc) WithDraw(c *gin.Context, req *dtoProcess.TaskWithDrawReq) (*dtoProcess.TaskWithDrawResp, error) {
	currentTaskEntity, getCurrentTaskErr := daoProcess.NewTaskDao().GetById(c, req.TaskID)
	if getCurrentTaskErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskWithDraw] daoTask Get currentTask fail, err:%v, req:%s", getCurrentTaskErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskWithDrawErr
	}
	if currentTaskEntity == nil || currentTaskEntity.ID == 0 {
		glog.Warnf(c, "[procTaskSvc.TaskWithDraw] currentTask not exist, req:%s", gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskNotExistErr
	}
	if currentTaskEntity.Step == 0 {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("开始位置无法撤回！")
	}
	if currentTaskEntity.IsFinished == constants.ProcessTaskStatusFinished {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("已经审批结束，无法撤回！")
	}
	if currentTaskEntity.UnCompleteNum != currentTaskEntity.MemberCount {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("已经有人审批过了，无法撤回！")
	}

	lastTaskEntity, getLastTaskErr := daoProcess.NewTaskDao().GetByCond(c, &daoProcess.TaskCond{
		ProcInstID: currentTaskEntity.ProcInstID,
		IsFinished: constants.ProcessTaskStatusUnfinished,
		OrderField: "claim_time desc",
	})
	if getLastTaskErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskWithDraw] daoTask Get lastTask fail, err:%v, req:%s", getLastTaskErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskWithDrawErr
	}
	if lastTaskEntity == nil || lastTaskEntity.ID == 0 {
		glog.Warnf(c, "[procTaskSvc.TaskWithDraw] lastTask not exist, req:%s", gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskNotExistErr
	}
	companyName, userID, userName := context.GetCompanyName(c), context.GetUserID(c), context.GetUserName(c)
	if lastTaskEntity.Assignee != userID {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("只能撤回本人审批过的任务！")
	}
	sub := currentTaskEntity.Step - lastTaskEntity.Step
	if math.Abs(float64(sub)) > 1 {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("只能撤回相邻的任务！")
	}
	pass := sub < 0

	executionEntity, getExecutionErr := daoProcess.NewExecutionDao().GetByCond(c, &daoProcess.ExecutionCond{
		ProcInstID: req.ProcInstID,
	})
	if getExecutionErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskWithDraw] daoTask Get execution fail, err:%v, req:%s", getExecutionErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskWithDrawErr
	}
	if executionEntity == nil || executionEntity.ID == 0 {
		return nil, errorCode.ProcTaskWithDrawErr.ResetMsg("流程实例执行记录不存在！")
	}
	execNodeList := executionEntity.NodeInfos

	moveStageParam := &flow.MoveStageParam{
		ExecNodeList: execNodeList,
		UserID:       userID,
		Company:      companyName,
		ProcInstID:   req.ProcInstID,
		Comment:      req.Comment,
		TaskID:       req.TaskID,
		Step:         currentTaskEntity.Step,
		Pass:         pass,
	}

	currentTaskEntity.IsFinished = constants.ProcessTaskStatusFinished
	// 启用事务
	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		if err := daoProcess.NewTaskDao().WithTx(tx).Update(c, currentTaskEntity); err != nil {
			return err
		}
		if err := flow.MoveStage(c, tx, moveStageParam); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[procTaskSvc.TaskWithDraw] daoTask Complete fail, err:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcTaskWithDrawErr
	}
	return &dtoProcess.TaskWithDrawResp{}, nil
}
