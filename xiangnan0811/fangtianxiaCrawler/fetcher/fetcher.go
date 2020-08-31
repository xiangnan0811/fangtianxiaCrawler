package fetcher

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/axgle/mahonia"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"
)

var rateLimiter = time.Tick(config.Qps * time.Millisecond)

func Fetch(url string) ([]byte, error) {
	<-rateLimiter

	// Add cookie
	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	// Set timeout
	client := &http.Client{
		Jar: jar,
	}
	request, _ := http.NewRequest(http.MethodGet, url, nil)

	rand.Seed(time.Now().Unix())
	request.Header.Add("User-Agent", config.UserAgentList[rand.Intn(len(config.UserAgentList))])
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
	}
	bodyString := mahonia.NewDecoder("gbk").ConvertString(string(body))

	return []byte(bodyString), nil
}
