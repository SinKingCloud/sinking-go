import React, {useMemo} from 'react';
import {createStyles} from "antd-style";

export type AnimationProps = {
    animate?: string;
    style?: React.CSSProperties;
    children?: React.ReactNode;
};

export const Animate = {
    None: "",
    FadeUp: "fadeInUp",
    FadeDown: "fadeInDown",
} as const;

const useStyles = createStyles(({css}): any => ({
    fadeInUp: css`
        @keyframes fadeInUp {
            0% {
                opacity: 0;
                transform: translateY(20px);
            }
            100% {
                opacity: 1;
                transform: translateY(0);
            }
        }
        animation-name: fadeInUp;
        animation-duration: 0.5s;
    `,
    fadeInDown: css`
        @keyframes fadeInDown {
            0% {
                opacity: 0;
                transform: translateY(-20px);
            }
            100% {
                opacity: 1;
                transform: translateY(0);
            }
        }
        animation-name: fadeInDown;
        animation-duration: 0.5s;
    `,
}));

const Animation: React.FC<AnimationProps> = React.memo(({animate, style, children}) => {
    const {styles} = useStyles();

    const animationClass = useMemo(() => {
        return animate && animate !== Animate.None ? styles?.[animate] : null;
    }, [animate, styles]);

    return (
        <div className={animationClass} style={style}>
            {children}
        </div>
    );
});

export default Animation;