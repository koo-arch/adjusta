'use client'
import React from 'react';
import { useAtom, useAtomValue } from 'jotai';
import { Button } from '@/components/ui/button';
import EventBasicForm from './EventBasicForm';
import FormStepper from './FormStepper';
import DraftCalendarPane from './DraftCalendarPane';
import EditCalendarPane from './EditCalendarPane';
import DraftDatesPanel from './DraftDatesPanel';
import EditDatesPanel from './EditDatesPanel';
import type { EventDraftDetail } from '@/features/events/types';
import { eventFormMessagesAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';
import { formStepAtomFamily } from '@/features/events/store/formStep';
import { ChevronLeft, ChevronRight } from 'lucide-react';

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
    const { formScope, submitLabel, isSubmitting, eventDetail } = props;
    const formErrors = useAtomValue(eventFormMessagesAtomFamily(formScope));
    const fieldErrors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const [step, setStep] = useAtom(formStepAtomFamily(formScope));
    const hasBasicErrors = !!(fieldErrors.title || fieldErrors.location || fieldErrors.description);
    const hasDatesErrors = !!(fieldErrors.selected_dates || fieldErrors.proposed_dates);

    return (
        <div className="space-y-6">
            {/* カレンダーは常設(md 以上で左)、隣のパネルだけステップで切り替える。
                ステッパーは md 以上で右端の縦レール、モバイルは上部の横一列。
                モバイルは ステッパー → パネル → カレンダー の縦積み */}
            <div className="flex flex-col gap-6 md:flex-row">
                <div className="order-3 min-w-0 md:order-1 md:flex-1">
                    {props.formType === 'draft' ? (
                        <DraftCalendarPane formScope={formScope} editingEvent={eventDetail} />
                    ) : (
                        <EditCalendarPane formScope={formScope} editingEvent={eventDetail} />
                    )}
                </div>
                {/* ページスクロールに追従させる(ヘッダー sticky 分のオフセット) */}
                <div className="order-2 md:sticky md:top-20 md:w-96 md:self-start">
                    {step === 'basic' ? (
                        <EventBasicForm formScope={formScope} />
                    ) : props.formType === 'draft' ? (
                        <DraftDatesPanel formScope={formScope} />
                    ) : (
                        <EditDatesPanel formScope={formScope} />
                    )}
                </div>
                <div className="order-1 md:order-3 md:sticky md:top-20 md:self-start">
                    <FormStepper
                        current={step}
                        hasBasicErrors={hasBasicErrors}
                        hasDatesErrors={hasDatesErrors}
                        onSelect={setStep}
                    />
                </div>
            </div>

            {/* スクロール位置に関係なく操作できるよう、送信バーは下部に固定する */}
            <div className="sticky bottom-0 z-10 border-t border-border bg-background py-3">
                <div className="flex flex-wrap items-center justify-end gap-x-4 gap-y-2">
                    {formErrors.length > 0 && (
                        <div className="min-w-0 space-y-1">
                            {formErrors.map((message) => (
                                <p key={message} className="text-sm text-destructive">
                                    {message}
                                </p>
                            ))}
                        </div>
                    )}
                    {step === 'dates' && hasBasicErrors && (
                        <button
                            type="button"
                            onClick={() => setStep('basic')}
                            className="text-sm text-destructive underline underline-offset-2 transition-opacity hover:opacity-80"
                        >
                            基本情報に入力エラーがあります
                        </button>
                    )}
                    {step === 'basic' ? (
                        <Button type="button" onClick={() => setStep('dates')}>
                            次へ
                            <ChevronRight className="size-4" />
                        </Button>
                    ) : (
                        <>
                            <Button type="button" variant="ghost" onClick={() => setStep('basic')}>
                                <ChevronLeft className="size-4" />
                                戻る
                            </Button>
                            <Button type="submit" disabled={isSubmitting}>
                                {submitLabel}
                            </Button>
                        </>
                    )}
                </div>
            </div>
        </div>
    );
};

export default EventForm;
