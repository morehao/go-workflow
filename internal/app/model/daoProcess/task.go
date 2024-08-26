package daoProcess

import (
	"fmt"
	"go-workflow/internal/app/model"
	"go-workflow/internal/pkg/constants"
	"go-workflow/internal/pkg/errorCode"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// TaskEntity 任务表表结构体
type TaskEntity struct {
	ActType       constants.ActionType `gorm:"column:act_type;comment:行为类型"`
	AgreeNum      int64                `gorm:"column:agree_num;comment:同意数量"`
	Assignee      string               `gorm:"column:assignee;comment:受派人"`
	ClaimTime     string               `gorm:"column:claim_time;comment:认领时间"`
	CreateTime    string               `gorm:"column:create_time;comment:创建时间"`
	ID            uint64               `gorm:"column:id;comment:自增ID;primaryKey"`
	IsFinished    int8                 `gorm:"column:is_finished;comment:是否完成"`
	MemberCount   int64                `gorm:"column:member_count;comment:成员数量"`
	NodeID        string               `gorm:"column:node_id;comment:节点ID"`
	ProcInstID    uint64               `gorm:"column:proc_inst_id;comment:流程实例ID"`
	Step          int64                `gorm:"column:step;comment:步骤"`
	UnCompleteNum int64                `gorm:"column:un_complete_num;comment:未完成数量"`
}

type TaskEntityList []TaskEntity

const TblNameTask = "task"

func (TaskEntity) TableName() string {
	return TblNameTask
}

type TaskCond struct {
	ID             uint64
	IDs            []uint64
	ProcInstID     uint64
	IsFinished     int8
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type TaskDao struct {
	model.Base
}

func NewTaskDao() *TaskDao {
	return &TaskDao{}
}

func (dao *TaskDao) WithTx(db *gorm.DB) *TaskDao {
	dao.Tx = db
	return dao
}

func (dao *TaskDao) Insert(c *gin.Context, entity *TaskEntity) error {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[TaskDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *TaskDao) BatchInsert(c *gin.Context, entityList TaskEntityList) error {
	db := dao.Db(c).Table(TblNameTask)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[TaskDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *TaskDao) Update(c *gin.Context, entity *TaskEntity) error {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *TaskDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *TaskDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *TaskDao) GetById(c *gin.Context, id uint64) (*TaskEntity, error) {
	var entity TaskEntity
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *TaskDao) GetByCond(c *gin.Context, cond *TaskCond) (*TaskEntity, error) {
	var entity TaskEntity
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *TaskDao) CountByCond(c *gin.Context, cond *TaskCond) (int64, error) {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)
	dao.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] CountByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return count, nil
}

func (dao *TaskDao) GetListByCond(c *gin.Context, cond *TaskCond) (TaskEntityList, error) {
	var entityList TaskEntityList
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *TaskDao) GetPageListByCond(c *gin.Context, cond *TaskCond) (TaskEntityList, int64, error) {
	db := dao.Db(c).Model(&TaskEntity{})
	db = db.Table(TblNameTask)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list TaskEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[TaskDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l TaskEntityList) ToMap() map[uint64]TaskEntity {
	m := make(map[uint64]TaskEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *TaskDao) BuildCondition(db *gorm.DB, cond *TaskCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameTask)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameTask)
		db.Where(query, cond.IDs)
	}
	if cond.ProcInstID > 0 {
		query := fmt.Sprintf("%s.proc_inst_id = ?", TblNameTask)
		db.Where(query, cond.ProcInstID)
	}
	if cond.IsFinished > 0 {
		query := fmt.Sprintf("%s.is_finished = ?", TblNameTask)
		db.Where(query, cond.IsFinished)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameTask)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameTask)
		db.Where(query, time.Unix(cond.CreatedAtEnd, 0))
	}
	if cond.IsDelete {
		db.Unscoped()
	}

	if cond.OrderField != "" {
		db.Order(cond.OrderField)
	}

	return
}
