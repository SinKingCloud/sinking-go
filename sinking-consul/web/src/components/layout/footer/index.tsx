import React from "react";

export type FooterProps = {
    children?: any;//子内容
};

const Footer: React.FC<FooterProps> = ({...props}) => {
    return <>
        {props?.children}
    </>
}

export default Footer