import React, { useState } from 'react';
import { Body, Message } from '@/components';
import { 
    Form, 
    Input, 
    Button, 
    Upload, 
    Avatar, 
    Space,
    Tabs,
    Switch,
    Slider,
    Select
} from 'antd';
import { 
    UserOutlined, 
    PhoneOutlined, 
    LockOutlined, 
    CameraOutlined,
    EditOutlined,
    SaveOutlined,
    SafetyOutlined,
    BellOutlined,
    EyeOutlined,
    GlobalOutlined,
    MailOutlined,
    WechatOutlined,
    QqOutlined,
    GithubOutlined,
    TwitterOutlined
} from '@ant-design/icons';
import { createStyles } from 'antd-style';
import { useModel } from 'umi';

const { TextArea } = Input;
const { TabPane } = Tabs;
const { Option } = Select;

const useStyles = createStyles(({ css, isDarkMode, token }) => ({
    pageWrapper: css`
        min-height: 100vh;
        background: ${isDarkMode 
            ? 'linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%)' 
            : 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
        };
        padding: 0;
        position: relative;
        overflow-x: hidden;
        
        &::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: 
                radial-gradient(circle at 20% 80%, rgba(120, 119, 198, 0.3) 0%, transparent 50%),
                radial-gradient(circle at 80% 20%, rgba(255, 119, 198, 0.3) 0%, transparent 50%),
                radial-gradient(circle at 40% 40%, rgba(120, 219, 255, 0.3) 0%, transparent 50%);
            pointer-events: none;
        }
    `,
    container: css`
        max-width: 1200px;
        margin: 0 auto;
        padding: 40px 20px;
        position: relative;
        z-index: 1;
        
        @media (max-width: 768px) {
            padding: 20px 16px;
        }
    `,
    profileCard: css`
        background: ${isDarkMode 
            ? 'rgba(255, 255, 255, 0.05)' 
            : 'rgba(255, 255, 255, 0.95)'
        };
        backdrop-filter: blur(20px);
        border-radius: 24px;
        border: 1px solid ${isDarkMode 
            ? 'rgba(255, 255, 255, 0.1)' 
            : 'rgba(255, 255, 255, 0.3)'
        };
        box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
        overflow: hidden;
        transition: all 0.3s ease;
        
        &:hover {
            transform: translateY(-5px);
            box-shadow: 0 30px 60px rgba(0, 0, 0, 0.15);
        }
    `,
    profileHeader: css`
        padding: 40px;
        text-align: center;
        position: relative;
        
        @media (max-width: 768px) {
            padding: 30px 20px;
        }
        
        &::before {
            content: '';
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            height: 120px;
            background: linear-gradient(135deg, 
                ${isDarkMode ? '#667eea' : '#667eea'} 0%, 
                ${isDarkMode ? '#764ba2' : '#764ba2'} 100%
            );
            border-radius: 24px 24px 0 0;
            z-index: -1;
        }
    `,
    avatarWrapper: css`
        position: relative;
        display: inline-block;
        margin-bottom: 20px;
        
        .avatar-container {
            position: relative;
            
            .ant-avatar {
                border: 4px solid rgba(255, 255, 255, 0.9);
                box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
                transition: all 0.3s ease;
                
                &:hover {
                    transform: scale(1.05);
                    box-shadow: 0 15px 40px rgba(0, 0, 0, 0.3);
                }
            }
            
            .upload-btn {
                position: absolute;
                bottom: 5px;
                right: 5px;
                width: 36px;
                height: 36px;
                border-radius: 50%;
                background: linear-gradient(135deg, #ff6b6b, #ee5a24);
                border: 3px solid white;
                display: flex;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                transition: all 0.3s ease;
                box-shadow: 0 4px 12px rgba(238, 90, 36, 0.4);
                
                &:hover {
                    transform: scale(1.1);
                    box-shadow: 0 6px 20px rgba(238, 90, 36, 0.6);
                }
                
                .anticon {
                    color: white;
                    font-size: 14px;
                }
            }
        }
    `,
    userInfo: css`
        color: white;
        
        .username {
            font-size: 28px;
            font-weight: 700;
            margin-bottom: 8px;
            text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
            
            @media (max-width: 768px) {
                font-size: 24px;
            }
        }
        
        .user-bio {
            font-size: 16px;
            opacity: 0.9;
            margin-bottom: 20px;
            max-width: 400px;
            margin-left: auto;
            margin-right: auto;
            line-height: 1.5;
        }
        
        .social-links {
            display: flex;
            justify-content: center;
            gap: 16px;
            
            .social-btn {
                width: 40px;
                height: 40px;
                border-radius: 50%;
                background: rgba(255, 255, 255, 0.2);
                border: 1px solid rgba(255, 255, 255, 0.3);
                display: flex;
                align-items: center;
                justify-content: center;
                color: white;
                font-size: 16px;
                transition: all 0.3s ease;
                cursor: pointer;
                
                &:hover {
                    background: rgba(255, 255, 255, 0.3);
                    transform: translateY(-2px);
                }
            }
        }
    `,
    tabsContainer: css`
        padding: 0 40px 40px;
        
        @media (max-width: 768px) {
            padding: 0 20px 30px;
        }
        
        .ant-tabs {
            .ant-tabs-tab {
                font-size: 16px;
                font-weight: 600;
                padding: 12px 24px;
                
                @media (max-width: 768px) {
                    font-size: 14px;
                    padding: 10px 16px;
                }
            }
            
            .ant-tabs-tab-active {
                color: ${token.colorPrimary} !important;
            }
            
            .ant-tabs-ink-bar {
                background: ${token.colorPrimary};
                height: 3px;
                border-radius: 2px;
            }
        }
    `,
    formSection: css`
        padding: 30px 0;
        
        .form-group {
            margin-bottom: 32px;
            
            .group-title {
                font-size: 18px;
                font-weight: 600;
                color: ${token.colorText};
                margin-bottom: 20px;
                display: flex;
                align-items: center;
                gap: 8px;
                
                .anticon {
                    color: ${token.colorPrimary};
                }
            }
            
            .form-grid {
                display: grid;
                grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
                gap: 20px;
                
                @media (max-width: 768px) {
                    grid-template-columns: 1fr;
                }
            }
        }
        
        .ant-form-item {
            margin-bottom: 24px;
            
            .ant-form-item-label > label {
                font-weight: 500;
                color: ${token.colorText};
                font-size: 14px;
            }
        }
        
        .ant-input, .ant-input-password, .ant-select-selector, .ant-input-affix-wrapper {
            height: 48px;
            border-radius: 12px;
            border: 2px solid ${isDarkMode ? 'rgba(255,255,255,0.1)' : '#e8e8e8'};
            background: ${isDarkMode ? 'rgba(255,255,255,0.05)' : '#fafafa'};
            transition: all 0.3s ease;
            
            &:hover, &:focus {
                border-color: ${token.colorPrimary};
                background: ${isDarkMode ? 'rgba(255,255,255,0.08)' : '#ffffff'};
                box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
            }
        }
        
        .ant-input-affix-wrapper {
            padding: 0 16px;
            
            .ant-input {
                height: auto;
                border: none;
                background: transparent;
                box-shadow: none;
                
                &:focus {
                    box-shadow: none;
                }
            }
        }
        
        textarea.ant-input {
            height: auto;
            min-height: 100px;
            padding: 12px 16px;
            resize: vertical;
        }
    `,
    actionButtons: css`
        display: flex;
        gap: 16px;
        justify-content: flex-end;
        margin-top: 32px;
        
        @media (max-width: 768px) {
            flex-direction: column;
        }
        
        .primary-btn {
            height: 48px;
            padding: 0 32px;
            border-radius: 12px;
            font-weight: 600;
            font-size: 16px;
            background: linear-gradient(135deg, ${token.colorPrimary}, ${token.colorPrimaryHover});
            border: none;
            box-shadow: 0 4px 16px rgba(24, 144, 255, 0.3);
            transition: all 0.3s ease;
            
            &:hover {
                transform: translateY(-2px);
                box-shadow: 0 8px 24px rgba(24, 144, 255, 0.4);
            }
        }
        
        .secondary-btn {
            height: 48px;
            padding: 0 32px;
            border-radius: 12px;
            font-weight: 600;
            font-size: 16px;
            background: transparent;
            border: 2px solid ${token.colorBorder};
            color: ${token.colorText};
            transition: all 0.3s ease;
            
            &:hover {
                border-color: ${token.colorPrimary};
                color: ${token.colorPrimary};
                transform: translateY(-2px);
            }
        }
    `,
    settingItem: css`
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 16px 0;
        border-bottom: 1px solid ${isDarkMode ? 'rgba(255,255,255,0.1)' : '#f0f0f0'};
        
        &:last-child {
            border-bottom: none;
        }
        
        .setting-info {
            .setting-title {
                font-weight: 600;
                color: ${token.colorText};
                margin-bottom: 4px;
            }
            
            .setting-desc {
                font-size: 13px;
                color: ${token.colorTextSecondary};
            }
        }
    `,
    verifySection: css`
        .verify-input {
            display: flex;
            gap: 12px;
            
            .ant-input {
                flex: 1;
            }
            
            .verify-btn {
                height: 48px;
                padding: 0 24px;
                border-radius: 12px;
                font-weight: 600;
                background: linear-gradient(135deg, #52c41a, #73d13d);
                border: none;
                color: white;
                white-space: nowrap;
                
                &:hover {
                    background: linear-gradient(135deg, #73d13d, #52c41a);
                }
                
                &:disabled {
                    background: #d9d9d9;
                }
            }
        }
    `,
    statsGrid: css`
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
        gap: 20px;
        margin: 30px 0;
        
        .stat-card {
            background: ${isDarkMode ? 'rgba(255,255,255,0.05)' : 'rgba(255,255,255,0.8)'};
            border-radius: 16px;
            padding: 24px;
            text-align: center;
            border: 1px solid ${isDarkMode ? 'rgba(255,255,255,0.1)' : 'rgba(0,0,0,0.05)'};
            transition: all 0.3s ease;
            
            &:hover {
                transform: translateY(-4px);
                box-shadow: 0 8px 24px rgba(0,0,0,0.1);
            }
            
            .stat-number {
                font-size: 32px;
                font-weight: 800;
                color: ${token.colorPrimary};
                margin-bottom: 8px;
                display: block;
            }
            
            .stat-label {
                font-size: 14px;
                color: ${token.colorTextSecondary};
                font-weight: 500;
            }
        }
    `
}));

