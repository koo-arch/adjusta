'use client'
import React from 'react';
import { toast } from 'react-toastify';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import StatusBadge from '@/components/StatusBadge';
import { formatJaDateSpan } from '@/lib/date/format';
import {
    PROPOSED_DATE_STATUS_COLORS,
    PROPOSED_DATE_STATUS_LABELS,
    SYNC_STATUS_COLORS,
    SYNC_STATUS_LABELS,
} from '@/features/events/status';
import type { EventDraftDetail, EventProposedDate } from '@/features/events/types';
import { Copy } from 'lucide-react';
import ConfirmButton from './ConfirmButton';

interface ProposedDatesSectionProps {
    eventID: string;
    detail: EventDraftDetail;
}

const sortByPriority = (dates: EventProposedDate[]) =>
    [...dates].sort((a, b) => a.priority - b.priority);

// 相手にそのまま貼れるテキスト(要件 5.6: 初期リリースの共有主導線)
const buildCopyText = (title: string, dates: EventProposedDate[]) =>
    [
        `【${title}】候補日程`,
        ...dates.map((date, index) => `第${index + 1}候補: ${formatJaDateSpan(date.start, date.end)}`),
    ].join('\n');

const ProposedDateRow: React.FC<{ date: EventProposedDate; index: number }> = ({ date, index }) => (
    <li className="space-y-1 py-3">
        <div className="flex flex-wrap items-center justify-between gap-2">
            <div className="flex min-w-0 items-center gap-3">
                <span className="shrink-0 text-sm font-medium text-muted-foreground">第{index + 1}候補</span>
                <span className="text-sm text-foreground">{formatJaDateSpan(date.start, date.end)}</span>
            </div>
            <div className="flex shrink-0 items-center gap-3">
                <StatusBadge
                    label={PROPOSED_DATE_STATUS_LABELS[date.status]}
                    circleColor={PROPOSED_DATE_STATUS_COLORS[date.status]}
                    textColor={PROPOSED_DATE_STATUS_COLORS[date.status]}
                    textSize="sm"
                    circleSize="sm"
                />
                <StatusBadge
                    label={SYNC_STATUS_LABELS[date.sync_status]}
                    circleColor={SYNC_STATUS_COLORS[date.sync_status]}
                    textColor={SYNC_STATUS_COLORS[date.sync_status]}
                    textSize="sm"
                    circleSize="sm"
                />
            </div>
        </div>
        {date.sync_status === 'sync_failed' && date.last_sync_error && (
            <p className="text-xs text-destructive">同期エラー: {date.last_sync_error}</p>
        )}
    </li>
);

const ProposedDatesSection: React.FC<ProposedDatesSectionProps> = ({ eventID, detail }) => {
    const dates = sortByPriority(detail.proposed_dates ?? []);
    const confirmedDate = dates.find((date) => date.id === detail.confirmed_date_id);
    const isConfirmed = detail.status === 'confirmed' && !!confirmedDate;
    const otherDates = isConfirmed ? dates.filter((date) => date.id !== confirmedDate.id) : dates;

    const handleCopy = async () => {
        try {
            await navigator.clipboard.writeText(buildCopyText(detail.title, dates));
            toast.success('候補日程をコピーしました');
        } catch {
            toast.error('コピーに失敗しました');
        }
    };

    return (
        <Card>
            <CardHeader className="flex-row flex-wrap items-start justify-between gap-2 space-y-0">
                <div className="space-y-1.5">
                    <CardTitle>候補日程</CardTitle>
                    {!isConfirmed && (
                        <CardDescription>候補を相手に共有し、決まったら日程を確定します。</CardDescription>
                    )}
                </div>
                <div className="flex shrink-0 flex-wrap items-center gap-2">
                    {dates.length > 0 && (
                        <Button
                            variant="ghost"
                            size="icon"
                            aria-label="候補日程をコピー"
                            title="候補日程をコピー"
                            className="text-muted-foreground hover:text-foreground [&_svg]:size-5"
                            onClick={handleCopy}
                        >
                            <Copy />
                        </Button>
                    )}
                    {detail.status !== 'cancelled' && (
                        <ConfirmButton eventID={eventID} detail={detail} isConfirmed={isConfirmed} />
                    )}
                </div>
            </CardHeader>
            <CardContent>
                {isConfirmed && (
                    <div className="mb-4 rounded-lg border border-green-200 bg-green-50 p-4">
                        <p className="text-sm font-medium text-green-700">確定日時</p>
                        <p className="mt-1 text-lg font-bold leading-snug text-gray-900">
                            {formatJaDateSpan(confirmedDate.start, confirmedDate.end)}
                        </p>
                    </div>
                )}
                {dates.length === 0 ? (
                    <p className="py-4 text-center text-sm text-muted-foreground">
                        候補日程がありません。編集画面から追加してください。
                    </p>
                ) : isConfirmed ? (
                    otherDates.length > 0 && (
                        <details>
                            <summary className="cursor-pointer text-sm text-muted-foreground hover:text-foreground">
                                他の候補日程を表示({otherDates.length}件)
                            </summary>
                            <ul className="mt-1 divide-y divide-border">
                                {otherDates.map((date) => (
                                    <ProposedDateRow
                                        key={date.id}
                                        date={date}
                                        index={dates.findIndex((d) => d.id === date.id)}
                                    />
                                ))}
                            </ul>
                        </details>
                    )
                ) : (
                    <ul className="divide-y divide-border">
                        {dates.map((date, index) => (
                            <ProposedDateRow key={date.id} date={date} index={index} />
                        ))}
                    </ul>
                )}
            </CardContent>
        </Card>
    );
};

export default ProposedDatesSection;
