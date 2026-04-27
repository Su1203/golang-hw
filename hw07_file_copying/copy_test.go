package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name   string
		offset int64
		limit  int64
	}{
		{
			name:   "copy full file",
			offset: 0,
			limit:  0,
		},
		{
			name:   "copy with limit",
			offset: 0,
			limit:  1000,
		},
		{
			name:   "copy with offset and limit",
			offset: 100,
			limit:  1000,
		},
		{
			name:   "copy near eof",
			offset: 6000,
			limit:  1000,
		},
	}

	fromPath := filepath.Join("testdata", "input.txt")
	sourceData, err := os.ReadFile(fromPath)
	if err != nil {
		t.Fatalf("failed to read source file: %v", err)
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			toPath := filepath.Join(t.TempDir(), "output.txt")

			if err := Copy(fromPath, toPath, tc.offset, tc.limit); err != nil {
				t.Fatalf("Copy() returned error: %v", err)
			}

			got, err := os.ReadFile(toPath)
			if err != nil {
				t.Fatalf("failed to read copied file: %v", err)
			}

			start := tc.offset
			end := int64(len(sourceData))
			if tc.limit > 0 && start+tc.limit < end {
				end = start + tc.limit
			}

			want := sourceData[start:end]
			wantQuoted := fmt.Sprintf("%q", string(want))
			gotQuoted := fmt.Sprintf("%q", string(got))

			if string(got) != string(want) {
				t.Fatalf(
					"copied content mismatch\nwant len: %d\ngot len: %d\nwant: %s\ngot: %s",
					len(want),
					len(got),
					wantQuoted,
					gotQuoted,
				)
			}
		})
	}

	t.Run("offset exceeds file size", func(t *testing.T) {
		toPath := filepath.Join(t.TempDir(), "output.txt")

		err := Copy(fromPath, toPath, 10000, 0)
		if err == nil {
			t.Fatal("Copy() error = nil, want non-nil")
		}

		if !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Fatalf("Copy() error = %v, want wrapped ErrOffsetExceedsFileSize", err)
		}
	})
}
