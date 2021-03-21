package main

import (
	"testing"
)

func TestFilenameWithoutExtension(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple_txt",
			input: "filename.txt",
			want:  "filename",
		},
		{
			name:  "simple_yaml",
			input: "config.yaml",
			want:  "config",
		},
		{
			name:  "full_linux_path",
			input: "/home/user/data/invoice_data.db",
			want:  "/home/user/data/invoice_data",
		},
		{
			name:  "file_without_ext",
			input: "invoice_data",
			want:  "invoice_data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := filenameWithoutExtension(tt.input)
			if output != tt.want {
				t.Errorf("oczekiwano: %q; otrzymano: %q", tt.want, output)
			}
		})
	}
}
