import {useRef, useState} from "react";
import type {Key} from "react";
import {App, Button, Form, Select, Tooltip, Typography} from "antd";
import {Body, PageTable, ProModal, Title} from "sinking-antd";
import type {PageTableProps, PageTableRef, ProModalRef} from "sinking-antd";
import ActionDropdown from "@/pages/components/action-dropdown";
import RecordTime from "@/pages/components/record-time";
import useStyles from "@/pages/components/resource-table/styles";
import {getParams} from "@/utils/page";
import {useEnum} from "@/utils/enum";
import {deleteNode, getNodeList, updateNode} from "@/service/admin/node";
import {ago} from "@/utils/time";

export default () => {
    const [enumData, enumLoading]: any = useEnum("node");
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const tableRef = useRef<PageTableRef<any> | null>(null);
    const modalRef = useRef<ProModalRef | null>(null);
    const [form] = Form.useForm();
    const [editRecords, setEditRecords] = useState<string[]>([]);
    const [editLoading, setEditLoading] = useState(false);
    const [openRowAddress, setOpenRowAddress] = useState("");
    const onlineStatusData: Record<string, string> = enumData?.online_status || {};
    const statusData: Record<string, string> = enumData?.status || {};

    const onEdit = (addresses: Key[], record?: any) => {
        if (!addresses.length) return;
        form.resetFields();
        if (record) form.setFieldsValue({status: String(record.status)});
        setEditRecords(addresses.map(String));
        modalRef.current?.show();
    };

    const onFormFinish = async (values: {status: string}) => {
        setEditLoading(true);
        try {
            const response = await updateNode({body: {addresses: editRecords, status: values.status}});
            if (response?.code !== 200) {
                message.error(response?.message || "编辑失败");
                return;
            }
            modalRef.current?.hide();
            tableRef.current?.reload();
            tableRef.current?.clearSelectedRows();
            message.success(response.message || "编辑成功");
        } finally {
            setEditLoading(false);
        }
    };

    const onDelete = (addresses: Key[]) => {
        if (!addresses.length) return;
        modal.confirm({
            title: "删除节点",
            content: `确定删除选中的 ${addresses.length} 个节点吗？`,
            okText: "删除",
            okType: "danger",
            cancelText: "取消",
            mask: {closable: true},
            onOk: async () => {
                const response = await deleteNode({body: {addresses: addresses.map(String)}});
                if (response?.code !== 200) {
                    message.error(response?.message || "删除失败");
                    throw new Error(response?.message || "删除失败");
                }
                tableRef.current?.reload();
                tableRef.current?.clearSelectedRows();
                message.success(response.message || "删除成功");
            },
        });
    };

    const renderCopy = (value: string) => <Typography.Text className="resource-copy"
        copyable={value ? {text: value} : false}>{value || "-"}</Typography.Text>;

    const columns: PageTableProps<any>["columns"] = [
        {title: "服务分组", dataIndex: "group", key: "group", width: 150, render: renderCopy},
        {title: "服务名称", dataIndex: "name", key: "name", width: 160, render: renderCopy},
        {title: "节点地址", dataIndex: "address", key: "address", width: 210, render: renderCopy},
        {
            title: "在线状态", dataIndex: "online_status", key: "online_status", width: 100,
            render: (value: number) => <span className={`resource-state is-${String(value) === "0" ? "success" : "error"}`}>
                <i/>{onlineStatusData[String(value)] || "未知状态"}
            </span>,
        },
        {
            title: "最后心跳", dataIndex: "last_heart", key: "last_heart", width: 140,
            render: (value: number) => <Tooltip title={value ? new Date(value * 1000).toLocaleString("zh-CN") : "-"}>
                <span className="resource-time">{value ? ago(new Date(value * 1000).toLocaleString("zh-CN")) : "-"}</span>
            </Tooltip>,
        },
        {
            title: "服务状态", dataIndex: "status", key: "status", width: 100,
            render: (value: number) => <span className={`resource-state is-${String(value) === "0" ? "success" : "warning"}`}>
                <i/>{statusData[String(value)] || "未知状态"}
            </span>,
        },
        {
            title: "相关时间", dataIndex: "create_time", key: "create_time", width: 220, sorter: true,
            render: (_, record) => <RecordTime createTime={record.create_time} updateTime={record.update_time}/>,
        },
        {
            title: "操作", key: "option", width: 80, align: "center", fixed: "right", className: "action-cell",
            render: (_, record) => <ActionDropdown
                open={openRowAddress === record.address}
                onOpenChange={(open) => setOpenRowAddress(open ? record.address : "")}
                menu={{items: [
                    {key: "edit", label: "编辑", onClick: () => onEdit([record.address], record)},
                    {key: "delete", label: "删除", danger: true, onClick: () => onDelete([record.address])},
                ]}}>
                <Button size="small" aria-label="节点操作">操作</Button>
            </ActionDropdown>,
        },
    ];

    return <Body loading={enumLoading}>
        <PageTable<any>
            ref={tableRef}
            ariaLabel="服务节点"
            className={styles.table}
            columns={columns}
            rowKey="address"
            rowClassName={(record) => openRowAddress === record.address ? "action-menu-open" : ""}
            request={async (params, sort) => {
                const response = await getNodeList({body: getParams({...params}, sort)});
                if (response?.code !== 200) throw new Error(response?.message || "获取节点列表失败");
                return {data: response.data?.list, total: response.data?.total ?? 0, success: true};
            }}
            defaultPageSize={10}
            paginationAffix={{offsetBottom: 15}}
            hero={{eyebrow: "SERVICE NODES", title: "服务节点"}}
            search={{
                placeholder: "搜索服务节点", maxLength: 200,
                defaultField: "keyword",
                fields: [
                    {value: "keyword", label: "全部", placeholder: "搜索分组、名称或地址"},
                    {value: "group", label: "分组", placeholder: "搜索服务分组"},
                    {value: "name", label: "名称", placeholder: "搜索服务名称"},
                    {value: "address", label: "地址", placeholder: "搜索节点地址"},
                ],
            }}
            filters={[
                {type: "select", name: "online_status", label: "在线状态", icon: "ApiOutlined", allLabel: "全部状态", valueEnum: onlineStatusData, showSelectedLabel: false},
                {type: "select", name: "status", label: "服务状态", icon: "FlagOutlined", allLabel: "全部状态", valueEnum: statusData, showSelectedLabel: false},
                {type: "date", name: "create_time", label: "创建时间", showSelectedLabel: false},
            ]}
            toolbar={{refresh: {ariaLabel: "刷新服务节点列表"}}}
            rowSelection={{preserveSelectedRowKeys: true, selections: true, actions: (keys) => [
                {key: "edit", type: "button", label: "批量编辑", icon: "EditOutlined", onClick: () => onEdit(keys)},
                {key: "delete", type: "button", label: "批量删除", icon: "DeleteOutlined", danger: true, onClick: () => onDelete(keys)},
            ]}}
        />
        <ProModal ref={modalRef} title={<Title>{editRecords.length > 1 ? "批量编辑节点" : "编辑节点"}</Title>}
            onOk={form.submit} width={380} modalProps={{confirmLoading: editLoading, forceRender: true} as any}>
            <Form form={form} layout="vertical" variant="filled" className={styles.form} onFinish={onFormFinish}>
                <Form.Item name="status" label="服务状态" rules={[{required: true, message: "请选择服务状态"}]}>
                    <Select placeholder="请选择服务状态" options={Object.entries(statusData).map(([value, label]) => ({value, label}))}/>
                </Form.Item>
            </Form>
        </ProModal>
    </Body>;
};
