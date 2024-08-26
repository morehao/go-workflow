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

// ProcDefEntity 审批流程定义表结构体
type ProcDefEntity struct {
	ID         uint64        `gorm:"column:id;comment:自增id;primaryKey"`
	Company    string        `gorm:"column:company;comment:公司名称"`
	Name       string        `gorm:"column:name;comment:流程名称"`
	Resource   *objFlow.Node `gorm:"column:resource;type:json;comment:流程配置"`
	UserID     string        `gorm:"column:userid;comment:用户id"`
	Username   string        `gorm:"column:username;comment:用户名称"`
	DeployTime string        `gorm:"column:deploy_time;comment:部署时间"`
	Version    uint64        `gorm:"column:version;comment:流程版本"`
}

type ProcDefEntityList []ProcDefEntity

const TblNameProcDef = "procdef"

func (ProcDefEntity) TableName() string {
	return TblNameProcDef
}

type ProcDefCond struct {
	ID             uint64
	IDs            []uint64
	Company        string
	Name           string
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ProcDefDao struct {
	model.Base
}

func NewProcDefDao() *ProcDefDao {
	return &ProcDefDao{}
}

func (dao *ProcDefDao) WithTx(db *gorm.DB) *ProcDefDao {
	dao.Tx = db
	return dao
}

func (dao *ProcDefDao) Insert(c *gin.Context, entity *ProcDefEntity) error {
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcDefDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcDefDao) BatchInsert(c *gin.Context, entityList ProcDefEntityList) error {
	db := dao.Db(c).Table(TblNameProcDef)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcDefDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ProcDefDao) Update(c *gin.Context, entity *ProcDefEntity) error {
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcDefDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ProcDefDao) Delete(c *gin.Context, id uint64) error {
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)
	if err := db.Delete(&ProcDefEntity{}, id).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcDefDao] Delete fail, id:%d", id)
	}
	return nil
}

func (dao *ProcDefDao) GetById(c *gin.Context, id uint64) (*ProcDefEntity, error) {
	var entity ProcDefEntity
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ProcDefDao) GetByCond(c *gin.Context, cond *ProcDefCond) (*ProcDefEntity, error) {
	var entity ProcDefEntity
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ProcDefDao) GetListByCond(c *gin.Context, cond *ProcDefCond) (ProcDefEntityList, error) {
	var entityList ProcDefEntityList
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ProcDefDao) GetPageListByCond(c *gin.Context, cond *ProcDefCond) (ProcDefEntityList, int64, error) {
	db := dao.Db(c).Model(&ProcDefEntity{})
	db = db.Table(TblNameProcDef)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ProcDefEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcDefDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ProcDefEntityList) ToMap() map[uint64]ProcDefEntity {
	m := make(map[uint64]ProcDefEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ProcDefDao) BuildCondition(db *gorm.DB, cond *ProcDefCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameProcDef)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameProcDef)
		db.Where(query, cond.IDs)
	}
	if cond.Company != "" {
		query := fmt.Sprintf("%s.company = ?", TblNameProcDef)
		db.Where(query, cond.Company)
	}
	if cond.Name != "" {
		query := fmt.Sprintf("%s.name = ?", TblNameProcDef)
		db.Where(query, cond.Name)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameProcDef)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameProcDef)
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
