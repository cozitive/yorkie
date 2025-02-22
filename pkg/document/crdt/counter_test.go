/*
 * Copyright 2020 The Yorkie Authors. All rights reserved.
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

package crdt_test

import (
	"math"
	"strconv"
	"testing"
	gotime "time"

	"github.com/stretchr/testify/assert"

	"github.com/yorkie-team/yorkie/pkg/document/crdt"
	"github.com/yorkie-team/yorkie/pkg/document/time"
)

func TestCounter(t *testing.T) {
	t.Run("new counter test", func(t *testing.T) {
		intCntWithInt32Value := crdt.NewCounter(crdt.IntegerCnt, int32(math.MaxInt32), time.InitialTicket)
		assert.Equal(t, crdt.IntegerCnt, intCntWithInt32Value.ValueType())

		intCntWithInt64Value := crdt.NewCounter(crdt.IntegerCnt, int64(math.MaxInt32+1), time.InitialTicket)
		assert.Equal(t, crdt.IntegerCnt, intCntWithInt64Value.ValueType())

		intCntWithIntValue := crdt.NewCounter(crdt.IntegerCnt, math.MaxInt32, time.InitialTicket)
		assert.Equal(t, crdt.IntegerCnt, intCntWithIntValue.ValueType())

		intCntWithDoubleValue := crdt.NewCounter(crdt.IntegerCnt, 0.5, time.InitialTicket)
		assert.Equal(t, crdt.IntegerCnt, intCntWithDoubleValue.ValueType())

		intCntWithUnsupportedValue := func() { crdt.NewCounter(crdt.IntegerCnt, "", time.InitialTicket) }
		assert.Panics(t, intCntWithUnsupportedValue)

		longCntWithInt32Value := crdt.NewCounter(crdt.LongCnt, int32(math.MaxInt32), time.InitialTicket)
		assert.Equal(t, crdt.LongCnt, longCntWithInt32Value.ValueType())

		longCntWithInt64Value := crdt.NewCounter(crdt.LongCnt, int64(math.MaxInt32+1), time.InitialTicket)
		assert.Equal(t, crdt.LongCnt, longCntWithInt64Value.ValueType())

		longCntWithIntValue := crdt.NewCounter(crdt.LongCnt, math.MaxInt32+1, time.InitialTicket)
		assert.Equal(t, crdt.LongCnt, longCntWithIntValue.ValueType())

		longCntWithDoubleValue := crdt.NewCounter(crdt.LongCnt, 0.5, time.InitialTicket)
		assert.Equal(t, crdt.LongCnt, longCntWithDoubleValue.ValueType())

		longCntWithUnsupportedValue := func() { crdt.NewCounter(crdt.LongCnt, "", time.InitialTicket) }
		assert.Panics(t, longCntWithUnsupportedValue)
	})

	t.Run("increase test", func(t *testing.T) {
		var x = 5
		var y int64 = 10
		var z = 3.14
		integer := crdt.NewCounter(crdt.IntegerCnt, x, time.InitialTicket)
		long := crdt.NewCounter(crdt.LongCnt, y, time.InitialTicket)
		double := crdt.NewCounter(crdt.IntegerCnt, z, time.InitialTicket)

		integerOperand := crdt.NewPrimitive(x, time.InitialTicket)
		longOperand := crdt.NewPrimitive(y, time.InitialTicket)
		doubleOperand := crdt.NewPrimitive(z, time.InitialTicket)

		// normal process test
		integer.Increase(integerOperand)
		integer.Increase(longOperand)
		integer.Increase(doubleOperand)
		assert.Equal(t, integer.Marshal(), "23")

		long.Increase(integerOperand)
		long.Increase(longOperand)
		long.Increase(doubleOperand)
		assert.Equal(t, long.Marshal(), "28")

		double.Increase(integerOperand)
		double.Increase(longOperand)
		double.Increase(doubleOperand)
		assert.Equal(t, double.Marshal(), "21")

		// error process test
		// TODO: it should be modified to error check
		// when 'Remove panic from server code (#50)' is completed.
		unsupportedTypePanicTest := func() {
			r := recover()
			assert.NotNil(t, r)
			assert.Equal(t, r, "unsupported type")
		}
		unsupportedTest := func(v interface{}) {
			defer unsupportedTypePanicTest()
			crdt.NewCounter(crdt.IntegerCnt, v, time.InitialTicket)
		}
		unsupportedTest("str")
		unsupportedTest(true)
		unsupportedTest([]byte{2})
		unsupportedTest(gotime.Now())

		assert.Equal(t, integer.Marshal(), "23")
		assert.Equal(t, long.Marshal(), "28")
		assert.Equal(t, double.Marshal(), "21")
	})

	t.Run("Counter value overflow test", func(t *testing.T) {
		integer := crdt.NewCounter(crdt.IntegerCnt, math.MaxInt32, time.InitialTicket)
		assert.Equal(t, integer.ValueType(), crdt.IntegerCnt)

		operand := crdt.NewPrimitive(1, time.InitialTicket)
		integer.Increase(operand)
		assert.Equal(t, integer.ValueType(), crdt.IntegerCnt)
		assert.Equal(t, integer.Marshal(), strconv.FormatInt(math.MinInt32, 10))
	})
}
