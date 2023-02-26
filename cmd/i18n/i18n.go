package main

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"os"
	"path"
)

func main() {
	wd, _ := os.Getwd()

	bundle := &i18n.Bundle{}
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.MustLoadMessageFile(path.Join(wd, "locales", "en.json"))

	loc := i18n.NewLocalizer(bundle, "en")
	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID:   "translateNotFound",
		PluralCount: 2,
	})
	fmt.Println(translation)
}
