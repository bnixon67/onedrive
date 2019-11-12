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

package onedrive

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/oauth2"
)

// randomBytesBase64 returns n bytes encoded in URL friendly base64.
func randomBytesBase64(n int) string {
	// buffer to store n bytes
	b := make([]byte, n)

	// get b random bytes
	_, err := rand.Read(b)
	if err != nil {
		log.Panic(err)
	}

	// convert to URL friendly base64
	return base64.URLEncoding.EncodeToString(b)
}

// readTokenFromFile reads the json encoded token from a file.
func readTokenFromFile(filename string) (*oauth2.Token, error) {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read json encoded token
	token := &oauth2.Token{}
	err = json.NewDecoder(file).Decode(token)

	return token, err
}

// writeTokenToFile writes a josn encoded token to a file.
// If file already exists, it is replaced.
func writeTokenToFile(fileName string, token *oauth2.Token) {
	// create file
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// write access token string
	json.NewEncoder(file).Encode(token)

	return
}
