package utils

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func Translate(messageID string, params map[string]interface{}, c *fiber.Ctx) string {
	var translate string
	var err error

	localizeConfig := &i18n.LocalizeConfig{
		MessageID: messageID,
	}
	if params != nil {
		localizeConfig.TemplateData = params
	}

	translate, err = fiberi18n.Localize(c, localizeConfig)
	if err != nil {
		return messageID
	}
	return translate
}
