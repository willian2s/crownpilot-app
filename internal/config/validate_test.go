package config

import (
	"strings"
	"testing"
)

func validConfig(env Environment) Config {
	docsEnabled := env == Local || env == Staging

	// Origens http de loopback só são válidas em Local; os outros ambientes
	// começam com allowlist vazia, como os defaults reais.
	var origins []string
	if env == Local {
		origins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	}

	// Fora de Local as três variáveis do Firebase são obrigatórias.
	var firebase Firebase
	if env != Local {
		firebase = testFirebase()
	}

	return Config{
		Environment: env,
		Port:        DefaultPort,
		Firebase:    firebase,
		CORS: CORS{
			AllowedOrigins: origins,
		},
		OpenAPI: OpenAPI{
			ExposeJSON: docsEnabled,
			ExposeUI:   docsEnabled,
		},
		Auth: Auth{
			ContractFixtures: env == Local,
		},
	}
}

func TestConfigValidate(t *testing.T) {
	t.Parallel()

	withOrigins := func(env Environment, origins ...string) Config {
		cfg := validConfig(env)
		cfg.CORS.AllowedOrigins = origins
		return cfg
	}

	withPort := func(port int) Config {
		cfg := validConfig(Local)
		cfg.Port = port
		return cfg
	}

	withOpenAPI := func(env Environment, json, ui bool) Config {
		cfg := validConfig(env)
		cfg.OpenAPI.ExposeJSON = json
		cfg.OpenAPI.ExposeUI = ui
		return cfg
	}

	withFixtures := func(env Environment) Config {
		cfg := validConfig(env)
		cfg.Auth.ContractFixtures = true
		return cfg
	}

	tests := []struct {
		name       string
		cfg        Config
		wantErr    []string
		wantNoErrs []string
	}{
		{
			name: "Local defaults are valid",
			cfg:  validConfig(Local),
		},
		{
			name: "Preview defaults are valid",
			cfg:  validConfig(Preview),
		},
		{
			name: "Staging defaults are valid",
			cfg:  validConfig(Staging),
		},
		{
			name: "Production defaults are valid",
			cfg:  validConfig(Production),
		},
		{
			name:    "zero value environment is invalid",
			cfg:     validConfig(Environment("")),
			wantErr: []string{"environment must be"},
		},
		{
			name:    "lowercase environment is invalid",
			cfg:     validConfig(Environment("local")),
			wantErr: []string{"environment must be"},
		},

		// Limites da porta
		{
			name:    "port zero is invalid",
			cfg:     withPort(0),
			wantErr: []string{"port must be between 1 and 65535"},
		},
		{
			name:    "negative port is invalid",
			cfg:     withPort(-1),
			wantErr: []string{"port must be between 1 and 65535"},
		},
		{
			name:    "port above maximum is invalid",
			cfg:     withPort(65536),
			wantErr: []string{"port must be between 1 and 65535"},
		},
		{
			name: "minimum port is valid",
			cfg:  withPort(1),
		},
		{
			name: "default port is valid",
			cfg:  withPort(5080),
		},
		{
			name: "maximum port is valid",
			cfg:  withPort(65535),
		},

		// CORS: listas vazias
		{
			name: "empty CORS list is valid in Local",
			cfg:  withOrigins(Local),
		},
		{
			name: "empty CORS list is valid in Preview",
			cfg:  withOrigins(Preview),
		},
		{
			name: "empty CORS list is valid in Staging",
			cfg:  withOrigins(Staging),
		},
		{
			name: "empty CORS list is valid in Production",
			cfg:  withOrigins(Production),
		},

		// CORS: origens permitidas
		{
			name: "localhost HTTP is valid in Local",
			cfg:  withOrigins(Local, "http://localhost:5173"),
		},
		{
			name: "IPv4 loopback HTTP is valid in Local",
			cfg:  withOrigins(Local, "http://127.0.0.1:5173"),
		},
		{
			name: "IPv6 loopback HTTP is valid in Local",
			cfg:  withOrigins(Local, "http://[::1]:5173"),
		},
		{
			name:    "localhost HTTP is invalid in Staging",
			cfg:     withOrigins(Staging, "http://localhost:5173"),
			wantErr: []string{"must use https"},
		},
		{
			name:    "localhost HTTP is invalid in Production",
			cfg:     withOrigins(Production, "http://localhost:5173"),
			wantErr: []string{"must use https"},
		},
		{
			name:    "non-loopback HTTP is invalid even in Local",
			cfg:     withOrigins(Local, "http://app.example.com"),
			wantErr: []string{"must use https"},
		},
		{
			name: "HTTPS origin is valid in Staging",
			cfg:  withOrigins(Staging, "https://app.example.com"),
		},
		{
			name: "HTTPS origin is valid in Production",
			cfg:  withOrigins(Production, "https://app.example.com"),
		},

		// CORS: origens malformadas
		{
			name:    "wildcard origin is invalid",
			cfg:     withOrigins(Local, "*"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "empty origin entry is invalid",
			cfg:     withOrigins(Local, ""),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "trailing slash is invalid",
			cfg:     withOrigins(Local, "https://app.example.com/"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "origin with path is invalid",
			cfg:     withOrigins(Local, "https://app.example.com/path"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "origin with query is invalid",
			cfg:     withOrigins(Local, "https://app.example.com?x=1"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "origin with fragment is invalid",
			cfg:     withOrigins(Local, "https://app.example.com#f"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "origin with credentials is invalid",
			cfg:     withOrigins(Local, "https://user@app.example.com"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "uppercase origin is invalid",
			cfg:     withOrigins(Local, "HTTPS://App.Example.com"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "FTP origin is invalid",
			cfg:     withOrigins(Local, "ftp://app.example.com"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name:    "origin without scheme is invalid",
			cfg:     withOrigins(Local, "app.example.com"),
			wantErr: []string{"is not a valid origin"},
		},
		{
			name: "duplicate origin is invalid",
			cfg: withOrigins(
				Local,
				"https://app.example.com",
				"https://app.example.com",
			),
			wantErr: []string{"duplicated"},
		},

		// OpenAPI
		{
			name:    "OpenAPI UI requires JSON in Local",
			cfg:     withOpenAPI(Local, false, true),
			wantErr: []string{"OpenAPI UI requires the OpenAPI JSON to be exposed"},
		},
		{
			name:    "OpenAPI JSON is forbidden in Preview",
			cfg:     withOpenAPI(Preview, true, false),
			wantErr: []string{"OpenAPI exposure is disabled by policy in Preview"},
		},
		{
			name:    "OpenAPI UI is forbidden in Preview",
			cfg:     withOpenAPI(Preview, false, true),
			wantErr: []string{"OpenAPI exposure is disabled by policy in Preview"},
		},
		{
			name:    "OpenAPI JSON is forbidden in Production",
			cfg:     withOpenAPI(Production, true, false),
			wantErr: []string{"OpenAPI exposure is disabled by policy in Production"},
		},
		{
			name:    "OpenAPI UI is forbidden in Production",
			cfg:     withOpenAPI(Production, false, true),
			wantErr: []string{"OpenAPI exposure is disabled by policy in Production"},
		},
		{
			name: "OpenAPI JSON and UI are valid in Staging",
			cfg:  withOpenAPI(Staging, true, true),
		},
		{
			name: "OpenAPI disabled is valid in Staging",
			cfg:  withOpenAPI(Staging, false, false),
		},

		// Fixtures de authentication
		{
			name:    "contract fixtures are forbidden in Preview",
			cfg:     withFixtures(Preview),
			wantErr: []string{"contract authentication fixtures are allowed only in Local"},
		},
		{
			name:    "contract fixtures are forbidden in Staging",
			cfg:     withFixtures(Staging),
			wantErr: []string{"contract authentication fixtures are allowed only in Local"},
		},
		{
			name:    "contract fixtures are forbidden in Production",
			cfg:     withFixtures(Production),
			wantErr: []string{"contract authentication fixtures are allowed only in Local"},
		},

		// Fail-closed: várias violações independentes
		{
			name: "multiple violations are joined",
			cfg: Config{
				Environment: Production,
				Port:        0,
				CORS: CORS{
					AllowedOrigins: []string{"*"},
				},
				Auth: Auth{
					ContractFixtures: true,
				},
			},
			wantErr: []string{"port must be between", "is not a valid origin", "only in Local"},
		},

		// Fail-closed: ambiente inválido precisa retornar cedo
		{
			name: "invalid environment returns before other validations",
			cfg: Config{
				Environment: Environment("Invalid"),
				Port:        0,
				CORS: CORS{
					AllowedOrigins: []string{"*"},
				},
				Auth: Auth{
					ContractFixtures: true,
				},
			},
			wantErr: []string{"environment must be Local, Preview, Staging or Production"},
			wantNoErrs: []string{
				"port must be between 1 and 65535",
				"is not a valid origin",
				"contract authentication fixtures",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.cfg.Validate()

			if len(tt.wantErr) == 0 {
				if err != nil {
					t.Fatalf("Validate() returned unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Fatalf(
					"Validate() returned nil error, want error containing %q",
					tt.wantErr,
				)
			}

			for _, wantErr := range tt.wantErr {
				if !strings.Contains(err.Error(), wantErr) {
					t.Errorf(
						"Validate() error = %q, want substring %q",
						err.Error(),
						wantErr,
					)
				}
			}

			for _, unwanted := range tt.wantNoErrs {
				if strings.Contains(err.Error(), unwanted) {
					t.Errorf(
						"Validate() error = %q, must not contain %q",
						err.Error(),
						unwanted,
					)
				}
			}
		})
	}
}
