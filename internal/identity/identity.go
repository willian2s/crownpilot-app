// Package identity define o contrato de identidade do CrownPilot.
//
// É um módulo funcional: não importa HTTP, configuração, Firebase ou banco
// (regra do depguard em .golangci.yml). Adapters como internal/platform/firebase
// produzem os tipos daqui; internal/httpapi os consome.
package identity

import (
	"errors"
	"strings"
	"time"
)

// Erros do boundary de authentication. Adapters os embrulham com %w e quem
// consome decide o status com errors.Is, sem conhecer o provider.
var (
	// ErrInvalidCredential cobre qualquer token recusado: ausente, malformado,
	// expirado, assinatura, issuer, audience ou kid incorretos, e revogado. É
	// um único erro de propósito: o 401 precisa ser indistinguível.
	ErrInvalidCredential = errors.New("identity: invalid credential")

	// ErrAuthenticationUnavailable indica que o provider não pôde ser
	// consultado (timeout, rede, resposta inesperada). Só a checagem sensível
	// o transforma em 503 authentication_unavailable.
	ErrAuthenticationUnavailable = errors.New("identity: authentication provider unavailable")

	// ErrReauthenticationRequired indica auth_time ausente ou antigo demais
	// para uma operação sensível (403 reauthentication_required).
	ErrReauthenticationRequired = errors.New("identity: recent authentication required")
)

// Janela de autenticação recente da spec 002 para operações sensíveis.
const (
	RecentAuthenticationMaxAge = 5 * time.Minute
	// ClockSkewTolerance aceita auth_time um pouco no futuro, para absorver a
	// diferença entre o relógio do Firebase e o do servidor.
	ClockSkewTolerance = 60 * time.Second
)

// AuthenticatedSubject é o subject externo produzido pela authentication
// depois que o provider verificou o token.
//
// Os campos são não exportados: fora deste pacote só é possível obter um
// subject pelo construtor, então um UID vindo de body, query ou header não
// consegue virar subject sem passar pela verificação. O Firebase UID nunca é
// identidade de domínio; o mapeamento para CrownPilotUserID chega na 002-15.
type AuthenticatedSubject struct {
	firebaseUID string
	authTime    time.Time
}

// NewAuthenticatedSubject cria o subject a partir de claims já verificadas.
// authTime pode ser zero quando o token não tem auth_time; nesse caso as
// operações sensíveis exigem nova autenticação.
func NewAuthenticatedSubject(firebaseUID string, authTime time.Time) (AuthenticatedSubject, error) {
	if strings.TrimSpace(firebaseUID) == "" || firebaseUID != strings.TrimSpace(firebaseUID) {
		return AuthenticatedSubject{}, ErrInvalidCredential
	}
	return AuthenticatedSubject{firebaseUID: firebaseUID, authTime: authTime}, nil
}

// FirebaseUID devolve o subject externo verificado. Não deve ir para logs,
// erros ou respostas HTTP.
func (s AuthenticatedSubject) FirebaseUID() string {
	return s.firebaseUID
}

// IsZero informa se o subject é o valor zero, ou seja, se não houve
// authentication.
func (s AuthenticatedSubject) IsZero() bool {
	return s.firebaseUID == ""
}

// RequireRecentAuthentication aplica a regra de operações sensíveis: auth_time
// presente, no máximo RecentAuthenticationMaxAge atrás e no máximo
// ClockSkewTolerance no futuro, em relação a now (o relógio do servidor).
func (s AuthenticatedSubject) RequireRecentAuthentication(now time.Time) error {
	switch {
	case s.authTime.IsZero():
		return ErrReauthenticationRequired
	case s.authTime.After(now.Add(ClockSkewTolerance)):
		return ErrReauthenticationRequired
	case now.Sub(s.authTime) > RecentAuthenticationMaxAge:
		return ErrReauthenticationRequired
	default:
		return nil
	}
}
