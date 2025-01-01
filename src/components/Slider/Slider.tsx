'use client'
import React from 'react';
import { Splide, SplideTrack, SplideSlide, SplideProps } from 'react-splide-ts';
import { Options } from '@splidejs/splide';
import 'react-splide-ts/css';

type SliderProps = SplideProps & Options & {
    leftAlignIfFew? : boolean;
}

const Slider: React.FC<SliderProps> = ({
    children,
    type = 'slide',
    rewind = false,
    pagination = false,
    gap = '1rem',
    focus = 0,
    perPage = 1,
    perMove = 1,
    padding,
    breakpoints,
    omitEnd = true,
    leftAlignIfFew = false,
    ...props
}) => {
    const childrenArray = React.Children.toArray(children);
    const isFewSlides = childrenArray.length <= perPage;

    return (
        <div>
            <Splide
                hasTrack={false} 
                options={{
                    type,
                    rewind,
                    pagination,
                    gap,
                    perPage,
                    perMove,
                    focus: isFewSlides && leftAlignIfFew ? 0 : focus,
                    breakpoints,
                    padding,
                    omitEnd,
                    ...props
                }}
                {...props}
            >
                <SplideTrack>
                    {React.Children.map(children, (child, index) => (
                        <SplideSlide key={index}>{child}</SplideSlide>
                    ))}
                </SplideTrack>
            </Splide>
        </div>
    );
}

export default Slider;