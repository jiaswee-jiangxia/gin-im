package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"goskeleton/app/model"
	_ "goskeleton/bootstrap"
	"io/ioutil"
	"os"
	"strings"
)

type Translation struct {
	model.BaseModel
	Lang string
	Code string
	Msg  string
}

// 执行翻译小工具 (go run translation/translation.go)
func main() {
	langs := []string{"en", "zh-CN", "zh-TW", "vi", "th", "kr", "jp"}
	content, err := ioutil.ReadFile("translation/source/response.json")
	if err != nil {
		fmt.Println("failed open")
		return
	}
	var text map[string]string
	err = json.Unmarshal(content, &text)
	var fContent string

	model.Setup()
	db := model.GetDB()
	for code, item := range text {
		errText := strings.ReplaceAll(item, "_", " ")
		humanText := cases.Title(language.Und, cases.NoLower).String(errText)
		errText = strings.ReplaceAll(humanText, " ", "")
		fContent += errText + " = \"" + item + "\"\n\t"

		for _, lang := range langs {
			count := db.Table("translations").
				Where("code = ?", code).
				Where("lang = ?", lang).
				RowsAffected

			if count == 0 {
				err := db.Table("translations").
					Create(&Translation{
						Lang: lang,
						Code: item,
						Msg:  humanText,
					}).Error

				if err != nil {
					fmt.Println("error:")
					fmt.Println(err)
					return
				}
			}
		}
	}
	f, err := os.Create("app/global/response/response.go")
	if err != nil {
		fmt.Println("error file:")
		fmt.Println(err)
		return
	}
	defer f.Close()
	_, err = f.WriteString("package consts\n\n" +
		"// 此文件不可被直接篡改，请运行翻译小工具，请参考 (translation/translation.go)\n" +
		"// 任何于此文件的直接篡改，可能会被小工具覆盖而使到系统崩溃。\n" +
		"var (\n\t" +
		fContent +
		")\n")
}
