// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package archive contains traversal utilities for ZIP and tar
// archives with gzip, bzip2, xz, and LZ4 compression.
package archive

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pierrec/lz4/v4"
)

// File exposes a common interface for files in an archive.
type File interface {
	Name() string
	Open() (io.ReadCloser, error)
	FileInfo() os.FileInfo
}

// WalkFunc is the type of function that is called for each file
// visited.
type WalkFunc func(File) error

// Walk traverses an archive from an io.Reader and executes the given
// walk function on each file. Supported archive and compression
// formats: ZIP, tar, gzip, bzip2, xz, and LZ4.
func Walk(r io.Reader, filename string, walk WalkFunc) error {
	return walkReader(r, filename, walk)
}

// WalkFile traverses an archive from a file and executes the given walk
// function on each file. Supported archive and compression formats:
// ZIP, tar, gzip, bzip2, xz, and LZ4.
func WalkFile(filename string, walk WalkFunc) error {
	if strings.HasSuffix(filename, ".zip") {
		return WalkZipFile(filename, walk)
	}
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return walkReader(f, filename, walk)
}

func walkReader(r io.Reader, filename string, walk WalkFunc) error {
	exts, err := splitExt(filename)
	if err != nil {
		return err
	}
	for _, ext := range exts {
		switch ext {
		case "zip":
			b, err := ioutil.ReadAll(r)
			if err != nil {
				return err
			}
			return WalkZip(bytes.NewReader(b), int64(len(b)), filename, walk)
		case "tar":
			return WalkTar(r, filename, walk)
		case "gz":
			gr, err := gzip.NewReader(r)
			if err != nil {
				return err
			}
			defer gr.Close()
			r = gr
		case "bz2":
			r = bzip2.NewReader(r)
		case "xz":
			xr, err := NewXZReader(r)
			if err != nil {
				return err
			}
			defer xr.Close()
			r = xr
		case "lz4":
			r = lz4.NewReader(r)
		default:
			return fmt.Errorf("archive: unsupported format: %q", ext)
		}
	}
	return fmt.Errorf("archive: no archive extension: %s", filename)
}

// splitExt splits a filename with compound extensions into recognized
// extensions.
func splitExt(filename string) ([]string, error) {
	name := filename
	var exts []string
	for {
		switch ext := filepath.Ext(name); ext {
		case ".zip", ".tar":
			return append(exts, ext[1:]), nil
		case ".tgz", ".tbz2", ".txz":
			return append(exts, ext[2:], "tar"), nil
		case ".gz", ".bz2", ".xz", ".lz4":
			exts = append(exts, ext[1:])
			name = name[:len(name)-len(ext)]
		case "":
			return nil, fmt.Errorf("archive: no archive extension in %q", filename)
		default:
			return nil, fmt.Errorf("archive: unrecognized extension %q in %q", ext, filename)
		}
	}
}
