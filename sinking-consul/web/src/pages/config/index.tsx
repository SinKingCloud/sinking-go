import React, {useRef, useState} from 'react';
import {Body, ProTable, Title, ProModal, AceEditor} from '@/components';
import {getData} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {App, Button, Dropdown, Form, Input, Select, Space, Spin, Row, Col} from 'antd';
import {ProModalRef} from "@/components/pro-modal";
import {createConfig, deleteConfig, getConfigInfo, getConfigList, updateConfig} from "@/service/admin/config";

const {Option} = Select;

export default (): React.ReactNode => {
    const [enumsData] = useEnums(["config"]);
    const {message, modal} = App.useApp();

    // 创建与编辑表单
    const createModalRef = useRef<ProModalRef>({} as ProModalRef);
    const editModalRef = useRef<ProModalRef>({} as ProModalRef);
    const batchModalRef = useRef<ProModalRef>({} as ProModalRef);
    const [createForm] = Form.useForm();
    const [editForm] = Form.useForm();
    const [batchForm] = Form.useForm();

    // 选中待编辑 keys（批量或单个）。格式：[{group, name}]
    const [editKeys, setEditKeys] = useState<any[]>([]);
    const [editBtnLoading, setEditBtnLoading] = useState(false);
    const [batchBtnLoading, setBatchBtnLoading] = useState(false);
    const [createBtnLoading, setCreateBtnLoading] = useState(false);
    const [editInfoLoading, setEditInfoLoading] = useState(false);

    // Ace 内容（单条编辑与创建使用）
    const [editAceContent, setEditAceContent] = useState<string>('');
    const [createAceContent, setCreateAceContent] = useState<string>('');
    const [editAceMode, setEditAceMode] = useState<string>('json');
    const [createAceMode, setCreateAceMode] = useState<string>('text');

    const mapTypeToAceMode = (t?: string) => {
        if (!t) return 'text';
        const lower = (t || '').toLowerCase();
        if (lower === 'json') return 'json';
        if (lower === 'yaml' || lower === 'yml') return 'yaml';
        if (lower === 'ini') return 'ini';
        return 'text';
    }

    // 删除配置（支持批量）
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

    // 提交编辑（单条）
    const onEditFinish = async (values: any) => {
        setEditBtnLoading(true);
        const body: any = {
            keys: editKeys,
        };
        if (values?.type) body.type = values.type;
        if (values?.status !== undefined && values?.status !== null) body.status = String(values.status);
        // 单条编辑允许修改内容
        body.content = editAceContent;
        await updateConfig({
            body,
            onSuccess: (r: any) => {
                editModalRef.current?.hide();
                tableRef.current?.refreshTableData();
                tableRef.current?.clearSelectedRows();
                message?.success(r?.message || '编辑成功');
            },
            onFail: (r: any) => message?.error(r?.message || '请求失败'),
            onFinally: () => setEditBtnLoading(false)
        });
    };

    // 提交批量编辑（仅状态）
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

    // 提交创建
    const onCreateFinish = async (values: any) => {
        setCreateBtnLoading(true);
        await createConfig({
            body: {
                group: values.group,
                name: values.name,
                type: values.type,
                status: values.status, // 后端为 int
                content: createAceContent,
            },
            onSuccess: (r: any) => {
                createModalRef.current?.hide();
                tableRef.current?.refreshTableData();
                message?.success(r?.message || '创建成功');
            },
            onFail: (r: any) => message?.error(r?.message || '请求失败'),
            onFinally: () => setCreateBtnLoading(false)
        });
    }

    // 表格
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
                                // 设置编辑 keys
                                const keys = [{group: record?.group, name: record?.name}];
                                setEditKeys(keys);
                                editForm?.setFieldsValue({
                                    group: record?.group,
                                    name: record?.name,
                                    type: record?.type,
                                    status: record?.status,
                                });
                                setEditAceMode(mapTypeToAceMode(record?.type));
                                setEditInfoLoading(true);
                                // 先展示弹窗 + loading
                                editModalRef.current?.show();
                                // 获取内容
                                await getConfigInfo({
                                    body: {group: record?.group, name: record?.name},
                                    onSuccess: (r: any) => setEditAceContent(r?.data?.content || ''),
                                    onFail: (r: any) => message?.error(r?.message || '读取配置失败'),
                                    onFinally: () => setEditInfoLoading(false)
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

    const rowKey = 'hash';

    return (
        <Body>
            <ProTable
                ref={tableRef}
                extraRefreshBtn={true}
                title={<Title>配置管理</Title>}
                rowKey={rowKey}
                columns={columns}
                defaultPage={1}
                defaultPageSize={20}
                rowSelection={{
                    rightExtra: (
                        <>
                            <Button type={"primary"} danger ghost onClick={() => {
                                const selected = tableRef?.current?.getSelectedRows() || [];
                                if (!selected.length) return message?.warning('请选择记录');
                                onDelete(selected.map((r: any) => ({group: r.group, name: r.name})));
                            }}>批量删除</Button>
                            <Button type={"primary"} ghost onClick={() => {
                                const selected = tableRef?.current?.getSelectedRows() || [];
                                if (!selected.length) return message?.warning('请选择记录');
                                const keys = selected.map((r: any) => ({group: r.group, name: r.name}));
                                setEditKeys(keys);
                                // 打开批量编辑，仅状态
                                batchForm?.resetFields();
                                batchModalRef.current?.show();
                            }}>批量编辑</Button>
                        </>
                    )
                }}
                request={(params, sort) => getData(params, sort, getConfigList)}
                paginationAffix={true}
                selectionAffix={true}
                extra={<Button key="create" type="primary" onClick={() => {
                    createForm?.resetFields();
                    setCreateAceContent('');
                    setCreateAceMode('text');
                    createModalRef.current?.show();
                }}>新 增</Button>}
            />

            {/* 单条编辑弹窗 */}
            <ProModal
                ref={editModalRef}
                title={<Title>编辑配置</Title>}
                onOk={editForm?.submit}
                width={800}
                modalProps={{
                    confirmLoading: editBtnLoading || editInfoLoading,
                    forceRender: true,
                    footer: editInfoLoading ? null : undefined,
                    style: {top: 30}
                } as any}
            >
                <Form form={editForm} layout="vertical" onFinish={onEditFinish}>
                    {editInfoLoading ? (
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
                                    <Form.Item label="配置分组" name="group">
                                        <Input disabled/>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item label="配置名称" name="name">
                                        <Input disabled/>
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Row gutter={16}>
                                <Col span={12}>
                                    <Form.Item name="type" label="配置类型">
                                        <Select placeholder="请选择类型" allowClear
                                                onChange={(v) => setEditAceMode(mapTypeToAceMode(v))}>
                                            {Object.entries(enumsData?.config?.type || {}).map(([key, value]) => (
                                                <Option key={key} value={key}>{value as any}</Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item name="status" label="状态">
                                        <Select placeholder="请选择状态" allowClear>
                                            {Object.entries(enumsData?.config?.status || {}).map(([key, value]) => (
                                                <Option key={key} value={parseInt(key)}>{value as any}</Option>
                                            ))}
                                        </Select>
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item label="配置内容">
                                <AceEditor
                                    value={editAceContent}
                                    mode={editAceMode}
                                    theme={'monokai'}
                                    width={'100%'}
                                    height={400}
                                    acePath={'/ace/'}
                                    onChange={(v: string) => setEditAceContent(v)}
                                />
                            </Form.Item>
                        </>
                    )}
                </Form>
            </ProModal>

            {/* 批量编辑弹窗（仅状态） */}
            <ProModal
                ref={batchModalRef}
                title={<Title>编辑配置</Title>}
                onOk={batchForm?.submit}
                width={350}
                modalProps={{confirmLoading: batchBtnLoading, forceRender: true} as any}
            >
                <Form form={batchForm} layout="vertical" onFinish={onBatchEditFinish}>
                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{required: true, message: '请选择状态'}]}
                    >
                        <Select placeholder="请选择状态">
                            {Object.entries(enumsData?.config?.status || {}).map(([key, value]) => (
                                <Option key={key} value={parseInt(key)}>{value as any}</Option>
                            ))}
                        </Select>
                    </Form.Item>
                </Form>
            </ProModal>

            {/* 创建弹窗 */}
            <ProModal
                ref={createModalRef}
                title={<Title>新增配置</Title>}
                onOk={createForm?.submit}
                width={800}
                modalProps={{confirmLoading: createBtnLoading, forceRender: true, style: {top: 30}} as any}
            >
                <Form form={createForm} layout="vertical" onFinish={onCreateFinish}>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="group" label="配置分组"
                                       rules={[{required: true, message: '请输入配置分组'}]}>
                                <Input placeholder="请输入配置分组"/>
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="name" label="配置名称"
                                       rules={[{required: true, message: '请输入配置名称'}]}>
                                <Input placeholder="请输入配置名称"/>
                            </Form.Item>
                        </Col>
                    </Row>
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item name="type" label="配置类型"
                                       rules={[{required: true, message: '请选择配置类型'}]}>
                                <Select placeholder="请选择类型"
                                        onChange={(v) => setCreateAceMode(mapTypeToAceMode(v))}>
                                    {Object.entries(enumsData?.config?.type || {}).map(([key, value]) => (
                                        <Option key={key} value={key}>{value as any}</Option>
                                    ))}
                                </Select>
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item name="status" label="状态" rules={[{required: true, message: '请选择状态'}]}>
                                <Select placeholder="请选择状态">
                                    {Object.entries(enumsData?.config?.status || {}).map(([key, value]) => (
                                        <Option key={key} value={parseInt(key)}>{value as any}</Option>
                                    ))}
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>
                    <Form.Item label="配置内容" name="content" required>
                        <AceEditor
                            value={createAceContent}
                            mode={createAceMode}
                            theme={'monokai'}
                            width={'100%'}
                            height={400}
                            acePath={'/ace/'}
                            onChange={(v: string) => setCreateAceContent(v)}
                        />
                    </Form.Item>
                </Form>
            </ProModal>
        </Body>
    );
}
