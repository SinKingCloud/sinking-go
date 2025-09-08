import React, {useRef, useEffect, useState, useCallback, useMemo} from 'react';
import Script, {preloadScript, isScriptLoaded} from '@/components/script';

// Ace Editor 配置接口
export interface AceEditorProps {
    value?: string;
    defaultValue?: string;
    mode?: string;
    theme?: string;
    width?: string | number;
    height?: string | number;
    fontSize?: number;
    tabSize?: number;
    readOnly?: boolean;
    showPrintMargin?: boolean;
    showGutter?: boolean;
    highlightActiveLine?: boolean;
    highlightSelectedWord?: boolean;
    wrapEnabled?: boolean;
    autoScrollEditorIntoView?: boolean;
    maxLines?: number;
    minLines?: number;
    placeholder?: string;
    className?: string;
    style?: React.CSSProperties;
    onChange?: (value: string, event?: any) => void;
    onSelectionChange?: (selection: any, event?: any) => void;
    onCursorChange?: (selection: any, event?: any) => void;
    onBlur?: (event?: any) => void;
    onFocus?: (event?: any) => void;
    onLoad?: (editor: any) => void;
    onBeforeLoad?: (ace: any) => void;
    commands?: Array<{
        name: string;
        bindKey: { win: string; mac: string };
        exec: (editor: any) => void;
    }>;
    annotations?: Array<{
        row: number;
        column: number;
        text: string;
        type: 'error' | 'warning' | 'info';
    }>;
    markers?: Array<{
        startRow: number;
        startCol: number;
        endRow: number;
        endCol: number;
        className: string;
        type: string;
    }>;
    enableBasicAutocompletion?: boolean;
    enableLiveAutocompletion?: boolean;
    enableSnippets?: boolean;
    showLineNumbers?: boolean;
    acePath?: string; // 自定义 Ace 资源路径
    loadingContent?: React.ReactNode; // 自定义加载内容
    onError?: (error: Error) => void; // 错误处理回调
}


