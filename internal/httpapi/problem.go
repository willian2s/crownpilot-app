package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/willian2s/crownpilot-app/internal/observability"
)

const problemType = "https://www.rfc-editor.org/rfc/rfc9457"

// writeProblem escreve uma resposta ProblemDetails com um código explícito.
//
// ProblemDetails e ProblemCode vêm de api.gen.go, gerados da spec. O schema não
// tem detail nem instance de propósito: é por eles que stack traces e URLs
// internas costumam vazar. Um código nunca carrega token, UID, Player Tag ou
// detalhe de infraestrutura.
func writeProblem(w http.ResponseWriter, r *http.Request, status int, code ProblemCode) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	// O corpo é uma struct fixa; erro de encode significa que o cliente foi
	// embora e não há mais nada útil a fazer.
	_ = json.NewEncoder(w).Encode(ProblemDetails{
		Type:    problemType,
		Title:   titleForStatus(status),
		Status:  status,
		Code:    code,
		TraceId: observability.RequestID(r.Context()),
	})
}

// defaultCode cobre só os status que o próprio transporte produz (404/405 do
// ServeMux, authentication, panics). Status de negócio como 409 ou 503 precisam
// ser escritos com código explícito, então um 503 nunca vira silenciosamente
// provider_unavailable, reservado ao lookup do Clash Royale.
func defaultCode(status int) ProblemCode {
	switch status {
	case http.StatusBadRequest:
		return ProblemCodeInvalidRequest
	case http.StatusUnauthorized:
		return ProblemCodeAuthenticationRequired
	case http.StatusForbidden:
		return ProblemCodeAuthorizationForbidden
	case http.StatusNotFound:
		return ProblemCodeResourceNotFound
	case http.StatusMethodNotAllowed:
		return ProblemCodeMethodNotAllowed
	default:
		return ProblemCodeInternalError
	}
}

func titleForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "The request is invalid."
	case http.StatusUnauthorized:
		return "Authentication is required."
	case http.StatusForbidden:
		return "The authenticated subject is not authorized."
	case http.StatusNotFound:
		return "The requested resource was not found."
	case http.StatusMethodNotAllowed:
		return "The request method is not allowed."
	case http.StatusConflict:
		return "The request conflicts with current state."
	case http.StatusUnprocessableEntity:
		return "The request failed semantic validation."
	case http.StatusTooManyRequests:
		return "The request was rate limited."
	case http.StatusServiceUnavailable:
		return "A required dependency is unavailable."
	default:
		return "The request could not be completed."
	}
}
