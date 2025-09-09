import React, {useState} from "react";
import {AceEditor, Body} from "@/components";
import {
    Button, Card, Space, Typography, Alert, App, Row, Col, Switch,
    InputNumber, Select, Slider, Divider, Tag, notification
} from "antd";

const {Text, Title} = Typography;
const {Option} = Select;

export default () => {
    const {message} = App.useApp();

    // AceEditor 基础状态
    const [aceCode, setAceCode] = useState(`// AceEditor 完整功能演示
/**
 * 这是一个展示 AceEditor 所有功能的示例
 * 支持语法高亮、智能补全、主题切换等功能
 */
class Calculator {
    constructor() {
        this.history = [];
    }

    // 加法运算
    add(a, b) {
        const result = a + b;
        this.history.push(\`\${a} + \${b} = \${result}\`);
        return result;
    }

    // 获取计算历史
    getHistory() {
        return this.history;
    }
}

// 使用示例
const calc = new Calculator();
console.log(calc.add(10, 20)); // 30
console.log(calc.getHistory());

// 异步函数示例
async function fetchData(url) {
    try {
        const response = await fetch(url);
        return await response.json();
    } catch (error) {
        console.error('获取数据失败:', error);
        throw error;
    }
}

// ES6+ 特性展示
const users = [
    { id: 1, name: 'Alice', age: 25 },
    { id: 2, name: 'Bob', age: 30 },
    { id: 3, name: 'Charlie', age: 35 }
];

// 使用解构和箭头函数
const getAdultUsers = (users) => users.filter(({ age }) => age >= 18);
const userNames = users.map(({ name }) => name);

console.log('成年用户:', getAdultUsers(users));
console.log('用户名列表:', userNames);`);

    // 编辑器配置状态
    const [aceMode, setAceMode] = useState('javascript');
    const [aceTheme, setAceTheme] = useState('monokai');
    const [aceFontSize, setAceFontSize] = useState(14);
    const [aceTabSize, setAceTabSize] = useState(4);
    const [aceReadOnly, setAceReadOnly] = useState(false);
    const [showPrintMargin, setShowPrintMargin] = useState(true);
    const [showGutter, setShowGutter] = useState(true);
    const [highlightActiveLine, setHighlightActiveLine] = useState(true);
    const [highlightSelectedWord, setHighlightSelectedWord] = useState(true);
    const [wrapEnabled, setWrapEnabled] = useState(false);
    const [autoScrollEditorIntoView, setAutoScrollEditorIntoView] = useState(false);
    const [showLineNumbers, setShowLineNumbers] = useState(true);
    const [enableBasicAutocompletion, setEnableBasicAutocompletion] = useState(true);
    const [enableLiveAutocompletion, setEnableLiveAutocompletion] = useState(true);
    const [enableSnippets, setEnableSnippets] = useState(true);
    const [editorWidth, setEditorWidth] = useState('100%');
    const [editorHeight, setEditorHeight] = useState(600);
    const [maxLines, setMaxLines] = useState<number | undefined>(undefined);
    const [minLines, setMinLines] = useState<number | undefined>(undefined);

    // 可用的模式和主题
    const modes = [
        'text', 'javascript', 'typescript', 'python', 'java', 'json', 'html', 'css', 'scss', 'less',
        'xml', 'yaml', 'markdown', 'sql', 'php', 'golang', 'rust', 'c_cpp', 'csharp', 'ruby'
    ];

    const themes = [
        'monokai', 'github', 'tomorrow', 'kuroir', 'twilight', 'xcode', 'textmate', 'solarized_dark',
        'solarized_light', 'terminal', 'chrome', 'eclipse', 'dreamweaver', 'cobalt'
    ];

    // 代码示例
    const codeExamples: Record<string, string> = {
        javascript: `// JavaScript ES6+ 示例
class DataProcessor {
    constructor(data) {
        this.data = data;
        this.processed = false;
    }

    async process() {
        if (this.processed) return this.data;
        
        // 模拟异步处理
        await new Promise(resolve => setTimeout(resolve, 100));
        
        this.data = this.data
            .filter(item => item.active)
            .map(item => ({
                ...item,
                processedAt: new Date().toISOString()
            }));
        
        this.processed = true;
        return this.data;
    }
}

// 使用示例
const processor = new DataProcessor([
    { id: 1, name: 'Item 1', active: true },
    { id: 2, name: 'Item 2', active: false },
    { id: 3, name: 'Item 3', active: true }
]);

processor.process().then(result => {
    console.log('处理结果:', result);
});`,

        python: `# Python 示例代码
import asyncio
from typing import List, Dict, Optional
from dataclasses import dataclass
from datetime import datetime

@dataclass
class User:
    id: int
    name: str
    email: str
    created_at: datetime = None
    
    def __post_init__(self):
        if self.created_at is None:
            self.created_at = datetime.now()

class UserService:
    def __init__(self):
        self.users: List[User] = []
    
    def add_user(self, name: str, email: str) -> User:
        user = User(
            id=len(self.users) + 1,
            name=name,
            email=email
        )
        self.users.append(user)
        return user
    
    def find_user(self, user_id: int) -> Optional[User]:
        return next((u for u in self.users if u.id == user_id), None)
    
    async def get_user_stats(self) -> Dict[str, int]:
        # 模拟异步操作
        await asyncio.sleep(0.1)
        return {
            'total_users': len(self.users),
            'active_users': len([u for u in self.users if u.email])
        }

# 使用示例
service = UserService()
user1 = service.add_user("Alice", "alice@example.com")
user2 = service.add_user("Bob", "bob@example.com")

print(f"创建用户: {user1.name} ({user1.email})")
print(f"用户总数: {len(service.users)}")`,

        json: `{
  "name": "ace-editor-demo",
  "version": "2.0.0",
  "description": "AceEditor 完整功能演示项目",
  "main": "index.js",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "preview": "vite preview",
    "test": "jest",
    "test:watch": "jest --watch",
    "lint": "eslint src --ext .ts,.tsx",
    "lint:fix": "eslint src --ext .ts,.tsx --fix"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "antd": "^5.12.0",
    "ace-builds": "^1.32.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "@types/react-dom": "^18.2.0",
    "@vitejs/plugin-react": "^4.2.0",
    "typescript": "^5.3.0",
    "vite": "^5.0.0",
    "eslint": "^8.55.0",
    "jest": "^29.7.0"
  },
  "keywords": [
    "react",
    "ace-editor",
    "typescript",
    "code-editor",
    "syntax-highlighting"
  ],
  "author": "AceEditor Demo Team",
  "license": "MIT",
  "repository": {
    "type": "git",
    "url": "https://github.com/example/ace-editor-demo"
  }
}`,

        html: `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AceEditor 演示页面</title>
    <link rel="stylesheet" href="styles.css">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
            padding: 24px;
        }
        .header {
            text-align: center;
            margin-bottom: 32px;
        }
        .editor-demo {
            border: 1px solid #d9d9d9;
            border-radius: 6px;
            overflow: hidden;
        }
    </style>
</head>
<body>
    <div class="container">
        <header class="header">
            <h1>🚀 AceEditor 功能演示</h1>
            <p>一个功能强大的代码编辑器组件</p>
        </header>
        
        <main>
            <section class="editor-demo">
                <div id="ace-editor"></div>
            </section>
            
            <section class="features">
                <h2>✨ 主要特性</h2>
                <ul>
                    <li>支持多种编程语言语法高亮</li>
                    <li>智能代码补全和片段</li>
                    <li>多种编辑器主题</li>
                    <li>自定义快捷键和命令</li>
                    <li>代码折叠和搜索替换</li>
                    <li>实时语法检查和错误提示</li>
                </ul>
            </section>
        </main>
    </div>
    
    <script src="ace-editor.js"></script>
    <script>
        // 初始化编辑器
        const editor = ace.edit("ace-editor");
        editor.setTheme("ace/theme/monokai");
        editor.session.setMode("ace/mode/javascript");
    </script>
</body>
</html>`
    };

    // 错误处理函数
    const handleError = (error: Error) => {
        notification.error({
            message: 'AceEditor 错误',
            description: error.message || '编辑器发生未知错误',
            duration: 5,
        });
        console.error('AceEditor Error:', error);
    };

    return (
        <Body>
            <Space direction="vertical" size="large" style={{width: '100%'}}>
                {/* 标题和介绍 */}
                <Card>
                    <div style={{textAlign: 'center'}}>
                        <Title level={2}>🚀 AceEditor 完整功能演示</Title>
                        <Text type="secondary">
                            展示 AceEditor 组件的所有配置选项和功能特性，包括语法高亮、主题切换、智能补全等
                        </Text>
                    </div>
                </Card>

                {/* 配置面板 */}
                <Card title="⚙️ 编辑器配置" size="small">
                    <Row gutter={[16, 16]}>
                        {/* 基础配置 */}
                        <Col span={24}>
                            <Title level={5}>基础配置</Title>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>编程语言:</Text>
                                <Select
                                    value={aceMode}
                                    onChange={setAceMode}
                                    style={{width: '100%'}}
                                    showSearch
                                >
                                    {modes.map(mode => (
                                        <Option key={mode} value={mode}>
                                            {mode.toUpperCase()}
                                        </Option>
                                    ))}
                                </Select>
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>编辑器主题:</Text>
                                <Select
                                    value={aceTheme}
                                    onChange={setAceTheme}
                                    style={{width: '100%'}}
                                    showSearch
                                >
                                    {themes.map(theme => (
                                        <Option key={theme} value={theme}>
                                            {theme.replace('_', ' ').toUpperCase()}
                                        </Option>
                                    ))}
                                </Select>
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>字体大小: {aceFontSize}px</Text>
                                <Slider
                                    min={10}
                                    max={24}
                                    value={aceFontSize}
                                    onChange={setAceFontSize}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>Tab 大小: {aceTabSize}</Text>
                                <Slider
                                    min={2}
                                    max={8}
                                    value={aceTabSize}
                                    onChange={setAceTabSize}
                                />
                            </Space>
                        </Col>

                        {/* 尺寸配置 */}
                        <Col span={24}>
                            <Divider/>
                            <Title level={5}>尺寸配置</Title>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>编辑器高度: {editorHeight}px</Text>
                                <Slider
                                    min={10}
                                    max={800}
                                    value={editorHeight}
                                    onChange={setEditorHeight}
                                />
                            </Space>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>最大行数:</Text>
                                <InputNumber
                                    min={10}
                                    max={100}
                                    value={maxLines}
                                    onChange={(value) => setMaxLines(value || undefined)}
                                    placeholder="不限制"
                                    style={{width: '100%'}}
                                />
                            </Space>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical" style={{width: '100%'}}>
                                <Text strong>最小行数:</Text>
                                <InputNumber
                                    min={1}
                                    max={20}
                                    value={minLines}
                                    onChange={(value) => setMinLines(value || undefined)}
                                    placeholder="不限制"
                                    style={{width: '100%'}}
                                />
                            </Space>
                        </Col>

                        {/* 功能开关 */}
                        <Col span={24}>
                            <Divider/>
                            <Title level={5}>功能开关</Title>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>只读模式:</Text>
                                <Switch
                                    checked={aceReadOnly}
                                    onChange={setAceReadOnly}
                                    checkedChildren="只读"
                                    unCheckedChildren="编辑"
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>显示打印边距:</Text>
                                <Switch
                                    checked={showPrintMargin}
                                    onChange={setShowPrintMargin}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>显示行号栏:</Text>
                                <Switch
                                    checked={showGutter}
                                    onChange={setShowGutter}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>显示行号:</Text>
                                <Switch
                                    checked={showLineNumbers}
                                    onChange={setShowLineNumbers}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>高亮当前行:</Text>
                                <Switch
                                    checked={highlightActiveLine}
                                    onChange={setHighlightActiveLine}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>高亮选中词:</Text>
                                <Switch
                                    checked={highlightSelectedWord}
                                    onChange={setHighlightSelectedWord}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>自动换行:</Text>
                                <Switch
                                    checked={wrapEnabled}
                                    onChange={setWrapEnabled}
                                />
                            </Space>
                        </Col>
                        <Col span={6}>
                            <Space direction="vertical">
                                <Text strong>自动滚动:</Text>
                                <Switch
                                    checked={autoScrollEditorIntoView}
                                    onChange={setAutoScrollEditorIntoView}
                                />
                            </Space>
                        </Col>

                        {/* 智能补全 */}
                        <Col span={24}>
                            <Divider/>
                            <Title level={5}>智能补全</Title>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical">
                                <Text strong>基础补全:</Text>
                                <Switch
                                    checked={enableBasicAutocompletion}
                                    onChange={setEnableBasicAutocompletion}
                                />
                            </Space>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical">
                                <Text strong>实时补全:</Text>
                                <Switch
                                    checked={enableLiveAutocompletion}
                                    onChange={setEnableLiveAutocompletion}
                                />
                            </Space>
                        </Col>
                        <Col span={8}>
                            <Space direction="vertical">
                                <Text strong>代码片段:</Text>
                                <Switch
                                    checked={enableSnippets}
                                    onChange={setEnableSnippets}
                                />
                            </Space>
                        </Col>

                        {/* 快速切换示例 */}
                        <Col span={24}>
                            <Divider/>
                            <Title level={5}>快速切换示例</Title>
                        </Col>
                        <Col span={24}>
                            <Space wrap>
                                {Object.keys(codeExamples).map(lang => (
                                    <Button
                                        key={lang}
                                        type={aceMode === lang ? 'primary' : 'default'}
                                        onClick={() => {
                                            setAceMode(lang);
                                            setTimeout(() => {
                                                setAceCode(codeExamples[lang as keyof typeof codeExamples]);
                                            }, 50);
                                        }}
                                    >
                                        {lang.toUpperCase()} 示例
                                    </Button>
                                ))}
                            </Space>
                        </Col>
                    </Row>
                </Card>

                {/* AceEditor 组件 */}
                <Card title="📝 代码编辑器" size="small">
                    <AceEditor
                        value={aceCode}
                        defaultValue=""
                        mode={aceMode}
                        theme={aceTheme}
                        width={editorWidth}
                        height={editorHeight}
                        fontSize={aceFontSize}
                        tabSize={aceTabSize}
                        readOnly={aceReadOnly}
                        showPrintMargin={showPrintMargin}
                        showGutter={showGutter}
                        highlightActiveLine={highlightActiveLine}
                        highlightSelectedWord={highlightSelectedWord}
                        wrapEnabled={wrapEnabled}
                        autoScrollEditorIntoView={autoScrollEditorIntoView}
                        maxLines={maxLines}
                        minLines={minLines}
                        placeholder="请输入代码..."
                        showLineNumbers={showLineNumbers}
                        enableBasicAutocompletion={enableBasicAutocompletion}
                        enableLiveAutocompletion={enableLiveAutocompletion}
                        enableSnippets={enableSnippets}
                        acePath="/ace/"
                        loadingContent={
                            <div style={{
                                width: typeof editorWidth === 'number' ? `${editorWidth}px` : editorWidth,
                                height: `${editorHeight}px`,
                                display: 'flex',
                                alignItems: 'center',
                                justifyContent: 'center',
                                border: '1px solid #d9d9d9',
                                borderRadius: '6px',
                                backgroundColor: '#fafafa'
                            }}>
                                <div style={{textAlign: 'center'}}>
                                    <div style={{fontSize: '32px', marginBottom: '16px'}}>⚡</div>
                                    <div style={{fontSize: '18px', fontWeight: 'bold', marginBottom: '8px'}}>
                                        正在加载 AceEditor
                                    </div>
                                    <div style={{fontSize: '14px', color: '#666'}}>
                                        动态加载编辑器资源中...
                                    </div>
                                </div>
                            </div>
                        }
                        onChange={(value) => {
                            setAceCode(value);
                        }}
                        onSelectionChange={(selection) => {
                            // console.log('Selection changed:', selection);
                        }}
                        onCursorChange={(selection) => {
                            // console.log('Cursor changed:', selection);
                        }}
                        onBlur={() => {
                            // console.log('Editor blurred');
                        }}
                        onFocus={() => {
                            // console.log('Editor focused');
                        }}
                        onLoad={(editor) => {
                            message.success('🎉 AceEditor 加载成功！');
                            //console.log('Editor loaded:', editor);
                        }}
                        onBeforeLoad={(ace) => {
                            //console.log('Ace loaded:', ace);
                        }}
                        onError={handleError}
                        commands={[
                            {
                                name: 'saveCode',
                                bindKey: {win: 'Ctrl-S', mac: 'Command-S'},
                                exec: (editor) => {
                                    const code = editor.getValue();
                                    message.success(`代码已保存（${code.split('\n').length} 行）`);
                                }
                            },
                            {
                                name: 'formatCode',
                                bindKey: {win: 'Ctrl-Shift-F', mac: 'Command-Shift-F'},
                                exec: (editor) => {
                                    message.info('代码格式化功能（演示）');
                                }
                            },
                            {
                                name: 'runCode',
                                bindKey: {win: 'Ctrl-R', mac: 'Command-R'},
                                exec: (editor) => {
                                    message.info('运行代码功能（演示）');
                                }
                            }
                        ]}
                        annotations={[
                            {
                                row: 2,
                                column: 0,
                                text: '这是一个类定义，展示了现代 JavaScript 语法',
                                type: 'info'
                            },
                            {
                                row: 10,
                                column: 4,
                                text: '建议添加参数类型检查',
                                type: 'warning'
                            }
                        ]}
                        markers={[]}
                        className="custom-ace-editor"
                        style={{
                            borderRadius: '8px'
                        }}
                    />
                </Card>

                {/* 功能说明 */}
                <Card title="📖 功能说明" size="small">
                    <Row gutter={[16, 16]}>
                        <Col span={12}>
                            <Alert
                                message="快捷键"
                                description={
                                    <Space direction="vertical">
                                        <Text><code>Ctrl+S / Cmd+S</code>: 保存代码</Text>
                                        <Text><code>Ctrl+Shift+F / Cmd+Shift+F</code>: 格式化代码</Text>
                                        <Text><code>Ctrl+R / Cmd+R</code>: 运行代码</Text>
                                        <Text><code>Ctrl+F / Cmd+F</code>: 查找</Text>
                                        <Text><code>Ctrl+H / Cmd+H</code>: 替换</Text>
                                        <Text><code>Ctrl+/ / Cmd+/</code>: 切换注释</Text>
                                    </Space>
                                }
                                type="info"
                                showIcon
                            />
                        </Col>
                        <Col span={12}>
                            <Alert
                                message="主要特性"
                                description={
                                    <Space direction="vertical">
                                        <Text>✨ 支持 20+ 种编程语言</Text>
                                        <Text>🎨 14+ 种编辑器主题</Text>
                                        <Text>🔧 完整的配置选项</Text>
                                        <Text>⚡ 动态资源加载</Text>
                                        <Text>🤖 智能代码补全</Text>
                                        <Text>🔍 语法检查和错误提示</Text>
                                    </Space>
                                }
                                type="success"
                                showIcon
                            />
                        </Col>
                    </Row>
                </Card>

                {/* 当前配置 */}
                <Card title="🔧 当前配置" size="small">
                    <Row gutter={[16, 8]}>
                        <Col span={6}><Tag color="blue">语言: {aceMode}</Tag></Col>
                        <Col span={6}><Tag color="green">主题: {aceTheme}</Tag></Col>
                        <Col span={6}><Tag color="orange">字体: {aceFontSize}px</Tag></Col>
                        <Col span={6}><Tag color="purple">Tab: {aceTabSize}</Tag></Col>
                        <Col span={6}><Tag color="red">只读: {aceReadOnly ? '是' : '否'}</Tag></Col>
                        <Col span={6}><Tag color="cyan">行号: {showLineNumbers ? '显示' : '隐藏'}</Tag></Col>
                        <Col span={6}><Tag color="magenta">换行: {wrapEnabled ? '开启' : '关闭'}</Tag></Col>
                        <Col span={6}><Tag color="gold">补全: {enableLiveAutocompletion ? '开启' : '关闭'}</Tag></Col>
                    </Row>
                </Card>
            </Space>
        </Body>
    );
}
