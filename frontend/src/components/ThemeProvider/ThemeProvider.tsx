'use client'
import React from 'react';
import { ThemeProvider as NextThemesProvider } from 'next-themes';
import type { ThemeProviderProps } from 'next-themes/dist/types';

const ThemeProvider: React.FC<ThemeProviderProps> = (props) => {
    return (
       <NextThemesProvider {...props} >
            {props.children}
       </NextThemesProvider>
    )
}

export default ThemeProvider;