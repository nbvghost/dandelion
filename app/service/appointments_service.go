package service

import (
	"dandelion/app/service/dao"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/nbvghost/gweb/tool"
)

type AppointmentService struct {
	dao.AppointmentDao
}

func (self AppointmentService) FindAppointmentOfPaging(Index int, CompanyID uint64) (List []dao.Appointment, Total int) {

	p := Orm.Where("CompanyID = ?", CompanyID).Order("ID desc")

	dao.SelectPaging(Index, p, &List, &Total)

	return
}
func (self AppointmentService) GetAppointment(ID uint64) *dao.Appointment {
	target := &dao.Appointment{}
	self.Get(Orm, ID, target)
	return target
}
func (self AppointmentService) SaveAppointment(appointment *dao.Appointment) *dao.ActionStatus {

	as := &dao.ActionStatus{}

	/*Prize := make(map[string]interface{})
	err := json.Unmarshal([]byte(appointment.Prize), &Prize)
	tool.CheckError(err)
	if Prize["Begin"] != nil && Prize["End"] != nil {
		Begin := Prize["Begin"].(float64)
		End := Prize["End"].(float64)

		if End < Begin {
			as.Success = false
			as.Message = "抽奖低消费数据出错,最多金额要大于等于最少金额"
			return as
		}

	} else {
		as.Success = false
		as.Message = "抽奖低消费数据出错"
		return as
	}*/

	Link := make(map[string]interface{})
	err := json.Unmarshal([]byte(appointment.Link), &Link)
	tool.CheckError(err)
	if Link["Show"] != nil {
		Show := Link["Show"].(bool)
		if Show {
			if Link["Name"] == nil || Link["Url"] == nil {
				as.Success = false
				as.Message = "设置链接时，链接名和地址不能为空。"
				return as
			}
		}
	}

	UseTime := make(map[string]interface{})
	err = json.Unmarshal([]byte(appointment.UseTime), &UseTime)
	tool.CheckError(err)
	if UseTime["Show"] != nil {
		Show := UseTime["Show"].(bool)
		if Show {
			if UseTime["Week"] == nil || UseTime["Begin"] == nil || UseTime["End"] == nil {
				as.Success = false
				as.Message = "设置预定时间时，请选择时间。"
				return as
			} else {
				Begin, _ := strconv.ParseUint(UseTime["Begin"].(string), 10, 64)
				End, _ := strconv.ParseUint(UseTime["End"].(string), 10, 64)
				if End < Begin {
					as.Success = false
					as.Message = "预定时间,时间范围不正确。"
					return as
				}
			}
		}
	}

	/*if appointment.IsPost {
		if appointment.IsPayment == false {
			as.Success = false
			as.Message = "选择邮寄时，必须选择在线支付。"
			return as
		}
	}*/
	/*if appointment.Invite < 0 {
		as.Success = false
		as.Message = "邀请奖励不能负值"
		return as
	}*/
	if appointment.Stock <= 0 {
		as.Success = false
		as.Message = "库存必须大于零"
		return as
	}

	if appointment.Orig <= 0 {
		as.Success = false
		as.Message = "原价必须大于零"
		return as
	}
	if appointment.Price <= 0 {
		as.Success = false
		as.Message = "现价必须大于零"
		return as
	}
	if appointment.Price > appointment.Orig {
		as.Success = false
		as.Message = "现价必须小于等于原价"
		return as
	}
	if strings.EqualFold(appointment.Name, "") {
		as.Success = false
		as.Message = "项目名称不能为空"
		return as
	}

	if appointment.ID == 0 {
		have := self.FindAppointmentByName(Orm, appointment.Name)
		if have.ID != 0 {
			as.Success = false
			as.Message = "项目名称已经存在，请更换项目名称。"
			return as
		}
		err = self.Add(Orm, appointment)
		as.Message = "项目添加成功"
	} else {
		have := self.FindAppointmentByName(Orm, appointment.Name)
		if have.ID != 0 && have.ID != appointment.ID {
			as.Success = false
			as.Message = "项目名称已经存在，请更换项目名称。"
			return as
		}
		err = self.ChangeModel(Orm, appointment.ID, appointment)
		as.Message = "项目修改成功"
	}

	if err != nil {
		as.Success = false
		as.Message = err.Error()
	} else {
		as.Success = true

	}

	return as

}
