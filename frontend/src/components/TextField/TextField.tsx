import React, { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';

const inputStyle = cva('block w-full mt-1 border', {
    variants: {
        inputSize: {
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
        inputSize: 'md',
        shape: 'md',
        error: false,
    },
});



interface TextFieldProps extends React.InputHTMLAttributes<HTMLInputElement>, VariantProps<typeof inputStyle> {
    label?: string;
    description?: string;
    placeholder?: string;
    helperText?: string;
}

const TextField = forwardRef<HTMLInputElement, TextFieldProps>(
    ({ label, description, placeholder, shape, inputSize, error, helperText, ...rest }, ref) => {
        return (
            <div>
                {label &&
                    <label className="font-medium text-md">{label}</label>
                }
                {description && (
                    <p className="text-gray-500 text-sm">{description}</p>
                )}
                <input
                    ref={ref}
                    placeholder={placeholder}
                    className={`${inputStyle( { shape, inputSize, error })} focus:outline-none focus:ring-2`}
                    {...rest}
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

TextField.displayName = 'TextField';

export default TextField;