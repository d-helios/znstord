package znstor

import (
	"encoding/json"
	"github.com/d-helios/itadm"
	"github.com/d-helios/stmf"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

/*
=======
Target Groups
=======
*/
func HandlerGetTargetGroupList(w http.ResponseWriter, r *http.Request) {
	targetGroups, err := stmf.ListTargetGroups()
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if targetGroups == nil {
		err = json.NewEncoder(w).Encode([]string{})
	} else {
		err = json.NewEncoder(w).Encode(targetGroups)
	}
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetTargetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetgroupName := vars["hostgroup"]

	targetGroup, err := stmf.GetTargetGroup(targetgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(targetGroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerAddTargetGroupMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetgroupName := vars["targetgroup"]
	memberName := vars["member"]

	targetgroup, err := stmf.GetTargetGroup(targetgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = targetgroup.AddMember(memberName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(targetgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerRemoveTargetGroupMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetgroupName := vars["targetgroup"]
	memberName := vars["member"]

	targetgroup, err := stmf.GetTargetGroup(targetgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	err = targetgroup.RemoveMember(memberName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(targetgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerDeleteTargetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetgroupName := vars["targetgroup"]

	targetgroup, err := stmf.GetTargetGroup(targetgroupName)

	err = targetgroup.Delete()
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerCreateTargetGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetgroupName := vars["targetgroup"]

	targetgroup, err := stmf.CreateTargetGroup(targetgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(targetgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

/*
=======
TargetPortGroups
=======
*/
func HandlerCreateTargetPortGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// tpg == target port group
	tpgName := vars["tpg"]

	var tpgOptions TpgCreateRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&tpgOptions)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	tpg, err := itadm.CreateTargetPortGroup(tpgName, tpgOptions.portals)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(tpg)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetTargetPortGroupList(w http.ResponseWriter, r *http.Request) {
	tpg, err := itadm.ListTargetPortGroups("")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if tpg == nil {
		// return empty string
		err = json.NewEncoder(w).Encode([]string{})
	} else {
		err = json.NewEncoder(w).Encode(tpg)
	}

	// check json.NewEncoder error
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetTargetPortGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpgName := vars["tpg"]
	tpg, err := itadm.GetTargetPortGroup(tpgName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(tpg)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerDeleteTargetPortGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpgName := vars["tpg"]
	tpg, err := itadm.GetTargetPortGroup(tpgName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = tpg.Delete(false)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerForceDeleteTargetPortGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tpgName := vars["tpg"]
	tpg, err := itadm.GetTargetPortGroup(tpgName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = tpg.Delete(true)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

/*
=======
Targets
=======
*/
func HandlerCreateTarget(w http.ResponseWriter, r *http.Request) {
	var targetOptions TargetCreateRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, requestPayloadMaxSize))
	err := decoder.Decode(&targetOptions)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	target, err := itadm.CreateTarget(
		targetOptions.Iqn,
		targetOptions.Alias,
		targetOptions.Tpg,
	)

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(target)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetTargetList(w http.ResponseWriter, r *http.Request) {
	targets, err := itadm.ListTargets("")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	if targets == nil {
		// return empty array
		err = json.NewEncoder(w).Encode([]string{})
	} else {
		err = json.NewEncoder(w).Encode(targets)
	}

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetName := vars["target"]

	target, err := itadm.GetTarget(targetName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(target)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerDeleteTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetName := vars["target"]

	target, err := itadm.GetTarget(targetName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = target.Delete(false)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerForceDeleteTarget(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	targetName := vars["target"]

	target, err := itadm.GetTarget(targetName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = target.Delete(true)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}
