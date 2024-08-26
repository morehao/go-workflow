package daoProcess

import (
	"fmt"
	"go-workflow/internal/app/model"
	"go-workflow/internal/pkg/errorCode"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// TaskHistoryEntity 任务历史记录表表结构体
type TaskHistoryEntity struct {
	ID            uint64 `gorm:"column:id;comment:自增ID;primaryKey"`
	NodeID        string `gorm:"column:node_id;comment:节点ID"`
	Step          int32  `gorm:"column:step;comment:步骤"`
	ProcInstID    uint64 `gorm:"column:proc_inst_id;comment:流程实例ID"`
	Assignee      string `gorm:"column:assignee;comment:受派人"`
	CreateTime    string `gorm:"column:create_time;comment:创建时间"`
	ClaimTime     string `gorm:"column:claim_time;comment:认领时间"`
	MemberCount   int8   `gorm:"column:member_count;comment:成员数量"`
	UnCompleteNum int8   `gorm:"column:un_complete_num;comment:未完成数量"`
	AgreeNum      int8   `gorm:"column:agree_num;comment:同意数量"`
	ActType       string `gorm:"column:act_type;comment:行为类型"`
	IsFinished    int8   `gorm:"column:is_finished;comment:是否完成"`
}

type TaskHistoryEntityList []TaskHistoryEntity

const TblNameTaskHistory = "task_history"

func (TaskHistoryEntity) TableName() string {
	return TblNameTaskHistory
}

type TaskHistoryCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type TaskHistoryDao struct {
	model.Base
}

func NewTaskHistoryDao() *TaskHistoryDao {
	return &TaskHistoryDao{}
}

func (dao *TaskHistoryDao) WithTx(db *gorm.DB) *TaskHistoryDao {
	dao.Tx = db
	return dao
}

func (dao *TaskHistoryDao) Insert(c *gin.Context, entity *TaskHistoryEntity) error {
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[TaskHistoryDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *TaskHistoryDao) BatchInsert(c *gin.Context, entityList TaskHistoryEntityList) error {
	db := dao.Db(c).Table(TblNameTaskHistory)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[TaskHistoryDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *TaskHistoryDao) Update(c *gin.Context, entity *TaskHistoryEntity) error {
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskHistoryDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *TaskHistoryDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskHistoryDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *TaskHistoryDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[TaskHistoryDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *TaskHistoryDao) GetById(c *gin.Context, id uint64) (*TaskHistoryEntity, error) {
	var entity TaskHistoryEntity
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskHistoryDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *TaskHistoryDao) GetByCond(c *gin.Context, cond *TaskHistoryCond) (*TaskHistoryEntity, error) {
	var entity TaskHistoryEntity
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskHistoryDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *TaskHistoryDao) GetListByCond(c *gin.Context, cond *TaskHistoryCond) (TaskHistoryEntityList, error) {
	var entityList TaskHistoryEntityList
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[TaskHistoryDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *TaskHistoryDao) GetPageListByCond(c *gin.Context, cond *TaskHistoryCond) (TaskHistoryEntityList, int64, error) {
	db := dao.Db(c).Model(&TaskHistoryEntity{})
	db = db.Table(TblNameTaskHistory)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[TaskHistoryDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list TaskHistoryEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[TaskHistoryDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l TaskHistoryEntityList) ToMap() map[uint64]TaskHistoryEntity {
	m := make(map[uint64]TaskHistoryEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *TaskHistoryDao) BuildCondition(db *gorm.DB, cond *TaskHistoryCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameTaskHistory)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameTaskHistory)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameTaskHistory)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameTaskHistory)
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
