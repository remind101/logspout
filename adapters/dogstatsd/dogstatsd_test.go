package dogstatsd

import (
	"fmt"
	"reflect"
	"testing"
)

func TestLabelsFromString(t *testing.T) {
	tests := []struct {
		in  string
		out map[string]bool
	}{
		{"", nil},
		{"empire.app.name", map[string]bool{"empire.app.name": true}},
		{"empire.app.name,empire.app.process", map[string]bool{"empire.app.name": true, "empire.app.process": true}},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			out := labelsFromString(tt.in)
			if !reflect.DeepEqual(out, tt.out) {
				t.Fatalf("%v; want %v", out, tt.out)
			}
		})
	}
}
