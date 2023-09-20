package regexhandlers_test

import (
	regexphandlers "github.com/BorisPlus/regexphandlers"
)

var (
	none    = regexphandlers.Params{}
	ids     = regexphandlers.Params{"parent_id", "child_name"}
	message = regexphandlers.Params{"id", "title", "text"}
)

func Handlers() regexphandlers.RegexpHandlers {
	return regexphandlers.NewRegexpHandlers(
		DefaultHandler{},
		*regexphandlers.NewRegexpHandler(
			`/api/version`,
			none,
			VersionHandler{},
		),
		*regexphandlers.NewRegexpHandler(
			`/api/get/{numeric}/{string}`, // "parent_id", "child_name"
			ids,
			GetHandler{},
		),
		*regexphandlers.NewRegexpHandler(
			`/api/get/{numeric}/{string}/{any}`, // "parent_id", "child_name"
			message,
			GetMessage{},
		),
	)
}
