package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/nbvghost/gweb/tool"
)

type AppointmentDao = BaseDao

func (AppointmentDao) FindAppointmentByName(DB *gorm.DB, Name string) *Appointment {
	appointment := &Appointment{}
	err := DB.Where("Name=?", Name).First(appointment).Error //SelectOne(user, "select * from User where Tel=?", Tel)
	tool.CheckError(err)
	return appointment
}
