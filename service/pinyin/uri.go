package pinyin

import (
	"fmt"
	"github.com/nbvghost/dandelion/service/cache"
	"regexp"
	"strings"
)

//[`~!@#$%^&*()_\-+=<>?:"{}|,.\/;'\\[\]·~！@#￥%……&*（）——\-+={}|《》？：“”【】、；‘’，。、]
var markRegexp = regexp.MustCompile(`[~!@#$%^&*()_\-+=<>?:"{},.;\|\/\\'\[\]·~！@#￥%……&*（）——\-+={}\|《》？：“”【】、；‘’，。、\s\x{0060}]+`)

//var markRegexp = regexp.MustCompile(`[[:punct:]]+`)

//var enRegexp = regexp.MustCompile("[a-zA-Z\u4E00-\u9FA5\u9FA6-\u9FFF\u3400-\u4DBF\u20000-\u2A6DF\u2A700-\u2B738\u2B740-\u2B81D\u2B820-\u2CEA1\u2CEB0-\u2EBE0\u30000-\u3134A\u2F00-\u2FD5\u2E80-\u2EF3\uF900-\uFAD9\u2F800-\u2FA1D\uE815-\uE86F\uE400-\uE5E8\uE600-\uE6CF\u31C0-\u31E3\u2FF0-\u2FFB\u3105-\u312F\u31A0-\u31BA\u3007]+")
var enRegexp = regexp.MustCompile("[a-zA-Z0-9]")
var gRegexp = regexp.MustCompile("-+")

//var enRegexp = regexp.MustCompile("[a-zA-Z\u4E00-\u9FA5\u9FA6-\u9FFF]+")

func (Service) AutoDetectUri(s string) string {

	var ll []string

	for _, v := range []rune(s) {
		word := string(v)
		if !markRegexp.MatchString(word) {
			if enRegexp.MatchString(word) {
				ll = append(ll, word)
			} else {
				py := cache.Cache.GetPinyin(word)
				if len(py) == 0 {
					ll = append(ll, fmt.Sprintf("-%x-", word[:]))
				} else {
					ll = append(ll, fmt.Sprintf("-%s-", py))
				}

			}
		} else {
			ll = append(ll, "-")
		}
	}

	txt := strings.Join(ll, "")
	txt = gRegexp.ReplaceAllString(txt, "-")
	return strings.ToLower(strings.Trim(txt, "-"))
}
