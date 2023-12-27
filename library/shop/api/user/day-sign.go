package user

import (
	"fmt"
	"github.com/nbvghost/dandelion/library/db"
	"log"
	"strconv"
	"time"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/result"
	"github.com/nbvghost/dandelion/service/journal"
	"github.com/nbvghost/dandelion/service/user"
	"github.com/pkg/errors"
)

type DaySign struct {
	UserService    user.UserService
	JournalService journal.JournalService
	User           *model.User `mapping:""`
}

func (m *DaySign) Handle(context constrain.IContext) (r constrain.IResult, err error) {

	userInfo := m.UserService.GetUserInfo(m.User.ID)

	now := userInfo.GetDaySignTime()
	today := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	//d,err:=time.ParseDuration("24h")
	//log.Println(err)

	fmt.Println(":", time.Now().Unix())
	fmt.Println(":", today.Unix())
	dayCount := float64(float64(time.Now().Unix()-today.Unix()) / 60 / 60 / 24) //天
	fmt.Println("天", dayCount)
	as := result.ActionResult{}
	if dayCount > 1 {
		//已经超过一天了，
		userInfo.SetDaySignTime(time.Now())
		userInfo.SetDaySignCount(1)
		as.Code = result.Success
		as.Message = "打卡成功，您的打卡已经超过一天了，打卡重新累计"
		//return &gweb.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("打卡成功，您的打卡已经超过一天了，打卡重新累计"), "OK", nil)}
	} else if dayCount <= 1 && dayCount >= 0 {
		//可以打卡
		userInfo.SetDaySignTime(time.Now())
		userInfo.SetDaySignCount(userInfo.GetDaySignCount() + 1)
		as.Code = result.Success
		as.Message = "打卡成功"
	} else {
		//负数
		//已经打过卡了
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("您今天已经打卡了"), "OK", nil)}, nil
	}

	if userInfo.GetDaySignCount() <= 0 {
		userInfo.SetDaySignCount(1)
	}

	//todo
	daySign := map[string]int{"1": 2}
	if len(daySign) > 0 {

		DaySign := daySign
		score, have := DaySign[strconv.Itoa(userInfo.GetDaySignCount())]
		if have {

		} else {
			score, have = DaySign["max"]
			if !have {
				log.Println("打卡data.json数据没有设置DaySign.max字段值")
			}

		}
		//err := m.JournalService.AddScoreJournal(db.Orm(), m.User.ID, "签到送积分", userInfo.DaySignTime.String()+"/"+strconv.Itoa(int(score))+"/"+strconv.Itoa(userInfo.DaySignCount), play.ScoreJournal_Type_DaySign, int64(score), extends.KV{Key: "UserInfoID", Value: userInfo.ID})
		err := m.JournalService.AddScoreJournal(db.Orm(), m.User.ID, "签到送积分", userInfo.GetDaySignTime().String()+"/"+strconv.Itoa(int(score))+"/"+strconv.Itoa(userInfo.GetDaySignCount()), model.ScoreJournal_Type_DaySign, int64(score))
		if err != nil {
			as.Code = result.Fail
			as.Message = err.Error()
		} else {
			err = userInfo.Update(db.Orm())
			if err != nil {
				return nil, err
			}
		}
		return &result.JsonResult{Data: &as}, nil

	} else {
		return &result.JsonResult{Data: (&result.ActionResult{}).SmartError(errors.New("暂时无法打卡"), "OK", nil)}, nil
	}
}
