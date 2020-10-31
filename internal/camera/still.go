package camera

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"time"
)

func TakePhoto(directory string) (string, error) {
	now := time.Now()
	imageName := fmt.Sprintf("%s.jpg", now.Format(time.RFC3339))
	imagePath := path.Join(directory, imageName)

	err := runCommand("raspistill", "-o", imagePath)
	if err != nil {
		return "", err
	}

	return imagePath, nil
}

func runCommand(program string, args ...string) error {
	cmd := exec.Command(program, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run %s: %s\n%s\n%s", err, program, stdout.String(), stderr.String())
	}

	return nil
}
