// Copyright 2024 gofrp (https://github.com/gofrp)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"io"
	"sync"
)

var (
	bufPool16k sync.Pool
	bufPool5k  sync.Pool
	bufPool2k  sync.Pool
	bufPool1k  sync.Pool
	bufPool    sync.Pool
)

func GetBuf(size int) []byte {
	var x interface{}
	if size >= 16*1024 {
		x = bufPool16k.Get()
	} else if size >= 5*1024 {
		x = bufPool5k.Get()
	} else if size >= 2*1024 {
		x = bufPool2k.Get()
	} else if size >= 1*1024 {
		x = bufPool1k.Get()
	} else {
		x = bufPool.Get()
	}
	if x == nil {
		return make([]byte, size)
	}
	buf := x.([]byte)
	if cap(buf) < size {
		return make([]byte, size)
	}
	return buf[:size]
}

func PutBuf(buf []byte) {
	size := cap(buf)
	if size >= 16*1024 {
		bufPool16k.Put(buf)
	} else if size >= 5*1024 {
		bufPool5k.Put(buf)
	} else if size >= 2*1024 {
		bufPool2k.Put(buf)
	} else if size >= 1*1024 {
		bufPool1k.Put(buf)
	} else {
		bufPool.Put(buf)
	}
}

// Join two io.ReadWriteCloser and do some operations.
func Join(c1 io.ReadWriteCloser, c2 io.ReadWriteCloser) (inCount int64, outCount int64) {
	var wait sync.WaitGroup
	pipe := func(to io.ReadWriteCloser, from io.ReadWriteCloser, count *int64) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		buf := GetBuf(16 * 1024)
		defer PutBuf(buf)
		*count, _ = io.CopyBuffer(to, from, buf)
	}

	wait.Add(2)
	go pipe(c1, c2, &inCount)
	go pipe(c2, c1, &outCount)
	wait.Wait()
	return
}
