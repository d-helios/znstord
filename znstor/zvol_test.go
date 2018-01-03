package znstor_test

import (
	"fmt"
	"github.com/d-helios/stmf"
	"github.com/d-helios/zfs"
	"github.com/d-helios/znstor"
	"github.com/twinj/uuid"
	"strings"
	"testing"
)

const mb_size = 1048576
const basepath = "tank"
