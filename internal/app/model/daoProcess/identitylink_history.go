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

// IdentitylinkHistoryEntity 身份链接历史记录表表结构体
type IdentitylinkHistoryEntity struct {
	Comment    string `gorm:"column:comment;comment:评论"`
	Company    string `gorm:"column:company;comment:公司"`
	Group      string `gorm:"column:group;comment:组"`
	ID         uint64 `gorm:"column:id;comment:自增ID;primaryKey"`
	ProcInstID uint64 `gorm:"column:proc_inst_id;comment:流程实例ID"`
	Step       int32  `gorm:"column:step;comment:步骤"`
	TaskID     uint64 `gorm:"column:task_id;comment:任务ID"`
	Type       string `gorm:"column:type;comment:类型"`
	UserID     string `gorm:"column:user_id;comment:用户ID"`
	UserName   string `gorm:"column:user_name;comment:用户名称"`
}

type IdentitylinkHistoryEntityList []IdentitylinkHistoryEntity

const TblNameIdentitylinkHistory = "identitylink_history"

func (IdentitylinkHistoryEntity) TableName() string {
	return TblNameIdentitylinkHistory
}

type IdentitylinkHistoryCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type IdentitylinkHistoryDao struct {
	model.Base
}

func NewIdentitylinkHistoryDao() *IdentitylinkHistoryDao {
	return &IdentitylinkHistoryDao{}
}

func (dao *IdentitylinkHistoryDao) WithTx(db *gorm.DB) *IdentitylinkHistoryDao {
	dao.Tx = db
	return dao
}

func (dao *IdentitylinkHistoryDao) Insert(c *gin.Context, entity *IdentitylinkHistoryEntity) error {
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[IdentitylinkHistoryDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *IdentitylinkHistoryDao) BatchInsert(c *gin.Context, entityList IdentitylinkHistoryEntityList) error {
	db := dao.Db(c).Table(TblNameIdentitylinkHistory)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[IdentitylinkHistoryDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *IdentitylinkHistoryDao) Update(c *gin.Context, entity *IdentitylinkHistoryEntity) error {
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkHistoryDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *IdentitylinkHistoryDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkHistoryDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *IdentitylinkHistoryDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkHistoryDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *IdentitylinkHistoryDao) GetById(c *gin.Context, id uint64) (*IdentitylinkHistoryEntity, error) {
	var entity IdentitylinkHistoryEntity
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkHistoryDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *IdentitylinkHistoryDao) GetByCond(c *gin.Context, cond *IdentitylinkHistoryCond) (*IdentitylinkHistoryEntity, error) {
	var entity IdentitylinkHistoryEntity
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkHistoryDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *IdentitylinkHistoryDao) GetListByCond(c *gin.Context, cond *IdentitylinkHistoryCond) (IdentitylinkHistoryEntityList, error) {
	var entityList IdentitylinkHistoryEntityList
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkHistoryDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *IdentitylinkHistoryDao) GetPageListByCond(c *gin.Context, cond *IdentitylinkHistoryCond) (IdentitylinkHistoryEntityList, int64, error) {
	db := dao.Db(c).Model(&IdentitylinkHistoryEntity{})
	db = db.Table(TblNameIdentitylinkHistory)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkHistoryDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list IdentitylinkHistoryEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkHistoryDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l IdentitylinkHistoryEntityList) ToMap() map[uint64]IdentitylinkHistoryEntity {
	m := make(map[uint64]IdentitylinkHistoryEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *IdentitylinkHistoryDao) BuildCondition(db *gorm.DB, cond *IdentitylinkHistoryCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameIdentitylinkHistory)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameIdentitylinkHistory)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameIdentitylinkHistory)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameIdentitylinkHistory)
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
