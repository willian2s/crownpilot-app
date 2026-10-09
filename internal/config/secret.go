package config

import "log/slog"

// redacted é o que qualquer impressão de um Secret mostra no lugar do valor.
const redacted = "[REDACTED]"

// Secret guarda um valor sensível, como a private key do Firebase.
//
// O campo é não exportado e o tipo implementa fmt.Stringer, fmt.GoStringer,
// slog.LogValuer e encoding.TextMarshaler. Assim fmt.Printf("%v", cfg),
// "%+v", "%#v", slog e encoding/json mostram "[REDACTED]". O valor real só sai
// por Reveal, que deixa explícito no código cada ponto onde o secret é usado.
type Secret struct {
	value string
}

// NewSecret embrulha um valor sensível.
func NewSecret(value string) Secret {
	return Secret{value: value}
}

// Reveal devolve o valor real. Use só onde ele precisa ir para o SDK, nunca
// para log, erro ou resposta HTTP.
func (s Secret) Reveal() string {
	return s.value
}

// IsZero informa se o secret está vazio, sem expor o conteúdo.
func (s Secret) IsZero() bool {
	return s.value == ""
}

func (s Secret) String() string   { return redacted }
func (s Secret) GoString() string { return redacted }

// LogValue faz o slog registrar "[REDACTED]" mesmo quando o Secret é passado
// diretamente como atributo de log.
func (s Secret) LogValue() slog.Value { return slog.StringValue(redacted) }

// MarshalText cobre encoding/json e outros encoders baseados em texto.
func (s Secret) MarshalText() ([]byte, error) { return []byte(redacted), nil }
