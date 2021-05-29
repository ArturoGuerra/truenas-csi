package truenasapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type (
	/* EXTERNAL USE STRUCTS */

	/* TargetGroupOpts */
	TargetGroup struct {
		Portal     int    `json:"portal"`
		Initiator  int    `json:"initiator"`
		Auth       string `json:"auth,omitempty"`
		AuthMethod string `json:"authmethod,omitempty"`
	}

	/* TargetOpts */
	TargetOpts struct {
		Name   string         `json:"name"`
		Alias  string         `json:"alias"`
		Mode   string         `json:"mode"`
		Groups []*TargetGroup `json:"groups"`
	}

	/* ExtentOpts */
	ExtentOpts struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Disk      string `json:"disk"`
		Path      string `json:"path"`
		Comment   string `json:"comment"`
		BlockSize int    `json:"blocksize"`
		Enabled   bool   `json:"enabled"`
		RO        bool   `json:"ro"`
	}

	/* TargetExtentOpts */
	TargetExtentOpts struct {
		LunID  int `json:"lunid"`
		Extent int `json:"extent"`
		Target int `json:"target"`
	}

	/* InitiatorOpts */
	InitiatorOpts struct {
		Initiators  []string `json:"initiators"`
		AuthNetwork []string `json:"auth_network"`
		Comment     string   `json:"comment"`
	}

	/* INTERNAL USE STRUCTS BELOW */

	/* Target */
	Target struct {
		ID     int            `json:"id"`
		Name   string         `json:"name"`
		Alias  string         `json:"alias"`
		Mode   string         `json:"mode"`
		Groups []*TargetGroup `json:"groups"`
	}

	/* Extent */
	Extent struct {
		ID         int    `json:"id"`
		Name       string `json:"name"`
		Serial     string `json:"serial"`
		Type       string `json:"type"`
		Path       string `json:"path"`
		BlockSize  int    `json:"blocksize"`
		PBlockSize bool   `json:"pblocksize"`
		Comment    string `json:"comment"`
		NAA        string `json:"naa"`
		RPM        string `json:"rpm"`
		RO         bool   `json:"ro"`
		Enabled    bool   `json:"enabled"`
		Vendor     string `json:"vendor"`
		Disk       string `json:"disk"`
		Locked     bool   `json:"locked"`
	}

	/* TargetExtent */
	TargetExtent struct {
		ID     int `json:"id"`
		LunID  int `json:"lunid"`
		Extent int `json:"extent"`
		Target int `json:"target"`
	}

	/* Global Target Global Configuration */
	Global struct {
		Basename           string   `json:"basename"`
		IsnsServers        []string `json:"isns_servers"`
		PoolAvailThreshold int      `json:"pool_avail_threshold"`
		Alua               bool     `json:"alua"`
	}

	PortalIP struct {
		IP   string `json:"ip"`
		Port int    `json:"port"`
	}

	Portal struct {
		ID                  int         `json:"id"`
		Tag                 string      `json:"tag"`
		Comment             string      `json:"comment"`
		Listen              []*PortalIP `json:"listen"`
		DiscoveryAuthmethod string      `json:"discovery_authmethod"`
		DiscoveryAuthgroup  string      `json:"discovery_authgroup"`
	}

	Initiator struct {
		ID          int      `json:"id"`
		Initiators  []string `json:"initiators"`
		AuthNetwork []string `json:"auth_network"`
		Comment     string   `json:"comment"`
	}

	/*
	 Common Interface
	 Handles top most layer api communication with other components
	*/

	ISCSIDeviceOpts struct {
		Name      string
		Alias     string
		Mode      string
		Portal    int
		Initiator int
		Disk      string
		Path      string
		BlockSize int
		Comment   string
		Enabled   bool
		LunID     int
		RO        bool
		Type      string
	}

	/* ISCSIDevice */
	ISCSIDevice struct {
		Name           string
		Node           *Node
		Mode           string
		Disk           string
		BlockSize      int
		Comment        string
		Enabled        bool
		RO             bool
		TargetID       int
		ExtentID       int
		TargetExtentID int
		IQN            string
		LunID          int
	}
)

func (c *client) GetISCSIID(volpath string) string {
	split := strings.Split(volpath, "/")
	return split[len(split)-1]
}

