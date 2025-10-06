import React from 'react';
import {createStyles} from "antd-style";

export type TitlePlacement = 'left' | 'right';
export type TitleSize = 'small' | 'normal' | 'larger';

export type TitleProps = {
    children?: any;//子内容
    placement?: TitlePlacement;//位置
    size?: TitleSize;//大小
    open?: boolean;//显示状态
    width?: number;//宽度
    height?: number;//高度
    radius?: number;//圆角
    space?: number;//相隔
};

const Title: React.FC<TitleProps> = ({...props}) => {
    const {
        placement = "left",
        size = "normal" as TitleSize,
        open = true,
        width = 4,
        height = 0,
        radius = -1,
        space = 7,
    } = props;

    const getHeight = (h: number, s: TitleSize) => {
        if (h > 0) {
            return h;
        }
        const sizeMap: Record<TitleSize, number> = {
            small: 15,
            normal: 20,
            larger: 25
        }
        return sizeMap[s] || 20;
    }

    const useStyles = createStyles(({token}): any => {
        return {
            center: {
                display: "flex",
                alignItems: "center",
                justifyContent: "start"
            },
            leftBar: {
                width: width,
                height: getHeight(height, size),
                borderTopRightRadius: radius < 0 ? token?.borderRadius : radius,
                borderBottomRightRadius: radius < 0 ? token?.borderRadius : radius,
                marginRight: space,
                backgroundColor: token?.colorPrimary,
            },
            rightBar: {
                width: width,
                height: getHeight(height, size),
                borderTopLeftRadius: radius < 0 ? token?.borderRadius : radius,
                borderBottomLeftRadius: radius < 0 ? token?.borderRadius : radius,
                marginLeft: space,
                backgroundColor: token?.colorPrimary,
            }
        };
    });

    const {styles: {center, leftBar, rightBar}} = useStyles();

    const left = <div className={center}>
        {open && <span className={leftBar}/>}
        {props?.children}
    </div>;

    const right = <div className={center}>
        {props?.children}
        {open && <span className={rightBar}/>}
    </div>;

    return (
        <>
            {placement == 'left' && left}
            {placement == 'right' && right}
        </>
    );
};

export default Title;