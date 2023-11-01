//go:build web

package save

import (
	"encoding/base64"
	"fmt"
	"syscall/js"
)

// the game only really needs the first few bytes so no reason to
// store all that data if we don't need it
const DataLen = 0x00FF

func LoadData(path string) {
	file, err := jsReadFile(path)
	if err != nil {
		// don't error out just because a save file is missing
		fmt.Printf("failed to load save file: %v\n", err)
	}

	// set up SRAM before creating a new engine
	for i := 0; i < DataLen; i++ {
		if i < len(file) {
			SRAM[i] = SRAMValue(file[i])
			continue
		}

		// 0xFF is the default value for the SRAM so set that up here
		SRAM[i] = 0xFF
	}
}

func SaveData(path string, data []byte) {
	jsWriteFile(path, data)
}

func jsWriteFile(fileName string, data []byte) {
	stringData := base64.StdEncoding.EncodeToString(data)
	fmt.Println("writing to local storage", fileName, stringData)

	global := js.Global()
	localStorage := global.Get("localStorage")
	_ = localStorage.Call("setItem", fileName, stringData)
}

func jsReadFile(fileName string) ([]byte, error) {
	fmt.Println("reading from local storage", fileName)
	global := js.Global()
	localStorage := global.Get("localStorage")
	fileContents := localStorage.Call("getItem", fileName)
	stringData := fileContents.String()
	return base64.StdEncoding.DecodeString(stringData)
}
