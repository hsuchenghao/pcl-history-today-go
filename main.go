package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type Response struct {
	TenacityResBody struct {
		List []struct {
			Title   string `json:"title"`
			Content string `json:"content"`
		}
	} `json:"tenacityapi_res_body"`
}

func createXaml(data Response) string {
	var cards []string

	for _, item := range data.TenacityResBody.List {
		title := regexp.MustCompile(`\d{4}年\d{1,2}月\d{1,2}日 `).ReplaceAllString(item.Title, "")
		card := `<local:MyCard Title="`+title+`" Margin="0,0,0,15" CanSwap="True" IsSwaped="True">` +
			`<StackPanel Margin="25,40,23,15">` +
			`<TextBlock TextWrapping="Wrap" Margin="0,0,0,4">`+item.Content+`</TextBlock>` +
			`</StackPanel></local:MyCard>`
		cards = append(cards, card)
	}

	return strings.Join(cards, "\n")
}

func main() {
	url := "https://api.tenacity.dev/212"
	params := "?api_appid=your_appid&api_sign=your_sign&needContent=1"

	resp, err := http.Get(url + params)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data Response
	json.Unmarshal(body, &data)

	xamlString := createXaml(data)

	filepath := "/www/wwwroot/history-today.lat/custom.xaml"
	err = ioutil.WriteFile(filepath, []byte(xamlString), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("脚本运行成功")

	// Script 2
	now := time.Now()
	versionString := "ver " + now.Format("06.01.02")

	iniFilepath := "/www/wwwroot/history-today.lat/custom.xaml.ini"
	err = ioutil.WriteFile(iniFilepath, []byte(versionString), 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("INI文件生成成功！")
}
