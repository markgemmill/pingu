package pkg

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/spf13/afero"
	"path"
	"strings"
	"time"
)

const PASS = "PASS"
const FAIL = "FAIL"

// StoreMaster holds the current and historic StoreRecords.
type StoreMaster struct {
	Url      string        `json:"url"`
	StoreId  string        `json:"store-id"`
	Current  StoreRecord   `json:"current"`
	Failures []StoreRecord `json:"failures"`
	Passes   []StoreRecord `json:"passes"`
}

func NewStoreMaster(url, storeId string) *StoreMaster {
	return &StoreMaster{
		Url:     url,
		StoreId: storeId,
	}
}

// Store manages the State of the StoreMaster record.
type Store struct {
	Url        string
	Name       string
	Path       string
	MaxRecords int
	Data       *StoreMaster
}

func sluggifyUrl(url string) string {
	url = strings.Replace(url, "://", "_", -1)
	url = strings.Replace(url, ":", "_", -1)
	url = strings.Replace(url, "//", "_", -1)
	url = strings.Replace(url, ".", "_", -1)
	url = strings.Replace(url, "/", "_", -1)
	url = strings.Replace(url, "?", "_", -1)
	url = strings.Replace(url, "&", "_", -1)
	url = strings.Replace(url, "=", "_", -1)
	return url
}

func getStoreId(url, name string) string {
	if name == "" {
		return fmt.Sprintf("%x", sha1.Sum([]byte(url)))
	}
	return name
}

func NewStore(url, name string) *Store {
	storeId := getStoreId(url, name)
	return &Store{
		Url:  url,
		Name: storeId,
		Path: path.Join(dirs.UserDataDir(), fmt.Sprintf("pingu-%s-log.json", storeId)),
		Data: NewStoreMaster(url, storeId),
	}
}

func (s *Store) Read() {
	/// first check if the file exists, and create it if it doesn't
	exists, err := afero.Exists(fs, s.Path)
	PanicOnError(err)

	if !exists {
		s.Write()
		return
	}

	content, err := afero.ReadFile(fs, s.Path)
	PanicOnError(err)

	err = json.Unmarshal(content, s.Data)
	PanicOnError(err)
}

func (s *Store) Write() {
	content, err := json.Marshal(s.Data)
	PanicOnError(err)

	err = afero.WriteFile(fs, s.Path, content, 0777)
	PanicOnError(err)
}

// Save either updates the current StoreRecord value,
// if the current check is a PASS or pushes the current StoreRecord to
// Passes history if the current check is a FAIL, and then creates a new
// StoreRecord for the current check failure.
//
// We toggle records between changes in status.
// There are 2 status' PASS and FAIL.
// When there is a change, we stash the current status to history,
// and start a new status.
func (s *Store) Save(status, message string) {
	// TODO: on first save we get an empty store record - need to fix this...
	currentTimestamp := time.Now()
	if s.Data.Current.Status == status {
		s.Data.Current.Count += 1
		s.Data.Current.Interval = currentTimestamp.Sub(s.Data.Current.Last).Seconds()
		s.Data.Current.Last = currentTimestamp
		s.Data.Current.Message = message
	} else {
		if status == PASS {
			s.Data.Failures = append(s.Data.Failures, s.Data.Current)
		} else if status == FAIL {
			s.Data.Passes = append(s.Data.Passes, s.Data.Current)
		}
		s.Data.Current = StoreRecord{
			Start:    currentTimestamp,
			Last:     currentTimestamp,
			Interval: 0.0,
			Count:    1,
			Status:   status,
			Message:  message,
		}
	}
}
