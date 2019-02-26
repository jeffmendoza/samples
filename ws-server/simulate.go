// Copyright 2016 Google, Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"time"
	"fmt"
)

/* Data change simulator */

func send (cons *[](chan bool)) {
	fmt.Printf("Simulated data changed\n")
	fmt.Printf("Sending to num of connections: %v\n", len(*cons))

	for _, con := range *cons {
		select {
		case con <- true:
		default:
		}
	}
}

func Simulate(da *DataAccess, cons *[](chan bool)) {
	for {
		da.ClearUsers()
		da.AddUser("First", 1)
		da.ChangeText("One")
		send(cons)
		time.Sleep(5 * time.Second)

		da.AddUser("Second", 2)
		da.ChangeText("Two")
		send(cons)
		time.Sleep(5 * time.Second)

		da.AddUser("Third", 3)
		da.ChangeText("Three")
		send(cons)
		time.Sleep(5 * time.Second)

		da.AddUser("Fourth", 4)
		da.ChangeText("Four")
		send(cons)
		time.Sleep(5 * time.Second)
	}
}
