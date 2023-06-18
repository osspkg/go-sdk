/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app

import (
	"os"

	"github.com/osspkg/go-sdk/log"
)

type _log struct {
	file    *os.File
	handler log.Logger
	conf    *Config
}

func newLog(conf *Config) *_log {
	file, err := os.OpenFile(conf.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return &_log{file: file, conf: conf}
}

func (v *_log) Handler(l log.Logger) {
	v.handler = l
	v.handler.SetOutput(v.file)
	v.handler.SetLevel(v.conf.Level)
}

func (v *_log) Close() error {
	if v.handler != nil {
		v.handler.Close()
	}
	return v.file.Close()
}
