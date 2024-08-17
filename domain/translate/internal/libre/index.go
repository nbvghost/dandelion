package libre


/*
	func (m *Html) Translate(query, from, to string) (string, error) {
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
		postParams["api_key"] = m.ApiKey

		postParamsBytes, err := json.Marshal(&postParams)
		if err != nil {
			return "", err
		}
		response, err := http.Post("http://translate.app.usokay.com/translate", "application/json", bytes.NewReader(postParamsBytes))
		if err != nil {
			return "", err
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()

		var r _Result
		if err = json.Unmarshal(body, &r); err != nil {
			return "", err
		}
		if response.StatusCode != 200 {
			return "", errors.New(r.Error)
		}
		return r.TranslatedText, nil
	}
*/