package svcProcDef

import (
	"go-workflow/internal/app/dto/dtoProcDef"
	"go-workflow/internal/app/flow"
	"go-workflow/internal/app/helper"
	"go-workflow/internal/app/model/daoProcDef"
	"go-workflow/internal/app/object/objCommon"
	"go-workflow/internal/app/object/objProcDef"
	"go-workflow/internal/pkg/context"
	"go-workflow/internal/pkg/errorCode"
	"time"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
)

type ProcDefSvc interface {
	Save(c *gin.Context, req *dtoProcDef.ProcDefSaveReq) (*dtoProcDef.ProcDefSaveResp, error)
	Delete(c *gin.Context, req *dtoProcDef.ProcDefDeleteReq) error
	Detail(c *gin.Context, req *dtoProcDef.ProcDefDetailReq) (*dtoProcDef.ProcDefDetailResp, error)
	PageList(c *gin.Context, req *dtoProcDef.ProcDefPageListReq) (*dtoProcDef.ProcDefPageListResp, error)
}

type procDefSvc struct {
}

var _ ProcDefSvc = (*procDefSvc)(nil)

func NewProcDefSvc() ProcDefSvc {
	return &procDefSvc{}
}

// Save 创建审批流程定义
func (svc *procDefSvc) Save(c *gin.Context, req *dtoProcDef.ProcDefSaveReq) (*dtoProcDef.ProcDefSaveResp, error) {
	if err := flow.IsValidProcessConfig(req.Resource); err != nil {
		glog.Errorf(c, "[svcProcdef.Save] flow.IsValidProcessConfig fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcDefSaveErr
	}

	companyName, userId, userName := context.GetCompanyName(c), context.GetDepartmentName(c), context.GetUserName(c)

	maxVersionEntity, getMaxVersionErr := daoProcDef.NewProcDefDao().GetByCond(c, &daoProcDef.ProcDefCond{
		Company:    companyName,
		Name:       req.Name,
		OrderField: "id desc",
	})
	if getMaxVersionErr != nil {
		glog.Errorf(c, "[svcProcdef.Save] daoProcdef GetByCond fail, err:%v, req:%s", getMaxVersionErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcDefSaveErr
	}
	var version uint64 = 1
	if maxVersionEntity != nil && maxVersionEntity.ID > 0 {
		version = maxVersionEntity.Version + 1
	}
	now := time.Now()
	insertEntity := &daoProcDef.ProcDefEntity{
		Company:    companyName,
		Name:       req.Name,
		Resource:   req.Resource,
		UserID:     userId,
		Username:   userName,
		DeployTime: gutils.TimeFormat(now, gutils.YYYY_MM_DD_HH_MM_SS),
		Version:    version,
	}
	txErr := helper.MysqlClient.Transaction(func(tx *gorm.DB) error {
		// 创建审批流程定义
		if err := daoProcDef.NewProcDefDao().WithTx(tx).Insert(c, insertEntity); err != nil {
			return err
		}
		if maxVersionEntity == nil || maxVersionEntity.ID <= 0 {
			return nil
		}
		// 迁移历史版本&&删除旧版本
		historyInsertEntity := &daoProcDef.ProcDefHistoryEntity{
			ID:         maxVersionEntity.ID,
			Company:    maxVersionEntity.Company,
			Name:       maxVersionEntity.Name,
			Resource:   maxVersionEntity.Resource,
			UserID:     maxVersionEntity.UserID,
			Username:   maxVersionEntity.Username,
			DeployTime: maxVersionEntity.DeployTime,
			Version:    maxVersionEntity.Version,
		}
		if err := daoProcDef.NewProcDefHistoryDao().WithTx(tx).Insert(c, historyInsertEntity); err != nil {
			return err
		}
		if err := daoProcDef.NewProcDefDao().WithTx(tx).Delete(c, maxVersionEntity.ID); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		glog.Errorf(c, "[svcProcdef.Save] txErr:%v, req:%s", txErr, gutils.ToJsonString(req))
		return nil, errorCode.ProcDefSaveErr
	}
	return &dtoProcDef.ProcDefSaveResp{
		ID: insertEntity.ID,
	}, nil
}

// Delete 删除审批流程定义
func (svc *procDefSvc) Delete(c *gin.Context, req *dtoProcDef.ProcDefDeleteReq) error {

	defEntity, getEntityErr := daoProcDef.NewProcDefDao().GetById(c, req.ID)
	if getEntityErr != nil {
		glog.Errorf(c, "[svcProcdef.Delete] daoProcdef GetById fail, err:%v, req:%s", getEntityErr, gutils.ToJsonString(req))
		return errorCode.ProcDefGetDetailErr
	}
	if defEntity == nil || defEntity.ID <= 0 {
		return errorCode.ProcDefNotExistErr
	}

	if err := daoProcDef.NewProcDefDao().Delete(c, req.ID); err != nil {
		glog.Errorf(c, "[svcProcdef.Delete] daoProcdef Delete fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return errorCode.ProcDefDeleteErr
	}
	return nil
}

// Detail 根据id获取审批流程定义
func (svc *procDefSvc) Detail(c *gin.Context, req *dtoProcDef.ProcDefDetailReq) (*dtoProcDef.ProcDefDetailResp, error) {
	detailEntity, err := daoProcDef.NewProcDefDao().GetById(c, req.ID)
	if err != nil {
		glog.Errorf(c, "[svcProcdef.ProcdefDetail] daoProcdef GetById fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcDefGetDetailErr
	}
	// 判断是否存在
	if detailEntity == nil || detailEntity.ID == 0 {
		return nil, errorCode.ProcDefNotExistErr
	}
	Resp := &dtoProcDef.ProcDefDetailResp{
		ID: detailEntity.ID,
		ProcDefBaseInfo: objProcDef.ProcDefBaseInfo{
			Name:     detailEntity.Name,
			Resource: detailEntity.Resource,
		},
		DeployTime: detailEntity.DeployTime,
		Version:    detailEntity.Version,
		OperatorBaseInfo: objCommon.OperatorBaseInfo{
			Company:  detailEntity.Company,
			UserID:   detailEntity.UserID,
			UserName: detailEntity.Username,
		},
	}
	return Resp, nil
}

// PageList 分页获取审批流程定义列表
func (svc *procDefSvc) PageList(c *gin.Context, req *dtoProcDef.ProcDefPageListReq) (*dtoProcDef.ProcDefPageListResp, error) {
	cond := &daoProcDef.ProcDefCond{
		Page:       req.Page,
		PageSize:   req.PageSize,
		OrderField: "id desc",
	}
	dataList, total, err := daoProcDef.NewProcDefDao().GetPageListByCond(c, cond)
	if err != nil {
		glog.Errorf(c, "[svcProcdef.ProcdefPageList] daoProcdef GetPageListByCond fail, err:%v, req:%s", err, gutils.ToJsonString(req))
		return nil, errorCode.ProcDefGetPageListErr
	}
	list := make([]dtoProcDef.ProcDefPageListItem, 0, len(dataList))
	for _, v := range dataList {
		list = append(list, dtoProcDef.ProcDefPageListItem{
			ID:         v.ID,
			Name:       v.Name,
			Company:    v.Company,
			DeployTime: v.DeployTime,
			Username:   v.Username,
			UserID:     v.UserID,
		})
	}
	return &dtoProcDef.ProcDefPageListResp{
		List:  list,
		Total: total,
	}, nil
}
