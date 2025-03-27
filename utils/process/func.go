package process

import (
	"fmt"
	"monitor/utils/win"
	"reflect"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

type ProcessInfoClass uint32

const (
	// ProcessInfoClass enumeration values that can be used as arguments to
	// NtQueryInformationProcess

	// ProcessBasicInformation returns a pointer to
	// the Process Environment Block (PEB) structure.
	ProcessBasicInformation ProcessInfoClass = 0

	// ProcessDebugPort returns a uint32 that is the port number for the
	// debugger of the process.
	ProcessDebugPort = 7

	// ProcessWow64Information returns whether a process is running under
	// WOW64.
	ProcessWow64Information = 26

	// ProcessImageFileName returns the image file name for the process, as a
	// UnicodeString struct.
	ProcessImageFileName = 27

	// ProcessBreakOnTermination returns a uintptr that tells if the process
	// is critical.
	ProcessBreakOnTermination = 29

	// ProcessSubsystemInformation returns the subsystem type of the process.
	ProcessSubsystemInformation = 75
)

type ProcessBasicInformationStruct struct {
	Reserved1       uintptr
	PebBaseAddress  uintptr
	Reserved2       [2]uintptr
	UniqueProcessID uintptr
	// Undocumented:
	InheritedFromUniqueProcessID uintptr
}

type Process struct {
	Size            uint32
	Usage           uint32
	ProcessID       uint32
	DefaultHeapID   uintptr
	ModuleID        uint32
	Threads         uint32
	ParentProcessID uint32
	PriClassBase    int32
	Flags           uint32
	FullPath        string
}

func GetProcessList() ([]*Process, error) {
	return getProcessList()
}

func GetCmdline(pid uint32) (string, string, error) {
	if pid == 0 { // 系统进程,无法读取
		return "", "", nil
	}

	h, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ, false, pid)
	if err != nil {
		if e, ok := err.(windows.Errno); ok && e == windows.ERROR_ACCESS_DENIED {
			return "", "", nil // 没权限,忽略这个进程
		}
		return "", "", err
	}
	defer windows.CloseHandle(h)

	var pbi struct {
		ExitStatus                   uint32
		PebBaseAddress               uintptr
		AffinityMask                 uintptr
		BasePriority                 int32
		UniqueProcessId              uintptr
		InheritedFromUniqueProcessId uintptr
	}
	pbiLen := uint32(unsafe.Sizeof(pbi))
	err = windows.NtQueryInformationProcess(h, windows.ProcessBasicInformation, unsafe.Pointer(&pbi), pbiLen, &pbiLen)
	if err != nil {
		return "", "", err
	}

	var addr uint64
	d := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&addr)),
		Len:  8, Cap: 8,
	}))
	err = windows.ReadProcessMemory(h, pbi.PebBaseAddress+32, // ntddk.h,ProcessParameters偏移32字节
		&d[0], uintptr(len(d)), nil)
	if err != nil {
		return "", "", err
	}

	var commandLine windows.NTUnicodeString
	Len := unsafe.Sizeof(commandLine)
	d = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(&commandLine)),
		Len:  int(Len), Cap: int(Len),
	}))
	err = windows.ReadProcessMemory(h, uintptr(addr+112), // winternl.h,分析文件偏移
		&d[0], Len, nil)
	if err != nil {
		return "", "", err
	}

	cmdData := make([]uint16, commandLine.Length/2)
	d = *(*[]byte)(unsafe.Pointer(&cmdData))
	err = windows.ReadProcessMemory(h, uintptr(unsafe.Pointer(commandLine.Buffer)),
		&d[0], uintptr(commandLine.Length), nil)
	if err != nil {
		return "", "", err
	}

	var exePath [windows.MAX_PATH]uint16
	var exePathLen uint32 = windows.MAX_PATH
	err = windows.QueryFullProcessImageName(h, 0, &exePath[0], &exePathLen)
	if err != nil {
		return "", "", err
	}
	fullPath := syscall.UTF16ToString(exePath[:exePathLen])
	return fullPath, windows.UTF16ToString(cmdData), nil
}

