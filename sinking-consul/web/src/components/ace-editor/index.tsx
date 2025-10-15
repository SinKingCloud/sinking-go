// noinspection JSIncompatibleTypesComparison

import React, {useRef, useEffect, useState, useCallback, useMemo, useLayoutEffect} from 'react';
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
    loadingContent?: React.ReactNode; // 仅首次加载时的自定义加载内容
    containerStyle?: React.CSSProperties; // 内部稳定容器样式
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
                                                 containerStyle,
                                                 loadingContent,
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
                                                 onError
                                             }) => {
    const containerRef = useRef<HTMLDivElement | null>(null);
    const editorRef = useRef<any>(null);
    const [aceLoaded, setAceLoaded] = useState(false);
    const [hasScriptLoaded, setHasScriptLoaded] = useState(false); // 仅首次加载展示 loading
    const [editorValue, setEditorValue] = useState(value || defaultValue);

    // 创建一个稳定的容器元素，避免 React 重新创建
    const stableContainer = useRef<HTMLDivElement | null>(null);
    const mountedRef = useRef(false);
    const [containerMounted, setContainerMounted] = useState(false);
    const outerContainerRef = useRef<any>(null);
    const commandNamesRef = useRef<string[]>([]);
    const markerIdsRef = useRef<number[]>([]);

    // 统一应用编辑器可变配置
    const applyEditorOptions = useCallback((ed: any, options: {
        fontSize: number;
        tabSize: number;
        readOnly: boolean;
        showPrintMargin: boolean;
        showGutter: boolean;
        showLineNumbers: boolean;
        highlightActiveLine: boolean;
        highlightSelectedWord: boolean;
        wrapEnabled: boolean;
        autoScrollEditorIntoView: boolean;
        maxLines?: number;
        minLines?: number;
        placeholder?: string;
    }) => {
        if (!ed) return;
        const {
            fontSize,
            tabSize,
            readOnly,
            showPrintMargin,
            showGutter,
            showLineNumbers,
            highlightActiveLine,
            highlightSelectedWord,
            wrapEnabled,
            autoScrollEditorIntoView,
            maxLines,
            minLines,
            placeholder
        } = options;

        ed.setFontSize(fontSize);
        ed.session.setTabSize(tabSize);
        ed.setReadOnly(readOnly);
        ed.setShowPrintMargin(showPrintMargin);
        ed.renderer.setShowGutter(showGutter && showLineNumbers);
        ed.setHighlightActiveLine(highlightActiveLine);
        ed.setHighlightSelectedWord(highlightSelectedWord);
        ed.session.setUseWrapMode(wrapEnabled);
        ed.setAutoScrollEditorIntoView(autoScrollEditorIntoView);
        if (maxLines) ed.setOption('maxLines', maxLines);
        if (minLines) ed.setOption('minLines', minLines);
        if (placeholder) ed.setOption('placeholder', placeholder);
    }, []);

    // 创建稳定容器（只执行一次）
    useLayoutEffect(() => {
        // 创建稳定的容器元素，只创建一次
        if (!stableContainer.current) {
            stableContainer.current = document.createElement('div');
            stableContainer.current.style.overflow = 'hidden';
        }

        return () => {
            // 不主动移除稳定容器，返回页面时可复用，避免丢失
        };
    }, []);

    // 挂载稳定容器到 React 容器中
    useLayoutEffect(() => {
        if (containerRef.current && stableContainer.current && !mountedRef.current) {
            // 设置初始尺寸
            const initialWidth = typeof width === 'number' ? `${width}px` : width;
            const initialHeight = typeof height === 'number' ? `${height}px` : height;

            // 设置外层容器初始尺寸
            if (outerContainerRef.current) {
                outerContainerRef.current.style.width = initialWidth;
                outerContainerRef.current.style.height = initialHeight;
            }

            // 设置稳定容器初始尺寸
            stableContainer.current.style.width = initialWidth;
            stableContainer.current.style.height = initialHeight;

            containerRef.current.appendChild(stableContainer.current);
            mountedRef.current = true;
            setContainerMounted(true);
        } else if (containerRef.current && stableContainer.current && !stableContainer.current.parentNode) {
            // 跨页面返回时，如果稳定容器没有挂载，重新挂载
            containerRef.current.appendChild(stableContainer.current);
            mountedRef.current = true;
            setContainerMounted(true);
        }
    });

    // 只加载核心脚本，扩展按需加载
    const coreScript = useMemo(() => `${acePath}/src-min/ace.js`, [acePath]);
    const isAcePresent = useMemo(() => typeof window !== 'undefined' && (window as any).ace, []);

    // 如果 ace 已经在全局存在（例如返回页面后），直接标记已加载
    useEffect(() => {
        if (isAcePresent && !aceLoaded) setAceLoaded(true);
    }, [isAcePresent, aceLoaded]);

    // React 容器样式 - 使用固定样式，通过稳定容器控制实际尺寸
    const containerBaseStyle = useMemo(() => ({
        width: '100%',
        height: '100%',
        overflow: 'hidden',
        position: 'relative' as const
    }), []);

    // 应用用户传入的 containerStyle 到稳定容器
    useLayoutEffect(() => {
        if (!stableContainer.current) return;
        const el = stableContainer.current as HTMLDivElement;
        el.style.overflow = 'hidden';
        if (containerStyle) {
            try {
                Object.assign(el.style, containerStyle);
            } catch {
            }
        }
    }, [containerStyle]);

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
            if (!stableContainer.current || !window.ace || !mountedRef.current) return;

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

                // 直接初始化编辑器，使用稳定容器
                if (!editorRef.current && stableContainer.current && window.ace) {
                    onBeforeLoad?.(window.ace);

                    const editor = window.ace.edit(stableContainer.current);
                    editorRef.current = editor;

                    // 基础配置
                    editor.setTheme(`ace/theme/${theme}`);
                    editor.session.setMode(`ace/mode/${mode}`);
                    editor.setValue(editorValue, -1);
                    applyEditorOptions(editor, {
                        fontSize,
                        tabSize,
                        readOnly,
                        showPrintMargin,
                        showGutter,
                        showLineNumbers,
                        highlightActiveLine,
                        highlightSelectedWord,
                        wrapEnabled,
                        autoScrollEditorIntoView,
                        maxLines,
                        minLines,
                        placeholder
                    });

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
                        const Range = window.ace.require('ace/range').Range as any;
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
    }, [aceLoaded, containerMounted]); // 依赖 aceLoaded 和容器挂载状态

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
        applyEditorOptions(editor, {
            fontSize,
            tabSize,
            readOnly,
            showPrintMargin,
            showGutter,
            showLineNumbers,
            highlightActiveLine,
            highlightSelectedWord,
            wrapEnabled,
            autoScrollEditorIntoView,
            maxLines,
            minLines,
            placeholder
        });

    }, [
        value, fontSize, tabSize, readOnly, showPrintMargin, showGutter,
        highlightActiveLine, highlightSelectedWord, wrapEnabled,
        autoScrollEditorIntoView, maxLines, minLines, placeholder, showLineNumbers
    ]);

    // 同步命令（commands）
    useEffect(() => {
        if (!editorRef.current) return;
        const editor = editorRef.current;

        // 移除旧命令
        if (commandNamesRef.current.length) {
            commandNamesRef.current.forEach((name) => {
                try {
                    editor.commands.removeCommand(name);
                } catch {
                }
            });
            commandNamesRef.current = [];
        }

        // 添加新命令
        if (Array.isArray(commands) && commands.length > 0) {
            commands.forEach((cmd) => {
                try {
                    editor.commands.addCommand(cmd);
                    commandNamesRef.current.push(cmd.name);
                } catch {
                }
            });
        }
    }, [commands]);

    // 同步注解（annotations）
    useEffect(() => {
        if (!editorRef.current) return;
        const editor = editorRef.current;
        try {
            editor.session.setAnnotations(Array.isArray(annotations) ? annotations : []);
        } catch {
        }
    }, [annotations]);

    // 同步标记（markers）
    useEffect(() => {
        if (!editorRef.current) return;
        const session = editorRef.current.session;

        // 清理旧的标记
        if (markerIdsRef.current.length) {
            markerIdsRef.current.forEach((id) => {
                try {
                    session.removeMarker(id);
                } catch {
                }
            });
            markerIdsRef.current = [];
        }

        if (Array.isArray(markers) && markers.length > 0) {
            try {
                const aceAny = (window as any).ace;
                const RangeCtor = aceAny?.require?.('ace/range')?.Range as any;
                if (RangeCtor) {
                    markers.forEach((m) => {
                        const range = new RangeCtor(m.startRow, m.startCol, m.endRow, m.endCol);
                        const id = session.addMarker(range, m.className, m.type);
                        markerIdsRef.current.push(id);
                    });
                }
            } catch {
            }
        }
    }, [markers]);

    // 处理尺寸变化，直接更新外层和稳定容器的样式
    useLayoutEffect(() => {
        const newWidth = typeof width === 'number' ? `${width}px` : width;
        const newHeight = typeof height === 'number' ? `${height}px` : height;

        // 更新外层容器尺寸
        if (outerContainerRef.current) {
            outerContainerRef.current.style.width = newWidth;
            outerContainerRef.current.style.height = newHeight;
        }

        // 更新稳定容器尺寸
        if (stableContainer.current) {
            stableContainer.current.style.width = newWidth;
            stableContainer.current.style.height = newHeight;
        }

        // 通知编辑器重新计算大小
        if (editorRef.current) {
            editorRef.current.resize();
        }
    }, [width, height]);

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
        <div className={className} style={style} ref={outerContainerRef}>
            {isAcePresent ? (
                <div ref={containerRef} style={containerBaseStyle as any}/>
            ) : (
                <Script
                    src={coreScript}
                    parallel={false}
                    timeout={10000}
                    retryCount={2}
                    cache={true}
                    onLoad={useCallback(() => {
                        setAceLoaded(true);
                        setHasScriptLoaded(true);
                    }, [])}
                    loading={hasScriptLoaded ? null : loadingContent}>
                    <div ref={containerRef} style={containerBaseStyle as any}/>
                </Script>
            )}
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
