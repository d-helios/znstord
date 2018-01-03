package zfs_test

import (
	"errors"
	"github.com/d-helios/znstord/zfs"
	"testing"
)

const (
	mb_size = 1048576

	dataset_name  = "rpool/test"
	snapshot_name = "snap1"
	clone_name    = "rpool/test_clone"
)

func TestCreateAndDestroy(t *testing.T) {
	// try to create dataset
	dataset, err := zfs.CreateFilesystem(dataset_name, "", 1*mb_size)
	if err != nil {
		t.Fatalf(err.Error())
	}

	createdDataset, err := zfs.GetDataset(dataset.Dataset)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := createdDataset.RefreshProps(); err != nil {
		t.Fatal(zfs.Error{
			Debug: err.Error(),
			Err:   err,
		})
	}

	if err := dataset.RefreshProps(); err != nil {
		t.Fatal(zfs.Error{
			Debug: err.Error(),
			Err:   err,
		})
	}

	if createdDataset.Dataset != dataset.Dataset ||
		createdDataset.Props.(*zfs.FsDataset).Creation != dataset.Props.(*zfs.FsDataset).Creation {
		t.Fatal(zfs.Error{
			Debug: "Dataset metainformation does not match.",
			Err:   errors.New("created and getted datasets not match"),
		}.Error())
	}

	// try to destroy dataset
	dataset.Destroy("")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDataset_Snapshots(t *testing.T) {
	// create dataset
	dataset, err := zfs.CreateFilesystem(dataset_name, "", 1*mb_size)
	if err != nil {
		t.Fatal(err.Error())
	}

	// create snapshot
	snapshot, err := dataset.Snapshot(snapshot_name)
	if err != nil {
		t.Fatal(err.Error())
	}

	// create volume from snapshot
	clone, err := zfs.CreateFromSnapshot(snapshot.Dataset, clone_name, "")
	if err != nil {
		t.Fatal(err.Error())
	}

	// destroy clone
	if err := clone.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}

	// destroy snapshot
	if err := snapshot.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}

	// destroy dataset
	if err := dataset.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}
}

func TestCreateVolume(t *testing.T) {
	// create volume
	volume, err := zfs.CreateVolume(dataset_name, "", true, 500*mb_size)
	if err != nil {
		t.Fatalf(err.Error())
	}

	if err := volume.RefreshProps(); err != nil {
		t.Fatal(
			zfs.Error{
				Err:   err,
				Debug: err.Error(),
			})
	}

	// check volume props
	if volume.Props.(*zfs.VolDataset).Volsize != 500*mb_size {
		t.Fatal(zfs.Error{
			Err: errors.New("zvol size mistmatch"),
		}.Error())
	}

	// destroy volume
	if err := volume.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}
}

func TestDataset_SetProp(t *testing.T) {
	// create volume
	volume, err := zfs.CreateVolume(dataset_name, "", true, 500*mb_size)
	if err != nil {
		t.Fatal(err.Error())
	}

	// set props
	if err := volume.SetProp("compression", "zle"); err != nil {
		t.Fatal(err.Error())
	}

	// refresh zfs props
	volume.RefreshProps()

	// check if volume updated
	if volume.Props.(*zfs.VolDataset).Compression != "zle" {
		t.Fatal(zfs.Error{
			Err: errors.New("Volume object does not updated"),
		}.Error())
	}

	// destroy volume
	if err := volume.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}
}

func TestDataset_Rename(t *testing.T) {
	// create volume
	volume, err := zfs.CreateVolume(dataset_name, "", true, 500*mb_size)
	if err != nil {
		t.Fatal(err.Error())
	}

	if err := volume.Rename(volume.Dataset + "_renamed"); err != nil {
		t.Fatal(err.Error())
	}

	// destroy volume
	if err := volume.Destroy(""); err != nil {
		t.Fatal(err.Error())
	}
}

func TestIsDatasetExist(t *testing.T) {
	err := zfs.IsDatasetExist("rpool")
	if err != nil {
		t.Fatalf("Dataset %s does not exists. Err: %s", "rpool", err.Error())
	}
}
