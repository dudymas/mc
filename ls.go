/*
 * Minio Client (C) 2015 Minio, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/minio/mc/pkg/client"
	"github.com/minio/mc/pkg/console"
	"github.com/minio/minio/pkg/probe"
)

/// ls - related internal functions

const (
	printDate = "2006-01-02 15:04:05 MST"
)

// ContentMessage container for content message structure.
type ContentMessage struct {
	Filetype string    `json:"type"`
	Time     time.Time `json:"lastModified"`
	Size     int64     `json:"size"`
	Name     string    `json:"name"`
}

// String colorized string message
func (c ContentMessage) String() string {
	message := console.Colorize("Time", fmt.Sprintf("[%s] ", c.Time.Format(printDate)))
	message = message + console.Colorize("Size", fmt.Sprintf("%6s ", humanize.IBytes(uint64(c.Size))))
	message = func() string {
		if c.Filetype == "folder" {
			return message + console.Colorize("Dir", fmt.Sprintf("%s", c.Name))
		}
		return message + console.Colorize("File", fmt.Sprintf("%s", c.Name))
	}()
	return message
}

// JSON jsonified content message
func (c ContentMessage) JSON() string {
	jsonMessageBytes, e := json.Marshal(c)
	fatalIf(probe.NewError(e), "Unable to marshal into JSON.")

	return string(jsonMessageBytes)
}

// parseContent parse client Content container into printer struct.
func parseContent(c *client.Content) ContentMessage {
	content := ContentMessage{}
	content.Time = c.Time.Local()

	// guess file type
	content.Filetype = func() string {
		if c.Type.IsDir() {
			return "folder"
		}
		return "file"
	}()

	content.Size = c.Size
	// Convert OS Type to match console file printing style.
	content.Name = func() string {
		switch {
		case runtime.GOOS == "windows":
			c.Name = strings.Replace(c.Name, "/", "\\", -1)
			c.Name = strings.TrimSuffix(c.Name, "\\")
		default:
			c.Name = strings.TrimSuffix(c.Name, "/")
		}
		if c.Type.IsDir() {
			switch {
			case runtime.GOOS == "windows":
				return fmt.Sprintf("%s\\", c.Name)
			default:
				return fmt.Sprintf("%s/", c.Name)
			}
		}
		return c.Name
	}()
	return content
}

// doList - list all entities inside a folder.
func doList(clnt client.Client, recursive, multipleArgs bool) *probe.Error {
	var err *probe.Error
	var parentContent *client.Content
	parentContent, err = clnt.Stat()
	if err != nil {
		return err.Trace(clnt.URL().String())
	}
	for contentCh := range clnt.List(recursive) {
		if contentCh.Err != nil {
			switch contentCh.Err.ToGoError().(type) {
			// handle this specifically for filesystem
			case client.BrokenSymlink:
				errorIf(contentCh.Err.Trace(), "Unable to list broken link.")
				continue
			case client.TooManyLevelsSymlink:
				errorIf(contentCh.Err.Trace(), "Unable to list too many levels link.")
				continue
			}
			if os.IsNotExist(contentCh.Err.ToGoError()) || os.IsPermission(contentCh.Err.ToGoError()) {
				if contentCh.Content != nil {
					if contentCh.Content.Type.IsDir() && (contentCh.Content.Type&os.ModeSymlink == os.ModeSymlink) {
						errorIf(contentCh.Err.Trace(), "Unable to list broken folder link.")
						continue
					}
				}
				errorIf(contentCh.Err.Trace(), "Unable to list.")
				continue
			}
			err = contentCh.Err.Trace()
			break
		}
		if multipleArgs && parentContent.Type.IsDir() {
			contentCh.Content.Name = filepath.Join(parentContent.Name, strings.TrimPrefix(contentCh.Content.Name, parentContent.Name))
		}
		Prints("%s\n", parseContent(contentCh.Content))
	}
	if err != nil {
		return err.Trace()
	}
	return nil
}
