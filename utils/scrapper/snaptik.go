package scrapper

import (
	"bytes"
	"errors"
	"strconv"
	"strings"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
)

type SnaptikResponse struct {
	Username    string   `json:"username"`
	Description string   `json:"description"`
	VideoUrl    []string `json:"video_url"`
	ImageUrl    []string `json:"image_url"`
}

func GetSnaptik(tiktok string) (response SnaptikResponse, err error) {
	defer TimeElapsed("Scrap Snaptik")()

	client := resty.New()
	resp, err := client.R().
		SetHeader("User-Agent", browser.Firefox()).
		Get("https://snaptik.app/ID")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()

	document, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return response, err
	}

	var param = map[string]string{
		"url": tiktok,
	}

	document.Find("form input").Each(func(index int, selector *goquery.Selection) {
		name, ok := selector.Attr("name")
		if ok && name == "lang" {
			value, _ := selector.Attr("value")
			param["lang"] = value
		}

		if ok && name == "token" {
			value, _ := selector.Attr("value")
			param["token"] = value
		}

		selector.Next()
	})

	resp, err = client.R().
		SetFormData(param).
		SetHeader("Origin", "https://snaptik.app").
		SetHeader("Referer", "https://snaptik.app/ID").
		SetHeader("User-Agent", browser.Firefox()).
		Post("https://snaptik.app/abc2.php")
	if err != nil {
		return response, err
	}

	defer resp.RawBody().Close()


	script := resp.String()
	splited := strings.Split(script, "}(")
	if len(splited) <= 1 {
		return response, errors.New("[404] Could not find executable script")
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

	document, err = goquery.NewDocumentFromReader(bytes.NewReader([]byte(parsedHtml)))
	if err != nil {
		return response, err
	}

	response = SnaptikResponse{}
	response.Username = document.Find("div.info > span").Text()
	response.Description = document.Find("div.info > div").Text()
	response.VideoUrl = make([]string, 0)
	response.ImageUrl = make([]string, 0)

	document.Find("div.video-links > a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if ok {
			if strings.Contains(href, "https://cdn.snaptik.app"){
				response.VideoUrl = append(response.VideoUrl, href)
			} else if strings.Contains(href, "/file.php?type=dl"){
				response.VideoUrl = append(response.VideoUrl, "https://snaptik.app/"+href)
			}
		}
		selection.Next()
	})

	document.Find("div.dl-footer").Each(func(i int, selection *goquery.Selection) {
		src, ok := selection.Find("a").Attr("href")
		if ok && strings.Contains(src, "https://tikcdn.net") {
			response.ImageUrl = append(response.ImageUrl, src)
		}
	})

	return response, nil
}