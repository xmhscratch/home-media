package segment

import (
	"bufio"
	"encoding/json"
	"home-media/sys/command"
	"home-media/sys/runtime"
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
	shell.ExitCode = 0
	shell.SetVar("FFMPEG_START_TIME", cmdFrag.Start)
	shell.ExitCode = 0
	shell.SetVar("FFMPEG_DURATION", cmdFrag.Duration)
	shell.ExitCode = 0
	shell.SetVar("FFMPEG_OUTPUT_FILE", cmdFrag.Output)

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
			arguments = append(arguments, `'`+shell.ReadVar("EXECUTOR")+`'`) //`../ci/segment.sh`
			arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_INPUT_FILE")+`'`)
			arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_START_TIME")+`'`)
			arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_DURATION")+`'`)
			arguments = append(arguments, `'`+shell.ReadVar("FFMPEG_OUTPUT_FILE")+`'`)
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
