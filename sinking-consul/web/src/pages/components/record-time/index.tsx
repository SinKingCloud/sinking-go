import {Space, Typography} from "antd";
import {createStyles} from "antd-style";

interface RecordTimeProps {
    createTime?: string;
    updateTime?: string;
}

const useStyles = createStyles(() => ({
    time: {
        width: "100%",
    },
    row: {
        display: "flex",
        alignItems: "center",
        gap: 4,
        minWidth: 0,
    },
    label: {
        flex: "none",
        color: "inherit",
        fontSize: 12,
        lineHeight: "18px",
    },
    value: {
        flex: 1,
        minWidth: 0,
        color: "inherit",
        fontSize: 12,
        lineHeight: "18px",
        whiteSpace: "nowrap",
        fontVariantNumeric: "tabular-nums",
    },
}));

const displayTime = (value?: string) => {
    const time = String(value || "").trim();
    return !time || /^(0000|0001)-/.test(time) ? "" : time;
};

export default ({createTime, updateTime}: RecordTimeProps) => {
    const {styles} = useStyles();
    const created = displayTime(createTime);
    const updated = displayTime(updateTime);

    return <Space orientation="vertical" size={2} className={styles.time}>
        {[{label: "创建于:", value: created}, {label: "更新于:", value: updated}].map((item) => (
            <div className={styles.row} key={item.label}>
                <Typography.Text className={styles.label}>{item.label}</Typography.Text>
                <Typography.Text className={styles.value} ellipsis={{tooltip: item.value || "-"}}>
                    {item.value || "-"}
                </Typography.Text>
            </div>
        ))}
    </Space>;
};
