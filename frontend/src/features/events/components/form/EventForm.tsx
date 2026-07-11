'use client'
import React from 'react';
import { useAtomValue } from 'jotai';
import { Button } from '@/components/ui/button';
import EventBasicForm from './EventBasicForm';
import DateSelectionSection from './DateSelectionSection';
import type { EventDraftDetail } from '@/features/events/types';
import { eventFormMessagesAtomFamily } from '@/features/events/store/errors';

interface EventFormBaseProps {
    formScope: string;
    submitLabel: string;
    isSubmitting?: boolean;
    eventDetail?: EventDraftDetail;
}

type DraftEventFormProps = EventFormBaseProps & {
    formType: 'draft';
};

type EditEventFormProps = EventFormBaseProps & {
    formType: 'edit';
};

type EventFormProps = DraftEventFormProps | EditEventFormProps;

const EventForm: React.FC<EventFormProps> = (props) => {
    const formErrors = useAtomValue(eventFormMessagesAtomFamily(props.formScope));
    const { formScope, submitLabel, isSubmitting, eventDetail } = props;

    return (
        <div className="space-y-6">
            {/* 基本情報 → 候補日程 の縦のステップ構成。入力欄は読みやすい幅に絞る */}
            <div className="max-w-xl">
                <EventBasicForm formScope={formScope} />
            </div>

            <div className="border-t border-border pt-6">
                {props.formType === 'draft' ? (
                    <DateSelectionSection
                        formType="draft"
                        formScope={formScope}
                        editingEvent={eventDetail}
                    />
                ) : (
                    <DateSelectionSection
                        formType="edit"
                        formScope={formScope}
                        editingEvent={eventDetail}
                    />
                )}
            </div>

            {/* スクロール位置に関係なく保存できるよう、送信バーは下部に固定する */}
            <div className="sticky bottom-0 z-10 border-t border-border bg-background py-3">
                <div className="flex flex-wrap items-center justify-end gap-x-6 gap-y-2">
                    {formErrors.length > 0 && (
                        <div className="min-w-0 space-y-1">
                            {formErrors.map((message) => (
                                <p key={message} className="text-sm text-destructive">
                                    {message}
                                </p>
                            ))}
                        </div>
                    )}
                    <Button type="submit" disabled={isSubmitting}>
                        {submitLabel}
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default EventForm;
