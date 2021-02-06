// Copyright (c) 2021 Andrew Archibald
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package archive

import (
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/ulikunitz/xz"
)

// NewXZReader returns a reader that decompresses XZ data using system
// xz, if in PATH, otherwise falling back to a slower Go implementation.
func NewXZReader(r io.Reader) (io.ReadCloser, error) {
	if _, err := exec.LookPath("xz"); err != nil {
		xr, err := xz.NewReader(r)
		if err != nil {
			return nil, err
		}
		return ioutil.NopCloser(xr), nil
	}

	rpipe, wpipe := io.Pipe()
	cmd := exec.Command("xz", "--decompress", "--stdout")
	cmd.Stdin = r
	cmd.Stdout = wpipe
	go func() {
		err := cmd.Run()
		wpipe.CloseWithError(err)
	}()
	return rpipe, nil
}
