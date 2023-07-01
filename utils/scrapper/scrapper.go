//from https://github.com/itzngga/Lara/blob/master/util/scrapper/scrapper.go

package scrapper

import (
	"botwa/utils"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
	"github.com/go-resty/resty/v2"

	"regexp"
	"time"
)

func NewCloudflareBypass() *resty.Client {
	httpTransport := &http.Client{}
	httpTransport.Transport = cloudflarebp.AddCloudFlareByPass(httpTransport.Transport)

	client := resty.New()
	client.SetTransport(httpTransport.Transport)

	return client
}

//type ScrapperResponse interface {
//	SnapInstaResponse | SnaptikResponse
//}
//
//func Scrape[T ScrapperResponse](link string) (T, error) {
//	scrapper := Scrapper{}
//	if strings.Contains(link, "instagram") {
//		result, err := scrapper.GetSnapInsta(link)
//
//		return result, err
//	}
//}

var innerHtml = regexp.MustCompile("\\.innerHTML = \"(.*?)\";")
//var innerToken = regexp.MustCompile("get_progressApi\\('/render\\.php\\?token=(.*?)'\\);")

func chip(d string, e int, f int) string {
	g := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+/"
	h := g[:e]
	i := g[:f]
	j := 0
	for c := len(d) - 1; c >= 0; c-- {
		b := string(d[c])
		index := strings.Index(h, b)
		if index != -1 {
			j += index * intPow(e, len(d)-1-c)
		}
	}
	k := ""
	for j > 0 {
		k = string(i[j%f]) + k
		j = (j - (j % f)) / f
	}
	if k == "" {
		return "0"
	}
	return k
}

func DecodeSnap(h string, u int, n string, t int, e int, r int) string {
	decodedString := ""
	for i := 0; i < len(h); i++ {
		s := ""
		for h[i:i+1] != n[e:e+1] {
			s += h[i : i+1]
			i++
		}
		for j := 0; j < len(n); j++ {
			s = strings.ReplaceAll(s, n[j:j+1], strconv.Itoa(j))
		}
		chipResult, _ := strconv.Atoi(chip(s, e, 10))
		decodedString += string(rune(chipResult - t))
	}
	return decodedString
}


func intPow(base, exponent int) int {
	result := 1
	for i := 0; i < exponent; i++ {
		result *= base
	}
	return result
}

func TimeElapsed(name string) func() {
	start := time.Now().In(time.Local)
	return func() {
		fmt.Printf("[%s] took %s\n", name, utils.HumanizeDuration(time.Since(start)))
	}
}