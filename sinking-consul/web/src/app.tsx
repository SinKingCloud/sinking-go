import {App} from "antd";
import React from "react";

/**
 * 挂载静态方法
 * @param container 容器
 */
export function rootContainer(container: React.ReactElement) {
    return (
        <App>
            {container}
        </App>
    )
}