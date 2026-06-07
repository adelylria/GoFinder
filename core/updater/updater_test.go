package updater

import "testing"

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name    string
		latest  string
		current string
		want    bool
	}{
		{name: "newer major", latest: "v2.0.0", current: "v1.9.9", want: true},
		{name: "newer minor", latest: "v1.2.0", current: "v1.1.9", want: true},
		{name: "same", latest: "v1.2.3", current: "1.2.3", want: false},
		{name: "same short version", latest: "v0.3", current: "v0.3", want: false},
		{name: "newer patch from short version", latest: "v0.3.1", current: "v0.3", want: true},
		{name: "older", latest: "v1.2.3", current: "v1.3.0", want: false},
		{name: "dev current", latest: "v1.2.3", current: "dev", want: false},
		{name: "hash current", latest: "v1.2.3", current: "abc123", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNewerVersion(tt.latest, tt.current); got != tt.want {
				t.Fatalf("isNewerVersion(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}