export default (): React.ReactNode => {
    const { styles } = useStyles();
    const [activeTab, setActiveTab] = useState('profile');
    const [profileForm] = Form.useForm();
    const [securityForm] = Form.useForm();
    const [notificationForm] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [avatarUrl, setAvatarUrl] = useState('');
    const [countdown, setCountdown] = useState(0);
    const user = useModel('user');

    // 头像上传
    const handleAvatarChange = (info: any) => {
        if (info.file.status === 'done') {
            setAvatarUrl(info.file.response?.url || '');
            Message.success('头像上传成功');
        } else if (info.file.status === 'error') {
            Message.error('头像上传失败');
        }
    };

    // 发送验证码
    const handleSendCode = () => {
        setCountdown(60);
        Message.success('验证码已发送');
        const timer = setInterval(() => {
            setCountdown((prev) => {
                if (prev <= 1) {
                    clearInterval(timer);
                    return 0;
                }
                return prev - 1;
            });
        }, 1000);
    };

    // 保存个人资料
    const handleProfileSave = async (values: any) => {
        setLoading(true);
        try {
            console.log('保存个人资料:', values);
            Message.success('个人资料保存成功');
        } catch (error) {
            Message.error('保存失败');
        } finally {
            setLoading(false);
        }
    };

    // 保存安全设置
    const handleSecuritySave = async (values: any) => {
        setLoading(true);
        try {
            console.log('保存安全设置:', values);
            Message.success('安全设置保存成功');
        } catch (error) {
            Message.error('保存失败');
        } finally {
            setLoading(false);
        }
    };

    // 保存通知设置
    const handleNotificationSave = async (values: any) => {
        setLoading(true);
        try {
            console.log('保存通知设置:', values);
            Message.success('通知设置保存成功');
        } catch (error) {
            Message.error('保存失败');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className={styles.pageWrapper}>
            <Body>
                <div className={styles.container}>
                    <div className={styles.profileCard}>
                        {/* 个人资料头部 */}
                        <div className={styles.profileHeader}>
                            <div className={styles.avatarWrapper}>
                                <Upload
                                    name="avatar"
                                    listType="picture-card"
                                    showUploadList={false}
                                    action="/api/upload/avatar"
                                    onChange={handleAvatarChange}
                                >
                                    <div className="avatar-container">
                                        <Avatar 
                                            size={120} 
                                            src={avatarUrl || user?.web?.avatar || 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=200&h=200&fit=crop&crop=face'} 
                                            icon={<UserOutlined />}
                                        />
                                        <div className="upload-btn">
                                            <CameraOutlined />
                                        </div>
                                    </div>
                                </Upload>
                            </div>
                            
                            <div className={styles.userInfo}>
                                <div className="username">
                                    {user?.web?.nickname || '用户名'}
                                </div>
                                <div className="user-bio">
                                    探索世界，记录生活，分享美好时光 ✨
                                </div>
                                <div className="social-links">
                                    <div className="social-btn">
                                        <WechatOutlined />
                                    </div>
                                    <div className="social-btn">
                                        <QqOutlined />
                                    </div>
                                    <div className="social-btn">
                                        <GithubOutlined />
                                    </div>
                                    <div className="social-btn">
                                        <TwitterOutlined />
                                    </div>
                                </div>
                            </div>
                        </div>

                        {/* 统计数据 */}
                        <div className={styles.statsGrid}>
                            <div className="stat-card">
                                <span className="stat-number">128</span>
                                <span className="stat-label">活跃天数</span>
                            </div>
                            <div className="stat-card">
                                <span className="stat-number">2.8K</span>
                                <span className="stat-label">获得点赞</span>
                            </div>
                            <div className="stat-card">
                                <span className="stat-number">95%</span>
                                <span className="stat-label">安全指数</span>
                            </div>
                            <div className="stat-card">
                                <span className="stat-number">LV.5</span>
                                <span className="stat-label">用户等级</span>
                            </div>
                        </div>

                        {/* 设置选项卡 */}
                        <div className={styles.tabsContainer}>
                            <Tabs activeKey={activeTab} onChange={setActiveTab}>
                                <TabPane tab="个人资料" key="profile">
                                    <div className={styles.formSection}>
                                        <Form
                                            form={profileForm}
                                            layout="vertical"
                                            onFinish={handleProfileSave}
                                            initialValues={{
                                                nickname: user?.web?.nickname || '',
                                                email: user?.web?.email || '',
                                                bio: '探索世界，记录生活，分享美好时光 ✨',
                                                location: '中国',
                                                website: '',
                                                birthday: '',
                                                gender: 'private'
                                            }}
                                        >
                                            <div className="form-group">
                                                <div className="group-title">
                                                    <UserOutlined />
                                                    基本信息
                                                </div>
                                                <div className="form-grid">
                                                    <Form.Item
                                                        label="昵称"
                                                        name="nickname"
                                                        rules={[
                                                            { required: true, message: '请输入昵称' },
                                                            { min: 2, max: 20, message: '昵称长度为2-20个字符' }
                                                        ]}
                                                    >
                                                        <Input 
                                                            prefix={<EditOutlined />} 
                                                            placeholder="请输入您的昵称" 
                                                        />
                                                    </Form.Item>
                                                    <Form.Item
                                                        label="邮箱"
                                                        name="email"
                                                        rules={[
                                                            { type: 'email', message: '请输入正确的邮箱格式' }
                                                        ]}
                                                    >
                                                        <Input 
                                                            prefix={<MailOutlined />} 
                                                            placeholder="请输入邮箱地址" 
                                                        />
                                                    </Form.Item>
                                                </div>
                                            </div>

                                            <div className="form-group">
                                                <div className="group-title">
                                                    <EditOutlined />
                                                    详细信息
                                                </div>
                                                <Form.Item
                                                    label="个人简介"
                                                    name="bio"
                                                >
                                                    <TextArea 
                                                        placeholder="介绍一下自己吧..." 
                                                        rows={4}
                                                        maxLength={200}
                                                        showCount
                                                    />
                                                </Form.Item>
                                                <div className="form-grid">
                                                    <Form.Item
                                                        label="所在地"
                                                        name="location"
                                                    >
                                                        <Input 
                                                            prefix={<GlobalOutlined />} 
                                                            placeholder="请输入所在地" 
                                                        />
                                                    </Form.Item>
                                                    <Form.Item
                                                        label="个人网站"
                                                        name="website"
                                                    >
                                                        <Input 
                                                            prefix={<GlobalOutlined />} 
                                                            placeholder="https://example.com" 
                                                        />
                                                    </Form.Item>
                                                    <Form.Item
                                                        label="生日"
                                                        name="birthday"
                                                    >
                                                        <Input 
                                                            type="date"
                                                        />
                                                    </Form.Item>
                                                    <Form.Item
                                                        label="性别"
                                                        name="gender"
                                                    >
                                                        <Select placeholder="请选择性别">
                                                            <Option value="male">男</Option>
                                                            <Option value="female">女</Option>
                                                            <Option value="private">保密</Option>
                                                        </Select>
                                                    </Form.Item>
                                                </div>
                                            </div>

                                            <div className={styles.actionButtons}>
                                                <Button className="secondary-btn">
                                                    重置
                                                </Button>
                                                <Button 
                                                    type="primary" 
                                                    htmlType="submit"
                                                    loading={loading}
                                                    className="primary-btn"
                                                    icon={<SaveOutlined />}
                                                >
                                                    保存资料
                                                </Button>
                                            </div>
                                        </Form>
                                    </div>
                                </TabPane>

                                <TabPane tab="安全设置" key="security">
                                    <div className={styles.formSection}>
                                        <Form
                                            form={securityForm}
                                            layout="vertical"
                                            onFinish={handleSecuritySave}
                                        >
                                            <div className="form-group">
                                                <div className="group-title">
                                                    <PhoneOutlined />
                                                    手机号码
                                                </div>
                                                <Form.Item
                                                    label="当前手机号"
                                                >
                                                    <Input 
                                                        value={user?.web?.phone ? `18888888888` : '未绑定'} 
                                                        disabled 
                                                        prefix={<PhoneOutlined />}
                                                    />
                                                </Form.Item>
                                                <Form.Item
                                                    label="新手机号"
                                                    name="newPhone"
                                                    rules={[
                                                        { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号' }
                                                    ]}
                                                >
                                                    <Input 
                                                        prefix={<PhoneOutlined />} 
                                                        placeholder="请输入新的手机号码" 
                                                    />
                                                </Form.Item>
                                                <Form.Item
                                                    label="验证码"
                                                    name="verifyCode"
                                                >
                                                    <div className={styles.verifySection}>
                                                        <div className="verify-input">
                                                            <Input placeholder="请输入验证码" />
                                                            <Button 
                                                                className="verify-btn"
                                                                disabled={countdown > 0}
                                                                onClick={handleSendCode}
                                                            >
                                                                {countdown > 0 ? `${countdown}s后重发` : '获取验证码'}
                                                            </Button>
                                                        </div>
                                                    </div>
                                                </Form.Item>
                                            </div>

                                            <div className="form-group">
                                                <div className="group-title">
                                                    <LockOutlined />
                                                    修改密码
                                                </div>
                                                <Form.Item
                                                    label="当前密码"
                                                    name="currentPassword"
                                                >
                                                    <Input.Password 
                                                        prefix={<LockOutlined />} 
                                                        placeholder="请输入当前密码" 
                                                    />
                                                </Form.Item>
                                                <div className="form-grid">
                                                    <Form.Item
                                                        label="新密码"
                                                        name="newPassword"
                                                        rules={[
                                                            { min: 6, max: 20, message: '密码长度为6-20个字符' }
                                                        ]}
                                                    >
                                                        <Input.Password 
                                                            prefix={<LockOutlined />} 
                                                            placeholder="请输入新密码" 
                                                        />
                                                    </Form.Item>
                                                    <Form.Item
                                                        label="确认密码"
                                                        name="confirmPassword"
                                                        dependencies={['newPassword']}
                                                        rules={[
                                                            ({ getFieldValue }) => ({
                                                                validator(_, value) {
                                                                    if (!value || getFieldValue('newPassword') === value) {
                                                                        return Promise.resolve();
                                                                    }
                                                                    return Promise.reject(new Error('两次输入的密码不一致'));
                                                                },
                                                            }),
                                                        ]}
                                                    >
                                                        <Input.Password 
                                                            prefix={<LockOutlined />} 
                                                            placeholder="请再次输入新密码" 
                                                        />
                                                    </Form.Item>
                                                </div>
                                            </div>

                                            <div className={styles.actionButtons}>
                                                <Button className="secondary-btn">
                                                    重置
                                                </Button>
                                                <Button 
                                                    type="primary" 
                                                    htmlType="submit"
                                                    loading={loading}
                                                    className="primary-btn"
                                                    icon={<SafetyOutlined />}
                                                >
                                                    保存设置
                                                </Button>
                                            </div>
                                        </Form>
                                    </div>
                                </TabPane>

                                <TabPane tab="通知设置" key="notifications">
                                    <div className={styles.formSection}>
                                        <Form
                                            form={notificationForm}
                                            layout="vertical"
                                            onFinish={handleNotificationSave}
                                            initialValues={{
                                                emailNotifications: true,
                                                pushNotifications: true,
                                                smsNotifications: false,
                                                marketingEmails: false,
                                                securityAlerts: true
                                            }}
                                        >
                                            <div className="form-group">
                                                <div className="group-title">
                                                    <BellOutlined />
                                                    通知偏好
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">邮件通知</div>
                                                        <div className="setting-desc">接收重要更新和消息提醒</div>
                                                    </div>
                                                    <Form.Item name="emailNotifications" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">推送通知</div>
                                                        <div className="setting-desc">浏览器推送通知</div>
                                                    </div>
                                                    <Form.Item name="pushNotifications" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">短信通知</div>
                                                        <div className="setting-desc">重要安全提醒短信</div>
                                                    </div>
                                                    <Form.Item name="smsNotifications" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">营销邮件</div>
                                                        <div className="setting-desc">产品更新和优惠信息</div>
                                                    </div>
                                                    <Form.Item name="marketingEmails" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">安全警报</div>
                                                        <div className="setting-desc">账户安全相关通知</div>
                                                    </div>
                                                    <Form.Item name="securityAlerts" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch />
                                                    </Form.Item>
                                                </div>
                                            </div>

                                            <div className="form-group">
                                                <div className="group-title">
                                                    <EyeOutlined />
                                                    隐私设置
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">个人资料可见性</div>
                                                        <div className="setting-desc">控制他人查看您的个人信息</div>
                                                    </div>
                                                    <Form.Item name="profileVisibility" style={{ margin: 0 }}>
                                                        <Select defaultValue="public" style={{ width: 120 }}>
                                                            <Option value="public">公开</Option>
                                                            <Option value="friends">好友</Option>
                                                            <Option value="private">私密</Option>
                                                        </Select>
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">在线状态</div>
                                                        <div className="setting-desc">显示您的在线状态</div>
                                                    </div>
                                                    <Form.Item name="showOnlineStatus" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch defaultChecked />
                                                    </Form.Item>
                                                </div>
                                                <div className={styles.settingItem}>
                                                    <div className="setting-info">
                                                        <div className="setting-title">活动记录</div>
                                                        <div className="setting-desc">记录您的活动历史</div>
                                                    </div>
                                                    <Form.Item name="activityTracking" valuePropName="checked" style={{ margin: 0 }}>
                                                        <Switch defaultChecked />
                                                    </Form.Item>
                                                </div>
                                            </div>

                                            <div className={styles.actionButtons}>
                                                <Button className="secondary-btn">
                                                    重置
                                                </Button>
                                                <Button 
                                                    type="primary" 
                                                    htmlType="submit"
                                                    loading={loading}
                                                    className="primary-btn"
                                                    icon={<SaveOutlined />}
                                                >
                                                    保存设置
                                                </Button>
                                            </div>
                                        </Form>
                                    </div>
                                </TabPane>
                            </Tabs>
                        </div>
                    </div>
                </div>
            </Body>
        </div>
    );
};
