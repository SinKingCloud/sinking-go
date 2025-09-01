import React from 'react';
import {createStyles} from "antd-style";

export type AnimationProps = {
    animate?: any;//主题
    style?: any;
    children?: any;//子内容
};

export const Animate = {
    None: "",
    FadeUp: "fadeInUp",
    FadeDown: "fadeInDown",
};

const useStyles = createStyles(({css}): any => {
    return {
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
    };
});

const Animation: React.FC<AnimationProps> = ({...props}) => {
    const {styles} = useStyles();
    return (
        <div className={props?.animate && props?.animate != Animate.None ? styles?.[props?.animate] : null}
             style={props?.style}>
            {props?.children}
        </div>
    );
};

export default Animation;