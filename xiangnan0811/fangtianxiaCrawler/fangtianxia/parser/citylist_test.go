package parser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/axgle/mahonia"
)

func TestParseCityList(t *testing.T) {
	c, err := ioutil.ReadFile("citylist_test_data.html")
	if err != nil {
		fmt.Println(err)
	}
	bodyString := mahonia.NewDecoder("gbk").ConvertString(string(c))
	result := ParseCityList([]byte(bodyString))

	// verify result
	const resultSize = 1250
	expectedUrls := []string{
		"https://newhouse.fang.com/house/s",
		"https://esf.fang.com",
		"https://sh.newhouse.fang.com/house/s",
		"https://sh.esf.fang.com",
		"https://tj.newhouse.fang.com/house/s",
		"https://tj.esf.fang.com",
	}

	if realReqSize := len(result.Requests); realReqSize != resultSize {
		t.Errorf("result should have %d requests; but had %d", resultSize, realReqSize)
	}

	for i, url := range expectedUrls {
		if result.Requests[i].Url != url {
			t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
		}
	}
}
