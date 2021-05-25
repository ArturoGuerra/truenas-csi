package truenasapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type (
	/* VolumeOpts (Name, type, size, sparse, sync, compression, etc) */
	VolumeOpts struct {
		Name         string `json:"name"`
		Type         string `json:"type"`
		VolSize      int64  `json:"volsize"`
		VolBlockSize string `json:"volblocksize"`
		Sparse       bool   `json:"sparse,omitempty"`
		Comments     string `json:"comments,omitempty"`
		Sync         string `json:"sync,omitempty"`
		Compression  string `json:"compression,omitempty"`
	}

	/* VolumeItem Generic struct for embeded items in volume */
	VolumeItem struct {
		Value    string `json:"value"`
		RawValue string `json:"rawvalue"`
		Parsed   string `json:"parsed"`
		Source   string `json:"source"`
	}

	/* VolumeItemInt Generic struct for embeded items in volume */
	VolumeItemInt struct {
		Value    string `json:"value"`
		RawValue string `json:"rawvalue"`
		Parsed   int    `json:"parsed"`
		Source   string `json:"source"`
	}

	/* VolumeItemBool Generic struct for embeded items in volume */
	VolumeItemBool struct {
		Value    string `json:"value"`
		RawValue string `json:"rawvalue"`
		Parsed   bool   `json:"parsed"`
		Source   string `json:"source"`
	}

	/* Volume (Size, Name, ID, Dataset, etc) */
	Volume struct {
		ID                  string         `json:"id" validate:"required"`
		Name                string         `json:"name" validate:"required"`
		Pool                string         `json:"pool" validate:"required"`
		Type                string         `json:"type" validate:"required"`
		Children            []*Volume      `json:"children" validate:"required"`
		Encrypted           bool           `json:"encrypted"`
		Comments            VolumeItem     `json:"comments" validate:"required"`
		ManagedBy           VolumeItem     `json:"managedby" validate:"required"`
		DeDuplication       VolumeItem     `json:"deduplication" validate:"required"`
		Sync                VolumeItem     `json:"sync" validate:"required"`
		Compression         VolumeItem     `json:"compression" validate:"required"`
		CompressRatio       VolumeItem     `json:"compressratio" validate:"required"`
		Origin              VolumeItem     `json:"origin" validate:"required"`
		Reservation         VolumeItem     `json:"reservation" validate:"required"`
		RefReservation      VolumeItemInt  `json:"refreservation" validate:"required"`
		Copies              VolumeItemInt  `json:"copies" validate:"required"`
		ReadOnly            VolumeItemBool `json:"readonly" validate:"required"`
		VolSize             VolumeItemInt  `json:"volsize" validate:"required"`
		VolBlockSize        VolumeItemInt  `json:"volblocksize" validate:"required"`
		KeyFormat           VolumeItem     `json:"key_format" validate:"required"`
		EncryptionAlgorithm VolumeItem     `json:"encryption_algorithm" validate:"required"`
		Used                VolumeItemInt  `json:"used" validate:"required"`
		Available           VolumeItemInt  `json:"available" validate:"required"`
		Pbkdf2iters         VolumeItem     `json:"pbkdf2iters" validate:"required"`
		Locked              bool           `json:"locked"`
	}
)

/*
   GetVolume take volume name/id returns volume object
   InternalError
   NotFoundError
*/
func (c *client) GetVolume(vol string) (*Volume, error) {
	url := c.parseurl(fmt.Sprintf("pool/dataset/id/%s", strings.ReplaceAll(vol, "/", "%2F")))
	resp, code, err := c.get(url)
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		v := new(Volume)
		// idk what would trigger this
		if err := json.Unmarshal(resp, v); err != nil {
			return nil, &InternalError{err}
		}

		// This error should only happen is struct validation fails and that should only happen if the volume is not found
		if err = c.Validate.Struct(v); err != nil {
			return nil, &InternalError{err}
		}

		return v, nil
	case 404:
		return nil, &NotFoundError{errors.New("Volume not found")}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}

	}

}

/*
CreateVolume takes volume spec returns volume object
Errors:
	InternalError
	AlreadyExistsError
*/
func (c *client) CreateVolume(vol VolumeOpts) (*Volume, error) {
	url := c.parseurl("pool/dataset")
	bytesData, err := json.Marshal(vol)
	if err != nil {
		return nil, &InternalError{err}
	}

	resp, code, err := c.post(url, bytes.NewBuffer(bytesData))
	if err != nil {
		return nil, &InternalError{err}
	}

	switch code {
	case 200:
		v := new(Volume)
		if err = json.Unmarshal(resp, v); err != nil {
			return nil, &InternalError{err}
		}

		if err = c.Validate.Struct(v); err != nil {
			return nil, &InternalError{err}
		}

		return v, nil

	case 422:
		return nil, &AlreadyExistsError{errors.New(string(resp))}

	default:
		return nil, &InternalError{fmt.Errorf("Error Code: %d\n", code)}
	}
}

/*
DeleteVolume takes volume name/id
Errors:
	ResourceBusyError
	InternalError
*/
func (c *client) DeleteVolume(vol string) error {
	todelvol, err := c.GetVolume(vol)
	if err != nil {
		return &InternalError{err}
	}

	if len(todelvol.Children) > 0 {
		return &ResourceBusyError{fmt.Errorf("Volume has %d children\n", len(todelvol.Children))}
	}

	url := c.parseurl(fmt.Sprintf("pool/dataset/id/%s", strings.ReplaceAll(vol, "/", "%2F")))
	_, code, err := c.delete(url)
	if err != nil {
		return &InternalError{err}
	}

	if code != 200 {
		return &InternalError{fmt.Errorf("Error Code: %d", code)}
	}

	return nil
}
