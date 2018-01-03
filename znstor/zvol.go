package znstor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"github.com/d-helios/znstord/stmf"
	"github.com/d-helios/znstord/zfs"
	"github.com/twinj/uuid"
	"time"
)

func CreateVolume(basepath string, zvol ZVolCreateRequest) (*stmf.LogicalUnit, error) {
	// check if alias is not empty
	if zvol.Alias == "" {
		return nil, &Error{
			Err:    errors.New("Alias not specified"),
			Debug:  fmt.Sprintf("request: %q", zvol),
			Stderr: "",
		}
	}

	// check if size is specified
	if zvol.VolSize == 0 {
		return nil, &Error{
			Err:    errors.New("VolSize not specified"),
			Debug:  fmt.Sprintf("request: %q", zvol),
			Stderr: "",
		}
	}

	if zvol.Serial == "" {
		zvol.Serial = uuid.NewV4().String()
	}

	// Append options
	zvolArgs := []string{}

	if zvol.Options.VolBlockSize == 0 {
		zvolArgs = append(zvolArgs, "-o volblocksize=8192")
	} else {
		zvolArgs = append(zvolArgs, "-o volblocksize="+strconv.FormatUint(
			zvol.Options.VolBlockSize,
			10,
		))
	}

	if zvol.Options.Compression != "" {
		zvolArgs = append(zvolArgs, "-o compression="+zvol.Options.Compression)
	}

	if zvol.Options.Dedup != "" {
		zvolArgs = append(zvolArgs, "-o dedup="+zvol.Options.Dedup)
	}

	if zvol.Options.Reservation != 0 {
		zvolArgs = append(zvolArgs, "-o reservation="+
			strconv.FormatUint(zvol.Options.Reservation, 10))
	}

	zvolArgs = append(zvolArgs, "-o custom:sflag="+sflagManaged)

	volName := zvol.Alias

	// create volume
	zfsVolume, err := zfs.CreateVolume(
		basepath+"/"+volName,
		strings.Join(zvolArgs, " "),
		zvol.Options.Thin,
		zvol.VolSize)

	if err != nil {
		return nil, err
	}

	stmfArgs := []string{"-p", "alias=" + zvol.Alias, "-p", "serial=" + zvol.Serial}

	stmfLu, err := stmf.CreateLu(
		zfsVolume.Dataset,
		strings.Join(stmfArgs, " "))

	if err != nil {
		return nil, err
	}

	return stmfLu, nil
}

func VolResize(lu_uuid string, newSize uint64) error {
	// check if new size is greater then old size
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return err
	}

	if lu.Size > newSize {
		return Error{
			Err:    errors.New(fmt.Sprintf("Can't decrise volume size.")),
			Debug:  "",
			Stderr: "",
		}
	}

	zfsVolume, err := zfs.GetDataset(lu.GetZvol())
	if err != nil {
		return err
	}

	if err := zfsVolume.SetProp("volsize", newSize); err != nil {
		return err
	}

	if err := zfsVolume.RefreshProps(); err != nil {
		return err
	}

	// change meta information for stmf lu
	lu.Modify(fmt.Sprintf("-s %d", zfsVolume.Props.(*zfs.VolDataset).Volsize))
	if err != nil {
		return err
	}

	return nil
}

func VolDestroy(lu_uuid string) error {

	lu, err := stmf.GetLu(lu_uuid)
	vol_mutex.Lock()
	err = lu.Delete(false)
	if err != nil {
		vol_mutex.Unlock()
		return err
	}
	// sleep 50ms to ensure we not catch setuation like this:
	// stmfadm[3928]: [ID 155448 user.error] transaction commit for provider_data_pg_sbd failed - object already exists
	time.Sleep(100)
	vol_mutex.Unlock()

	zfsVolume, err := zfs.GetDataset(lu.GetZvol())
	if err != nil {
		return err
	}

	return zfsVolume.Destroy("")
}

func VolSnapshot(lu_uuid, snapname string) (*zfs.Dataset, error) {
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return nil, err
	}

	zfsVolume, err := zfs.GetDataset(lu.GetZvol())
	if err != nil {
		return nil, err
	}
	return zfsVolume.Snapshot(snapname)
}

func ListSnapshot(lu_uuid string) ([]*zfs.Dataset, error) {
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return nil, err
	}
	return zfs.ListDatasets(zfs.Snapshot, lu.GetZvol(), true, 1)
}

func GetSnapshot(lu_uuid, snapshotName string) (*zfs.Dataset, error) {
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return nil, err
	}
	return zfs.GetDataset(lu.GetZvol() + "@" + snapshotName)
}

func VolRollback(lu_uuid, snapshotName string) error {
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return err
	}
	zfsVolume, err := zfs.GetDataset(lu.GetZvol())
	if err != nil {
		return err
	}

	saved_alias := lu.Alias

	// delete stmf lu with keepViews option
	err = lu.Delete(true)
	if err != nil {
		return err
	}

	err = zfsVolume.Rollback(zfsVolume.Dataset + "@" + snapshotName)
	if err != nil {
		return err
	}

	// create logical unit with saved guid option and alias
	_, err = stmf.CreateLu(zfsVolume.Dataset, "-p guid="+lu_uuid+" -p alias="+saved_alias)
	if err != nil {
		return err
	}

	return nil
}

func VolCloneFromSnapshot(lu_uuid, snapname, cloneAlias string) (*stmf.LogicalUnit, error) {
	lu, err := stmf.GetLu(lu_uuid)
	if err != nil {
		return nil, err
	}
	basepath := strings.Join(
		strings.Split(lu.GetZvol(), "/")[0:len(strings.Split(lu.GetZvol(), "/"))-1], "/")

	cloneName := basepath + "/" + cloneAlias
	cloneZfsVolume, err := zfs.CreateFromSnapshot(lu.GetZvol()+"@"+snapname, cloneName, "")
	if err != nil {
		return nil, err
	}

	clonedLu, err := stmf.CreateLu(cloneZfsVolume.Dataset, "-p alias="+cloneAlias)
	if err != nil {
		return nil, err
	}

	return clonedLu, nil
}
