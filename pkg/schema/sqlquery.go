package schema

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type SqlQuery struct {
	Query  string `json:"query"`
	Params []any  `json:"params"`
}

func ExtractMustacheKey(s string) string {
	l := log.WithFields(log.Fields{
		"pkg": "schema",
		"fn":  "ExtractMustacheKey",
	})
	l.Debug("Extracting mustache key")
	var key string
	for _, k := range strings.Split(s, "{{") {
		if strings.Contains(k, "}}") {
			key = strings.Split(k, "}}")[0]
			break
		}
	}
	return key
}

func ExtractMustacheKeys(s string) []string {
	l := log.WithFields(log.Fields{
		"pkg": "schema",
		"fn":  "ExtractMustacheKeys",
	})
	l.Debug("Extracting mustache keys")
	keys := []string{}
	for _, k := range strings.Split(s, "{{") {
		if strings.Contains(k, "}}") {
			keys = append(keys, strings.Split(k, "}}")[0])
		}
	}
	return keys
}

func ReplaceParams(bd []byte, params []any) []any {
	for i, v := range params {
		sv := fmt.Sprintf("%s", v)
		if sv == "{{pushx_payload}}" {
			params[i] = bd
		} else if strings.Contains(sv, "{{") {
			key := ExtractMustacheKey(sv)
			params[i] = gjson.GetBytes(bd, key).String()
		}
	}
	return params
}

func ReplaceJSONKey(query string, k string, v string) string {
	l := log.WithFields(log.Fields{
		"pkg": "schema",
		"fn":  "ReplaceJSONKey",
		"k":   k,
		"v":   v,
	})
	l.Debug("Replacing JSON key")
	return strings.ReplaceAll(query, "{{"+k+"}}", v)
}

func ReplaceParamsString(bd []byte, params string) string {
	l := log.WithFields(log.Fields{
		"pkg": "schema",
		"fn":  "ReplaceParamsString",
	})
	l.Debug("Replacing params string")
	s := strings.ReplaceAll(params, "{{pushx_payload}}", string(bd))
	keys := ExtractMustacheKeys(s)
	for _, k := range keys {
		jv := gjson.GetBytes(bd, k)
		s = ReplaceJSONKey(s, k, jv.String())
	}
	return s
}
