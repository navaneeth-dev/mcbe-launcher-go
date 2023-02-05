package main

import (
	"fmt"
	"log"
	"os/exec"
	"syscall"
)

// const STORE_OWNERSHIP_CHECK_OFFSET = 0x16FB7C6
const STORE_OWNERSHIP_CHECK_OFFSET = 0xDF79F6

var (
	k32                = syscall.NewLazyDLL("kernel32.dll")
	WriteProcessMemory = k32.NewProc("WriteProcessMemory")
	FreeConsole        = k32.NewProc("FreeConsole")
)

func main() {
	launchMinecraft()

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

	log.Printf("Base: 0x%06X 0x%X", process.ModBaseAddr, process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)
	fmt.Printf("Patched! %x\n", process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)
	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}

func launchMinecraft() {
	_, err := exec.Command("C:\\Windows\\system32\\cmd.exe", "/c", "start Minecraft://").Output()
	if err != nil {
		panic(err)
	}
}
