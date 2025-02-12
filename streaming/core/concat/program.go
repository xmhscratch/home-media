package concat

import (
	"bufio"
	"encoding/json"
	"home-media/streaming/core/command"
	"home-media/streaming/core/runtime"
)

func Main(shell *runtime.Shell, streamManager *runtime.StreamManager) {
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
	shell.SetVar("INPUT", cmdFrag.Input)
	shell.ExitCode = 0
	shell.SetVar("OUTPUT", cmdFrag.Output)

	func() {
		var commandName = shell.ReadVar("EXECUTOR") // `../ci/concat.sh`
		var arguments []string
		arguments = append(arguments, shell.ReadVar("INPUT"))
		arguments = append(arguments, shell.ReadVar("OUTPUT"))
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
