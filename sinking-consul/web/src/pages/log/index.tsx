import React, {useMemo} from "react";
import {Tooltip, Typography} from "antd";
import {Body, PageTable} from "sinking-antd";
import type {PageTableProps} from "sinking-antd";
import {getParams} from "@/utils/page";
import {useEnum} from "@/utils/enum";
import {getLog} from "@/service/admin/person";
import useStyles from "./styles";

interface LogRecord {
    id: number;
    ip: string;
    type: number;
    title: string;
    content: string;
    create_time: string;
}

const emptyEnum = {};

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const [enumData, enumLoading]: any = useEnum("log");
    const typeData: Record<string, string> = enumData?.type || emptyEnum;

    const columns = useMemo<PageTableProps<LogRecord>["columns"]>(() => [
        {
            title: "ID",
            dataIndex: "id",
            key: "id",
            width: 90,
            sorter: true,
        },
        {
            title: "操作IP",
            dataIndex: "ip",
            key: "ip",
            width: 170,
            render: (value: string) => (
                <Typography.Text className="log-copy" copyable={value ? {text: value} : false}>
                    {value || "-"}
                </Typography.Text>
            ),
        },
        {
            title: "操作类型",
            dataIndex: "type",
            key: "type",
            width: 120,
            render: (value: number) => (
                <span className={`log-type is-${value}`}>
                    <i/>
                    {typeData[String(value)] || "未知类型"}
                </span>
            ),
        },
        {
            title: "操作标题",
            dataIndex: "title",
            key: "title",
            width: 180,
            render: (value: string) => (
                <Tooltip title={value || "-"}><span className="log-text">{value || "-"}</span></Tooltip>
            ),
        },
        {
            title: "操作内容",
            dataIndex: "content",
            key: "content",
            width: 300,
            render: (value: string) => (
                <Tooltip title={value || "-"}><span className="log-text">{value || "-"}</span></Tooltip>
            ),
        },
        {
            title: "操作时间",
            dataIndex: "create_time",
            key: "create_time",
            width: 180,
            sorter: true,
            render: (value: string) => (
                <Tooltip title={value || "-"}><span className="log-time">{value || "-"}</span></Tooltip>
            ),
        },
    ], [typeData]);

    return (
        <Body loading={enumLoading}>
            <PageTable<LogRecord>
                ariaLabel="操作日志"
                className={styles.logTable}
                columns={columns}
                rowKey="id"
                request={async (params, sort) => {
                    const response = await getLog({body: getParams({...params}, sort)});
                    if (response?.code !== 200) {
                        throw new Error(response?.message || "获取列表失败");
                    }
                    return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
                }}
                defaultPageSize={10}
                paginationAffix={{offsetBottom: 15}}
                hero={{eyebrow: "OPERATION AUDIT", title: "操作日志"}}
                search={{
                    placeholder: "搜索日志",
                    ariaLabel: "搜索操作日志",
                    maxLength: 200,
                    defaultField: "keyword",
                    fields: [
                        {value: "keyword", label: "全部", placeholder: "搜索 IP、标题或内容"},
                        {value: "ip", label: "IP", placeholder: "搜索操作 IP"},
                        {value: "title", label: "标题", placeholder: "搜索标题"},
                        {value: "content", label: "内容", placeholder: "搜索内容"},
                    ],
                }}
                filters={[
                    {
                        type: "select",
                        name: "type",
                        label: "日志类型筛选",
                        icon: "TagsOutlined",
                        allLabel: "全部类型",
                        valueEnum: typeData,
                    },
                    {type: "date", name: "create_time", label: "操作时间筛选"},
                ]}
                toolbar={{refresh: {ariaLabel: "刷新操作日志列表"}}}
            />
        </Body>
    );
};
