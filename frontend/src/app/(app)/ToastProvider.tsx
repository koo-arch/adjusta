'use client'
import React from 'react';
import "react-toastify/dist/ReactToastify.css";
import { Flip, ToastContainer } from 'react-toastify';

interface ToastProviderProps {
    children: React.ReactNode;
}

const ToastProvider: React.FC<ToastProviderProps> = ({ children }) => {
    return (
        <div>
            {children}
            <ToastContainer
                position="top-right"
                autoClose={5000}
                hideProgressBar
                newestOnTop={false}
                closeOnClick
                rtl={false}
                pauseOnFocusLoss
                draggable
                pauseOnHover
                transition={Flip}
            />
        </div>
    )
}

export default ToastProvider;