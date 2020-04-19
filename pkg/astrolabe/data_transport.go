/*
 * Copyright 2019 the Astrolabe contributors
 * SPDX-License-Identifier: Apache-2.0
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

package astrolabe

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	credentials2 "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vmware-tanzu/astrolabe/gen/models"
	"net"
	"strconv"
	"time"
)

// DataTransport is our internal interface representing the data transport for Protected Entity
// data, metadata or combined info
// DataTransport contains parameters for the transport but does not actually move data
// DataTransport is used in two ways:
//		Each ProtectedEntity exports a set of DataTransports for accessing its data, metadata, and combined streams
//		These exported DataTransports are used to form the JSON
//
//		When we copy from a ProtectedEntity, the DataTransports of the source may be used by the ProtectedEntity to
//		return a stream.  This is most useful for remote ProtectedEntities.

type DataTransport struct {
	// The type of this data source, e.g. S3, VADP
	transportType string
	params        map[string]string
}

func NewDataTransport(transportType string, params map[string]string) DataTransport {
	return DataTransport{
		transportType: transportType,
		params:        params,
	}
}

const(
	S3TransportType = "s3"
	S3URLParam = "url"
	S3HostParam = "host"
	S3BucketParam = "bucket"
	S3KeyParam = "key"
)

type S3Config struct {
	Port int			`json:"port,omitempty"`
	Host net.IP			`json:"host,omitempty"`
	AccessKey string	`json:"accessKey,omitempty"`
	Secret string		`json:"secret,omitempty"`
	Prefix string		`json:"prefix,omitempty"`
	URLBase string		`json:"urlBase,omitempty"`
}

func getLocalIP() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("Could not determine local IP address")
}

func (this S3Config) getURL() (string) {
	return "http://" + this.Host.String() + ":" + strconv.Itoa(this.Port) + "/"+this.Prefix+"/"
}

func newDataTransportForS3URL(url string) DataTransport {
	return DataTransport{
		transportType: S3TransportType,
		params: map[string]string{
			S3URLParam: url,
		},
	}
}

func newDataTransportForS3(host string, bucket string, key string) DataTransport {
	url := "http://" + host + "/" + bucket + "/" + key
	return DataTransport{
		transportType: S3TransportType,
		params: map[string]string{
			S3URLParam:    url,
			S3HostParam:   host,
			S3BucketParam: bucket,
			S3KeyParam:    key,
		},
	}
}

const (
	DataExt = ""
	MDExt = ".md"
	CombinedExt = ".zip"
	PEInfoExt = ".peinfo"
)

func NewS3DataTransportForPEID(peid ProtectedEntityID, s3Config S3Config) DataTransport {
	return NewS3TransportForPEID(peid, DataExt, s3Config)
}

func NewS3MDTransportForPEID(peid ProtectedEntityID, s3Config S3Config) DataTransport {
	return NewS3TransportForPEID(peid, MDExt, s3Config)
}

func NewS3CombinedTransportForPEID(peid ProtectedEntityID, s3Config S3Config) DataTransport {
	return NewS3TransportForPEID(peid, CombinedExt, s3Config)
}

func NewS3PEInfoTransportForPEID(peid ProtectedEntityID, s3Config S3Config) DataTransport {
	return NewS3TransportForPEID(peid, PEInfoExt, s3Config)
}


func NewS3TransportForPEID(peid ProtectedEntityID, ext string, s3Config S3Config) DataTransport {
	credentials := credentials2.NewStaticCredentials(s3Config.AccessKey, s3Config.Secret, "")
	s3ForcePathStyle := true
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(s3Config.getURL()),
		Credentials: credentials,
		Region: aws.String("us-west-2"),
		S3ForcePathStyle: &s3ForcePathStyle,
	},
	)

	if err != nil {

	}
	// Create S3 service client
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(peid.peType),
		Key:    aws.String(peid.String() + ext),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	return newDataTransportForS3URL(urlStr)
}
func (this DataTransport) GetTransportType() string {
	return this.transportType
}

func (this DataTransport) GetParam(key string) (string, bool) {
	val, ok := this.params[key]
	return val, ok
}

func (this DataTransport) getModelDataTransport() models.DataTransport {
	return models.DataTransport{
		TransportType: this.transportType,
		Params: this.params,
	}
}

func newDataTransportForModelTransport(transport models.DataTransport) DataTransport {
	return DataTransport{
		transportType: transport.TransportType,
		params:        transport.Params,
	}
}

func (this DataTransport) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.getModelDataTransport())
}

func (this *DataTransport) UnmarshalJSON(data []byte) error {
	jsonStruct := models.DataTransport{}
	err := json.Unmarshal(data, &jsonStruct)
	if err != nil {
		return err
	}
	this.transportType = jsonStruct.TransportType
	this.params = jsonStruct.Params
	return nil
}
