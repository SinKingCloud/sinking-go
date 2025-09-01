import React, { useRef } from 'react';
import { Body, ProTable, Title } from "@/components";
import { Tag, Space, Typography } from "antd";
import { PlusOutlined, MinusOutlined } from '@ant-design/icons';
import { getData } from "@/utils/page";
import { getPayLogList } from "@/service/user/pay";
import { createStyles } from 'antd-style';

const { Text } = Typography;

const useStyles = createStyles(({ css, isDarkMode }) => ({
    moneyIncome: css`
        color: #52c41a;
        font-weight: 500;
    `,
    moneyExpense: css`
        color: #ff4d4f;
        font-weight: 500;
    `,
}));

export default (): React.ReactNode => {
    const { styles } = useStyles();
    const ref = useRef({});

    const columns: any[] = [
        {
            title: 'ID',
            dataIndex: 'id',
            tip: '记录ID',
            hideInSearch: true,
            sorter: true,
            width: 80,
            align: 'center',
        },
        {
            title: '类型',
            dataIndex: 'type',
            tip: '收支类型',
            width: 100,
            valueEnum: {
                0: {
                    text: '充值',
                    status: 'success',
                },
                1: {
                    text: '扣除',
                    status: 'error',
                },
            },
            render: (_: any, record: any) => {
                const isIncome = record.type === 0;
                return (
                    <Tag 
                        icon={isIncome ? <PlusOutlined /> : <MinusOutlined />} 
                        color={isIncome ? 'success' : 'error'}
                    >
                        {isIncome ? '充值' : '扣除'}
                    </Tag>
                );
            },
        },
        {
            title: '标题',
            dataIndex: 'title',
            tip: '操作标题',
            valueType: 'text',
            ellipsis: true,
            width: 120,
        },
        {
            title: '内容',
            dataIndex: 'content',
            tip: '详细内容',
            hideInSearch: true,
            valueType: 'text',
            ellipsis: true,
            width: 200,
        },
        {
            title: '金额',
            dataIndex: 'money',
            tip: '变动金额',
            hideInSearch: true,
            width: 120,
            render: (_: any, record: any) => {
                const isIncome = record.type === 0;
                const amount = parseFloat(record.money || 0).toFixed(2);
                return (
                    <Text className={isIncome ? styles.moneyIncome : styles.moneyExpense}>
                        {isIncome ? '+' : '-'}￥{amount}
                    </Text>
                );
            },
        },
        {
            title: '创建时间',
            valueType: 'dateTime',
            dataIndex: 'create_time',
            tip: '记录创建时间',
            sorter: true,
            hideInSearch: true,
            width: 160,
        },
        {
            title: '创建时间',
            valueType: 'dateRange',
            dataIndex: 'create_time',
            tip: '按时间范围筛选',
            hideInTable: true,
            transform: (value: any) => {
                return {
                    create_time_start: value[0]?.format ? value[0].format('YYYY-MM-DD HH:mm:ss') : value[0],
                    create_time_end: value[1]?.format ? value[1].format('YYYY-MM-DD HH:mm:ss') : value[1],
                };
            },
        },
    ] as any;

    return (
        <Body>
            <ProTable
                ref={ref}
                title={<Title>余额明细</Title>}
                extraRefreshBtn={true}
                rowKey={'id'}
                columns={columns}
                defaultPage={1}
                defaultPageSize={20}
                search={{
                    layout: 'vertical',
                }}
                request={(params, sort) => {
                    return getData(params, sort, getPayLogList)
                }}
                paginationAffix={true}
                scroll={{ x: 800 }}
            />
        </Body>
    );
};
