'use client'
import React from 'react';
import { toast } from 'react-toastify';
import { Button } from '@/components/ui/button';
import StatusBadge from '@/components/StatusBadge';
import { formatJaDateSpan } from '@/lib/date/format';
import { SYNC_STATUS_COLORS, SYNC_STATUS_LABELS } from '@/features/events/status';
import type { EventDraftDetail, EventProposedDate } from '@/features/events/types';
import { CalendarCheck, Copy } from 'lucide-react';
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
    <li className="py-2.5">
        <div className="flex flex-wrap items-center justify-between gap-x-3 gap-y-2">
            <div className="flex min-w-0 items-center gap-3">
                <span
                    aria-label={`第${index + 1}候補`}
                    className="grid size-8 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary"
                >
                    {index + 1}
                </span>
                <span className="text-sm font-medium text-foreground">
                    {formatJaDateSpan(date.start, date.end)}
                </span>
            </div>
            <div className="flex shrink-0 items-center gap-3">
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
            <p className="mt-1 pl-11 text-xs text-destructive">同期エラー: {date.last_sync_error}</p>
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
        <section className="border-t border-border pt-6">
            <div className="flex flex-wrap items-center justify-between gap-2">
                <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">候補日程</h2>
                <div className="flex shrink-0 items-center gap-2">
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
            </div>
            {!isConfirmed && dates.length > 0 && (
                <p className="mt-1 text-sm text-muted-foreground">
                    候補を相手に共有し、決まったら日程を確定します。
                </p>
            )}
            <div className="mt-4">
                {isConfirmed && (
                    <div className="mb-4 flex items-center gap-4 rounded-lg border border-green-200 bg-green-50 p-4">
                        <span className="grid size-10 shrink-0 place-items-center rounded-full bg-green-100 text-green-600">
                            <CalendarCheck className="size-5" />
                        </span>
                        <div className="min-w-0">
                            <p className="text-sm font-medium text-green-700">確定日時</p>
                            <p className="mt-0.5 text-lg font-bold leading-snug text-gray-900">
                                {formatJaDateSpan(confirmedDate.start, confirmedDate.end)}
                            </p>
                        </div>
                    </div>
                )}
                {dates.length === 0 ? (
                    <p className="text-sm text-muted-foreground">
                        候補日程がありません。編集画面から追加してください。
                    </p>
                ) : isConfirmed ? (
                    otherDates.length > 0 && (
                        <details>
                            <summary className="cursor-pointer text-sm text-muted-foreground hover:text-foreground">
                                他の候補日程を表示({otherDates.length}件)
                            </summary>
                            <ul className="mt-1">
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
                    <ul>
                        {dates.map((date, index) => (
                            <ProposedDateRow key={date.id} date={date} index={index} />
                        ))}
                    </ul>
                )}
            </div>
        </section>
    );
};

export default ProposedDatesSection;
