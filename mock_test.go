// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"archive/tar"
	"archive/zip"
	"io"
	"strings"
)

type mockFile struct {
	Name, Body string
}

var testFiles = []mockFile{
	{".git/HEAD", "ref: refs/heads/main\n"},
	{"cmd/hello/main.go", "package main\n\nimport \"fmt\"\n\nfunc hello() {\n\tfmt.Println(\"Hallo Welt!\")\n}\n\nfunc main() {\n\thello()\n}\n"},
	{"cmd/hello/main_test.go", "package main\n\nfunc TestHello(t *testing.T) {\n\tif h := hello(); h != \"Hallo Welt!\" {\n\t\tt.Error(\"bad salutation\")\n\t}\n}\n"},
	{"README.md", "# Project\n\nHello, World!\n"},
	{"LICENSE", "Copyright (c) 2021 Acme Corp.\n"},
	{"go.mod", "module example.com/hello\n\ngo1.16"},
}

func makeZip(w io.Writer, files []mockFile) error {
	zw := zip.NewWriter(w)
	defer zw.Close()
	for _, f := range files {
		fw, err := zw.Create(f.Name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(fw, strings.NewReader(f.Body)); err != nil {
			return err
		}
	}
	return nil
}

func makeTar(w io.Writer, files []mockFile) error {
	tw := tar.NewWriter(w)
	defer tw.Close()
	for _, f := range files {
		h := &tar.Header{
			Name:     f.Name,
			Typeflag: tar.TypeReg,
			Mode:     0600,
			Size:     int64(len(f.Body)),
		}
		if err := tw.WriteHeader(h); err != nil {
			return err
		}
		if _, err := io.Copy(tw, strings.NewReader(f.Body)); err != nil {
			return err
		}
	}
	return nil
}
