/*
Copyright 2025 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package source

import (
	"context"
	"text/template"

	nomad "github.com/hashicorp/nomad/api"
	log "github.com/sirupsen/logrus"

	"sigs.k8s.io/external-dns/endpoint"
)

// nomadSource is an implementation of Source for Nomad services.
type nomadSource struct {
	client           *nomad.Client
	namespace        string
	annotationFilter string

	fqdnTemplate             *template.Template
	combineFQDNAnnotation    bool
	ignoreHostnameAnnotation bool
}

// NewNomadSource creates a new nomadSource.
func NewNomadSource(ctx context.Context, nomadClient *nomad.Client, namespace, annotationFilter, fqdnTemplate string, combineFqdnAnnotation bool, ignoreHostnameAnnotation bool) (Source, error) {
	tmpl, err := parseTemplate(fqdnTemplate)
	if err != nil {
		return nil, err
	}

	return &nomadSource{
		client:                   nomadClient,
		namespace:                namespace,
		annotationFilter:         annotationFilter,
		fqdnTemplate:             tmpl,
		combineFQDNAnnotation:    combineFqdnAnnotation,
		ignoreHostnameAnnotation: ignoreHostnameAnnotation,
	}, nil
}

// Endpoints collects endpoints of all nested Sources and returns them in a single slice.
func (e *nomadSource) Endpoints(ctx context.Context) ([]*endpoint.Endpoint, error) {
	log.Info("Fetching endpoints from Nomad")

	return []*endpoint.Endpoint{}, nil
}

func (e *nomadSource) AddEventHandler(ctx context.Context, handler func()) {
	log.Info("Adding event handler for Nomad")
}
