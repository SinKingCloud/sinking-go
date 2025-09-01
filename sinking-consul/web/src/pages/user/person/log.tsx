import React from 'react';
import {Body, ProTable, Title} from '@/components';
import {getData} from "@/utils/page";
import {getLogList} from "@/service/user/log";
import {useEnums} from "@/utils/enum";

export default (): React.ReactNode => {

    const [enumsData, loading] = useEnums(["user_log", "user_order", "mt_notify"]);
    const typeEnum = Object.fromEntries(Object.entries(enumsData?.user_log?.types || {}).map(([key, value]) => [
        key,
        {text: value, color: 'blue'}
    ]));

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
            dataIndex: 'request_ip',
            tip: '请求来源IP地址',
            valueType: 'text',
            copyable: true,
        },
        {
            title: '操作类型',
            dataIndex: 'type',
            tip: '操作事件类型',
            valueEnum: typeEnum,
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
                defaultPage={1}
                defaultPageSize={20}
                request={(params, sort) => {
                    return getData(params, sort, getLogList)
                }}
            />
        </Body>
    );
}
;
