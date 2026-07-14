import React from 'react';
import Link from 'next/link';
import { Card } from '@/components/ui/card';
import StatusBadge from '@/components/common/StatusBadge/StatusBadge';
import { formatJaDateSpan } from '@/lib/date/format';
import { Calendar } from 'lucide-react';

interface MiniEventCardProps {
    title: string;
    start: Date;
    end: Date;
    needs_attention?: boolean;
    href: string;
}

// ダッシュボードの右パネル用のコンパクトな行カード
const MiniEventCard: React.FC<MiniEventCardProps> = ({ title, start, end, needs_attention, href }) => {
    return (
        <Link
            href={href}
            className="block rounded-lg focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
        >
            <Card className="px-3 py-2 transition-shadow hover:shadow-md">
                <div className="flex items-center justify-between gap-2">
                    <p className="min-w-0 truncate text-sm font-medium text-foreground">{title}</p>
                    {needs_attention && (
                        <div className="shrink-0">
                            <StatusBadge label="要対応" color="red" textSize="sm" dotSize="sm" />
                        </div>
                    )}
                </div>
                <div className="mt-1 flex items-center gap-1.5 text-xs text-muted-foreground">
                    <Calendar className="size-3.5 shrink-0" />
                    {formatJaDateSpan(start, end)}
                </div>
            </Card>
        </Link>
    )
}

export default MiniEventCard;
