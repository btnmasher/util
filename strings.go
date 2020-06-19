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

import (
	"bytes"
	"fmt"
	"sync"
)

// ChunkJoinStrings takes a list of individual parameters and joins them to strings
// separated by sep, limited by the maxlength. For each item, if appending the item
// would breach the maxlength, it instead starts to build a new string. Once all of
// the strings are built, it returns the list of strings.
func ChunkJoinStrings(params []string, maxlength int, sep string) []string {
	var buffer bytes.Buffer
	currlen := 0
	joined := []string{}
	iterate := false

	for i, param := range params {
		// Check if we have enough room to write the item
		if currlen+len(param) < maxlength {
			buffer.WriteString(param)
			currlen += len(param)
		} else { // Not enough room, reiterate for the next item
			iterate = true
		}

		// Check if last item or if we can fit a space
		if i+1 < len(params) && currlen+len(sep) < maxlength {
			buffer.WriteString(sep)
			currlen++
		} else { // Not enough room, reiterate for the next item
			iterate = true
		}

		if iterate {
			currlen = 0
			iterate = false
			joined = append(joined, buffer.String())
			buffer.Reset()
		}
	}

	if buffer.Len() > 0 { // Finished iterating without hitting max length on the current pass.
		joined = append(joined, buffer.String())
	}

	return joined
}

// ConcurrentMapString is a simple map[string]string wrapped with a concurrent-safe API
type ConcurrentMapString struct {
	data map[string]string
	sync.RWMutex
}

// NewConcurrentMapString initializes and returns a pointer to a new ConcurrentMapString instance.
func NewConcurrentMapString() *ConcurrentMapString {
	m := &ConcurrentMapString{
		data: make(map[string]string),
	}
	return m
}

// ForEach will call the provided function for each entry in the ConcurrentMapString
func (m *ConcurrentMapString) ForEach(do func(string, string)) {
	m.RLock()
	defer m.RUnlock()

	for key, val := range m.data {
		do(key, val)
	}
}

// Length returns the length of the underlying map.
func (m *ConcurrentMapString) Length() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.data)
}

// Add is used to add a key/value to the map.
// Returns an error if the key already exists.
func (m *ConcurrentMapString) Add(key string, value string) error {
	m.Lock()
	defer m.Unlock()

	_, exists := m.data[key]

	if exists {
		return fmt.Errorf("ConcurrentMapString: Cannot add map entry, key already exists: %q", key)
	}

	m.data[key] = value
	return nil
}

// Del is used to remove a key/value from the map.
// Returns an error if the key does not exist.
func (m *ConcurrentMapString) Del(key string) error {
	m.Lock()
	defer m.Unlock()

	_, exists := m.data[key]

	if !exists {
		return fmt.Errorf("ConcurrentMapString: Cannot delete map entry, key does not exist: %q", key)
	}

	delete(m.data, key)

	return nil
}

// Get is used to get a key/value from the map.
// Returns an error if the key does not exist.
func (m *ConcurrentMapString) Get(key string) (string, error) {
	m.RLock()
	defer m.RUnlock()

	v, exists := m.data[key]

	if !exists {
		return "", fmt.Errorf("ConcurrentMapString: Cannot get map value, key does not exist: %q", key)
	}

	return v, nil
}

// Set is used to change an existing key/value in the map.
// Returns an error if the key does not exist.
func (m *ConcurrentMapString) Set(key string, value string) error {
	m.Lock()
	defer m.Unlock()

	_, exists := m.data[key]

	if !exists {
		return fmt.Errorf("ConcurrentMapString: Cannot set map value, key does not exist: %q", key)
	}

	m.data[key] = value

	return nil
}

// Exists is used by external callers to check if a value
// exists in the map and returns a boolean with the result.
func (m *ConcurrentMapString) Exists(key string) bool {
	m.RLock()
	defer m.RUnlock()

	_, exists := m.data[key]
	return exists
}
