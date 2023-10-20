package main

import (
	"bytes"
	"text/template"
)

var errorsTemplate = `
{{ range .Errors }}

{{ if .HasComment }}{{ .ColumnComment }}{{ end -}}
func Is{{.CamelValue}}(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == {{ .TableName }}_{{ .Value }}.String() && e.Code == {{ .HTTPCode }} 
}
{{- end }}

{{ range .Errors }}

{{ if .HasComment }}{{ .ColumnComment }}{{ end -}}
func Err{{ .CamelValue }}(format string, args ...interface{}) *errors.Error {
	 return errors.New({{ .HTTPCode }}, {{ .TableName }}_{{ .Value }}.String(), fmt.Sprintf(format, args...))
}
{{- end }}

{{ range .Errors }}
{{ if .HasComment }}{{ .ColumnComment }}{{ end -}}
func Err{{ .CamelValue }}WithMeta(err error, msg string) *errors.Error {
	return Err{{ .CamelValue }}(msg).WithMetadata(locationErrMeta(err))
}
{{- end }}

{{ range .Errors }}
{{ if .HasComment }}{{ .ColumnComment }}{{ end -}}
func TryErr{{ .CamelValue }}Wrap(err error, msg string) *errors.Error {
	if e, ok := err.(*errors.Error); !ok {
		return Err{{ .CamelValue }}(msg).WithMetadata(locationErrMeta(err))
	} else {
		e.Reason = msg + " " + e.Reason
		return e.WithMetadata(locationErrMeta(err))
	}
}

func TryErr{{ .CamelValue }}Wrapf(err error, format string, args ...interface{}) *errors.Error {
	if e, ok := err.(*errors.Error); !ok {
		return Err{{ .CamelValue }}(format, args).WithMetadata(locationErrMeta(err))
	} else {
		e.Reason = fmt.Sprintf(format, args) + " " + e.Reason
		return e.WithMetadata(locationErrMeta(err))
	}
}
{{- end }}

func locationErrMeta(err error) map[string]string {
	if err == nil {
		return map[string]string{locKey: location()}
	}

	if ee := new(errors.Error); errors.As(err, &ee) {
		m := ee.GetMetadata()
		var hasLocation bool
		for k, v := range m {
			if k == locKey && len(v) > 0 {
				hasLocation = true
				break
			}
		}

		if !hasLocation {
			m[locKey] = location()
		}
		return m
	}

	return map[string]string{locKey: location(), "error": err.Error()}
}

func location() string {
	pc := make([]uintptr, 6)
	n := runtime.Callers(1, pc)
	if n <= 1 {
		return ""
	}
	currentFuncPc := runtime.FuncForPC(pc[0])
	currentFile, _ := currentFuncPc.FileLine(pc[0])
	for i := 1; i < n; i++ {
		file, line := runtime.FuncForPC(pc[i]).FileLine(pc[i])

		if filepath.Dir(currentFile) == filepath.Dir(file) {
			continue
		}

		var split [2]int
		for i := 0; i < len(file) || i < len(currentFile); i++ {
			if file[i] == '/' {
				split[0], split[1] = i, split[0]
			}
			if file[i] != currentFile[i] {
				break
			}
		}

		if split[1]+1 > len(file) {
			return ""
		}
		return fmt.Sprintf(" %s:%d ", file[split[1]+1:], line)
	}

	return ""
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
