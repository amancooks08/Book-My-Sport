package service

import (
	"net/http"
	"reflect"
	"testing"
)


func TestRegisterCustomer(t *testing.T) {
	type args struct {
		deps dependencies
	}
	tests := []struct {
		name string
		args args
		want http.HandlerFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RegisterCustomer(tt.args.deps); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RegisterCustomer() = %v, want %v", got, tt.want)
			}
		})
	}
}
