import React, { useState } from 'react';
import { Body, Title, Message } from '@/components';
import { Card, Row, Col, InputNumber, Button, Radio, Space, Divider, Typography, Alert } from 'antd';
import { WechatOutlined, AlipayOutlined, WalletOutlined } from '@ant-design/icons';
import { createStyles } from 'antd-style';
import { useModel } from 'umi';
import { createRechargeOrder, submitPayment } from '@/service/user/pay';
import { setPayJumpUrl } from '@/utils/pay';

const { Text } = Typography;

const useStyles = createStyles(({ css, isDarkMode, token }) => ({
    container: css`
        max-width: 1000px;
        margin: 0 auto;
    `,
    amountCard: css`
        .ant-card-body {
            padding: 24px;
        }
    `,
    amountButton: css`
        width: 100%;
        height: 60px;
        font-size: 18px;
        font-weight: 500;
        margin-bottom: 12px;
        border-radius: 8px;
        transition: all 0.3s ease;
        
        &:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
        
        &.selected {
            border-color: ${token.colorPrimary};
            background-color: ${isDarkMode ? 'rgba(22, 119, 255, 0.1)' : 'rgba(22, 119, 255, 0.05)'};
        }
    `,
    customInput: css`
        width: 100%;
        height: 60px;
        font-size: 16px;
        
        .ant-input-number-input {
            height: 58px;
            font-size: 18px;
            font-weight: 500;
        }
    `,
    payMethodCard: css`
        .ant-radio-group {
            width: 100%;
        }
        
        .ant-radio-button-wrapper {
            height: 80px;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 16px;
            border-radius: 8px;
            margin-bottom: 12px;
            
            &:not(:last-child) {
                margin-right: 0;
            }
            
            &:hover {
                border-color: ${token.colorPrimary};
            }
        }
        
        .ant-radio-button-wrapper-checked {
            background-color: ${isDarkMode ? 'rgba(22, 119, 255, 0.15)' : 'rgba(22, 119, 255, 0.08)'};
            border-color: ${token.colorPrimary};
        }
    `,
    payIcon: css`
        font-size: 28px;
        margin-right: 12px;
    `,
    submitButton: css`
        width: 100%;
        height: 56px;
        font-size: 18px;
        font-weight: 500;
        border-radius: 8px;
        margin-top: 24px;
    `,
    balanceInfo: css`
        padding: 16px;
        border-radius: 8px;
        background-color: ${isDarkMode ? 'rgba(255, 255, 255, 0.04)' : 'rgba(0, 0, 0, 0.02)'};
        margin-bottom: 24px;
        display: flex;
        align-items: center;
        justify-content: space-between;
    `,
    tips: css`
        margin-top: 24px;
        
        .ant-alert {
            border-radius: 8px;
        }
    `
}));

