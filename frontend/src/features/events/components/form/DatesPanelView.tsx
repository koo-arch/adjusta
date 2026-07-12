'use client'
import React, { useState } from 'react';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';
import DraggableDateList from '@/features/events/components/form/DraggableDateList';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

interface DatesPanelViewProps<T extends SelectedDate | ProposedDate> {
    title: string;
    location: string;
    dates: T[];
    onDatesChange: React.Dispatch<React.SetStateAction<T[]>>;
    // 既定値で候補を 1 件追加し、その ID を返す(追加直後の行を編集モードで開く)
    onAdd: () => string;
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
    const [editingId, setEditingId] = useState<string | null>(null);

    const handleAdd = () => {
        setEditingId(onAdd());
    };

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
                        並び順がそのまま優先順位になります。
                        {/* 操作説明は入力手段(タッチ/マウス)で出し分ける */}
                        <span className="[@media(pointer:coarse)]:hidden">ドラッグで入れ替えられます</span>
                        <span className="hidden [@media(pointer:coarse)]:inline">長押しで入れ替えられます</span>
                    </p>
                </div>
                <Button
                    type="button"
                    variant="ghost"
                    className="text-primary hover:text-primary-dark"
                    onClick={handleAdd}
                >
                    <Plus className="size-4" />
                    日時を追加
                </Button>
            </div>
            {error && <p className="text-sm text-destructive">{error}</p>}
            {dates.length > 0 ? (
                <DraggableDateList
                    dates={dates}
                    onDatesChange={onDatesChange}
                    editingId={editingId}
                    onEditingIdChange={setEditingId}
                />
            ) : (
                <p className="rounded-md border border-dashed border-input px-3 py-4 text-center text-sm text-muted-foreground">
                    カレンダーで範囲を選択するか、「日時を追加」から登録できます
                </p>
            )}
        </div>
    );
};

export default DatesPanelView;
