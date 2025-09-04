import React from 'react';
import {Body, ProTable, Title} from '@/components';
import {getData} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {getClusterList} from "@/service/admin/cluster";
import {ago} from "@/utils/time";

export default (): React.ReactNode => {

    const [enumsData, loading] = useEnums(["cluster"]);

    const columns: any[] = [
        {
            title: '集群地址',
            dataIndex: 'address',
            tip: '集群访问地址',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '在线状态',
            dataIndex: 'status',
            tip: '集群在线状态',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.cluster?.status || {}).map(([key, value]) => {
                return [
                    key,
                    {text: value, color: key === '1' ? 'red' : 'green'}
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
    ] as any;

    return (
        <Body loading={loading}>
            <ProTable
                extraRefreshBtn={true}
                title={<Title>集群管理</Title>}
                pageInTable={true}
                rowKey={'address'}
                columns={columns}
                request={(params, sort) => {
                    return getData(params, sort, getClusterList)
                }}
            />
        </Body>
    );
};
