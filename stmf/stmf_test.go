package stmf_test

import (
	"errors"
	"fmt"
	"github.com/d-helios/znstord/stmf"
	"github.com/d-helios/znstord/zfs"
	"testing"
)

const (
	mb_size = 1048576

	volumePool   = "rpool/vol"
	zvol1        = "vol_1"
	snapshotName = "snap1"
	zvolClone1   = "vol_clone_1"

	hostGroup1   = "host_group1"
	targetGroup1 = "target_group1"
)

func TestCreateVolumePool(t *testing.T) {
	_, err := zfs.CreateFilesystem(volumePool, "-o custom:alias=volumePool", 500*mb_size)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestCreateZvol1(t *testing.T) {
	_, err := zfs.CreateVolume(volumePool+"/"+zvol1, "-o compression=lz4", true, 256*mb_size)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestCreateLu(t *testing.T) {
	_, err := stmf.CreateLu(volumePool+"/"+zvol1, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLogicalUnit_Modify(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = lu.Modify("-p alias=MyVol")
	if err != nil {
		t.Fatalf(err.Error())
	}

	if lu.Alias != "MyVol" {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. can't change Logical Volume alias"),
			Debug:  fmt.Sprintf("LogicalUnit: %q", lu),
			Stderr: ""}.Error())
	}
}

func TestLogicalUnit_Offline(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = lu.Offline()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if lu.OperationalStatus != "Offline" {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. can't change Logical Volume Operational Status"),
			Debug:  fmt.Sprintf("LogicalUnit: %q.", lu),
			Stderr: ""}.Error())
	}
}

func TestLogicalUnit_Online(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf(err.Error())
	}
	err = lu.Online()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if lu.OperationalStatus != "Online" {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. can't change Logical Volume Operational Status"),
			Debug:  fmt.Sprintf("LogicalUnit: %q.", lu),
			Stderr: ""}.Error())
	}
}

func TestCreateHostGroup(t *testing.T) {
	hg, err := stmf.CreateHostGroup(hostGroup1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if hg.HostGroup != hostGroup1 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. HostGroup name does not match provided while creation."),
			Debug:  fmt.Sprintf("HostGroup: %q.", hg),
			Stderr: ""}.Error())
	}
}

func TestHostGroup_AddMember(t *testing.T) {
	hg, err := stmf.GetHostGroup(hostGroup1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, host := range []string{"iqn.1986-03.com.sun:aa:e20000000000.59235955",
		"iqn.1986-03.com.sun:aa:e20000000000.59235956",
		"iqn.1986-03.com.sun:aa:e20000000000.59235957"} {
		err := hg.AddMember(host)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
	if len(hg.Members) != 3 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. HostGroup member count mismatch."),
			Debug:  fmt.Sprintf("HostGroup: %q.", hg),
			Stderr: ""}.Error())
	}
}

func TestHostGroup_RemoveMember(t *testing.T) {
	hg, err := stmf.GetHostGroup(hostGroup1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	member_1 := hg.Members[0]
	err = hg.RemoveMember(member_1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(hg.Members) != 2 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. HostGroup member count mismatch."),
			Debug:  fmt.Sprintf("HostGroup: %q.", hg),
			Stderr: ""}.Error())
	}
}

func TestCreateTargetGroup(t *testing.T) {
	tg, err := stmf.CreateTargetGroup(targetGroup1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if tg.TargetGroup != targetGroup1 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. TargetGroup name does not match provided while creation."),
			Debug:  fmt.Sprintf("HostGroup: %q.", tg),
			Stderr: ""}.Error())
	}
}

func TestListTargetGroups(t *testing.T) {
	tgs, err := stmf.ListTargetGroups()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(tgs) == 0 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("STMF. TargetGroup count mismatch."),
			Debug:  fmt.Sprintf("HostGroup: %q.", tgs),
			Stderr: ""}.Error())
	}
}

func TestLogicalUnit_AddView(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	l_view, err := lu.AddView(hostGroup1, targetGroup1, -1)
	if err != nil {
		t.Fatal(err.Error())
	}

	if l_view.TargetGroup != targetGroup1 || l_view.HostGroup != hostGroup1 {
		t.Fatalf(stmf.Error{
			Err:    errors.New("View mistmatch."),
			Debug:  fmt.Sprintf("Views: %q", l_view),
			Stderr: "",
		}.Error())
	}
}

func TestLogicalUnit_RemoveView(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	view, err := lu.GetViewEntry(hostGroup1, targetGroup1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = lu.RemoveView(view.ViewEntry)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestHostGroup_Delete(t *testing.T) {
	hg, err := stmf.GetHostGroup(hostGroup1)
	if err != nil {
		t.Fatalf("Couldn't get hostgroup: %s. Err: %s", hostGroup1, err.Error())
	}
	err = hg.Delete()
	if err != nil {
		t.Fatalf("Couldn't delete hostgroup: %s. Err: %s", hostGroup1, err.Error())
	}

	if hg.HostGroup != "" {
		t.Fatalf("HostGroup pointer is not nil. Pointer: %q", hg)
	}
}

func TestTargetGroup_Delete(t *testing.T) {
	tg, err := stmf.GetTargetGroup(targetGroup1)
	if err != nil {
		t.Fatalf("Couldn't get targetgroup: %s. Err: %s", targetGroup1, err.Error())
	}
	err = tg.Delete()
	if err != nil {
		t.Fatalf("Couldn't delete targetgroup: %s. Err: %s", targetGroup1, err.Error())
	}

	if tg.TargetGroup != "" {
		t.Fatalf("TargetGroup pointer is not nil. Pointer: %q", tg)
	}
}

func TestLogicalUnit_Delete(t *testing.T) {
	lu, err := stmf.GetLuByZvol(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf("Couldn't find logical unit by zvol: %s. Err: %s", volumePool+"/"+zvol1, err.Error())
	}

	err = lu.Delete(false)
	if err != nil {
		t.Fatalf("Couldn't delete LU %s. Err: %s", lu.LUName, err.Error())
	}
}

func TestDestroyZvol_1(t *testing.T) {
	dataset, err := zfs.GetDataset(volumePool + "/" + zvol1)
	if err != nil {
		t.Fatalf("Couldn't get volume %s. Err: %s", volumePool+"/"+zvol1, err.Error())
	}

	err = dataset.Destroy("")
	if err != nil {
		t.Fatalf("Couldn't destroy Volume %s. Err: %s", volumePool+"/"+zvol1, err.Error())
	}
}

func TestDestroyVolumePool(t *testing.T) {
	dataset, err := zfs.GetDataset(volumePool)
	if err != nil {
		t.Fatalf("Couldn't get volume pool %s. Err: %s ", volumePool, err.Error())
	}

	err = dataset.Destroy("")
	if err != nil {
		t.Fatalf("Couldn't destroy volume pool %s. Err: %s", volumePool, err.Error())
	}
}
