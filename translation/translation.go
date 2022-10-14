package translation

import (
	"encoding/json"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"goskeleton/app/model"
	"log"
)

// Bundle var
var Bundle *i18n.Bundle
var Loc *i18n.Localizer

// Localizer struct
type Localizer struct {
	Localizer *i18n.Localizer
	Language  string
}

// Setup func
func Setup() {
	Bundle = i18n.NewBundle(language.English)
	Bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	langs, err := model.GetLanguageList()
	if err != nil {
		log.Fatalf("translation.Setup [GetLanguageList] err: %v", err)
	}

	for _, lang := range langs {
		//trans, err := model.GetTranslationByLocale("api", lang.Locale)
		//if err != nil {
		//	log.Fatalf("translation.Setup [GetTranslationByLocale] err: %v", err)
		//}
		//
		//var tranStr string
		//
		//for _, tran := range trans {
		//	tranStr = tranStr + tran.Name + " = " + "\"" + tran.Value + "\"\n"
		//}
		Bundle.MustLoadMessageFile("translation/source/" + lang.Locale + ".json")
	}
}

func SetNewLocalizer(locale string) {
	Loc = i18n.NewLocalizer(Bundle, locale)
}

func Trans(text string) string {
	translation, err := Loc.Localize(&i18n.LocalizeConfig{
		MessageID: text,
	})
	if err != nil {
		return text
	}
	return translation
}
