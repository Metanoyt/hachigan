package domain

import "testing"

func TestDeriveHealth(t *testing.T) {
	tests := []struct {
		name    string
		desired int
		ready   int
		want    HealthStatus
	}{
		{"all ready", 3, 3, HealthOK},
		{"partially ready", 3, 1, HealthWarning},
		{"none ready", 3, 0, HealthCritical},
		{"no desired", 0, 0, HealthUnknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DeriveHealth(tt.desired, tt.ready); got != tt.want {
				t.Fatalf("DeriveHealth() = %q, want %q", got, tt.want)
			}
		})
	}
}
