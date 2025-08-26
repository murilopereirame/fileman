package handler

import (
	"errors"
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
		nil,
	}

	fileHandler := FileHandler{
		clock: mockClock,
	}

	files, err := fileHandler.ListFiles(mockFS, "foo/bar")
	assert.Equal(t, 1, files.Len())
	assert.Equal(t, mockedResult, files.Front().Value)
	assert.Equal(t, nil, err)
}

func TestHandleListFilesReadDirError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClock := mocks.NewMockClock(ctrl)

	mockError := errors.New("foo")
	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{}, mockError).Times(1)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	files, err := fileHandler.ListFiles(mockFS, "foo/bar")
	assert.Equal(t, 0, files.Len())
	assert.Equal(t, mockError, err)
}

func TestListFilesFileHasError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClock := mocks.NewMockClock(ctrl)

	mockError := errors.New("error")

	mockEntry := mocks.NewMockDirEntry(ctrl)
	mockEntry.EXPECT().Name().Return("file1.txt").Times(1)
	mockEntry.EXPECT().Info().Return(nil, mockError)

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{
		mockEntry,
	}, nil).Times(1)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	files, _ := fileHandler.ListFiles(mockFS, "foo/bar")
	file, _ := files.Front().Value.(*File)

	assert.Equal(t, 1, files.Len())
	assert.Equal(t, mockError, file.error)
	assert.Equal(t, "file1.txt", file.name)
	assert.Equal(t, "", file.path)
	assert.Equal(t, false, file.isDir)
	assert.Equal(t, int64(0), file.createdAt)
	assert.Equal(t, float64(0), file.age)
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

	result, errs := fileHandler.DeleteOldFiles(mockFS, "foo/bar", 7)
	assert.Equal(t, 0, len(errs))
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "foo/bar/file1.txt", result[0])
}

func TestDeleteOldFileHandlesListingError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClock := mocks.NewMockClock(ctrl)

	mockError := errors.New("foo")

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{}, mockError).Times(1)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	result, errs := fileHandler.DeleteOldFiles(mockFS, "foo/bar", 7)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, mockError, errs[0])
	assert.Equal(t, 0, len(result))
}

func TestDeleteOldFileHandlesFileWithError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockError := errors.New("error")
	fileToBeDeletedCreatedAt := time.Unix(1755561600, 0)

	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().CalculateAge(int64(1755561600)).Return(7.1).Times(1)

	mockFileInfoToBeDeleted := mocks.NewMockFileInfo(ctrl)
	mockFileInfoToBeDeleted.EXPECT().ModTime().Return(fileToBeDeletedCreatedAt).Times(2)
	mockFileInfoToBeDeleted.EXPECT().IsDir().Return(false).Times(1)

	mockEntryToBeDeleted := mocks.NewMockDirEntry(ctrl)
	mockEntryToBeDeleted.EXPECT().Name().Return("file1.txt").Times(2)
	mockEntryToBeDeleted.EXPECT().Info().Return(mockFileInfoToBeDeleted, nil)

	mockEntryToBeAlsoDeletedWithError := mocks.NewMockDirEntry(ctrl)
	mockEntryToBeAlsoDeletedWithError.EXPECT().Name().Return("file2.txt").Times(1)
	mockEntryToBeAlsoDeletedWithError.EXPECT().Info().Return(nil, mockError)

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{
		mockEntryToBeDeleted,
		mockEntryToBeAlsoDeletedWithError,
	}, nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file1.txt").Return(nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file2.txt").Times(0)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	result, errs := fileHandler.DeleteOldFiles(mockFS, "foo/bar", 7)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, mockError, errs[0])
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "foo/bar/file1.txt", result[0])
}

func TestDeleteOldFileFailsOnDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockError := errors.New("error")

	mockClock := mocks.NewMockClock(ctrl)
	mockClock.EXPECT().CalculateAge(int64(1755561600)).Return(7.1).Times(1)
	mockClock.EXPECT().CalculateAge(int64(1755475200)).Return(8.0).Times(1)

	fileToBeDeletedCreatedAt := time.Unix(1755561600, 0)
	fileToBeKeptCreatedAt := time.Unix(1755475200, 0)

	mockFileInfoToBeDeleted := mocks.NewMockFileInfo(ctrl)
	mockEntryToBeDeleted := mocks.NewMockDirEntry(ctrl)

	mockFileInfoToBeDeleted.EXPECT().ModTime().Return(fileToBeDeletedCreatedAt).Times(2)
	mockFileInfoToBeDeleted.EXPECT().IsDir().Return(false)
	mockEntryToBeDeleted.EXPECT().Name().Return("file1.txt").Times(2)
	mockEntryToBeDeleted.EXPECT().Info().Return(mockFileInfoToBeDeleted, nil)

	mockFileInfoToBeDeletedWithError := mocks.NewMockFileInfo(ctrl)
	mockEntryToBeDeletedWithError := mocks.NewMockDirEntry(ctrl)

	mockFileInfoToBeDeletedWithError.EXPECT().ModTime().Return(fileToBeKeptCreatedAt).Times(2)
	mockFileInfoToBeDeletedWithError.EXPECT().IsDir().Return(false)
	mockEntryToBeDeletedWithError.EXPECT().Name().Return("file2.txt").Times(2)
	mockEntryToBeDeletedWithError.EXPECT().Info().Return(mockFileInfoToBeDeletedWithError, nil)

	mockFS := mocks.NewMockFileSystem(ctrl)
	mockFS.EXPECT().ReadDir("foo/bar").Return([]fs.DirEntry{
		mockEntryToBeDeleted,
		mockEntryToBeDeletedWithError,
	}, nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file1.txt").Return(nil).Times(1)
	mockFS.EXPECT().DeleteFile("foo/bar/file2.txt").Return(mockError).Times(1)

	fileHandler := FileHandler{
		clock: mockClock,
	}

	result, errs := fileHandler.DeleteOldFiles(mockFS, "foo/bar", 7)
	assert.Equal(t, 1, len(errs))
	assert.Equal(t, mockError, errs[0])
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "foo/bar/file1.txt", result[0])
}
