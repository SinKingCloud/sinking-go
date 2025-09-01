import React from 'react';
import {theme} from "antd";
import {createStyles} from "antd-style";

export type TitlePlacement = 'left' | 'right';
export type TitleSize = 'small' | 'normal' | 'larger';

export type TitleProps = {
    children?: any;//子内容
    placement?: TitlePlacement;//位置
    size?: TitleSize;//大小
    open?: boolean;//显示状态
    width?: number;//宽度
    radius?: number;//圆角
    space?: number;//相隔
};


const useStyles = createStyles(({isDarkMode, token}): any => {
    return {
        center: {
            display: "flex",
            alignItems: "center",
            justifyContent: "start"
        },
    };
});

const Title: React.FC<TitleProps> = ({...props}) => {
    const {
        placement = "left",
        size = "normal",
        open = true,
        width = 5,
        height = 0,
        radius = 5,
        space = 5,
    } = props;

    const {token} = theme?.useToken();

    const getHeight = (h, s) => {
        if (h > 0) {
            return h;
        }
        const sizeMap = {
            small: 15,
            normal: 20,
            larger: 25
        }
        return sizeMap?.[s] || 20;
    }

    const {styles: {center}} = useStyles();

    const left = <div className={center}>
        {open && <span style={{
            width: width,
            height: getHeight(height, size),
            borderTopRightRadius: radius,
            borderBottomRightRadius: radius,
            marginRight: space,
            backgroundColor: token?.colorPrimary,
        }}/>}
        {props?.children}
    </div>;

    const right = <div className={center}>
        {props?.children}
        {open && <span style={{
            width: width,
            height: getHeight(height, size),
            borderTopLeftRadius: radius,
            borderBottomLeftRadius: radius,
            marginLeft: space,
            backgroundColor: token?.colorPrimary
        }}/>}
    </div>;

    return (
        <>
            {placement == 'left' && left}
            {placement == 'right' && right}
        </>
    );
};

export default Title;