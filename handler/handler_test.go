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

func TestDeleteOldFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().CalculateAge(int64(1755561600)).Return(7.1).Times(1)
	mockClock.EXPECT().CalculateAge(int64(1755907200)).Return(3.0).Times(1)

	fileToBeDeletedCreatedAt := time.Unix(1755561600, 0)
	fileToBeKeptCreatedAt := time.Unix(1755907200, 0)

	mockFileInfoToBeDeleted := mocks.NewMockFileInfo(ctrl)
	mockEntryToBeDeleted := mocks.NewMockDirEntry(ctrl)

	mockFileInfoToBeDeleted.EXPECT().ModTime().Return(fileToBeDeletedCreatedAt).Times(2)
	mockFileInfoToBeDeleted.EXPECT().IsDir().Return(false)
	mockEntryToBeDeleted.EXPECT().Name().Return("file1.txt").Times(2)
	mockEntryToBeDeleted.EXPECT().Info().Return(mockFileInfoToBeDeleted, nil)

	mockFileInfoToBeKept := mocks.NewMockFileInfo(ctrl)
	mockEntryToBeKept := mocks.NewMockDirEntry(ctrl)

	mockFileInfoToBeKept.EXPECT().ModTime().Return(fileToBeKeptCreatedAt).Times(2)
	mockFileInfoToBeKept.EXPECT().IsDir().Return(false)
	mockEntryToBeKept.EXPECT().Name().Return("file2.txt").Times(2)
	mockEntryToBeKept.EXPECT().Info().Return(mockFileInfoToBeKept, nil)

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{
		mockEntryToBeDeleted,
		mockEntryToBeKept,
	}, nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file1.txt").Return(nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file2.txt").Times(0)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	result := fileHandler.DeleteOldFiles(mockFS, "foo/bar", 7)
	assert.Equal(t, len(result), 0)
}
