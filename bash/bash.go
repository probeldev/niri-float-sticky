package bash

import (
	"bufio"
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os/exec"
)

func RunCommand(command string) ([]byte, error) {
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
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}

	// Возвращаем stdout как результат
	return stdout.Bytes(), nil
}

func RunAndListenCommand(command string) (<-chan []byte, error) {
	cmd := exec.Command("sh", "-c", command)

	var stderr bytes.Buffer
	stdoutReader, stdoutWriter := io.Pipe()
	cmd.Stdout = stdoutWriter
	cmd.Stderr = &stderr

	linesCh := make(chan []byte)

	go func() {
		defer close(linesCh)

		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			linesCh <- scanner.Bytes()
		}
		if err := scanner.Err(); err != nil {
			log.Errorf("%v: %s", err, stderr.String())
		}
	}()

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("%w: %s", err, stderr.String())
	}

	go func() {
		if err := cmd.Wait(); err != nil {
			log.Errorf("%v: %s", err, stderr.String())
		}
		if err := stdoutWriter.Close(); err != nil {
			log.Error(err)
		}
	}()

	return linesCh, nil
}
