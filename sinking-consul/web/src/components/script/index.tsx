import React, {useEffect, useRef, useState, useCallback, useMemo} from 'react';

// 脚本加载状态枚举
export enum ScriptLoadStatus {
    IDLE = 'idle',
    LOADING = 'loading',
    LOADED = 'loaded',
    ERROR = 'error'
}

// 脚本加载配置接口
export interface ScriptConfig {
    src: string;                    // 脚本URL
    async?: boolean;               // 是否异步加载
    defer?: boolean;               // 是否延迟执行
    crossOrigin?: 'anonymous' | 'use-credentials'; // 跨域设置
    integrity?: string;            // 完整性校验
    noModule?: boolean;           // 是否为非模块脚本
    type?: string;                // 脚本类型
    charset?: string;             // 字符编码
    timeout?: number;             // 加载超时时间(ms)
    retryCount?: number;          // 重试次数
    retryDelay?: number;          // 重试延迟(ms)
    cache?: boolean;              // 是否缓存
    removeOnUnmount?: boolean;    // 组件卸载时是否移除脚本
    onLoad?: () => void;          // 加载成功回调
    onError?: (error: Error) => void; // 加载失败回调
    onTimeout?: () => void;       // 超时回调
}

// 脚本加载组件属性
export interface ScriptProps extends Omit<ScriptConfig, 'src'> {
    src: string | string[];       // 支持单个或多个脚本
    children?: React.ReactNode;   // 子组件
    fallback?: React.ReactNode;   // 加载失败时的备用内容
    loading?: React.ReactNode;    // 加载中的显示内容
    parallel?: boolean;           // 多脚本是否并行加载
}

// 全局脚本缓存
const scriptCache = new Map<string, Promise<void>>();
const loadedScripts = new Set<string>();

/**
 * 创建脚本元素
 */
const createScriptElement = (config: ScriptConfig): HTMLScriptElement => {
    const script = document.createElement('script');

    script.src = config.src;
    if (config.async !== undefined) script.async = config.async;
    if (config.defer !== undefined) script.defer = config.defer;
    if (config.crossOrigin) script.crossOrigin = config.crossOrigin;
    if (config.integrity) script.integrity = config.integrity;
    if (config.noModule !== undefined) script.noModule = config.noModule;
    if (config.type) script.type = config.type;
    if (config.charset) script.charset = config.charset;

    return script;
};

/**
 * 加载单个脚本
 */
const loadScript = (config: ScriptConfig): Promise<void> => {
    const {src, timeout = 10000, retryCount = 2, retryDelay = 1000, cache = true} = config;

    // 检查缓存
    if (cache && loadedScripts.has(src)) {
        return Promise.resolve();
    }

    if (cache && scriptCache.has(src)) {
        return scriptCache.get(src)!;
    }

    const loadPromise = new Promise<void>((resolve, reject) => {
        let retryAttempts = 0;

        const attemptLoad = () => {
            const script = createScriptElement(config);
            let timeoutId: NodeJS.Timeout | null = null;

            const cleanup = () => {
                if (timeoutId) clearTimeout(timeoutId as any);
                script.removeEventListener('load', onLoad);
                script.removeEventListener('error', onError);
            };

            const onLoad = () => {
                cleanup();
                loadedScripts.add(src);
                config.onLoad?.();
                resolve();
            };

            const onError = () => {
                cleanup();
                document.head.removeChild(script);

                if (retryAttempts < retryCount) {
                    retryAttempts++;
                    setTimeout(attemptLoad, retryDelay);
                } else {
                    const error = new Error(`Failed to load script: ${src}`);
                    config.onError?.(error);
                    reject(error);
                }
            };

            const onTimeout = () => {
                cleanup();
                document.head.removeChild(script);

                if (retryAttempts < retryCount) {
                    retryAttempts++;
                    setTimeout(attemptLoad, retryDelay);
                } else {
                    const error = new Error(`Script load timeout: ${src}`);
                    config.onTimeout?.();
                    reject(error);
                }
            };

            script.addEventListener('load', onLoad);
            script.addEventListener('error', onError);

            if (timeout > 0) {
                timeoutId = setTimeout(onTimeout, timeout);
            }

            document.head.appendChild(script);
        };

        attemptLoad();
    });

    if (cache) {
        scriptCache.set(src, loadPromise);
    }

    return loadPromise;
};

/**
 * Script 组件
 */
const Script: React.FC<ScriptProps> = ({
    src,
    children,
    fallback,
    loading,
    parallel = true,
    removeOnUnmount = false,
    ...config
}) => {
    const [status, setStatus] = useState<ScriptLoadStatus>(ScriptLoadStatus.IDLE);
    const loadedScriptsRef = useRef<Set<string>>(new Set());

    // 标准化脚本配置
    const scriptConfigs = useMemo(() => {
        const urls = Array.isArray(src) ? src : [src];
        return urls.map(url => ({
            ...config,
            src: url
        }));
    }, [src, config]);

    // 加载脚本的核心逻辑
    const loadScripts = useCallback(async () => {
        if (scriptConfigs.length === 0) return;

        setStatus(ScriptLoadStatus.LOADING);

        try {
            if (parallel) {
                // 并行加载
                await Promise.all(scriptConfigs.map(loadScript));
            } else {
                // 串行加载
                for (const scriptConfig of scriptConfigs) {
                    await loadScript(scriptConfig);
                }
            }

            scriptConfigs.forEach(({src}) => {
                loadedScriptsRef.current.add(src);
            });

            setStatus(ScriptLoadStatus.LOADED);
        } catch (err) {
            setStatus(ScriptLoadStatus.ERROR);
        }
    }, [scriptConfigs, parallel]);

    // 组件挂载时加载脚本
    useEffect(() => {
        loadScripts();
    }, [loadScripts]);

    // 组件卸载时清理
    useEffect(() => {
        return () => {
            if (removeOnUnmount) {
                loadedScriptsRef.current.forEach(src => {
                    const script = document.querySelector(`script[src="${src}"]`);
                    if (script) {
                        document.head.removeChild(script);
                        loadedScripts.delete(src);
                        scriptCache.delete(src);
                    }
                });
            }
        };
    }, [removeOnUnmount]);

    // 渲染逻辑
    switch (status) {
        case ScriptLoadStatus.LOADING:
            return loading ? <>{loading}</> : null;

        case ScriptLoadStatus.ERROR:
            return fallback ? <>{fallback}</> : null;

        case ScriptLoadStatus.LOADED:
            return children ? <>{children}</> : null;

        default:
            return null;
    }
};

// 工具函数：预加载脚本
export const preloadScript = (src: string, config?: Partial<ScriptConfig>): Promise<void> => {
    return loadScript({src, ...config});
};

// 工具函数：预加载多个脚本
export const preloadScripts = (
    scripts: string[],
    config?: Partial<ScriptConfig>,
    parallel = true
): Promise<void[]> => {
    const scriptConfigs = scripts.map(src => ({src, ...config}));

    if (parallel) {
        return Promise.all(scriptConfigs.map(loadScript));
    } else {
        return scripts.reduce(
            (promise, src) => promise.then(() => loadScript({src, ...config})),
            Promise.resolve()
        ).then(() => []);
    }
};

// 工具函数：检查脚本是否已加载
export const isScriptLoaded = (src: string): boolean => {
    return loadedScripts.has(src);
};

// 工具函数：清除脚本缓存
export const clearScriptCache = (src?: string): void => {
    if (src) {
        scriptCache.delete(src);
        loadedScripts.delete(src);
    } else {
        scriptCache.clear();
        loadedScripts.clear();
    }
};

export default Script;
