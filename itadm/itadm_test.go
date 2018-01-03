package itadm_test

import (
	"errors"
	"fmt"
	"github.com/d-helios/itadm"
	"testing"
)

var tpgs = []string{
	"tpg1",
	"tpg2",
	"tpg3",
	"tpg4",
}

var ipaddrs = [][]string{
	[]string{"192.168.2.100"},
	[]string{"192.168.2.101"},
	[]string{"192.168.3.100", "192.168.4.100"},
	[]string{"192.168.3.101", "192.168.4.101", "192.168.5.101", "192.168.6.101"},
}

var targets = []string{
	"iqn.1986-03.com.sun:9f0f1c727925",
	"iqn.1986-03.com.sun:9f0f1c727926",
	"iqn.1986-03.com.sun:9f0f1c727927",
	"iqn.1986-03.com.sun:9f0f1c727928",
}

func TestCreateTargetPortGroup(t *testing.T) {
	for index := range tpgs {
		tpg, err := itadm.CreateTargetPortGroup(tpgs[index], ipaddrs[index])
		if err != nil {
			t.Fatalf(err.Error())
		}
		if tpg.TargetPortGroup != tpgs[index] {
			t.Fatalf(itadm.Error{
				Err: errors.New(
					fmt.Sprintf("TargetPortGroup name is not "+
						"equal to specified one. TGP: %q", tpg)),
				Debug:  "",
				Stderr: "",
			}.Error())
		}
		if tpg.Count != uint64(len(ipaddrs[index])) {
			t.Fatalf(itadm.Error{
				Err: errors.New(
					fmt.Sprintf("TargetPortGroup count "+
						"not equal to specified. TGP: %q", tpg)),
				Debug:  "",
				Stderr: "",
			}.Error())
		}
	}
}

func TestCreateTarget(t *testing.T) {
	for index := range targets {
		alias := targets[index][len(targets[index])-4 : len(targets[index])]
		tpg, err := itadm.GetTargetPortGroup(tpgs[index])
		if err != nil {
			t.Fatalf(err.Error())
		}
		target, err := itadm.CreateTarget(targets[index], alias, tpg.TargetPortGroup)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if target.IQN != targets[index] {
			t.Fatalf(itadm.Error{
				Err: errors.New(
					fmt.Sprintf("Target IQN not match specified."+
						".Target: %q", target)),
				Debug:  "",
				Stderr: "",
			}.Error())
		}
	}
}

func TestTarget_Delete(t *testing.T) {
	targets, err := itadm.ListTargets("")
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, target := range targets {
		err := target.Delete(true)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}

func TestTargetPortGroup_Delete(t *testing.T) {
	tpgs, err := itadm.ListTargetPortGroups("")
	if err != nil {
		t.Fatalf(err.Error())
	}

	for _, tpg := range tpgs {
		err := tpg.Delete(true)
		if err != nil {
			t.Fatalf(err.Error())
		}
	}
}
