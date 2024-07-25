//go:build integration

package tests

import (
	"context"
	"github.com/radium-rtf/coderunner_lib/internal/tests/files"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPython(t *testing.T) {
	t.Parallel()

	type test struct {
		cmd        string
		outputFile string
		name       string
		t          string
		dataName   string
		code       int64
	}

	ctx := context.Background()
	tests := []test{
		{

			cmd:        "cat input.txt | python3 main.py",
			name:       "input is ok",
			outputFile: "output.txt",
			code:       0,
			t:          "python",
			dataName:   "input",
		},
		{
			cmd:        "python3 main.py",
			name:       "output is ok",
			outputFile: "output.txt",
			code:       0,
			t:          "python",
			dataName:   "output",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dirFiles, err := files.Parse(tt.t, tt.dataName)
			require.NoError(t, err)

			output := string(files.Find(dirFiles, tt.outputFile).Content.GetBytes())

			sandbox, err := runner.NewSandbox(ctx, tt.cmd, pythonProfile, limits, dirFiles)
			require.NoError(t, err)
			defer sandbox.Close()

			err = sandbox.Start()
			require.NoError(t, err)

			statusCode, err := sandbox.Wait()
			require.NoError(t, err)
			require.Equal(t, tt.code, statusCode)

			stdout, err := sandbox.ShowStdout()
			require.NoError(t, err)
			require.Equal(t, output, stdout)

			std, err := sandbox.ShowStd()
			require.NoError(t, err)
			require.Equal(t, output, std.StdOut)
			require.Equal(t, "", std.StdErr)
		})
	}
}
