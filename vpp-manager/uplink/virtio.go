// Copyright (C) 2020 Cisco Systems Inc.
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

package uplink

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/projectcalico/vpp-dataplane/vpp-manager/config"
	"github.com/projectcalico/vpp-dataplane/vpp-manager/utils"
	log "github.com/sirupsen/logrus"
)

type VirtioDriver struct {
	*UplinkDriverData
}

func (d *VirtioDriver) PreconfigureLinux() (err error) {
	if d.conf.IsUp {
		// Set interface down if it is up, bind it to a VPP-friendly driver
		err := utils.SafeSetInterfaceDownByName(d.params.MainInterface)
		if err != nil {
			return err
		}
	}
	if d.conf.DoSwapDriver {
		if d.conf.PciId == "" {
			log.Warnf("PCI ID not found, not swapping drivers")
		} else {
			err = utils.SwapDriver(d.conf.PciId, d.params.NewDriverName, true)
			if err != nil {
				log.Warnf("Failed to swap driver to %s: %v", d.params.NewDriverName, err)
			}
		}
	}
	return nil
}

func (d *VirtioDriver) RestoreLinux() {
	if d.conf.PciId != "" && d.conf.Driver != "" {
		err := utils.SwapDriver(d.conf.PciId, d.conf.Driver, false)
		if err != nil {
			log.Warnf("Error swapping back driver to %s for %s: %v", d.conf.Driver, d.conf.PciId, err)
		}
	}
	if !d.conf.IsUp {
		return
	}
	// This assumes the link has kept the same name after the rebind.
	// It should be always true on systemd based distros
	link, err := utils.SafeSetInterfaceUpByName(d.params.MainInterface)
	if err != nil {
		log.Warnf("Error setting %s up: %v", d.params.MainInterface, err)
		return
	}

	// Re-add all adresses and routes
	d.restoreLinuxIfConf(link)
}

func (d *VirtioDriver) CreateMainVppInterface() (err error) {
	swIfIndex, err := d.vpp.CreateVirtio(d.conf.PciId, &d.conf.HardwareAddr)
	if err != nil {
		return errors.Wrapf(err, "Error creating VIRTIO interface")
	}
	log.Infof("Created VIRTIO interface %d", swIfIndex)

	if swIfIndex != config.DataInterfaceSwIfIndex {
		return fmt.Errorf("Created VIRTIO interface has wrong swIfIndex %d!", swIfIndex)
	}
	return nil
}
