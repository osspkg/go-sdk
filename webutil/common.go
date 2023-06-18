/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package webutil

import (
	"strings"
	"time"

	"github.com/osspkg/go-sdk/errors"
)

const (
	statusOn  = 1
	statusOff = 0

	defaultTimeout         = 10 * time.Second
	defaultShutdownTimeout = 1 * time.Second
	defaultNetwork         = "tcp"
)

var (
	errServAlreadyRunning = errors.New("server already running")
	errServAlreadyStopped = errors.New("server already stopped")
	errEpollEmptyEvents   = errors.New("epoll empty event")
	errFailContextKey     = errors.New("context key is not found")
)

var (
	defaultEOF = []byte("\r\n")

	networkType = map[string]struct{}{
		"tcp":        {},
		"tcp4":       {},
		"tcp6":       {},
		"unixpacket": {},
		"unix":       {},
	}
)

/**********************************************************************************************************************/

const urlSplitSeparate = "/"

func urlSplit(uri string) []string {
	vv := strings.Split(strings.ToLower(uri), urlSplitSeparate)
	for i := 0; i < len(vv); i++ {
		if len(vv[i]) == 0 {
			copy(vv[i:], vv[i+1:])
			vv = vv[:len(vv)-1]
			i--
		}
	}
	return vv
}
