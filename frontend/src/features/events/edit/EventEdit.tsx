'use client'
import React, { useEffect, useRef } from 'react';
import { useAtom } from 'jotai';
import {  proposedDatesAtom, sendProposedDatesAtom } from '@/atoms/calendar';
import { isConfirmedAtomFamily } from '@/atoms/event';
import { useFetchEventDetail } from '@/hooks/event/useFetchEventDetail';
import { useForm, SubmitHandler, FormProvider } from 'react-hook-form';
import { useParams, useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api/client';
import { type DiscriminatedEventForm, DiscriminatedEventFormResolver } from '../zod';
import EventForm from '../EventForm';

const EventEdit = () => {
    const params = useParams<{ id: string }>();
    const router = useRouter();
    const { eventDetail, isLoading, error } = useFetchEventDetail(params.id);
    const [, setProposedDates] = useAtom(proposedDatesAtom);
    const [sendProposedDates] = useAtom(sendProposedDatesAtom);
    const [isConfirmed, setIsConfirmed] = useAtom(isConfirmedAtomFamily(params.id));
    const initializedRef = useRef(false);

    const method = useForm<DiscriminatedEventForm>({
        resolver: DiscriminatedEventFormResolver,
        defaultValues: {
            id: null,
            form_type: "edit",
        }
    }
    );
    const { handleSubmit, setValue } = method;

    useEffect(() => {
        if (eventDetail) {
            setValue('proposed_dates', sendProposedDates);
        }
    }, [eventDetail, sendProposedDates, setValue]);
    
    useEffect(() => {
        if (!eventDetail || initializedRef.current) {
            return;
        }

        setProposedDates(eventDetail.proposed_dates);
        setIsConfirmed(eventDetail.status === 'confirmed');
        initializedRef.current = true;
    }, [eventDetail, setIsConfirmed, setProposedDates]);

    const putEventUpdate = async (data: DiscriminatedEventForm) => {
        return apiClient.put<void, DiscriminatedEventForm>(`/api/calendar/event/draft/${params.id}`, data);
    }

    const onSubmit: SubmitHandler<DiscriminatedEventForm> = (data) => {
        if (data.form_type !== 'edit') {
            return;
        }

        // statusをisConfirmedの値に応じて変更
        const updateData: DiscriminatedEventForm = {
            ...data,
            status: isConfirmed ? "confirmed" : "active",
        }

        putEventUpdate(updateData)
            .then(res => {
                console.log(res);
            })
            .catch(err => {
                console.log(err);
            })
    }

    useEffect(() => {
        if (!isLoading && (!eventDetail || error)) {
            router.replace('/schedule/draft');
        }
    }, [error, eventDetail, isLoading, router]);

    if (isLoading) {
        return <p>Loading...</p>;
    }

    if (error || !eventDetail) {
        return null;
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
