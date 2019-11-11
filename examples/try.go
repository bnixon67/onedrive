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
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"

	"github.com/bnixon67/onedrive"
)

const (
	msBase        = "https://login.microsoftonline.com/common/oauth2/v2.0"
	msAuthURL     = msBase + "/authorize"
	msTokenURL    = msBase + "/token"
	myRedirectURL = "https://login.microsoftonline.com/common/oauth2/nativeclient"
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

func showResponse(v interface{}) {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		fmt.Println("error: err")
	}
	fmt.Printf("==========\n%s\n===========\n", b)
}

func main() {
	ctx := context.Background()

	conf := &oauth2.Config{
		ClientID: "c32f556d-11cc-45ce-9b73-37f701abf48c",
		// TODO: need offline_access? AuthCodeURL offline?
		Scopes: []string{"Files.Read.All", "offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  msAuthURL,
			TokenURL: msTokenURL,
		},
		RedirectURL: myRedirectURL,
	}

	// try to get a token from the file
	token, err := readTokenFromFile("token.json")

	if token == nil {
		// could not get token from file

		// generate random state to detect Cross-Site Request Forgery
		state := randomBytesBase64(32)

		// get authentication URL for offline access
		authURL := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)

		// instruct the user to vist the authentication URL
		fmt.Println("Vist the following URL in a browser to authenticate this application")
		fmt.Println("After authentication, copy the response URL from the browser")
		fmt.Println(authURL)

		// read the response URL
		fmt.Println("Enter the response URL:")
		responseString, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		responseString = strings.TrimSpace(responseString)

		// parse the response URL
		responseURL, err := url.Parse(responseString)
		if err != nil {
			log.Fatal(err)
		}

		// get and compare state to prevent Cross-Site Request Forgery
		responseState := responseURL.Query().Get("state")
		if responseState != state {
			log.Fatalln("state mismatch, potenial Cross-Site Request Forgery (CSRF)")
		}

		// get authorization code
		code := responseURL.Query().Get("code")

		// exchange authorize code for token
		token, err = conf.Exchange(ctx, code)
		if err != nil {
			log.Fatal(err)
		}

		// save the token to a file
		writeTokenToFile("token.json", token)
	}

	// create HTTP client using the provided token
	client := conf.Client(ctx, token)

	/*
		drive, err := onedrive.GetMyDrive(client)
		if err != nil {
			log.Fatal(err)
		}
		showResponse(drive)
	*/

	/*
		drives, err := onedrive.ListMyDrives(client)
		if err != nil {
			log.Fatal(err)
		}
		showResponse(drives)

	*/
	recentFiles, err := onedrive.ListRecentFiles(client)
	if err != nil {
		log.Fatal(err)
	}
	showResponse(recentFiles)

	/*
		resp, err := client.Get("https://graph.microsoft.com/v1.0/me/drive")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(body))
	*/
}
