// Copyright 2022 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package drone converts Drone pipelines to Harness pipelines.
package drone

import (
	"bytes"
	"errors"
	"io"
	"os"

	harness "github.com/drone/spec/dist/go"

	"github.com/drone/go-convert/internal/store"
	"github.com/ghodss/yaml"
)

// Converter converts a Drone pipeline to a Harness
// v1 pipeline.
type Converter struct {
	kubeEnabled   bool
	kubeNamespace string
	kubeConnector string
	dockerhubConn string
	identifiers   *store.Identifiers

	// // as we walk the yaml, we store a
	// // a snapshot of the current node and
	// // its parents.
	// config *drone.Pipeline
	// stage  *drone.Stage
}

// New creates a new Converter that converts a Drone
// pipeline to a Harness v1 pipeline.
func New(options ...Option) *Converter {
	d := new(Converter)

	// create the unique identifier store. this store
	// is used for registering unique identifiers to
	// prevent duplicate names, unique index violations.
	d.identifiers = store.New()

	// loop through and apply the options.
	for _, option := range options {
		option(d)
	}

	// set the default kubernetes namespace.
	if d.kubeNamespace == "" {
		d.kubeNamespace = "default"
	}

	// set the runtime to kubernetes if the kubernetes
	// connector is configured.
	if d.kubeConnector != "" {
		d.kubeEnabled = true
	}

	return d
}

// Convert downgrades a v1 pipeline.
func (d *Converter) Convert(r io.Reader) ([]byte, error) {
	// src, err := drone.Parse(r)
	// if err != nil {
	// 	return nil, err
	// }
	// d.config = src // push the drone config to the state
	// return d.convert()
	return nil, errors.New("not implemented")
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertBytes(b []byte) ([]byte, error) {
	return d.Convert(
		bytes.NewBuffer(b),
	)
}

// ConvertString downgrades a v1 pipeline.
func (d *Converter) ConvertString(s string) ([]byte, error) {
	return d.Convert(
		bytes.NewBufferString(s),
	)
}

// ConvertFile downgrades a v1 pipeline.
func (d *Converter) ConvertFile(p string) ([]byte, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return d.Convert(f)
}

// converts converts a Drone pipeline to a Harness pipeline.
func (d *Converter) convert() ([]byte, error) {

	// create the harness pipeline
	pipeline := &harness.Pipeline{
		Version: 1,
		// Default: convertDefault(d.config),
	}

	// for _, steps := range d.config.Pipelines.Default {
	// 	if steps.Stage != nil {
	// 		// TODO support for fast-fail
	// 		d.stage = steps.Stage // push the stage to the state
	// 		stage := d.convertStage()
	// 		pipeline.Stages = append(pipeline.Stages, stage)
	// 	}
	// }

	// marshal the harness yaml
	out, err := yaml.Marshal(pipeline)
	if err != nil {
		return nil, err
	}

	return out, nil
}
