'use client'
import React from 'react';
import { cn } from '@/lib/utils';
import type { EventFormStep } from '@/features/events/store/formStep';
import { CalendarDays, NotebookPen, type LucideIcon } from 'lucide-react';

interface FormStepperProps {
    current: EventFormStep;
    hasBasicErrors?: boolean;
    hasDatesErrors?: boolean;
    onSelect: (step: EventFormStep) => void;
}

const STEPS: { key: EventFormStep; label: string; icon: LucideIcon }[] = [
    { key: 'basic', label: '基本情報', icon: NotebookPen },
    { key: 'dates', label: '候補日程', icon: CalendarDays },
];

// アイコンのステップレール(lg 未満は横一列、lg 以上は右端の縦レール)。
// クリックで行き来でき、入力エラーのあるステップは赤いドットで示す
const FormStepper: React.FC<FormStepperProps> = ({ current, hasBasicErrors, hasDatesErrors, onSelect }) => {
    const hasError = (step: EventFormStep) =>
        step === 'basic' ? !!hasBasicErrors : !!hasDatesErrors;

    return (
        <nav aria-label="入力ステップ" className="flex flex-row gap-2 lg:flex-col">
            {STEPS.map((step) => {
                const isActive = current === step.key;
                const isError = hasError(step.key);
                const title = isError ? `${step.label}(入力エラーあり)` : step.label;
                const Icon = step.icon;
                return (
                    <button
                        key={step.key}
                        type="button"
                        onClick={() => onSelect(step.key)}
                        aria-label={title}
                        title={title}
                        aria-current={isActive ? 'step' : undefined}
                        className={cn(
                            'relative grid size-10 place-items-center rounded-md transition-colors',
                            isActive
                                ? 'bg-primary/10 text-primary'
                                : 'text-muted-foreground hover:bg-accent hover:text-foreground',
                        )}
                    >
                        <Icon className="size-5" />
                        {isError && (
                            <span
                                aria-hidden
                                className="absolute right-1.5 top-1.5 size-2 rounded-full bg-destructive"
                            />
                        )}
                    </button>
                );
            })}
        </nav>
    );
};

export default FormStepper;
