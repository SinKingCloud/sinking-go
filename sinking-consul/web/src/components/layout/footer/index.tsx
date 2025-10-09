import React from "react";

export type FooterProps = {
    children?: any;//子内容
};

const Footer: React.FC<FooterProps> = React.memo(({children}) => {
    return <>{children}</>
});

export default Footer;