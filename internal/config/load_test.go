package config

import (
	"reflect"
	"strings"
	"testing"
)

// env monta uma LookupFunc a partir de um map, substituindo os.LookupEnv.
func env(vars map[string]string) LookupFunc {
	return func(key string) (string, bool) {
		v, ok := vars[key]
		return v, ok
	}
}

func TestLoadDefaults(t *testing.T) {
	tests := []struct {
		name string
		env  Environment
		want Config
	}{
		{
			name: "Local enables docs, fixtures and localhost origins",
			env:  Local,
			want: Config{
				Environment: Local,
				Port:        DefaultPort,
				CORS:        CORS{AllowedOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"}},
				OpenAPI:     OpenAPI{ExposeJSON: true, ExposeUI: true},
				Auth:        Auth{ContractFixtures: true},
			},
		},
		{
			name: "Preview disables everything",
			env:  Preview,
			want: Config{Environment: Preview, Port: DefaultPort},
		},
		{
			name: "Staging exposes docs only",
			env:  Staging,
			want: Config{Environment: Staging, Port: DefaultPort, OpenAPI: OpenAPI{ExposeJSON: true, ExposeUI: true}},
		},
		{
			name: "Production disables everything",
			env:  Production,
			want: Config{Environment: Production, Port: DefaultPort},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(env(map[string]string{EnvEnvironment: string(tt.env)}))
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestLoadOverrides(t *testing.T) {
	got, err := Load(env(map[string]string{
		EnvEnvironment:      "Staging",
		EnvPort:             "8080",
		EnvCORSOrigins:      "https://app.example.com, https://admin.example.com",
		EnvOpenAPIJSON:      "true",
		EnvOpenAPIUI:        "false",
		EnvContractFixtures: "false",
	}))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	want := Config{
		Environment: Staging,
		Port:        8080,
		CORS:        CORS{AllowedOrigins: []string{"https://app.example.com", "https://admin.example.com"}},
		OpenAPI:     OpenAPI{ExposeJSON: true},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Load() = %+v, want %+v", got, want)
	}
}

func TestLoadEmptyOriginsClearsLocalDefault(t *testing.T) {
	got, err := Load(env(map[string]string{EnvEnvironment: "Local", EnvCORSOrigins: ""}))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(got.CORS.AllowedOrigins) != 0 {
		t.Errorf("AllowedOrigins = %v, want empty", got.CORS.AllowedOrigins)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name    string
		vars    map[string]string
		wantErr []string // trechos que precisam aparecer todos no erro
	}{
		{
			name:    "missing environment",
			vars:    map[string]string{},
			wantErr: []string{EnvEnvironment},
		},
		{
			name: "environment names are exact",
			// "Development" era alias do ASP.NET Core; o processo Go não tem alias.
			vars:    map[string]string{EnvEnvironment: "Development"},
			wantErr: []string{EnvEnvironment},
		},
		{
			name:    "port is not a number",
			vars:    map[string]string{EnvEnvironment: "Local", EnvPort: "http"},
			wantErr: []string{EnvPort},
		},
		{
			name: "all parse errors are reported together",
			vars: map[string]string{
				EnvEnvironment: "Local",
				EnvPort:        "x",
				EnvOpenAPIUI:   "yes please",
			},
			wantErr: []string{EnvPort, EnvOpenAPIUI},
		},
		{
			name:    "parsed values still go through Validate",
			vars:    map[string]string{EnvEnvironment: "Production", EnvContractFixtures: "true"},
			wantErr: []string{"only in Local"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(env(tt.vars))
			if err == nil {
				t.Fatalf("Load() error = nil, want error")
			}
			for _, s := range tt.wantErr {
				if !strings.Contains(err.Error(), s) {
					t.Errorf("Load() error = %q, want it to contain %q", err, s)
				}
			}
			if !reflect.DeepEqual(got, Config{}) {
				t.Errorf("Load() returned %+v on error, want zero Config", got)
			}
		})
	}
}
