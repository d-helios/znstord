package stmf

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"log"
)

// wrapper for cmdStmfadm
func cmdStmfadm(arg ...string) ([][]string, error) {
	c := command{Command: "stmfadm"}
	return c.Run(":", arg...)
}

// wrapper for stmfha
func cmdStmfha(arg ...string) ([][]string, error) {
	c := command{Command: "/opt/HAC/RSF-1/bin/stmfha"}
	return c.Run(":", arg...)
}

// ListLUs - stmfadm list-logicalunit [logicalunit]
func ListLUs(luName string) ([]*LogicalUnit, error) {
	args := []string{"list-lu", "-v"}

	if luName != "" {
		args = append(args, luName)
	}

	out, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	var luArray []*LogicalUnit
	var lu *LogicalUnit

	// if operation complete successfuly, but no volumes found
	// return empty array
	if len(out) == 0 {
		return luArray, nil
	}

	// stmfadm list-logicalunit output at least 16 lines
	if len(out) < 17 {
		return nil, &Error{
			Debug:  "stmfadm " + strings.Join(args, " "),
			Stderr: fmt.Sprintf("%q\n", out),
		}
	}

	// fill luArray
	for i := 0; i < len(out); i++ {
		if out[i][0] == "LU Name" {
			lu, i = &LogicalUnit{LUName: out[i][1]}, i+1
			lu.OperationalStatus, i = out[i][1], i+1
			lu.ProviderName, i = out[i][1], i+1
			lu.Alias, i = out[i][1], i+1
			lu.ViewEntryCount, err = strconv.ParseUint(out[i][1], 10, 16)

			if err != nil {
				return nil, err
			}
			i = i + 1

			lu.DataFile, i = out[i][1], i+1
			lu.MetaFile, i = out[i][1], i+1
			lu.Size, err = strconv.ParseUint(out[i][1], 10, 64)

			if err != nil {
				return nil, err
			}
			i = i + 1

			lu.BlockSize, err = strconv.ParseUint(out[i][1], 10, 16)

			if err != nil {
				return nil, err
			}
			i = i + 1

			lu.ManagementURL, i = out[i][1], i+1
			lu.VendorID, i = out[i][1], i+1
			lu.ProductID, i = out[i][1], i+1
			lu.SerialNum, i = out[i][1], i+1
			lu.WriteProtect, i = out[i][1], i+1
			lu.WriteCacheModeSelect, i = out[i][1], i+1
			lu.WritebackCache, i = out[i][1], i+1
			lu.AccessState = out[i][1]
			luArray = append(luArray, lu)
		}
	}
	return luArray, nil
}

// GetLu - get logical unit. stmfadm list-lu wwid
func GetLu(luName string) (*LogicalUnit, error) {
	LUs, err := ListLUs(luName)

	if err != nil {
		return nil, err
	}
	return LUs[0], nil
}

// GetLuByZvol - get logical unit by zvol.
func GetLuByZvol(zvol string) (*LogicalUnit, error) {
	LUs, err := ListLUs("")
	if err != nil {
		return nil, err
	}

	for _, lu := range LUs {
		if lu.DataFile == RDSK_DEFAULT_PREFIX+zvol {
			return lu, nil
		}
	}

	return nil, &Error{
		Err:    err,
		Debug:  fmt.Sprintf("LogicalUnit not found. ZVOL: %s. LUs: %q", zvol, LUs),
		Stderr: "",
	}
}

// GetLuByAlias - get logical unit by alias
// Returns the first match.
func GetLuByAlias(alias string) (*LogicalUnit, error) {
	LUs, err := ListLUs("")
	if err != nil {
		return nil, err
	}

	for _, lu := range LUs {
		if lu.Alias == alias {
			return lu, nil
		}
	}

	return nil, &Error{
		Err:    err,
		Debug:  fmt.Sprintf("LogicalUnit not found. Alias: %s", alias),
		Stderr: "",
	}
}

// CreateLu - create logical unit. zvol is only supported backend
func CreateLu(zvol, opts string) (*LogicalUnit, error) {
	// default options is -p blk=4096
	args := []string{"create-lu", "-p", "blk=" + defaultVolBlockSize}

	if opts != "" {
		optionLists := strings.Split(opts, " ")
		args = append(args, optionLists...)
	}

	args = append(args, RDSK_DEFAULT_PREFIX+zvol)

	output, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	lu, err := GetLu(output[0][1])

	if err != nil {
		return nil, err
	}

	if StmfHAEnabled {
		if err := BackupSTMFConfiguration(strings.Split(zvol, "/")[0]); err != nil {
			log.Printf("Can't backup pool configuration. Err: %s", err.Error())
		}
	}

	return lu, nil
}

// Offline - set logical unit Offline
func (logicalunit *LogicalUnit) Offline() error {
	args := []string{"offline-lu"}

	args = append(args, logicalunit.LUName)

	_, err := cmdStmfadm(args...)
	if err != nil {
		return err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return err
	}
	*logicalunit = *tmpLu
	return nil
}

