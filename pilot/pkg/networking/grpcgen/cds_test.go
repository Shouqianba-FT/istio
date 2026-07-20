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

	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	networking "istio.io/api/networking/v1alpha3"
)

func TestApplyRoundRobinLoadBalancerWarmup(t *testing.T) {
	c := &cluster.Cluster{}
	loadbalancer := &networking.LoadBalancerSettings{
		Warmup: &networking.WarmupConfiguration{
			Duration: durationpb.New(5 * time.Second),
		},
	}

	applyRoundRobinLoadBalancer(c, loadbalancer)

	lbConfig, ok := c.LbConfig.(*cluster.Cluster_RoundRobinLbConfig_)
	if !ok {
		t.Fatalf("expected round robin lb config, got %T", c.LbConfig)
	}
	slowStart := lbConfig.RoundRobinLbConfig.GetSlowStartConfig()
	if slowStart.GetSlowStartWindow().GetSeconds() != 5 {
		t.Fatalf("expected slow start window 5s, got %v", slowStart.GetSlowStartWindow())
	}
	if slowStart.GetAggression().GetDefaultValue() != 1 {
		t.Fatalf("expected default aggression 1, got %v", slowStart.GetAggression().GetDefaultValue())
	}
	if slowStart.GetAggression().GetRuntimeKey() != "istio.slowstart.aggression" {
		t.Fatalf("unexpected aggression runtime key %q", slowStart.GetAggression().GetRuntimeKey())
	}
	if slowStart.GetMinWeightPercent().GetValue() != 10 {
		t.Fatalf("expected default min weight percent 10, got %v", slowStart.GetMinWeightPercent().GetValue())
	}
}

func TestApplyRoundRobinLoadBalancerWarmupExplicitValues(t *testing.T) {
	c := &cluster.Cluster{}
	loadbalancer := &networking.LoadBalancerSettings{
		Warmup: &networking.WarmupConfiguration{
			Duration:       durationpb.New(10 * time.Second),
			Aggression:     wrapperspb.Double(0.5),
			MinimumPercent: wrapperspb.Double(2.5),
		},
	}

	applyRoundRobinLoadBalancer(c, loadbalancer)

	lbConfig, ok := c.LbConfig.(*cluster.Cluster_RoundRobinLbConfig_)
	if !ok {
		t.Fatalf("expected round robin lb config, got %T", c.LbConfig)
	}
	slowStart := lbConfig.RoundRobinLbConfig.GetSlowStartConfig()
	if slowStart.GetSlowStartWindow().GetSeconds() != 10 {
		t.Fatalf("expected slow start window 10s, got %v", slowStart.GetSlowStartWindow())
	}
	if slowStart.GetAggression().GetDefaultValue() != 0.5 {
		t.Fatalf("expected aggression 0.5, got %v", slowStart.GetAggression().GetDefaultValue())
	}
	if slowStart.GetMinWeightPercent().GetValue() != 2.5 {
		t.Fatalf("expected min weight percent 2.5, got %v", slowStart.GetMinWeightPercent().GetValue())
	}
}
