package znstor

import (
	"encoding/json"
	"github.com/d-helios/stmf"
	"github.com/gorilla/mux"
	"net/http"
)

func HandlerCreateHostGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]

	hostgroup, err := stmf.CreateHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(hostgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
	}
}

func HandlerGetHostGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]

	hostgroup, err := stmf.GetHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(hostgroup)

	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

}

func HandlerGetHostGroupList(w http.ResponseWriter, r *http.Request) {
	hostgroup, err := stmf.ListHostGroup("")
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
	if hostgroup == nil {
		// return empty array if pointer is nil
		err = json.NewEncoder(w).Encode([]string{})
	} else {
		err = json.NewEncoder(w).Encode(hostgroup)
	}

	// send error Econding error
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerAddHostGroupMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]
	memberName := vars["member"]

	hostgroup, err := stmf.GetHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = hostgroup.AddMember(memberName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(hostgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerRemoveHostGroupMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]
	memberName := vars["member"]

	hostgroup, err := stmf.GetHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = hostgroup.RemoveMember(memberName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(hostgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerDeleteHostGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]

	hostgroup, err := stmf.GetHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = hostgroup.Delete()
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}

func HandlerAddMultiGroupMember(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostgroupName := vars["hostgroup"]
	memberName := vars["member"]

	hostgroup, err := stmf.GetHostGroup(hostgroupName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = hostgroup.AddMultiHostGroupMember(memberName)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}

	err = json.NewEncoder(w).Encode(hostgroup)
	if err != nil {
		sendMessage(w, http.StatusBadRequest, traceFunctionName(), err.Error())
		return
	}
}
