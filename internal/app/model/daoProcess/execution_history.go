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

// ExecutionHistoryEntity 执行实例历史记录表表结构体
type ExecutionHistoryEntity struct {
	ID          uint64 `gorm:"column:id;comment:自增ID;primaryKey"`
	Rev         int32  `gorm:"column:rev;comment:修订版本号"`
	ProcInstID  uint64 `gorm:"column:proc_inst_id;comment:流程实例ID"`
	ProcDefID   uint64 `gorm:"column:proc_def_id;comment:流程定义ID"`
	ProcDefName string `gorm:"column:proc_def_name;comment:流程定义名称"`
	NodeInfos   string `gorm:"column:node_infos;comment:节点信息"`
	IsActive    int8   `gorm:"column:is_active;comment:是否活跃"`
	StartTime   string `gorm:"column:start_time;comment:开始时间"`
}

type ExecutionHistoryEntityList []ExecutionHistoryEntity

const TblNameExecutionHistory = "execution_history"

func (ExecutionHistoryEntity) TableName() string {
	return TblNameExecutionHistory
}

type ExecutionHistoryCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ExecutionHistoryDao struct {
	model.Base
}

func NewExecutionHistoryDao() *ExecutionHistoryDao {
	return &ExecutionHistoryDao{}
}

func (dao *ExecutionHistoryDao) WithTx(db *gorm.DB) *ExecutionHistoryDao {
	dao.Tx = db
	return dao
}

func (dao *ExecutionHistoryDao) Insert(c *gin.Context, entity *ExecutionHistoryEntity) error {
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ExecutionHistoryDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ExecutionHistoryDao) BatchInsert(c *gin.Context, entityList ExecutionHistoryEntityList) error {
	db := dao.Db(c).Table(TblNameExecutionHistory)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ExecutionHistoryDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ExecutionHistoryDao) Update(c *gin.Context, entity *ExecutionHistoryEntity) error {
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionHistoryDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ExecutionHistoryDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionHistoryDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ExecutionHistoryDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionHistoryDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *ExecutionHistoryDao) GetById(c *gin.Context, id uint64) (*ExecutionHistoryEntity, error) {
	var entity ExecutionHistoryEntity
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionHistoryDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ExecutionHistoryDao) GetByCond(c *gin.Context, cond *ExecutionHistoryCond) (*ExecutionHistoryEntity, error) {
	var entity ExecutionHistoryEntity
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionHistoryDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ExecutionHistoryDao) GetListByCond(c *gin.Context, cond *ExecutionHistoryCond) (ExecutionHistoryEntityList, error) {
	var entityList ExecutionHistoryEntityList
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionHistoryDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ExecutionHistoryDao) GetPageListByCond(c *gin.Context, cond *ExecutionHistoryCond) (ExecutionHistoryEntityList, int64, error) {
	db := dao.Db(c).Model(&ExecutionHistoryEntity{})
	db = db.Table(TblNameExecutionHistory)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionHistoryDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ExecutionHistoryEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionHistoryDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ExecutionHistoryEntityList) ToMap() map[uint64]ExecutionHistoryEntity {
	m := make(map[uint64]ExecutionHistoryEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ExecutionHistoryDao) BuildCondition(db *gorm.DB, cond *ExecutionHistoryCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameExecutionHistory)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameExecutionHistory)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameExecutionHistory)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameExecutionHistory)
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
