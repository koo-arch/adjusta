import type { EventStatus, ProposedDateStatus, SyncStatus } from './types';

export type StatusColor = 'blue' | 'yellow' | 'green' | 'red' | 'gray';

export const EVENT_STATUS_LABELS: Record<EventStatus, string> = {
    draft: '下書き',
    active: '調整中',
    confirmed: '確定',
    cancelled: 'キャンセル',
};

// DESIGN.md「Status(ドメインステータス色)」に対応する StatusBadge の色名
export const EVENT_STATUS_COLORS: Record<EventStatus, StatusColor> = {
    draft: 'blue',
    active: 'yellow',
    confirmed: 'green',
    cancelled: 'red',
};

export const PROPOSED_DATE_STATUS_LABELS: Record<ProposedDateStatus, string> = {
    active: '調整中',
    confirmed: '確定',
    not_selected: '非選択',
    cancelled: 'キャンセル',
};

export const PROPOSED_DATE_STATUS_COLORS: Record<ProposedDateStatus, StatusColor> = {
    active: 'yellow',
    confirmed: 'green',
    not_selected: 'gray',
    cancelled: 'red',
};

export const SYNC_STATUS_LABELS: Record<SyncStatus, string> = {
    not_synced: '未同期',
    pending_sync: '同期待ち',
    synced: '同期済み',
    sync_failed: '同期失敗',
};

export const SYNC_STATUS_COLORS: Record<SyncStatus, StatusColor> = {
    not_synced: 'gray',
    pending_sync: 'yellow',
    synced: 'green',
    sync_failed: 'red',
};
