/**
 * table查询数据格式化
 * @param params 查询参数
 * @param sort 排序参数
 */
export function getParams(params: any, sort: any): any {
    function getOrder(obj: any): any {
        for (const key in obj) {
            const temp: any = {
                order_by_field: key
            };
            if (obj[key] == "ascend") {
                temp.order_by_type = "asc";
            } else {
                temp.order_by_type = "desc";
            }
            return temp;
        }
        return {};
    }

    params.page = params.current;
    params.page_size = params.pageSize;
    delete params.current;
    delete params.pageSize;
    return Object.assign(params, getOrder(sort));
}

/**
 * 获取分页
 * @param params 参数
 * @param sort 排序
 * @param request 请求
 */
export async function getData(params: any, sort: any, request: any): Promise<any> {
    const fetchParams = getParams(params, sort)
    const data = await request({
        body: {
            ...fetchParams
        }
    });
    return {
        data: data?.data?.list === undefined || data?.data?.list === null || data?.data?.list?.length <= 0 ? [] : data?.data?.list,
        success: data?.code === 200,
        total: data?.data?.total,
    };
}