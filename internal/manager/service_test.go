package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewManager(t *testing.T) {
	srcRootDir, err := os.MkdirTemp("", "managerSrcDir_*")
	require.NoError(t, err)
	t.Log("source root dir: ", srcRootDir)
	targetRootDir, err := os.MkdirTemp("", "managerTargetDir_*")
	require.NoError(t, err)
	t.Log("target root dir: ", targetRootDir)
	writer := new(strings.Builder)
	reader := new(strings.Reader)

	type testCase struct {
		title          string
		input          ServiceInitInput
		expectedErrStr string
	}
	testCases := []testCase{
		{
			title: "valid input=>success expected",
			input: ServiceInitInput{
				SourceRootDir:          srcRootDir,
				TargetRootDir:          targetRootDir,
				ListingDirPathsWriter:  writer,
				ListingFilePathsWriter: writer,
				ListingErrorsLogWriter: writer,
				FilePathsReader:        reader,
				FileCopyPipelineLength: 2,
				FileBackupLogWriter:    writer,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			m := NewService(tc.input)
			assert.NoError(t, err)
			assert.NotNil(t, m)
			assertedService, ok := m.(*service)
			assert.True(t, ok)
			assert.Len(t, assertedService.FileBackupWorkers, tc.input.FileCopyPipelineLength)
			assert.Equal(t, cap(assertedService.TasksPipeline), tc.input.FileCopyPipelineLength)
		})
	}
}

func TestListSources(t *testing.T) {
	srcRootDir, err := os.MkdirTemp("", "msgListSrcDir_*")
	t.Log("source root dir: ", srcRootDir)
	require.NoError(t, err)
	type setupFunc func(t *testing.T, dirPath string) string

	testCases := []struct {
		title             string
		testDirName       string
		setupDirFunc     setupFunc
		expectedDirPaths  func(string) string
		expectedFilePaths string
		expectedErrorsLog string
	}{
		{
			title:       "empty source=>expect empty output",
			testDirName: "emptySrc",
			setupDirFunc: func(t *testing.T, dirName string) string {
				fullPath := filepath.Join(srcRootDir, dirName)
				t.Log("srcPath: ", fullPath)
				err = os.MkdirAll(fullPath, 0755)
				t.Log("err: ", err)
				require.NoError(t, err)
				return fullPath
			},
			expectedDirPaths:  func(testRoot string) string {return fmt.Sprintf("%s\n", testRoot)},
			expectedFilePaths: "",
			expectedErrorsLog: "",
		}, {
			title:       "single dir",
			testDirName: "singleDir",
			setupDirFunc: func(t *testing.T, testDirName string) string {
				fullPath := filepath.Join(srcRootDir, testDirName)
				t.Log("srcPath: ", fullPath)
				err = os.MkdirAll(fullPath, 0755)
				require.NoError(t, err)
				pathOne := filepath.Join(fullPath, "one")
				err = os.MkdirAll(pathOne, 0755)
				require.NoError(t, err)
				return fullPath
			},
			expectedDirPaths:  func(testRoot string) string {
				return fmt.Sprintf("%s\n%s" + string(filepath.Separator) + "%s\n", testRoot, testRoot, "one")
			},
			expectedFilePaths: "",
			expectedErrorsLog: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			var dirPaths = &strings.Builder{}
			var filePaths = &strings.Builder{}
			var errorsLog = &strings.Builder{}
			testRoot := tc.setupDirFunc(t, tc.testDirName)
			api := NewService(ServiceInitInput{
				SourceRootDir:          testRoot,
				ListingDirPathsWriter:  dirPaths,
				ListingFilePathsWriter: filePaths,
				ListingErrorsLogWriter: errorsLog,
			})

			err = api.ListSources()

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedDirPaths(testRoot), dirPaths.String())
			assert.Equal(t, tc.expectedFilePaths, filePaths.String())
			assert.Equal(t, tc.expectedErrorsLog, errorsLog.String())
		})
	}
}
