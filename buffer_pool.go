/*
   Copyright (c) 2020, btnmasher
   All rights reserved.

   Redistribution and use in source and binary forms, with or without modification, are permitted provided that
   the following conditions are met:

   1. Redistributions of source code must retain the above copyright notice, this list of conditions and the
      following disclaimer.

   2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and
      the following disclaimer in the documentation and/or other materials provided with the distribution.

   3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or
      promote products derived from this software without specific prior written permission.

   THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED
   WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A
   PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR
   ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED
   TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
   HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
   NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
   POSSIBILITY OF SUCH DAMAGE.
*/

package util

import "bytes"

// BufferPool holds the Buffers in a Channel as a queue.
type BufferPool struct {
	Buffers chan *bytes.Buffer
}

// NewBufferPool creates a new object pool of bytes.Buffer.
func NewBufferPool(max int) *BufferPool {
	return &BufferPool{
		Buffers: make(chan *bytes.Buffer, max),
	}
}

// Warmup fills the BufferPool with the specified number of objects
// up to one below the maximum capacity of the internal channel
func (pool *BufferPool) Warmup(num, length int) {
	for i := 0; i < num; i++ {
		select {
		case pool.Buffers <- &bytes.Buffer{}: // Add the new buffer to the pool.
		default: // We're full now because we got blocked trying to add that buffer.
			return
		}
	}
}

// New takes a Buffer from the pool.
func (pool *BufferPool) New() (buf *bytes.Buffer) {
	select {
	case buf = <-pool.Buffers:
	default:
		buf = &bytes.Buffer{}
	}
	return
}

// Recycle returns a BUffer to the pool.
func (pool *BufferPool) Recycle(buf *bytes.Buffer) {
	buf.Reset()
	select {
	case pool.Buffers <- buf:
	default:
		// let it go, let it go...
	}
}
