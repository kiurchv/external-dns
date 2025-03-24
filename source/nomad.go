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
	"sort"
	"strings"
	"text/template"

	nomad "github.com/hashicorp/nomad/api"
	log "github.com/sirupsen/logrus"

	"sigs.k8s.io/external-dns/endpoint"
)

const (
	tagPrefix = "external-dns"
)

// nomadSource is an implementation of Source for Nomad services.
type nomadSource struct {
	client    *nomad.Client
	namespace string

	fqdnTemplate             *template.Template
	combineFQDNAnnotation    bool
	ignoreHostnameAnnotation bool
}

// NewNomadSource creates a new nomadSource.
func NewNomadSource(ctx context.Context, nomadClient *nomad.Client, namespace, fqdnTemplate string, combineFqdnAnnotation bool, ignoreHostnameAnnotation bool) (Source, error) {
	tmpl, err := parseTemplate(fqdnTemplate)
	if err != nil {
		return nil, err
	}

	return &nomadSource{
		client:                   nomadClient,
		namespace:                namespace,
		fqdnTemplate:             tmpl,
		combineFQDNAnnotation:    combineFqdnAnnotation,
		ignoreHostnameAnnotation: ignoreHostnameAnnotation,
	}, nil
}

// Endpoints collects endpoints of all nested Sources and returns them in a single slice.
func (ns *nomadSource) Endpoints(ctx context.Context) ([]*endpoint.Endpoint, error) {
	log.Info("Fetching endpoints from Nomad", "namespace:", ns.namespace)

	namespace := ns.namespace
	if namespace == "" {
		namespace = "*"
	}
	opts := &nomad.QueryOptions{
		Namespace: namespace,
	}
	opts = opts.WithContext(ctx)

	serviceLists, _, err := ns.client.Services().List(opts)
	if err != nil {
		return nil, err
	}

	endpoints := []*endpoint.Endpoint{}

	for _, serviceList := range serviceLists {
		for _, service := range serviceList.Services {
			annotations := tagsToAnnotations(service.Tags)

			controller, ok := annotations[controllerAnnotationKey]
			if ok && controller != controllerAnnotationValue {
				log.Debugf("Skipping service %s/%s because controller value does not match, found: %s, required: %s",
					serviceList.Namespace, service.ServiceName, controller, controllerAnnotationValue)
				continue
			}

			// TBD: Implement logic to create endpoints from Nomad service information. See source/services.go for (rough) reference.
		}
	}

	for _, ep := range endpoints {
		sort.Sort(ep.Targets)
	}

	return endpoints, nil
}

func tagsToAnnotations(tags []string) map[string]string {
	annotations := make(map[string]string, len(tags))
	for _, tag := range tags {
		if strings.HasPrefix(tag, tagPrefix) {
			if parts := strings.SplitN(tag, "=", 2); len(parts) == 2 {
				left, right := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
				key := "external-dns.alpha.kubernetes.io/" + strings.TrimPrefix(left, tagPrefix+".")
				annotations[key] = right
			}
		}
	}
	return annotations
}

func (ns *nomadSource) AddEventHandler(ctx context.Context, handler func()) {}
