/*
Copyright 2019 Bill Nixon

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published
by the Free Software Foundation, either version 3 of the License,
or (at your option) any later version.

This program is distributed in the hope that it will be useful, but
WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/bnixon67/onedrive"
)

func showResponse(v interface{}) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println("error: err")
	}
	fmt.Println("=== BEGIN ===")
	fmt.Println(string(b))
	fmt.Println("=== END ===")
}

func main() {
	oneDriveClient := onedrive.New(".token.json")

	resp, err := oneDriveClient.Get("https://graph.microsoft.com/v1.0/me/drive")
	if err != nil {
		log.Fatal(err)
	}

	var out bytes.Buffer
	json.Indent(&out, resp, "", " ")
	fmt.Println("=== BEGIN ===")
	fmt.Println(out.String())
	fmt.Println("=== END ===")
}
