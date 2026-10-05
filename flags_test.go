package main

import (
	"path/filepath"
	"strings"
	"testing"
	"tg-cli/connection"
)

func strPtr(s string) *string { return &s }

func TestCheckFlags(t *testing.T) {
	conn := connection.NewConnection()
	missingFile := filepath.Join(t.TempDir(), "missing.txt")

	tests := []struct {
		name    string
		flags   flags
		wantOk  bool
		wantErr string
	}{
		{
			name:    "no flags starts tui",
			flags:   flags{chatIdFlag: strPtr(""), fileFlag: strPtr(""), photoFlag: strPtr(""), captionFlag: strPtr("")},
			wantOk:  false,
			wantErr: "",
		},
		{
			name:    "file without chat id",
			flags:   flags{chatIdFlag: strPtr(""), fileFlag: strPtr("/tmp/a"), photoFlag: strPtr(""), captionFlag: strPtr("")},
			wantOk:  true,
			wantErr: "error: -chat is required when using -f or -p",
		},
		{
			name:    "caption without chat id",
			flags:   flags{chatIdFlag: strPtr(""), fileFlag: strPtr(""), photoFlag: strPtr(""), captionFlag: strPtr("cap")},
			wantOk:  true,
			wantErr: "error: -chat is required when using -f or -p",
		},
		{
			name:    "both file and photo",
			flags:   flags{chatIdFlag: strPtr("123"), fileFlag: strPtr("/tmp/a"), photoFlag: strPtr("/tmp/b"), captionFlag: strPtr("")},
			wantOk:  true,
			wantErr: "error: exactly one of -f or -p must be provided, not both",
		},
		{
			name:    "chat id without file and photo",
			flags:   flags{chatIdFlag: strPtr("123"), fileFlag: strPtr(""), photoFlag: strPtr(""), captionFlag: strPtr("cap")},
			wantOk:  true,
			wantErr: "error: exactly one of -f or -p must be provided, not both",
		},
		{
			name:    "invalid chat id",
			flags:   flags{chatIdFlag: strPtr("abc"), fileFlag: strPtr(missingFile), photoFlag: strPtr(""), captionFlag: strPtr("")},
			wantOk:  true,
			wantErr: "invalid syntax",
		},
		{
			name:    "file does not exist",
			flags:   flags{chatIdFlag: strPtr("123"), fileFlag: strPtr(missingFile), photoFlag: strPtr(""), captionFlag: strPtr("")},
			wantOk:  true,
			wantErr: "no such file or directory",
		},
		{
			name:    "photo does not exist",
			flags:   flags{chatIdFlag: strPtr("123"), fileFlag: strPtr(""), photoFlag: strPtr(missingFile), captionFlag: strPtr("")},
			wantOk:  true,
			wantErr: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok, err := checkFlags(conn, tt.flags)

			if ok != tt.wantOk {
				t.Errorf("ok = %v, want %v", ok, tt.wantOk)
			}
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("error = %q, want containing %q", err.Error(), tt.wantErr)
			}
		})
	}
}
