package httpapi

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"

	"github.com/willian2s/crownpilot-app/api/openapi"
)

// handleOpenAPI serve a fonte canônica embutida, byte a byte. Não há
// conversão nem documento gerado em runtime que possa divergir de v1.json.
func handleOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = io.WriteString(w, openapi.V1())
}

// Swagger UI é só visualização de /openapi/v1.json (ADR 005), nunca fonte do
// contrato. Os assets vêm do jsDelivr com versão exata e hash SRI: se o CDN
// entregar qualquer byte diferente, o browser se recusa a executar.
const (
	swaggerUIVersion = "5.33.0"
	swaggerUICSSSRI  = "sha384-Ov4/wv3j2bmct8cDc5X4ngJZohVPzEmc6uDPH8WeljUxO5vtoykvMEfbu9Vh6RaW"
	swaggerUIJSSRI   = "sha384-YDALVcy8kj8yltLBVi1vBiBAUqdxvus673gM8XKwiy6aDUJFXivF/KCufekjYbVf"
	swaggerUIBase    = "https://cdn.jsdelivr.net/npm/swagger-ui-dist@" + swaggerUIVersion
)

// validatorUrl "none" desliga o validador externo padrão do Swagger UI, que
// enviaria a URL da spec (por exemplo, a de Staging) para validator.swagger.io.
const docsInitScript = `SwaggerUIBundle({ url: "/openapi/v1.json", dom_id: "#swagger-ui", validatorUrl: "none" });`

const docsPage = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>CrownPilot API documentation</title>
<link rel="stylesheet" href="` + swaggerUIBase + `/swagger-ui.css" integrity="` + swaggerUICSSSRI + `" crossorigin="anonymous">
</head>
<body>
<div id="swagger-ui"></div>
<script src="` + swaggerUIBase + `/swagger-ui-bundle.js" integrity="` + swaggerUIJSSRI + `" crossorigin="anonymous"></script>
<script>` + docsInitScript + `</script>
</body>
</html>
`

// docsCSP libera só o CDN fixado e o script inline exato acima, identificado
// pelo seu hash SHA-256. O hash é calculado uma vez, na inicialização do
// pacote, a partir da própria constante: editar o script atualiza o hash.
var docsCSP = func() string {
	sum := sha256.Sum256([]byte(docsInitScript))
	scriptHash := "'sha256-" + base64.StdEncoding.EncodeToString(sum[:]) + "'"
	return "default-src 'none'; " +
		"script-src https://cdn.jsdelivr.net " + scriptHash + "; " +
		// Swagger UI aplica estilos inline em runtime.
		"style-src https://cdn.jsdelivr.net 'unsafe-inline'; " +
		"img-src 'self' data:; " +
		"connect-src 'self'; " +
		"frame-ancestors 'none'"
}()

func handleDocs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Security-Policy", docsCSP)
	_, _ = io.WriteString(w, docsPage)
}
