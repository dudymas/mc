/*
 * Minio Client (C) 2014, 2015 Minio, Inc.
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
	"fmt"
	"strings"
	"time"

	"github.com/minio/cli"
	"github.com/minio/minio/pkg/probe"
)

// Share documents via URL.
var shareUpload = cli.Command{
	Name:   "upload",
	Usage:  "Share link that can be used to upload files to private bucket",
	Action: mainShareUpload,
	CustomHelpTemplate: `NAME:
   mc share {{.Name}} - {{.Usage}}

USAGE:
   mc share {{.Name}} TARGET [DURATION] [Content-Type]

   DURATION = NN[h|m|s] [DEFAULT=168h]

EXAMPLES:
   1. Generate Curl upload command, with a default expiry of 7 days.
      $ mc share {{.Name}} https://s3.amazonaws.com/backup/2006-Mar-1/backup.tar.gz

   2. Generate Curl upload command to upload files to a folder, with expiry of 120 hours
      $ mc share {{.Name}} https://s3.amazonaws.com/backup/2007-Mar-2/... 120h

   3. Generate Curl upload command to upload with expiry of 2 hours with content-type image/png
      $ mc share {{.Name}} https://s3.amazonaws.com/backup/2007-Mar-2/... 2h image/png

`,
}

func checkShareUploadSyntax(ctx *cli.Context) {
	args := ctx.Args()
	if !args.Present() || args.First() == "help" {
		cli.ShowCommandHelpAndExit(ctx, "upload", 1) // last argument is exit code
	}
	if len(args) > 3 {
		cli.ShowCommandHelpAndExit(ctx, "upload", 1) // last argument is exit code
	}
	url := stripRecursiveURL(strings.TrimSpace(args.Get(0)))
	if !isObjectKeyPresent(url) {
		fatalIf(errDummy().Trace(), fmt.Sprintf("Upload location needs object key ‘%s’.", strings.TrimSpace(args.Get(0))))
	}
	if strings.HasSuffix(strings.TrimSpace(args.Get(0)), "/") {
		fatalIf(errDummy().Trace(), fmt.Sprintf("Upload location cannot end with ‘/’. Did you mean ‘%s’.", url+recursiveSeparator))
	}
}

func mainShareUpload(ctx *cli.Context) {
	var expires time.Duration
	var err error
	shareDataDirSetup()
	checkShareUploadSyntax(ctx)
	setSharePalette(ctx.GlobalString("colors"))

	args := ctx.Args()
	config := mustGetMcConfig()
	if strings.TrimSpace(args.Get(1)) == "" {
		expires = time.Duration(604800) * time.Second
	} else {
		expires, err = time.ParseDuration(strings.TrimSpace(args.Get(1)))
		if err != nil {
			fatalIf(probe.NewError(err), "Unable to parse time argument.")
		}
	}
	contentType := strings.TrimSpace(args.Get(2))
	targetURL := getAliasURL(strings.TrimSpace(args.Get(0)), config.Aliases)

	e := doShareUploadURL(stripRecursiveURL(targetURL), isURLRecursive(targetURL), expires, contentType)
	fatalIf(e.Trace(targetURL), "Unable to generate URL for upload.")
}

// doShareURL share files from target
func doShareUploadURL(targetURL string, recursive bool, expires time.Duration, contentType string) *probe.Error {
	shareDate := time.Now().UTC()
	sURLs, err := loadSharedURLsV3()
	if err != nil {
		return err.Trace()
	}

	clnt, err := url2Client(targetURL)
	if err != nil {
		return err.Trace()
	}
	m, err := clnt.ShareUpload(recursive, expires, contentType)
	if err != nil {
		return err.Trace()
	}
	Key := targetURL
	if recursive {
		Key = Key + recursiveSeparator
		m["key"] = m["key"] + "<FILE>"
	}
	shareMessage := ShareMessage{
		Expiry:     expires,
		UploadInfo: m,
		Key:        Key,
	}
	shareMessageV3 := ShareMessageV3{
		Expiry:     expires,
		UploadInfo: m,
		Key:        Key,
	}
	sURLs.URLs = append(sURLs.URLs, struct {
		Date    time.Time
		Message ShareMessageV3
	}{
		Date:    shareDate,
		Message: shareMessageV3,
	})
	Prints("%s\n", shareMessage)
	saveSharedURLsV3(sURLs)
	return nil
}
