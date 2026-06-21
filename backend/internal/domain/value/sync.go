package value

type SyncStatus string

const (
	SyncStatusNotSynced SyncStatus = "not_synced"
	SyncStatusPending   SyncStatus = "pending_sync"
	SyncStatusSynced    SyncStatus = "synced"
	SyncStatusFailed    SyncStatus = "sync_failed"
)
