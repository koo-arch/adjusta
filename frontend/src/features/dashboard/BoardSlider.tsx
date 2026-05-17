import React from 'react';
import Slider from '@/components/Slider';

interface BoardSliderProps {
    children: React.ReactNode;
}

const BoardSlider: React.FC<BoardSliderProps> = ({ children }) => {
    return (
        <div>
            <Slider
                perPage={5}
                perMove={1}
                focus={0}
                gap="1rem"
                padding="2rem"
                breakpoints={{
                    1024: {
                        perPage: 4
                    },
                    768: {
                        perPage: 3
                    },
                    640: {
                        perPage: 2
                    }
                }}
                omitEnd
                leftAlignIfFew
            >
                {children}
            </Slider>
        </div>
    )
}

export default BoardSlider;