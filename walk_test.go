// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"reflect"
	"testing"
)

func TestSplitExt(t *testing.T) {
	tests := []struct {
		filename string
		exts     []string
		err      bool
	}{
		{"/home/user/.config/data.tar.gz", []string{"gz", "tar"}, false},
		{"archive.20060102.2.zip", []string{"zip"}, false},
		{"files_20060102-150405.bak.tar.xz", []string{"xz", "tar"}, false},
		{"export.tar.lz4", []string{"lz4", "tar"}, false},
		{"takeout-20060102T150405Z-001.tgz", []string{"gz", "tar"}, false},
		{"data.txz", []string{"xz", "tar"}, false},
		{"weird.zip.gz.lz4.xz.gz", []string{"gz", "xz", "lz4", "gz", "zip"}, false},
		{"archive", nil, true},
		{"log.txt.gz", nil, true},
		{".gitignore", nil, true},
	}
	for _, test := range tests {
		exts, err := splitExt(test.filename)
		if (err != nil) != test.err {
			t.Errorf("splitExt(%q) got err: %v, want err: %t", test.filename, err, test.err)
			continue
		}
		if !reflect.DeepEqual(exts, test.exts) {
			t.Errorf("splitExt(%q) = %q, want %q", test.filename, exts, test.exts)
		}
	}
}
