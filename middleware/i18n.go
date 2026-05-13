package middleware

import (
	"github.com/azharf99/portofolio-api/pkg/i18n"
	"github.com/gin-gonic/gin"
)

func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Priority: Query param "lang", then "Accept-Language" header
		lang := c.Query("lang")
		if lang == "" {
			lang = c.GetHeader("Accept-Language")
		}

		localizer := i18n.GetLocalizer(lang, "en") // fallback to en
		c.Set("localizer", localizer)
		c.Next()
	}
}
