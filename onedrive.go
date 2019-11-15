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

// Package onedrive provides access to Microsoft OneDrive API via Microsoft Graph
// See https://docs.microsoft.com/en-us/graph/api/resources/onedrive?view=graph-rest-1.0
package onedrive

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func init() {
	// log with date, time, file name, and line number
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func (c *OneDriveClient) Get(url string) (body []byte, err error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}

	if codeIsError(resp.StatusCode) {
		resError := RespError{}

		err = json.Unmarshal(body, &resError)
		if err != nil {
			return nil, err
		}

		return nil, &resError
	}

	return body, err
}

func (c *OneDriveClient) GetMyDrive() (drive Drive, err error) {
	body, err := c.Get("https://graph.microsoft.com/v1.0/me/drive")
	if err != nil {
		return Drive{}, err
	}

	err = json.Unmarshal(body, &drive)

	return drive, err
}

// ListMyDrives retrieve a list of Drives available for the current user
func (c *OneDriveClient) ListMyDrives() (drives Drives, err error) {
	body, err := c.Get("https://graph.microsoft.com/v1.0/me/drives")
	if err != nil {
		return Drives{}, err
	}

	err = json.Unmarshal(body, &drives)

	return drives, err
}

func (c *OneDriveClient) ListRecentFiles() (driveItems DriveItems, err error) {
	body, err := c.Get("https://graph.microsoft.com/v1.0/me/drive/recent")
	if err != nil {
		return DriveItems{}, err
	}

	err = json.Unmarshal(body, &driveItems)

	return driveItems, err
}

type OneDriveClient struct {
	httpClient *http.Client
}

const (
	msBase        = "https://login.microsoftonline.com/common/oauth2/v2.0"
	msAuthURL     = msBase + "/authorize"
	msTokenURL    = msBase + "/token"
	myRedirectURL = "https://login.microsoftonline.com/common/oauth2/nativeclient"
)

// New create an initialized OneDriveClient using the token from tokenFileName.
// If tokenFileName doesn't exist, then a token is requested and saved in the file.
// User interaction is required to request a token for the first time.
func New(tokenFileName string) *OneDriveClient {
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

	client := &OneDriveClient{}

	// try to get a token from the file
	token, _ := readTokenFromFile(tokenFileName)

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
		writeTokenToFile(tokenFileName, token)
	}

	// create HTTP client using the provided token
	client.httpClient = conf.Client(ctx, token)

	return client
}