/* GetISCSIDevice gets (Target, Extent, and TargetExtent Association) */
func (c *client) GetISCSIDevice(vol string) (*ISCSIDevice, error) {
	return nil, nil
	/*	targets, err := c.getTargets()
		if err != nil {
			return nil, err
		}

		extents, err := c.getExtents()
		if err != nil {
			return nil, err
		}

		var extent *Extent
		for _, v := range extents {
			if v.Name == vol {
				extent = v
				break
			}
		}

		var target *Target
		for _, v := range targets {
			if v.Name == vol {
				target = v
				break
			}
		}

		if extent == nil || target == nil {
			return nil, &NotFoundError{errors.New("Unable to find ISCSI Device")}
		}

		targetextent, err := c.getTargetExtent(target.ID, extent.ID)
		if err != nil {
			return nil, err
		}

		return &ISCSIDevice{
			Name:           target.Name,
			Alias:          target.Alias,
			Mode:           target.Mode,
			Portal:         target.Groups[0].Portal,
			Initiator:      target.Groups[0].Initiator,
			Disk:           extent.Disk,
			Path:           extent.Path,
			BlockSize:      extent.BlockSize,
			Comment:        extent.Comment,
			Enabled:        extent.Enabled,
			RO:             extent.RO,
			LunID:          targetextent.LunID,
			TargetID:       target.ID,
			ExtentID:       extent.ID,
			TargetExtentID: targetextent.ID,
		}, nil
	*/
}

/* CreateISCSIDevice creates (Target, Extent, and TargetExtent Association) */
func (c *client) CreateISCSIDevice(dev ISCSIDeviceOpts) (*ISCSIDevice, error) {
	return nil, nil
	/*
		group := &TargetGroup{
			Portal:    dev.Portal,
			Initiator: dev.Initiator,
		}

		targetopts := TargetOpts{
			Name:  dev.Name,
			Alias: dev.Alias,
			Mode:  dev.Mode,
			Groups: []*TargetGroup{
				group,
			},
		}

		extentopts := ExtentOpts{
			Name:      dev.Name,
			Type:      dev.Type,
			Disk:      dev.Disk,
			Path:      dev.Path,
			Comment:   dev.Comment,
			BlockSize: dev.BlockSize,
			Enabled:   dev.Enabled,
			RO:        dev.RO,
		}

		target, err := c.createTarget(targetopts)
		if err != nil {
			return nil, err
		}

		extent, err := c.createExtent(extentopts)
		if err != nil {
			return nil, err
		}

		targetextentopts := TargetExtentOpts{
			LunID:  dev.LunID,
			Target: target.ID,
			Extent: extent.ID,
		}

		targetextent, err := c.createTargetExtent(targetextentopts)
		if err != nil {
			return nil, err
		}

		return &ISCSIDevice{
			Name:           target.Name,
			Alias:          target.Alias,
			Mode:           target.Mode,
			Portal:         string(target.Groups[0].Portal),
			Initiator:      target.Groups[0].Initiator,
			Disk:           extent.Disk,
			Path:           extent.Path,
			BlockSize:      extent.BlockSize,
			Comment:        extent.Comment,
			Enabled:        extent.Enabled,
			RO:             extent.RO,
			LunID:          targetextent.LunID,
			TargetID:       target.ID,
			ExtentID:       extent.ID,
			TargetExtentID: targetextent.ID,
		}, nil
	*/
}

/* DeleteISCSIDevice deletes (Target, Extent, and TargetExtent Association) */
func (c *client) DeleteISCSIDevice(volID string) error {
	/*
		iscsi, err := c.GetISCSIDevice(volID)
		if err != nil {
			return err
		}

		if err := c.deleteTargetExtent(iscsi.TargetExtentID); err != nil {
			return err
		}

		if err := c.deleteTarget(iscsi.TargetID); err != nil {
			return err
		}

		if err := c.deleteExtent(iscsi.ExtentID); err != nil {
			return err
		}
	*/
	return nil
}

func (c *client) getTargets() ([]*Target, error) {
	url := c.parseurl("iscsi/target")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var t []*Target
		if err = json.Unmarshal(resp, &t); err != nil {
			return nil, &InternalError{err}
		}

		return t, nil
	case 404:
		return nil, &NotFoundError{errors.New("Target not found")}
	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d", code)}
	}
}

func (c *client) createTarget(tgt TargetOpts) (*Target, error) {
	url := c.parseurl("iscsi/target")
	bytesData, err := json.Marshal(tgt)
	if err != nil {
		return nil, &InternalError{err}
	}

	resp, code, err := c.post(url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var t Target
		if err = json.Unmarshal(resp, &t); err != nil {
			return nil, &InternalError{err}
		}

		return &t, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d", code)}
	}
}

func (c *client) deleteTarget(tgt int) error {
	url := c.parseurl(fmt.Sprintf("iscsi/target/id/%d", tgt))
	resp, code, err := c.delete(url)
	if err != nil {
		return &InternalError{err}
	}

	switch code {
	case 200:
		return nil
	case 404:
		return &NotFoundError{errors.New("Target not found")}
	default:
		return &InternalError{fmt.Errorf("Error Code: %d Message: %s", code, string(resp))}
	}
}

func (c *client) getExtents() ([]*Extent, error) {
	url := c.parseurl("iscsi/extent")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var e []*Extent
		if err = json.Unmarshal(resp, &e); err != nil {
			return nil, &InternalError{err}
		}

		return e, nil
	case 404:
		return nil, &NotFoundError{errors.New("extent not found")}
	default:
		return nil, &InternalError{fmt.Errorf("error code: %d", code)}
	}
}

