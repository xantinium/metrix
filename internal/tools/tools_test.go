package tools_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/xantinium/metrix/internal/tools"
)

func TestFloatToStr(t *testing.T) {
	tests := []struct {
		value float64
		want  string
	}{
		{
			value: 0,
			want:  "0",
		},
		{
			value: 0.5,
			want:  "0.5",
		},
		{
			value: -1.62,
			want:  "-1.62",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("value=%f", tt.value), func(t *testing.T) {
			if got := tools.FloatToStr(tt.value); got != tt.want {
				t.Errorf("FloatToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToStr(t *testing.T) {
	tests := []struct {
		value int64
		want  string
	}{
		{
			value: 0,
			want:  "0",
		},
		{
			value: 5,
			want:  "5",
		},
		{
			value: -82,
			want:  "-82",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("value=%d", tt.value), func(t *testing.T) {
			if got := tools.IntToStr(tt.value); got != tt.want {
				t.Errorf("IntToStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToFloat(t *testing.T) {
	tests := []struct {
		value   string
		want    float64
		wantErr bool
	}{
		{
			value: "0",
			want:  0,
		},
		{
			value: "0.5",
			want:  0.5,
		},
		{
			value: "-1.62",
			want:  -1.62,
		},
		{
			value:   "5,2",
			wantErr: true,
		},
		{
			value:   "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("value=%s", tt.value), func(t *testing.T) {
			got, err := tools.StrToFloat(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToFloat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStrToInt(t *testing.T) {
	tests := []struct {
		value   string
		want    int
		wantErr bool
	}{
		{
			value: "0",
			want:  0,
		},
		{
			value: "5",
			want:  5,
		},
		{
			value: "-82",
			want:  -82,
		},
		{
			value:   "12.6",
			wantErr: true,
		},
		{
			value:   "invalid",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("value=%s", tt.value), func(t *testing.T) {
			got, err := tools.StrToInt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("StrToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("StrToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompression(t *testing.T) {
	tests := []struct {
		data []byte
	}{
		{
			data: []byte{},
		},
		{
			data: []byte("some value"),
		},
		{
			data: []byte(`{"id":"Alloc","value":12.5}`),
		},
	}

	for _, tt := range tests {
		compressedData, err := tools.Compress(tt.data)
		require.NoError(t, err)

		var got []byte
		got, err = tools.Decompress(compressedData)
		require.NoError(t, err)

		require.Equal(t, tt.data, got)
	}
}
