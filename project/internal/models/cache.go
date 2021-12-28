package models

type CacheCommand string

const (
	CacheCommandRemove CacheCommand = "REMOVE"
	CacheCommandPurge  CacheCommand = "PURGE"
)

type CacheMsg struct {
	Command CacheCommand `json:"command"`
	Key     interface{}  `json:"key"`
}
