package Utils

import (
	"testing"
)

func TestAddress_String(t *testing.T) {
	type fields struct {
		Addr string
		Port uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ipv4",
			fields: fields{
				Addr: "127.0.0.1",
				Port: 8080,
			},
			want: "127.0.0.1:8080",
		},
		{
			name: "ipv6",
			fields: fields{
				Addr: "2001:db8:a0b:12f0::1",
				Port: 9000,
			},
			want: "2001:db8:a0b:12f0::1:9000",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, _ := NewAddress(tt.fields.Addr, tt.fields.Port)
			if got := a.String(); got != tt.want {
				t.Errorf("Address.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
