package httpapi

// Os tipos de api.gen.go são gerados a partir de api/openapi/v1.json. Depois de
// editar a spec, rode `go generate ./...` e versione o resultado; o CI falha
// se o arquivo gerado divergir da fonte.
//go:generate go tool -modfile=../../tools/go.mod oapi-codegen -config ../../api/openapi/oapi-codegen.yaml ../../api/openapi/v1.json
