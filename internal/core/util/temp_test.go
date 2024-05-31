package util_test

import (
	"io/fs"
	"myapp/internal/config"
	"myapp/internal/core/util"
	"myapp/internal/core/util/mock"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestNewTemper_DirectoryExists(t *testing.T) {
	virtualFS := mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/": &fstest.MapFile{Mode: fs.ModeDir},
			"testDir/": &fstest.MapFile{Mode: fs.ModeDir},
		},
	}

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}
	temper, err := util.NewTemper(tempConfig, &virtualFS)

	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if temper == nil {
		t.Fatal("expected temper to be non-nil")
	}
}

func TestNewTemper_DirectoryNotExist(t *testing.T) {
	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/": &fstest.MapFile{Mode: fs.ModeDir},
		},
	}

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}
	temper, err := util.NewTemper(tempConfig, virtualFS)

	if _, ok := virtualFS.Mfs["testDir"]; !ok {
		t.Fatalf("directory not create")
	}
	if err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if temper == nil {
		t.Fatal("expected temper to be non-nil")
	}
}

func TestNewTemper_MaxTempedID(t *testing.T) {
	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}

	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/":           &fstest.MapFile{Mode: fs.ModeDir},
			"testDir/temp1-file": &fstest.MapFile{Data: []byte("data")},
			"testDir/temp3-file": &fstest.MapFile{Data: []byte("data")},
			"testDir/temp2-file": &fstest.MapFile{Data: []byte("data")},
		},
	}

	temper, _ := util.NewTemper(tempConfig, virtualFS)

	if temper.MaxTempedID != 3 {
		t.Fatalf("expected '3', got %d", temper.MaxTempedID)
	}
}

func TestNewTemper_TempedIDs(t *testing.T) {

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}

	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/":            &fstest.MapFile{Mode: fs.ModeDir},
			"testDir/temp1-file":  &fstest.MapFile{Data: []byte("data")},
			"testDir/temp3-file":  &fstest.MapFile{Data: []byte("data")},
			"testDir/temp3-1file": &fstest.MapFile{Data: []byte("data")},
			"testDir/temp2-file":  &fstest.MapFile{Data: []byte("data")},
		},
	}

	temper, _ := util.NewTemper(tempConfig, virtualFS)

	expectedTempedIDs := map[int]bool{
		1: true,
		2: true,
		3: true,
	}

	if !reflect.DeepEqual(temper.TempedIDs, expectedTempedIDs) {
		t.Fatal("existed IDs is not defined correctly")
	}

}

func TestReadTempFile_Success(t *testing.T) {

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}

	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/":      &fstest.MapFile{Mode: fs.ModeDir},
			"testDir/temp1": &fstest.MapFile{Data: []byte("data")},
		},
	}

	temper, _ := util.NewTemper(tempConfig, virtualFS)

	data := temper.ReadTempFile("testDir/temp1")
	if string(data) != "data" {
		t.Fatalf("expected 'data', got %s", data)
	}
}

func TestReadTempFile_Error(t *testing.T) {
	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/": &fstest.MapFile{Mode: fs.ModeDir},
		},
	}

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}
	temper, _ := util.NewTemper(tempConfig, virtualFS)

	data := temper.ReadTempFile("testDir/temp1")
	if data != nil {
		t.Fatal("expected nil, got data")
	}
}

func TestSaveTempDataByID(t *testing.T) {
	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate": &fstest.MapFile{Mode: fs.ModeDir},
			"testDir": &fstest.MapFile{Mode: fs.ModeDir},
		},
	}

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}
	temper, _ := util.NewTemper(tempConfig, virtualFS)

	err := temper.SaveTempDataByID([]byte("data"), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	found := false
	for name := range virtualFS.Mfs {
		if strings.Contains(name, "testDir/temp1") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected temp file to be created")
	}
}

func TestRemoveTemp(t *testing.T) {

	expectedMapFS := fstest.MapFS{
		"migrate/": &fstest.MapFile{Mode: fs.ModeDir},
	}

	virtualFS := &mock.MockOSFileSystem{
		Mfs: fstest.MapFS{
			"migrate/":            &fstest.MapFile{Mode: fs.ModeDir},
			"testDir/temp1-file":  &fstest.MapFile{Data: []byte("data")},
			"testDir/temp3-file":  &fstest.MapFile{Data: []byte("data")},
			"testDir/temp3-1file": &fstest.MapFile{Data: []byte("data")},
			"testDir/temp2-file":  &fstest.MapFile{Data: []byte("data")},
		},
	}

	tempConfig := &config.TempConfig{TempDir: "testDir", TempFilePattern: "temp"}
	temper, _ := util.NewTemper(tempConfig, virtualFS)
	err := temper.RemoveTemp()
	if err != nil {
		t.Fatalf("err of remove temp %v", err)
	}

	if !reflect.DeepEqual(expectedMapFS, virtualFS.Mfs) {
		t.Fatal("not all files were deleted, or unnecessary ones were deleted")
	}
}
