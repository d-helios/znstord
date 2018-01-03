package stmf

const (
	// RDSK_DEFAULT_PREFIX - default representation of zfs structure under /dev/ directory.
	RDSK_DEFAULT_PREFIX = "/dev/zvol/rdsk/"
	defaultVolBlockSize = "4096"
)

// Enable/Disable RSF-1 Comstar features
var StmfHAEnabled = true

// LogicalUnit - representation of stmfadm list-lu -v output.
type LogicalUnit struct {
	LUName               string `json:"LUName"`
	OperationalStatus    string `json:"OperationalStatus"`
	ProviderName         string `json:"ProviderName"`
	Alias                string `json:"Alias"`
	ViewEntryCount       uint64 `json:"ViewEntryCount"`
	DataFile             string `json:"DataFile"`
	MetaFile             string `json:"MetaFile"`
	Size                 uint64 `json:"Size"`
	BlockSize            uint64 `json:"BlockSize"`
	ManagementURL        string `json:"ManagementURL"`
	VendorID             string `json:"VendorID"`
	ProductID            string `json:"ProductID"`
	SerialNum            string `json:"SerialNum"`
	WriteProtect         string `json:"WriteProtect"`
	WriteCacheModeSelect string `json:"WriteCacheModeSelect"` // This field doesn't present in OpenSolaris
	WritebackCache       string `json:"WritebackCache"`
	AccessState          string `json:"AccessState"`
}

// TargetGroup - representation of stmfadm list-tg
type TargetGroup struct {
	TargetGroup     string   `json:"TargetGroup"`
	TargetPortGroup []string `json:"TargetPortGroup"`
}

// View - representation of stmadm list-view -l wwid
type View struct {
	ViewEntry   uint64 `json:"ViewEntry"`
	HostGroup   string `json:"HostGroup"`
	TargetGroup string `json:"TargetGroup"`
	LUN         uint64 `json:"LUN"`
}

// HostGroup - representation of smtfadm list-hg output
type HostGroup struct {
	HostGroup string   `json:"HostGroup"`
	Members   []string `json:"Members"`
}
