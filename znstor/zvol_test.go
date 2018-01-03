package znstor_test

import (
	"fmt"
	"github.com/d-helios/znstord/stmf"
	"github.com/d-helios/znstord/zfs"
	"github.com/d-helios/znstord/znstor"
	"github.com/twinj/uuid"
	"strings"
	"testing"
)

const mb_size = 1048576
const basepath = "tank"
