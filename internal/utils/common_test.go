package utils

import (
	"testing"
)

func TestIsIPValid(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name  string
		args  args
		valid bool
		ip    string
	}{
		{
			name: "valid ip address shouod pass",
			args: args{
				ip: "200.106.141.15",
			},
			valid: true,
			ip:    "200.106.141.15",
		},
		{
			name: "invalid valid ip address length shouod not pass",
			args: args{
				ip: "200.106.141",
			},
			valid: false,
			ip:    "200.106.141",
		},
		{
			name: "invalid valid ip address shouod not pass",
			args: args{
				ip: "200.106.141.er",
			},
			valid: false,
			ip:    "200.106.141.er",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, ip := IsIPValid(tt.args.ip)
			if valid != tt.valid {
				t.Errorf("ValidIP() got = %v, valid %v", valid, tt.valid)
			}
			if ip != tt.ip {
				t.Errorf("ValidIP() got1 = %v, valid %v", ip, tt.ip)
			}
		})
	}
}

func TestIsStringValid(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		valid bool
		value string
	}{
		{
			name: "non empty string value should pass",
			args: args{
				s: "a valid string",
			},
			valid: true,
			value: "a valid string",
		},
		{
			name: "empty string value should fail",
			args: args{
				s: "",
			},
			valid: false,
			value: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsStringValid(tt.args.s)
			if got != tt.valid {
				t.Errorf("ValidString() got = %v, want %v", got, tt.valid)
			}
			if got1 != tt.value {
				t.Errorf("ValidString() got1 = %v, want %v", got1, tt.value)
			}
		})
	}
}

func TestIsFloat64Valid(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name  string
		args  args
		valid bool
		value float64
	}{
		{
			name: "valid decimal value should pass",
			args: args{
				s: "57.565675676",
			},
			valid: true,
			value: 57.565675676,
		},
		{
			name: "invalid valid decimal value should pass",
			args: args{
				s: "hello",
			},
			valid: false,
			value: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := IsFloat64Valid(tt.args.s)
			if got != tt.valid {
				t.Errorf("ValidFloat64() got = %v, want %v", got, tt.valid)
			}
			if got1 != tt.value {
				t.Errorf("ValidFloat64() got1 = %v, want %v", got1, tt.value)
			}
		})
	}
}
