'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { useAtom } from 'jotai';
import {  proposedDatesAtom, sendProposedDatesAtom } from '@/atoms/calendar';
import { isConfirmedAtomFamily } from '@/atoms/event';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useForm, SubmitHandler, FormProvider } from 'react-hook-form';
import { useParams } from 'next/navigation';
import type { EventUpdateForm } from '../type';
import EventForm from '../EventForm';

const EventEdit = () => {
    const params = useParams<{ id: string }>();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);
    const [proposedDates, setProposedDates] = useAtom(proposedDatesAtom);
    const [sendProposedDates] = useAtom(sendProposedDatesAtom);
    const [isConfirmed] = useAtom(isConfirmedAtomFamily(params.id));

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
        // statusをisConfirmedの値に応じて変更
        const updateData = {
            ...data,
            status: isConfirmed ? "confirmed" : "pending",
        }

        putEventUpdate(updateData)
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
                        <EventForm
                            formType="edit"
                            eventDetail={eventDetail}
                        />
                    </form>
                </FormProvider>
            }
        </div>
    )
}

export default EventEdit;