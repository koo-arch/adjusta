'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { useAtom } from 'jotai';
import { ProposedDate, proposedDatesAtom, proposedEventsAtom, sendProposedDatesAtom } from '@/atoms/calendar';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useForm, SubmitHandler, FormProvider } from 'react-hook-form';
import { useParams } from 'next/navigation';
import type { EventUpdateForm } from '../type';
import EventBasicForm from '../EventBasicForm';
import SelectableCalendar from '@/features/calendar/SelectableCalendar';
import DraggableEventList from '../DraggableEventList';

const EventEdit = () => {
    const params = useParams<{ id: string }>();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);
    const [proposedDates, setProposedDates] = useAtom(proposedDatesAtom);
    const [sendProposedDates] = useAtom(sendProposedDatesAtom);
    const method = useForm<EventUpdateForm>({
        defaultValues: {
            id: params.id,
        }
    }
    );
    const { handleSubmit, setValue, formState: { errors } } = method;

    useEffect(() => {
        if (eventDetail) {
            setValue('proposed_dates', sendProposedDates);
        }
    }), [proposedDates, setValue];
    
    useEffect(() => {
        if (eventDetail) {
            setProposedDates(eventDetail.proposed_dates);
        }
    }, [eventDetail, setProposedDates]);

    const putEventUpdate = async (data: EventUpdateForm) => {
        return await axios.put(`api/calendar/event/draft/${params.id}`, data);
    }

    const onSubmit: SubmitHandler<EventUpdateForm> = (data) => {
        putEventUpdate(data)
            .then(res => {
                console.log(res);
            })
            .catch(err => {
                console.log(err);
            })
    }

    if (isLoading) {
        return <p>Loading...</p>;
    }

    return (
        <div>
            {eventDetail &&
                <FormProvider {...method}>
                    <form onSubmit={handleSubmit(onSubmit)}>
                        <EventBasicForm
                            title={eventDetail.title}
                            description={eventDetail.description}
                            location={eventDetail.location}
                        />
                        <SelectableCalendar
                            editingEvent={eventDetail.proposed_dates}
                            dateAtom={proposedDatesAtom}
                            eventAtom={proposedEventsAtom}
                        />
                        <DraggableEventList<ProposedDate> atom={proposedDatesAtom}/>
                        <button type="submit">Submit</button>
                    </form>
                </FormProvider>
            }
        </div>
    )
}

export default EventEdit;