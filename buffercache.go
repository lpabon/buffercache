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
	"errors"
)

var (
	ErrKeyNotFound = errors.New("Key not found")
)

type BufferCache interface {
	Get(key uint64, buf []byte) (err error)
	Set(key uint64, buf []byte) (err error)
	Invalidate(key uint64) (err error)
}
