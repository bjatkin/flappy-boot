//go:build local

package save

import (
	"fmt"
	"os"
)

// the game only really needs the first few bytes so no reason to
// store all that data if we don't need it
const DataLen = 0x00FF

func LoadData(path string) {
	file, err := os.ReadFile(path)
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
	err := os.WriteFile(path, data, 0o0664)
	if err != nil {
		fmt.Printf("failed to save data: %v\n", err)
	}
}
