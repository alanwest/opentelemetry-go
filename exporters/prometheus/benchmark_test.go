// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheus

import (
	"context"
	"fmt"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/sdk/metric"
)

func benchmarkCollect(b *testing.B, n int) {
	ctx := context.Background()
	exporter := New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("testmeter")

	registry := prometheus.NewRegistry()
	err := registry.Register(exporter.Collector)
	require.NoError(b, err)

	for i := 0; i < n; i++ {
		counter, err := meter.SyncFloat64().Counter(fmt.Sprintf("foo_%d", i))
		require.NoError(b, err)
		counter.Add(ctx, float64(i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := registry.Gather()
		require.NoError(b, err)
	}
}

func BenchmarkCollect1(b *testing.B)     { benchmarkCollect(b, 1) }
func BenchmarkCollect10(b *testing.B)    { benchmarkCollect(b, 10) }
func BenchmarkCollect100(b *testing.B)   { benchmarkCollect(b, 100) }
func BenchmarkCollect1000(b *testing.B)  { benchmarkCollect(b, 1000) }
func BenchmarkCollect10000(b *testing.B) { benchmarkCollect(b, 10000) }
