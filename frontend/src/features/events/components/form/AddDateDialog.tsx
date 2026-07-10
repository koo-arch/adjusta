'use client'
import React, { useState } from 'react';
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
import DateTimePicker from '@/components/DateTimePicker';
import { Plus } from 'lucide-react';

interface AddDateDialogProps {
    onAdd: (date: { start: Date; end: Date }) => void;
}

// カレンダーのドラッグ選択と併設する、日時の直接入力手段(ui-review P2 #7)
const AddDateDialog: React.FC<AddDateDialogProps> = ({ onAdd }) => {
    const [isOpen, setIsOpen] = useState(false);
    const [start, setStart] = useState<Date | null>(null);
    const [end, setEnd] = useState<Date | null>(null);
    const [error, setError] = useState<string | null>(null);

    const reset = () => {
        setStart(null);
        setEnd(null);
        setError(null);
    };

    const handleOpenChange = (open: boolean) => {
        setIsOpen(open);
        if (!open) {
            reset();
        }
    };

    const handleAdd = () => {
        if (!start || !end) {
            setError('開始日時と終了日時を入力してください');
            return;
        }
        if (end.getTime() <= start.getTime()) {
            setError('終了日時は開始日時より後に設定してください');
            return;
        }
        onAdd({ start, end });
        handleOpenChange(false);
    };

    return (
        <Dialog open={isOpen} onOpenChange={handleOpenChange}>
            <DialogTrigger asChild>
                <Button type="button" variant="ghost" className="text-primary hover:text-primary-dark">
                    <Plus className="size-4" />
                    日時を追加
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-md">
                <DialogHeader>
                    <DialogTitle>日時を追加</DialogTitle>
                    <DialogDescription>候補にする開始日時と終了日時を入力してください。</DialogDescription>
                </DialogHeader>
                <div className="space-y-4">
                    <DateTimePicker
                        label="開始日時"
                        selected={start}
                        onChange={(date: Date | null) => {
                            setStart(date);
                            setError(null);
                        }}
                    />
                    <DateTimePicker
                        label="終了日時"
                        selected={end}
                        onChange={(date: Date | null) => {
                            setEnd(date);
                            setError(null);
                        }}
                    />
                    {error && <p className="text-sm text-destructive">{error}</p>}
                </div>
                <DialogFooter>
                    <Button type="button" onClick={handleAdd}>
                        追加する
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>
    );
};

export default AddDateDialog;
