export const insertLineBreakAtMarker = (text: string, marker: string, maxLength: number): string => {
    if (text.length > maxLength && text.includes(marker)) {
        return text.replace(marker, `${marker}\n`);
    }
    return text;
}