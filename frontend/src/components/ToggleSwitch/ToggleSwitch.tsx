'use client';
import React from 'react';
import { cva, cx, type VariantProps } from 'class-variance-authority';

const toggleSwitchStyles = cva(
    'relative inline-flex items-center rounded-full transition-colors cursor-pointer bg-gray-200 dark:bg-gray-700',
    {
        variants: {
            size: {
                sm: 'w-9 h-5',
                md: 'w-11 h-6',
                lg: 'w-14 h-7',
            },
            color: {
                primary: 'peer-checked:bg-primary-600 peer-focus:ring-primary-300 dark:peer-focus:ring-primary-800',
                secondary: 'peer-checked:bg-secondary-600 peer-focus:ring-secondary-300 dark:peer-focus:ring-secondary-800',
                danger: 'peer-checked:bg-red-600 peer-focus:ring-red-300 dark:peer-focus:ring-red-800',
                warning: 'peer-checked:bg-yellow-600 peer-focus:ring-yellow-300 dark:peer-focus:ring-yellow-800',
                success: 'peer-checked:bg-green-600 peer-focus:ring-green-300 dark:peer-focus:ring-green-800',
                indigo: 'peer-checked:bg-indigo-600 peer-focus:ring-indigo-300 dark:peer-focus:ring-indigo-800',
            },
        },
        defaultVariants: {
            size: 'md',
            color: 'indigo',
        },
    }
);

const knobStyles = cva(
    'after:content-[""] after:absolute after:bg-white after:border-gray-300 after:border after:rounded-full after:transition-all dark:border-gray-600',
    {
        variants: {
            size: {
                sm: 'after:w-4 after:h-4 after:top-0.5 after:start-[2px]',
                md: 'after:w-5 after:h-5 after:top-0.5 after:start-[2px]',
                lg: 'after:w-6 after:h-6 after:top-0.5 after:start-[3px]',
            },
            checked: {
                true: 'after:translate-x-full',
                false: '',
            },
        },
        defaultVariants: {
            size: 'md',
            checked: false,
        },
    }
);

interface ToggleSwitchProps extends VariantProps<typeof toggleSwitchStyles>, VariantProps<typeof knobStyles> {
    checked: boolean;
    onChange: (checked: boolean) => void;
    label?: string;
    size?: 'sm' | 'md' | 'lg';
    color?: 'primary' | 'secondary' | 'danger' | 'warning' | 'success' | 'indigo';
}

const ToggleSwitch: React.FC<ToggleSwitchProps> = ({
    checked,
    onChange,
    label,
    size,
    color
}) => {
    return (
        <label className="inline-flex items-center space-x-2 cursor-pointer">
            {label && <span className="text-sm font-medium">{label}</span>}
            <input
                type="checkbox"
                className="sr-only peer"
                checked={checked}
                onChange={(e) => onChange(e.target.checked)}
            />
            <div className={toggleSwitchStyles({ size, color })}>
                <span className={knobStyles({ size, checked })} />
            </div>
        </label>
    );
};

export default ToggleSwitch;