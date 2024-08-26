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

// ProcInstEntity 审批流程实例表结构体
type ProcInstEntity struct {
	Candidate     string `gorm:"column:candidate;comment:候选人"`
	Company       string `gorm:"column:company;comment:公司名称"`
	Department    string `gorm:"column:department;comment:部门名称"`
	Duration      uint64 `gorm:"column:duration;comment:持续时间"`
	EndTime       string `gorm:"column:end_time;comment:结束时间"`
	ID            uint64 `gorm:"column:id;comment:自增ID;primaryKey"`
	IsFinished    int8   `gorm:"column:is_finished;comment:是否完成"`
	NodeID        string `gorm:"column:node_id;comment:当前节点ID"`
	ProcDefID     uint64 `gorm:"column:proc_def_id;comment:流程定义ID"`
	ProcDefName   string `gorm:"column:proc_def_name;comment:流程定义名称"`
	StartTime     string `gorm:"column:start_time;comment:开始时间"`
	StartUserID   string `gorm:"column:start_user_id;comment:发起用户ID"`
	StartUserName string `gorm:"column:start_user_name;comment:发起用户名称"`
	TaskID        uint64 `gorm:"column:task_id;comment:任务ID"`
	Title         string `gorm:"column:title;comment:流程实例标题"`
}

type ProcInstEntityList []ProcInstEntity

const TblNameProcInst = "proc_inst"

func (ProcInstEntity) TableName() string {
	return TblNameProcInst
}

type ProcInstCond struct {
	ID             uint64
	IDs            []uint64
	StartUserID    string
	Company        string
	Candidate      string
	GroupList      []string
	DepartmentList []string
	IsFinished     int8
	IsDelete       bool
	Page           int
	PageSize       int
	CreatedAtStart int64
	CreatedAtEnd   int64
	OrderField     string
}

type ProcInstDao struct {
	model.Base
}

func NewProcInstDao() *ProcInstDao {
	return &ProcInstDao{}
}

func (dao *ProcInstDao) WithTx(db *gorm.DB) *ProcInstDao {
	dao.Tx = db
	return dao
}

func (dao *ProcInstDao) Insert(c *gin.Context, entity *ProcInstEntity) error {
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)
	if err := db.Create(entity).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcInstDao] Insert fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcInstDao) BatchInsert(c *gin.Context, entityList ProcInstEntityList) error {
	db := dao.Db(c).Table(TblNameProcInst)
	if err := db.Create(entityList).Error; err != nil {
		return errorCode.ErrorDbInsert.Wrapf(err, "[ProcInstDao] BatchInsert fail, entityList:%s", gutils.ToJsonString(entityList))
	}
	return nil
}

func (dao *ProcInstDao) Update(c *gin.Context, entity *ProcInstEntity) error {
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)
	if err := db.Where("id = ?", entity.ID).Updates(entity).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcInstDao] Update fail, entity:%s", gutils.ToJsonString(entity))
	}
	return nil
}

func (dao *ProcInstDao) UpdateMap(c *gin.Context, id uint64, updateMap map[string]interface{}) error {
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)
	if err := db.Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcInstDao] UpdateMap fail, id:%d, updateMap:%s", id, gutils.ToJsonString(updateMap))
	}
	return nil
}

func (dao *ProcInstDao) Delete(c *gin.Context, id, deletedBy uint64) error {
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)
	updatedField := map[string]interface{}{
		"deleted_time": time.Now(),
		"deleted_by":   deletedBy,
	}
	if err := db.Where("id = ?", id).Updates(updatedField).Error; err != nil {
		return errorCode.ErrorDbUpdate.Wrapf(err, "[ProcInstDao] Delete fail, id:%d, deletedBy:%d", id, deletedBy)
	}
	return nil
}

func (dao *ProcInstDao) GetById(c *gin.Context, id uint64) (*ProcInstEntity, error) {
	var entity ProcInstEntity
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)
	if err := db.Where("id = ?", id).Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcInstDao] GetById fail, id:%d", id)
	}
	return &entity, nil
}

func (dao *ProcInstDao) GetByCond(c *gin.Context, cond *ProcInstCond) (*ProcInstEntity, error) {
	var entity ProcInstEntity
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entity).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcInstDao] GetById fail, cond:%s", gutils.ToJsonString(cond))
	}
	return &entity, nil
}

func (dao *ProcInstDao) GetListByCond(c *gin.Context, cond *ProcInstCond) (ProcInstEntityList, error) {
	var entityList ProcInstEntityList
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)

	dao.BuildCondition(db, cond)

	if err := db.Find(&entityList).Error; err != nil {
		return nil, errorCode.ErrorDbFind.Wrapf(err, "[ProcInstDao] GetListByCond fail, cond:%s", gutils.ToJsonString(cond))
	}
	return entityList, nil
}

func (dao *ProcInstDao) GetPageListByCond(c *gin.Context, cond *ProcInstCond) (ProcInstEntityList, int64, error) {
	db := dao.Db(c).Model(&ProcInstEntity{})
	db = db.Table(TblNameProcInst)

	dao.BuildCondition(db, cond)

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcInstDao] GetPageListByCond count fail, cond:%s", gutils.ToJsonString(cond))
	}
	if cond.PageSize > 0 && cond.Page > 0 {
		db.Offset((cond.Page - 1) * cond.PageSize).Limit(cond.PageSize)
	}
	var list ProcInstEntityList
	if err := db.Find(&list).Error; err != nil {
		return nil, 0, errorCode.ErrorDbFind.Wrapf(err, "[ProcInstDao] GetPageListByCond find fail, cond:%s", gutils.ToJsonString(cond))
	}
	return list, count, nil
}

func (l ProcInstEntityList) ToMap() map[uint64]ProcInstEntity {
	m := make(map[uint64]ProcInstEntity)
	for _, v := range l {
		m[v.ID] = v
	}
	return m
}

func (dao *ProcInstDao) BuildCondition(db *gorm.DB, cond *ProcInstCond) {
	if cond.ID > 0 {
		query := fmt.Sprintf("%s.id = ?", TblNameProcInst)
		db.Where(query, cond.ID)
	}
	if len(cond.IDs) > 0 {
		query := fmt.Sprintf("%s.id in (?)", TblNameProcInst)
		db.Where(query, cond.IDs)
	}
	if cond.Company != "" {
		query := fmt.Sprintf("%s.company = ?", TblNameProcInst)
		db.Where(query, cond.Company)
	}
	if cond.StartUserID != "" {
		query := fmt.Sprintf("%s.start_user_id = ?", TblNameProcInst)
		db.Where(query, cond.StartUserID)
	}
	if cond.Candidate != "" {
		query := fmt.Sprintf("%s.candidate = ?", TblNameProcInst)
		db.Where(query, cond.Candidate)
	}
	if len(cond.GroupList) > 0 && len(cond.DepartmentList) > 0 {
		query := fmt.Sprintf("%s.candidate in (?)", TblNameProcInst)
		db.Where(query, cond.GroupList)
	}
	if cond.IsFinished > 0 {
		query := fmt.Sprintf("%s.is_finished = ?", TblNameProcInst)
		db.Where(query, cond.IsFinished)
	}
	if cond.CreatedAtStart > 0 {
		query := fmt.Sprintf("%s.created_at >= ?", TblNameProcInst)
		db.Where(query, time.Unix(cond.CreatedAtStart, 0))
	}
	if cond.CreatedAtEnd > 0 {
		query := fmt.Sprintf("%s.created_at <= ?", TblNameProcInst)
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
