import React, {useRef, useState} from 'react';
import {Body, ProTable, Title, ProModal, ProModalRef} from 'sinking-antd';
import {getData} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {App, Button, Dropdown, Form, Input, Select, Spin, Row, Col} from 'antd';
import {createConfig, deleteConfig, getConfigInfo, getConfigList, updateConfig} from "@/service/admin/config";
import defaultSettings from "../../../config/defaultSettings";
import AceEditor from "@/components/ace-editor";
import {createStyles} from "antd-style";

const AcePath = defaultSettings?.basePath + "ace/" || "/ace/";

const useStyles = createStyles(({token}): any => {
    return {
        ace: {
            ".ace_editor": {
                borderRadius: token.borderRadius + "px !important",
            }
        },
    };
});

export default (): React.ReactNode => {
    const [enumsData] = useEnums(["config"]);
    const {message, modal} = App.useApp();
    const {styles: {ace}} = useStyles();

    const formModalRef = useRef<ProModalRef>({} as ProModalRef);
    const batchModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [form] = Form.useForm();
    const [batchForm] = Form.useForm();
    const [isEditMode, setIsEditMode] = useState(false);
    const [editKeys, setEditKeys] = useState<any[]>([]);
    const [formBtnLoading, setFormBtnLoading] = useState(false);
    const [batchBtnLoading, setBatchBtnLoading] = useState(false);
    const [formInfoLoading, setFormInfoLoading] = useState(false);
    const [aceContent, setAceContent] = useState<string>('');
    const [aceMode, setAceMode] = useState<string>('text');

    const mapTypeToAceMode = (t?: string) => {
        if (!t) return 'text';
        return (t || '').toLowerCase();
    }

    const onDelete = (records: any[]) => {
        const keys = records?.map((r: any) => typeof r === 'string' ? r.split(':') : [r.group, r.name])
            .map(([group, name]: any[]) => ({group, name}));
        modal.confirm({
            title: '删除配置',
            content: `确定删除选中的 ${keys?.length || 0} 条配置吗？`,
            okText: '确 定',
            okType: 'danger',
            cancelText: '取 消',
            maskClosable: true,
            onOk: async () => {
                await deleteConfig({
                    body: {keys},
                    onSuccess: (r: any) => {
                        tableRef.current?.refreshTableData();
                        tableRef.current?.clearSelectedRows();
                        message?.success(r?.message || '删除成功');
                    },
                    onFail: (r: any) => message?.error(r?.message || '请求失败')
                })
            }
        } as any)
    }

    const onFormFinish = async (values: any) => {
        setFormBtnLoading(true);
        const body: any = isEditMode ? {keys: editKeys, ...values} : {...values};
        if (values?.status !== undefined && values?.status !== null) {
            body.status = String(values.status);
        }
        if (aceContent.trim() !== "") {
            body.content = aceContent;
        } else {
            delete body.content;
        }
        const apiCall = isEditMode ? updateConfig : createConfig;
        const successMsg = isEditMode ? '编辑成功' : '创建成功';
        await apiCall({
            body,
            onSuccess: (r: any) => {
                formModalRef.current?.hide();
                tableRef.current?.refreshTableData();
                if (isEditMode) {
                    tableRef.current?.clearSelectedRows();
                }
                message?.success(r?.message || successMsg);
            },
            onFail: (r: any) => message?.error(r?.message || '请求失败'),
            onFinally: () => setFormBtnLoading(false)
        });
    };

    const onBatchEditFinish = async (values: any) => {
        setBatchBtnLoading(true);
        await updateConfig({
            body: {
                keys: editKeys,
                status: values?.status !== undefined && values?.status !== null ? String(values.status) : undefined,
            },
            onSuccess: (r: any) => {
                batchModalRef.current?.hide();
                tableRef.current?.refreshTableData();
                tableRef.current?.clearSelectedRows();
                message?.success(r?.message || '编辑成功');
            },
            onFail: (r: any) => message?.error(r?.message || '请求失败'),
            onFinally: () => setBatchBtnLoading(false)
        });
    };

    const tableRef = React.useRef<any>(null);
    const columns: any[] = [
        {
            title: '配置分组',
            dataIndex: 'group',
            tip: '配置分组',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '配置名称',
            dataIndex: 'name',
            tip: '配置名称',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '配置类型',
            dataIndex: 'type',
            tip: '配置类型',
            valueType: 'select',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.config?.type || {}).map(([key, value]) => [key, {text: value}]))
        },
        {
            title: '哈希',
            dataIndex: 'hash',
            tip: '内容哈希',
            valueType: 'text',
            hideInSearch: true,
            copyable: true,
        },
        {
            title: '状态',
            dataIndex: 'status',
            tip: '状态',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.config?.status || {}).map(([key, value]) => [key, {
                text: value,
                color: key === '0' ? 'green' : 'red'
            }]))
        },
        {
            title: '创建时间',
            valueType: 'dateRange',
            dataIndex: 'create_time',
            tip: '创建时间',
            sorter: true,
            transform: (value: any) => ({
                create_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                create_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
            }),
        },
        {
            title: '更新时间',
            valueType: 'dateRange',
            dataIndex: 'update_time',
            tip: '更新时间',
            sorter: true,
            hideInSearch: true,
            transform: (value: any) => ({
                update_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                update_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
            }),
        },
        {
            title: '操作',
            valueType: 'option',
            hideInSearch: true,
            render: (text: any, record: any) => [
                <Dropdown key={`${record?.group}:${record?.name}`} menu={{
                    items: [
                        {
                            key: 'edit',
                            label: <a onClick={async () => {
                                setIsEditMode(true);
                                const keys = [{group: record?.group, name: record?.name}];
                                setEditKeys(keys);
                                form?.setFieldsValue(record);
                                setAceMode(mapTypeToAceMode(record?.type));
                                setFormInfoLoading(true);
                                formModalRef.current?.show();
                                await getConfigInfo({
                                    body: {group: record?.group, name: record?.name},
                                    onSuccess: (r: any) => setAceContent(r?.data?.content || ''),
                                    onFail: (r: any) => message?.error(r?.message || '读取配置失败'),
                                    onFinally: () => setFormInfoLoading(false)
                                });
                            }}>编 辑</a>
                        },
                        {
                            key: 'delete',
                            label: <a onClick={() => onDelete([{group: record?.group, name: record?.name}])}>删 除</a>
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
                title={<Title>配置管理</Title>}
                rowKey={'hash'}
                columns={columns}
                defaultPage={1}
                defaultPageSize={20}
                rowSelection={{
                    rightExtra: [
                        <Button key="delete" type="primary" danger ghost onClick={() => {
                            const selected = tableRef?.current?.getSelectedRows() || [];
                            if (!selected.length) return message?.warning('请选择记录');
                            onDelete(selected.map((r: any) => ({group: r.group, name: r.name})));
                        }}>批量删除</Button>,
                        <Button key="edit" type="primary" ghost onClick={() => {
                            const selected = tableRef?.current?.getSelectedRows() || [];
                            if (!selected.length) return message?.warning('请选择记录');
                            const keys = selected.map((r: any) => ({group: r.group, name: r.name}));
                            setEditKeys(keys);
                            batchForm?.resetFields();
                            batchModalRef.current?.show();
                        }}>批量编辑</Button>
                    ]
                }}
                request={(params, sort) => getData(params, sort, getConfigList)}
                paginationAffix={true}
                selectionAffix={true}
                extra={<Button key="create" type="primary" onClick={() => {
                    setIsEditMode(false);
                    form?.resetFields();
                    setAceContent('');
                    setAceMode('text');
                    formModalRef.current?.show();
                }}>新 增</Button>}/>

            <ProModal
                ref={formModalRef}
                title={<Title>{isEditMode ? '编辑配置' : '新增配置'}</Title>}
                onOk={form?.submit}
                width={800}
                modalProps={{
                    confirmLoading: formBtnLoading || formInfoLoading,
                    forceRender: true,
                    footer: formInfoLoading ? null : undefined,
                    style: {top: 30}
                } as any}>
                <Form form={form} layout="vertical" onFinish={onFormFinish}>
                    {formInfoLoading ? (
                        <div style={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            height: 260,
                            width: '100%'
                        }}>
                            <Spin size="large"/>
                        </div>
                    ) : (
                        <>
                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        label="配置分组"
                                        name="group"
                                        rules={isEditMode ? [] : [{required: true, message: '请输入配置分组'}]}
                                    >
                                        <Input
                                            disabled={isEditMode}
                                            placeholder={isEditMode ? '' : '请输入配置分组'}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        label="配置名称"
                                        name="name"
                                        rules={isEditMode ? [] : [{required: true, message: '请输入配置名称'}]}
                                    >
                                        <Input
                                            disabled={isEditMode}
                                            placeholder={isEditMode ? '' : '请输入配置名称'}
                                        />
                                    </Form.Item>
                                </Col>
                            </Row>
                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item
                                        name="type"
                                        label="配置类型"
                                        rules={[{required: true, message: '请选择配置类型'}]}
                                    >
                                        <Select
                                            placeholder="请选择类型"
                                            onChange={(v) => setAceMode(mapTypeToAceMode(v))}
                                        >
                                            {Object.entries(enumsData?.config?.type || {}).map(([key, value]) => (
                                                <Select.Option key={key} value={key}>{value as any}</Select.Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        name="status"
                                        label="状态"
                                        rules={[{required: true, message: '请选择状态'}]}
                                    >
                                        <Select
                                            placeholder="请选择状态"
                                        >
                                            {Object.entries(enumsData?.config?.status || {}).map(([key, value]) => (
                                                <Select.Option key={key}
                                                               value={parseInt(key)}>{value as any}</Select.Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                            </Row>
                            <Form.Item label="配置内容" name={isEditMode ? undefined : "content" as any}>
                                <AceEditor
                                    value={aceContent}
                                    mode={aceMode}
                                    showPrintMargin={false}
                                    theme={'monokai'}
                                    width={'100%'}
                                    height={400}
                                    acePath={AcePath}
                                    onChange={(v: string) => setAceContent(v)}
                                    className={ace}
                                />
                            </Form.Item>
                        </>
                    )}
                </Form>
            </ProModal>

            <ProModal
                ref={batchModalRef}
                title={<Title>编辑配置</Title>}
                onOk={batchForm?.submit}
                width={350}
                modalProps={{confirmLoading: batchBtnLoading, forceRender: true} as any}>
                <Form form={batchForm} layout="vertical" onFinish={onBatchEditFinish}>
                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{required: true, message: '请选择状态'}]}
                    >
                        <Select placeholder="请选择状态">
                            {Object.entries(enumsData?.config?.status || {}).map(([key, value]) => (
                                <Select.Option key={key} value={parseInt(key)}>{value as any}</Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                </Form>
            </ProModal>
        </Body>
    );
}
