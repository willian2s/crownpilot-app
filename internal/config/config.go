// Package config carrega e valida a configuração de runtime não secreta do
// processo da API.
//
// O startup é fail-closed: cmd/crownpilot-api chama Load antes de abrir
// qualquer listener, e qualquer erro de parse ou de política encerra o
// processo. Load só lê por meio de uma função de lookup injetada, então os
// testes nunca tocam no ambiente real.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Environment é o ambiente de runtime. Só as quatro constantes abaixo são
// válidas; o valor zero "" é inválido de propósito, para que uma configuração
// ausente nunca seja confundida com um ambiente real.
type Environment string

const (
	Local      Environment = "Local"
	Preview    Environment = "Preview"
	Staging    Environment = "Staging"
	Production Environment = "Production"
)

func (e Environment) valid() bool {
	switch e {
	case Local, Preview, Staging, Production:
		return true
	}
	return false
}

// DefaultPort é a porta local canônica, a mesma do baseline .NET.
// Plataformas de hosting a sobrescrevem via PORT.
const DefaultPort = 5080

// Nomes das variáveis de ambiente lidas por Load.
const (
	EnvEnvironment      = "CROWNPILOT_ENVIRONMENT"
	EnvPort             = "PORT"
	EnvCORSOrigins      = "CROWNPILOT_CORS_ALLOWED_ORIGINS"
	EnvOpenAPIJSON      = "CROWNPILOT_OPENAPI_EXPOSE_JSON"
	EnvOpenAPIUI        = "CROWNPILOT_OPENAPI_EXPOSE_UI"
	EnvContractFixtures = "CROWNPILOT_AUTH_CONTRACT_FIXTURES"
)

// Config é a configuração de runtime validada. Não contém secrets.
type Config struct {
	Environment Environment
	Port        int
	CORS        CORS
	OpenAPI     OpenAPI
	Auth        Auth
}

// CORS lista as origens exatas que podem chamar a API a partir de um browser.
type CORS struct {
	AllowedOrigins []string
}

// OpenAPI controla a exposição de /openapi/v1.json (JSON) e /docs (UI).
type OpenAPI struct {
	ExposeJSON bool
	ExposeUI   bool
}

// Auth guarda as configurações do boundary de authentication desta task.
// ContractFixtures habilita tokens bearer determinísticos de Local que provam
// o pipeline 401/403/200; não são validação Firebase (task 002-14).
type Auth struct {
	ContractFixtures bool
}

// LookupFunc tem a assinatura de os.LookupEnv. O booleano diferencia uma
// variável não definida (usa o default do ambiente) de uma definida como "".
type LookupFunc func(key string) (string, bool)

// Load faz o parse da configuração via lookup, aplica os defaults de cada
// ambiente e devolve o resultado de Validate. Em qualquer erro, devolve um
// Config zero.
func Load(lookup LookupFunc) (Config, error) {
	rawEnv, ok := lookup(EnvEnvironment)
	env := Environment(rawEnv)
	if !ok || !env.valid() {
		// Os defaults dependem do ambiente, então nada mais pode ser lido.
		return Config{}, fmt.Errorf("%s must be Local, Preview, Staging or Production", EnvEnvironment)
	}

	// A documentação só fica ligada por padrão onde a spec permite.
	docsDefault := env == Local || env == Staging
	cfg := Config{
		Environment: env,
		Port:        DefaultPort,
		CORS:        CORS{AllowedOrigins: defaultOrigins(env)},
		OpenAPI:     OpenAPI{ExposeJSON: docsDefault, ExposeUI: docsDefault},
		Auth:        Auth{ContractFixtures: env == Local},
	}

	// Acumula todos os erros de parse em vez de parar no primeiro, para que um
	// deploy mal configurado mostre todos os problemas numa única falha.
	var errs []error
	if v, ok := lookup(EnvPort); ok {
		port, err := strconv.Atoi(v)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s must be an integer", EnvPort))
		}
		cfg.Port = port
	}
	if v, ok := lookup(EnvCORSOrigins); ok {
		cfg.CORS.AllowedOrigins = parseList(v)
	}
	errs = append(errs,
		parseBool(lookup, EnvOpenAPIJSON, &cfg.OpenAPI.ExposeJSON),
		parseBool(lookup, EnvOpenAPIUI, &cfg.OpenAPI.ExposeUI),
		parseBool(lookup, EnvContractFixtures, &cfg.Auth.ContractFixtures),
	)
	if err := errors.Join(errs...); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate aplica a política de ambientes da spec 002. Reporta todas as
