package rest

import (
    "html/template"
    "os"
)

//CreateTemplate CreateTemplate
func CreateTemplate() {
    name := "zhw"
    tmpl, err := template.New("test").Parse("hello,{{.}}")
    if err != nil {
        panic(err)
    }
    err = tmpl.Execute(os.Stdout, name)
    if err != nil {
        panic(err)
    }
}
