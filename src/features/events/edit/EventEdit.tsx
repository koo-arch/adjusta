'use client'
import React, { useEffect } from 'react';
import axios from '@/lib/axios/public';
import { useAtom } from 'jotai';
import {  proposedDatesAtom, sendProposedDatesAtom } from '@/atoms/calendar';
import { isConfirmedAtomFamily } from '@/atoms/event';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useForm, SubmitHandler, FormProvider } from 'react-hook-form';
import { useParams } from 'next/navigation';
import { type DiscriminatedEventForm, DiscriminatedEventFormResolver } from '../zod';
import EventForm from '../EventForm';

const EventEdit = () => {
    const params = useParams<{ slug: string }>();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.slug);
    const [proposedDates, setProposedDates] = useAtom(proposedDatesAtom);
    const [sendProposedDates] = useAtom(sendProposedDatesAtom);
    const [isConfirmed] = useAtom(isConfirmedAtomFamily(params.slug));

    const method = useForm<DiscriminatedEventForm>({
        resolver: DiscriminatedEventFormResolver,
        defaultValues: {
            id: null,
            form_type: "edit",
            slug: params.slug,
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

    useEffect(() => {
        console.log("errors", errors);
    }, [errors]);

    const putEventUpdate = async (data: DiscriminatedEventForm) => {
        return await axios.put(`api/calendar/event/draft/${params.slug}`, data);
    }

    const onSubmit: SubmitHandler<DiscriminatedEventForm> = (data) => {
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
                        <EventForm eventDetail={eventDetail} />
                    </form>
                </FormProvider>
            }    
        </div>
    )
}

export default EventEdit;