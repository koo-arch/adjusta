export interface AcccountsEvents {
    accountID : string;
    email : string;
    events : Event[];
}

interface Event {
    id : string;
    summary : string;
    color : string;
    start : string;
    end : string;
}