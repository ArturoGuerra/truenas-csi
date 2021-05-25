package truenasapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

type (
	/* EXTERNAL USE STRUCTS */

	/* ISCSIDevice */
	ISCSIDevice struct {
		Name           string
		Alias          string
		Mode           string
		Portal         int
		Initiator      int
		Disk           string
		Path           string
		BlockSize      int
		Comment        string
		Enabled        bool
		RO             bool
		LunID          int
		TargetID       int
		ExtentID       int
		TargetExtentID int
	}

	/* ISCSIDeviceOpts */
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

	/* INTERNAL USE STRUCTS BELOW */

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
)

/* GetISCSIDevice gets (Target, Extent, and TargetExtent Association) */
func (c *client) GetISCSIDevice(tid, eid, teid int) (*ISCSIDevice, error) {
	target, err := c.getTarget(tid)
	if err != nil {
		return nil, err
	}

	extent, err := c.getExtent(eid)
	if err != nil {
		return nil, err
	}

	targetextent, err := c.getTargetExtent(teid)
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

}

/* CreateISCSIDevice creates (Target, Extent, and TargetExtent Association) */
func (c *client) CreateISCSIDevice(dev ISCSIDeviceOpts) (*ISCSIDevice, error) {
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
}

/* DeleteISCSIDevice deletes (Target, Extent, and TargetExtent Association) */
func (c *client) DeleteISCSIDevice(tid, eid, teid int) error {
	if err := c.deleteTargetExtent(teid); err != nil {
		return err
	}

	if err := c.deleteTarget(tid); err != nil {
		return err
	}

	if err := c.deleteExtent(eid); err != nil {
		return err
	}

	return nil
}

func (c *client) getTarget(tgt int) (*Target, error) {
	url := c.parseurl(fmt.Sprintf("iscsi/target/id/%d", tgt))
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		t := new(Target)
		if err = json.Unmarshal(resp, t); err != nil {
			return nil, &InternalError{err}
		}

		return t, nil
	case 404:
		return nil, &NotFoundError{errors.New("Target not found")}
	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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
		t := new(Target)
		if err = json.Unmarshal(resp, t); err != nil {
			return nil, &InternalError{err}
		}

		return t, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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

func (c *client) getExtent(ext int) (*Extent, error) {
	url := c.parseurl(fmt.Sprintf("iscsi/extent/id/%d", ext))
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		e := new(Extent)
		if err = json.Unmarshal(resp, e); err != nil {
			return nil, &InternalError{err}
		}

		return e, nil
	case 404:
		return nil, &NotFoundError{errors.New("Extent not found")}
	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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
		e := new(Extent)
		if err = json.Unmarshal(resp, e); err != nil {
			return nil, &InternalError{err}
		}

		return e, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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

func (c *client) getTargetExtent(tgtext int) (*TargetExtent, error) {
	url := c.parseurl(fmt.Sprintf("iscsi/targetextent/id/%d", tgtext))
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		te := new(TargetExtent)
		if err = json.Unmarshal(resp, te); err != nil {
			return nil, &InternalError{err}
		}

		return te, nil
	case 404:
		return nil, &NotFoundError{errors.New("TargetExtent not found")}
	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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
		te := new(TargetExtent)
		if err = json.Unmarshal(resp, te); err != nil {
			return nil, &InternalError{err}
		}

		return te, nil
	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
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
