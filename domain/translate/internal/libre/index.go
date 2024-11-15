package libre

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func New() *Index {
	return &Index{}
}

type Index struct{}

func (m *Index) Translate(query []string, from, to string) (map[int]string, error) {
	translateTexts, err := m.translateBase(query, from, to)
	if err != nil {
		return nil, err
	}
	translateMap := make(map[int]string)
	for i := range translateTexts {
		translateMap[i] = translateTexts[i]
	}
	return translateMap, nil
}
func (m *Index) translateBase(query []string, from, to string) ([]string, error) {
	//###
	//POST http://translate.app.usokay.com/translate
	//Content-Type: application/json
	//
	//{
	//	"q": "name",
	//	"source": "en",
	//	"target": "zh",
	//	"format": "text",
	//	"alternatives": 0,
	//	"api_key": "ba07e09c-6e8c-4c1f-b3e0-88091934d51f"
	//}

	postParams := make(map[string]any)
	//q := url.QueryEscape(query)
	postParams["q"] = query
	postParams["source"] = from
	postParams["target"] = to
	postParams["format"] = "text"
	postParams["alternatives"] = 0
	//postParams["api_key"] = "ba07e09c-6e8c-4c1f-b3e0-88091934d51f"

	postParamsBytes, err := json.Marshal(&postParams)
	if err != nil {
		return nil, err
	}
	response, err := http.Post("http://translate.app.usokay.com/translate", "application/json", bytes.NewReader(postParamsBytes))
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var r struct {
		Error          string
		TranslatedText []string
	}
	if err = json.Unmarshal(body, &r); err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, errors.New(r.Error)
	}
	return r.TranslatedText, nil
}
