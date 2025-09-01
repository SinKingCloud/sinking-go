import React from 'react';
import { Body } from "@/components";
import { createStyles } from "antd-style";
import { Row, Col, Card, Button, Typography, Space, Statistic, Avatar, Divider, Menu } from "antd";
import { 
    TrophyOutlined, 
    SafetyOutlined, 
    ThunderboltOutlined, 
    TeamOutlined,
    BarChartOutlined,
    BellOutlined,
    CheckCircleOutlined,
    ArrowRightOutlined,
    HomeOutlined,
    CustomerServiceOutlined,
    QuestionCircleOutlined,
    FileTextOutlined,
    UserOutlined,
    LoginOutlined
} from "@ant-design/icons";
import {historyPush} from "@/utils/route";

const { Title, Paragraph, Text } = Typography;

const useStyles = createStyles(({ css, token, responsive, isDarkMode }): any => {
    const primaryColor = token.colorPrimary;
    const textColor = isDarkMode ? "#fff" : "rgba(0, 0, 0, 0.85)";
    const secondaryColor = isDarkMode ? "rgba(255, 255, 255, 0.65)" : "rgba(0, 0, 0, 0.65)";
    
    return {
        navbar: css`
            background: ${token.colorBgContainer};
            border-bottom: 1px solid ${token.colorBorderSecondary};
            padding: 0 20px;
            position: sticky;
            top: 0;
            z-index: 1000;
            backdrop-filter: blur(8px);
            box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
        `,
        navContent: css`
            max-width: 1200px;
            margin: 0 auto;
            display: flex;
            align-items: center;
            justify-content: space-between;
            height: 64px;
        `,
        logo: css`
            display: flex;
            align-items: center;
            font-size: 20px;
            font-weight: 700;
            color: ${primaryColor};
            text-decoration: none;
            
            .anticon {
                margin-right: 8px;
                font-size: 24px;
            }
            
            ${responsive.sm} {
                font-size: 18px;
                
                .anticon {
                    font-size: 20px;
                }
            }
        `,
        navMenu: css`
            border: none;
            background: transparent;
            
            .ant-menu-item {
                color: ${textColor};
                font-weight: 500;
                
                &:hover {
                    color: ${primaryColor};
                }
            }
            
            .ant-menu-item-selected {
                color: ${primaryColor};
                background: transparent;
                
                &::after {
                    border-bottom-color: ${primaryColor};
                }
            }
            
            ${responsive.md} {
                display: none;
            }
        `,
        navActions: css`
            display: flex;
            align-items: center;
            gap: 12px;
            
            .ant-btn {
                border-radius: 6px;
                font-weight: 500;
            }
            
            ${responsive.sm} {
                gap: 8px;
                
                .ant-btn {
                    padding: 0 12px;
                }
            }
        `,
        mobileMenu: css`
            display: none;
            
            ${responsive.md} {
                display: block;
            }
        `,
        hero: css`
            background: linear-gradient(135deg, ${primaryColor}15 0%, ${primaryColor}08 100%);
            padding: 60px 20px;
            text-align: center;
            border-radius: 12px;
            margin-bottom: 40px;
            margin-top: 20px;
            
            ${responsive.md} {
                padding: 50px 20px;
                margin-top: 10px;
            }
        `,
        heroTitle: css`
            font-size: 48px;
            font-weight: 700;
            color: ${textColor};
            margin-bottom: 20px;
            line-height: 1.2;
            
            ${responsive.md} {
                font-size: 36px;
            }
            
            ${responsive.sm} {
                font-size: 28px;
            }
        `,
        heroSubtitle: css`
            font-size: 20px;
            color: ${secondaryColor};
            margin-bottom: 40px;
            max-width: 600px;
            margin-left: auto;
            margin-right: auto;
            
            ${responsive.md} {
                font-size: 18px;
            }
        `,
        ctaButtons: css`
            .ant-btn {
                height: 50px;
                font-size: 16px;
                padding: 0 30px;
                margin: 0 10px;
                border-radius: 8px;
                
                ${responsive.sm} {
                    width: 100%;
                    margin: 5px 0;
                }
            }
        `,
        featureCard: css`
            text-align: center;
            padding: 30px 20px;
            border: 1px solid ${token.colorBorderSecondary};
            border-radius: 12px;
            transition: all 0.3s ease;
            height: 100%;
            
            &:hover {
                border-color: ${primaryColor};
                box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1);
                transform: translateY(-2px);
            }
        `,
        featureIcon: css`
            font-size: 48px;
            color: ${primaryColor};
            margin-bottom: 20px;
        `,
        featureTitle: css`
            font-size: 18px;
            font-weight: 600;
            color: ${textColor};
            margin-bottom: 15px;
        `,
        featureDesc: css`
            color: ${secondaryColor};
            line-height: 1.6;
        `,
        statsSection: css`
            background: ${token.colorBgContainer};
            padding: 60px 20px;
            border-radius: 12px;
            text-align: center;
            margin: 40px 0;
        `,
        sectionTitle: css`
            font-size: 32px;
            font-weight: 600;
            color: ${textColor};
            margin-bottom: 50px;
            text-align: center;
            
            ${responsive.md} {
                font-size: 28px;
            }
        `,
        advantageCard: css`
            padding: 30px;
            border-radius: 12px;
            background: ${token.colorBgContainer};
            border: 1px solid ${token.colorBorderSecondary};
            height: 100%;
            
            .ant-typography {
                margin-bottom: 10px;
            }
        `,
        advantageIcon: css`
            font-size: 24px;
            color: ${primaryColor};
            margin-right: 10px;
        `,
        pricingCard: css`
            padding: 40px 30px;
            border-radius: 12px;
            text-align: center;
            border: 2px solid ${token.colorBorderSecondary};
            transition: all 0.3s ease;
            height: 100%;
            
            &:hover {
                border-color: ${primaryColor};
                transform: translateY(-5px);
                box-shadow: 0 15px 40px rgba(0, 0, 0, 0.1);
            }
            
            &.featured {
                border-color: ${primaryColor};
                position: relative;
                
                &::before {
                    content: "推荐";
                    position: absolute;
                    top: -10px;
                    right: 20px;
                    background: ${primaryColor};
                    color: white;
                    padding: 5px 15px;
                    border-radius: 20px;
                    font-size: 12px;
                }
            }
        `,
        footer: css`
            background: ${token.colorBgContainer};
            padding: 40px 20px;
            text-align: center;
            border-radius: 12px;
            margin-top: 40px;
            border-top: 1px solid ${token.colorBorderSecondary};
        `
    };
});

