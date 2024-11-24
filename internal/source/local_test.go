package source

import (
	"os"
	"testing"

	"github.com/jonashiltl/openchangelog/internal"
	"github.com/jonashiltl/openchangelog/internal/store"
	"github.com/jonashiltl/openchangelog/internal/xcache"
)

func createLocalSource(t *testing.T) (*localSource, *os.File) {
	cache := xcache.NewMemoryCache()
	file, err := os.CreateTemp("", "release-*.md")
	if err != nil {
		t.Error(err)
	}

	source := NewLocalSourceFromStore(store.LocalSource{
		Path: file.Name(),
	}, cache).(*localSource)
	return source, file
}

func createLocalSourceDir(t *testing.T) (*localSource, string) {
	cache := xcache.NewMemoryCache()
	dir, err := os.MkdirTemp("", "notes-*")
	if err != nil {
		t.Error(err)
	}

	source := NewLocalSourceFromStore(store.LocalSource{
		Path: dir,
	}, cache).(*localSource)
	return source, dir
}

func createNTempFiles(t *testing.T, n int, dir string) {
	for range n {
		_, err := os.CreateTemp(dir, "*.md")
		if err != nil {
			t.Error(err)
		}
	}
}

func TestOpenAndCacheFile(t *testing.T) {
	source, file := createLocalSource(t)
	defer os.Remove(file.Name())
	file.Write([]byte("first line"))

	raw, err := source.openAndCacheFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	if raw.hasChanged == false {
		t.Error("expected hasChanged to be true initially")
	}

	raw, err = source.openAndCacheFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	if raw.hasChanged == true {
		t.Error("expected hasChanged to be false on reopen")
	}

	file.Write([]byte("\nsecond line"))
	raw, err = source.openAndCacheFile(file.Name())
	if err != nil {
		t.Error(err)
	}
	if raw.hasChanged == false {
		t.Error("expected hasChanged to be true after new line was added")
	}
}

func TestOpenAndCacheFileNotExists(t *testing.T) {
	source, file := createLocalSource(t)
	defer os.Remove(file.Name())
	_, err := source.openAndCacheFile(file.Name() + "agagjafja")
	if err == nil {
		t.Error("expected error to be returned")
	}
}

func TestLoadDir(t *testing.T) {
	source, dir := createLocalSourceDir(t)
	defer os.RemoveAll(dir)
	createNTempFiles(t, 2, dir)

	loaded, err := source.loadDir(dir, internal.NoPagination())
	if err != nil {
		t.Error(err)
	}

	if loaded.HasMore == true {
		t.Error("expected hasMore to be false")
	}

	if len(loaded.Raw) != 2 {
		t.Errorf("expected 2 raw notes, but got %d", len(loaded.Raw))
	}
}

func TestLoadDirPagination(t *testing.T) {
	source, dir := createLocalSourceDir(t)
	defer os.RemoveAll(dir)

	createNTempFiles(t, 10, dir)

	tests := []struct {
		pageSize        int
		page            int
		expectedHasMore bool
		expectedLen     int
	}{
		{
			pageSize:        9,
			page:            1,
			expectedHasMore: true,
			expectedLen:     9,
		},
		{
			pageSize:        10,
			page:            1,
			expectedHasMore: false,
			expectedLen:     10,
		},
		{
			pageSize:        5,
			page:            2,
			expectedHasMore: false,
			expectedLen:     5,
		},
		{
			pageSize:        3,
			page:            4,
			expectedHasMore: false,
			expectedLen:     1,
		},
	}

	for _, test := range tests {
		loaded, err := source.loadDir(dir, internal.NewPagination(test.pageSize, test.page))
		if err != nil {
			t.Error(err)
		}

		if loaded.HasMore != test.expectedHasMore {
			t.Errorf("expected hasMore to be %t, but got %t", test.expectedHasMore, loaded.HasMore)
		}
		if len(loaded.Raw) != test.expectedLen {
			t.Errorf("expected %d release notes, but got %d", test.expectedLen, len(loaded.Raw))
		}
	}
}

func TestCalculatePaginationIndices(t *testing.T) {
	tests := []struct {
		page          internal.Pagination
		totalFiles    int
		expectedStart int
		expectedEnd   int
	}{
		{
			page:          internal.NoPagination(),
			totalFiles:    10,
			expectedStart: 0,
			expectedEnd:   10,
		},
		{
			page:          internal.NewPagination(10, 1),
			totalFiles:    20,
			expectedStart: 0,
			expectedEnd:   10,
		},
		{
			page:          internal.NewPagination(10, 2),
			totalFiles:    20,
			expectedStart: 10,
			expectedEnd:   20,
		},
		{
			page:          internal.NewPagination(10, 2),
			totalFiles:    15,
			expectedStart: 10,
			expectedEnd:   15,
		},
	}

	for _, test := range tests {
		start, end := calculatePaginationIndices(test.page, test.totalFiles)
		if start != test.expectedStart {
			t.Errorf("expected start to be %d, but got %d", test.expectedStart, start)
		}
		if end != test.expectedEnd {
			t.Errorf("expected end to be %d, but got %d", test.expectedEnd, end)
		}
	}
}
