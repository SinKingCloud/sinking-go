import React from 'react';
import {Body, Message, Title} from "@/components";
import {createStyles} from "antd-style";
import {useModel} from "umi";
import {App, Avatar, Button, Card, Col, Row, Skeleton, Statistic} from "antd";

const useStyles = createStyles(({css, responsive, isDarkMode}): any => {
    const color = isDarkMode ? "#fff" : "rgba(0, 0, 0, 0.85)"
    return {
        pageHeaderContent: css`
            display: flex;
            align-items: center;
            padding: 10px;
        `,
        right: css`
            display: flex;
            align-items: center;
        `,
        align: css`
            width: 100%;
            display: flex;
            align-items: center;

            ${responsive.md} {
                display: flex;
                flex-direction: column;

                .ant-col-16 {
                    max-width: 100%;
                    display: flex;
                    flex-wrap: wrap;
                    justify-content: center
                }

                .ant-col-8 {
                    max-width: 100%;
                }
            }
        `,
        avatar: css`
            height: 65px;
            width: 65px;
            margin-right: 15px;
            flex-shrink: 0;

            ${responsive.md && responsive.xl && responsive.lg && responsive.sm} {
                margin-bottom: 10px;
            }
        `,
        box: css`
            height: 65px;
        `,
        title: css`
            margin-top: 8px;
            font-size: 20px;
            color: ${color};
            line-height: 1;

            ${responsive.md && responsive.xl && responsive.lg && responsive.sm} {
                text-align: center;
                margin-top: 0;
                line-height: 1;
            }
        `,
        content: css`
            font-size: 15px;
            line-height: 1;
            color: ${color};

            ${responsive.md && responsive.xl && responsive.lg && responsive.sm} {
                line-height: 20px;
            }
        `,
        money: css`
            margin-bottom: 4px;
            color: rgba(0, 0, 0, 0.45);
            font-size: 14px;
            padding: 0 30px;
            box-sizing: border-box;
            display: flex;
            align-items: center;

            ${responsive.md} {
                padding: 0 30px;
                box-sizing: border-box;
                margin-top: 10px;

                .ant-statistic-title {
                    font-size: 14px;
                }

                .ant-statistic-content {
                    font-size: 22px !important;
                }
            }
        `,
    };
});

const PageHeaderContent: React.FC = () => {
    const {
        styles: {box, avatar, title, content}
    } = useStyles();
    const user = useModel("user");
    const loading = user && Object.keys(user).length;
    if (!loading) {
        return <Skeleton avatar paragraph={{rows: 1}} active/>;
    }
    return (
        <>
            <Avatar src={user?.web?.avatar} className={avatar}/><br/>
            <div className={box}>
                <div className={title}>
                    你好，
                    {user?.web?.nick_name}！
                </div>
                <br/>
                <div style={{maxWidth: "450px"}} className={content}>
                    尊敬的 <b>{user?.web?.nick_name}</b> ,欢迎回来，我们已等候多时！
                </div>
            </div>
        </>
    );
};
const ExtraContent: React.FC = () => {
    const user = useModel("user")
    const {styles: {money}} = useStyles();
    return (
        <>
            <div className={money}>
                <Statistic title="余额" prefix={"￥"} value={parseFloat(user?.web?.money || 0).toFixed(2)}/>
            </div>
            <div className={money}>
                <Statistic title="身份"
                           value={user?.web?.is_master ? '管理员' : (user?.web?.is_admin ? '站长' : '会员')}/>
            </div>
        </>
    )
};


export default (): React.ReactNode => {
    const {
        styles: {
            align,
            pageHeaderContent,
            right,
        }
    } = useStyles();
    const {message, modal} = App.useApp()
    return (
        <Body space={false}>
            <Row gutter={[10, 10]}>
                <Col span={24}>
                    <Card variant={"borderless"}>
                        <Row className={align}>
                            <Col span={16} className={pageHeaderContent}>
                                <PageHeaderContent/>
                            </Col>
                            <Col span={8} className={right}>
                                <ExtraContent/>
                            </Col>
                        </Row>
                    </Card>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                <Button onClick={() => {
                                    message.success("测试按钮点击事件");
                                }}>测试</Button>
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
                <Col span={24}>
                    <Row gutter={[10, 10]}>
                        <Col xl={16} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>我的订阅</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                        <Col xl={8} lg={24} md={24} sm={24} xs={24}>
                            <Card title={<Title>系统公告</Title>} variant={"borderless"}>
                                测试
                            </Card>
                        </Col>
                    </Row>
                </Col>
            </Row>
        </Body>
    );
};