const LandingPage: React.FC = () => {
    const { styles } = useStyles();

    const menuItems = [
        {
            key: 'home',
            icon: <HomeOutlined />,
            label: '首页',
        },
        {
            key: 'features',
            icon: <BarChartOutlined />,
            label: '功能特色',
        },
        {
            key: 'pricing',
            icon: <FileTextOutlined />,
            label: '订阅方案',
        },
        {
            key: 'help',
            icon: <QuestionCircleOutlined />,
            label: '帮助中心',
        },
        {
            key: 'contact',
            icon: <CustomerServiceOutlined />,
            label: '联系我们',
        }
    ];

    const features = [
        {
            icon: <ThunderboltOutlined className={styles.featureIcon} />,
            title: "实时策略推送",
            description: "毫秒级策略信号推送，确保您第一时间获得交易机会，不错过任何盈利时机"
        },
        {
            icon: <SafetyOutlined className={styles.featureIcon} />,
            title: "安全可靠",
            description: "银行级安全加密，多重风控体系，保障您的交易数据和资金安全"
        },
        {
            icon: <BarChartOutlined className={styles.featureIcon} />,
            title: "专业分析",
            description: "由资深交易员和量化团队提供的专业策略分析，助您提高交易胜率"
        },
        {
            icon: <BellOutlined className={styles.featureIcon} />,
            title: "多端同步",
            description: "支持手机、电脑、MT5平台多端同步，随时随地接收交易信号"
        }
    ];

    const advantages = [
        "✓ 7×24小时不间断监控市场",
        "✓ 支持外汇、黄金、原油等多种品种",
        "✓ 智能风险控制系统",
        "✓ 历史收益率透明公开",
        "✓ 专业客服团队支持",
        "✓ 30天免费试用期"
    ];

    return (
        <Body breadCrumb={false}>
            {/* Navigation Bar */}
            <div className={styles.navbar}>
                <div className={styles.navContent}>
                    <div className={styles.logo}>
                        <BarChartOutlined />
                        MT5策略系统
                    </div>
                    
                    <Menu 
                        mode="horizontal" 
                        items={menuItems}
                        className={styles.navMenu}
                        selectedKeys={['home']}
                    />
                    
                    <div className={styles.navActions}>
                        <Button 
                            icon={<UserOutlined />}
                            onClick={()=>{
                                historyPush("login")
                            }}
                        >
                            登录
                        </Button>
                        <Button 
                            type="primary" 
                            icon={<LoginOutlined />}
                            href="/user/login"
                        >
                            注册
                        </Button>
                    </div>
                </div>
            </div>

            {/* Hero Section */}
            <div className={styles.hero}>
                <Title className={styles.heroTitle}>
                    MetaTrader5 策略订阅系统
                </Title>
                <Paragraph className={styles.heroSubtitle}>
                    专业的外汇交易策略订阅平台，为您提供实时、精准的交易信号
                    让每一次交易都更加智能和profitable
                </Paragraph>
                <div className={styles.ctaButtons}>
                    <Button type="primary" size="large" icon={<ArrowRightOutlined />}>
                        立即开始
                    </Button>
                    <Button size="large">
                        了解更多
                    </Button>
                </div>
            </div>

            {/* Features Section */}
            <Title className={styles.sectionTitle}>为什么选择我们</Title>
            <Row gutter={[24, 24]} style={{ marginBottom: 60 }}>
                {features.map((feature, index) => (
                    <Col xl={6} lg={8} md={12} sm={24} xs={24} key={index}>
                        <Card className={styles.featureCard} bordered={false}>
                            {feature.icon}
                            <Title level={4} className={styles.featureTitle}>
                                {feature.title}
                            </Title>
                            <Paragraph className={styles.featureDesc}>
                                {feature.description}
                            </Paragraph>
                        </Card>
                    </Col>
                ))}
            </Row>

            {/* Stats Section */}
            <div className={styles.statsSection}>
                <Title className={styles.sectionTitle}>平台数据</Title>
                <Row gutter={[24, 24]}>
                    <Col xl={6} lg={12} md={12} sm={24} xs={24}>
                        <Statistic 
                            title="活跃用户" 
                            value={12000} 
                            suffix="+" 
                            valueStyle={{ color: '#1890ff', fontSize: '32px' }}
                        />
                    </Col>
                    <Col xl={6} lg={12} md={12} sm={24} xs={24}>
                        <Statistic 
                            title="策略成功率" 
                            value={85.6} 
                            suffix="%" 
                            valueStyle={{ color: '#52c41a', fontSize: '32px' }}
                        />
                    </Col>
                    <Col xl={6} lg={12} md={12} sm={24} xs={24}>
                        <Statistic 
                            title="月均收益率" 
                            value={15.8} 
                            suffix="%" 
                            valueStyle={{ color: '#f5222d', fontSize: '32px' }}
                        />
                    </Col>
                    <Col xl={6} lg={12} md={12} sm={24} xs={24}>
                        <Statistic 
                            title="累计信号" 
                            value={50000} 
                            suffix="+" 
                            valueStyle={{ color: '#722ed1', fontSize: '32px' }}
                        />
                    </Col>
                </Row>
            </div>

            {/* Advantages Section */}
            <Title className={styles.sectionTitle}>核心优势</Title>
            <Row gutter={[24, 24]} style={{ marginBottom: 60 }}>
                <Col xl={12} lg={24} md={24} sm={24} xs={24}>
                    <Card className={styles.advantageCard}>
                        <Title level={3}>
                            <TrophyOutlined className={styles.advantageIcon} />
                            专业团队保障
                        </Title>
                        <Space direction="vertical" size="small" style={{ width: '100%' }}>
                            {advantages.map((advantage, index) => (
                                <Text key={index} style={{ display: 'block', fontSize: '16px' }}>
                                    {advantage}
                                </Text>
                            ))}
                        </Space>
                    </Card>
                </Col>
                <Col xl={12} lg={24} md={24} sm={24} xs={24}>
                    <Card className={styles.advantageCard}>
                        <Title level={3}>
                            <TeamOutlined className={styles.advantageIcon} />
                            用户评价
                        </Title>
                        <Space direction="vertical" size="large" style={{ width: '100%' }}>
                            <div>
                                <Avatar size="large" src="https://joeschmoe.io/api/v1/random">U</Avatar>
                                <div style={{ marginTop: 10 }}>
                                    <Text strong>交易者张先生</Text>
                                    <Paragraph style={{ margin: '5px 0 0 0' }}>
                                        "使用这个平台3个月，收益率提升了40%，信号准确及时，非常满意！"
                                    </Paragraph>
                                </div>
                            </div>
                            <Divider />
                            <div>
                                <Avatar size="large" src="https://joeschmoe.io/api/v1/random">L</Avatar>
                                <div style={{ marginTop: 10 }}>
                                    <Text strong>投资人李女士</Text>
                                    <Paragraph style={{ margin: '5px 0 0 0' }}>
                                        "专业的策略分析和及时的客服支持，让我的交易更加自信和profitable。"
                                    </Paragraph>
                                </div>
                            </div>
                        </Space>
                    </Card>
                </Col>
            </Row>

            {/* Pricing Section */}
            <Title className={styles.sectionTitle}>选择您的订阅方案</Title>
            <Row gutter={[24, 24]} style={{ marginBottom: 60 }}>
                <Col xl={8} lg={8} md={24} sm={24} xs={24}>
                    <Card className={styles.pricingCard}>
                        <Title level={4}>基础版</Title>
                        <Title level={2} style={{ color: '#1890ff' }}>
                            ¥299<small style={{ fontSize: '14px' }}>/月</small>
                        </Title>
                        <Space direction="vertical" size="small" style={{ width: '100%', marginBottom: 20 }}>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 基础策略信号</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 5个交易品种</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 邮件推送</Text>
                        </Space>
                        <Button size="large" style={{ width: '100%' }}>选择基础版</Button>
                    </Card>
                </Col>
                <Col xl={8} lg={8} md={24} sm={24} xs={24}>
                    <Card className={`${styles.pricingCard} featured`}>
                        <Title level={4}>专业版</Title>
                        <Title level={2} style={{ color: '#1890ff' }}>
                            ¥599<small style={{ fontSize: '14px' }}>/月</small>
                        </Title>
                        <Space direction="vertical" size="small" style={{ width: '100%', marginBottom: 20 }}>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 高级策略信号</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 15个交易品种</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 实时推送</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 专家分析报告</Text>
                        </Space>
                        <Button type="primary" size="large" style={{ width: '100%' }}>选择专业版</Button>
                    </Card>
                </Col>
                <Col xl={8} lg={8} md={24} sm={24} xs={24}>
                    <Card className={styles.pricingCard}>
                        <Title level={4}>旗舰版</Title>
                        <Title level={2} style={{ color: '#1890ff' }}>
                            ¥999<small style={{ fontSize: '14px' }}>/月</small>
                        </Title>
                        <Space direction="vertical" size="small" style={{ width: '100%', marginBottom: 20 }}>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 全套策略信号</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 所有交易品种</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 多端同步</Text>
                            <Text><CheckCircleOutlined style={{ color: '#52c41a' }} /> 1对1专属顾问</Text>
                        </Space>
                        <Button size="large" style={{ width: '100%' }}>选择旗舰版</Button>
                    </Card>
                </Col>
            </Row>

            {/* Footer */}
            <div className={styles.footer}>
                <Title level={4}>立即开始您的盈利之旅</Title>
                <Paragraph>
                    加入我们的专业交易社区，获得稳定收益
                </Paragraph>
                <Button type="primary" size="large" icon={<ArrowRightOutlined />}>
                    免费试用30天
                </Button>
            </div>
        </Body>
    );
};

export default LandingPage;
