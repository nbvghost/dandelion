package area

import (
	"bufio"
	"io"
	"os"
	"strings"

	"gorm.io/gorm"

	"github.com/nbvghost/dandelion/entity"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/dao"
	"github.com/nbvghost/dandelion/library/singleton"
	"github.com/nbvghost/tool/object"
)

func LoadArea() error {
	//D:\projects\rent\server
	f, err := os.Open("area_code_2022.csv")
	if err != nil {
		return err
	}
	read := bufio.NewReader(f)

	err = dao.DeleteBy(singleton.Orm().Session(&gorm.Session{AllowGlobalUpdate: true}), entity.Area, nil)
	if err != nil {
		return err
	}

	for {
		l, _, err := read.ReadLine()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		line := string(l)

		fields := strings.Split(line, ",")

		area := model.Area{
			Code:  object.ParseUint(fields[0]),
			Name:  fields[1],
			Level: object.ParseUint(fields[2]),
			PCode: object.ParseUint(fields[3]),
		}
		err = dao.Create(singleton.Orm(), &area)
		if err != nil {
			return err
		}
	}
}
