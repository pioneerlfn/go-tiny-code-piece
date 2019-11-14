package leap

import (
	"reflect"
	"testing"
)

func TestLeap(t *testing.T) {
	tests := map[string]struct{
		year int
		isLeap bool
	}{"2000":{2000, true},
		"2019":{2019, false},
		"2020":{2020, true},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := IsLeap(tt.year)
			if !reflect.DeepEqual(tt.isLeap, got) {
				t.Errorf("want=%v, got=%v\n", tt.isLeap, got)
			}
		})
	}
}

