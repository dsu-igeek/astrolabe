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

package server

import (
	"context"
	"github.com/labstack/echo"
	"github.com/vmware-tanzu/astrolabe/pkg/astrolabe"
	"net/http"
	"strings"
)

type Astrolabe struct {
	petm         *DirectProtectedEntityManager
	api_services map[string]*ServiceAPI
	s3_services  map[string]*ServiceS3
}

func NewAstrolabe(confDirPath string) *Astrolabe {

	pem := NewProtectedEntityManager(confDirPath)
	api_services := make(map[string]*ServiceAPI)
	s3_services := make(map[string]*ServiceS3)
	for _, curService := range pem.ListEntityTypeManagers() {
		serviceName := curService.GetTypeName()
		api_services[serviceName] = NewServiceAPI(curService)
		s3_services[serviceName] = NewServiceS3(curService)
	}

	retAstrolabe := Astrolabe{
		api_services: api_services,
		s3_services:  s3_services,
	}

	return &retAstrolabe
}

func NewProtectedEntityManager(confDirPath string) (astrolabe.ProtectedEntityManager) {
	dpem := NewDirectProtectedEntityManagerFromConfigDir(confDirPath)
	var pem astrolabe.ProtectedEntityManager
	pem = dpem
	return pem
}

func NewAstrolabeRepository() *Astrolabe {
	return nil
}

func (this *Astrolabe) Get(c echo.Context) error {
	var servicesList strings.Builder
	needsComma := false
	for serviceName := range this.api_services {
		if needsComma {
			servicesList.WriteString(",")
		}
		servicesList.WriteString(serviceName)
		needsComma = true
	}
	return c.String(http.StatusOK, servicesList.String())
}


const fileSuffix = ".pe.json"

func getProtectedEntityForIDStr(petm astrolabe.ProtectedEntityTypeManager, idStr string,
	echoContext echo.Context) (astrolabe.ProtectedEntityID, astrolabe.ProtectedEntity, error) {
	var id astrolabe.ProtectedEntityID
	var pe astrolabe.ProtectedEntity
	var err error

	id, err = astrolabe.NewProtectedEntityIDFromString(idStr)
	if err != nil {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is invalid "+err.Error())
		return id, pe, err
	}
	if id.GetPeType() != (petm).GetTypeName() {
		echoContext.String(http.StatusBadRequest, "id = "+idStr+" is not type "+petm.GetTypeName())
		return id, pe, err
	}
	pe, err = (petm).GetProtectedEntity(context.Background(), id)
	if err != nil {
		echoContext.String(http.StatusNotFound, "Could not retrieve id "+id.String()+" error = "+err.Error())
		return id, pe, err
	}
	if pe == nil {
		echoContext.String(http.StatusInternalServerError, "pe was nil for "+id.String())
		return id, pe, err
	}
	return id, pe, nil
}
