'use client'
import React from 'react';
import { Description, Dialog, DialogPanel, DialogTitle, DialogBackdrop } from '@headlessui/react';
import Button from '@/components/Button';

interface ModalProps {
    isOpen: boolean;
    onClose: () => void;
    title?: string;
    description?: string;
    children?: React.ReactNode;
    actions?: React.ReactNode;
    hideCloseButton?: boolean;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children, description, actions, hideCloseButton }) => {
    return (
        <Dialog open={isOpen} className="relative z-10 focus:outline-none" onClose={onClose}>
            <DialogBackdrop 
                transition
                className="fixed inset-0 bg-black/30 duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0"
            />
            <div className="fixed inset-0 z-10 w-screen overflow-y-auto">
                <div className="flex min-h-full items-center justify-center p-4">
                    <DialogPanel
                        transition
                        className="w-full max-w-md rounded-xl bg-white p-6 duration-300 ease-out data-[closed]:transform-[scale(95%)] data-[closed]:opacity-0"
                    >
                        <DialogTitle as="h3" className="text-inherit font-medium text-lg">{title}</DialogTitle>
                        {description && (
                            <Description className="text-sm text-gray-500">
                                {description}
                            </Description>
                        )}

                        {children && (
                            <div className="mt-4">
                                {children}
                            </div>
                        )}

                        {actions && (
                            <div className="mt-4 flex justify-end space-x-2">
                                {!hideCloseButton && (
                                    <Button
                                        onClick={onClose}
                                        variant="outline"
                                        intent="clear"
                                    >
                                        Close
                                    </Button>
                                )}
                                {actions}
                            </div>
                        )}

                    </DialogPanel>
                </div>
            </div>
        </Dialog>
    );
}

export default Modal;