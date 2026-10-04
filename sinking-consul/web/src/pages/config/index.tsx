import React, {useEffect, useRef, useState} from 'react';
import {Body, PageTable, Title, ProModal, useTheme} from 'sinking-antd';
import type {PageTableProps, PageTableRef, ProModalRef} from 'sinking-antd';
import {getParams} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {App, Button, Form, Input, Select, Spin, Row, Col, Tooltip, Typography} from 'antd';
import {createConfig, deleteConfig, getConfigInfo, getConfigList, updateConfig} from "@/service/admin/config";
import defaultSettings from "../../../config/defaultSettings";
import AceEditor from "@/components/ace-editor";
import ActionDropdown from '@/pages/components/action-dropdown';
import RecordTime from '@/pages/components/record-time';
import useStyles from './styles';

const AcePath = `${defaultSettings?.basePath || '/'}ace/`;

export default (): React.ReactNode => {
    const [enumsData, enumLoading] = useEnums(["config"]);
    const {message, modal} = App.useApp();
    const {styles} = useStyles();
    const theme = useTheme();
    const typeData = (enumsData?.config?.type || {}) as unknown as Record<string, string>;
    const statusData = (enumsData?.config?.status || {}) as unknown as Record<string, string>;

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
    const requestRef = useRef(0);
    const [openRowKey, setOpenRowKey] = useState('');

    useEffect(() => () => { requestRef.current += 1; }, []);

    const mapTypeToAceMode = (t?: string) => {
        if (!t) return 'text';
        return (t || '').toLowerCase();
    }

    const openCreate = () => {
        requestRef.current += 1;
        setIsEditMode(false);
        setFormInfoLoading(false);
        form.resetFields();
        setAceContent('');
        setAceMode('text');
        formModalRef.current?.show();
    };

    const openEdit = async (record: any) => {
        const requestId = ++requestRef.current;
        setIsEditMode(true);
        setEditKeys([{group: record.group, name: record.name}]);
        form.resetFields();
        form.setFieldsValue(record);
        setAceMode(mapTypeToAceMode(record.type));
        setAceContent('');
        setFormInfoLoading(true);
        formModalRef.current?.show();
        await getConfigInfo({
            body: {group: record.group, name: record.name},
            onSuccess: (r: any) => {
                if (requestRef.current === requestId) setAceContent(r?.data?.content || '');
            },
            onFail: (r: any) => {
                if (requestRef.current !== requestId) return;
                formModalRef.current?.hide();
                message.error(r?.message || '读取配置失败');
            },
            onFinally: () => {
                if (requestRef.current === requestId) setFormInfoLoading(false);
            }
        });
    };

    const onDelete = (records: any[]) => {
        const keys = records?.map((r: any) => typeof r === 'string' ? r.split(':') : [r.group, r.name])
            .map(([group, name]: any[]) => ({group, name}));
        modal.confirm({
            title: '删除配置',
            content: `确定删除选中的 ${keys?.length || 0} 条配置吗？`,
            okText: '确定',
            okType: 'danger',
            cancelText: '取消',
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

    const tableRef = useRef<PageTableRef<any> | null>(null);
    const columns: PageTableProps<any>['columns'] = [
        {
            title: '配置分组',
            dataIndex: 'group',
            key: 'group',
            width: 150,
            render: (value: string) => <Typography.Text className="config-copy" copyable={value ? {text: value} : false}>{value || '-'}</Typography.Text>,
        },
        {
            title: '配置名称',
            dataIndex: 'name',
            key: 'name',
            width: 190,
            render: (value: string, record: any) => (
                <Typography.Text className="config-copy" copyable={value ? {text: value} : false}>
                    <button className="config-name" type="button" onClick={() => openEdit(record)}>{value || '-'}</button>
                </Typography.Text>
            ),
        },
        {
            title: '配置类型',
            dataIndex: 'type',
            key: 'type',
            width: 100,
            render: (value: string) => <span className="config-text">{typeData[value] || value || '-'}</span>,
        },
        {
            title: '哈希',
            dataIndex: 'hash',
            key: 'hash',
            width: 210,
            render: (value: string) => (
                <Typography.Text className="config-copy config-hash" copyable={value ? {text: value} : false}>
                    <Tooltip title={value || '-'}><span className="config-hash-value">{value || '-'}</span></Tooltip>
                </Typography.Text>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 90,
            render: (value: number) => <span className={`config-state ${Number(value) === 0 ? 'enabled' : 'disabled'}`}><i/>{statusData[String(value)] || '未知状态'}</span>,
        },
        {
            title: '相关时间',
            dataIndex: 'create_time',
            key: 'create_time',
            width: 220,
            sorter: true,
            render: (_: any, record: any) => <RecordTime createTime={record.create_time} updateTime={record.update_time}/>,
        },
        {
            title: '操作',
            key: 'action',
            width: 80,
            align: 'center',
            fixed: 'right',
            className: 'action-cell',
            render: (_: any, record: any) => (
                <ActionDropdown open={openRowKey === JSON.stringify([record.group, record.name])}
                    onOpenChange={(open) => setOpenRowKey(open ? JSON.stringify([record.group, record.name]) : '')}
                    menu={{
                    items: [
                        {
                            key: 'edit',
                            label: '编辑',
                            onClick: () => openEdit(record),
                        },
                        {
                            key: 'delete',
                            label: '删除',
                            danger: true,
                            onClick: () => onDelete([record]),
                        },
                    ]
                }}>
                    <Button size="small" aria-label="配置操作">操作</Button>
                </ActionDropdown>
            ),
        },
    ];

    return (
        <Body loading={enumLoading}>
            <PageTable<any>
                ref={tableRef}
                ariaLabel="配置管理"
                className={styles.table}
                hero={{eyebrow: 'CONFIGURATION CENTER', title: '配置管理', action: {label: '新增配置', onClick: openCreate}}}
                rowKey={(record) => JSON.stringify([record.group, record.name])}
                rowClassName={(record) => openRowKey === JSON.stringify([record.group, record.name]) ? 'action-menu-open' : ''}
                columns={columns}
                defaultPage={1}
                defaultPageSize={10}
                search={{
                    placeholder: '搜索配置',
                    maxLength: 200,
                    defaultField: 'keyword',
                    fields: [
                        {value: 'keyword', label: '全部', placeholder: '搜索分组、名称、哈希或内容'},
                        {value: 'group', label: '分组', placeholder: '搜索配置分组'},
                        {value: 'name', label: '名称', placeholder: '搜索配置名称'},
                        {value: 'hash', label: '哈希', placeholder: '搜索内容哈希'},
                        {value: 'content', label: '内容', placeholder: '搜索配置内容'},
                    ],
                }}
                filters={[
                    {type: 'select', name: 'type', label: '配置类型筛选', icon: 'CodeOutlined', allLabel: '全部类型', valueEnum: typeData},
                    {type: 'select', name: 'status', label: '配置状态筛选', icon: 'FlagOutlined', allLabel: '全部状态', valueEnum: statusData},
                    {type: 'date', name: 'create_time', label: '创建时间筛选'},
                ]}
                toolbar={{refresh: {ariaLabel: '刷新配置列表'}}}
                rowSelection={{
                    preserveSelectedRowKeys: true,
                    selections: true,
                    actions: (_keys, selected) => [
                        {key: 'edit', type: 'button', label: '批量编辑', onClick: () => {
                            setEditKeys(selected.map((r: any) => ({group: r.group, name: r.name})));
                            batchForm.resetFields();
                            batchModalRef.current?.show();
                        }},
                        {key: 'delete', type: 'button', label: '批量删除', danger: true, onClick: () => onDelete(selected)},
                    ]
                }}
                request={async (params, sort) => {
                    const response = await getConfigList({body: getParams({...params}, sort)});
                    if (response?.code !== 200) throw new Error(response?.message || '获取配置列表失败');
                    return {data: response.data?.list || [], total: response.data?.total ?? 0, success: true};
                }}
                paginationAffix={{offsetBottom: 15}}
                pagination={{unit: '条'}}/>

            <ProModal
                ref={formModalRef}
                title={<Title>{isEditMode ? '编辑配置' : '新增配置'}</Title>}
                onOk={form?.submit}
                width={800}
                onCancel={() => {
                    requestRef.current += 1;
                    formModalRef.current?.hide();
                }}
                modalProps={{
                    rootClassName: styles.modal,
                    confirmLoading: formBtnLoading || formInfoLoading,
                    forceRender: true,
                    footer: formInfoLoading ? null : undefined,
                    okText: '保存',
                    cancelText: '取消',
                    closable: !formBtnLoading,
                    keyboard: !formBtnLoading,
                    mask: {closable: !formBtnLoading},
                    cancelButtonProps: {disabled: formBtnLoading},
                    style: {top: 100, paddingBottom: 100},
                } as any}>
                <Form form={form} layout="vertical" variant="filled" onFinish={onFormFinish}>
                    {formInfoLoading ? (
                        <div className={styles.loading}>
                            <Spin/>
                        </div>
                    ) : (
                        <>
                            <Row gutter={16}>
                                <Col xs={24} sm={12}>
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
                                <Col xs={24} sm={12}>
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
                                <Col xs={24} sm={12}>
                                    <Form.Item
                                        name="type"
                                        label="配置类型"
                                        rules={[{required: true, message: '请选择配置类型'}]}
                                    >
                                        <Select
                                            placeholder="请选择类型"
                                            onChange={(v) => setAceMode(mapTypeToAceMode(v))}
                                            options={Object.entries(typeData).map(([value, label]) => ({value, label}))}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col xs={24} sm={12}>
                                    <Form.Item
                                        name="status"
                                        label="状态"
                                        rules={[{required: true, message: '请选择状态'}]}
                                    >
                                        <Select
                                            placeholder="请选择状态"
                                            options={Object.entries(statusData).map(([value, label]) => ({value: Number(value), label}))}
                                        />
                                    </Form.Item>
                                </Col>
                            </Row>
                            <Form.Item label="配置内容" name={isEditMode ? undefined : "content" as any}>
                                <AceEditor
                                    value={aceContent}
                                    mode={aceMode}
                                    showPrintMargin={false}
                                    theme={theme?.isDarkTheme?.() ? 'monokai' : 'chrome'}
                                    fontSize={theme?.isCompactTheme?.() ? 12 : 14}
                                    width={'100%'}
                                    height={400}
                                    acePath={AcePath}
                                    onChange={(v: string) => setAceContent(v)}
                                    className={styles.editor}
                                />
                            </Form.Item>
                        </>
                    )}
                </Form>
            </ProModal>

            <ProModal
                ref={batchModalRef}
                title={<Title>批量编辑配置</Title>}
                onOk={batchForm?.submit}
                width={350}
                modalProps={{
                    rootClassName: styles.modal,
                    confirmLoading: batchBtnLoading,
                    forceRender: true,
                    okText: '保存',
                    cancelText: '取消',
                    closable: !batchBtnLoading,
                    keyboard: !batchBtnLoading,
                    mask: {closable: !batchBtnLoading},
                    cancelButtonProps: {disabled: batchBtnLoading},
                    style: {top: 100, paddingBottom: 100},
                } as any}>
                <Form form={batchForm} layout="vertical" variant="filled" onFinish={onBatchEditFinish}>
                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{required: true, message: '请选择状态'}]}
                    >
                        <Select
                            placeholder="请选择状态"
                            options={Object.entries(statusData).map(([value, label]) => ({value: Number(value), label}))}
                        />
                    </Form.Item>
                </Form>
            </ProModal>
        </Body>
    );
}
