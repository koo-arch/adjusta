'use client'
import React from 'react';
import Link from 'next/link';
import Image from 'next/image';
import UserButton from '../UserButton';
import { usePathname } from 'next/navigation';
import { Disclosure, DisclosureButton, DisclosurePanel } from '@headlessui/react';
import { Bars3Icon, XMarkIcon } from '@heroicons/react/20/solid';

const navigation = [
    { name: 'Home', href: '/' },
    { name: 'About', href: '/about' },
    { name: 'Contact', href: '/contact' },
]

const classNames = (...classes: string[]) => {
    return classes.filter(Boolean).join(' ')
}

const Header: React.FC = () => {
    const pathname = usePathname();

    return (
        <Disclosure as="nav" className="lg:pb-12">
            {({ open }) => (
                <>
                    <div className="mx-auto max-w-screen-2xl px-4 md:px-8">
                        <div className='relative flex items-center justify-between py-4 md:py-8'>
                            <div className="absolute inset-y-0 left-0 flex items-center md:hidden">

                                <DisclosureButton className="relative inline-flex items-center justify-center rounded-md p-2">
                                    <span className="absolute -inset-0.5"></span>
                                    <span className="sr-only">Open main menu</span>
                                    {open ? (
                                        <XMarkIcon className="h-6 w-6" aria-hidden="true" />
                                    ) : (
                                        <Bars3Icon className="h-6 w-6" aria-hidden="true" />
                                    )}
                                </DisclosureButton>
                            </div>
                            <div className="flex flex-1 items-center justify-center md:items-stretch md:justify-start">
                                <div className="flex flex-shrink-0 items-center">
                                    <Image
                                        className="h-6 w-auto"
                                        height={24}
                                        width={24}
                                        src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=500"
                                        alt="Your Company"
                                    />
                                </div>
                                <div className="hidden md:ml-6 md:block">
                                    <div className="flex gap-12">
                                        {navigation.map((item) => (
                                            <Link key={item.name} href={item.href} className={classNames(
                                                item.href === pathname ? 'text-indigo-500' : 'hover:text-indigo-500 active:text-indigo-700',
                                                "text-lg font-semibold transiton duration-100"
                                            )}
                                            >
                                                    {item.name}
                                            </Link>
                                        ))}
                                    </div>
                                </div>
                            </div>
                            <UserButton classNames={classNames} />
                        </div>
                    </div>

                    <DisclosurePanel 
                        transition
                        className="sm:hidden origin-top transition duration-200 ease-out data-[closed]:-translate-x-6 data-[closed]:opacity-0"
                    >
                        <div className="space-y-1 px-4 pb-3 pt-2">
                            {navigation.map((item) => (
                                <DisclosureButton
                                    key={item.name}
                                    as={Link}
                                    href={item.href}
                                    className={classNames(
                                        item.href === pathname ? 'text-indigo-500' : 'hover:text-indigo-500 active:text-indigo-700',
                                        "block text-lg font-semibold"
                                    )}
                                >
                                    {item.name}
                                </DisclosureButton>
                            ))}
                        </div>
                    </DisclosurePanel>
                </>
            )}
        </Disclosure>
    )
}

export default Header;