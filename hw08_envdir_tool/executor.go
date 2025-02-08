package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204

	// Установка переменных окружения
	command.Env = os.Environ() // Сначала копируем текущее окружение

	// Применяем наши переменные окружения
	for name, value := range env {
		if value.NeedRemove {
			command.Env = removeEnv(command.Env, name)
		} else {
			command.Env = append(command.Env, name+"="+value.Value)
		}
	}

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	command.Start()

	command.Wait()

	return command.ProcessState.ExitCode()
}

func removeEnv(env []string, name string) []string {
	result := make([]string, 0, len(env))
	prefix := name + "="
	for _, e := range env {
		if len(e) > len(prefix) && e[:len(prefix)] != prefix {
			result = append(result, e)
		}
	}
	return result
}
