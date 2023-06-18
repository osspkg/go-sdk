/*
 *  Copyright (c) 2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package app_test

import "testing"

type (
	//Simple model
	Simple struct{}
	//Config model
	Config1 struct {
		Env string `yaml:"env"`
	}
	Config2 struct {
		Env string `yaml:"env"`
	}
)

func TestUnit_Invoke(t *testing.T) {

}
