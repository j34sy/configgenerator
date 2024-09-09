package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j34sy/configgenerator/pkg/devices"
)

func WriteMLSwitchConfigFile(mlSwitch *devices.MLSwitch, dirPath string) {
	fileName := mlSwitch.Name + "." + mlSwitch.Domain + ".cisco"

	filePath := filepath.Join(dirPath, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating or clearing file:", err)
		return
	}
	file.Close()

	file, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	file.WriteString("hostname " + mlSwitch.Name + "\n")
	file.WriteString("ip domain-name " + mlSwitch.Domain + "\n")

}
