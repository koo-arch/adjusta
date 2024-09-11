export interface EventDraftDetail {
    id: string;
    title: string;
    description: string;
    allDay: boolean;
    location: string;
    url: string;
    status: string;
    proposed_dates: ProposedDate[];
}

export interface ProposedDate {
    id: string;
    start_date: Date;
    end_date: Date;
    priority: number;
    is_finalized: boolean;
}