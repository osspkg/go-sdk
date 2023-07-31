/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app

// Config config model
type Config struct {
	Env     string `yaml:"env"`
	Level   uint32 `yaml:"level"`
	LogFile string `yaml:"log"`
}
