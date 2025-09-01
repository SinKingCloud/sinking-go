import React, {useRef, useState} from 'react';
import {Body, ProModal, ProTable, Title} from "@/components";
import {Button, DatePicker} from "antd";
import {getData} from "@/utils/page";
import {getLogList} from "@/service/user/log";
import {ProModalRef} from "@/components/pro-modal";

export default (): React.ReactNode => {
    const columns: any[] = [
        {
            title: 'ID',
            dataIndex: 'id',
            tip: '记录ID',
            hideInSearch: true,
            sorter: true,
            //valueType: 'digit',
        },
        {
            title: '操作时间',
            valueType: 'dateRange',
            dataIndex: 'create_time',
            tip: '请求接口时间',
            hideInTable: true,
            transform: (value: any) => {
                return {
                    create_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                    create_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
                };
            },
        },
        {
            title: '请求IP',
            dataIndex: 'request_ip',
            tip: '请求者IP',
            copyable: true,
            valueType: 'text',
        },
        {
            title: '进度',
            dataIndex: 'progress',
            tip: '任务进度',
            valueType: 'progress',
            //hideInSearch: true,
        },
        {
            title: '评分',
            dataIndex: 'rate',
            tip: '服务评分',
            valueType: 'rate',
        },
        {
            title: '事件类型',
            dataIndex: 'type',
            tip: '事件类型',
            valueEnum: {
                0: {
                    text: '登陆',
                    color: 'success',
                },
                1: {
                    text: '查看',
                    color: 'warning',
                },
                2: {
                    text: '删除',
                    color: 'error',
                },
                3: {
                    text: '修改',
                    color: 'default',
                },
                4: {
                    text: '创建',
                    color: 'success',
                },
            }
        },
        {
            title: '标题',
            dataIndex: 'title',
            tip: '标题',
            valueType: 'text',
        },
        {
            title: '内容',
            dataIndex: 'content',
            tip: '内容',
            hideInSearch: true,
            valueType: 'text',
        },
        {
            title: '是否启用',
            dataIndex: 'enabled',
            tip: '启用状态',
            valueType: 'switch',
            hideInSearch: false,
            props: {
                defaultChecked: true,
            },
        },
        {
            title: '完成率',
            dataIndex: 'percent',
            tip: '完成百分比',
            valueType: 'percent',
            hideInSearch: true,
        },
        {
            title: '操作时间',
            valueType: 'dateTime',
            dataIndex: 'create_time',
            tip: '操作时间',
            sorter: true,
            hideInSearch: true,
            editable: false,
        },

    ] as any;

    const ref = useRef({});
    const [selectedRows, setSelectedRows] = useState<any[]>([]);

    const modalRef = useRef<ProModalRef>({})

    return (
        <Body>
            <ProModal ref={modalRef} maskClosable={false} title={<Title>日志列表</Title>}>
                <ProTable
                    //  title={<Title>日志列表</Title>}
                    rowKey={'id'}

                    columns={columns}
                    defaultPage={1}
                    defaultPageSize={20}
                    extraRefreshBtn={false}
                    request={(params, sort) => {
                        return getData(params, sort, getLogList)
                    }}

                    // search={{
                    //     layout: 'vertical',
                    // }}
                />
            </ProModal>
            <ProTable
                ref={ref}
                title={<Title>日志列表</Title>}
                extraRefreshBtn={true}
                rowKey={'id'}
                columns={columns}
                defaultPage={1}
                defaultPageSize={20}
                extra={<Button type={"primary"} onClick={() => {
                    modalRef.current?.show?.();
                }}>新建</Button>}
                rowSelection={{
                    fixed: true,
                    onChange: (_, selectedRows) => {
                        setSelectedRows(selectedRows)
                    },
                    //leftExtra: selectedRows?.length,
                    // hideExtra: false,
                    // type: "radio",
                    rightExtra: <>
                        <Button type={"primary"} ghost>批量编辑</Button>
                        <Button type={"primary"} danger ghost onClick={() => {
                            console.log(selectedRows)
                        }}>批量删除</Button>
                    </>
                }}
                // search={{
                //     layout: 'vertical',
                // }}
                //tableProps={{size: "large"} as any}
                //pageHidden={false}
                //pageInTable={true}
                //paginationProps={{pageSizeOptions: [10, 20, 50, 100, 200]}}
                request={(params, sort) => {
                    return getData(params, sort, getLogList)
                }}
                paginationAffix={true}
                selectionAffix={true}
            />
        </Body>
    );
};