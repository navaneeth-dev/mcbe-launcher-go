package main

import (
	"fmt"
	"log"
)

const STORE_OWNERSHIP_CHECK_OFFSET = 0x120E7A6

func main() {
	process, err := ProcessByName("Minecraft.Windows")
	if err != nil {
		log.Panicf("Minecraft running? Error: %s", err.Error())
		fmt.Scanln()
		return
	}
	log.Printf("Base: 0x%06X 0x%X", process.ModBaseAddr, process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)

	err = process.WriteBytesSigned(process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET, []int8{-21})
	if err != nil {
		fmt.Printf("Failed to patch: %s", err.Error())
		fmt.Scanln()
		return
	}

	fmt.Printf("Patched! %x\n", process.ModBaseAddr+STORE_OWNERSHIP_CHECK_OFFSET)
	fmt.Println("Press ENTER to exit...")
	fmt.Scanln()
}