// Online - set logical unit Online
func (logicalunit *LogicalUnit) Online() error {
	args := []string{"online-lu"}

	args = append(args, logicalunit.LUName)

	_, err := cmdStmfadm(args...)
	if err != nil {
		return err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return err
	}
	*logicalunit = *tmpLu

	return nil
}

// Modify logical unit stmf properties or zvol properties
func (logicalunit *LogicalUnit) Modify(opts string) error {
	args := []string{"modify-lu"}

	if opts == "" {
		// nothing to modify. exit
		return nil
	}

	for _, parameter := range strings.Split(opts, " ") {
		args = append(args, parameter)
	}

	args = append(args, logicalunit.LUName)
	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return err
	}
	*logicalunit = *tmpLu

	return nil
}

// Delete - delete logical unit.
// keepViews flag is nessessory to save views while modifine volume
// options, for example to resize volume.
// Ex:
//  * 	Delete lu with option keepViews
//	*	Resize zvol
//	*	Create lu with same options (my be using stmfadm import ???)
func (logicalunit *LogicalUnit) Delete(keepViews bool) error {
	args := []string{"delete-lu"}

	if keepViews {
		args = append(args, "-k")
	}

	args = append(args, logicalunit.LUName)

	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	logicalunit = nil

	return nil
}

// AddView - export volume.
func (logicalunit *LogicalUnit) AddView(hostGroup, targetGroup string, lun int64) (*View, error) {
	args := []string{"add-view"}

	if hostGroup != "" {
		args = append(args, "-h")
		args = append(args, hostGroup)
	}

	if targetGroup != "" {
		args = append(args, "-t")
		args = append(args, targetGroup)
	}

	if lun >= 0 {
		args = append(args, "-n")
		args = append(args, strconv.Itoa(int(lun)))
	}

	args = append(args, logicalunit.LUName)

	_, err := cmdStmfadm(args...)
	if err != nil {
		return nil, err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return nil, err
	}
	*logicalunit = *tmpLu

	return logicalunit.GetViewEntry(hostGroup, targetGroup)
}

// GetViewEntry - get volume export entry
func (logicalunit *LogicalUnit) GetViewEntry(hostgroup, targetgroup string) (*View, error) {
	Views, err := logicalunit.ListView()
	if err != nil {
		return nil, err
	}

	// If hostgroup not specified replace by All
	if hostgroup == "" {
		hostgroup = "All"
	}

	if targetgroup == "" {
		targetgroup = "All"
	}

	// if targetgroup not specified replace by All
	for _, view := range Views {
		if view.HostGroup == hostgroup &&
			view.TargetGroup == targetgroup {
			return &view, nil
		}
	}
	return nil, &Error{
		Err:    errors.New("View Entry not found"),
		Debug:  fmt.Sprintf("LU %s exported to %q", logicalunit.LUName, Views),
		Stderr: "",
	}
}

// RemoveView - unexport lu.
func (logicalunit *LogicalUnit) RemoveView(viewEntryNumber uint64) error {
	args := []string{"remove-view", "-l"}

	args = append(args, logicalunit.LUName)

	args = append(args, strconv.FormatUint(viewEntryNumber, 10))

	_, err := cmdStmfadm(args...)
	if err != nil {
		return err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return err
	}
	*logicalunit = *tmpLu

	return nil
}

// RemoveAllView - unexport lu from all hosts
func (logicalunit *LogicalUnit) RemoveAllView() error {
	args := []string{"remove-view", "-a", "-l"}

	args = append(args, logicalunit.LUName)

	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	tmpLu, err := GetLu(logicalunit.LUName)
	if err != nil {
		return err
	}
	*logicalunit = *tmpLu

	return nil
}

// ListView - lu export list
func (logicalunit *LogicalUnit) ListView() ([]View, error) {
	args := []string{"list-view", "-l"}
	args = append(args, logicalunit.LUName)

	out, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	var Views []View

	var view *View

	var viewNum uint64
	var lunNum uint64

	for i := 0; i < len(out); i++ {
		if out[i][0] == "View Entry" {
			viewNum, err = strconv.ParseUint(out[i][1], 10, 16)
			if err != nil {
				return nil, err
			}

			i = i + 1

			view = &View{ViewEntry: viewNum}
			view.HostGroup, i = out[i][1], i+1
			view.TargetGroup, i = out[i][1], i+1
			lunNum, err = strconv.ParseUint(out[i][1], 10, 16)

			if err != nil {
				return nil, &Error{
					Err:    err,
					Debug:  "Can't convert LUNum to Uint, in function ListView",
					Stderr: "",
				}
			}
			view.LUN = lunNum
			Views = append(Views, *view)
		}
	}

	return Views, nil
}

