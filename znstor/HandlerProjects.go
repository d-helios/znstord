package znstor

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/d-helios/znstord/zfs"
	"github.com/gorilla/mux"
	"github.com/jinzhu/copier"
	"io"
)

// List available projects within domain
func HandlerGetProjectList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	basepath := poolName + "/" + domainName

	datasets, err := zfs.ListDatasets(zfs.Filesystem, basepath, true, 1)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	var projects = make([]Project, 0)

	for _, dataset := range datasets {
		var project Project

		// get last dataset path
		project.Dataset = strings.Join(
			strings.Split(dataset.Dataset, "/")[len(strings.Split(dataset.Dataset, "/"))-1:len(strings.Split(dataset.Dataset, "/"))],
			"")

		// exclude domain dataset
		if basepath != dataset.Dataset {
			projects = append(projects, project)
		}
	}

	err = json.NewEncoder(w).Encode(projects)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	return
}

// Get specified project
func HandlerGetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	dataset, err := zfs.GetDataset(basepath)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	var project Project

	project.Dataset = strings.Join(
		strings.Split(dataset.Dataset, "/")[len(strings.Split(dataset.Dataset, "/"))-1:len(strings.Split(dataset.Dataset, "/"))],
		"")

	dataset.RefreshProps()
	copier.Copy(&project.Options, dataset.Props.(*zfs.FsDataset))

	err = json.NewEncoder(w).Encode(project)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	return
}

// Get specified project
func HandlerProjectExists(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	err := zfs.IsDatasetExist(basepath)
	if err != nil {
		sendMessage(w, http.StatusNoContent, traceFunctionName(), "")
	} else {
		sendMessage(w, http.StatusOK, traceFunctionName(), "")
	}
}

// Create project
func HandlerCreateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	var zfsOptions FilesystemRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&zfsOptions)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	// Check if quota is specified
	if zfsOptions.Quota == 0 {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), "")
		return
	}

	args := []string{}
	if zfsOptions.Alias != "" {
		args = append(args, "-o custom:alias="+zfsOptions.Alias)
	}
	if zfsOptions.Reservation != 0 {
		args = append(args, "-o reservation="+strconv.FormatUint(zfsOptions.Reservation, 10))
	}
	if zfsOptions.Dedup != "" {
		args = append(args, "-o dedup="+zfsOptions.Dedup)
	}
	if zfsOptions.Compression != "" {
		args = append(args, "-o compression="+zfsOptions.Compression)
	}
	if zfsOptions.Atime != "" {
		args = append(args, "-o atime="+zfsOptions.Atime)
	}
	if zfsOptions.Refquota != 0 {
		args = append(args, "-o refquota="+strconv.FormatUint(zfsOptions.Refquota, 10))
	}
	if zfsOptions.Refreservation != 0 {
		args = append(args, "-o refreservation="+strconv.FormatUint(zfsOptions.Refreservation, 10))
	}

	dataset, err := zfs.CreateFilesystem(basepath, strings.Join(args, " "), zfsOptions.Quota)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	var project Project

	project.Dataset = strings.Join(
		strings.Split(dataset.Dataset, "/")[len(strings.Split(dataset.Dataset, "/"))-1:len(strings.Split(dataset.Dataset, "/"))],
		"")

	dataset.RefreshProps()
	copier.Copy(&project.Options, dataset.Props.(*zfs.FsDataset))

	err = json.NewEncoder(w).Encode(project)

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerModifyProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	var zfsOptions FilesystemRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&zfsOptions)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	dataset, err := zfs.GetDataset(basepath)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if zfsOptions.Alias != "" {
		err := dataset.SetProp("custom:alias", zfsOptions.Alias)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}

	if zfsOptions.Quota != 0 {
		err := dataset.SetProp("quota", zfsOptions.Quota)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Reservation != 0 {
		err := dataset.SetProp("reservation", zfsOptions.Reservation)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Dedup != "" {
		err := dataset.SetProp("dedup", zfsOptions.Dedup)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Compression != "" {
		err := dataset.SetProp("compression", zfsOptions.Compression)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Atime != "" {
		err := dataset.SetProp("atime", zfsOptions.Atime)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Refquota != 0 {
		err := dataset.SetProp("refquota", zfsOptions.Refquota)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}
	if zfsOptions.Refreservation != 0 {
		err := dataset.SetProp("refreservation", zfsOptions.Refreservation)
		if err != nil {
			sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
			return
		}
	}

	var project Project

	project.Dataset = strings.Join(
		strings.Split(dataset.Dataset, "/")[len(strings.Split(dataset.Dataset, "/"))-1:len(strings.Split(dataset.Dataset, "/"))],
		"")

	dataset.RefreshProps()
	copier.Copy(&project.Options, dataset.Props.(*zfs.FsDataset))

	err = json.NewEncoder(w).Encode(project)

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerDestroyProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	dataset, err := zfs.GetDataset(basepath)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = dataset.Destroy("")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	sendMessage(w, http.StatusOK, traceFunctionName(), "")
}

func HandlerForceDestroyProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	domainName := vars["domain"]
	poolName := vars["pool"]
	projectName := vars["project"]
	basepath := poolName + "/" + domainName + "/" + projectName

	dataset, err := zfs.GetDataset(basepath)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = dataset.Destroy("-R")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	sendMessage(w, http.StatusOK, traceFunctionName(), "")
}
