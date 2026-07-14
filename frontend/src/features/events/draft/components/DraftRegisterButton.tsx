'use client'
import React from 'react';
import Link from 'next/link';
import { Button } from '@/components/ui/button';
import { Plus } from 'lucide-react';

const DraftRegisterButton = () => {
    return (
        <Button variant="ghost" size="icon" asChild>
            <Link href="/events/new" aria-label="イベントを作成" title="イベントを作成">
                <Plus className="text-primary" />
            </Link>
        </Button>
    )
}

export default DraftRegisterButton;
