/*
 * Minio Go Library for Amazon S3 Legacy v2 Signature Compatible Cloud Storage (C) 2015 Minio, Inc.
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

package minio

import (
	"encoding/xml"
	"time"
)

// listAllMyBucketsResult container for listBuckets response
type listAllMyBucketsResult struct {
	// Container for one or more buckets.
	Buckets struct {
		Bucket []BucketStat
	}
	Owner owner
}

// owner container for bucket owner information
type owner struct {
	DisplayName string
	ID          string
}

// commonPrefix container for prefix response
type commonPrefix struct {
	Prefix string
}

// listBucketResult container for listObjects response
type listBucketResult struct {
	CommonPrefixes []commonPrefix // A response can contain CommonPrefixes only if you have specified a delimiter
	Contents       []ObjectStat   // Metadata about each object returned
	Delimiter      string

	// Encoding type used to encode object keys in the response.
	EncodingType string

	// A flag that indicates whether or not ListObjects returned all of the results
	// that satisfied the search criteria.
	IsTruncated bool
	Marker      string
	MaxKeys     int64
	Name        string

	// When response is truncated (the IsTruncated element value in the response
	// is true), you can use the key name in this field as marker in the subsequent
	// request to get next set of objects. Object storage lists objects in alphabetical
	// order Note: This element is returned only if you have delimiter request parameter
	// specified. If response does not include the NextMaker and it is truncated,
	// you can use the value of the last Key in the response as the marker in the
	// subsequent request to get the next set of object keys.
	NextMarker string
	Prefix     string
}

// multiPartUpload container for multipart session
type multiPartUpload struct {
	// Date and time at which the multipart upload was initiated.
	Initiated time.Time `type:"timestamp" timestampFormat:"iso8601"`

	Initiator initiator
	Owner     owner

	StorageClass string

	// Key of the object for which the multipart upload was initiated.
	Key string

	// Upload ID that identifies the multipart upload.
	UploadID string `xml:"UploadId"`
}

// listMultipartUploadsResult container for ListMultipartUploads response
type listMultipartUploadsResult struct {
	Bucket             string
	KeyMarker          string
	UploadIDMarker     string `xml:"UploadIdMarker"`
	NextKeyMarker      string
	NextUploadIDMarker string `xml:"NextUploadIdMarker"`
	EncodingType       string
	MaxUploads         int64
	IsTruncated        bool
	Uploads            []multiPartUpload `xml:"Upload"`
	Prefix             string
	Delimiter          string
	CommonPrefixes     []commonPrefix // A response can contain CommonPrefixes only if you specify a delimiter
}

// initiator container for who initiated multipart upload
type initiator struct {
	ID          string
	DisplayName string
}

// partMetadata container for particular part of an object
type partMetadata struct {
	// Part number identifies the part.
	PartNumber int

	// Date and time the part was uploaded.
	LastModified time.Time

	// Entity tag returned when the part was uploaded, usually md5sum of the part
	ETag string

	// Size of the uploaded part data.
	Size int64
}

// listObjectPartsResult container for ListObjectParts response.
type listObjectPartsResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`

	Initiator initiator
	Owner     owner

	StorageClass         string
	PartNumberMarker     int
	NextPartNumberMarker int
	MaxParts             int

	// Indicates whether the returned list of parts is truncated.
	IsTruncated bool
	Parts       []partMetadata `xml:"Part"`

	EncodingType string
}

// initiateMultipartUploadResult container for InitiateMultiPartUpload response.
type initiateMultipartUploadResult struct {
	Bucket   string
	Key      string
	UploadID string `xml:"UploadId"`
}

// completeMultipartUploadResult container for completed multipart upload response.
type completeMultipartUploadResult struct {
	Location string
	Bucket   string
	Key      string
	ETag     string
}

// completePart sub container lists individual part numbers and their md5sum, part of completeMultipartUpload.
type completePart struct {
	XMLName xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ Part" json:"-"`

	// Part number identifies the part.
	PartNumber int
	ETag       string
}

// completeMultipartUpload container for completing multipart upload
type completeMultipartUpload struct {
	XMLName xml.Name       `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CompleteMultipartUpload" json:"-"`
	Parts   []completePart `xml:"Part"`
}

// createBucketConfiguration container for bucket configuration
type createBucketConfiguration struct {
	XMLName  xml.Name `xml:"http://s3.amazonaws.com/doc/2006-03-01/ CreateBucketConfiguration" json:"-"`
	Location string   `xml:"LocationConstraint"`
}

type grant struct {
	Grantee struct {
		ID           string
		DisplayName  string
		EmailAddress string
		Type         string
		URI          string
	}
	Permission string
}

type accessControlPolicy struct {
	AccessControlList struct {
		Grant []grant
	}
	Owner owner
}
