package dynamicvalue

import (
	"bytes"
	"net/http"
	"strings"
	"text/template"

	"github.com/pkg/errors"

	"github.com/dpb587/ssoca/auth"
)

type templateValue struct {
	template *template.Template
}

func NewTemplateValue(tpl *template.Template) Value {
	return templateValue{
		template: tpl,
	}
}

func CreateTemplateValue(value string) (Value, error) {
	tpl, err := template.New("configvalue").Funcs(template.FuncMap{
		"join":  strings.Join,
		"split": strings.Split,
	}).Parse(value)
	if err != nil {
		return templateValue{}, errors.Wrapf(err, "parsing template: %s", value)
	}

	return NewTemplateValue(tpl), nil
}

func MustCreateTemplateValue(value string) Value {
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