func (c *client) createExtent(ext ExtentOpts) (*Extent, error) {
	url := c.parseurl("iscsi/extent")
	bytesData, err := json.Marshal(ext)
	if err != nil {
		return nil, &InternalError{err}
	}

	resp, code, err := c.post(url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var e Extent
		if err = json.Unmarshal(resp, &e); err != nil {
			return nil, &InternalError{err}
		}

		return &e, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("error code: %d", code)}
	}

}

func (c *client) deleteExtent(ext int) error {
	url := c.parseurl(fmt.Sprintf("iscsi/extent/id/%d", ext))
	resp, code, err := c.delete(url)
	if err != nil {
		return &InternalError{err}
	}

	switch code {
	case 200:
		return nil
	case 404:
		return &NotFoundError{errors.New("Extent not found")}
	default:
		return &InternalError{fmt.Errorf("Error Code: %d Message: %s", code, string(resp))}
	}
}

func (c *client) getTargetExtent(tgt, ext int) (*TargetExtent, error) {
	url := c.parseurl("iscsi/targetextent")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var t []*TargetExtent
		if err = json.Unmarshal(resp, &t); err != nil {
			return nil, &InternalError{err}
		}

		for _, v := range t {
			if v.Target == tgt && v.Extent == ext {
				return v, nil
			}
		}

		return nil, &NotFoundError{errors.New("unable to find target extent")}
	case 404:
		return nil, &NotFoundError{errors.New("targetExtent not found")}
	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d", code)}
	}
}

func (c *client) createTargetExtent(tgtext TargetExtentOpts) (*TargetExtent, error) {
	url := c.parseurl("iscsi/targetextent")
	bytesData, err := json.Marshal(tgtext)
	if err != nil {
		return nil, &InternalError{err}
	}

	resp, code, err := c.post(url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var te TargetExtent
		if err = json.Unmarshal(resp, &te); err != nil {
			return nil, &InternalError{err}
		}

		return &te, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d", code)}
	}
}

func (c *client) deleteTargetExtent(tgtext int) error {
	url := c.parseurl(fmt.Sprintf("iscsi/targetextent/id/%d", tgtext))
	resp, code, err := c.delete(url)
	if err != nil {
		return &InternalError{err}
	}

	switch code {
	case 200:
		return nil

	case 404:
		return &NotFoundError{errors.New("TargetExtent not found")}
	default:
		return &InternalError{fmt.Errorf("Error Code: %d Message: %s", code, string(resp))}
	}
}

func (c *client) getGlobal() (*Global, error) {
	url := c.parseurl("iscsi/global")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	if code == 200 {
		var global Global
		if err := json.Unmarshal(resp, &global); err != nil {
			return nil, &InternalError{err}
		}

		return &global, nil
	}

	return nil, &InternalError{fmt.Errorf("Error Code: %d", code)}

}

func (c *client) getPortal() (*Portal, error) {
	url := c.parseurl("iscsi/portal")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var portals []*Portal
		if err := json.Unmarshal(resp, &portals); err != nil {
			return nil, &InternalError{err}
		}

		for _, p := range portals {
			if p.Comment == "truenas-csi" {
				return p, nil
			}
		}

		return nil, &NotFoundError{errors.New("couldn't find portal with name truenas-csi")}
	default:
		return nil, &InternalError{fmt.Errorf("error code: %d", code)}
	}
}

func (c *client) getIntiator(name string) (*Initiator, error) {
	url := c.parseurl("iscsi/initiator")
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var initiators []*Initiator
		if err := json.Unmarshal(resp, &initiators); err != nil {
			return nil, &InternalError{err}
		}

		for _, i := range initiators {
			if i.Comment == name {
				return i, nil
			}
		}

		return nil, &NotFoundError{errors.New("initiator not found")}
	default:
		return nil, &InternalError{fmt.Errorf("error code: %d", code)}
	}
}

func (c *client) createInitiator(opts InitiatorOpts) (*Initiator, error) {
	url := c.parseurl("iscsi/initiator")
	bytesData, err := json.Marshal(opts)
	if err != nil {
		return nil, &InternalError{err}
	}

	resp, code, err := c.post(url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		var initiator Initiator
		if err := json.Unmarshal(resp, &initiator); err != nil {
			return nil, &InternalError{err}
		}

		return &initiator, nil
	default:
		return nil, &InternalError{fmt.Errorf("error code: %d", code)}
	}
}

func (c *client) deleteInitiator(name string) error {
	initiator, err := c.getIntiator(name)
	if err != nil {
		switch err := err.(type) {
		case *NotFoundError:
			return nil
		default:
			return &InternalError{err}
		}
	}

	url := c.parseurl(fmt.Sprintf("iscsi/initiator/id/%d", initiator.ID))
	_, code, err := c.delete(url)
	if err != nil {
		return &InternalError{err}
	}

	switch code {
	case 200:
		return nil
	case 404:
		return nil
	default:
		return &InternalError{err}
	}
}