func getProcessList() ([]*Process, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}

	data := make([]*Process, 0, 30)
	for {
		// procName := syscall.UTF16ToString(procEntry.ExeFile[:])
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ|windows.PROCESS_TERMINATE, false, procEntry.ProcessID)
		if err != nil {
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return data, err
			}
			continue
		}

		var exePath [windows.MAX_PATH]uint16
		var exePathLen uint32 = windows.MAX_PATH
		fullPath := ""
		err = windows.QueryFullProcessImageName(handle, 0, &exePath[0], &exePathLen)
		if err != nil {
			windows.CloseHandle(handle)
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return data, err
			}
			// continue
		} else {
			fullPath = syscall.UTF16ToString(exePath[:exePathLen])
		}

		data = append(data, &Process{
			Size:            procEntry.Size,
			Usage:           procEntry.Usage,
			ProcessID:       procEntry.ProcessID,
			DefaultHeapID:   procEntry.DefaultHeapID,
			ModuleID:        procEntry.ModuleID,
			Threads:         procEntry.Threads,
			ParentProcessID: procEntry.ParentProcessID,
			PriClassBase:    procEntry.PriClassBase,
			Flags:           procEntry.Flags,
			FullPath:        fullPath,
		})

		windows.CloseHandle(handle)
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return data, err
		}
	}
	return data, err
}

func CloseProcess(find func(ProcName string, FullPath string) bool) error {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return err
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return err
	}

	for {
		procName := syscall.UTF16ToString(procEntry.ExeFile[:])

		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ|windows.PROCESS_TERMINATE, false, procEntry.ProcessID)
		if err != nil {
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return err
			}
			continue
		}

		var exePath [windows.MAX_PATH]uint16
		var exePathLen uint32 = windows.MAX_PATH
		err = windows.QueryFullProcessImageName(handle, 0, &exePath[0], &exePathLen)
		if err != nil {
			windows.CloseHandle(handle)
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return err
			}
			continue
		}

		fullPath := syscall.UTF16ToString(exePath[:exePathLen])

		if find(procName, fullPath) {
			windows.TerminateProcess(handle, 1)
		}
		windows.CloseHandle(handle)

		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return err
		}
	}

	return nil
}

func IsProcess64BitByName(name string) (uint32, bool, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return 0, false, fmt.Errorf("IsProcess64BitByName => CreateToolhelp32Snapshot err = %w", err)
	}
	defer windows.CloseHandle(snapshot)
	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return 0, false, err
	}

	for {
		procName := syscall.UTF16ToString(procEntry.ExeFile[:])
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ|windows.PROCESS_TERMINATE, false, procEntry.ProcessID)
		if err != nil {
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return 0, false, err
			}
			continue
		}
		if procName == name {
			isWow64 := false
			err := windows.IsWow64Process(handle, &isWow64)
			windows.CloseHandle(handle)
			return procEntry.ProcessID, isWow64, err
		}
		windows.CloseHandle(handle)
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return 0, false, err
		}
	}

	return 0, false, nil
}

func GetPpids() (map[uint32]uint32, error) {
	return getPpids()
}

func getPpids() (map[uint32]uint32, error) {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)

	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}

	data := make(map[uint32]uint32, 50)

	for {
		handle, err := windows.OpenProcess(windows.PROCESS_QUERY_INFORMATION|windows.PROCESS_VM_READ|windows.PROCESS_TERMINATE, false, procEntry.ProcessID)
		if err != nil {
			err = windows.Process32Next(snapshot, &procEntry)
			if err != nil {
				return data, err
			}
			continue
		}
		if procEntry.ParentProcessID > 0 {
			data[procEntry.ProcessID] = procEntry.ParentProcessID
		}

		windows.CloseHandle(handle)
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return data, err
		}
	}
	return data, err
}

func FindWindows(pid uint32) []win.HWND {
	var rh []win.HWND = nil
	var pidt uint32 = 0
	win.EnumWindows(func(hwnd win.HWND, lParam uintptr) bool {
		pidt = 0
		win.GetWindowThreadProcessId(hwnd, &pidt)
		if pid == pidt {
			rh = append(rh, hwnd)
		}
		return true
	}, 0)
	return rh
}
