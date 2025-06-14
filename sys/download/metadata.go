package download

import (
	"bufio"
	"encoding/json"
	"home-media/sys/command"
	"home-media/sys/runtime"
)

func MetadataShell(shell *runtime.Shell, streamManager *runtime.StreamManager) {
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
	shell.SetVar("EXEC_BIN", cmdFrag.ExecBin)
	shell.ExitCode = 0
	shell.SetVar("FFMPEG_INPUT_FILE", cmdFrag.Input)

	func() {
		var commandName = shell.ReadVar("EXEC_BIN")
		var arguments []string
		arguments = append(arguments, shell.ReadVar("FFMPEG_INPUT_FILE"))
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
