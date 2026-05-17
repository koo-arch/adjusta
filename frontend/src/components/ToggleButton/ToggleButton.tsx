'use client'
import React, { useState } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

const toggleButtonStyles = cva('font-medium text-center w-full', {
    variants: {
        variant: {
            solid: 'bg-indigo-500 text-white hover:bg-indigo-600 active:bg-indigo-700',
            outline: 'border border-indigo-500 text-indigo-500 hover:bg-indigo-500 hover:text-white active:bg-indigo-600',
        },
        size: {
            sm: "px-2 py-1 text-sm",
            md: "px-3 py-1.5 text-base",
            lg: "px-4 py-2 text-lg",
            xl: "px-5 py-2.5 text-xl",
        },
        position: {
            left: 'rounded-l-md', // 左端のボタン
            right: 'rounded-r-md', // 右端のボタン
            middle: '', // 中央のボタン（丸めなし）
        },
    },
    defaultVariants: {
        variant: 'outline',
        size: 'md',
        position: 'middle',
    },
});

interface ToggleButtonProps<T> extends VariantProps<typeof toggleButtonStyles> {
    options: T[];
    selected: T;
    onToggle: (selected: T) => void;
    renderLabel: (option: T) => string;
}

const ToggleButton = <T extends unknown>({
    options,
    onToggle,
    renderLabel,
    size,
}: ToggleButtonProps<T>) => {
    const [selected, setSelected] = useState(options[0]);

    const handleToggle = (option: T) => {
        setSelected(option);
        onToggle(option);
    }
    return (
        <div className="inline-flex w-full shadow-sm">
            {options.map((option, index) => {
                const position =
                    index === 0 ? 'left' : index === options.length - 1 ? 'right' : 'middle';
                return (
                    <button
                        key={index}
                        className={`${toggleButtonStyles({
                            variant: selected === option ? 'solid' : 'outline',
                            size,
                            position,
                        })} flez-grow`}
                        onClick={() => handleToggle(option)}
                    >
                        {renderLabel(option)}
                    </button>
                );
            })}
        </div>
    );
};


export default ToggleButton;