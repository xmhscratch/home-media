package duration

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
	shell.SetVar("FFMPEG_INPUT_FILE", cmdFrag.Input)

	func() {
		var pipelineWaitgroup []func() error
		pipeReader1, pipeWriter1, err := runtime.NewPipe()
		if err != nil {
			shell.HandleError(err)
			return
		}

		func() {
			var commandName = `echo`
			var arguments []string
			arguments = append(arguments, `ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 '`+shell.ReadVar("FFMPEG_INPUT_FILE")+`'`)
			var command = shell.Command(commandName, arguments...)
			streamManager := streamManager.Clone()
			streamManager.Add(`1`, pipeWriter1, true)
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
			if err := command.Start(); err != nil {
				shell.HandleError(err)
				return
			}
			pipelineWaitgroup = append(pipelineWaitgroup, func() error {
				defer streamManager.Destroy()
				return command.Wait()
			})

		}()
		func() {
			var commandName = `sh`
			var arguments []string
			var command = shell.Command(commandName, arguments...)
			streamManager := streamManager.Clone()
			streamManager.Add(`0`, pipeReader1, false)
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
			if err := command.Start(); err != nil {
				shell.HandleError(err)
				return
			}
			pipelineWaitgroup = append(pipelineWaitgroup, func() error {
				defer streamManager.Destroy()
				return command.Wait()
			})

		}()
		for i, wait := range pipelineWaitgroup {
			if err := wait(); err != nil {
				shell.HandleError(err)
			}
			if i < (len(pipelineWaitgroup) - 1) {
				shell.ExitCode = 0
			}
		}

	}()
}
