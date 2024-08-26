package flow

import (
	"errors"
	"go-workflow/internal/app/model/daoProcess"
	"go-workflow/internal/app/object/objFlow"
	"go-workflow/internal/pkg/constants"
	"time"

	"github.com/morehao/go-tools/gutils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MoveStageParam struct {
	UserID       string
	Username     string
	Company      string
	Comment      string
	Candidate    string
	TaskID       uint64
	ProcInstID   uint64
	Step         int64
	Pass         bool
	FinishStatus int8
	NotifyStatus int8
	ExecNodeList objFlow.ExecNodeList
}

func MoveStage(c *gin.Context, tx *gorm.DB, param *MoveStageParam) error {
	// 添加上一步的参与人
	if err := AddParticipant(c, tx, param.UserID, param.Username, param.Company, param.Comment, param.TaskID, param.ProcInstID, param.Step); err != nil {
		return err
	}
	if param.Pass {
		param.Step++
		if param.Step-1 > int64(len(param.ExecNodeList)) {
			return errors.New("已经结束无法流转到下一个节点")
		}
	} else {
		param.Step--
		if param.Step < 0 {
			return errors.New("处于开始位置，无法回退到上一个节点")
		}
	}
	// 指定下一步执行人
	if len(param.Candidate) > 0 {
		param.ExecNodeList[param.Step].Approver = param.Candidate
	}
	// 判断下一流程： 如果是审批人是：抄送人
	// fmt.Printf("下一审批人类型：%s\n", param.ExecNodeList[param.Step].AproverType)
	// fmt.Println(param.ExecNodeList[param.Step].AproverType == flow.NodeTypes[flow.NOTIFIER])
	if param.ExecNodeList[param.Step].ApproverType == string(constants.IdentityTypeNotifier) {
		// 生成新的任务
		taskEntity := daoProcess.TaskEntity{
			NodeID:     string(constants.IdentityTypeNotifier),
			Step:       param.Step,
			ProcInstID: param.ProcInstID,
			IsFinished: param.FinishStatus,
		}
		if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, &taskEntity); err != nil {
			return err
		}
		// 添加抄送人
		if err := AddNotifier(c, tx, param.ExecNodeList[param.Step].Approver, param.Company, param.ProcInstID, param.Step); err != nil {
			return err
		}
		tempParam := &MoveStageParam{
			UserID:       param.UserID,
			Username:     param.Username,
			Company:      param.Company,
			Comment:      param.Comment,
			Candidate:    param.Candidate,
			TaskID:       param.TaskID,
			ProcInstID:   param.ProcInstID,
			Step:         param.Step,
			Pass:         param.Pass,
			FinishStatus: param.FinishStatus,
			NotifyStatus: param.NotifyStatus,
			ExecNodeList: param.ExecNodeList,
		}
		return MoveStage(c, tx, tempParam)
	}
	if param.Pass {
		nextStageParam := &MoveToNextStageParam{
			ExecNodeList: param.ExecNodeList,
			UserID:       param.UserID,
			Company:      param.Company,
			TaskID:       param.TaskID,
			ProcInstID:   param.ProcInstID,
			Step:         param.Step,
			FinishStatus: param.FinishStatus,
		}
		return MoveToNextStage(c, tx, nextStageParam)
	}
	prevStageParam := &MoveToPrevStageParam{
		ExecNodeList: param.ExecNodeList,
		Company:      param.Company,
		ProcInstID:   param.ProcInstID,
		Step:         param.Step,
	}
	return MoveToPrevStage(c, tx, prevStageParam)
}

type MoveToNextStageParam struct {
	ExecNodeList []objFlow.ExecNode
	UserID       string
	Company      string
	TaskID       uint64
	ProcInstID   uint64
	Step         int64
	FinishStatus int8
}

func MoveToNextStage(c *gin.Context, tx *gorm.DB, param *MoveToNextStageParam) error {
	var currentTime = gutils.TimeFormat(time.Now(), gutils.YYYY_MM_DD_HH_MM_SS)
	taskEntity := &daoProcess.TaskEntity{
		NodeID:        param.ExecNodeList[param.Step].NodeID,
		Step:          param.Step,
		CreateTime:    currentTime,
		ProcInstID:    param.ProcInstID,
		MemberCount:   param.ExecNodeList[param.Step].MemberCount,
		UnCompleteNum: param.ExecNodeList[param.Step].MemberCount,
		ActType:       param.ExecNodeList[param.Step].ActType,
	}
	procInstUpdateMap := map[string]any{
		"node_id":   param.ExecNodeList[param.Step].NodeID,
		"candidate": param.ExecNodeList[param.Step].Approver,
	}
	if (param.Step + 1) != int64(len(param.ExecNodeList)) { // 下一步不是【结束】
		// 生成新的任务
		if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, taskEntity); err != nil {
			return err
		}
		// 添加candidate group
		if err := AddCandidateGroup(c, tx, param.ExecNodeList[param.Step].Approver, param.Company, taskEntity.ID, param.ProcInstID, param.Step); err != nil {
			return err
		}
		// 更新流程实例
		procInstUpdateMap["task_id"] = taskEntity.ID
		if err := daoProcess.NewProcInstDao().WithTx(tx).UpdateMap(c, param.ProcInstID, procInstUpdateMap); err != nil {
			return err
		}
	} else { // 最后一步直接结束
		// 生成新的任务
		taskEntity.IsFinished = constants.ProcessTaskStatusFinished
		taskEntity.ClaimTime = currentTime
		if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, taskEntity); err != nil {
			return err
		}
		// 删除候选用户组
		if err := daoProcess.NewIdentitylinkDao().WithTx(tx).DeleteCandidateByProcInstID(c, param.ProcInstID); err != nil {
			return err
		}
		// 更新流程实例
		procInstUpdateMap["task_id"] = taskEntity.ID
		procInstUpdateMap["end_time"] = currentTime
		procInstUpdateMap["is_finished"] = constants.ProcessTaskStatusFinished
		procInstUpdateMap["candidate"] = "审批结束"
		if err := daoProcess.NewProcInstDao().WithTx(tx).UpdateMap(c, param.ProcInstID, procInstUpdateMap); err != nil {
			return err
		}
	}
	return nil
}

