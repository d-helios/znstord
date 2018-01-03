package znstor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"github.com/d-helios/znstord/stmf"
)

func IsValueInArray(value interface{}, list []interface{}) bool {
	for _, curval := range list {
		if value == curval {
			return true
		}
	}
	return false
}

func checkOption(prop string, validOpts []string) bool {
	for _, b := range validOpts {
		if b == prop {
			return true
		}
	}
	return false
}

func checkRequiredParameters(dictToCheck map[string]interface{}, requiredParameters []string) bool {
	numberOfMatches := 0
	targetNumberOfMatches := len(requiredParameters)

	for _, key := range requiredParameters {
		if _, ok := dictToCheck[key]; ok {
			numberOfMatches++
		}
	}

	if numberOfMatches == targetNumberOfMatches {
		return true
	} else {
		return false
	}
}

func traceFunctionName() string {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	return fmt.Sprintf("%s:%d", f.Name(), line)
}

func IsVolumeBelongsToProject(basepath string, lu stmf.LogicalUnit) bool {
	volLength := len(strings.Split(lu.GetZvol(), "/"))
	vol_basepath := strings.Join(
		strings.Split(lu.GetZvol(), "/")[0:volLength-1],
		"/")

	if basepath != vol_basepath {
		log.Printf("\n===\nVolume %s, not belongs to %s project\n", lu.GetZvol(), basepath)
		log.Printf("lu: %q\n", lu)
		return false
	}

	return true
}

func sendMessage(w http.ResponseWriter, httpStatusCode uint64, subject, message string) {
	msg := RespMsg{
		Subject: subject,
		Msg:     message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(httpStatusCode))

	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Println("ERROR:", err)
	} else {
		log.Printf("sendMessage: %q", msg)
	}
}
