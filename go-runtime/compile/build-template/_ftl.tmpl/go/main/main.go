// Code generated by FTL. DO NOT EDIT.
package main

import (
	"context"

	"github.com/TBD54566975/ftl/backend/common/plugin"
	"github.com/TBD54566975/ftl/go-runtime/server"
	"github.com/TBD54566975/ftl/protos/xyz/block/ftl/v1/ftlv1connect"

	"ftl/{{.Name}}"
)

func main() {
  verbConstructor := server.NewUserVerbServer("{{.Name}}",
{{- range .Verbs}}
    server.Handle({{$.Name}}.{{.Name|camel}}),
{{- end}}
  )
  plugin.Start(context.Background(), "{{.Name}}", verbConstructor, ftlv1connect.VerbServiceName, ftlv1connect.NewVerbServiceHandler)
}