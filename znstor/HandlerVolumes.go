package znstor

import (
	"encoding/json"
	"github.com/d-helios/znstord/stmf"
	"github.com/d-helios/znstord/zfs"
	"github.com/gorilla/mux"
	"github.com/twinj/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// HandlerGetVolumeList - get volume list

func HandlerGetVolumeJobStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	statusUuid := vars["uuid"]

	status, err := ioutil.ReadFile(asyncResultDir + statusUuid + ".status")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	sendMessage(w, http.StatusOK, traceFunctionName(), string(status))
}

func HandlerGetVolumeList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	// lock mutex to perform atomic view
	vol_mutex.Lock()

	lus, err := stmf.ListLUs("")
	if err != nil {
		vol_mutex.Unlock()
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	// unlock mutex
	vol_mutex.Unlock()

	// only lu's associeted with project
	var projectLus = make([]stmf.LogicalUnit, 0)

	for _, lu := range lus {
		if IsVolumeBelongsToProject(basepath, *lu) {
			projectLus = append(projectLus, *lu)
		}
	}

	err = json.NewEncoder(w).Encode(projectLus)

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Get Volume
func HandlerGetVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName
	volumeName := vars["volume"]

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if IsVolumeBelongsToProject(basepath, *lu) {
		err := json.NewEncoder(w).Encode(lu)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
		return
	}

	sendMessage(w, http.StatusBadRequest, traceFunctionName(),
		"Volume not found in specified project")
	return
}

// CreateVolume
func HandlerCreateVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]

	var reqJson ZVolCreateRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := CreateVolume(basepath, reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(lu))
	err = json.NewEncoder(w).Encode(lu)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// VolDestroy Volume
func HandlerDestroyVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(),
			"Volume not found in specified project")
		return
	}

	requestUuid := uuid.NewV4().String()

	go func() {
		// status uuid
		// request status file
		statusFile := asyncResultDir + requestUuid + ".status"

		// start logging
		ioutil.WriteFile(statusFile, []byte(asyncOptStatusInProgress), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// destroy volume
		if err := VolDestroy(lu.LUName); err != nil {
			// log operation failed
			ioutil.WriteFile(statusFile, []byte(err.Error()), 0644)
		} else {
			// log operation successfully
			ioutil.WriteFile(statusFile, []byte(asyncOptStatusCompletedSuccefully), 0644)
		}
	}()

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(requestUuid))
	sendMessage(w, http.StatusAccepted, traceFunctionName(), requestUuid)
}

// Create VolumeSnapshot
func HandlerCreateVolumeSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	snapshotName := vars["snapshot"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	snapshot, err := VolSnapshot(lu.LUName, snapshotName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(snapshot))
	err = json.NewEncoder(w).Encode(snapshot)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Get Volume VolSnapshot List
func HandlerGetVolumeSnapshotList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	snapshots, err := zfs.ListDatasets(zfs.Snapshot, lu.GetZvol(), true, 1)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if snapshots == nil {
		// return empty array
		log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode([]string{}))
		err = json.NewEncoder(w).Encode([]string{})
	} else {
		log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(lu))
		err = json.NewEncoder(w).Encode(snapshots)
	}

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Get Volume VolSnapshot
func HandlerGetVolumeSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	snapshotName := vars["snapshot"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	snapshot, err := zfs.GetDataset(lu.GetZvol() + "@" + snapshotName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(snapshot))
	err = json.NewEncoder(w).Encode(snapshot)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// VolDestroy Volume VolSnapshot
func HandlerDestroyVolumeSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	snapshotName := vars["snapshot"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	snapshot, err := zfs.GetDataset(lu.GetZvol() + "@" + snapshotName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	requestUuid := uuid.NewV4().String()

	go func() {
		// status uuid
		// request status file
		statusFile := asyncResultDir + requestUuid + ".status"

		// start logging
		ioutil.WriteFile(statusFile, []byte(asyncOptStatusInProgress), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// destroy snapshot
		if err := snapshot.Destroy(""); err != nil {
			// log operation failed
			ioutil.WriteFile(statusFile, []byte(err.Error()), 0644)
		} else {
			// log operation successfully
			ioutil.WriteFile(statusFile, []byte(asyncOptStatusCompletedSuccefully), 0644)
		}
	}()
	sendMessage(w, http.StatusAccepted, traceFunctionName(), requestUuid)
}

// VolRollback Volume VolSnapshot
func HandlerRollbackVolumeSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	snapshotName := vars["snapshot"]
	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	err = VolRollback(lu.LUName, snapshotName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(lu))
	err = json.NewEncoder(w).Encode(lu)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Clone Volume From VolSnapshot
func HandlerCloneVolumeFromSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	snapshotName := vars["snapshot"]
	basepath := poolName + "/" + domainName + "/" + projectName

	var reqJson ZVolCloneRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if reqJson.Alias == "" {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Alias not specified")
		return
	}

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	clone, err := VolCloneFromSnapshot(lu.LUName, snapshotName, reqJson.Alias)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(clone))
	err = json.NewEncoder(w).Encode(clone)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// VolResize Volume
func HandlerResizeVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	basepath := poolName + "/" + domainName + "/" + projectName

	var reqJson ZvolResizeRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	volume, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *volume) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	err = VolResize(volume.LUName, reqJson.VolSize)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(volume))
	err = json.NewEncoder(w).Encode(volume)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Set Volume Compression
func HandlerSetVolumeCompression(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]
	compressionType := vars["compression"]

	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	zvol, err := zfs.GetDataset(lu.GetZvol())
	err = zvol.SetProp("compression", compressionType)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	zvol.RefreshProps()
	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(zvol))
	err = json.NewEncoder(w).Encode(lu)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

// Export Volume
func HandlerExportVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]

	var reqJson ExportRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	view, err := lu.AddView(reqJson.Hostgroup, reqJson.Targetgroup, reqJson.Lun)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	log.Printf("===\n%s: %s\n\n", traceFunctionName(), json.NewEncoder(os.Stdout).Encode(view))
	err = json.NewEncoder(w).Encode(view)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

// UnExport Volume
func HandlerUnExportVolume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]

	var reqJson ExportRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&reqJson)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	if lu.ViewEntryCount > 0 {
		viewNumber, err := lu.GetViewEntry(reqJson.Hostgroup, reqJson.Targetgroup)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
		if err = lu.RemoveView(viewNumber.ViewEntry); err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}

		json.NewEncoder(w).Encode(lu)
		return

	} else {
		sendMessage(w, http.StatusNotFound, traceFunctionName(), err.Error())
		return
	}
}

func HandlerGetVolumeExports(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	volumeName := vars["volume"]

	basepath := poolName + "/" + domainName + "/" + projectName

	lu, err := stmf.GetLu(volumeName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if !IsVolumeBelongsToProject(basepath, *lu) {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "Volume not found in specified project")
		return
	}

	if lu.ViewEntryCount == 0 {
		json.NewEncoder(w).Encode([]string{})
		return
	}

	views, err := lu.ListView()
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	json.NewEncoder(w).Encode(views)
}
