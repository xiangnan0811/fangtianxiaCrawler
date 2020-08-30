package fetcher

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"time"

	"golang.org/x/net/publicsuffix"

	"github.com/xiangnan0811/fangtianxiaCrawler_distributed/config"

	"golang.org/x/text/encoding/unicode"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
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
		Timeout: 20 * time.Second,
		Jar:     jar,
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

	bodyReader := bufio.NewReader(response.Body)
	e := MyDetermineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())

	return ioutil.ReadAll(utf8Reader)
}

func MyDetermineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
