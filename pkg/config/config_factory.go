/*
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */

package config

import (
	"flag"
	"log"
	"os"

	"go.opentelemetry.io/collector/confmap/converter/expandconverter"
	"go.opentelemetry.io/collector/confmap/converter/overwritepropertiesconverter"
	"go.opentelemetry.io/collector/confmap/provider/envprovider"
	"go.opentelemetry.io/collector/confmap/provider/fileprovider"
	"go.opentelemetry.io/collector/confmap/provider/yamlprovider"
	"go.opentelemetry.io/collector/service"

	"go.opentelemetry.io/collector/confmap"
)

const (
	envKey = "AOT_CONFIG_CONTENT"
)

func GetConfigProvider(flags *flag.FlagSet) service.ConfigProvider {
	// aws-otel-collector supports loading yaml config from Env Var
	// including SSM parameter store for ECS use case
	loc := getConfigFlag(flags)
	if configContent, ok := os.LookupEnv(envKey); ok {
		log.Printf("Reading AOT config from environment: %v\n", configContent)
		loc = []string{"env:" + envKey}
	}

	// generate the MapProviders for the Config Provider Settings
	providers := []confmap.Provider{fileprovider.New(), envprovider.New(), yamlprovider.New()}

	mapProviders := make(map[string]confmap.Provider, len(providers))
	for _, provider := range providers {
		mapProviders[provider.Scheme()] = provider
	}

	// create Config Provider Settings
	settings := service.ConfigProviderSettings{
		ResolverSettings: confmap.ResolverSettings{
			URIs:       loc,
			Providers:  mapProviders,
			Converters: []confmap.Converter{expandconverter.New(), overwritepropertiesconverter.New(getSetFlag(flags))},
		},
	}

	// get New config Provider
	config_provider, err := service.NewConfigProvider(settings)

	if err != nil {
		log.Panicf("Err on creating Config Provider: %v\n", err)
	}

	return config_provider
}
