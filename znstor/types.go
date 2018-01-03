package znstor

import (
	"sync"
)

// Constants
const (
	asyncResultDir                    = "/tmp/"
	asyncOptStatusInProgress          = "In Progress"
	asyncOptStatusCompletedSuccefully = "Completed Successfully"
	requestPayloadMaxSize             = 8192
	sflagManaged                      = "managed_by_znstor"
	sflagDeleting                     = "deleting"
)

var (
	privatePoolList = []string{"rpool", "zpool"}
	vol_mutex       = &sync.Mutex{}
)

// Configuration structure
type ServerData struct {
	Listen string   `json:"listen"` // Listen address. ex: 127.0.0.1
	Auth   AuthData `json:"auth"`
	StmfHa bool     `json:"stmfhaEnabled"`
}

type AuthData struct {
	UserName     string `json:"username"`
	UserPassword string `json:"password"`
}

// END

// Responce Structures
type RespMsg struct {
	Subject string `json:"subject"`
	Msg     string `json:"message"`
}

// Options volume create / update functions
type ZVolOptions struct {
	VolBlockSize uint64 `json:"volblocksize,omitempty"`
	Reservation  uint64 `json:"reservation,omitempty"`
	Dedup        string `json:"dedup,omitempty"`
	Compression  string `json:"compression,omitempty"`
	Thin         bool   `json:"thin,omitempty"`
}

// Options to represent filesystem dataset
type ZFilesystemOptions struct {
	Used                 uint64  `json:"used"`
	Available            uint64  `json:"available"`
	Referenced           uint64  `json:"referenced"`
	Quota                uint64  `json:"quota"`
	Refquota             uint64  `json:"refquota"`
	Reservation          uint64  `json:"reservation"`
	Refreservation       uint64  `json:"refreservation"`
	Origin               string  `json:"origin"`
	Compressratio        float32 `json:"compressratio"`
	Compression          string  `json:"compression"`
	Dedup                string  `json:"dedup"`
	Usedbysnapshots      uint64  `json:"usedbysnapshots"`
	Usedbydataset        uint64  `json:"usedbydataset"`
	Usedbychildren       uint64  `json:"usedbychildren"`
	Usedbyrefreservation uint64  `json:"usedbyrefreservation"`
}

// Projects representation
type Project struct {
	Dataset string             `json:"project"`
	Options ZFilesystemOptions `json:"options"`
}

// Filesystem representation
type Filesystem struct {
	Dataset string             `json:"filesystem"`
	Options ZFilesystemOptions `json:"options"`
}

/*
====================
 JSon Payload.
====================
*/
type ZVolCreateRequest struct {
	Alias   string      `json:"alias"`
	VolSize uint64      `json:"volsize"`
	Guid    string      `json:"guid,omitempty"`
	Serial  string      `json:"serial,omitempty"`
	Options ZVolOptions `json:"options,omitempty"`
}

type ZVolCloneRequest struct {
	Alias  string `json:"alias"`
	Serial string `json:"serial,omitempty"`
}

type ZvolResizeRequest struct {
	VolSize uint64 `json:"volsize"`
}

type FilesystemCloneRequest struct {
	Alias   string `json:"alias"`
	Dataset string `json:"filesystem"`
}

type FilesystemRequest struct {
	Alias          string `json:"alias,omitempty"`
	Quota          uint64 `json:"quota"`
	Refquota       uint64 `json:"refquota,omitempty"`
	Reservation    uint64 `json:"reservation,omitempty"`
	Refreservation uint64 `json:"refreservation,omitempty"`
	Compression    string `json:"compression,omitempty"`
	Dedup          string `json:"dedup,omitempty"`
	Atime          string `json:"atime,omitempty"`
}

type TpgCreateRequest struct {
	portals []string `json:"portals,omitempty"`
}

type TargetCreateRequest struct {
	Iqn   string `json:"iqn,omitempty"`
	Alias string `json:"alias,omitempty"`
	Tpg   string `json:"tpgs,omitempty"`
}

type ExportRequest struct {
	Hostgroup   string `json:"hostgroup,omitempty"`
	Targetgroup string `json:"targetgroup,omitempty"`
	Lun         int64  `json:"lun,omitempty"`
}
