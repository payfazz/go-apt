package fazzconfig

import (
	"bufio"
	"github.com/payfazz/go-apt/pkg/fazzconfig/configsource"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestConfigReader_Get(t *testing.T) {
	err := os.Setenv("B", "Y1")
	require.NoError(t, err)
	err = os.Setenv("C", "Y2")
	require.NoError(t, err)

	path := "./.env"
	err = createTestFile(path)
	require.NoError(t, err)
	defer os.Remove(path)

	reader := NewReader(
		configsource.FromEnvFile(path),
		configsource.FromEnv(),
		configsource.FromMap(map[string]string{
			"A": "Z1", "B": "Z2", "C": "Z3", "D": "Z4",
		}),
	)

	require.Equal(t, "X1", reader.Get("A"))
	require.Equal(t, "X2", reader.Get("B"))
	require.Equal(t, "Y2", reader.Get("C"))
	require.Equal(t, "Z4", reader.Get("D"))
	require.Empty(t, reader.Get("E"))
}

func createTestFile(path string) error {
	file, err := os.Create(path)
	w := bufio.NewWriter(file)
	_, err = w.WriteString("A=X1\nB=X2\n")
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	_ = file.Close()
	return nil
}
