package scrapper

import (
	"encoding/json"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
)


type Y2MateRes struct{
	Status    string `json:"status"`
	Mess      string `json:"mess"`
	Page      string `json:"page"`
	Vid       string `json:"vid"`
	Extractor string `json:"extractor"`
	Title     string `json:"title"`
	T         int    `json:"t"`
	A         string `json:"a"`
	Links     struct{
		Mp4 map[string]struct{
			Size  string `json:"size"`
			F     string `json:"f"`
			Q     string `json:"q"`
			QText string `json:"q_text"`
			K     string `json:"k"`
		}
		Mp3 map[string]struct{
			Size  string `json:"size"`
			F     string `json:"f"`
			Q     string `json:"q"`
			QText string `json:"q_text"`
			K     string `json:"k"`
		}
	} `json:"links"`
}

type Y2MateCVRes struct {
	Status   string `json:"status"`
	Mess     string `json:"mess"`
	CStatus  string `json:"c_status"`
	Vid      string `json:"vid"`
	Title    string `json:"title"`
	Ftype    string `json:"ftype"`
	Fquality string `json:"fquality"`
	Dlink    string `json:"dlink"`
}

func Y2Mate(url string) (result *Y2MateRes, err error){
	client := resty.New()
	res, err := client.R().
				SetFormData(map[string]string{
					"k_query": url,
					"k_page": "home",
					"hl": "en",
					"q_auto": "0",
				}).
				SetHeader("Origin", "https://www.y2mate.com").
				SetHeader("Referer", "https://www.y2mate.com").
				SetHeader("USer-Agent", browser.MacOSX()).
				Post("https://www.y2mate.com/mates/analyzeV2/ajax")
	if err != nil {
		return nil, err
	}

	defer res.RawBody().Close()

	var R Y2MateRes

	err = json.Unmarshal(res.Body(), &R)
	
	if err != nil {
		return nil, err
	}
	
	return &R, nil
}

func Y2MateDownloadUrl(vid string, key string) (result *Y2MateCVRes, err error){
	client := resty.New()
	res, err := client.R().
				SetFormData(map[string]string{
					"vid": vid,
					"k": key,
				}).
				SetHeader("Origin", "https://www.y2mate.com").
				SetHeader("Referer", "https://www.y2mate.com/youtube/"+vid).
				SetHeader("USer-Agent", browser.MacOSX()).
				Post("https://www.y2mate.com/mates/convertV2/index")
	if err != nil {
		return nil, err
	}
	defer res.RawBody().Close()

	var R Y2MateCVRes

	err = json.Unmarshal(res.Body(), &R)
	if err != nil {
		return nil, err
	}

	return &R, nil
}
