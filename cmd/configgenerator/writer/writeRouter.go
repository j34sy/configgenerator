package writer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/j34sy/configgenerator/pkg/devices"
)

func writeRouterConfigFile(router *devices.Router, dirPath string) {
	fileName := router.Name + "." + router.Domain + ".cisco"

	filePath := filepath.Join(dirPath, fileName)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error creating or clearing file:", err)
		return
	}
	defer file.Close()

}
