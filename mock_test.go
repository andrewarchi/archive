// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"io"
)

type mockFile struct {
	Name string
	Body []byte
}

func makeTar(files []mockFile) ([]byte, error) {
	var b bytes.Buffer
	tw := tar.NewWriter(&b)
	for _, f := range files {
		h := &tar.Header{
			Name:     f.Name,
			Typeflag: tar.TypeReg,
			Mode:     0600,
			Size:     int64(len(f.Body)),
		}
		if err := tw.WriteHeader(h); err != nil {
			return nil, err
		}
		if _, err := io.Copy(tw, bytes.NewReader(f.Body)); err != nil {
			return nil, err
		}
	}
	tw.Close()
	return b.Bytes(), nil
}

func makeZip(files []mockFile) ([]byte, error) {
	var b bytes.Buffer
	zw := zip.NewWriter(&b)
	for _, f := range files {
		fw, err := zw.Create(f.Name)
		if err != nil {
			return nil, err
		}
		if _, err := io.Copy(fw, bytes.NewReader(f.Body)); err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}
