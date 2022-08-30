package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == {{ .Name }}_{{ .Value }}.String() && e.Code == {{ .HTTPCode }} 
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}(format string, args ...interface{}) *errors.Error {
	 return errors.New({{ .HTTPCode }}, {{ .Name }}_{{ .Value }}.String(), fmt.Sprintf(format, args...))
}

{{ if .HasComment }}{{ .Comment }}{{ end -}}
func Error{{ .CamelValue }}WithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	if errors.Is(err, sql.ErrNoRows) {
		return Error{{ .CamelValue }}("数据不存在！").WithMetadata(m)
	}
	return ErrorSystemError("系统繁忙！").WithMetadata(m)
}

{{- end }}

func SystemErrorWithMeta(err error) *errors.Error {
	m := map[string]string{"location": location(), "error": err.Error()}
	return errors.New(500, ErrorReason_SystemError.String(), "系统繁忙！").WithMetadata(m)
}

func location() string {
	_, lcErr, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	_, lcCaller, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}

	if !strings.HasSuffix(lcErr, "errors.go") || strings.HasSuffix(lcCaller, "errors.go") {
		return ""
	}

	var split [2]int
	for i := 0; i < len(lcCaller) || i < len(lcErr); i++ {
		if lcCaller[i] == '/' {
			split[0], split[1] = i, split[0]
		}
		if lcCaller[i] != lcErr[i] {
			break
		}
	}

	if split[1]+1 > len(lcCaller) {
		return ""
	}
	return fmt.Sprintf("%s:%d\n", lcCaller[split[1]+1:], line)
}
`

type errorInfo struct {
	Name       string
	Value      string
	HTTPCode   int
	CamelValue string
	Comment    string
	HasComment bool
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("errors").Parse(errorsTemplate)
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, e); err != nil {
		panic(err)
	}
	return buf.String()
}
