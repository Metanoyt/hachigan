package kubernetes

import "testing"

func TestApplicationName(t *testing.T) {
	tests := []struct {
		name     string
		labels   map[string]string
		fallback string
		want     string
	}{
		{"prefers recommended label", map[string]string{"app.kubernetes.io/name": "api", "app": "legacy"}, "deploy", "api"},
		{"uses app label", map[string]string{"app": "legacy"}, "deploy", "legacy"},
		{"falls back", nil, "deploy", "deploy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applicationName(tt.labels, tt.fallback); got != tt.want {
				t.Fatalf("applicationName() = %q, want %q", got, tt.want)
			}
		})
	}
}
