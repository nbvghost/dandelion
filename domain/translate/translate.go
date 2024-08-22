package translate

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/nbvghost/dandelion/domain/translate/internal"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"html/template"
	"regexp"
	"strings"
	"sync"

	"github.com/nbvghost/dandelion/constrain"
	"github.com/nbvghost/dandelion/entity/model"
	"github.com/nbvghost/dandelion/library/contexext"
	"github.com/nbvghost/dandelion/library/db"
)

var varRegexp = regexp.MustCompile(`\$\{[\S\x20]+}`)
var notTranslateRegexp = regexp.MustCompile(`(?i)^(([\d+-.$])|([\d|+|-|.|$|/|\\])|(kg)|([:\s]))+$`)

func CheckNotTranslate(text string) bool {
	return notTranslateRegexp.MatchString(text)
}

type nodeID string

func newNodeID(seqID int) nodeID {
	return nodeID(fmt.Sprintf("node%d", seqID))
}

type Html struct {
	translate    internal.Translate
	LanguageCode map[string]string
	sync.RWMutex
}

func (m *Html) html(node *html.Node) ([]byte, error) {
	var buf bytes.Buffer
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		err := html.Render(&buf, c)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil

}

type translateInfo struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

func (m *Html) Translate(query []string, from, to string) (map[int]string, error) {
	return m.translate.Translate(query, from, to)
}
func (m *Html) TranslateHtml(context constrain.IContext, docBytes []byte) ([]byte, error) {
	contextValue := contexext.FromContext(context)
	var err error
	var node *html.Node
	node, err = html.Parse(bytes.NewBuffer(docBytes))
	if err != nil {
		return nil, err
	}

	//todo 当前暂定en为内容的默认语言,后面走数据库
	if strings.EqualFold(contextValue.Lang, "en") {
		return docBytes, nil
	}

	if _, ok := m.LanguageCode[contextValue.Lang]; !ok {
		//当前语言不在翻译列表中
		return docBytes, nil
	}

	var seqID int
	setMap := map[string]nodeID{}
	willTranslateTexts := map[string]translateInfo{}
	//var willTranslateLocker sync.RWMutex

	//把翻译的文本替换成变量占位符，并把要翻译的放在willTranslateTexts里
	extractFunc := func(text string) string {
		texts := strings.Split(text, "\n")
		for i := range texts {
			if varRegexp.MatchString(texts[i]) {
				continue
			}
			if CheckNotTranslate(texts[i]) {
				continue
			}
			v := strings.TrimSpace(texts[i])
			var nid nodeID
			var ok bool
			if nid, ok = setMap[v]; !ok {
				seqID++
				nid = newNodeID(seqID)
				setMap[v] = nid
				willTranslateTexts[string(nid)] = translateInfo{
					Src: v,
				}
			}
			texts[i] = fmt.Sprintf("{{.%s.Dst}}", string(nid))
		}
		return strings.Join(texts, "\n")
	}

	//提取要翻译的文字
	var f func(*html.Node)
	f = func(n *html.Node) {
		var noTranslate bool
		for _, v := range n.Attr {
			if strings.EqualFold(v.Key, "no-translate") {
				noTranslate = true
				break
			}
		}
		if noTranslate {
			return
		}

		if n.Type == html.TextNode && !strings.EqualFold(n.Parent.Data, "style") && !strings.EqualFold(n.Parent.Data, "script") {
			text := strings.TrimSpace(n.Data)
			if len(text) > 0 {
				n.Data = extractFunc(text)
			}
		}
		if strings.EqualFold(n.Data, "input") {
			for i, v := range n.Attr {
				if strings.EqualFold(v.Key, "placeholder") {
					n.Attr[i].Val = extractFunc(n.Attr[i].Val)
					break
				}
			}
		}
		if strings.EqualFold(n.Data, "html") && n.Type == html.ElementNode {
			if len(n.Attr) > 0 {
				var hasLangAttr bool
				for i, v := range n.Attr {
					if strings.EqualFold(v.Key, "lang") {
						n.Attr[i].Val = contextValue.Lang
						hasLangAttr = true
						break
					}
				}
				if !hasLangAttr {
					n.Attr = append(n.Attr, html.Attribute{
						Key: "lang",
						Val: contextValue.Lang,
					})
				}
			} else {
				n.Attr = append(n.Attr, html.Attribute{
					Key: "lang",
					Val: contextValue.Lang,
				})
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	tx := db.Orm().Begin(&sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	defer func() {
		if err != nil || recover() != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var translateList []string

	//TODO 缓存在cache server 里，或者放在redis
	var translateModelList []model.Translate
	tx.Model(model.Translate{}).Where(`"TextType"=? and "LangType"=?`, "en", contextValue.Lang).Find(&translateModelList)

	for k := range willTranslateTexts {
		var has bool
		v := willTranslateTexts[k]
		//todo 优化二分查找或者其它方法
		for _, e := range translateModelList {
			if strings.EqualFold(v.Src, e.Text) {
				v.Dst = e.LangText
				willTranslateTexts[k] = v
				has = true
				break
			}
		}
		if !has {
			translateList = append(translateList, v.Src)
		}
	}

	if len(translateList) > 0 {
		var translateMap map[int]string
		translateMap, err = m.translate.Translate(translateList, "en", contextValue.Lang)
		if err != nil {
			context.Logger().Error("translate error", zap.Error(err))
			return nil, err
		}

		for index := range translateList {
			v := translateList[index]
			translateText := translateMap[index]
			for k, vv := range willTranslateTexts {
				if vv.Src == v {
					vv.Dst = translateText
					willTranslateTexts[k] = vv
					break
				}
			}
			if len(translateText) == 0 {
				translateText = v
			}
			if err = tx.Model(model.Translate{}).Create(&model.Translate{
				Text:     v,
				TextType: "en",
				LangType: contextValue.Lang,
				LangText: translateText,
			}).Error; err != nil {
				context.Logger().Error("write db.Translate error", zap.Error(err))
				return nil, err
			}
		}
	}

	/*//var translateText string
	for _, v := range translateList {
		if len(v) > 0 {

			var translateText string
			translateText, err = m.Translate(v, "en", contextValue.Lang)
			if err != nil {
				context.Logger().Error("translate error", zap.Error(err))
				return nil, err
			}

			//willTranslateLocker.Lock()
			for k, vv := range willTranslateTexts {
				if vv.Src == v {
					vv.Dst = translateText
					willTranslateTexts[k] = vv
					break
				}
			}
			//willTranslateLocker.Unlock()

			if len(translateText) == 0 {
				translateText = v
			}

			if err = tx.Model(model.Translate{}).Create(&model.Translate{
				Text:     v,
				TextType: "en",
				LangType: contextValue.Lang,
				LangText: translateText,
			}).Error; err != nil {
				context.Logger().Error("write db.Translate error", zap.Error(err))
				return nil, err
			}
		}
	}*/

	docBytes, err = m.html(node)
	if err != nil {
		return nil, err
	}

	var t *template.Template
	t, err = template.New("default").Parse(string(docBytes))
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(nil)
	err = t.Execute(buffer, willTranslateTexts)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func NewTranslate(languageCode map[string]string) (*Html, error) {
	translate, err := internal.New()
	if err != nil {
		return nil, err
	}
	return &Html{LanguageCode: languageCode, translate: translate}, nil
}