type MoveToPrevStageParam struct {
	ExecNodeList []objFlow.ExecNode
	Company      string
	ProcInstID   uint64
	Step         int64
}

// MoveToPrevStage MoveToPrevStage
// 驳回
func MoveToPrevStage(c *gin.Context, tx *gorm.DB, param *MoveToPrevStageParam) error {
	// 生成新的任务
	taskEntity := &daoProcess.TaskEntity{
		NodeID:        param.ExecNodeList[param.Step].NodeID,
		Step:          param.Step,
		CreateTime:    gutils.TimeFormat(time.Now(), gutils.YYYY_MM_DD_HH_MM_SS),
		ProcInstID:    param.ProcInstID,
		MemberCount:   param.ExecNodeList[param.Step].MemberCount,
		UnCompleteNum: param.ExecNodeList[param.Step].MemberCount,
		ActType:       param.ExecNodeList[param.Step].ActType,
	}
	if err := daoProcess.NewTaskDao().WithTx(tx).Insert(c, taskEntity); err != nil {
		return err
	}
	procInstUpdateMap := map[string]any{
		"node_id":   param.ExecNodeList[param.Step].NodeID,
		"candidate": param.ExecNodeList[param.Step].Approver,
		"task_id":   taskEntity.ID,
	}
	if err := daoProcess.NewProcInstDao().WithTx(tx).UpdateMap(c, param.ProcInstID, procInstUpdateMap); err != nil {
		return err
	}

	if param.Step == 0 { // 流程回到起始位置，注意起始位置为0,
		if err := AddCandidateUser(c, tx, param.ExecNodeList[param.Step].Approver, param.Company, taskEntity.ID, param.ProcInstID, param.Step); err != nil {
			return err
		}
		return nil
	}
	// 添加candidate group
	if err := AddCandidateGroup(c, tx, param.ExecNodeList[param.Step].Approver, param.Company, taskEntity.ID, param.ProcInstID, param.Step); err != nil {
		return err
	}
	return nil
}

// AddParticipant 添加任务参与人
func AddParticipant(c *gin.Context, tx *gorm.DB, userID, username, company, comment string, taskID, procInstID uint64, step int64) error {
	insertEntity := &daoProcess.IdentitylinkEntity{
		Type:       constants.IdentityTypeParticipant,
		UserID:     userID,
		UserName:   username,
		ProcInstID: procInstID,
		Step:       step,
		Company:    company,
		TaskID:     taskID,
		Comment:    comment,
	}
	return daoProcess.NewIdentitylinkDao().WithTx(tx).Insert(c, insertEntity)
}

// AddNotifier 添加抄送人候选用户组
func AddNotifier(c *gin.Context, tx *gorm.DB, group, company string, procInstID uint64, step int64) error {
	existNotifierEntity, getNotifierErr := daoProcess.NewIdentitylinkDao().GetByCond(c, &daoProcess.IdentitylinkCond{
		ProcInstID: procInstID,
		Company:    company,
		Group:      group,
		Type:       constants.IdentityTypeNotifier,
	})
	if getNotifierErr != nil {
		return getNotifierErr
	}
	if existNotifierEntity != nil && existNotifierEntity.ID > 0 {
		return nil
	}

	insertEntity := &daoProcess.IdentitylinkEntity{
		Group:      group,
		Type:       constants.IdentityTypeNotifier,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return daoProcess.NewIdentitylinkDao().WithTx(tx).Insert(c, insertEntity)
}

// AddCandidateGroup 添加候选用户组
func AddCandidateGroup(c *gin.Context, tx *gorm.DB, group, company string, taskID, procInstID uint64, step int64) error {
	if err := daoProcess.NewIdentitylinkDao().WithTx(tx).DeleteCandidateByProcInstID(c, procInstID); err != nil {
		return err
	}

	insertEntity := &daoProcess.IdentitylinkEntity{
		Group:      group,
		Type:       constants.IdentityTypeCandidate,
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	if err := daoProcess.NewIdentitylinkDao().WithTx(tx).Insert(c, insertEntity); err != nil {
		return err
	}
	return nil
}

// AddCandidateUser 添加候选用户
func AddCandidateUser(c *gin.Context, tx *gorm.DB, userID, company string, taskID, procInstID uint64, step int64) error {
	if err := daoProcess.NewIdentitylinkDao().WithTx(tx).DeleteCandidateByProcInstID(c, procInstID); err != nil {
		return err
	}
	insertEntity := &daoProcess.IdentitylinkEntity{
		UserID:     userID,
		Type:       constants.IdentityTypeCandidate,
		TaskID:     taskID,
		Step:       step,
		ProcInstID: procInstID,
		Company:    company,
	}
	return daoProcess.NewIdentitylinkDao().WithTx(tx).Insert(c, insertEntity)
}
