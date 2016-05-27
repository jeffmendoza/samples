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

/* Abstract data for future persistance. */

/* Exmaple

   da = NewDataAccess()
   da.ChangeText("Foo")
   da.AddUser("Bob", 1)
   data = da.Get()
   fmt.Printf("%v", data.Users[0].Name)
   fmt.Printf("%v", data.Users[0].ID)
   fmt.Printf("%v", data.Text)
*/

type User struct {
	Name string
	ID   int
}

type Data struct {
	Users []User
	Text string
}

type DataAccess struct {
	data Data
}

func (da DataAccess) Get() Data {
	return da.data
}

func (da *DataAccess) AddUser(name string, id int) {
	da.data.Users = append(da.data.Users, User{name, id})
}

func (da *DataAccess) ClearUsers() {
	da.data.Users = nil
}

func (da *DataAccess) ChangeText(text string) {
	da.data.Text = text
}

func NewDataAccess() DataAccess {
	var da DataAccess
	return da
}
