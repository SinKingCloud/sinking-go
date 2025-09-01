import React from "react";
import {createStyles} from "antd-style";

const useStyles = createStyles(({responsive, css}): any => {
    return {
        box: {
            display: "flex",
            alignItems: "center",
            justifyContent: "space-between"
        },
        right: css`
            z-index: 2;
        `,
        left: css`
            z-index: 2;
        `
    };
});
export type HeaderProps = {
    right?: any;//右侧
    left?: any;//左侧
};

const Header: React.FC<HeaderProps> = (props) => {
    const {styles: {box, right, left}} = useStyles();
    return <div className={box}>
        {props?.left && <div className={left}>
            {props.left}
        </div>}
        {props?.right && <div className={right}>
            {props.right}
        </div>}
    </div>
}

export default Header;