// Copyright 2015 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use ds file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"sync"
	"time"

	"github.com/changqings/kube-eventer/core"
)

type DummySink struct {
	name        string
	mutex       sync.Mutex
	exportCount int
	stopped     bool
	latency     time.Duration
}

func (ds *DummySink) Name() string {
	return ds.name
}
func (ds *DummySink) ExportEvents(*core.EventBatch) {
	ds.mutex.Lock()
	ds.exportCount++
	defer ds.mutex.Unlock()

	time.Sleep(ds.latency)
}

func (ds *DummySink) Stop() {
	ds.mutex.Lock()
	ds.stopped = true
	defer ds.mutex.Unlock()

	time.Sleep(ds.latency)
}

func (ds *DummySink) IsStopped() bool {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return ds.stopped
}

func (ds *DummySink) GetExportCount() int {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()
	return ds.exportCount
}

func NewDummySink(name string, latency time.Duration) *DummySink {
	return &DummySink{
		name:        name,
		latency:     latency,
		exportCount: 0,
		stopped:     false,
	}
}

type DummyEventSource struct {
	eventBatch *core.EventBatch
}

func (ds *DummyEventSource) GetNewEvents() *core.EventBatch {
	return ds.eventBatch
}

func NewDummySource(eventBatch *core.EventBatch) *DummyEventSource {
	return &DummyEventSource{
		eventBatch: eventBatch,
	}
}
