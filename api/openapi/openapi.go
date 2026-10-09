// Package openapi embute a fonte canônica do contrato, v1.json.
//
// A diretiva go:embed copia o arquivo para dentro do binário em tempo de
// compilação. O processo serve exatamente esses bytes em /openapi/v1.json:
// não existe cópia, conversão nem binário auxiliar que possa divergir da fonte.
package openapi

import _ "embed" // necessário para a diretiva //go:embed

//go:embed v1.json
var v1 string

// V1 devolve o documento OpenAPI v1 como está no repositório. É uma função
// sobre uma string (imutável em Go) para que nenhum pacote altere o conteúdo.
func V1() string {
	return v1
}
