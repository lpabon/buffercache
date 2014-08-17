//
// Copyright (c) 2014 The buffercache Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buffercache

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewClockCache(t *testing.T) {
	cachesize := uint64(1024 * 1024)
	blocksize := uint64(1024)
	b := NewClockCache(cachesize, blocksize)
	assert.Equal(t, int(cachesize/blocksize), len(b.cacheblocks))
}

func TestClockSimple(t *testing.T) {
	var err error
	var buffer = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tmpbuffer := make([]byte, len(buffer))
	key := uint64(1)

	b := NewClockCache(uint64(2*len(buffer)), uint64(len(buffer)))
	err = b.Set(key, buffer)
	assert.NoError(t, err)

	err = b.Get(key, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer, tmpbuffer))

	err = b.Invalidate(key)
	assert.NoError(t, err)
	err = b.Get(key, tmpbuffer)
	assert.Error(t, err)
}

func TestClockCopiedBuffers(t *testing.T) {
	var err error
	var buffer = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	key := uint64(1)

	b := NewClockCache(uint64(2*len(buffer)), uint64(len(buffer)))
	err = b.Set(key, buffer)
	assert.NoError(t, err)

	tmpbuffer := make([]byte, len(buffer))
	err = b.Get(key, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer, tmpbuffer))

	buffer[0] = 0
	buffer[1] = 0
	tmpbuffer = make([]byte, len(buffer))
	err = b.Get(key, tmpbuffer)
	assert.NoError(t, err)
	assert.False(t, reflect.DeepEqual(buffer, tmpbuffer))
	assert.Equal(t, byte(1), tmpbuffer[0])
	assert.Equal(t, byte(2), tmpbuffer[1])
}

func TestClockEvictions(t *testing.T) {
	var err error
	var buffer1 = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var buffer2 = []byte{0, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var buffer3 = []byte{0, 0, 3, 4, 5, 6, 7, 8, 9, 10}
	var key1 = uint64(1)
	var key2 = uint64(2)
	var key3 = uint64(3)
	tmpbuffer := make([]byte, len(buffer1))

	b := NewClockCache(uint64(2*len(buffer1)), uint64(len(buffer1)))
	err = b.Set(key1, buffer1)
	assert.NoError(t, err)
	err = b.Set(key2, buffer2)
	assert.NoError(t, err)

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key1, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer1, tmpbuffer))

	err = b.Set(key3, buffer3)
	assert.NoError(t, err)

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key1, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer1, tmpbuffer))

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key2, tmpbuffer)
	assert.Error(t, err)

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key3, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer3, tmpbuffer))

}

func TestClockDoubleSets(t *testing.T) {
	var err error
	var buffer1 = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var buffer2 = []byte{0, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	var buffer3 = []byte{0, 0, 3, 4, 5, 6, 7, 8, 9, 10}
	var key1 = uint64(1)
	var key2 = uint64(2)
	tmpbuffer := make([]byte, len(buffer1))

	b := NewClockCache(uint64(2*len(buffer1)), uint64(len(buffer1)))
	err = b.Set(key1, buffer1)
	assert.NoError(t, err)
	err = b.Set(key2, buffer2)
	assert.NoError(t, err)

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key1, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer1, tmpbuffer))

	err = b.Set(key1, buffer3)
	assert.NoError(t, err)

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key1, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer3, tmpbuffer))

	tmpbuffer = make([]byte, len(buffer1))
	err = b.Get(key2, tmpbuffer)
	assert.NoError(t, err)
	assert.True(t, reflect.DeepEqual(buffer2, tmpbuffer))
}
