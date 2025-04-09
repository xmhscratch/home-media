package download

import (
	"bufio"
	"encoding/json"
	"home-media/sys/command"
	"home-media/sys/runtime"
)

func ExtractShell(shell *runtime.Shell, streamManager *runtime.StreamManager) {
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
	shell.SetVar("FFMPEG_INPUT_FILE", cmdFrag.Input)
	shell.ExitCode = 0
	shell.SetVar("STREAM_INDEX", cmdFrag.StreamIndex)
	shell.ExitCode = 0
	shell.SetVar("LANG_CODE", cmdFrag.LangCode)
	shell.ExitCode = 0
	shell.SetVar("FFMPEG_OUTPUT_FILE", cmdFrag.Output)

	func() {
		var commandName = shell.ReadVar("EXECUTOR")
		var arguments []string
		arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_INPUT_FILE")+`'`)
		arguments = append(arguments, shell.ReadVar("STREAM_INDEX"))
		arguments = append(arguments, shell.ReadVar("LANG_CODE"))
		arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_OUTPUT_FILE")+`'`)
		// fmt.Println(commandName, arguments)
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
