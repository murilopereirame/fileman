package config

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestParseValidConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configHandler := ConfigHandler{
		config: "testdata/config_valid.json",
	}

	config, err := configHandler.Load()

	assert.Equal(t, nil, err)
	assert.Equal(t, "* * * * *", config.Cron)
	assert.Equal(t, []WatchedDirectory{
		{
			Path: "foo/bar",
			Age:  1.5,
		}, {
			Path: "bar/foo",
			Age:  2.0,
		},
	}, config.WatchedDirectories)
}

func TestParseInvalidConfig(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	configHandler := ConfigHandler{
		config: "",
	}

	config, err := configHandler.Load()

	assert.Equal(t, config.Cron, "")
	assert.Equal(t, []WatchedDirectory(nil), config.WatchedDirectories)
	assert.Error(t, err)
}
