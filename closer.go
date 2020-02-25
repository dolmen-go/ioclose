/*
Copyright 2020 Olivier Mengué

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package ioclose provides utilities for objects implementing or using interface io.Closer.
package ioclose

// Closer is the same as io.Closer.
type Closer interface {
	Close() error
}

// CloserFunc allows to convert a simple func() error into an io.Closer
type CloserFunc func() error

func (fn CloserFunc) Close() error {
	if fn == nil {
		return nil
	}
	return fn()
}

// CloseAll calls the .Close method on each argument and returns the first non-nil error.
//
// This helper can be used for example to close multiple prepared statements.
func CloseAll(closers ...interface {
	Close() error
}) (err error) {
	for _, c := range closers {
		if c == nil {
			continue
		}
		e := c.Close()
		if e != nil && err == nil {
			err = e
		}
	}
	return
}

// Closers is a list of closers to call all closers at once.
type Closers struct {
	closers []func() error
}

// Close calls all closers in reverse order.
//
// The first non-nil error is returned.
func (c *Closers) Close() (err error) {
	for i := len(c.closers) - 1; i >= 0; i-- {
		f := c.closers[i]
		if f == nil {
			continue
		}
		e := f()
		if e != nil && err == nil {
			err = e
		}
	}
	c.closers = nil
	return
}

// Close calls all closers in reverse order.
//
// After a call to CloseDefered is defered you can still safely Append more closers.
//
// If *perr is nil, it will contain the result of the first non-nil Close error.
func (c *Closers) CloseDefered(perr *error) {
	err := c.Close()
	if perr != nil && *perr == nil {
		*perr = err
	}
}

func (c *Closers) Append(closers ...interface {
	Close() error
}) {
	for _, cl := range closers {
		if cl == nil {
			continue
		}
		c.closers = append(c.closers, cl.Close)
	}
}

func (c *Closers) AppendFunc(closers ...func() error) {
	for _, f := range closers {
		if f == nil {
			continue
		}
		c.closers = append(c.closers, f)
	}
}
