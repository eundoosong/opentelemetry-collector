// Copyright  OpenTelemetry Authors
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

package exporterhelper

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configerror"
	"go.opentelemetry.io/collector/config/configmodels"
)

// FactoryOption apply changes to ExporterOptions.
type FactoryOption func(o *factory)

// CreateDefaultConfig is the equivalent of component.ExporterFactory.CreateDefaultConfig()
type CreateDefaultConfig func() configmodels.Exporter

// CreateTraceExporter is the equivalent of component.ExporterFactory.CreateTraceExporter()
type CreateTraceExporter func(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.TraceExporter, error)

// CreateMetricsExporter is the equivalent of component.ExporterFactory.CreateMetricsExporter()
type CreateMetricsExporter func(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.MetricsExporter, error)

// CreateMetricsExporter is the equivalent of component.ExporterFactory.CreateLogsExporter()
type CreateLogsExporter func(context.Context, component.ExporterCreateParams, configmodels.Exporter) (component.LogsExporter, error)

// factory is the factory for Jaeger gRPC exporter.
type factory struct {
	cfgType               configmodels.Type
	createDefaultConfig   CreateDefaultConfig
	createTraceExporter   CreateTraceExporter
	createMetricsExporter CreateMetricsExporter
	createLogsExporter    CreateLogsExporter
}

var _ component.LogsExporterFactory = new(factory)

// WithTraces overrides the default "error not supported" implementation for CreateTraceReceiver.
func WithTraces(createTraceExporter CreateTraceExporter) FactoryOption {
	return func(o *factory) {
		o.createTraceExporter = createTraceExporter
	}
}

// WithMetrics overrides the default "error not supported" implementation for CreateMetricsReceiver.
func WithMetrics(createMetricsExporter CreateMetricsExporter) FactoryOption {
	return func(o *factory) {
		o.createMetricsExporter = createMetricsExporter
	}
}

// WithLogs overrides the default "error not supported" implementation for CreateLogsReceiver.
func WithLogs(createLogsExporter CreateLogsExporter) FactoryOption {
	return func(o *factory) {
		o.createLogsExporter = createLogsExporter
	}
}

// NewFactory returns a component.ExporterFactory that only supports all types.
func NewFactory(
	cfgType configmodels.Type,
	createDefaultConfig CreateDefaultConfig,
	options ...FactoryOption) component.ExporterFactory {
	f := &factory{
		cfgType:             cfgType,
		createDefaultConfig: createDefaultConfig,
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

// Type gets the type of the Exporter config created by this factory.
func (f *factory) Type() configmodels.Type {
	return f.cfgType
}

// CreateDefaultConfig creates the default configuration for processor.
func (f *factory) CreateDefaultConfig() configmodels.Exporter {
	return f.createDefaultConfig()
}

// CreateTraceExporter creates a component.TraceExporter based on this config.
func (f *factory) CreateTraceExporter(
	ctx context.Context,
	params component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.TraceExporter, error) {
	if f.createTraceExporter != nil {
		return f.createTraceExporter(ctx, params, cfg)
	}
	return nil, configerror.ErrDataTypeIsNotSupported
}

// CreateMetricsExporter creates a consumer.MetricsConsumer based on this config.
func (f *factory) CreateMetricsExporter(
	ctx context.Context,
	params component.ExporterCreateParams,
	cfg configmodels.Exporter) (component.MetricsExporter, error) {
	if f.createMetricsExporter != nil {
		return f.createMetricsExporter(ctx, params, cfg)
	}
	return nil, configerror.ErrDataTypeIsNotSupported
}

// CreateLogsExporter creates a metrics processor based on this config.
func (f *factory) CreateLogsExporter(
	ctx context.Context,
	params component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.LogsExporter, error) {
	if f.createLogsExporter != nil {
		return f.createLogsExporter(ctx, params, cfg)
	}
	return nil, configerror.ErrDataTypeIsNotSupported
}
