package process

import (
	"maps"
	"slices"
	"testing"
)

// go test -v -run TestGetProcessList  -tags "dev"
func TestGetProcessList(t *testing.T) {
	processs, _ := GetProcessList()

	for _, vv := range processs {
		t.Log(vv.FullPath, vv.ProcessID, vv.ParentProcessID)
	}
}

// go test -v -run TestGetCmdline  -tags "dev"
func TestGetCmdline(t *testing.T) {
	fullpath, _, _ := GetCmdline(11432)
	t.Log(fullpath)
}

// go test -v -run TestGetProcessParentProcessMap  -tags "dev"
func TestGetProcessParentProcessMap(t *testing.T) {
	processs, _ := GetProcessList()
	processsMap := map[uint32]*Process{}

	parentProcessMap := map[uint32][]uint32{}

	var f func(pid uint32, processsMap map[uint32]*Process) uint32

	f = func(pid uint32, processsMap map[uint32]*Process) uint32 {
		if vv, ok := processsMap[pid]; ok {
			if vv.ParentProcessID == 0 {
				return vv.ProcessID
			}
			return f(vv.ParentProcessID, processsMap)
		}
		return pid
	}

	for _, vv := range processs {
		processsMap[vv.ProcessID] = vv
		if vv.ParentProcessID == 0 {
			parentProcessMap[vv.ProcessID] = []uint32{}
		}
	}

	for _, vv := range processs {
		if vv.ParentProcessID != 0 {
			topParentPid := f(vv.ParentProcessID, processsMap)
			parentProcessMap[topParentPid] = append(parentProcessMap[topParentPid], vv.ProcessID)
		}
	}
	t.Log(slices.Sorted(maps.Keys(parentProcessMap)))

	t.Log(slices.Sorted(maps.Keys(processsMap)))
}

// go test -v -run TestIsProcess64BitByName  -tags "dev"
func TestIsProcess64BitByName(t *testing.T) {
	t.Log(IsProcess64BitByName("WeChat.exe"))
}

// go test -v -run TestIsProcess64BitByNameTest  -tags "dev"
func TestIsProcess64BitByNameTest(t *testing.T) {
	t.Log(IsProcess64BitByName("Test1.exe"))
}
