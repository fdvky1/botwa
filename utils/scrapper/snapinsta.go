package scrapper

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
)

type SnapInstaResponse struct {
	Username    string   `json:"username"`
	Avatar      string   `json:"avatar"`
	ResultMedia []string `json:"result_media"`
}

func GetSnapInsta(instagram string) (response *SnapInstaResponse, err error) {
	defer TimeElapsed("Scrap SnapInsta")()

	client := NewCloudflareBypass()
	resp, err := client.R().
		SetFormData(map[string]string{
			"url":    instagram,
			"action": "post",
			"lang":   "id",
		}).
		SetHeader("Origin", "https://snapinsta.app").
		SetHeader("Referer", "https://snapinsta.app/id").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snapinsta.app/action2.php")
	if err != nil {
		return nil, err
	}
	
	defer resp.RawBody().Close()
	script := resp.String()
	splited := strings.Split(script, "}(")
	if len(splited) <= 1 {
		return nil, errors.New("[404] Could not find executable script")
	}
	splited = strings.Split(strings.Split(splited[1], "))")[0], ",")
	h := strings.ReplaceAll(splited[0], "\"", "")
	u, _ := strconv.Atoi(splited[1])
	n := strings.ReplaceAll(splited[2], "\"", "")
	t, _ := strconv.Atoi(splited[3])
	e, _ := strconv.Atoi(splited[4])
	r, _ := strconv.Atoi(splited[5])

	dec := DecodeSnap(h, u, n, t, e, r)
	html := innerHtml.FindStringSubmatch(dec)[1]
	parsedHtml := strings.ReplaceAll(html, `\"`, "")

	document, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
	if err != nil {
		return response, err
	}

	response = &SnapInstaResponse{
		ResultMedia: make([]string, 0),
	}
	
	response.Username = document.Find("div.download-top > div").First().Text()
	avatar, _ := document.Find("div.download-top > div > img").First().Attr("src")
	response.Avatar = avatar

	document.Find("div.download-bottom > a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			response.ResultMedia = append(response.ResultMedia, href)
		}
	})

	return response, nil
}