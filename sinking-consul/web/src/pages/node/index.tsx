import React, {useRef, useState} from 'react';
import {Body, ProTable, Title, ProModal} from '@/components';
import {getData} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {getNodeList, updateNode} from "@/service/admin/node";
import {ago} from "@/utils/time";
import {Button, Select, Form, App, Dropdown} from 'antd';
import {ProModalRef} from "@/components/pro-modal";

export default (): React.ReactNode => {
    const [enumsData] = useEnums(["node"]);
    const {message} = App.useApp();

    /**
     * 编辑表单
     */
    const modalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [editRecords, setEditRecords] = useState([]);
    const [editBtnLoading, setEditBtnLoading] = useState(false);
    const onFormFinish = async (values) => {
        setEditBtnLoading(true);
        await updateNode({
            body: {
                addresses: editRecords,
                status: values?.status?.toString()
            },
            onSuccess: (r: any) => {
                modalRef.current?.hide();
                tableRef.current?.refreshTableData();
                tableRef.current?.clearSelectedRows();
                message?.success(r?.message || '编辑成功');
            },
            onFail: (r: any) => {
                message?.error(r?.message || "请求失败");
            }
        });
        setEditBtnLoading(false);
    };

    /**
     * 表格
     */
    const tableRef = React.useRef<any>(null);
    const columns: any[] = [
        {
            title: '服务分组',
            dataIndex: 'group',
            tip: '服务分组',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '服务名称',
            dataIndex: 'name',
            tip: '服务名称',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '节点地址',
            dataIndex: 'address',
            tip: '节点访问地址',
            valueType: 'text',
            hideInSearch: true,
            copyable: true,
        },
        {
            title: '在线状态',
            dataIndex: 'online_status',
            tip: '节点在线状态',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.node?.online_status || {}).map(([key, value]) => {
                return [
                    key,
                    {text: value, color: key == '0' ? 'green' : 'red'}
                ]
            })),
        },
        {
            title: '最后心跳',
            dataIndex: 'last_heart',
            tip: '最后心跳时间戳',
            hideInSearch: true,
            render: (text: any) => {
                if (!text) return '-';
                return ago(new Date(text * 1000).toLocaleString('zh-CN'));
            },
        },
        {
            title: '服务状态',
            dataIndex: 'status',
            tip: '服务状态',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.node?.status || {}).map(([key, value]) => {
                return [
                    key,
                    {text: value, color: key == '0' ? 'green' : 'yellow'}
                ]
            })),
        },
        {
            title: '创建时间',
            valueType: 'dateRange',
            dataIndex: 'create_time',
            tip: '创建时间',
            sorter: true,
            transform: (value: any) => {
                return {
                    create_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                    create_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
                };
            },
        },
        {
            title: '更新时间',
            valueType: 'dateRange',
            dataIndex: 'update_time',
            tip: '更新时间',
            sorter: true,
            hideInSearch: true,
            transform: (value: any) => {
                return {
                    update_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                    update_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
                };
            },
        },
        {
            title: '操作',
            valueType: 'option',
            key: 'option',
            render: (text: any, record: any) => [
                <Dropdown key={record?.address} menu={{
                    items: [
                        {
                            key: "edit",
                            label: <a onClick={() => {
                                form?.setFieldsValue(record);
                                setEditRecords([record?.address]);
                                modalRef.current?.show();
                            }}>编 辑</a>
                        },
                    ]
                }} trigger={['click']} placement="bottom" arrow={true}>
                    <Button size="small">操作</Button>
                </Dropdown>
            ],
        },
    ] as any;


    return (
        <Body>
            <ProTable
                ref={tableRef}
                extraRefreshBtn={true}
                title={<Title>服务节点</Title>}
                rowKey={'address'}
                columns={columns}
                defaultPage={1}
                defaultPageSize={20}
                rowSelection={{
                    rightExtra: (
                        <Button type={"primary"} ghost onClick={() => {
                            form?.resetFields();
                            setEditRecords(tableRef?.current?.getSelectedRowKeys());
                            modalRef.current?.show();
                        }}>批量编辑</Button>
                    )
                }}
                request={(params, sort) => {
                    return getData(params, sort, getNodeList)
                }}
                paginationAffix={true}
                selectionAffix={true}
            />

            <ProModal
                ref={modalRef}
                title="编辑节点"
                onOk={form?.submit}
                width={350}
                modalProps={{
                    confirmLoading: editBtnLoading,
                }}
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={onFormFinish}
                >
                    <Form.Item
                        name="status"
                        label="服务状态"
                        rules={[{required: true, message: '请选择服务状态'}]}
                    >
                        <Select
                            placeholder="请选择状态"
                            options={Object.entries(enumsData?.node?.status || {}).map(([key, value]) => ({
                                label: value,
                                value: parseInt(key)
                            })) as any}
                        />
                    </Form.Item>
                </Form>
            </ProModal>
        </Body>
    );
};
