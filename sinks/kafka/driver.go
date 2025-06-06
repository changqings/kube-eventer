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

package kafka

import (
	"encoding/json"
	"net/url"
	"sync"
	"time"

	"github.com/changqings/kube-eventer/util"

	"k8s.io/klog"

	kafka_common "github.com/changqings/kube-eventer/common/kafka"
	event_core "github.com/changqings/kube-eventer/core"
	"github.com/changqings/kube-eventer/metrics/core"
	kube_api "k8s.io/api/core/v1"
)

type KafkaSinkPoint struct {
	EventValue     interface{}
	EventTimestamp time.Time
	EventTags      map[string]string
}

type kafkaSink struct {
	kafka_common.KafkaClient
	sync.RWMutex
	KubeClusterName string
}

func getEventValue(event *kube_api.Event) ([]byte, error) {
	bytes, err := json.Marshal(event)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func eventToPoint(event *kube_api.Event, kubeClusterName string) (*KafkaSinkPoint, error) {
	value, err := getEventValue(event)
	if err != nil {
		return nil, err
	}
	point := KafkaSinkPoint{
		EventTimestamp: util.GetLastEventTimestamp(event).UTC(),
		EventValue:     value,
		EventTags: map[string]string{
			"eventID":         string(event.UID),
			"kubeClusterName": kubeClusterName,
		},
	}
	if event.InvolvedObject.Kind == "Pod" {
		point.EventTags[core.LabelPodId.Key] = string(event.InvolvedObject.UID)
		point.EventTags[core.LabelPodName.Key] = event.InvolvedObject.Name
	}
	point.EventTags[core.LabelHostname.Key] = event.Source.Host
	return &point, nil
}

func (sink *kafkaSink) ExportEvents(eventBatch *event_core.EventBatch) {
	sink.Lock()
	defer sink.Unlock()

	for _, event := range eventBatch.Events {
		point, err := eventToPoint(event, sink.KubeClusterName)
		if err != nil {
			klog.Warningf("Failed to convert event to point: %v", err)
		}

		err = sink.ProduceKafkaMessage(*point)
		if err != nil {
			klog.Errorf("Failed to produce event message: %s", err)
		}
	}
}

func NewKafkaSink(uri *url.URL, KuberClusterName string) (event_core.EventSink, error) {
	client, err := kafka_common.NewKafkaClient(uri, kafka_common.EventsTopic)
	if err != nil {
		return nil, err
	}

	return &kafkaSink{
		KafkaClient:     client,
		KubeClusterName: KuberClusterName,
	}, nil
}
