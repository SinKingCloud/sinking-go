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
            const result = await requestPromise;
            this.pendingRequests.delete(name);
            return result;
        } catch (error) {
            this.pendingRequests.delete(name);
            throw error;
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
        const response = await get(`/admin/system/enum?name=${name}`);
        if (response.code !== 200) {
            return {};
        }
        const allEnumData = response.data || {};
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
        const cached = this.cache.get(name);
        return Date.now() - cached.timestamp < this.CACHE_EXPIRE_TIME;
    }

    /**
     * 直接从缓存获取数据，不发起网络请求
     */
    getCachedData(name: string): Record<string, string> | null {
        const cached = this.cache.get(name);
        const now = Date.now();
        if (cached) {
            const age = now - cached.timestamp;
            const isValid = age < this.CACHE_EXPIRE_TIME;
            if (isValid) {
                return cached.data;
            }
        }
        return null;
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
    const [enumData, setEnumData] = useState<Record<string, string>>({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (!name) {
            setLoading(false);
            return;
        }

        let mounted = true;

        const fetchEnum = async () => {
            try {
                setLoading(true);
                setError(null);
                const data = await enumService.getEnum(name);

                if (mounted) {
                    setEnumData(data);
                }
            } catch (err) {
                if (mounted) {
                    setError(err instanceof Error ? err : new Error('获取枚举数据失败'));
                }
            } finally {
                if (mounted) {
                    setLoading(false);
                }
            }
        };

        fetchEnum();

        return () => {
            mounted = false;
        };
    }, [name]);

    return [enumData, loading, error];
};

/**
 * 获取多个枚举的Hook
 * @param names 枚举名称数组
 * @returns [enumsData, loading, error] - enumsData是Record<string, Record<string, string>>格式
 */
export const useEnums = (names: string[]): [Record<string, Record<string, string>>, boolean, Error | null] => {
    const [enumsData, setEnumsData] = useState<Record<string, Record<string, string>>>({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (!names || names.length === 0) {
            setLoading(false);
            setEnumsData({});
            return;
        }
        let mounted = true;
        const fetchEnums = async () => {
            try {
                setError(null);
                // 检查哪些枚举需要获取（未缓存的）
                const needFetchNames: string[] = [];
                const cachedData: Record<string, Record<string, string>> = {};
                // 分离已缓存和需要获取的枚举
                for (const name of names) {
                    const cachedItem = enumService.getCachedData(name);
                    if (cachedItem) {
                        // 从缓存直接获取数据，不发起网络请求
                        cachedData[name] = cachedItem;
                    } else {
                        // 缓存不存在或已过期，需要重新获取
                        needFetchNames.push(name);
                    }
                }
                // 如果有需要获取的枚举，显示loading
                if (needFetchNames.length > 0) {
                    setLoading(true);
                    // 先设置已缓存的数据
                    if (Object.keys(cachedData).length > 0 && mounted) {
                        setEnumsData(cachedData);
                    }
                    // 并发获取需要的枚举
                    const results = await enumService.getEnums(needFetchNames);
                    if (mounted) {
                        // 合并缓存数据和新获取的数据
                        const finalData = {...cachedData, ...results};
                        setEnumsData(finalData);
                    }
                } else {
                    // 全部都有缓存，直接设置数据
                    if (mounted) {
                        setEnumsData(cachedData);
                    }
                }

            } catch (err) {
                if (mounted) {
                    setError(err instanceof Error ? err : new Error('获取枚举数据失败'));
                }
            } finally {
                if (mounted) {
                    setLoading(false);
                }
            }
        };

        fetchEnums();

        return () => {
            mounted = false;
        };
    }, [JSON.stringify(names.sort())]); // 排序确保数组顺序不影响缓存

    return [enumsData, loading, error];
};
