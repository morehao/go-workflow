package daoProcDef

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

// ProcDefHistoryEntity 审批流程定义历史记录表结构体
type ProcDefHistoryEntity struct {
	ID         uint64        `gorm:"column:id;comment:自增id;primaryKey"`
	Company    string        `gorm:"column:company;comment:公司名称"`
	Name       string        `gorm:"column:name;comment:流程名称"`
	Resource   *objFlow.Node `gorm:"column:resource;type:json;comment:流程配置"`
	UserID     string        `gorm:"column:userid;comment:用户id"`
	Username   string        `gorm:"column:username;comment:用户名称"`
	DeployTime string        `gorm:"column:deploy_time;comment:部署时间"`
	Version    uint64        `gorm:"column:version;comment:流程版本"`
}

type ProcDefHistoryEntityList []ProcDefHistoryEntity

const TblNameProcDefHistory = "procdef_history"

func (ProcDefHistoryEntity) TableName() string {
	return TblNameProcDefHistory
}

type ProcDefHistoryCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ProcDefHistoryDao struct {
	model.Base
}

func NewProcDefHistoryDao() *ProcDefHistoryDao {
	return &ProcDefHistoryDao{}
}

func (dao *ProcDefHistoryDao) WithTx(db *gorm.DB) *ProcDefHistoryDao {
	dao.Tx = db
	return dao
}

func (dao *ProcDefHistoryDao) Insert(c *gin.Context, entity *ProcDefHistoryEntity) error {
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcDefHistoryDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcDefHistoryDao) BatchInsert(c *gin.Context, entityList ProcDefHistoryEntityList) error {
	db := dao.Db(c).Table(TblNameProcDefHistory)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcDefHistoryDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ProcDefHistoryDao) Update(c *gin.Context, entity *ProcDefHistoryEntity) error {
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefHistoryDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcDefHistoryDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefHistoryDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ProcDefHistoryDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefHistoryDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *ProcDefHistoryDao) GetById(c *gin.Context, id uint64) (*ProcDefHistoryEntity, error) {
	var entity ProcDefHistoryEntity
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefHistoryDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ProcDefHistoryDao) GetByCond(c *gin.Context, cond *ProcDefHistoryCond) (*ProcDefHistoryEntity, error) {
	var entity ProcDefHistoryEntity
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefHistoryDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ProcDefHistoryDao) GetListByCond(c *gin.Context, cond *ProcDefHistoryCond) (ProcDefHistoryEntityList, error) {
	var entityList ProcDefHistoryEntityList
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefHistoryDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ProcDefHistoryDao) GetPageListByCond(c *gin.Context, cond *ProcDefHistoryCond) (ProcDefHistoryEntityList, int64, error) {
	db := dao.Db(c).Model(&ProcDefHistoryEntity{})
	db = db.Table(TblNameProcDefHistory)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefHistoryDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ProcDefHistoryEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefHistoryDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ProcDefHistoryEntityList) ToMap() map[uint64]ProcDefHistoryEntity {
	m := make(map[uint64]ProcDefHistoryEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ProcDefHistoryDao) BuildCondition(db *gorm.DB, cond *ProcDefHistoryCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameProcDefHistory)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameProcDefHistory)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameProcDefHistory)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameProcDefHistory)
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
