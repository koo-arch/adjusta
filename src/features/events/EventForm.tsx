'use client'
import React from 'react';
import Button from '@/components/Button';
import { useMediaQuery } from 'react-responsive';
import CalendarForm from './CalendarForm';
import EventBasicForm from './EventBasicForm';
import SelectEventList from './SelectEventList';
import type { EventDraftDetail } from '@/hooks/event/type';

interface EventFormProps {
    formType: 'draft' | 'edit';
    eventDetail?: EventDraftDetail;
}


const EventForm: React.FC<EventFormProps> = ({ formType, eventDetail }) => {
    const isMobile = useMediaQuery({ maxWidth: 768 });

    if (formType == 'edit' && !eventDetail) {
        return null;
    }

    return (
        <div>
            <div className="mx-auto grid grid-cols-1 md:grid-cols-10 gap-6 mb-4">
                <div className="md:col-span-4 space-y-6">
                    <section>
                        <EventBasicForm
                            description={eventDetail?.description}
                            location={eventDetail?.location}
                        />
                    </section>
                    {isMobile && (
                        <section>
                            <CalendarForm
                                formType={formType}
                                editingEvent={eventDetail}
                            />
                        </section>
                    )}
                    <section>
                        <SelectEventList
                            formType={formType}
                        />
                    </section>
                </div>
                {!isMobile && (
                    <section className="md:col-span-6">
                        <CalendarForm
                            formType={formType}
                            editingEvent={eventDetail}
                        />
                    </section>
                )}
            </div>
            
            <div className="border-t border-gray-300 mb-4"></div>
            <div className="py-5 flex items-center justify-center">
                <Button type="submit">登録する</Button>
            </div>
        </div>
    );
};

export default EventForm;