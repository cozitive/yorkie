/*
 * Copyright 2022 The Yorkie Authors. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package crdt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	t.Run("escape normal string", func(t *testing.T) {
		str := `hello world`
		expected := `hello world`
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})

	t.Run("escape string with doublequote", func(t *testing.T) {
		str := `hello world"`
		expected := `hello world\\"`
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})

	t.Run("escape string with backslash", func(t *testing.T) {
		str := `hello world\`
		expected := `hello world\\\\`
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})

	t.Run("escape string with control character in switch case", func(t *testing.T) {
		str := "hello world\n"
		expected := `hello world\\n`
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})

	t.Run("escape string with control character not in switch case", func(t *testing.T) {
		str := "hello world\u0000"
		expected := `hello world\\u0000`
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})

	t.Run("escape string with unicode character", func(t *testing.T) {
		str := "hello world\u1234"
		expected := "hello world\u1234"
		actual := EscapeString(str)
		assert.Equal(t, expected, actual)
	})
}
