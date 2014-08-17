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
	"github.com/lpabon/godbc"
	"sync"
)

type ClockCacheBlock struct {
	key  uint64
	mru  bool
	used bool
	data []byte
}

type ClockCache struct {
	cacheblocks []ClockCacheBlock
	keymap      map[uint64]uint64
	index       uint64
	lock        sync.Mutex
}

func NewClockCache(cachesize, blocksize uint64) *ClockCache {

	b := &ClockCache{}
	numblocks := cachesize / blocksize
	b.cacheblocks = make([]ClockCacheBlock, numblocks)
	b.keymap = make(map[uint64]uint64)

	for i := 0; i < len(b.cacheblocks); i++ {
		b.cacheblocks[i].data = make([]byte, blocksize)
		godbc.Check(b.cacheblocks[i].data != nil)
	}

	godbc.Ensure(b.cacheblocks != nil)
	godbc.Ensure(b.keymap != nil)

	return b
}

func (c *ClockCache) remove(index uint64) {
	delete(c.keymap, c.cacheblocks[index].key)

	c.cacheblocks[index].mru = false
	c.cacheblocks[index].used = false
	c.cacheblocks[index].key = 0
}

func (c *ClockCache) Set(key uint64, buf []byte) (err error) {

	c.lock.Lock()
	defer c.lock.Unlock()

	// Yes i know its the same as Invalides.. I'll fix it later!
	// :-)
	if index, ok := c.keymap[key]; ok {
		c.remove(index)
	}

	for {
		for ; c.index < uint64(len(c.cacheblocks)); c.index++ {
			if c.cacheblocks[c.index].mru {
				c.cacheblocks[c.index].mru = false
			} else {
				if c.cacheblocks[c.index].used {
					c.remove(c.index)
				}

				err = nil

				c.cacheblocks[c.index].key = key
				c.cacheblocks[c.index].mru = false
				c.cacheblocks[c.index].used = true

				godbc.Check(len(buf) == len(c.cacheblocks[c.index].data))
				copy(c.cacheblocks[c.index].data, buf)

				c.keymap[key] = c.index
				c.index++

				return
			}
		}
		c.index = 0
	}
}

func (c *ClockCache) Get(key uint64, buf []byte) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if index, ok := c.keymap[key]; ok {
		c.cacheblocks[index].mru = true
		copy(buf, c.cacheblocks[index].data)

		return nil
	} else {
		return ErrKeyNotFound
	}
}

func (c *ClockCache) Invalidate(key uint64) (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if index, ok := c.keymap[key]; ok {
		c.remove(index)
	}

	return nil
}
