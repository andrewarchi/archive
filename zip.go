// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
)

type zipFile struct {
	f *zip.File
}

func (zf zipFile) Name() string                 { return zf.f.Name }
func (zf zipFile) Open() (io.ReadCloser, error) { return zf.f.Open() }
func (zf zipFile) FileInfo() os.FileInfo        { return zf.f.FileInfo() }

func walkZip(zr *zip.Reader, filename string, walk WalkFunc) error {
	for _, f := range zr.File {
		if err := walk(zipFile{f}); err != nil {
			return fmt.Errorf("archive: walk %s:%s: %w", filename, f.Name, err)
		}
	}
	return nil
}

// WalkZip traverses a ZIP archive from an io.ReaderAt and executes the
// given walk function on each file.
func WalkZip(r io.ReaderAt, size int64, filename string, walk WalkFunc) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return err
	}
	return walkZip(zr, filename, walk)
}

// WalkZipFile traverses a ZIP archive from a file and executes the
// given walk function on each file.
func WalkZipFile(filename string, walk WalkFunc) error {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer zr.Close()
	return walkZip(&zr.Reader, filename, walk)
}

// OpenSingleFileZip opens a zip containing a single file for reading
// and returns the filename of the contained file.
func OpenSingleFileZip(filename string) (io.ReadCloser, string, error) {
	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, "", err
	}
	if len(zr.File) != 1 {
		return nil, "", fmt.Errorf("archive: zip has %d files: %q", len(zr.File), filename)
	}
	f, err := zr.File[0].Open()
	if err != nil {
		return nil, "", err
	}
	return &singleFileZipReader{zr, f}, zr.File[0].Name, nil
}

type singleFileZipReader struct {
	zr *zip.ReadCloser
	f  io.ReadCloser
}

func (z *singleFileZipReader) Read(p []byte) (n int, err error) {
	return z.f.Read(p)
}

func (z *singleFileZipReader) Close() error {
	err1 := z.f.Close()
	err2 := z.zr.Close()
	if err1 != nil {
		return err1
	}
	return err2
}
