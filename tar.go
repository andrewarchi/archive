// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type tarFile struct {
	r *tar.Reader
	h *tar.Header
}

func (tf tarFile) Name() string                 { return tf.h.Name }
func (tf tarFile) Open() (io.ReadCloser, error) { return ioutil.NopCloser(tf.r), nil }
func (tf tarFile) FileInfo() os.FileInfo        { return tf.h.FileInfo() }

// WalkTar traverses a tar archive from an io.Reader and executes the
// given walk function on each file.
func WalkTar(r io.Reader, filename string, walk WalkFunc) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		if err := walk(tarFile{tr, header}); err != nil {
			return fmt.Errorf("archive: walk %s:%s: %w", filename, header.Name, err)
		}
	}
	return nil
}

// WalkTarFile traverses a tar archive from a file and executes the
// given walk function on each file.
func WalkTarFile(filename string, walk WalkFunc) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	return WalkTar(f, filename, walk)
}
