// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

type testMakeArchive func(w io.Writer, files []mockFile) error
type testWalker func(b *bytes.Buffer, filename string, walk WalkFunc) error

func testWalk(t *testing.T, filename string, makeArchive testMakeArchive, walker testWalker) {
	var b bytes.Buffer
	if err := makeArchive(&b, testFiles); err != nil {
		t.Fatal(err)
	}
	i := -1
	var walkErr error
	err := walker(&b, filename, func(f File) error {
		i++
		if walkErr != nil {
			t.Errorf("#%d: traversed after err: %q", i, f.Name())
		}
		if i > len(testFiles) {
			t.Errorf("#%d: traversed after %d files: %q", i, len(testFiles), f.Name())
			return nil
		}
		if f.Name() != testFiles[i].Name {
			t.Errorf("#%d: got name %q, want %q", i, f.Name(), testFiles[i].Name)
			return nil
		}
		fr, err := f.Open()
		if err != nil {
			walkErr = err
			return err
		}
		defer func() {
			if err := fr.Close(); err != nil {
				walkErr = err
			}
		}()
		var b bytes.Buffer
		if _, err := io.Copy(&b, fr); err != nil {
			walkErr = err
			return err
		}
		if body := b.String(); body != testFiles[i].Body {
			t.Errorf("#%d: got body %s, want %s", i, body, testFiles[i].Body)
		}
		return nil
	})
	if walkErr == nil && err != nil {
		t.Error(err)
	}
	if walkErr != nil && err != walkErr {
		t.Errorf("returned err %v not consistent with %v", err, walkErr)
	}
	if i+1 != len(testFiles) {
		t.Errorf("traversed %d files, want %d", i+1, len(testFiles))
	}
}

func TestWalkZip(t *testing.T) {
	walker := func(b *bytes.Buffer, filename string, walk WalkFunc) error {
		return WalkZip(bytes.NewReader(b.Bytes()), int64(b.Len()), filename, walk)
	}
	testWalk(t, "archive.zip", makeZip, walker)
}

func TestWalkTar(t *testing.T) {
	walker := func(b *bytes.Buffer, filename string, walk WalkFunc) error {
		return WalkTar(b, filename, walk)
	}
	testWalk(t, "archive.tar", makeTar, walker)
}

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
