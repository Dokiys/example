package cel_go

import (
	"fmt"
	"log"
	"testing"

	"github.com/Dokiys/go_test/go/pkg/github.com/google/cel-go/example"
	"github.com/google/cel-go/cel"
)

func TestCel(t *testing.T) {
	env, err := cel.NewEnv(
		cel.Types(&example.Person{}),
		cel.Variable("people", cel.ObjectType("cel_go.example.Person")),
	)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}
	ast, iss := env.Compile(`"Hello world! I'm " + people.name + "."`)
	// Check iss for compilation errors.
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalln(err)
	}
	out, _, err := prg.Eval(map[string]any{
		"people": &example.Person{Name: "zhangsan"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(out)
}
