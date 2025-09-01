import React, {forwardRef} from "react";
import Body from "@/components/layout/body";
import Footer from "@/components/layout/footer";
import Header from "@/components/layout/header";
import Sider from "@/components/layout/sider";
import 'dayjs/locale/zh-cn';
import {Theme} from "@/components";

import SkLayout, {LayoutProps, SinKingRef} from "@/components/layout/sinking";

const Layout: React.FC<LayoutProps> = forwardRef<SinKingRef>((props: any, ref) => {
    return (<Theme>
        <SkLayout ref={ref} {...props}/>
    </Theme>);
});

export {
    Body,
    Footer,
    Header,
    Sider,
}

export default Layout;