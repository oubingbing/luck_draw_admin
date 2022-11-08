package models

import (
	"github.com/jinzhu/gorm"
	"luck-admin/enums"
	"time"
)

type Inbox struct {
	gorm.Model
	UserId 		int64 		`gorm:"column:user_id"`
	ObjectType 	int8   		`gorm:"column:object_type"`		//消息对象类型，1=抽奖，2=共享会员
	ObjectId    int64   	`gorm:"column:object_id"`
	Content     string   	`gorm:"column:content"`
	ReadAt      time.Time	`gorm:"column:read_at"`
}

type InboxPageSli []enums.InboxPage

func (Inbox) TableName() string  {
	return "inbox"
}

func (inbox *Inbox)Store(db *gorm.DB) (int64,error) {
	createResult := db.Create(inbox)
	return createResult.RowsAffected,createResult.Error
}

func (inbox *Inbox)Update(db *gorm.DB,id interface{},data map[string]interface{}) error {
	err := db.Table(inbox.TableName()).Where("id = ?",id).Updates(data).Error
	return err
}

func (inbox *Inbox)Page(db *gorm.DB,userId interface{},page *PageParam) (InboxPageSli,error) {
	var inboxList InboxPageSli
	err :=  Page(db,inbox.TableName(),page).
		Joins("left join activity on activity.id = inbox.object_id").
		Where("inbox.user_id = ?",userId).
		Select("inbox.id,inbox.user_id,object_type,object_id,inbox.content,read_at,activity.name,activity.attachments").
		Order("id desc").
		Find(&inboxList).Error

	return inboxList,err
}

func (inbox *Inbox)CountUnRead(db *gorm.DB,userId interface{}) (int,error) {
	count := 0
	err :=  db.Table(inbox.TableName()).
		Where("user_id = ?",userId).
		Where("read_at is null").
		Count(&count).Error

	return count,err
}