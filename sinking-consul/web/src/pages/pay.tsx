import React from 'react';
import {Button, Col, Result, Row,} from 'antd';
import {getPayJumpUrl} from "@/utils/pay";
import {Body} from "@/components";

export default (): React.ReactNode => {
    return (
        <Body breadCrumb={false}>
            <Row type="flex" justify="center" align="start" style={{minHeight: '100vh'}}>
                <Col>
                    <Result
                        style={{
                            marginTop: "20vh"
                        }}
                        status="success"
                        title="您的订单已支付成功!"
                        subTitle="商品到账可能会有1-10分钟延迟, 如超时未到账请联系客服。"
                        extra={[
                            <Button type="primary" key="console" onClick={() => {
                                window.location.href = getPayJumpUrl() || "/";
                            }}>
                                返回上一页
                            </Button>,
                        ]}
                    />
                </Col>
            </Row>
        </Body>
    );
};
