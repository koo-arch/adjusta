'use client'
import React, { useId, useState } from 'react';
import { format } from 'date-fns';
import { ja } from 'date-fns/locale';
import { Button } from '@/components/ui/button';
import { Calendar } from '@/components/ui/calendar';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { cn } from '@/lib/utils';
import { CalendarIcon } from 'lucide-react';

interface DateTimePickerProps {
    label?: string;
    selected: Date | null;
    onChange: (date: Date | null) => void;
    error?: boolean;
    helperText?: string;
}

// shadcn calendar(日付)+ time 入力(時刻)を 1 つの popover にまとめた日時ピッカー
const isValidDate = (date: Date | null): date is Date =>
    date !== null && !Number.isNaN(date.getTime());

export const DateTimePicker: React.FC<DateTimePickerProps> = ({
    label,
    selected: selectedProp,
    onChange,
    error,
    helperText,
}) => {
    const [isOpen, setIsOpen] = useState(false);
    // time 入力の入力途中の値。不完全な値でも表示を巻き戻さず、完成した時だけ反映する
    const [timeDraft, setTimeDraft] = useState<string | null>(null);
    const triggerId = useId();
    // Invalid Date が渡されても描画で落ちないよう未選択として扱う
    const selected = isValidDate(selectedProp) ? selectedProp : null;

    const handleOpenChange = (open: boolean) => {
        setIsOpen(open);
        if (!open) {
            setTimeDraft(null);
        }
    };

    const handleSelectDate = (date: Date | undefined) => {
        if (!date) {
            return;
        }
        const next = new Date(date);
        // 時刻は既存の選択値を引き継ぐ(未選択なら 00:00)
        if (selected) {
            next.setHours(selected.getHours(), selected.getMinutes(), 0, 0);
        }
        if (Number.isNaN(next.getTime())) {
            return;
        }
        onChange(next);
        // 同日クリックでの解除は扱わず、日付確定として閉じずに時刻調整へ進める
    };

    const handleTimeChange = (value: string) => {
        setTimeDraft(value);
        const [hours, minutes] = value.split(':').map(Number);
        // 入力途中(空・欠けたセグメント)は反映しない
        if (!Number.isInteger(hours) || !Number.isInteger(minutes)) {
            return;
        }
        const base = selected ?? new Date();
        const next = new Date(base);
        next.setHours(hours, minutes, 0, 0);
        if (Number.isNaN(next.getTime())) {
            return;
        }
        onChange(next);
    };

    return (
        <div className="space-y-1.5">
            {label && <Label htmlFor={triggerId}>{label}</Label>}
            <Popover open={isOpen} onOpenChange={handleOpenChange}>
                <PopoverTrigger asChild>
                    <Button
                        id={triggerId}
                        type="button"
                        variant="outline"
                        aria-invalid={error || undefined}
                        className={cn(
                            'w-full justify-start font-normal',
                            !selected && 'text-muted-foreground',
                            error && 'border-destructive focus-visible:ring-destructive',
                        )}
                    >
                        <CalendarIcon className="size-4" />
                        {selected ? format(selected, 'M月d日(E) H:mm', { locale: ja }) : '日時を選択'}
                    </Button>
                </PopoverTrigger>
                <PopoverContent align="start" className="w-auto p-0">
                    <Calendar
                        mode="single"
                        locale={ja}
                        selected={selected ?? undefined}
                        defaultMonth={selected ?? undefined}
                        onSelect={handleSelectDate}
                        required={false}
                    />
                    <div className="border-t border-border p-3">
                        <Input
                            type="time"
                            aria-label="時刻"
                            value={timeDraft ?? (selected ? format(selected, 'HH:mm') : '')}
                            onChange={(e) => handleTimeChange(e.target.value)}
                            onBlur={() => setTimeDraft(null)}
                        />
                    </div>
                </PopoverContent>
            </Popover>
            {helperText && (
                <p className={cn('text-sm', error ? 'text-destructive' : 'text-muted-foreground')}>
                    {helperText}
                </p>
            )}
        </div>
    );
};
