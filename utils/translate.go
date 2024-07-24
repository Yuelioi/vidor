package utils

import (
	"context"
	"fmt"

	"github.com/michimani/deepl-sdk-go"
	"github.com/michimani/deepl-sdk-go/params"
	"github.com/michimani/deepl-sdk-go/types"
)

func translate() {
	client, err := deepl.NewClient()
	// client.AuthenticationKey = ""
	if err != nil {
		fmt.Println(err)
		return
	}

	text := []string{
		"Привіт",
		"Це приклад тексту",
	}
	params := &params.TranslateTextParams{
		TargetLang: types.TargetLangEN,
		Text:       text,
	}

	res, errRes, err := client.TranslateText(context.TODO(), params)

	if err != nil {
		fmt.Println(err)
	}

	if errRes != nil {
		fmt.Println("ErrorResponse", errRes.Message)
	}

	for i := range res.Translations {
		fmt.Printf("%s -> %s\n", text[i], res.Translations[i].Text)
	}
}
