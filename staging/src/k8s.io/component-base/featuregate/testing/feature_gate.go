/*
Copyright 2017 The Kubernetes Authors.

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

package testing

import (
	"fmt"
	"testing"

	"k8s.io/component-base/featuregate"
)

// SetFeatureGateDuringTest sets the specified gate to the specified value, and returns a function that restores the original value.
// Failures to set or restore cause the test to fail.
//
// Example use:
//
// defer featuregatetesting.SetFeatureGateDuringTest(t, featuregate.Default, features.<FeatureName>, true)()
func SetFeatureGateDuringTest(tb testing.TB, gate featuregate.FeatureGate, f featuregate.Feature, value bool) func() {
	originalValue := gate.Enabled(f)

	// Specially handle AllAlpha and AllBeta
	var cleanups []func()
	if f == "AllAlpha" || f == "AllBeta" {
		// Iterate over individual gates so their individual values get restored
		for k, v := range gate.(featuregate.MutableFeatureGate).GetAll() {
			if k == "AllAlpha" || k == "AllBeta" {
				continue
			}
			if (f == "AllAlpha" && v.PreRelease == featuregate.Alpha) || (f == "AllBeta" && v.PreRelease == featuregate.Beta) {
				cleanups = append(cleanups, SetFeatureGateDuringTest(tb, gate, k, value))
			}
		}
	}

	if err := gate.(featuregate.MutableFeatureGate).Set(fmt.Sprintf("%s=%v", f, value)); err != nil {
		tb.Errorf("error setting %s=%v: %v", f, value, err)
	}

	return func() {
		if err := gate.(featuregate.MutableFeatureGate).Set(fmt.Sprintf("%s=%v", f, originalValue)); err != nil {
			tb.Errorf("error restoring %s=%v: %v", f, originalValue, err)
		}
		for _, cleanup := range cleanups {
			cleanup()
		}
	}
}
