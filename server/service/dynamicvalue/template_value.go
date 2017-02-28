package dynamicvalue

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	"github.com/dpb587/ssoca/auth"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type templateValue struct {
	template *template.Template
}

func NewTemplateValue(tpl *template.Template) templateValue {
	return templateValue{
		template: tpl,
	}
}

func CreateTemplateValue(value string) (templateValue, error) {
	tpl, err := template.New("configvalue").Funcs(template.FuncMap{
		"join":  strings.Join,
		"split": strings.Split,
	}).Parse(value)
	if err != nil {
		bosherr.WrapErrorf(err, "Parsing template: %s", value)
	}

	return NewTemplateValue(tpl), nil
}

func MustCreateTemplateValue(value string) templateValue {
	must, err := CreateTemplateValue(value)
	if err != nil {
		panic(err)
	}

	return must
}

func (cv templateValue) Evaluate(req *http.Request, token *auth.Token) (string, error) {
	data := struct {
		Request http.Request
		Token   auth.Token
	}{
		Request: *req,
	}

	if token != nil {
		data.Token = *token
	}

	out := new(bytes.Buffer)
	err := cv.template.Execute(out, data)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
