'use client'
import React from 'react';
import DraggableList from '@/components/DraggableList';
import { Button } from '@/components/ui/button';
import { DateTimePicker } from '@/components/common/DateTimePicker/DateTimePicker';
import { formatJaDateOnly, formatJaTimeSpan } from '@/lib/date/format';
import { Check, Pencil, X } from 'lucide-react';
import type { ProposedDate, SelectedDate } from '@/features/events/store/dates';

interface DraggableDateListProps<T extends SelectedDate | ProposedDate> {
    dates: T[];
    onDatesChange: React.Dispatch<React.SetStateAction<T[]>>;
    // 編集中の行 ID(親が保持。追加直後の行を編集モードで開くため)
    editingId: string | null;
    onEditingIdChange: (id: string | null) => void;
}

const DraggableDateList = <T extends SelectedDate | ProposedDate>({
    dates,
    onDatesChange,
    editingId,
    onEditingIdChange,
 }: DraggableDateListProps<T>) => {
    const handleReorder = (newDates: T[]) => {
        onDatesChange(newDates);
    }

    const handleDelete = (id: string) => {
        if (editingId === id) {
            onEditingIdChange(null);
        }
        onDatesChange(dates.filter(date => date.id !== id));
    }

    const handleTimeChange = (id: string, field: 'start' | 'end', value: Date | null) => {
        // Invalid Date を atom に入れると表示側の format が落ちるため弾く
        if (!value || Number.isNaN(value.getTime())) {
            return;
        }
        onDatesChange(dates.map((date) => (date.id === id ? { ...date, [field]: value } : date)));
    }

    return (
        <DraggableList
            items={dates}
            onReorder={handleReorder}
            renderItem={(date: T, index) => {
                const isEditing = editingId === date.id;
                const isInvalid = date.end.getTime() <= date.start.getTime();

                if (isEditing) {
                    return (
                        <div className="space-y-2">
                            <div className="flex items-center justify-between gap-2">
                                <span
                                    aria-label={`第${index + 1}候補`}
                                    className="grid size-7 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary"
                                >
                                    {index + 1}
                                </span>
                                <Button
                                    type="button"
                                    variant="ghost"
                                    size="icon"
                                    className="h-8 w-8 text-primary hover:text-primary-dark"
                                    aria-label="編集を完了"
                                    title="編集を完了"
                                    onClick={() => onEditingIdChange(null)}
                                >
                                    <Check />
                                </Button>
                            </div>
                            <DateTimePicker
                                label="開始"
                                selected={date.start}
                                onChange={(value) => handleTimeChange(date.id, 'start', value)}
                            />
                            <DateTimePicker
                                label="終了"
                                selected={date.end}
                                onChange={(value) => handleTimeChange(date.id, 'end', value)}
                            />
                            {isInvalid && (
                                <p className="text-xs text-destructive">終了日時は開始日時より後にしてください</p>
                            )}
                        </div>
                    );
                }

                return (
                    <div className="flex items-center justify-between gap-2">
                        <div className="flex min-w-0 items-center gap-3">
                            <span
                                aria-label={`第${index + 1}候補`}
                                className="grid size-7 shrink-0 place-items-center rounded-full bg-primary/10 text-sm font-semibold text-primary"
                            >
                                {index + 1}
                            </span>
                            <div className="min-w-0">
                                <p className="text-sm font-medium text-foreground">{formatJaDateOnly(date.start)}</p>
                                <p className="text-sm text-muted-foreground">{formatJaTimeSpan(date.start, date.end)}</p>
                            </div>
                        </div>
                        <div className="flex shrink-0 items-center">
                            <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 text-muted-foreground hover:text-foreground"
                                aria-label="日時を編集"
                                title="日時を編集"
                                onClick={() => onEditingIdChange(date.id)}
                            >
                                <Pencil />
                            </Button>
                            <Button
                                type="button"
                                variant="ghost"
                                size="icon"
                                className="h-8 w-8 text-muted-foreground hover:text-destructive"
                                aria-label="この日程を削除"
                                onClick={() => handleDelete(date.id)}
                            >
                                <X />
                            </Button>
                        </div>
                    </div>
                );
            }}
            getKey={(date: T) => date.id}
            disabledIds={editingId ? [editingId] : undefined}
        />
    )
}

export default DraggableDateList;
