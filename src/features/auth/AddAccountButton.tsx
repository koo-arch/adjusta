'use client'
import React from 'react';
import Button from '@/components/Button';
import Image from 'next/image';

const AddAccountButton: React.FC = () => {
    return (
        <Button
            shape='full'
            variant='outline'
            intent='clear'
            size='md'
            to="http://localhost:8080/api/google/add-account"
            startIcon={
                <Image
                    src="https://www.svgrepo.com/show/475656/google-color.svg"
                    loading="lazy"
                    alt="google logo"
                    height={24}
                    width={24}
                ></Image>
            }
        >
            Add Account
        </Button>
    )
}

export default AddAccountButton;