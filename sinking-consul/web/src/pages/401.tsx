import React from "react";
import {Body} from "@/components";
import {Button, Result} from "antd";
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
                status="403"
                title="403"
                subTitle="您无权限访问此页面,请联系管理员"
                extra={<Button type="primary" onClick={() => {
                    historyPush("user.index");
                }}>返回首页</Button>}
            />
        </Body>
    );
};