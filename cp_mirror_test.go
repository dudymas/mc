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
	"os"

	"github.com/minio/mc/pkg/console"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestCopyURLType(c *C) {
	sourceURLs := []string{server.URL + "/bucket/object1"}
	targetURL := server.URL + "/bucket/test"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeA)

	sourceURLs = []string{server.URL + "/bucket/object1"}
	targetURL = server.URL + "/bucket"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeB)

	sourceURLs = []string{server.URL + "/bucket/..."}
	targetURL = server.URL + "/bucket"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeC)

	sourceURLs = []string{server.URL + "/bucket/...", server.URL + "/bucket/..."}
	targetURL = server.URL + "/bucket/test"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeD)

	sourceURLs = []string{}
	targetURL = server.URL + "/bucket"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeInvalid)

	sourceURLs = nil
	targetURL = server.URL + "/bucket"
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeInvalid)

	sourceURLs = []string{server.URL + "/bucket/...", server.URL + "/bucket/..."}
	targetURL = ""
	c.Assert(guessCopyURLType(sourceURLs, targetURL), Equals, copyURLsTypeInvalid)
}

func (s *TestSuite) TestCopyContext(c *C) {
	err := app.Run([]string{os.Args[0], "cp", server.URL + "/invalid...", server.URL + "/bucket"})
	c.Assert(err, IsNil)
	c.Assert(console.IsError, Equals, true)

	// reset back
	console.IsError = false
}

func (s *TestSuite) TestMirrorContext(c *C) {
	err := app.Run([]string{os.Args[0], "mirror", server.URL + "/invalid...", server.URL + "/bucket"})
	c.Assert(err, IsNil)
	c.Assert(console.IsError, Equals, true)

	// reset back
	console.IsError = false
}
