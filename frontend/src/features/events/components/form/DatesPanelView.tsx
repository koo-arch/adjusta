'use client'
import React from 'react';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';
import DraggableDateList from '@/features/events/components/form/DraggableDateList';
import AddDateDialog from '@/features/events/components/form/AddDateDialog';

interface DatesPanelViewProps<T extends SelectedDate | ProposedDate> {
    title: string;
    location: string;
    dates: T[];
    onDatesChange: React.Dispatch<React.SetStateAction<T[]>>;
    onAdd: (date: { start: Date; end: Date }) => void;
    onEditBasic: () => void;
    error?: string;
}

// 候補日程ステップの表示専用パネル(基本情報サマリー + 候補リスト)
const DatesPanelView = <T extends SelectedDate | ProposedDate>({
    title,
    location,
    dates,
    onDatesChange,
    onAdd,
    onEditBasic,
    error,
}: DatesPanelViewProps<T>) => {
    return (
        <div className="space-y-4">
            {/* 確認ステップの代替: 送信前に基本情報を見渡せるサマリー */}
            <div className="rounded-md bg-muted/60 px-3 py-2">
                <div className="flex items-start justify-between gap-2">
                    <div className="min-w-0 text-sm">
                        <p className="truncate font-medium text-foreground">
                            {title || <span className="text-muted-foreground">タイトル未入力</span>}
                        </p>
                        {location && <p className="truncate text-muted-foreground">{location}</p>}
                    </div>
                    <button
                        type="button"
                        onClick={onEditBasic}
                        className="shrink-0 text-sm text-primary transition-colors hover:text-primary-dark"
                    >
                        編集
                    </button>
                </div>
            </div>
            <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                    <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">候補日程</h2>
                    <p className="mt-1 text-sm text-muted-foreground">
                        カレンダーから選ぶか、「日時を追加」で登録します
                    </p>
                </div>
                <AddDateDialog onAdd={onAdd} />
            </div>
            {error && <p className="text-sm text-destructive">{error}</p>}
            {dates.length > 0 ? (
                <DraggableDateList dates={dates} onDatesChange={onDatesChange} />
            ) : (
                <div className="rounded-md border border-dashed border-input py-10 text-center text-sm text-muted-foreground">
                    候補日程がまだありません。
                    <br />
                    カレンダーで範囲を選択するか、「日時を追加」から登録してください。
                </div>
            )}
        </div>
    );
};

export default DatesPanelView;
