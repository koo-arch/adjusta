'use client'
import React, { useMemo, useState } from 'react';
import Button from '@/components/Button';
import IconButton from '@/components/IconButton';
import ToggleButton from '@/components/ToggleButton';
import Modal from '@/components/Modal';
import DropdownSelect from '@/components/DropdownSelect';
import { formatJaDateSpan } from '@/lib/date/format';
import DateTimePicker from '@/components/DateTimePicker';
import type { EventDraftDetail } from '@/features/events/types';
import { useConfirmEventMutation } from '@/features/events/detail/hooks/useConfirmEventMutation';
import { MdEditCalendar } from 'react-icons/md';
import { FaRegCalendarCheck } from 'react-icons/fa6';

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
    const [isDropdownSelected, setIsDropdownSelected] = useState(true);
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

    const handleToggle = (selected: string) => {
        setIsDropdownSelected(selected === '候補日程を選択');
        resetConfirmDate();
    };

    const handleSelectProposedDate = (date: EventDraftDetail['proposed_dates'][number] | null) => {
        resetMutationErrorState();

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
            selectionMode: isDropdownSelected ? 'dropdown' : 'manual',
        });
        if (confirmed) {
            setIsOpen(false);
        }
    };

    return (
        <>
            <IconButton
                iconColor={isConfirmed? 'primary': 'success'}
                iconSize={"lg"}
                onClick={() => setIsOpen(true)}
            >
                {isConfirmed ? (
                    <MdEditCalendar />
                ) : (
                    <FaRegCalendarCheck />
                )}
            </IconButton>
            <Modal
                isOpen={isOpen}
                onClose={() => setIsOpen(false)}
                title={isConfirmed ? '日程を変更' : '日程を確定'}
                description="確定させる日程を選択してください。"
                actions={
                    <Button
                        variant='solid'
                        intent='primary'
                        size='md'
                        type='button'
                        onClick={handleSubmit}
                        disabled={confirmEventMutation.isPending}
                    >
                        確定
                    </Button>
                }
            >
                {confirmEventMutation.errors.formErrors.length > 0 && (
                    <div className="mb-4 space-y-2">
                        {confirmEventMutation.errors.formErrors.map((message) => (
                            <p key={message} className="text-sm text-red-500">
                                {message}
                            </p>
                        ))}
                    </div>
                )}
                <div className="mb-4">
                    <ToggleButton
                        options={['候補日程を選択', '手動で入力']}
                        selected={isDropdownSelected ? '候補日程を選択' : '手動で入力'}
                        onToggle={handleToggle}
                        renderLabel={(option) => option}
                    />
                </div>
                {isDropdownSelected ? 
                    <DropdownSelect
                        label='日程'
                        options={proposedDates}
                        value={selectedProposedDate}
                        renderLabel={(date) => 
                            date && (
                                <>
                                    {`${formatJaDateSpan(date.start, date.end)}`}
                                </>
                            )
                        }
                        onChange={handleSelectProposedDate}
                        error={!!errors.confirm_date}
                        helperText={errors.confirm_date}
                    /> : 
                    <div>
                        <DateTimePicker
                            label='開始日時'
                            selected={confirmDate.start}
                            onChange={(date) => {
                                setConfirmDate((prev) => ({ ...prev, start: date }));
                                resetMutationErrorState();
                            }}
                            error={!!errors['confirm_date.start']}
                            helperText={errors['confirm_date.start']}
                        />
                        <DateTimePicker
                            label='終了日時'
                            selected={confirmDate.end}
                            onChange={(date) => {
                                setConfirmDate((prev) => ({ ...prev, end: date }));
                                resetMutationErrorState();
                            }}
                            error={!!errors['confirm_date.end']}
                            helperText={errors['confirm_date.end']}
                        />
                    </div>}
            </Modal>
        </>
    )
}

export default ConfirmButton;
