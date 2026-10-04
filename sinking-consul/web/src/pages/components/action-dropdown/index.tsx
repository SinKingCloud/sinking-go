import {forwardRef} from "react";
import {Dropdown} from "antd";
import type {DropdownProps} from "antd";
import {createStyles} from "antd-style";
import {useTheme} from "sinking-antd";

const useStyles = createStyles(({css, token}, {compact}: {compact: boolean}) => ({
    menu: css`
        && {
            min-width: ${compact ? 94 : 104}px;
            max-width: calc(100vw - 24px);
            padding: 0;
            border-radius: ${token.borderRadius}px;
        }

        && .ant-dropdown-menu {
            padding: 3px !important;
            border-radius: ${token.borderRadius}px !important;
            background: ${token.colorBgElevated};
            box-shadow: ${token.boxShadowSecondary};
        }

        && .ant-dropdown-menu-item {
            min-height: ${compact ? 28 : 30}px;
            padding: 0 ${compact ? 7 : 8}px !important;
            border-radius: ${token.borderRadiusSM}px !important;
            color: ${token.colorTextSecondary};
            font-size: ${compact ? 11 : 12}px;
        }

        && .ant-dropdown-menu-item-danger { color: ${token.colorError}; }
        && .ant-dropdown-menu-item-disabled { color: ${token.colorTextDisabled}; }
    `,
}));

// 统一操作菜单外观，触发按钮和菜单内容沿用 Ant Design API。
const ActionDropdown = forwardRef<HTMLElement, DropdownProps>(({rootClassName,
    trigger = ["click"], placement = "bottom", arrow = {pointAtCenter: true}, ...props}, ref) => {
    const theme = useTheme();
    const {styles, cx} = useStyles({compact: Boolean(theme?.isCompactTheme?.())});

    return <Dropdown {...props} ref={ref}
        trigger={trigger} placement={placement} arrow={arrow}
        rootClassName={cx(styles.menu, rootClassName)}/>;
});

ActionDropdown.displayName = "ActionDropdown";

export default ActionDropdown;
