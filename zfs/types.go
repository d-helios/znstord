package zfs

/*
Basic properties list available both in OpenSolaris and OracleSolaris.
*/

const (
	// Filesystem zfs dataset
	Filesystem = "filesystem"

	// Volume zfs dataset aka zvol
	Volume = "volume"

	// Snapshot zfs dataset
	Snapshot = "snapshot"

	// RootPool - is a os root dataset. should not be modified
	RootPool = "rpool" // TODO: Need to add FreeBsd compability

	// OpenSolaris - for OpenZfs compability behaviar
	OpenSolaris = "OpenSolaris"

	// OracleSolaris - for original zfs compability behaviar
	OracleSolaris = "OracleSolaris"
)

// NumericProps - list of the zfs numeric properties.
var NumericProps = []string{
	"available",
	"copies",
	"recordsize",
	"quota",
	"referenced",
	"refquota",
	"reservation",
	"refreservation",
	"used",
	"usedbychildren",
	"usedbydataset",
	"usedbyrefreservation",
	"usedbysnapshots",
	"volblocksize",
	"volsize"}

// FloatProps - list of the zfs floating properties.
var FloatProps = []string{
	"compressratio",
}

// Dataset - common dataset structure.
type Dataset struct {
	Dataset string      `json:"dataset"`
	Props   interface{} `json:"options"`
}

// FsDataset - zfs filesystem dataset.
type FsDataset struct {
	Type                 string  `json:"type"`
	Creation             string  `json:"creation"`
	Used                 uint64  `json:"used"`
	Available            uint64  `json:"available"`
	Referenced           uint64  `json:"referenced"`
	Compressratio        float32 `json:"compressratio"`
	Mounted              string  `json:"mounted"`
	Origin               string  `json:"origin"`
	Quota                uint64  `json:"quota"`
	Reservation          uint64  `json:"reservation"`
	Recordsize           uint64  `json:"recordsize"`
	Mountpoint           string  `json:"mountpoint"`
	Sharenfs             string  `json:"sharenfs"`
	Checksum             string  `json:"checksum"`
	Compression          string  `json:"compression"`
	Atime                string  `json:"atime"`
	Devices              string  `json:"devices"`
	Exec                 string  `json:"exec"`
	Setuid               string  `json:"setuid"`
	Readonly             string  `json:"readonly"`
	Zoned                string  `json:"zoned"`
	Snapdir              string  `json:"snapdir"`
	Aclmode              string  `json:"aclmode"`
	Aclinherit           string  `json:"aclinherit"`
	Canmount             string  `json:"canmount"`
	Xattr                string  `json:"xattr"`
	Copies               uint64  `json:"copies"`
	Version              string  `json:"version"`
	Utf8only             string  `json:"utf8only"`
	Normalization        string  `json:"normalization"`
	Casesensitivity      string  `json:"casesensitivity"`
	Vscan                string  `json:"vscan"`
	Nbmand               string  `json:"nbmand"`
	Sharesmb             string  `json:"sharesmb"`
	Refquota             uint64  `json:"refquota"`
	Refreservation       uint64  `json:"refreservation"`
	Primarycache         string  `json:"primarycache"`
	Secondarycache       string  `json:"secondarycache"`
	Usedbysnapshots      uint64  `json:"usedbysnapshots"`
	Usedbydataset        uint64  `json:"usedbydataset"`
	Usedbychildren       uint64  `json:"usedbychildren"`
	Usedbyrefreservation uint64  `json:"usedbyrefreservation"`
	Logbias              string  `json:"logbias"`
	Dedup                string  `json:"dedup"`
	Mlslabel             string  `json:"mlslabel"`
	Sync                 string  `json:"sync"`
	SFlag                string  `json:"service_flag"`
}

// VolDataset - zfs zvol dataset
type VolDataset struct {
	Type                 string  `json:"type"`
	Creation             string  `json:"creation"`
	Used                 uint64  `json:"used"`
	Available            uint64  `json:"available"`
	Referenced           uint64  `json:"referenced"`
	Compressratio        float32 `json:"compressratio"`
	Reservation          uint64  `json:"reservation"`
	Volsize              uint64  `json:"volsize"`
	Volblocksize         uint64  `json:"volblocksize"`
	Checksum             string  `json:"checksum"`
	Compression          string  `json:"compression"`
	Origin               string  `json:"origin"`
	Readonly             string  `json:"readonly"`
	Copies               uint64  `json:"copies"`
	Refreservation       uint64  `json:"refreservation"`
	Primarycache         string  `json:"primarycache"`
	Secondarycache       string  `json:"secondarycache"`
	Usedbysnapshots      uint64  `json:"usedbysnapshots"`
	Usedbydataset        uint64  `json:"usedbydataset"`
	Usedbychildren       uint64  `json:"usedbychildren"`
	Usedbyrefreservation uint64  `json:"usedbyrefreservation"`
	Logbias              string  `json:"logbias"`
	Dedup                string  `json:"dedup"`
	Mlslabel             string  `json:"mlslabel"`
	Sync                 string  `json:"sync"`
	Alias                string  `json:"alias"`
	SFlag                string  `json:"service_flag"`
}

// SnapDataset  - zfs snapshot dataset.
type SnapDataset struct {
	Type            string  `json:"type"`
	Creation        string  `json:"creation"`
	Used            uint64  `json:"used"`
	Referenced      uint64  `json:"referenced"`
	Compressratio   float32 `json:"compressratio"`
	Devices         string  `json:"devices"`
	Exec            string  `json:"exec"`
	Setuid          string  `json:"setuid"`
	Xattr           string  `json:"xattr"`
	Version         string  `json:"version"`
	Utf8only        string  `json:"utf8only"`
	Normalization   string  `json:"normalization"`
	Casesensitivity string  `json:"casesensitivity"`
	Nbmand          string  `json:"nbmand"`
	Primarycache    string  `json:"primarycache"`
	Secondarycache  string  `json:"secondarycache"`
	DeferDestroy    string  `json:"defer_destroy"`
	Userrefs        string  `json:"userrefs"`
	Mlslabel        string  `json:"mlslabel"`
	Alias           string  `json:"custom:alias"`
	SFlag           string  `json:"service_flag"`
}
