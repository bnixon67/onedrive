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

import "fmt"

type InnerError struct {
	RequestId string `json:"request-id,omitempty"`
	Date      string `json:"date,omitempty"`
}

type Err struct {
	Code       string      `json:"code,omitempty"`
	Message    string      `json:"message,omitempty"`
	InnerError *InnerError `json:"innerError,omitempty"`
}

type RespError struct {
	Err *Err `json:"error,omitempty"`
}

func (e *RespError) Error() string {
	return fmt.Sprintf("Code: %s Message: %s RequestId: %s Date: %s\n",
		e.Err.Code, e.Err.Message,
		e.Err.InnerError.RequestId, e.Err.InnerError.Date)
}

func codeIsError(code int) bool {
	// Microsoft Graph error responses and resource types
	// https://docs.microsoft.com/en-us/graph/errors
	var errorCodes = []int{400, 401, 403, 404, 405, 406, 409, 410,
		411, 412, 413, 415, 416, 422, 423, 429, 500, 501, 503,
		504, 507, 509}

	for _, n := range errorCodes {
		if code == n {
			return true
		}
	}
	return false
}
