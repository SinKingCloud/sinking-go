import React, {useEffect, useState} from 'react';
import {Body, Icon, Title} from "@/components";
import {Card, Row, Col, Avatar, Typography} from 'antd';
import {
    UserOutlined,
    ClusterOutlined,
    NodeIndexOutlined,
    SettingOutlined,
} from '@ant-design/icons';
import {useModel} from 'umi';
import {getOverviewInfo} from '@/service/admin/system';
import {createStyles} from 'antd-style';
import {historyPush} from "@/utils/route";
import {Logo} from "@/components/icon";

const useStyles: any = createStyles(({token, css}) => ({
    welcomeCard: css`
        background: linear-gradient(135deg, ${token.colorPrimary}, ${token.colorPrimaryActive});
        border: none;

        .ant-card-body {
            padding: 24px;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 120px;
        }

        .welcome-content {
            color: white;
            text-align: center;
        }

        .welcome-title {
            color: white !important;
            margin-bottom: 8px;
            margin-top: 0;
        }

        .welcome-text {
            color: rgba(255, 255, 255, 0.85);
            margin-bottom: 0;
        }

        @media (max-width: 768px) {
            .ant-card-body {
                padding: 16px;
            }
        }
    `,
    statCard: css`
        height: 100%;

        .ant-card-body {
            padding: 20px;
        }


        .stat-content {
            text-align: center;
        }

        .stat-number {
            font-size: 32px;
            font-weight: 600;
            line-height: 1.2;
        }

        .stat-label {
            color: ${token.colorTextSecondary};
            margin-bottom: 12px;
        }

        .stat-details {
            display: flex;
            justify-content: space-around;
            margin-top: 16px;
            padding-top: 16px;
            border-top: 1px solid ${token.colorBorderSecondary};
        }

        .stat-detail-item {
            text-align: center;
            flex: 1;
        }

        .stat-detail-number {
            font-size: 18px;
            font-weight: 500;
            display: block;
        }

        .stat-detail-label {
            font-size: 12px;
            color: ${token.colorTextTertiary};
            margin-top: 4px;
        }

        @media (max-width: 768px) {
            .stat-number {
                font-size: 24px;
            }

            .stat-detail-number {
                font-size: 16px;
            }

            .ant-card-body {
                padding: 16px;
            }
        }
    `,
    userCard: css`
        height: 100%;

        .ant-card-body {
            padding: 24px;
        }

        .user-header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            padding: 16px;
            background: linear-gradient(135deg, ${token.colorPrimaryBg}, ${token.colorBgContainer});
            border-radius: ${token.borderRadius}px;
            border: 1px solid ${token.colorBorderSecondary};
        }

        .user-avatar {
            margin-right: 16px;
            background: ${token.colorPrimary};
            box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }

        .user-info {
            flex: 1;
        }

        .user-name {
            font-weight: 600;
            font-size: 25px;
        }

        .user-status {
            color: ${token.colorTextSecondary};
        }

        .user-details {
            margin-top: 20px;
        }

        .detail-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px 16px;
            margin-bottom: 8px;
            background: ${token.colorFillAlter};
            border-radius: ${token.borderRadius}px;
            transition: all 0.2s ease;

            &:hover {
                background: ${token.colorFillContent};
            }

            &:last-child {
                margin-bottom: 0;
            }
        }

        .detail-label {
            color: ${token.colorTextSecondary};
            font-weight: 500;
            display: flex;
            align-items: center;
        }

        .detail-value {
            font-weight: 600;
            color: ${token.colorText};
        }

        @media (max-width: 768px) {
            .user-header {
                flex-direction: column;
                text-align: center;
            }

            .user-avatar {
                margin-right: 0;
                margin-bottom: 12px;
            }

            .ant-card-body {
                padding: 16px;
            }
        }
    `,
    systemCard: css`
        height: 100%;

        .ant-card-body {
            padding: 24px;
        }

        .system-header {
            display: flex;
            align-items: center;
            justify-content: center;
            margin-bottom: 20px;
            padding: 16px;
            background: linear-gradient(135deg, ${token.colorPrimaryBg}, ${token.colorBgContainer});
            border-radius: ${token.borderRadius}px;
            border: 1px solid ${token.colorBorderSecondary};
        }

        .system-icon {
            font-size: 27px;
            margin-right: 16px;
            vertical-align: top;
            color: ${token.colorPrimary};
        }

        .system-title {
            font-size: 25px;
            font-weight: 700;
            color: ${token.colorText};
            margin: 0;
        }

        .system-details {
            margin-top: 20px;
        }

        .detail-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 12px 16px;
            margin-bottom: 8px;
            background: ${token.colorFillAlter};
            border-radius: ${token.borderRadius}px;
            transition: all 0.2s ease;

            &:hover {
                background: ${token.colorFillContent};
            }

            &:last-child {
                margin-bottom: 0;
            }
        }

        .detail-label {
            color: ${token.colorTextSecondary};
            font-weight: 500;
            display: flex;
            align-items: center;
        }

        .detail-value {
            font-weight: 600;
            color: ${token.colorText};
        }

        @media (max-width: 768px) {
            .ant-card-body {
                padding: 16px;
            }
        }
    `
}));

