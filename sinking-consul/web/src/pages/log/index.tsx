import React from 'react';
import {Body, ProTable, Title} from '@/components';
import {getData} from "@/utils/page";
import {useEnums} from "@/utils/enum";
import {getLog} from "@/service/admin/person";

export default (): React.ReactNode => {

    const [enumsData, loading] = useEnums(["log"]);

    const columns: any[] = [
        {
            title: 'ID',
            dataIndex: 'id',
            tip: '日志ID',
            hideInSearch: true,
            sorter: true,
        },
        {
            title: '操作IP',
            dataIndex: 'ip',
            tip: '请求来源IP地址',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '操作类型',
            dataIndex: 'type',
            tip: '操作事件类型',
            valueEnum: Object.fromEntries(Object.entries(enumsData?.log?.type || {}).map(([key, value]) => {
                return [
                    key,
                    {text: value, color: 'blue'}
                ]
            })),
        },
        {
            title: '操作标题',
            dataIndex: 'title',
            tip: '操作事件标题',
            hideInSearch: true,
        },
        {
            title: '操作内容',
            dataIndex: 'content',
            tip: '详细操作内容',
            hideInSearch: true,
            ellipsis: true,
        },
        {
            title: '操作时间',
            valueType: 'dateRange',
            dataIndex: 'create_time',
            tip: '操作时间',
            sorter: true,
            transform: (value: any) => {
                return {
                    create_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                    create_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
                };
            },
        },
    ] as any;

    return (
        <Body loading={loading}>
            <ProTable
                extraRefreshBtn={true}
                title={<Title>操作日志</Title>}
                rowKey={'id'}
                columns={columns}
                request={(params, sort) => {
                    return getData(params, sort, getLog)
                }}
            />
        </Body>
    );
}
;
