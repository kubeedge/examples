/*
Copyright 2024 The KubeEdge Authors.

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

// Package inference runs edge anomaly detection offline.
package inference

import (
	"math"
	"sync"

	"k8s.io/klog/v2"

	"github.com/kubeedge/examples/predictive-maintenance/mapper/pkg/driver"
)

const (
	DefaultWindowSize      = 20
	DefaultZScoreThreshold = 2.5
)

// AnomalyResult holds detection output.
type AnomalyResult struct {
	IsAnomaly         bool    `json:"isAnomaly"`
	VibrationZScore   float64 `json:"vibrationZScore"`
	TemperatureZScore float64 `json:"temperatureZScore"`
	Confidence        float64 `json:"confidence"`
}

// Detector is a thread-safe anomaly detector.
type Detector struct {
	mu              sync.RWMutex
	window          []driver.SensorReading
	windowSize      int
	zScoreThreshold float64
}

// NewDetector returns a default detector.
func NewDetector() *Detector {
	return &Detector{
		window:          make([]driver.SensorReading, 0, DefaultWindowSize),
		windowSize:      DefaultWindowSize,
		zScoreThreshold: DefaultZScoreThreshold,
	}
}

// NewDetectorWithParams returns a custom detector.
func NewDetectorWithParams(windowSize int, threshold float64) *Detector {
	return &Detector{
		window:          make([]driver.SensorReading, 0, windowSize),
		windowSize:      windowSize,
		zScoreThreshold: threshold,
	}
}

// Analyze classifies a reading as anomalous or not.
func (d *Detector) Analyze(reading driver.SensorReading) AnomalyResult {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.window) < 5 {
		d.window = append(d.window, reading)
		return AnomalyResult{}
	}

	vibMean, vibStd, tempMean, tempStd := calcStats(d.window)

	vibZ := zScore(reading.Vibration, vibMean, vibStd)
	tempZ := zScore(reading.Temperature, tempMean, tempStd)
	maxZ := math.Max(math.Abs(vibZ), math.Abs(tempZ))

	// Update window after analysis to avoid masking effect
	if len(d.window) >= d.windowSize {
		copy(d.window, d.window[1:])
		d.window[d.windowSize-1] = reading
	} else {
		d.window = append(d.window, reading)
	}

	isAnomaly := maxZ > d.zScoreThreshold
	confidence := math.Min(1.0, maxZ/d.zScoreThreshold/2)

	if isAnomaly {
		klog.Warningf("ANOMALY DETECTED: vibration_z=%.2f, temperature_z=%.2f, confidence=%.2f",
			vibZ, tempZ, confidence)
	}

	return AnomalyResult{
		IsAnomaly:         isAnomaly,
		VibrationZScore:   math.Round(vibZ*100) / 100,
		TemperatureZScore: math.Round(tempZ*100) / 100,
		Confidence:        math.Round(confidence*100) / 100,
	}
}

// WindowStats returns current window statistics.
func (d *Detector) WindowStats() (vibMean, vibStd, tempMean, tempStd float64, n int) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.window) == 0 {
		return
	}
	n = len(d.window)
	vibMean, vibStd, tempMean, tempStd = calcStats(d.window)
	return
}

func calcStats(window []driver.SensorReading) (vibMean, vibStd, tempMean, tempStd float64) {
	var m2Vib, m2Temp float64
	for i, r := range window {
		n := float64(i + 1)
		dv := r.Vibration - vibMean
		dt := r.Temperature - tempMean
		vibMean += dv / n
		tempMean += dt / n
		m2Vib += dv * (r.Vibration - vibMean)
		m2Temp += dt * (r.Temperature - tempMean)
	}
	n := float64(len(window))
	if n > 0 {
		vibStd = math.Sqrt(m2Vib / n)
		tempStd = math.Sqrt(m2Temp / n)
	}
	return
}

func zScore(value, mean, std float64) float64 {
	// Add noise floor to prevent division by zero and allow spike detection in stable data
	std = math.Max(std, 1e-4)
	return (value - mean) / std
}
