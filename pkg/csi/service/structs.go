package service

type (
	// Params Example: StorageClass Parameters on k8s
	Params struct {
		Dataset     string `json:"dataset"`
		FSType      string `json:"fstype"`
		Compression string `json:"compression"`
	}

	/* VolumeContext
	Volume Name
	FSType
	*/
	VolumeContext struct {
		Name    string `json:"name"`
		Dataset string `json:"dataset"`
		FSType  string `json:"fstype"`
	}

	/* PublishContext
	Full Iqn path
	portal address
	lun id iqn lun
	*/
	PublishContext struct {
		IQN    string `json:"iqn"`
		Portal string `json:"portal"`
		LunID  int    `json:"lunid"`
	}
)
