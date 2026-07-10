import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import StatusBadge from '@/components/StatusBadge';
import { formatJaDate } from '@/lib/date/format';
import { SYNC_STATUS_COLORS, SYNC_STATUS_LABELS } from '@/features/events/status';
import type { EventDraftDetail } from '@/features/events/types';
import { MapPin } from 'lucide-react';

interface EventInfoSectionProps {
    detail: EventDraftDetail;
}

const formatSyncedAt = (value?: string) => {
    if (!value) {
        return null;
    }
    const parsed = new Date(value);
    return Number.isNaN(parsed.getTime()) ? null : formatJaDate(parsed);
};

const EventInfoSection: React.FC<EventInfoSectionProps> = ({ detail }) => {
    const syncedAt = formatSyncedAt(detail.last_synced_at);

    return (
        <Card>
            <CardHeader>
                <CardTitle>基本情報</CardTitle>
            </CardHeader>
            <CardContent>
                <dl className="space-y-4">
                    <div>
                        <dt className="text-sm font-medium text-muted-foreground">場所</dt>
                        <dd className="mt-1 flex items-center gap-1.5 text-sm text-foreground">
                            <MapPin className="size-4 shrink-0 text-muted-foreground" />
                            {detail.location || '未設定'}
                        </dd>
                    </div>
                    <div>
                        <dt className="text-sm font-medium text-muted-foreground">説明</dt>
                        <dd className="mt-1 whitespace-pre-wrap break-words text-sm text-foreground">
                            {detail.description || '未設定'}
                        </dd>
                    </div>
                    <div>
                        <dt className="text-sm font-medium text-muted-foreground">Google Calendar 同期</dt>
                        <dd className="mt-1 space-y-1">
                            <div className="flex flex-wrap items-center gap-3">
                                <StatusBadge
                                    label={SYNC_STATUS_LABELS[detail.sync_status]}
                                    circleColor={SYNC_STATUS_COLORS[detail.sync_status]}
                                    textColor={SYNC_STATUS_COLORS[detail.sync_status]}
                                    textSize="sm"
                                    circleSize="sm"
                                />
                                {syncedAt && (
                                    <span className="text-xs text-muted-foreground">最終同期: {syncedAt}</span>
                                )}
                            </div>
                            {detail.sync_status === 'sync_failed' && detail.last_sync_error && (
                                <p className="text-xs text-destructive">同期エラー: {detail.last_sync_error}</p>
                            )}
                        </dd>
                    </div>
                </dl>
            </CardContent>
        </Card>
    );
};

export default EventInfoSection;
