import React, { forwardRef } from 'react';
import { Field, Label, Description, Textarea } from '@headlessui/react';
import { cva, type VariantProps } from 'class-variance-authority';

const textareaStyle = cva('block w-full mt-1 border', {
    variants: {
        areaSize: {
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
        areaSize: 'md',
        shape: 'md',
        error: false,
    },
});

interface TextAreaProps extends React.TextareaHTMLAttributes<HTMLTextAreaElement>, VariantProps<typeof textareaStyle> {
    label?: string;
    description?: string;
    placeholder?: string;
    helperText?: string;
    rows?: number;
}

const TextArea = forwardRef<HTMLTextAreaElement, TextAreaProps>(
    ({ label, description, placeholder, shape, areaSize, error, helperText, rows, ...rest }, ref) => {
        return (
            <div>
                {label &&
                    <label className="font-medium text-md">{label}</label>
                }
                {description && (
                    <p className="text-gray-500 text-sm">{description}</p>
                )}
                <textarea
                    ref={ref}
                    placeholder={placeholder}
                    className={`${textareaStyle({ shape, areaSize, error })} focus:outline-none focus:ring-2`}
                    {...rest}
                    rows={rows || 3}
                />
                {helperText && (
                    <p className={`mt-1 text-sm ${error ? 'text-red-500' : 'text-gray-500'}`}>
                        {helperText}
                    </p>
                )}
            </div>
        );
    }
);

export default TextArea;