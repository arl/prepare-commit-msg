package main

import "testing"

func Test_commonPrefix(t *testing.T) {
	tests := []struct {
		name  string
		paths []string
		want  string
	}{
		{
			name:  "empty",
			paths: nil,
			want:  "",
		},
		{
			name:  "single file",
			paths: []string{"file"},
			want:  "file",
		},
		{
			name:  "2 files at root",
			paths: []string{"file1", "file2"},
			want:  "",
		},
		{
			name:  "2 files in dir",
			paths: []string{"dir/file1", "dir/file2"},
			want:  "dir",
		},
		{
			name:  "3 files in dir/dir",
			paths: []string{"dir/dir/file1", "dir/dir/file2", "dir/dir/file3"},
			want:  "dir/dir",
		},
		{
			name:  "3 files in different dirs",
			paths: []string{"dir/dir/file1", "dir/dir2/file2", "file3"},
			want:  "",
		},
		{
			name:  "3 files, 2 in same dir",
			paths: []string{"dir/dir/file1", "dir/dir2/file2", "dir/dir2"},
			want:  "dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := commonPrefix(tt.paths); got != tt.want {
				t.Errorf("commonPrefix() = %q, want %q", got, tt.want)
			}
		})
	}
}
