import {get} from "@/utils/request";
import {useState, useEffect} from "react";

/**
 * 枚举缓存接口
 */
interface EnumCache {
    data: Record<string, string>;
    timestamp: number;
}

/**
 * 枚举服务类 - 提供高性能的枚举数据获取和缓存功能
 */
class EnumService {
    private cache = new Map<string, EnumCache>();
    private readonly CACHE_EXPIRE_TIME = 5 * 60 * 1000; // 5分钟缓存过期时间
    private readonly pendingRequests = new Map<string, Promise<any>>();

    /**
     * 获取单个枚举数据
     * @param name 枚举标识名
     * @returns 原始枚举数据 Record<string, string>
     */
    async getEnum(name: string): Promise<Record<string, string>> {
        // 检查缓存
        const cached = this.getFromCache(name);
        if (cached) {
            return cached;
        }

        // 检查是否有正在进行的请求，避免重复请求
        if (this.pendingRequests.has(name)) {
            return await this.pendingRequests.get(name)!;
        }

        // 创建请求Promise
        const requestPromise = this.fetchEnum(name);
        this.pendingRequests.set(name, requestPromise);
        try {
            return await requestPromise;
        } finally {
            this.pendingRequests.delete(name);
        }
    }

    /**
     * 批量获取多个枚举数据
     * @param names 枚举标识名数组
     * @returns 包含所有枚举数据的对象
     */
    async getEnums(names: string[]): Promise<Record<string, Record<string, string>>> {
        const promises = names.map(name =>
            this.getEnum(name).then(data => ({
                name,
                data
            }))
        );

        const results = await Promise.all(promises);
        return results.reduce((acc, {name, data}) => {
            acc[name] = data;
            return acc;
        }, {} as Record<string, Record<string, string>>);
    }

    /**
     * 从缓存获取数据
     */
    private getFromCache(name: string): Record<string, string> | null {
        const cached = this.cache.get(name);
        if (cached && (Date.now() - cached.timestamp) < this.CACHE_EXPIRE_TIME) {
            return cached.data;
        }

        // 清理过期缓存
        if (cached) {
            this.cache.delete(name);
        }

        return null;
    }

    /**
     * 请求枚举数据
     */
    private async fetchEnum(name: string): Promise<Record<string, string>> {
        const response = await get(`/admin/system/enum?name=${encodeURIComponent(name)}`);
        if (response?.code !== 200) {
            return {};
        }
        const allEnumData = response?.data || {};
        this.cache.set(name, {
            data: allEnumData,
            timestamp: Date.now()
        });
        return allEnumData;
    }

    /**
     * 检查缓存是否存在且有效
     */
    isCached(name: string): boolean {
        return this.getFromCache(name) !== null;
    }

    /**
     * 直接从缓存获取数据，不发起网络请求
     */
    getCachedData(name: string): Record<string, string> | null {
        return this.getFromCache(name);
    }

    /**
     * 清空缓存
     */
    clearCache(name?: string): void {
        if (name) {
            this.cache.delete(name);
        } else {
            this.cache.clear();
        }
    }

    /**
     * 获取缓存状态
     */
    getCacheInfo(): Array<{ name: string; timestamp: number; expired: boolean }> {
        const info: Array<{ name: string; timestamp: number; expired: boolean }> = [];

        this.cache.forEach((value, key) => {
            info.push({
                name: key,
                timestamp: value.timestamp,
                expired: (Date.now() - value.timestamp) >= this.CACHE_EXPIRE_TIME
            });
        });

        return info;
    }
}

// 导出单例实例
export const enumService = new EnumService();

// 导出便捷方法
export const getEnum = (name: string) => enumService.getEnum(name);

export const getEnums = (names: string[]) => enumService.getEnums(names);

// ==================== React Hooks ====================

/**
 * 获取单个枚举的Hook
 * @param name 枚举标识名
 * @returns [enumData, loading, error] - enumData是原始的Record<string, string>格式
 */
export const useEnum = (name: string): [Record<string, string>, boolean, Error | null] => {
    const [enumData, setEnumData] = useState<Record<string, string>>(() => enumService.getCachedData(name) || {});
    const [loading, setLoading] = useState(() => Boolean(name && !enumService.isCached(name)));
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        const cached = name ? enumService.getCachedData(name) : null;
        setError(null);
        if (!name || cached) {
            setEnumData(cached || {});
            setLoading(false);
            return;
        }
        let mounted = true;
        setLoading(true);
        enumService.getEnum(name).then(data => {
            if (mounted) setEnumData(data);
        }).catch(err => {
            if (mounted) setError(err instanceof Error ? err : new Error('获取枚举数据失败'));
        }).finally(() => {
            if (mounted) setLoading(false);
        });
        return () => { mounted = false; };
    }, [name]);

    return [enumData, loading, error];
};

/** 批量读取枚举，名称顺序不影响请求，且不修改调用方的数组。 */
export const useEnums = (names: string[]): [Record<string, Record<string, string>>, boolean, Error | null] => {
    const namesKey = JSON.stringify([...new Set(names || [])].filter(Boolean).sort());
    const readCache = () => Object.fromEntries((JSON.parse(namesKey) as string[])
        .map(name => [name, enumService.getCachedData(name)])
        .filter(([, data]) => data !== null)) as Record<string, Record<string, string>>;
    const [enumsData, setEnumsData] = useState(readCache);
    const [loading, setLoading] = useState(() => (JSON.parse(namesKey) as string[]).some(name => !enumService.isCached(name)));
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        const requestedNames = JSON.parse(namesKey) as string[];
        const cachedData = readCache();
        const missing = requestedNames.filter(name => !cachedData[name]);
        setError(null);
        setEnumsData(previous => Object.keys(previous).length === Object.keys(cachedData).length &&
            Object.keys(cachedData).every(name => previous[name] === cachedData[name]) ? previous : cachedData);
        setLoading(missing.length > 0);
        if (!missing.length) return;
        let mounted = true;
        enumService.getEnums(missing).then(data => {
            if (mounted) setEnumsData({...cachedData, ...data});
        }).catch(err => {
            if (mounted) setError(err instanceof Error ? err : new Error('获取枚举数据失败'));
        }).finally(() => {
            if (mounted) setLoading(false);
        });
        return () => { mounted = false; };
    }, [namesKey]);

    return [enumsData, loading, error];
};