export default (): React.ReactNode => {
    const {styles} = useStyles();
    const userModel = useModel('user');
    const webModel = useModel('web');
    const [overviewData, setOverviewData] = useState<any>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        setLoading(true);
        getOverviewInfo({
            onSuccess: (r) => {
                setOverviewData(r?.data || {});
            }
        }).finally(() => {
            setLoading(false);
        });
    }, []);

    // 格式化时间
    const formatTime = (timeStr: string) => {
        if (!timeStr) return '未知';
        return new Date(timeStr).toLocaleString('zh-CN');
    };

    return (
        <Body loading={loading}>
            {/* 欢迎卡片 */}
            <Card className={styles.welcomeCard}>
                <div className="welcome-content">
                    <Typography.Title level={3} className="welcome-title">
                        欢迎使用 {webModel?.info?.name}
                    </Typography.Title>
                    <Typography.Text className="welcome-text">
                        分布式配置管理系统，为您提供高效的服务注册与发现能力
                    </Typography.Text>
                </div>
            </Card>

            <Row gutter={[16, 16]}>
                {/* 系统概览统计 */}
                <Col xs={24} sm={8} lg={8}>
                    <Card className={styles.statCard} title={<Title>集群状态</Title>}
                          extra={<a onClick={() => {
                              historyPush("cluster")
                          }}><ClusterOutlined/>
                          </a>} variant={"borderless"}>
                        <div className="stat-content">
                            <div className="stat-number">{overviewData?.cluster?.total || 0}</div>
                            <div className="stat-label">集群总数</div>
                            <div className="stat-details">
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#52c41a'}}>
                                        {overviewData?.cluster?.online || 0}
                                    </span>
                                    <div className="stat-detail-label">在线</div>
                                </div>
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#ff4d4f'}}>
                                        {overviewData?.cluster?.offline || 0}
                                    </span>
                                    <div className="stat-detail-label">离线</div>
                                </div>
                            </div>
                        </div>
                    </Card>
                </Col>

                <Col xs={24} sm={8} lg={8}>
                    <Card className={styles.statCard} title={<Title>节点状态</Title>}
                          extra={<a onClick={() => {
                              historyPush("node")
                          }}><NodeIndexOutlined/></a>} variant={"borderless"}>
                        <div className="stat-content">
                            <div className="stat-number">{overviewData?.node?.total || 0}</div>
                            <div className="stat-label">节点总数</div>
                            <div className="stat-details">
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#52c41a'}}>
                                        {overviewData?.node?.online || 0}
                                    </span>
                                    <div className="stat-detail-label">在线</div>
                                </div>
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#ff4d4f'}}>
                                        {overviewData?.node?.offline || 0}
                                    </span>
                                    <div className="stat-detail-label">离线</div>
                                </div>
                            </div>
                        </div>
                    </Card>
                </Col>

                <Col xs={24} sm={8} lg={8}>
                    <Card className={styles.statCard} title={<Title>配置状态</Title>}
                          extra={<a onClick={() => {
                              historyPush("config")
                          }}><SettingOutlined/></a>} variant={"borderless"}>
                        <div className="stat-content">
                            <div className="stat-number">{overviewData?.config?.total || 0}</div>
                            <div className="stat-label">配置总数</div>
                            <div className="stat-details">
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#52c41a'}}>
                                        {overviewData?.config?.normal || 0}
                                    </span>
                                    <div className="stat-detail-label">正常</div>
                                </div>
                                <div className="stat-detail-item">
                                    <span className="stat-detail-number" style={{color: '#f3c402'}}>
                                        {overviewData?.config?.abnormal || 0}
                                    </span>
                                    <div className="stat-detail-label">暂停</div>
                                </div>
                            </div>
                        </div>
                    </Card>
                </Col>

                {/* 用户信息卡片 */}
                <Col xs={24} sm={24} md={12}>
                    <Card className={styles.userCard} title={<Title>账户信息</Title>} variant={"borderless"}>
                        <div className="user-header">
                            <Avatar size={64} icon={<UserOutlined/>} className="user-avatar"/>
                            <div className="user-info">
                                <div level={4} className="user-name">
                                    {userModel.web?.account || '管理员'}
                                </div>
                            </div>
                        </div>
                        <div className="user-details">
                            <div className="detail-item">
                                <span className="detail-label">
                                    <UserOutlined style={{marginRight: 8}}/>
                                    登录IP
                                </span>
                                <span className="detail-value">{userModel.web?.login_ip || '未知'}</span>
                            </div>
                            <div className="detail-item">
                                <span className="detail-label">
                                    <ClusterOutlined style={{marginRight: 8}}/>
                                    登录时间
                                </span>
                                <span className="detail-value">{formatTime(userModel.web?.login_time)}</span>
                            </div>
                        </div>
                    </Card>
                </Col>

                {/* 系统信息卡片 */}
                <Col xs={24} sm={24} md={12}>
                    <Card className={styles.systemCard} title={<Title>系统信息</Title>} variant={"borderless"}>
                        <div className="system-header">
                            <Icon type={Logo} className="system-icon"/>
                            <Typography.Title level={4}
                                              className="system-title">{webModel?.info?.name}</Typography.Title>
                        </div>
                        <div className="system-details">
                            <div className="detail-item">
                                <span className="detail-label">
                                    <NodeIndexOutlined style={{marginRight: 8}}/>
                                    监听地址
                                </span>
                                <span className="detail-value">{overviewData?.application?.listen || ":5678"}</span>
                            </div>
                            <div className="detail-item">
                                <span className="detail-label">
                                    <ClusterOutlined style={{marginRight: 8}}/>
                                    外网地址
                                </span>
                                <span className="detail-value">{overviewData?.application?.address || ":5678"}</span>
                            </div>
                            <div className="detail-item">
                                <span className="detail-label">
                                    <SettingOutlined style={{marginRight: 8}}/>
                                    部署环境
                                </span>
                                <span className="detail-value">{overviewData?.application?.mode || "release"}</span>
                            </div>
                        </div>
                    </Card>
                </Col>
            </Row>
        </Body>
    );
};
