/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app

import "github.com/osspkg/go-sdk/errors"

var (
	errDepBuilderNotRunning = errors.New("dependencies builder is not running")
	errDepNotRunning        = errors.New("dependencies are not running yet")
	errServiceUnknown       = errors.New("unknown service")
	errBadFileFormat        = errors.New("is not a supported file format")
)
