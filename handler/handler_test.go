package handler

import (
	"fileman/mocks"
	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"
	"io/fs"
	"testing"
	"time"
)

func TestListFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().CalculateAge(gomock.Any()).Return(3.0)

	fileCreatedAt := time.Unix(1755907200, 0)
	mockFileInfo := mocks.NewMockFileInfo(ctrl)
	mockFileInfo.EXPECT().ModTime().Return(fileCreatedAt).Times(2)
	mockFileInfo.EXPECT().IsDir().Return(false)

	mockEntry := mocks.NewMockDirEntry(ctrl)
	mockEntry.EXPECT().Name().Return("file1.txt").Times(2)
	mockEntry.EXPECT().Info().Return(mockFileInfo, nil)

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{
		mockEntry,
	}, nil).Times(1)

	mockedResult := &File{
		1755907200,
		3.0,
		"file1.txt",
		"foo/bar/file1.txt",
		false,
		[]string{},
	}

	fileHandler := FileHandler{
		clock: mockClock,
	}

	files := fileHandler.ListFiles(mockFS, "foo/bar")
	assert.Equal(t, files.Len(), 1)
	assert.Equal(t, files.Front().Value, mockedResult)
}
