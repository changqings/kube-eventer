// Copyright 2015 Google Inc. All Rights Reserved.
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

package logsink

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	kube_api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/changqings/kube-eventer/core"
)

func TestSimpleWrite(t *testing.T) {
	now := time.Now()
	event := kube_api.Event{
		Message:        "bzium",
		Count:          251,
		LastTimestamp:  metav1.NewTime(now),
		FirstTimestamp: metav1.NewTime(now),
	}
	batch := core.EventBatch{
		Timestamp: now,
		Events:    []*kube_api.Event{&event},
	}

	log := batchToString(&batch)
	fmt.Printf(log)

	assert.True(t, strings.Contains(log, "bzium"))
	assert.True(t, strings.Contains(log, "251"))
	assert.True(t, strings.Contains(log, fmt.Sprintf("%s", now)))
}
