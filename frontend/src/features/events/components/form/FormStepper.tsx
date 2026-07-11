'use client'
import React from 'react';
import { cn } from '@/lib/utils';
import type { EventFormStep } from '@/features/events/store/formStep';
import { ChevronRight } from 'lucide-react';

interface FormStepperProps {
    current: EventFormStep;
    hasBasicErrors?: boolean;
    hasDatesErrors?: boolean;
    onSelect: (step: EventFormStep) => void;
}

const STEPS: { key: EventFormStep; label: string }[] = [
    { key: 'basic', label: '基本情報' },
    { key: 'dates', label: '候補日程' },
];

// クリックで行き来できるステップ表示。エラーのあるステップは色で示す
const FormStepper: React.FC<FormStepperProps> = ({ current, hasBasicErrors, hasDatesErrors, onSelect }) => {
    const hasError = (step: EventFormStep) =>
        step === 'basic' ? !!hasBasicErrors : !!hasDatesErrors;

    return (
        <nav aria-label="入力ステップ" className="flex items-center gap-2">
            {STEPS.map((step, index) => {
                const isActive = current === step.key;
                const isError = hasError(step.key);
                return (
                    <React.Fragment key={step.key}>
                        {index > 0 && <ChevronRight className="size-4 text-muted-foreground" aria-hidden />}
                        <button
                            type="button"
                            onClick={() => onSelect(step.key)}
                            aria-current={isActive ? 'step' : undefined}
                            className={cn(
                                'flex items-center gap-2 text-sm font-medium transition-colors',
                                isActive ? 'text-primary' : 'text-muted-foreground hover:text-foreground',
                                isError && 'text-destructive',
                            )}
                        >
                            <span
                                className={cn(
                                    'grid size-6 shrink-0 place-items-center rounded-full border text-xs',
                                    isActive
                                        ? 'border-primary bg-primary text-primary-foreground'
                                        : 'border-border bg-card',
                                    isError && 'border-destructive bg-destructive text-destructive-foreground',
                                )}
                            >
                                {index + 1}
                            </span>
                            {step.label}
                        </button>
                    </React.Fragment>
                );
            })}
        </nav>
    );
};

export default FormStepper;