// violações unidas num único erro, ou nil quando a configuração é aceitável.
func (c Config) Validate() error {
	if !c.Environment.valid() {
		return errors.New("environment must be Local, Preview, Staging or Production")
	}

	var errs []error
	if c.Port < 1 || c.Port > 65535 {
		errs = append(errs, errors.New("port must be between 1 and 65535"))
	}

	seen := make(map[string]bool, len(c.CORS.AllowedOrigins))
	for _, origin := range c.CORS.AllowedOrigins {
		if seen[origin] {
			errs = append(errs, fmt.Errorf("CORS origin %q is duplicated", origin))
		}
		seen[origin] = true
		if err := validateOrigin(origin, c.Environment); err != nil {
			errs = append(errs, err)
		}
	}

	if c.OpenAPI.ExposeUI && !c.OpenAPI.ExposeJSON {
		errs = append(errs, errors.New("OpenAPI UI requires the OpenAPI JSON to be exposed"))
	}
	if (c.Environment == Preview || c.Environment == Production) &&
		(c.OpenAPI.ExposeJSON || c.OpenAPI.ExposeUI) {
		errs = append(errs, fmt.Errorf("OpenAPI exposure is disabled by policy in %s", c.Environment))
	}

	if c.Auth.ContractFixtures && c.Environment != Local {
		errs = append(errs, errors.New("contract authentication fixtures are allowed only in Local"))
	}

	return errors.Join(errs...)
}

// Addr é o endereço de escuta para net/http. Host vazio escuta em todas as
// interfaces, o que containers precisam para receber tráfego encaminhado.
func (c Config) Addr() string {
	return ":" + strconv.Itoa(c.Port)
}

func defaultOrigins(env Environment) []string {
	if env == Local {
		return []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}
	return nil
}

// parseList separa um valor por vírgulas. Itens vazios são mantidos de
// propósito ("a,,b") para que Validate os rejeite, em vez de descartar entrada
// silenciosamente.
func parseList(v string) []string {
	if v == "" {
		return nil
	}
	items := strings.Split(v, ",")
	for i, item := range items {
		items[i] = strings.TrimSpace(item)
	}
	return items
}

// parseBool sobrescreve *dst só quando key está definida; senão mantém o default.
func parseBool(lookup LookupFunc, key string, dst *bool) error {
	v, ok := lookup(key)
	if !ok {
		return nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fmt.Errorf("%s must be true or false", key)
	}
	*dst = b
	return nil
}

// validateOrigin aceita só a origem serializada exata que um browser envia no
// header Origin: scheme://host[:port] em minúsculas, sem path, barra final,
// query, fragment ou credenciais. Qualquer outra forma nunca casaria com um
// request real, então é rejeitada em vez de ignorada silenciosamente.
func validateOrigin(origin string, env Environment) error {
	invalid := fmt.Errorf("CORS origin %q is not a valid origin", origin)
	if origin == "" || origin == "*" {
		return invalid
	}
	u, err := url.Parse(origin)
	if err != nil || u.Host == "" || u.User != nil ||
		origin != u.Scheme+"://"+u.Host || origin != strings.ToLower(origin) {
		return invalid
	}

	switch u.Scheme {
	case "https":
		return nil
	case "http":
		if env == Local && isLoopback(u.Hostname()) {
			return nil
		}
		return fmt.Errorf("CORS origin %q must use https (http only for localhost in Local)", origin)
	default:
		return invalid
	}
}

func isLoopback(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}
