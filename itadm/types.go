package itadm

type TargetPortGroup struct {
	TargetPortGroup string   `json:"tpg"`
	Count           uint64   `json:"count"`
	Portals         []string `json:"portals"`
}

type Target struct {
	IQN              string `json:"iqn"`
	State            string `json:"state"`
	Sessions         uint64 `json:"sessions"`
	Alias            string `json:"alias"`
	Auth             string `json:"auth"`
	TargetChapUser   string `json:"chapuser"`
	TargetChapSecret string `json:"chapsecret"`
	TpgTags          string `json:"tpg-tags"`
}