export default () => {
    const { styles } = useStyles();
    const [amount, setAmount] = useState<number>(100);
    const [customAmount, setCustomAmount] = useState<number | null>(null);
    const [payType, setPayType] = useState<number>(2); // 默认选择微信支付
    const [loading, setLoading] = useState(false);
    const user = useModel('user');

    const presetAmounts = [50, 100, 200, 500, 1000, 2000];

    const handleAmountSelect = (value: number) => {
        setAmount(value);
        setCustomAmount(null);
    };

    const handleCustomAmountChange = (value: number | null) => {
        setCustomAmount(value);
        if (value && value > 0) {
            setAmount(value);
        }
    };

    const handleSubmit = async () => {
        if (!amount || amount <= 0) {
            Message.error('请选择充值金额');
            return;
        }

        setLoading(true);
        try {
            // 创建充值订单
            const orderRes = await createRechargeOrder({
                body: {
                    money: amount,
                    type: payType
                }
            });

            if (orderRes.code === 0) {
                const tradeNo = orderRes.data;
                
                // 保存当前页面URL，支付完成后返回
                setPayJumpUrl(window.location.pathname);
                
                // 提交支付
                const payRes = await submitPayment({
                    body: {
                        trade_no: tradeNo
                    }
                });

                if (payRes.code === 0) {
                    // 跳转到支付页面
                    window.location.href = payRes.data;
                } else {
                    Message.error((payRes as any).msg || '创建支付链接失败');
                }
            } else {
                Message.error((orderRes as any).msg || '创建订单失败');
            }
        } catch (error) {
            Message.error('充值失败，请稍后重试');
        } finally {
            setLoading(false);
        }
    };

    return (
        <Body>
            <div className={styles.container}>
                <Row gutter={[16, 16]}>
                    <Col span={24}>
                        <div className={styles.balanceInfo}>
                            <Space size={12}>
                                <WalletOutlined style={{ fontSize: 24 }} />
                                <div>
                                    <Text type="secondary">当前余额</Text>
                                    <div>
                                        <Text strong style={{ fontSize: 24 }}>
                                            ￥{parseFloat(user?.web?.money || 0).toFixed(2)}
                                        </Text>
                                    </div>
                                </div>
                            </Space>
                        </div>
                    </Col>

                    <Col xl={14} lg={24} md={24} sm={24} xs={24}>
                        <Card 
                            title={<Title>选择充值金额</Title>} 
                            className={styles.amountCard}
                            variant="borderless"
                        >
                            <Row gutter={[12, 0]}>
                                {presetAmounts.map((value) => (
                                    <Col span={8} key={value}>
                                        <Button
                                            className={`${styles.amountButton} ${amount === value && !customAmount ? 'selected' : ''}`}
                                            onClick={() => handleAmountSelect(value)}
                                        >
                                            ￥{value}
                                        </Button>
                                    </Col>
                                ))}
                            </Row>

                            <Divider>自定义金额</Divider>

                            <InputNumber
                                className={styles.customInput}
                                min={1}
                                max={99999}
                                placeholder="请输入充值金额"
                                prefix="￥"
                                value={customAmount}
                                onChange={handleCustomAmountChange}
                                formatter={(value) => {
                                    if (!value) return '';
                                    return `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',');
                                }}
                                parser={(value) => {
                                    const parsedValue = value!.replace(/\$\s?|(,*)/g, '');
                                    return parsedValue === '' ? 0 : Number(parsedValue);
                                }}
                            />
                        </Card>
                    </Col>

                    <Col xl={10} lg={24} md={24} sm={24} xs={24}>
                        <Card 
                            title={<Title>选择支付方式</Title>} 
                            className={styles.payMethodCard}
                            variant="borderless"
                        >
                            <Radio.Group value={payType} onChange={(e) => setPayType(e.target.value)}>
                                <Space direction="vertical" style={{ width: '100%' }}>
                                    <Radio.Button value={2} style={{ width: '100%' }}>
                                        <Space>
                                            <WechatOutlined className={styles.payIcon} style={{ color: '#07C160' }} />
                                            <span>微信支付</span>
                                        </Space>
                                    </Radio.Button>
                                    <Radio.Button value={1} style={{ width: '100%' }}>
                                        <Space>
                                            <AlipayOutlined className={styles.payIcon} style={{ color: '#1677FF' }} />
                                            <span>支付宝支付</span>
                                        </Space>
                                    </Radio.Button>
                                </Space>
                            </Radio.Group>

                            <Button
                                type="primary"
                                size="large"
                                className={styles.submitButton}
                                loading={loading}
                                onClick={handleSubmit}
                            >
                                立即充值 ￥{amount || 0}
                            </Button>
                        </Card>

                        <div className={styles.tips}>
                            <Alert
                                message="温馨提示"
                                description={
                                    <ul style={{ marginBottom: 0, paddingLeft: 20 }}>
                                        <li>充值金额将实时到账</li>
                                        <li>如遇到充值问题，请联系客服</li>
                                        <li>请确保支付金额与订单金额一致</li>
                                    </ul>
                                }
                                type="info"
                                showIcon
                            />
                        </div>
                    </Col>
                </Row>
            </div>
        </Body>
    );
};