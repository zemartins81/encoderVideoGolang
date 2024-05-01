package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zemartins81/encoderVideoGolang/framework/utils"
)

func TestIsJson(t *testing.T) {
	json := `{
                "id=": "525b5fd9-700d-4feb-89c0-415a1e6e148c",
                "file_path": "emilly.mp4",
                "status": "pending"
            }`
	err := utils.IsJson(json)
	require.Nil(t, err)

	json = `wes`
	err = utils.IsJson(json)
	require.Error(t, err)
}
