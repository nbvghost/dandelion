package express

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/nbvghost/dandelion/app/service/dao"
	"github.com/nbvghost/gweb/tool/encryption"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/nbvghost/glog"
)

type ExpressTemplateService struct {
	dao.BaseDao
}

func (b ExpressTemplateService) GetExpressInfo(OrdersID uint64, LogisticCode, ShipperName string) map[string]interface{} {

	shipperMap := make(map[string]string)
	shipperMap["中国邮政"] = "YZPY"
	shipperMap["EMS"] = "EMS"
	shipperMap["顺丰快递"] = "SF"
	shipperMap["中通快递"] = "ZTO"
	shipperMap["圆通快递"] = "YTO"
	shipperMap["申通快递"] = "STO"
	shipperMap["韵达快递"] = "YD"
	shipperMap["百世汇通"] = "HTKY"
	shipperMap["天天快递"] = "HHTT"
	shipperMap["国通快递"] = "GTO"
	shipperMap["宅急送"] = "ZJS"

	ShipperNameCode := shipperMap[ShipperName]

	requestData := "{'OrderCode':'" + strconv.Itoa(int(OrdersID)) + "','ShipperCode':'" + ShipperNameCode + "','LogisticCode':'" + LogisticCode + "'}"

	DataSign := base64.StdEncoding.EncodeToString([]byte(strings.ToLower(encryption.Md5ByString(requestData + "8d8ef028-000f-4f3e-8475-bc90d5772002"))))

	postData := url.Values{}

	postData.Add("RequestData", url.PathEscape(requestData))
	postData.Add("EBusinessID", "1334134")
	postData.Add("RequestType", "1002")
	postData.Add("DataType", "2")
	postData.Add("DataSign", url.PathEscape(string(DataSign)))

	//fmt.Println(postData.Encode())

	result := make(map[string]interface{})

	resp, err := http.PostForm("http://api.kdniao.cc/Ebusiness/EbusinessOrderHandle.aspx", postData)
	if err != nil {
		return result
	}
	defer resp.Body.Close()
	//resp, err := http.PostForm("http://sandboxapi.kdniao.cc:8080/kdniaosandbox/gateway/exterfaceInvoke.json", postData)

	bsdfsd, errs := ioutil.ReadAll(resp.Body)
	glog.Error(errs)
	json.Unmarshal(bsdfsd, &result)

	result["ShipperName"] = ShipperName
	result["ShipperCode"] = ShipperNameCode
	result["ShipperNo"] = LogisticCode

	//fmt.Println(errs)
	//fmt.Println(result)
	return result
}
func (b ExpressTemplateService) GetExpressTemplateByName(Name string) dao.ExpressTemplate {
	Orm := dao.Orm()
	var list dao.ExpressTemplate
	Orm.Model(&dao.ExpressTemplate{}).Where("Name=?", Name).Find(&list)
	return list
}
func (b ExpressTemplateService) GetExpressTemplateByOID(OID uint64) dao.ExpressTemplate {
	Orm := dao.Orm()
	var list dao.ExpressTemplate
	Orm.Model(&dao.ExpressTemplate{}).Where("OID=?", OID).Find(&list)
	return list
}
func (b ExpressTemplateService) SaveExpressTemplate(target *dao.ExpressTemplate) error {
	Orm := dao.Orm()
	have := b.GetExpressTemplateByName(target.Name)
	if have.ID == 0 {
		return b.Save(Orm, target)
	} else {
		if have.ID == target.ID {
			return b.Save(Orm, target)
		} else {
			return errors.New("名称已经存在")
		}

	}
}
