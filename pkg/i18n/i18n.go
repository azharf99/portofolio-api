package i18n

import (
	"encoding/json"
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func Init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load files
	_, err := bundle.LoadMessageFile("locales/en.json")
	if err != nil {
		fmt.Printf("Warning: failed to load en.json: %v\n", err)
	}
	_, err = bundle.LoadMessageFile("locales/id.json")
	if err != nil {
		fmt.Printf("Warning: failed to load id.json: %v\n", err)
	}
	_, err = bundle.LoadMessageFile("locales/ru.json")
	if err != nil {
		fmt.Printf("Warning: failed to load ru.json: %v\n", err)
	}
}

func GetLocalizer(langs ...string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, langs...)
}

func T(localizer *i18n.Localizer, messageID string) string {
	translation, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return translation
}
