package download

import (
	"bufio"
	"encoding/json"
	"home-media/sys/command"
	"home-media/sys/runtime"
)

func DownloadShell(shell *runtime.Shell, streamManager *runtime.StreamManager) {
	var cmdFrag *command.CommandFrags = &command.CommandFrags{}

	if stream, err := streamManager.Get(`0`); err != nil {
		shell.HandleError(err)
		return
	} else {
		var reader = bufio.NewReader(stream)
		message, _, _ := reader.ReadLine()
		if err = json.Unmarshal(message, cmdFrag); err != nil {
			shell.HandleError(err)
			return
		}
	}

	shell.ExitCode = 0
	shell.SetVar("EXECUTOR", cmdFrag.ExecBin)
	shell.ExitCode = 0
	shell.SetVar("DOWNLOAD_URL", cmdFrag.DownloadURL)
	shell.ExitCode = 0
	shell.SetVar("OUTPUT_DIR", cmdFrag.Output)
	shell.ExitCode = 0
	shell.SetVar("BASE_URL", cmdFrag.BaseURL)
	shell.ExitCode = 0
	shell.SetVar("ROOT_DIR", cmdFrag.RootDir)

	func() {
		var commandName = shell.ReadVar("EXECUTOR") // `../ci/download.sh`
		var arguments []string
		arguments = append(arguments, shell.ReadVar("DOWNLOAD_URL"))
		arguments = append(arguments, shell.ReadVar("OUTPUT_DIR"))
		arguments = append(arguments, shell.ReadVar("BASE_URL"))
		arguments = append(arguments, shell.ReadVar("ROOT_DIR"))
		var command = shell.Command(commandName, arguments...)
		streamManager := streamManager.Clone()
		defer streamManager.Destroy()
		if stream, err := streamManager.Get(`0`); err != nil {
			shell.HandleError(err)
			return
		} else {
			command.Stdin = stream
		}
		if stream, err := streamManager.Get(`1`); err != nil {
			shell.HandleError(err)
			return
		} else {
			command.Stdout = stream
		}
		if stream, err := streamManager.Get(`2`); err != nil {
			shell.HandleError(err)
			return
		} else {
			command.Stderr = stream
		}
		if err := command.Run(); err != nil {
			shell.HandleError(err)
			return
		}
		shell.ExitCode = command.ProcessState.ExitCode()

	}()
}
