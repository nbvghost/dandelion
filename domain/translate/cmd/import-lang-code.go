package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/nbvghost/tool/object"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
)

type Lang struct {
	Code   string
	EnName string
	ZhName string
}

func readBaidu() map[string]Lang {
	langs := make(map[string]Lang)

	file, err := os.ReadFile("domain/translate/baidu.xml")
	if err != nil {
		return langs
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		log.Println(err)
		return langs
	}

	doc.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}

		zhName1 := strings.TrimSpace(selection.Children().Eq(0).Text())
		code1 := strings.TrimSpace(selection.Children().Eq(1).Text())

		if len(code1) == 0 {
			return
		}
		langs[code1] = Lang{
			Code:   code1,
			EnName: "",
			ZhName: zhName1,
		}

		zhName2 := strings.TrimSpace(selection.Children().Eq(2).Text())
		code2 := strings.TrimSpace(selection.Children().Eq(3).Text())

		if len(code2) == 0 {
			return
		}
		langs[code2] = Lang{
			Code:   code2,
			EnName: "",
			ZhName: zhName2,
		}

		zhName3 := strings.TrimSpace(selection.Children().Eq(4).Text())
		code3 := strings.TrimSpace(selection.Children().Eq(5).Text())
		if len(code3) == 0 {
			return
		}
		langs[code1] = Lang{
			Code:   code3,
			EnName: "",
			ZhName: zhName3,
		}

	})
	return langs
}

func readAliyun() map[string]Lang {
	langs := make(map[string]Lang)

	file, err := os.ReadFile("domain/translate/aliyun.xml")
	if err != nil {
		return langs
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		log.Println(err)
		return langs
	}

	var ignoreCode = []string{
		"eo", "tw", "lt", "mh", "os", "et", "umb", "el",
		"to", "sq", "sl", "ti", "cv", "sm", "om", "fj",
		"xh", "tk", "pa", "rw", "sk", "yue", "tt", "sw",
		"gl", "mk", "hy", "kg", "bi", "ka", "wo", "zu",
		"ts", "ab", "ay", "ty", "sn", "cy", "kl", "lv", "ht", "bs"}

	doc.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			return
		}

		children := selection.Children()
		l := children.Length()

		code := strings.TrimSpace(selection.Children().Eq(l - 1).Text())
		enName := strings.TrimSpace(selection.Children().Eq(l - 2).Text())
		zhName := strings.TrimSpace(selection.Children().Eq(l - 3).Text())

		if len(code) == 0 {
			return
		}

		if lo.IndexOf[string](ignoreCode, code) > -1 {
			return
		}

		langs[code] = Lang{
			Code:   code,
			EnName: enName,
			ZhName: zhName,
		}
	})
	return langs
}

func readVolcengine() map[string]Lang {
	langs := make(map[string]Lang)

	file, err := os.ReadFile("domain/translate/volcengine.xml")
	if err != nil {
		return langs
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(file))
	if err != nil {
		log.Println(err)
		return langs
	}

	doc.Find("table tbody tr").Each(func(i int, selection *goquery.Selection) {

		code := strings.TrimSpace(selection.Children().Eq(0).Text())
		zhName := strings.TrimSpace(selection.Children().Eq(1).Text())
		enName := strings.TrimSpace(selection.Children().Eq(2).Text())

		if len(code) == 0 {
			return
		}

		langs[code] = Lang{
			Code:   code,
			EnName: enName,
			ZhName: zhName,
		}
	})
	return langs
}

func main() {

	ableMap := make(map[string]Lang)

	aLangs := readAliyun()
	vLangs := readVolcengine()

	for k := range aLangs {
		_, ok := vLangs[k]
		if ok {
			ableMap[k] = aLangs[k]
		}
	}

	///log.Println(ableMap)

	for k := range ableMap {
		id := getUint(k)
		item := ableMap[k]
		println(fmt.Sprintf(`INSERT INTO "Language" ("ID", "Code", "Name", "ChineseName") VALUES ('%d', '%s', '%s', '%s');`, id, k, item.EnName, item.ZhName))
		//log.Println(aLangs[k])
	}

	log.Println("-----------------------------------------------------------")
	/*baiduData:=readBaidu()
	for k := range ableMap {
		var has = false
		for kk := range baiduData {
			if ableMap[k].ZhName==baiduData[kk].ZhName{
				has=true
				break
			}
		}
		if has{
			log.Println(ableMap[k])
		}
	}*/

}

func getUint(st string) uint {
	vs := []byte(st)
	var p = ""
	for i := range vs {
		p += fmt.Sprintf("%d", vs[i])
	}
	return object.ParseUint(p)
}
