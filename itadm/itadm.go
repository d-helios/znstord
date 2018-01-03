package itadm

import (
	"strconv"
	"strings"
)

// wrapper for itadm
func cmdItadm(arg ...string) ([][]string, error) {
	c := command{Command: "itadm"}
	return c.Run(" ", arg...)
}

// ListTargetPortGroups - list itadm target port group (TPG)
func ListTargetPortGroups(tpg string) ([]*TargetPortGroup, error) {
	args := []string{"list-tpg", "-v"}

	if tpg != "" {
		args = append(args, tpg)
	}

	out, err := cmdItadm(args...)
	if err != nil {
		return nil, err
	}

	for len(out) == 0 {
		return nil, err
	}

	var tpgs []*TargetPortGroup
	var tpgTmp *TargetPortGroup

	for i := 0; i < len(out); i++ {
		if out[i][0] == "TARGET" {
			continue
		}
		if out[i][0] != "portals:" {
			count, err := strconv.ParseUint(out[i][1], 10, 64)
			if err != nil {
				return nil, err
			}
			tpgTmp = &TargetPortGroup{
				TargetPortGroup: out[i][0],
				Count:           count,
			}
		}
		if out[i][0] == "portals:" {
			tpgTmp.Portals = strings.Split(out[i][1], ",")
			tpgs = append(tpgs, tpgTmp)
		}
	}

	return tpgs, nil
}

// GetTargetPortGroup - get specified target port group.
func GetTargetPortGroup(tpg string) (*TargetPortGroup, error) {
	tpgs, err := ListTargetPortGroups(tpg)
	if err != nil {
		return nil, err
	}
	return tpgs[0], nil
}

// CreateTargetPortGroup - create target port group.
func CreateTargetPortGroup(tpg string, ipaddrs []string) (*TargetPortGroup, error) {
	args := []string{"create-tpg", tpg}
	args = append(args, ipaddrs...)

	_, err := cmdItadm(args...)
	if err != nil {
		return nil, err
	}

	return GetTargetPortGroup(tpg)
}

// Delete - delete target port group
func (tpg *TargetPortGroup) Delete(force bool) error {
	args := []string{"delete-tpg"}

	if force {
		args = append(args, "-f")
	}

	args = append(args, tpg.TargetPortGroup)

	_, err := cmdItadm(args...)
	if err != nil {
		return err
	}
	return nil
}

// CreateTarget - Create Target.
// TODO: add interface to chap authentication.
func CreateTarget(iqn, alias, targetPortGroup string) (*Target, error) {
	args := []string{"create-target", "-a", "default"}

	if alias != "" {
		args = append(args, "-l", alias)
	}

	if targetPortGroup != "" {
		args = append(args, "-n", iqn)
	}

	if iqn != "" {
		args = append(args, "-t", targetPortGroup)
	}

	out, err := cmdItadm(args...)
	if err != nil {
		return nil, err
	}

	created_target := out[0][1]

	return GetTarget(created_target)
}

// ListTargets - list available targets
func ListTargets(iqn string) ([]*Target, error) {
	args := []string{"list-target", "-v"}

	if iqn != "" {
		args = append(args, iqn)
	}

	out, err := cmdItadm(args...)
	if err != nil {
		return nil, err
	}

	for len(out) == 0 {
		return nil, err
	}

	var targets []*Target
	var targetTmp *Target

	for i := 0; i < len(out); i++ {
		if out[i][0] == "TARGET" {
			continue
		}
		// check iqn.
		if out[i][0][0:4] == "iqn." {
			sessionCount, err := strconv.ParseUint(out[i][2], 10, 64)
			if err != nil {
				return nil, err
			}
			targetTmp = &Target{
				IQN:      out[i][0],
				State:    out[i][1],
				Sessions: sessionCount,
			}
		}
		if out[i][0] == "alias:" {
			targetTmp.Alias = out[i][1]
		}
		if out[i][0] == "auth:" {
			targetTmp.Auth = out[i][1]
		}
		if out[i][0] == "targetchapuser:" {
			targetTmp.TargetChapUser = out[i][1]
		}
		if out[i][0] == "targetchapsecret:" {
			targetTmp.TargetChapSecret = out[i][1]
		}
		if out[i][0] == "tpg-tags:" {
			targetTmp.TpgTags = out[i][1]
			targets = append(targets, targetTmp)
		}
	}

	return targets, nil
}

// GetTarget - Get specified target
func GetTarget(iqn string) (*Target, error) {
	target, err := ListTargets(iqn)
	if err != nil {
		return nil, err
	}
	return target[0], nil
}

// Delete - delete Target
func (target *Target) Delete(force bool) error {
	args := []string{"delete-target"}

	if force {
		args = append(args, "-f")
	}

	args = append(args, target.IQN)

	_, err := cmdItadm(args...)

	return err
}