const AceEditor: React.FC<AceEditorProps> = ({
                                                 value = '',
                                                 defaultValue = '',
                                                 mode = 'text',
                                                 theme = 'monokai',
                                                 width = '100%',
                                                 height = 400,
                                                 fontSize = 14,
                                                 tabSize = 4,
                                                 readOnly = false,
                                                 showPrintMargin = true,
                                                 showGutter = true,
                                                 highlightActiveLine = true,
                                                 highlightSelectedWord = true,
                                                 wrapEnabled = false,
                                                 autoScrollEditorIntoView = false,
                                                 maxLines,
                                                 minLines,
                                                 placeholder = '',
                                                 className = '',
                                                 style = {},
                                                 onChange,
                                                 onSelectionChange,
                                                 onCursorChange,
                                                 onBlur,
                                                 onFocus,
                                                 onLoad,
                                                 onBeforeLoad,
                                                 commands = [],
                                                 annotations = [],
                                                 markers = [],
                                                 enableBasicAutocompletion = true,
                                                 enableLiveAutocompletion = true,
                                                 enableSnippets = true,
                                                 showLineNumbers = true,
                                                 acePath = '/ace',
                                                 loadingContent,
                                                 onError
                                             }) => {
    const containerRef = useRef<HTMLDivElement>(null);
    const editorRef = useRef<any>(null);
    const [aceLoaded, setAceLoaded] = useState(false);
    const [editorValue, setEditorValue] = useState(value || defaultValue);

    // 只加载核心脚本，扩展按需加载
    const coreScript = useMemo(() => `${acePath}/src-min/ace.js`, [acePath]);

    // 动态加载模式文件
    const loadMode = useCallback(async (modeName: string) => {
        if (!modeName || modeName === 'text') return;

        const scriptUrl = `${acePath}/src-min/mode-${modeName}.js`;

        if (isScriptLoaded(scriptUrl)) return;

        try {
            await preloadScript(scriptUrl, {
                cache: true,
                timeout: 5000,
                retryCount: 2
            });
        } catch (error) {
            onError?.(error as Error);
        }
    }, [acePath]);

    // 动态加载主题文件
    const loadTheme = useCallback(async (themeName: string) => {
        if (!themeName) return;

        const scriptUrl = `${acePath}/src-min/theme-${themeName}.js`;

        if (isScriptLoaded(scriptUrl)) return;

        try {
            await preloadScript(scriptUrl, {
                cache: true,
                timeout: 5000,
                retryCount: 2
            });
        } catch (error) {
            onError?.(error as Error);
        }
    }, [acePath]);


    // 当 Ace 加载完成后初始化编辑器
    useEffect(() => {
        if (!aceLoaded || editorRef.current) return;

        const initialize = async () => {
            if (!containerRef.current || !window.ace) return;

            // 预先加载初始的模式和主题
            const modeUrl = `${acePath}/src-min/mode-${mode}.js`;
            const themeUrl = `${acePath}/src-min/theme-${theme}.js`;

            try {
                // 并行加载所有需要的资源
                const loadTasks = [];

                // 加载模式文件
                if (mode !== 'text' && !isScriptLoaded(modeUrl)) {
                    loadTasks.push(preloadScript(modeUrl, {cache: true}));
                }

                // 加载主题文件
                if (!isScriptLoaded(themeUrl)) {
                    loadTasks.push(preloadScript(themeUrl, {cache: true}));
                }

                // 加载扩展文件
                if (enableBasicAutocompletion || enableLiveAutocompletion) {
                    const langToolsUrl = `${acePath}/src-min/ext-language_tools.js`;
                    if (!isScriptLoaded(langToolsUrl)) {
                        loadTasks.push(preloadScript(langToolsUrl, {cache: true}));
                    }
                }

                if (enableSnippets) {
                    const searchboxUrl = `${acePath}/src-min/ext-searchbox.js`;
                    if (!isScriptLoaded(searchboxUrl)) {
                        loadTasks.push(preloadScript(searchboxUrl, {cache: true}));
                    }
                }

                if (loadTasks.length > 0) {
                    await Promise.all(loadTasks);
                }

                // 直接初始化编辑器
                if (!editorRef.current && containerRef.current && window.ace) {
                    onBeforeLoad?.(window.ace);

                    const editor = window.ace.edit(containerRef.current);
                    editorRef.current = editor;

                    // 基础配置
                    editor.setTheme(`ace/theme/${theme}`);
                    editor.session.setMode(`ace/mode/${mode}`);
                    editor.setValue(editorValue, -1);
                    editor.setFontSize(fontSize);
                    editor.session.setTabSize(tabSize);
                    editor.setReadOnly(readOnly);
                    editor.setShowPrintMargin(showPrintMargin);
                    editor.renderer.setShowGutter(showGutter && showLineNumbers);
                    editor.setHighlightActiveLine(highlightActiveLine);
                    editor.setHighlightSelectedWord(highlightSelectedWord);
                    editor.session.setUseWrapMode(wrapEnabled);
                    editor.setAutoScrollEditorIntoView(autoScrollEditorIntoView);

                    if (maxLines) editor.setOption('maxLines', maxLines);
                    if (minLines) editor.setOption('minLines', minLines);
                    if (placeholder) editor.setOption('placeholder', placeholder);

                    // 自动完成配置
                    if (enableBasicAutocompletion || enableLiveAutocompletion || enableSnippets) {
                        window.ace.require('ace/ext/language_tools');
                        editor.setOptions({
                            enableBasicAutocompletion,
                            enableLiveAutocompletion,
                            enableSnippets
                        });
                    }

                    // 事件监听
                    editor.on('change', () => {
                        const newValue = editor.getValue();
                        setEditorValue(newValue);
                        onChange?.(newValue);
                    });

                    if (onSelectionChange) {
                        editor.on('changeSelection', (e: any) => {
                            onSelectionChange(editor.getSelection(), e);
                        });
                    }

                    if (onCursorChange) {
                        editor.on('changeCursor', (e: any) => {
                            onCursorChange(editor.getSelection(), e);
                        });
                    }

                    if (onBlur) editor.on('blur', onBlur);
                    if (onFocus) editor.on('focus', onFocus);

                    // 添加自定义命令
                    commands.forEach(command => {
                        editor.commands.addCommand(command);
                    });

                    // 设置注释和标记
                    if (annotations.length > 0) {
                        editor.session.setAnnotations(annotations);
                    }

                    if (markers.length > 0) {
                        const Range = window.ace.require('ace/range').Range;
                        markers.forEach(marker => {
                            const range = new Range(marker.startRow, marker.startCol, marker.endRow, marker.endCol);
                            editor.session.addMarker(range, marker.className, marker.type);
                        });
                    }

                    onLoad?.(editor);
                }

            } catch (error) {
                onError?.(error as Error);
            }
        };

        initialize();
    }, [aceLoaded]); // 只依赖 aceLoaded，避免无限循环

    // 更新编辑器配置
    useEffect(() => {
        if (!editorRef.current) return;

        const editor = editorRef.current;

        // 更新值（避免无限循环）
        if (value !== undefined && value !== editor.getValue()) {
            const cursorPosition = editor.getCursorPosition();
            editor.setValue(value, -1);
            editor.moveCursorToPosition(cursorPosition);
            setEditorValue(value);
        }

        // 更新其他配置
        editor.setFontSize(fontSize);
        editor.session.setTabSize(tabSize);
        editor.setReadOnly(readOnly);
        editor.setShowPrintMargin(showPrintMargin);
        editor.renderer.setShowGutter(showGutter && showLineNumbers);
        editor.setHighlightActiveLine(highlightActiveLine);
        editor.setHighlightSelectedWord(highlightSelectedWord);
        editor.session.setUseWrapMode(wrapEnabled);
        editor.setAutoScrollEditorIntoView(autoScrollEditorIntoView);

        if (maxLines) editor.setOption('maxLines', maxLines);
        if (minLines) editor.setOption('minLines', minLines);
        if (placeholder) editor.setOption('placeholder', placeholder);

    }, [
        value, fontSize, tabSize, readOnly, showPrintMargin, showGutter,
        highlightActiveLine, highlightSelectedWord, wrapEnabled,
        autoScrollEditorIntoView, maxLines, minLines, placeholder, showLineNumbers
    ]);

    // 处理模式变化
    useEffect(() => {
        if (editorRef.current && mode) {
            const updateMode = async () => {
                await loadMode(mode);
                if (editorRef.current) {
                    editorRef.current.session.setMode(`ace/mode/${mode}`);
                }
            };
            updateMode();
        }
    }, [mode]);

    // 处理主题变化
    useEffect(() => {
        if (editorRef.current && theme) {
            const updateTheme = async () => {
                await loadTheme(theme);
                if (editorRef.current) {
                    editorRef.current.setTheme(`ace/theme/${theme}`);
                }
            };
            updateTheme();
        }
    }, [theme]);

    // 组件卸载时清理
    useEffect(() => {
        return () => {
            if (editorRef.current) {
                editorRef.current.destroy();
                editorRef.current = null;
            }
        };
    }, []);

    return (
        <div className={className} style={style}>
            <Script
                src={coreScript}
                parallel={false}
                timeout={10000}
                retryCount={2}
                cache={true}
                onLoad={() => setAceLoaded(true)}
                loading={loadingContent}
            >
                <div
                    ref={containerRef}
                    style={{
                        width: typeof width === 'number' ? `${width}px` : width,
                        height: typeof height === 'number' ? `${height}px` : height,
                        
                        overflow: 'hidden'
                    }}
                />
            </Script>
        </div>
    );
};

export default AceEditor;

// 声明全局 ace 对象
declare global {
    interface Window {
        ace: any;
    }
}
