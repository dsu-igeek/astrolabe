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
	"archive/zip"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
)

func ZipProtectedEntity(ctx context.Context, pe ProtectedEntity, writer io.Writer) error {
	zipWriter := zip.NewWriter(writer)
	peInfo, err := pe.GetInfo(ctx)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to get info for %s", pe.GetID().String())
	}
	jsonBuf, err := json.Marshal(peInfo)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to get marshal info for %s", pe.GetID().String())
	}
	peInfoWriter, err := zipWriter.Create(pe.GetID().String() + ".peinfo")
	if (err != nil) {
		return errors.Wrapf(err, "Failed to create peinfo zip writer for %s", pe.GetID().String())
	}
	var bytesWritten int64
	jsonWritten, err := peInfoWriter.Write(jsonBuf)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to write info json for %s", pe.GetID().String())
	}
	bytesWritten += int64(jsonWritten)


	mdReader, err := pe.GetMetadataReader(ctx)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to get metadata reader for %s", pe.GetID().String())
	}
	if mdReader != nil {
		peMDWriter, err := zipWriter.Create(pe.GetID().String() + ".md")
		if (err != nil) {
			return errors.Wrapf(err, "Failed to create metadata zip writer for %s", pe.GetID().String())
		}
		mdWritten, err := io.Copy(peMDWriter, mdReader)
		if (err != nil) {
			return errors.Wrapf(err, "Failed to write metadata for %s", pe.GetID().String())
		}
		bytesWritten += mdWritten
	}

	dataReader, err := pe.GetDataReader(ctx)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to get data reader for %s", pe.GetID().String())
	}
	if dataReader != nil {
		peDataWriter, err := zipWriter.Create(pe.GetID().String() + ".data")
		if (err != nil) {
			return errors.Wrapf(err, "Failed to create data zip writer for %s", pe.GetID().String())
		}
		dataWritten, err := io.Copy(peDataWriter, dataReader)
		if (err != nil) {
			return errors.Wrapf(err, "Failed to write data for %s", pe.GetID().String())
		}
		bytesWritten += dataWritten
	}
	components, err := pe.GetComponents(ctx)
	if (err != nil) {
		return errors.Wrapf(err, "Failed to get components for %s", pe.GetID().String())
	}

	for _,curComponent := range components {
		componentWriter, err := zipWriter.Create("components/" + curComponent.GetID().String() + ".zip")
		if (err != nil) {
			return errors.Wrapf(err, "Failed to get create component writer for component %s of %s",
				curComponent.GetID().String(), pe.GetID().String())
		}
		err = ZipProtectedEntity(ctx, curComponent, componentWriter)
		if (err != nil) {
			return errors.Wrapf(err, "Failed to write zip for component %s of %s",
				curComponent.GetID().String(), pe.GetID().String())
		}
	}
	err = zipWriter.Flush()
	if (err != nil) {
		return errors.Wrapf(err, "Zipwriter flush failed for for %s", pe.GetID().String())
	}
	err = zipWriter.Close()
	if (err != nil) {
		return errors.Wrapf(err, "Zipwriter close failed for for %s", pe.GetID().String())
	}
	return nil
}
