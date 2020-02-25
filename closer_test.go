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

package ioclose_test

import (
	"errors"
	"fmt"
	"io"

	"github.com/dolmen-go/ioclose"
)

var (
	// Check bi-directional equivalence
	_ io.Closer      = ioclose.Closer(nil)
	_ ioclose.Closer = io.Closer(nil)
)

type PrintOnClose string

func (p PrintOnClose) Close() error {
	fmt.Println(string(p))
	return nil
}

type ErrorOnClose string

func (e ErrorOnClose) Close() error {
	fmt.Println(string(e))
	return errors.New(string(e))
}

func ExampleCloseAll() {

}

func ExampleClosers_CloseDefered() {
	var err error
	defer func() {
		fmt.Println("error:", err)
	}()

	var all ioclose.Closers
	defer all.CloseDefered(&err)
	all.Append(PrintOnClose("closer1"))
	all.Append(ErrorOnClose("closer2"))
	all.Append(PrintOnClose("closer3"))
	all.AppendFunc(func() error {
		fmt.Println("closer4")
		return nil
	})

	// Output:
	// closer4
	// closer3
	// closer2
	// closer1
	// error: closer2
}
