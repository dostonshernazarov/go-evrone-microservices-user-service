package app

import (
	"context"
)

type ctxKeyLocalization int

const (
	EnvironmentProduction                    = "production"
	EnvironmentDevelop                       = "develop"
	CtxKeyLocalization    ctxKeyLocalization = 0
)

func GetLocalizationFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(CtxKeyLocalization).(string); ok {
		return lang
	}
	return ""
}
