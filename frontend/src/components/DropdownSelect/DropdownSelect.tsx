'use client'
import React, { useState } from 'react';
import { Label, Listbox, ListboxButton, ListboxOptions, ListboxOption } from '@headlessui/react';
import { CheckIcon, ChevronUpDownIcon } from '@heroicons/react/20/solid';
import { cva, type VariantProps } from 'class-variance-authority';

const listboxStyle = cva('block w-full mt-1 border', {
    variants: {
        selectSize: {
            sm: 'px-2 py-1 text-sm',
            md: 'px-3 py-1.5 text-base',
            lg: 'px-4 py-2 text-lg',
            xl: 'px-5 py-2.5 text-xl',
        },
        shape: {
            rounded: 'rounded',
            md: 'rounded-md',
            lg: 'rounded-lg',
            full: 'rounded-full',
        },
        error: {
            true: 'border-red-500 focus:ring-red-500',
            false: 'focus:ring-indigo-500',
        },
    },
    defaultVariants: {
        selectSize: 'md',
        shape: 'md',
        error: false,
    },
});

interface DropdownSelectProps<T> extends VariantProps<typeof listboxStyle> {
    label?: string;
    options: T[];
    onChange: (item: T | null) => void;
    renderLabel: (item: T | null) => React.ReactNode;
    defaultSelected?: T;
    error?: boolean;
    helperText?: string;
    placeholder?: string;
}

const DropdownSelect = <T extends unknown>({
    label,
    options,
    onChange,
    renderLabel,
    defaultSelected,
    selectSize,
    shape,
    helperText,
    error,
    placeholder = '未選択', // プレースホルダーのデフォルト
}: DropdownSelectProps<T>) => {
    const [selected, setSelected] = useState<T | null>(defaultSelected || null);

    const handleSelect = (item: T | null) => {
        setSelected(item);
        onChange(item);
    }

    return (
        <div>
            <Listbox value={selected} onChange={handleSelect}>
                {label && <Label className="font-medium text-md">{label}</Label>}
                <div className="relative mt-2">
                    <ListboxButton
                        className={`${listboxStyle({ selectSize, shape, error })} cursor-default rounded-md bg-white py-1.5 pl-3 pr-10 text-left shadow-sm focus:outline-none focus:ring-2`}
                    >
                        <span className="block truncate">{selected ? renderLabel(selected) : placeholder}</span>
                        <span className="pointer-events-none absolute inset-y-0 right-0 ml-3 flex items-center pr-2">
                            <ChevronUpDownIcon aria-hidden="true" className="h-5 w-5 text-gray-400" />
                        </span>
                    </ListboxButton>

                    <ListboxOptions
                        transition
                        className="absolute z-10 mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black ring-opacity-5 transition duration-100 ease-out [--anchor-gap:var(--spacing-1)] focus:outline-none data-[closed]:scale-95 data-[closed]:opacity-0 sm:text-sm"
                    >
                        {/* 未選択オプション */}
                        <ListboxOption
                            key="null"
                            value={null}
                            className={({ focus }) =>
                                `relative cursor-default select-none py-2 pl-3 pr-9 ${focus ? 'bg-indigo-600 text-white' : 'text-gray-900'
                                }`
                            }
                        >
                            {({ focus }) => (
                                <>
                                    <span className={`block truncate ${focus ? 'font-semibold' : 'font-normal'}`}>
                                        {placeholder}
                                    </span>
                                    {selected === null && (
                                        <span className="absolute inset-y-0 right-0 flex items-center pr-4">
                                            <CheckIcon className="h-5 w-5 text-indigo-600" aria-hidden="true" />
                                        </span>
                                    )}
                                </>
                            )}
                        </ListboxOption>

                        {/* 通常のオプション */}
                        {options.map((option, index) => (
                            <ListboxOption
                                key={index}
                                value={option}
                                className={({ focus }) =>
                                    `relative cursor-default select-none py-2 pl-3 pr-9 ${focus ? 'bg-indigo-600 text-white' : 'text-gray-900'
                                    }`
                                }
                            >
                                {({ focus }) => (
                                    <>
                                        <span className={`block truncate ${focus ? 'font-semibold' : 'font-normal'}`}>
                                            {renderLabel(option)}
                                        </span>
                                        {selected === option && (
                                            <span className="absolute inset-y-0 right-0 flex items-center pr-4">
                                                <CheckIcon className="h-5 w-5 text-indigo-600" aria-hidden="true" />
                                            </span>
                                        )}
                                    </>
                                )}
                            </ListboxOption>
                        ))}
                    </ListboxOptions>
                </div>
            </Listbox>
            {helperText && (
                <p className={`mt-1 text-sm ${error ? 'text-red-500' : 'text-gray-500'}`}>
                    {helperText}
                </p>
            )}
        </div>
    );
};

export default DropdownSelect;