import type { EventStatus } from './types';

export const EVENT_STATUS_LABELS: Record<EventStatus, string> = {
    draft: '下書き',
    active: '調整中',
    confirmed: '確定',
    cancelled: 'キャンセル',
};

// DESIGN.md「Status(ドメインステータス色)」に対応する StatusBadge の色名
export const EVENT_STATUS_COLORS: Record<EventStatus, 'blue' | 'yellow' | 'green' | 'red'> = {
    draft: 'blue',
    active: 'yellow',
    confirmed: 'green',
    cancelled: 'red',
};
