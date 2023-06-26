package dopTools

import (
	"strconv"
	"testing"
)

func TestValidateIin(t *testing.T) {
	tests := []struct {
		v    string
		want bool
	}{
		{
			v:    "870504300822",
			want: true,
		},
		{
			v:    "870504300821",
			want: false,
		},
		{
			v:    "870504310822",
			want: false,
		},
	}
	for i, tt := range tests {
		t.Run("Case-"+strconv.Itoa(i+1), func(t *testing.T) {
			if got := ValidateIin(tt.v); got != tt.want {
				t.Errorf("ValidateIin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		v    string
		want bool
	}{
		{
			v:    "asd@gmail..com",
			want: false,
		},
		{
			v:    "asd.@gmail.com",
			want: false,
		},
		{
			v:    "asd@gmail.com",
			want: true,
		},
		{
			v:    "asd@gm.ail.com",
			want: true,
		},
		{
			v:    "asd.dsa@gm.ail.com",
			want: true,
		},
		{
			v:    "asd-dsa@gm.ail.com",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.v, func(t *testing.T) {
			if got := ValidateEmail(tt.v); got != tt.want {
				t.Errorf("ValidateEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
