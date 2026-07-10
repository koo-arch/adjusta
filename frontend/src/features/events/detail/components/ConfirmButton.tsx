'use client'
import React, { useMemo, useState } from 'react';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select';
import { Tabs, TabsList, TabsTrigger } from '@/components/ui/tabs';
import DateTimePicker from '@/components/DateTimePicker';
import { formatJaDateSpan } from '@/lib/date/format';
import type { EventDraftDetail } from '@/features/events/types';
import { useConfirmEventMutation } from '@/features/events/detail/hooks/useConfirmEventMutation';
import { CalendarCheck } from 'lucide-react';

interface ConfirmButtonProps {
    eventID: string;
    detail: EventDraftDetail;
    isConfirmed: boolean;
}

interface ConfirmDateInput {
    id: string | null;
    google_event_id?: string;
    start: Date | null;
    end: Date | null;
    priority: number;
}

type SelectionMode = 'dropdown' | 'manual';

const buildEmptyConfirmDate = (googleEventID?: string): ConfirmDateInput => ({
    id: null,
    google_event_id: googleEventID,
    start: null,
    end: null,
    priority: 0,
});

const ConfirmButton: React.FC<ConfirmButtonProps> = ({ eventID, detail, isConfirmed }) => {
    const confirmEventMutation = useConfirmEventMutation(eventID);
    const proposedDates = detail.proposed_dates;
    const confirmedGoogleEventID = detail.confirmed_google_event_id ?? detail.google_event_id;
    const [isOpen, setIsOpen] = useState(false);
    const [selectionMode, setSelectionMode] = useState<SelectionMode>('dropdown');
    const [confirmDate, setConfirmDate] = useState<ConfirmDateInput>(buildEmptyConfirmDate(confirmedGoogleEventID));
    const errors = confirmEventMutation.errors.fieldErrors;

    const selectedProposedDate = useMemo(
        () => proposedDates.find((date) => date.id === confirmDate.id) ?? null,
        [confirmDate.id, proposedDates],
    );

    const resetMutationErrorState = () => {
        confirmEventMutation.reset();
    };

    const resetConfirmDate = () => {
        setConfirmDate(buildEmptyConfirmDate(confirmedGoogleEventID));
        resetMutationErrorState();
    };

    const handleOpenChange = (open: boolean) => {
        setIsOpen(open);
        if (!open) {
            setSelectionMode('dropdown');
            resetConfirmDate();
        }
    };

    const handleModeChange = (value: string) => {
        setSelectionMode(value === 'manual' ? 'manual' : 'dropdown');
        resetConfirmDate();
    };

    const handleSelectProposedDate = (id: string) => {
        resetMutationErrorState();
        const date = proposedDates.find((proposed) => proposed.id === id);

        if (!date) {
            setConfirmDate(buildEmptyConfirmDate(confirmedGoogleEventID));
            return;
        }

        setConfirmDate({
            id: date.id,
            google_event_id: date.google_event_id,
            start: date.start,
            end: date.end,
            priority: date.priority,
        });
    };

    const handleSubmit = async () => {
        const confirmed = await confirmEventMutation.submit({
            confirm_date: {
                id: confirmDate.id,
                google_event_id: confirmDate.google_event_id,
                start: confirmDate.start,
                end: confirmDate.end,
                priority: confirmDate.priority,
            },
            selectionMode,
        });
        if (confirmed) {
            handleOpenChange(false);
        }
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleOpenChange}>
            <DialogTrigger asChild>
                {/* 確定はこの画面の主目的なのでラベル付きで置く(ui-guidelines)。
                    フラット構成に合わせ、塗りつぶしではなく primary 色のテキストボタンにする */}
                <Button variant="ghost" className="text-primary hover:text-primary-dark">
                    <CalendarCheck className="size-4" />
                    {isConfirmed ? '確定日程を変更' : '日程を確定'}
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md">
                <DialogHeader>
                    <DialogTitle>{isConfirmed ? '日程を変更' : '日程を確定'}</DialogTitle>
                    <DialogDescription>
                        確定させる日程を選択してください。確定するとメインカレンダーに本予定として登録されます。
                    </DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <Tabs value={selectionMode} onValueChange={handleModeChange}>
                        <TabsList className="grid w-full grid-cols-2">
                            <TabsTrigger value="dropdown">候補から選択</TabsTrigger>
                            <TabsTrigger value="manual">手動で入力</TabsTrigger>
                        </TabsList>
                    </Tabs>
                    {confirmEventMutation.errors.formErrors.length > 0 && (
                        <div className="space-y-2">
                            {confirmEventMutation.errors.formErrors.map((message) => (
                                <p key={message} className="text-sm text-destructive">
                                    {message}
                                </p>
                            ))}
                        </div>
                    )}
                    {selectionMode === 'dropdown' ? (
                        <div className="space-y-2">
                            <Label>日程</Label>
                            <Select
                                value={selectedProposedDate?.id ?? undefined}
                                onValueChange={handleSelectProposedDate}
                            >
                                <SelectTrigger aria-invalid={!!errors.confirm_date}>
                                    <SelectValue placeholder="候補日程を選択" />
                                </SelectTrigger>
                                <SelectContent>
                                    {proposedDates.map((date) => (
                                        <SelectItem key={date.id} value={date.id}>
                                            {formatJaDateSpan(date.start, date.end)}
                                        </SelectItem>
                                    ))}
                                </SelectContent>
                            </Select>
                            {errors.confirm_date && (
                                <p className="text-sm text-destructive">{errors.confirm_date}</p>
                            )}
                        </div>
                    ) : (
                        <div className="space-y-4">
                            <DateTimePicker
                                label="開始日時"
                                selected={confirmDate.start}
                                onChange={(date: Date | null) => {
                                    setConfirmDate((prev) => ({ ...prev, start: date }));
                                    resetMutationErrorState();
                                }}
                                error={!!errors['confirm_date.start']}
                                helperText={errors['confirm_date.start']}
                            />
                            <DateTimePicker
                                label="終了日時"
                                selected={confirmDate.end}
                                onChange={(date: Date | null) => {
                                    setConfirmDate((prev) => ({ ...prev, end: date }));
                                    resetMutationErrorState();
                                }}
                                error={!!errors['confirm_date.end']}
                                helperText={errors['confirm_date.end']}
                            />
                        </div>
                    )}
                </div>
                <DialogFooter>
                    <Button type="button" onClick={handleSubmit} disabled={confirmEventMutation.isPending}>
                        確定する
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};

export default ConfirmButton;
