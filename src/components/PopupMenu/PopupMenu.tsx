'use client'
import React from 'react';
import { Menu, MenuItems, MenuItem, MenuButton } from '@headlessui/react';

interface PopupMenuProps {
    items: Array<{ label: string; onClick: () => void }>;
    position?: { top: number; left: number };
    buttonRef?: React.RefObject<HTMLButtonElement>;
}

const classNames = (...classes: string[]) => {
    return classes.filter(Boolean).join(' ');
}

const PopupMenu: React.FC<PopupMenuProps> = ({ items, position={ top: 0, left: 0 }, buttonRef }) => {

    return (
        <Menu as="div" className="absolute z-10" style={position}>
            <MenuButton ref={buttonRef}>
                <span className="sr-only">open event menu</span>
            </MenuButton>
            <MenuItems
                transition
                className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 transition focus:outline-none data-[closed]:scale-95 data-[closed]:transform data-[closed]:opacity-0 data-[enter]:duration-100 data-[leave]:duration-75 data-[enter]:ease-out data-[leave]:ease-in"
            >
                {items.map((item, index) => (
                    <MenuItem key={index}>
                        {({ focus }) => (
                            <a
                                onClick={() => {
                                    item.onClick();
                                }}
                                className={classNames(focus ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700')}
                            >
                                {item.label}
                            </a>
                        )}
                    </MenuItem>
                ))}
                <MenuItem>
                    {({ focus }) => (
                        <a
                            className={classNames(focus ? 'bg-gray-100' : '', 'block px-4 py-2 text-sm text-gray-700')}
                        >
                            キャンセル
                        </a>
                    )}
                </MenuItem>
            </MenuItems>
        </Menu>
    );
};

export default PopupMenu;