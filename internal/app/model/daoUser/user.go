package daoUser

import (
	"fmt"
	"go-workflow/internal/app/model"
	"go-workflow/internal/pkg/errorCode"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morehao/go-tools/gutils"
	"gorm.io/gorm"
)

// UserEntity 用户表结构体
type UserEntity struct {
	ID             uint64 `gorm:"column:id;comment:自增ID;primaryKey"`
	CompanyID      uint64 `gorm:"column:company_id;comment:公司id"`
	CompanyName    string `gorm:"column:company_name;comment:公司名称"`
	DepartmentID   uint64 `gorm:"column:department_id;comment:部门id"`
	DepartmentName string `gorm:"column:department_name;comment:部门名称"`
	Name           string `gorm:"column:name;comment:姓名"`
}

type UserEntityList []UserEntity

const TblNameUser = "user"

func (UserEntity) TableName() string {
	return TblNameUser
}

type UserCond struct {
	ID             uint64
	IDs            []uint64
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type UserDao struct {
	model.Base
}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (dao *UserDao) WithTx(db *gorm.DB) *UserDao {
	dao.Tx = db
	return dao
}

func (dao *UserDao) Insert(ctx *gin.Context, entity *UserEntity) error {
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[UserDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *UserDao) BatchInsert(ctx *gin.Context, entityList UserEntityList) error {
	db := dao.Db(ctx).Table(TblNameUser)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[UserDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *UserDao) Update(ctx *gin.Context, entity *UserEntity) error {
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[UserDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *UserDao) UpdateMap(ctx *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[UserDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *UserDao) Delete(ctx *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[UserDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *UserDao) GetById(ctx *gin.Context, id uint64) (*UserEntity, error) {
	var entity UserEntity
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[UserDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *UserDao) GetByCond(ctx *gin.Context, cond *UserCond) (*UserEntity, error) {
	var entity UserEntity
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[UserDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *UserDao) GetListByCond(ctx *gin.Context, cond *UserCond) (UserEntityList, error) {
	var entityList UserEntityList
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[UserDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *UserDao) GetPageListByCond(ctx *gin.Context, cond *UserCond) (UserEntityList, int64, error) {
	db := dao.Db(ctx).Model(&UserEntity{})
	db = db.Table(TblNameUser)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[UserDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list UserEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[UserDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l UserEntityList) ToMap() map[uint64]UserEntity {
	m := make(map[uint64]UserEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *UserDao) BuildCondition(db *gorm.DB, cond *UserCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameUser)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameUser)
		db.Where(query, cond.IDs)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameUser)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameUser)
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
