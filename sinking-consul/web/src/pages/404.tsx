import React from "react";
import {Body} from "@/components";
import { Button, Result} from "antd";
import {createStyles} from "antd-style";
import {historyPush} from "@/utils/route";

const useStyles = createStyles(({token}): any => {
    return {
        body: {
            width: "100%",
            height: "100vh",
            paddingTop: "10vh",
        },
    };
});

export default () => {
    const {styles: {body}} = useStyles();
    return (
        <Body className={body} breadCrumb={false}>
            <Result
                status="404"
                title="404"
                subTitle="对不起, 您访问的页面不存在."
                extra={<Button type="primary" onClick={() => {
                    historyPush("user.index");
                }}>返回首页</Button>}
            />
        </Body>
    );
};