// ListHostGroup - list host groups, stmfadm list-hg
func ListHostGroup(hg string) ([]*HostGroup, error) {
	args := []string{"list-hg", "-v"}

	if hg != "" {
		args = append(args, hg)
	}

	out, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	var hostGroup *HostGroup
	var hostGroups []*HostGroup

	for i := 0; i < len(out); i++ {
		if out[i][0] == "Host Group" {
			hostGroup = &HostGroup{HostGroup: strings.Trim(out[i][1], " ")}
			hostGroups = append(hostGroups, hostGroup)
		}

		if out[i][0] == "\tMember" {
			hostGroup.Members = append(hostGroup.Members, strings.Trim(strings.Join(out[i][1:], ":"), " "))
		}
	}

	return hostGroups, err
}

// GetHostGroup - get specified hostgroup
func GetHostGroup(hg string) (*HostGroup, error) {
	HGs, err := ListHostGroup(hg)

	if err != nil {
		return nil, err
	}

	return HGs[0], nil
}

// CreateHostGroup - create host group. stmfadm create-hg <hostgroup>
func CreateHostGroup(hg string) (*HostGroup, error) {
	args := []string{"create-hg", hg}

	if hg == "All" {
		return nil, &Error{
			Err:    errors.New("TargetName - All, not allowed"),
			Debug:  fmt.Sprintf("TargetGroup: %s", hg),
			Stderr: "",
		}
	}

	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return nil, err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return nil, err
		}
	}

	hostGroup, hgErr := GetHostGroup(hg)

	if hgErr != nil {
		return nil, hgErr
	}

	return hostGroup, nil
}

// AddMember - add member to hostgroup.
func (hg *HostGroup) AddMember(member string) error {
	args := []string{"add-hg-member", "-g", hg.HostGroup, member}

	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	}

	tmpHg, err := GetHostGroup(hg.HostGroup)
	if err != nil {
		return err
	}
	*hg = *tmpHg

	return nil
}

// AddMultiHostGroupMember - add initiator to one or more hostgroups.
// This feature not available in OpenSolaris.
func (hg *HostGroup) AddMultiHostGroupMember(member string) error {
	args := []string{"add-hg-member", "-g", hg.HostGroup, "-F", member}

	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	}

	tmpHg, err := GetHostGroup(hg.HostGroup)
	if err != nil {
		return err
	}
	*hg = *tmpHg

	return nil
}

// RemoveMember - remove member from hostgroup
func (hg *HostGroup) RemoveMember(member string) error {
	args := []string{"remove-hg-member", "-g", hg.HostGroup, member}

	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	}

	tmpHg, err := GetHostGroup(hg.HostGroup)
	if err != nil {
		return err
	}
	*hg = *tmpHg

	return nil
}

// Delete Host Group Member
func (hg *HostGroup) Delete() error {
	args := []string{"delete-hg", hg.HostGroup}

	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return err
		}
	}

	*hg = HostGroup{}

	return nil
}

// CreateTargetGroup - create target group.
func CreateTargetGroup(tg string) (*TargetGroup, error) {
	args := []string{"create-tg", tg}

	if tg == "All" {
		return nil, &Error{
			Err:    errors.New("TargetName - All, not allowed"),
			Debug:  fmt.Sprintf("TargetGroup: %s", tg),
			Stderr: "",
		}
	}
	if StmfHAEnabled {
		if _, err := cmdStmfha(args...); err != nil {
			return nil, err
		}
	} else  {
		if _, err := cmdStmfha(args...); err != nil {
			return nil, err
		}
	}

	return &TargetGroup{TargetGroup: tg}, nil
}

// AddMember - add member to target group
func (tg *TargetGroup) AddMember(tpg string) error {
	args := []string{"add-tg-member", "-g", tg.TargetGroup, tpg}
	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	return nil
}

// RemoveMember - remove target group member
func (tg *TargetGroup) RemoveMember(tpg string) error {
	args := []string{"remove-tg-member", "-g", tg.TargetGroup, tpg}
	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	return nil
}

// Delete - delete target group member
func (tg *TargetGroup) Delete() error {
	args := []string{"delete-tg", tg.TargetGroup}

	_, err := cmdStmfadm(args...)

	if err != nil {
		return err
	}

	*tg = TargetGroup{}

	return nil
}

// ListTargetGroups - list target groups
func ListTargetGroups() ([]*TargetGroup, error) {
	args := []string{"list-tg"}

	out, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	var tgList []*TargetGroup

	for _, target := range out {
		tgList = append(tgList, &TargetGroup{TargetGroup: target[1]})
	}

	return tgList, nil
}

// GetTargetGroup - get specified target groups
func GetTargetGroup(tg string) (*TargetGroup, error) {
	args := []string{"list-tg", tg}

	out, err := cmdStmfadm(args...)

	if err != nil {
		return nil, err
	}

	var tgList []*TargetGroup

	for _, target := range out {
		tgList = append(tgList, &TargetGroup{TargetGroup: target[1]})
	}

	if len(tgList) > 0 {
		return tgList[0], nil
	}

	return nil, nil
}

// BackupSTMFConfiguration - Backup Configuration (RSF-1 Cluster)
func BackupSTMFConfiguration(pool string) error {
	args := []string{"backup", pool}
	_, err := cmdStmfha(args...)

	if err != nil {
		return err
	}
	return nil
}
