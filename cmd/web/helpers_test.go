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

func TestPrepareTextStyle(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "kapitaliki",
			input: "{{{1 stycznia 1500}}} roku",
			want:  `<span class="newthought">1 stycznia 1500</span> roku`,
		},
		{
			name:  "pogrubienia",
			input: "to jest {{istotna}} informacja",
			want:  `to jest <strong>istotna</strong> informacja`,
		},
		{
			name:  "italiki",
			input: "to jest {wyróżniony} tekst",
			want:  `to jest <em>wyróżniony</em> tekst`,
		},
		{
			name:  "złamanie_linii",
			input: `po tym słowie \\ następuje nowa linia tekstu`,
			want:  `po tym słowie </p><p> następuje nowa linia tekstu`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := prepareTextStyle(tt.input, false)
			if output != tt.want {
				t.Errorf("oczekiwano: %q; otrzymano: %q", tt.want, output)
			}
		})
	}
}
