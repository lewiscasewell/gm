package util

import (
	"fmt"
	"os"
	"testing"
)

func TestGetFileUrl(t *testing.T) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		t.Error(err)
	}

	type testCase struct {
		name string
		file string
		want string
	}
	tests := []testCase{
		{
			name: "test for profile file url",
			file: "profile",
			want: fmt.Sprintf("%s/gm/profile.json", dirname),
		},
		{
			name: "test for wallet file url",
			file: "wallet",
			want: fmt.Sprintf("%s/gm/wallet.json", dirname),
		},
		{
			name: "test an invalid file url",
			file: "invalid",
			want: "",
		},
		{
			name: "test an empty file url",
			file: "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetFileUrl(tt.file)
			if got != tt.want {
				t.Errorf("GetFileUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
