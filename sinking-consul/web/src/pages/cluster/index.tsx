import {Tooltip, Typography} from "antd";
import {Body, PageTable} from "sinking-antd";
import type {PageTableProps} from "sinking-antd";
import RecordTime from "@/pages/components/record-time";
import useStyles from "@/pages/components/resource-table/styles";
import {getParams} from "@/utils/page";
import {useEnum} from "@/utils/enum";
import {getClusterList} from "@/service/admin/cluster";
import {ago} from "@/utils/time";

export default () => {
    const [enumData, enumLoading]: any = useEnum("cluster");
    const {styles} = useStyles();
    const statusData: Record<string, string> = enumData?.status || {};

    const columns: PageTableProps<any>["columns"] = [
        {
            title: "集群地址", dataIndex: "address", key: "address", width: 260,
            render: (value: string) => <Typography.Text className="resource-copy"
                copyable={value ? {text: value} : false}>{value || "-"}</Typography.Text>,
        },
        {
            title: "在线状态", dataIndex: "status", key: "status", width: 120,
            render: (value: number) => <span className={`resource-state is-${String(value) === "0" ? "success" : "error"}`}>
                <i/>{statusData[String(value)] || "未知状态"}
            </span>,
        },
        {
            title: "最后心跳", dataIndex: "last_heart", key: "last_heart", width: 180,
            render: (value: number) => <Tooltip title={value ? new Date(value * 1000).toLocaleString("zh-CN") : "-"}>
                <span className="resource-time">{value ? ago(new Date(value * 1000).toLocaleString("zh-CN")) : "-"}</span>
            </Tooltip>,
        },
        {
            title: "相关时间", dataIndex: "create_time", key: "create_time", width: 220, sorter: true,
            render: (_, record) => <RecordTime createTime={record.create_time} updateTime={record.update_time}/>,
        },
    ];

    return <Body loading={enumLoading}>
        <PageTable<any>
            ariaLabel="集群管理"
            className={styles.table}
            columns={columns}
            rowKey="address"
            request={async (params, sort) => {
                const response = await getClusterList({body: getParams({...params}, sort)});
                if (response?.code !== 200) throw new Error(response?.message || "获取集群列表失败");
                return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
            }}
            defaultPageSize={10}
            paginationAffix={{offsetBottom: 15}}
            hero={{eyebrow: "CLUSTER MANAGER", title: "集群管理"}}
            search={{
                placeholder: "搜索集群", maxLength: 200,
                defaultField: "keyword",
                fields: [
                    {value: "keyword", label: "全部", placeholder: "搜索集群地址"},
                    {value: "address", label: "地址", placeholder: "搜索集群地址"},
                ],
            }}
            filters={[
                {type: "select", name: "status", label: "在线状态", icon: "ApiOutlined", allLabel: "全部状态", valueEnum: statusData},
                {type: "date", name: "create_time", label: "创建时间筛选"},
            ]}
            toolbar={{refresh: {ariaLabel: "刷新集群列表"}}}
        />
    </Body>;
};
