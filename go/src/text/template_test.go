package text

import (
	"html/template"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	p := make(map[string]interface{})
	p["Name"] = "zhangsan"
	p["Age"] = "10"
	tmpl, err := template.New("test").Parse("Name: {{.Name}}, Age: {{.Age}}")
	if err != nil {
		t.Fatal(err)
	}

	_ = tmpl.Execute(os.Stdout, p)
}
