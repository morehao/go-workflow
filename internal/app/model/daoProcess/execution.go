package daoProcess

import (
	"fmt"
	"go-workflow/internal/app/model"
	"go-workflow/internal/app/object/objFlow"
	"go-workflow/internal/pkg/errorCode"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// ExecutionEntity 执行实例表表结构体
type ExecutionEntity struct {
	ID          uint64               `gorm:"column:id;comment:自增ID;primaryKey"`
	Rev         uint64               `gorm:"column:rev;comment:修订版本号"`
	ProcInstID  uint64               `gorm:"column:proc_inst_id;comment:流程实例ID"`
	ProcDefID   uint64               `gorm:"column:proc_def_id;comment:流程定义ID"`
	ProcDefName string               `gorm:"column:proc_def_name;comment:流程定义名称"`
	NodeInfos   objFlow.ExecNodeList `gorm:"column:node_infos;comment:节点信息"`
	IsActive    int8                 `gorm:"column:is_active;comment:是否活跃"`
	StartTime   string               `gorm:"column:start_time;comment:开始时间"`
}

type ExecutionEntityList []ExecutionEntity

const TblNameExecution = "execution"

func (ExecutionEntity) TableName() string {
	return TblNameExecution
}

type ExecutionCond struct {
	ID             uint64
	IDs            []uint64
	ProcInstID     uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ExecutionDao struct {
	model.Base
}

func NewExecutionDao() *ExecutionDao {
	return &ExecutionDao{}
}

func (dao *ExecutionDao) WithTx(db *gorm.DB) *ExecutionDao {
	dao.Tx = db
	return dao
}

func (dao *ExecutionDao) Insert(c *gin.Context, entity *ExecutionEntity) error {
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ExecutionDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ExecutionDao) BatchInsert(c *gin.Context, entityList ExecutionEntityList) error {
	db := dao.Db(c).Table(TblNameExecution)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ExecutionDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ExecutionDao) Update(c *gin.Context, entity *ExecutionEntity) error {
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ExecutionDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ExecutionDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ExecutionDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *ExecutionDao) GetById(c *gin.Context, id uint64) (*ExecutionEntity, error) {
	var entity ExecutionEntity
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ExecutionDao) GetByCond(c *gin.Context, cond *ExecutionCond) (*ExecutionEntity, error) {
	var entity ExecutionEntity
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ExecutionDao) GetListByCond(c *gin.Context, cond *ExecutionCond) (ExecutionEntityList, error) {
	var entityList ExecutionEntityList
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ExecutionDao) GetPageListByCond(c *gin.Context, cond *ExecutionCond) (ExecutionEntityList, int64, error) {
	db := dao.Db(c).Model(&ExecutionEntity{})
	db = db.Table(TblNameExecution)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ExecutionEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ExecutionDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ExecutionEntityList) ToMap() map[uint64]ExecutionEntity {
	m := make(map[uint64]ExecutionEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ExecutionDao) BuildCondition(db *gorm.DB, cond *ExecutionCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameExecution)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameExecution)
		db.Where(query, cond.IDs)
	}
	if cond.ProcInstID > 0 {
		query := fmt.Sprintf("%s.proc_inst_id = ?", TblNameExecution)
		db.Where(query, cond.ProcInstID)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameExecution)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameExecution)
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
