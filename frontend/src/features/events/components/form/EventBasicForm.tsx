'use client'
import React, { useId } from 'react';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import {
    descriptionAtomFamily,
    locationAtomFamily,
    titleAtomFamily,
} from '@/features/events/store/calendar';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { cn } from '@/lib/utils';
import { clearEditedEventFieldStateAtomFamily, mergedEventFormErrorsAtomFamily } from '@/features/events/store/errors';

interface EventBasicFormProps {
    formScope: string;
}

const FieldError: React.FC<{ id: string; message?: string }> = ({ id, message }) => {
    if (!message) {
        return null;
    }
    return (
        <p id={id} className="text-sm text-destructive">
            {message}
        </p>
    );
};

const EventBasicForm: React.FC<EventBasicFormProps> = ({
    formScope,
}) => {
    const [title, setTitle] = useAtom(titleAtomFamily(formScope));
    const [description, setDescription] = useAtom(descriptionAtomFamily(formScope));
    const [location, setLocation] = useAtom(locationAtomFamily(formScope));
    const errors = useAtomValue(mergedEventFormErrorsAtomFamily(formScope));
    const clearEditedFieldState = useSetAtom(clearEditedEventFieldStateAtomFamily(formScope));
    const titleId = useId();
    const locationId = useId();
    const descriptionId = useId();

    return (
        <section className="space-y-4">
            <div>
                <h2 className="text-lg font-bold leading-snug tracking-normal text-gray-900">基本情報</h2>
                <p className="mt-1 text-sm text-muted-foreground">
                    日程調整するイベントのタイトルや詳細を入力してください
                </p>
            </div>
            <div className="space-y-2">
                <Label htmlFor={titleId}>
                    タイトル
                    <span className="ml-1 text-xs font-normal text-destructive">必須</span>
                </Label>
                <Input
                    id={titleId}
                    value={title}
                    placeholder="例: チーム定例MTG"
                    aria-invalid={!!errors.title}
                    aria-describedby={errors.title ? `${titleId}-error` : undefined}
                    className={cn(errors.title && 'border-destructive focus-visible:ring-destructive')}
                    onChange={(e) => {
                        setTitle(e.target.value);
                        clearEditedFieldState('title');
                    }}
                />
                <FieldError id={`${titleId}-error`} message={errors.title} />
            </div>
            <div className="space-y-2">
                <Label htmlFor={locationId}>場所</Label>
                <Input
                    id={locationId}
                    value={location}
                    placeholder="例: 会議室A / オンライン"
                    aria-invalid={!!errors.location}
                    aria-describedby={errors.location ? `${locationId}-error` : undefined}
                    className={cn(errors.location && 'border-destructive focus-visible:ring-destructive')}
                    onChange={(e) => {
                        setLocation(e.target.value);
                        clearEditedFieldState('location');
                    }}
                />
                <FieldError id={`${locationId}-error`} message={errors.location} />
            </div>
            <div className="space-y-2">
                <Label htmlFor={descriptionId}>説明</Label>
                <Textarea
                    id={descriptionId}
                    value={description}
                    rows={4}
                    placeholder="イベントの目的や補足があれば入力"
                    aria-invalid={!!errors.description}
                    aria-describedby={errors.description ? `${descriptionId}-error` : undefined}
                    className={cn(errors.description && 'border-destructive focus-visible:ring-destructive')}
                    onChange={(e) => {
                        setDescription(e.target.value);
                        clearEditedFieldState('description');
                    }}
                />
                <FieldError id={`${descriptionId}-error`} message={errors.description} />
            </div>
        </section>
    );
}

export default EventBasicForm;
