// Copyright (C) 2021 Cisco Systems Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type CalicoVppEventType string

const (
	ChanSize = 500

	PeerNodeStateChanged CalicoVppEventType = "PeerNodeStateChanged"
	FelixConfChanged     CalicoVppEventType = "FelixConfChanged"
	IpamConfChanged      CalicoVppEventType = "IpamConfChanged"
	BGPConfChanged       CalicoVppEventType = "BGPConfChanged"

	ConnectivityAdded   CalicoVppEventType = "ConnectivityAdded"
	ConnectivityDeleted CalicoVppEventType = "ConnectivityDeleted"

	SRv6PolicyAdded   CalicoVppEventType = "SRv6PolicyAdded"
	SRv6PolicyDeleted CalicoVppEventType = "SRv6PolicyDeleted"

	PodAdded   CalicoVppEventType = "PodAdded"
	PodDeleted CalicoVppEventType = "PodDeleted"

	LocalPodAddressAdded   CalicoVppEventType = "LocalPodAddressAdded"
	LocalPodAddressDeleted CalicoVppEventType = "LocalPodAddressDeleted"

	TunnelAdded   CalicoVppEventType = "TunnelAdded"
	TunnelDeleted CalicoVppEventType = "TunnelDeleted"

	BGPPeerAdded     CalicoVppEventType = "BGPPeerAdded"
	BGPPeerDeleted   CalicoVppEventType = "BGPPeerDeleted"
	BGPPeerUpdated   CalicoVppEventType = "BGPPeerUpdated"
	BGPSecretChanged CalicoVppEventType = "BGPSecretChanged"

	BGPDefinedSetAdded   CalicoVppEventType = "BGPDefinedSetAdded"
	BGPDefinedSetDeleted CalicoVppEventType = "BGPDefinedSetDeleted"

	BGPPathAdded   CalicoVppEventType = "BGPPathAdded"
	BGPPathDeleted CalicoVppEventType = "BGPPathDeleted"

	NetAddedOrUpdated CalicoVppEventType = "NetAddedOrUpdated"
	NetDeleted        CalicoVppEventType = "NetDeleted"
	NetsSynced        CalicoVppEventType = "NetsSynced"

	IpamPoolUpdate CalicoVppEventType = "IpamPoolUpdate"
	IpamPoolRemove CalicoVppEventType = "IpamPoolRemove"

	WireguardPublicKeyChanged CalicoVppEventType = "WireguardPublicKeyChanged"
)

var (
	ThePubSub *PubSub
)

type CalicoVppEvent struct {
	Type CalicoVppEventType

	Old interface{}
	New interface{}
}

type PubSubHandlerRegistration struct {
	/* Name for the registration, for logging & debugging */
	name string
	/* Channel where to send events */
	channel chan CalicoVppEvent
	/* Receive only these events. If empty we'll receive all */
	expectedEvents map[CalicoVppEventType]bool
	/* Receive all events */
	expectAllEvents bool
}

func (reg *PubSubHandlerRegistration) ExpectEvents(eventTypes ...CalicoVppEventType) {
	for _, eventType := range eventTypes {
		reg.expectedEvents[eventType] = true
	}
	reg.expectAllEvents = false
}

type PubSub struct {
	log                        *log.Entry
	pubSubHandlerRegistrations []*PubSubHandlerRegistration
}

func RegisterHandler(channel chan CalicoVppEvent, name string) *PubSubHandlerRegistration {
	reg := &PubSubHandlerRegistration{
		channel:         channel,
		name:            name,
		expectedEvents:  make(map[CalicoVppEventType]bool),
		expectAllEvents: true, /* By default receive everything, unless we ask for a filter */
	}
	ThePubSub.pubSubHandlerRegistrations = append(ThePubSub.pubSubHandlerRegistrations, reg)
	return reg
}

func redactPassword(event CalicoVppEvent) string {
	switch event.Type {
	case BGPPeerAdded:
		return string(event.Type)
	default:
		return fmt.Sprintf("%+v", event)
	}

}

func SendEvent(event CalicoVppEvent) {
	ThePubSub.log.Debugf("Broadcasting event %s", redactPassword(event))
	for _, reg := range ThePubSub.pubSubHandlerRegistrations {
		if reg.expectAllEvents || reg.expectedEvents[event.Type] {
			reg.channel <- event
		}
	}
}

func NewPubSub(log *log.Entry) *PubSub {
	return &PubSub{
		log:                        log,
		pubSubHandlerRegistrations: make([]*PubSubHandlerRegistration, 0),
	}
}
