'use client'
import React from 'react';
import { useAtomValue } from 'jotai';
import Button from '@/components/Button';
import { useMediaQuery } from 'react-responsive';
import CalendarForm from './CalendarForm';
import EventBasicForm from './EventBasicForm';
import SelectEventList from './SelectEventList';
import type { EventDraftDetail } from '@/features/events/types';
import { eventFormMessagesAtomFamily } from '@/features/events/store/errors';

interface EventFormBaseProps {
    formScope: string;
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
    const isMobile = useMediaQuery({ maxWidth: 768 });
    const formErrors = useAtomValue(eventFormMessagesAtomFamily(props.formScope));
    const { formScope, isSubmitting, eventDetail } = props;

    return (
        <div>
            <div className="mx-auto grid grid-cols-1 md:grid-cols-10 gap-6 mb-4">
                <div className="md:col-span-4 space-y-6">
                    <section>
                        <EventBasicForm formScope={formScope} />
                    </section>
                    {isMobile && (
                        <section>
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
                        </section>
                    )}
                    <section>
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
                    </section>
                </div>
                {!isMobile && (
                    <section className="md:col-span-6">
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
                    </section>
                )}
            </div>
            
            <div className="border-t border-gray-300 mb-4"></div>
            {formErrors.length > 0 && (
                <div className="mb-4 space-y-2">
                    {formErrors.map((message) => (
                        <p key={message} className="text-sm text-red-500 text-center">
                            {message}
                        </p>
                    ))}
                </div>
            )}
            <div className="py-5 flex items-center justify-center">
                <Button type="submit" disabled={isSubmitting}>登録する</Button>
            </div>
        </div>
    );
};

export default EventForm;
