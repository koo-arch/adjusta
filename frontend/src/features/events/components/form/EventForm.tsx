'use client'
import React from 'react';
import { useAtomValue } from 'jotai';
import { Button } from '@/components/ui/button';
import CalendarForm from './CalendarForm';
import EventBasicForm from './EventBasicForm';
import SelectEventList from './SelectEventList';
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
            {/* モバイルは 基本情報 → カレンダー → 選択日程 の縦積み、
                md 以上は左(基本情報+選択日程)/右(カレンダー)の 2 カラム(CSS のみで分岐) */}
            <div className="grid grid-cols-1 gap-8 md:grid-cols-10 md:grid-rows-[auto_1fr] md:gap-x-6 md:gap-y-8">
                <div className="md:col-span-4">
                    <EventBasicForm formScope={formScope} />
                </div>
                <div className="md:col-span-6 md:col-start-5 md:row-span-2 md:row-start-1">
                    {props.formType === 'draft' ? (
                        <CalendarForm
                            formType="draft"
                            formScope={formScope}
                            editingEvent={eventDetail}
                        />
                    ) : (
                        <CalendarForm
                            formType="edit"
                            formScope={formScope}
                            editingEvent={eventDetail}
                        />
                    )}
                </div>
                <div className="md:col-span-4 md:col-start-1">
                    {props.formType === 'draft' ? (
                        <SelectEventList
                            formType="draft"
                            formScope={formScope}
                        />
                    ) : (
                        <SelectEventList
                            formType="edit"
                            formScope={formScope}
                        />
                    )}
                </div>
            </div>

            <div className="space-y-4 border-t border-border pt-4">
                {formErrors.length > 0 && (
                    <div className="space-y-2">
                        {formErrors.map((message) => (
                            <p key={message} className="text-sm text-destructive">
                                {message}
                            </p>
                        ))}
                    </div>
                )}
                <div className="flex justify-end">
                    <Button type="submit" disabled={isSubmitting}>
                        {submitLabel}
                    </Button>
                </div>
            </div>
        </div>
    );
};

export default EventForm;
