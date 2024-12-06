import React from 'react';
import { insertLineBreakAtMarker } from '@/lib/text/breakWords';

interface WrapTextProps {
    text: string;
    maxLength: number;
    marker: string;
}

const WrapText: React.FC<WrapTextProps> = ({ text, maxLength, marker }) => {
    const brokenText = insertLineBreakAtMarker(text, marker, maxLength);
    return (
        <div>
            {brokenText.split('\n').map((word, index) => (
                <React.Fragment key={index}>
                    {word}
                    <br />
                </React.Fragment>
            ))}
        </div>
    )
}

export default WrapText;