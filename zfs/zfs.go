package zfs

import (
	"strconv"
	"strings"
)

// List datasets starting from the BaseDsPath (ex: tank/nfs )
func ListDatasets(datasetType, baseDsPath string, recursive bool, depth uint64) ([]*Dataset, error) {
	args := []string{"list", "-H", "-t", datasetType, "-o", "name"}

	if recursive {
		// append "-r" argument if recursive is specified
		args = append(args, "-r")
	}

	if depth > 0 {
		args = append(args, []string{"-d", strconv.FormatUint(depth, 10)}...)
	}

	if baseDsPath != "" {
		args = append(args, baseDsPath)
	}

	out, err := cmdZfs(args...)

	if err != nil {
		return nil, err
	}

	var dataset []*Dataset

	// enumerate datasets and fill array
	var ds *Dataset

	for _, data := range out {
		ds = &Dataset{
			Dataset: data[0],
		}

		dataset = append(dataset, ds)
	}
	return dataset, nil
}

// Get dataset
func GetDataset(dataset string) (*Dataset, error) {
	datasets, err := ListDatasets("all", dataset, false, 0)
	if err != nil {
		return nil, err
	}

	return datasets[0], nil
}

// Create Volume
func CreateVolume(datasetName, options string, thin bool, size uint64) (*Dataset, error) {
	args := []string{"-V", strconv.FormatUint(size, 10)}

	if thin {
		args = append(args, "-s")
	}

	if options != "" {
		args = append(args, options)
	}

	ds, err := CreateDataset(datasetName, strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// Is Dataset Exists
func IsDatasetExist(datasetName string) error {
	args := []string{"-a", "/dev/zvol/rdsk/" + datasetName}
	return cmdTest(args...)
}

// Create Filesystem
func CreateFilesystem(dataset, options string, size uint64) (*Dataset, error) {
	args := []string{"-o quota=" + strconv.FormatUint(size, 10)}

	if options != "" {
		args = append(args, options)
	}

	ds, err := CreateDataset(dataset, strings.Join(args, " "))
	if err != nil {
		return nil, err
	}
	return ds, nil
}

// Create dataset from snapshot (clone dataset)
func CreateFromSnapshot(snapshot, clone, options string) (*Dataset, error) {
	args := []string{"clone"}

	// add options, for ex: -o mountpoint=/mnt/a -o recordsize=8192
	if options != "" {
		optionList := strings.Split(strings.TrimSpace(options), " ")
		args = append(args, optionList...)
	}

	// append filesystem name
	args = append(args, snapshot, clone)

	_, err := cmdZfs(args...)

	if err != nil {
		return nil, err
	}

	ds := &Dataset{
		Dataset: clone,
	}

	return ds, nil
}

// Create dataset
func CreateDataset(dataset, options string) (*Dataset, error) {
	args := []string{"create"}

	// add options, for ex: -o mountpoint=/mnt/a -o recordsize=8192
	if options != "" {
		optionList := strings.Split(strings.TrimSpace(options), " ")
		args = append(args, optionList...)
	}

	args = append(args, dataset)

	_, err := cmdZfs(args...)

	if err != nil {
		return nil, err
	}

	ds := &Dataset{
		Dataset: dataset,
	}

	return ds, nil
}

// Share NFS
func (dataset *Dataset) ShareNfs(ro, rw, root string) error {
	args := []string{}
	OsRelease, err := GetOsRelease()
	if err != nil {
		return err
	}
	if OsRelease == OpenSolaris {
		args = []string{"ro=" + ro +
			",rw=" + rw +
			",root=" + root +
			",sec=sys"}
		dataset.SetProp("sharenfs", strings.Join(args, " "))
		if err != nil {
			return err
		}
		return dataset.RefreshProps()

	}
	if OsRelease == OracleSolaris {
		err := dataset.SetProp("share.nfs", "on")
		if err != nil {
			return err
		}

		err = dataset.SetProp("share.nfs.ro", ro)
		if err != nil {
			return err
		}

		err = dataset.SetProp("share.nfs.rw", rw)
		if err != nil {
			return err
		}

		err = dataset.SetProp("share.nfs.root", root)
		if err != nil {
			return err
		}

		err = dataset.SetProp("share.nfs.sec", "sys")
		if err != nil {
			return err
		}
	}
	return dataset.RefreshProps()
}

// simple nfs share method
func (dataset *Dataset) ShareNfsRW(hosts string) error {
	return dataset.ShareNfs("", hosts, hosts)
}

// unshare dataset
func (dataset *Dataset) UnshareNfs() error {
	OsRelease, err := GetOsRelease()
	if err != nil {
		return err
	}

	if OsRelease == OpenSolaris {
		return dataset.SetProp("sharenfs", "off")
	}

	if OsRelease == OracleSolaris {
		return dataset.SetProp("share.nfs", "off")
	}
	return nil
}

/*
Dataset methods
*/
func (dataset *Dataset) Destroy(options string) error {
	args := []string{"destroy"}

	// add options, for ex: -fnpRrv filesystem|volume|snapshot
	if options != "" {
		optionList := strings.Split(strings.TrimSpace(options), " ")
		args = append(args, optionList...)
	}

	// append filesystem name
	args = append(args, dataset.Dataset)

	_, err := cmdZfs(args...)

	if err != nil {
		return err
	}
	return nil
}

func (dataset *Dataset) Rename(newName string) error {
	args := []string{"rename"}

	// append filesystem name
	args = append(args, dataset.Dataset)
	args = append(args, newName)

	_, err := cmdZfs(args...)

	if err != nil {
		return err
	}

	//ds, err := GetDataset(newName)
	dataset.Dataset = newName

	return nil
}

func (dataset *Dataset) Promote() error {
	args := []string{"promote"}

	// append filesystem name
	args = append(args, dataset.Dataset)

	_, err := cmdZfs(args...)
	if err != nil {
		return err
	}
	return nil
}

func (dataset *Dataset) Rollback(snapname string) error {
	args := []string{"rollback"}

	// append filesystem name
	args = append(args, snapname)

	_, err := cmdZfs(args...)

	if err != nil {
		return err
	}

	return nil
}

func (dataset *Dataset) Snapshot(snapname string) (*Dataset, error) {
	args := []string{"snapshot"}

	// append filesystem name
	args = append(args, dataset.Dataset+"@"+snapname)

	_, err := cmdZfs(args...)

	if err != nil {
		return nil, err
	}

	ds, err := GetDataset(dataset.Dataset + "@" + snapname)
	if err != nil {
		return nil, err
	}

	return ds, err
}

// fill dataset properties
func (dataset *Dataset) RefreshProps() error {
	args := []string{"get", "-Hp", "-o", "property,value", "all"}
	args = append(args, dataset.Dataset)

	out, err := cmdZfs(args...)

	if err != nil {
		return err
	}

	dict := make(map[string]string)

	for j := range out {
		if out[j][0] == "custom:alias" {
			dict["alias"] = out[j][1]
			continue
		}

		switch out[j][0] {
		case "custom:alias":
			dict["alias"] = out[j][1]
			continue
		case "custom:sflag":
			dict["sflag"] = out[j][1]
			continue
		}

		// BUGS: If snapshot have not clones, it's reported "" value, instead of "none"
		// Not implemented in Oracle Solaris
		if out[j][0] == "clones" {
			if len(out[j]) <= 1 {
				dict[out[j][0]] = "none"
			} else {
				dict[out[j][0]] = out[j][1]
			}
			continue
		}
		dict[out[j][0]] = out[j][1]
	}

	switch dict["type"] {
	case Filesystem:
		dataset.Props = &FsDataset{}
	case Volume:
		dataset.Props = &VolDataset{}
	case Snapshot:
		dataset.Props = &SnapDataset{}
	}

	err = FillDataset(dataset.Props, dict)
	if err != nil {
		return err
	}

	return nil
}

func (dataset *Dataset) SetProp(parameter string, value interface{}) error {

	options := ""
	if IsItNumericProp(parameter, NumericProps) {
		options = parameter + "=" + strconv.FormatUint(value.(uint64), 10)
	} else {
		options = parameter + "=" + value.(string)
	}

	args := []string{"set", options, dataset.Dataset}

	_, err := cmdZfs(args...)
	if err != nil {
		return err
	}

	return nil
}
