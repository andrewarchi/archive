// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"io"
	"os"
	"path/filepath"
)

// Extract returns a WalkFunc that extracts the files in an archive to
// the given directory.
func Extract(filename, dir string) WalkFunc {
	return func(f File) error {
		out := filepath.Join(dir, f.Name())
		if f.FileInfo().IsDir() {
			return os.MkdirAll(out, 0700)
		}
		r, err := f.Open()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(out), 0700); err != nil {
			return err
		}
		w, err := os.Create(out)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, r)
		return err
	}
}
