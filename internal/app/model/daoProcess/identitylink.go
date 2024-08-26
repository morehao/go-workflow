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

// IdentitylinkEntity 身份链接表表结构体
type IdentitylinkEntity struct {
	Comment    string                 `gorm:"column:comment;comment:评论"`
	Company    string                 `gorm:"column:company;comment:公司"`
	Group      string                 `gorm:"column:group;comment:组"`
	ID         uint64                 `gorm:"column:id;comment:自增ID;primaryKey"`
	ProcInstID uint64                 `gorm:"column:proc_inst_id;comment:流程实例ID"`
	Step       int64                  `gorm:"column:step;comment:步骤"`
	TaskID     uint64                 `gorm:"column:task_id;comment:任务ID"`
	Type       constants.IdentityType `gorm:"column:type;comment:类型"`
	UserID     string                 `gorm:"column:user_id;comment:用户ID"`
	UserName   string                 `gorm:"column:user_name;comment:用户名称"`
}

type IdentitylinkEntityList []IdentitylinkEntity

const TblNameIdentitylink = "identitylink"

func (IdentitylinkEntity) TableName() string {
	return TblNameIdentitylink
}

type IdentitylinkCond struct {
	ID             uint64
	IDs            []uint64
	TaskID         uint64
	IsDelete       bool
	Company        string
	Group          string
	UserID         string
	Type           constants.IdentityType
	ProcInstID     uint64
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type IdentitylinkDao struct {
	model.Base
}

func NewIdentitylinkDao() *IdentitylinkDao {
	return &IdentitylinkDao{}
}

func (dao *IdentitylinkDao) WithTx(db *gorm.DB) *IdentitylinkDao {
	dao.Tx = db
	return dao
}

func (dao *IdentitylinkDao) Insert(c *gin.Context, entity *IdentitylinkEntity) error {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[IdentitylinkDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *IdentitylinkDao) BatchInsert(c *gin.Context, entityList IdentitylinkEntityList) error {
	db := dao.Db(c).Table(TblNameIdentitylink)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[IdentitylinkDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *IdentitylinkDao) Update(c *gin.Context, entity *IdentitylinkEntity) error {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *IdentitylinkDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *IdentitylinkDao) DeleteCandidateByProcInstID(c *gin.Context, procInstID uint64) error {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	if err := db.Where("proc_inst_id = ? and type = ?", procInstID, constants.IdentityTypeCandidate).Delete(&IdentitylinkEntity{}).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[IdentitylinkDao] DeleteCandidateByProcInstID fail, procInstID:%d", procInstID)
	}
	return nil
}

func (dao *IdentitylinkDao) GetById(c *gin.Context, id uint64) (*IdentitylinkEntity, error) {
	var entity IdentitylinkEntity
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *IdentitylinkDao) GetByCond(c *gin.Context, cond *IdentitylinkCond) (*IdentitylinkEntity, error) {
	var entity IdentitylinkEntity
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *IdentitylinkDao) GetListByCond(c *gin.Context, cond *IdentitylinkCond) (IdentitylinkEntityList, error) {
	var entityList IdentitylinkEntityList
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *IdentitylinkDao) GetPageListByCond(c *gin.Context, cond *IdentitylinkCond) (IdentitylinkEntityList, int64, error) {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list IdentitylinkEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l IdentitylinkEntityList) ToMap() map[uint64]IdentitylinkEntity {
	m := make(map[uint64]IdentitylinkEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *IdentitylinkDao) CountByCond(c *gin.Context, cond *IdentitylinkCond) (int64, error) {
	db := dao.Db(c).Model(&IdentitylinkEntity{})
	db = db.Table(TblNameIdentitylink)
	dao.BuildCondition(db, cond)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return 0, errorCode.ErrorDbFind.Wrapf(err, "[IdentitylinkDao] CountByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return count, nil
}

func (dao *IdentitylinkDao) BuildCondition(db *gorm.DB, cond *IdentitylinkCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameIdentitylink)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameIdentitylink)
		db.Where(query, cond.IDs)
	}
	if cond.UserID != "" {
		query := fmt.Sprintf("%s.user_id = ?", TblNameIdentitylink)
		db.Where(query, cond.UserID)
	}
	if cond.TaskID > 0 {
		query := fmt.Sprintf("%s.task_id = ?", TblNameIdentitylink)
		db.Where(query, cond.TaskID)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameIdentitylink)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameIdentitylink)
		db.Where(query, time.Unix(cond.CreatedAtEnd, 0))
	}
	if cond.Company != "" {
		query := fmt.Sprintf("%s.company = ?", TblNameIdentitylink)
		db.Where(query, cond.Company)
	}
	if cond.Group != "" {
		query := fmt.Sprintf("%s.group = ?", TblNameIdentitylink)
		db.Where(query, cond.Group)
	}
	if cond.ProcInstID > 0 {
		query := fmt.Sprintf("%s.proc_inst_id = ?", TblNameIdentitylink)
		db.Where(query, cond.ProcInstID)
	}
	if cond.Type != "" {
		query := fmt.Sprintf("%s.type = ?", TblNameIdentitylink)
		db.Where(query, cond.Type)
	}
	if cond.IsDelete {
		db.Unscoped()
	}

	if cond.OrderField != "" {
		db.Order(cond.OrderField)
	}

	return
}
