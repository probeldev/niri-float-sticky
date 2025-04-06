package bash

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunCommand(command string) (string, error) {
	// Создаем команду для выполнения в shell
	cmd := exec.Command("sh", "-c", command)

	// Буферы для захвата stdout и stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Выполняем команду
	err := cmd.Run()

	// Если есть ошибка, возвращаем stderr как часть ошибки
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}

	// Возвращаем stdout как результат
	return stdout.String(), nil
}
