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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Identity struct {
	EMail       string `json:"email,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Id          string `json:"id,omitempty"`
}

// using pointer to Identity so Unmarshal creates a nil on empty
// (see https://stackoverflow.com/questions/33447334/golang-json-marshal-how-to-omit-empty-nested-struct)
type IdentitySet struct {
	Application *Identity `json:"application,omitempty"`
	Device      *Identity `json:"device,omitempty"`
	User        *Identity `json:"user,omitempty"`
}

type Quota struct {
	// Total allowed storage space, in bytes. Read-only.
	Total int64 `json:"total,omitempty"`

	// Total space used, in bytes. Read-only.
	Used int64 `json:"used,omitempty"`

	// Total space remaining before reaching the quota limit, in bytes. Read-only.
	Remaining int64 `json:"remaining,omitempty"`

	// Total space consumed by files in the recycle bin, in bytes. Read-only.
	Deleted int64 `json:"deleted,omitempty"`

	// Enumeration value that indicates the state of the storage space. Read-only.
	// normal	The drive has plenty of remaining quota left.
	// nearing	Remaining quota is less than 10% of total quota space.
	// critical	Remaining quota is less than 1% of total quota space.
	// exceeded	The used quota has exceeded the total quota.
	// 		New files or folders cannot be added to the drive until it is under
	// 		the total quota amount or more storage space is purchased.
	State string `json:"state,omitempty"`
}

type Drive struct {
	// The unique identifier of the drive. Read-only.
	Id string `json:"id"`

	// Identity of the user, device, or application which created the item. Read-only.
	CreatedBy *IdentitySet `json:"createdBy,omitempty"`

	// Identity of the user, device, or application which created the item. Read-only.
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// Provide a user-visible description of the drive. Read-write.
	Description string `json:"description,omitempty"`

	// Describes the type of drive represented by this resource. Read-only.
	// 	OneDrive personal drives will return personal.
	// 	OneDrive for Business will return business.
	// 	SharePoint document libraries will return documentLibrary.
	DriveType string `json:"driveType"`

	// Identity of the user, device, and application which last modified the item. Read-only.
	LastModifiedBy *IdentitySet `json:"lastModifiedBy,omitempty"`

	// Date and time the item was last modified. Read-only.
	LastModifiedDateTime string `json:"LastModifiedDateTime,omitempty"`

	// The name of the item. Read-write.
	Name string `json:"name,omitempty"`

	// Optional. The user account that owns the drive. Read-only.
	Owner *IdentitySet `json:"owner,omitempty"`

	// Optional. Information about the drive's storage space quota. Read-only.
	Quota *Quota `json:"quota,omitempty"`

	// URL that displays the resource in the browser. Read-only.
	WebURL string `json:"webUrl,omitempty"`
}

type FileSystemInfo struct {
	CreatedDateTime      string `json:"createdDateTime,omitempty"`
	LastAccessedDateTime string `json:"lastAccessedDateTime,omitempty"`
	LastModifiedDateTime string `json:"lastModifiedDateTime,omitempty"`
}

type ParentReference struct {
	// Unique identifier of the drive instance that contains the item. Read-only.
	DriveId string `json:"driveID,omitempty"`

	// Identifies the type of drive. See drive resource for values.
	DriveType string `json:"driveType,omitempty"`

	Id string `json:"id,omitempty"`
}

type Package struct {
	// A string indicating the type of package.
	// While oneNote is the only currently defined value, you should expect
	// other package types to be returned and handle them accordingly.
	Type string `json:"type,omitempty"`
}

type RemoteItem struct {
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// Unique identifier for the remote item in its drive. Read-only.
	Id string `json:"id,omitempty"`

	LastModifiedDateTime string `json:"lastModifiedDateTime,omitempty"`

	// Optional. Filename of the remote item. Read-only.
	Name string `json:"name,omitempty"`

	// Size of the remote item. Read-only.
	Size int64 `json:"size,omitempty"`

	WebDavURL string `json:"webDavUrl,omitempty"`

	// URL that displays the resource in the browser. Read-only.
	WebURL string `json:"webUrl,omitempty"`

	CreatedBy *IdentitySet `json:"createdBy,omitempty"`

	File *File `json:"file,omitempty"`

	// Information about the remote item from the local file system. Read-only.
	FileSystemInfo FileSystemInfo `json:"fileSystemInfo,omitempty"`

	LastModifiedBy *IdentitySet `json:"lastModifiedBy,omitempty"`

	// If present, indicates that this item is a package instead of a folder or file.
	// Packages are treated like files in some contexts and folders in others. Read-only.
	Package *Package `json:"package,omitempty"`

	// Properties of the parent of the remote item. Read-only.
	ParentReference ParentReference `json:"parentReference,omitempty"`

	SharepointIds *SharepointIds `json:"sharepointIds,omitempty"`
}

type File struct {
	// The MIME type for the file.
	// This is determined by logic on the server and might not be
	// the value provided when the file was uploaded. Read-only.
	MimeType string `json:"mimeType,omitempty"`
}

type SharepointIds struct {
	ListId           string `json:"listId,omitempty"`
	ListItemId       string `json:"listItemId,omitempty"`
	ListItemUniqueId string `json:"listItemUniqueId,omitempty"`
	SiteId           string `json:"siteId,omitempty"`
	SiteUrl          string `json:"siteUrl,omitempty"`
	WebId            string `json:"webId,omitempty"`
}

type DriveItem struct {
	// Date and time of item creation. Read-only.
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// The unique identifier of the item within the Drive. Read-only.
	Id string `json:"id,omitempty"`

	// Date and time the item was last modified. Read-only.
	LastModifiedDateTime string `json:"lastModifiedDateTime,omitempty"`

	// The name of the item (filename and extension). Read-write.
	Name string `json:"name,omitempty"`

	// URL that displays the resource in the browser. Read-only.
	WebURL string `json:"webUrl,omitempty"`

	// Size of the remote item. Read-only.
	Size int64 `json:"size,omitempty"`

	// Identity of the user, device, and application which created the item. Read-only.
	CreatedBy *IdentitySet `json:"createdBy,omitempty"`

	// Identity of the user, device, and application which last modified the item. Read-only.
	LastModifiedBy *IdentitySet `json:"lastModifiedBy,omitempty"`

	File *File `json:"file,omitempty"`

	// File system information on client. Read-write.
	FileSystemInfo FileSystemInfo `json:"fileSystemInfo,omitempty"`

	// Remote item data, if the item is shared from a drive other than the one being accessed.
	// Read-only.
	RemoteItem RemoteItem `json:"remoteItem,omitempty"`

	// eTag for the entire item (metadata + content). Read-only.
	ETag string `json:"eTag,omitempty"`

	// Parent information, if the item has a parent. Read-write.
	ParentReference *ParentReference `json:"parentReference,omitempty"`

	SharepointIds *SharepointIds `json:"sharepointIds,omitempty"`
}

type DriveItems struct {
	Value []DriveItem `json:"value"`
}

type Drives struct {
	Value []Drive `json:"value"`
}

func GetMyDrive(client *http.Client) (drive Drive, err error) {
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me/drive")
	if err != nil {
		return drive, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return drive, err
	}

	err = json.Unmarshal(body, &drive)

	return drive, err
}

func GetDriveItem(client *http.Client) (drive Drive, err error) {
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me/drive")
	if err != nil {
		return drive, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return drive, err
	}

	err = json.Unmarshal(body, &drive)

	return drive, err
}

func ListMyDrives(client *http.Client) (drives Drives, err error) {
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me/drives")
	if err != nil {
		return drives, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return drives, err
	}

	err = json.Unmarshal(body, &drives)

	return drives, err
}

func ListRecentFiles(client *http.Client) (driveItems DriveItems, err error) {
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me/drive/recent?$top=1")
	if err != nil {
		return driveItems, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return driveItems, err
	}

	fmt.Println(string(body))

	err = json.Unmarshal(body, &driveItems)

	return driveItems, err
}
