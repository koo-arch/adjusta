'use client'
import React, { useRef } from 'react';
import { Input } from '@/components/ui/input';
import { cn } from '@/lib/utils';
import { Search, X } from 'lucide-react';

interface EventSearchFormProps {
    defaultValue: string;
    onSearch: (value: string) => void;
}

const EventSearchForm: React.FC<EventSearchFormProps> = ({ defaultValue, onSearch }) => {
    const inputRef = useRef<HTMLInputElement>(null);

    const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        onSearch(inputRef.current?.value.trim() ?? '');
    };

    const handleKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
        // IME 変換確定の Enter でフォーム送信しない
        if (event.key === 'Enter' && event.nativeEvent.isComposing) {
            event.preventDefault();
        }
    };

    const handleClear = () => {
        if (inputRef.current) {
            inputRef.current.value = '';
        }
        onSearch('');
    };

    // 入力は 1 つだけなので Enter の暗黙送信で検索が発火する(検索ボタンは置かない)
    return (
        <form role="search" onSubmit={handleSubmit} className="w-full md:w-auto">
            <div className="relative md:w-64">
                <Search
                    className="pointer-events-none absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground"
                    aria-hidden
                />
                {/* 非制御 + key で「URL の確定値が変わったときだけ」入力を振り直す */}
                <Input
                    key={defaultValue}
                    ref={inputRef}
                    type="text"
                    name="title"
                    defaultValue={defaultValue}
                    placeholder="タイトルで検索"
                    aria-label="タイトルで検索"
                    onKeyDown={handleKeyDown}
                    className={cn('h-10 pl-9', defaultValue !== '' && 'pr-9')}
                />
                {defaultValue !== '' && (
                    <button
                        type="button"
                        onClick={handleClear}
                        aria-label="検索条件をクリア"
                        className="absolute right-2 top-1/2 -translate-y-1/2 rounded-sm p-1.5 text-muted-foreground transition-colors hover:text-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring"
                    >
                        <X className="size-4" />
                    </button>
                )}
            </div>
        </form>
    );
};

export default EventSearchForm;
