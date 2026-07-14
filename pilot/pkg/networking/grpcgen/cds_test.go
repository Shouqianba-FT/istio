// Copyright Istio Authors
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

package grpcgen

import (
	"testing"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	networking "istio.io/api/networking/v1alpha3"
)

func TestSetSlowStartConfigWarmup(t *testing.T) {
	warmupDuration := durationpb.New(30 * time.Second)
	cfg := setSlowStartConfig(&networking.LoadBalancerSettings{
		Warmup: &networking.WarmupConfiguration{
			Duration:       warmupDuration,
			Aggression:     wrapperspb.Double(0.5),
			MinimumPercent: wrapperspb.Double(5),
		},
	})

	if cfg.SlowStartWindow != warmupDuration {
		t.Fatalf("unexpected slow start window: got %v want %v", cfg.SlowStartWindow, warmupDuration)
	}
	if cfg.Aggression.GetDefaultValue() != 0.5 {
		t.Fatalf("unexpected aggression: got %v want 0.5", cfg.Aggression.GetDefaultValue())
	}
	if cfg.Aggression.GetRuntimeKey() != "istio.slowstart.aggression" {
		t.Fatalf("unexpected aggression runtime key: %q", cfg.Aggression.GetRuntimeKey())
	}
	if cfg.MinWeightPercent.GetValue() != 5 {
		t.Fatalf("unexpected minimum weight percent: got %v want 5", cfg.MinWeightPercent.GetValue())
	}
}

func TestSetSlowStartConfigWarmupDefaults(t *testing.T) {
	warmupDuration := durationpb.New(30 * time.Second)
	cfg := setSlowStartConfig(&networking.LoadBalancerSettings{
		Warmup: &networking.WarmupConfiguration{
			Duration: warmupDuration,
		},
	})

	if cfg.SlowStartWindow != warmupDuration {
		t.Fatalf("unexpected slow start window: got %v want %v", cfg.SlowStartWindow, warmupDuration)
	}
	if cfg.Aggression.GetDefaultValue() != 1 {
		t.Fatalf("unexpected default aggression: got %v want 1", cfg.Aggression.GetDefaultValue())
	}
	if cfg.MinWeightPercent.GetValue() != 10 {
		t.Fatalf("unexpected default minimum weight percent: got %v want 10", cfg.MinWeightPercent.GetValue())
	}
}

func TestSetSlowStartConfigWarmupDurationSecs(t *testing.T) {
	warmupDuration := durationpb.New(30 * time.Second)
	cfg := setSlowStartConfig(&networking.LoadBalancerSettings{
		WarmupDurationSecs: warmupDuration,
	})

	if cfg.SlowStartWindow != warmupDuration {
		t.Fatalf("unexpected slow start window: got %v want %v", cfg.SlowStartWindow, warmupDuration)
	}
	if cfg.Aggression != nil {
		t.Fatalf("unexpected aggression for legacy warmupDurationSecs: %v", cfg.Aggression)
	}
	if cfg.MinWeightPercent != nil {
		t.Fatalf("unexpected minimum weight percent for legacy warmupDurationSecs: %v", cfg.MinWeightPercent)
	}
}
