import { theme } from 'ant-design-vue'

const BRAND = '#6366f1'

// ant-design-vue 4 的组件 token 命名与 antd v5 不同，
// 以下键名均以 es/<component>/style/ 源码为准
const sharedToken = {
  colorPrimary: BRAND,
  colorInfo: BRAND,
  colorSuccess: '#10b981',
  colorWarning: '#f59e0b',
  colorError: '#ef4444',
  borderRadius: 8,
  borderRadiusLG: 12,
  borderRadiusSM: 6,
  fontFamily: `-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', Arial, sans-serif`
}

const sharedComponents = {
  Menu: {
    // 菜单主题固定用 light（配合算法自动适配暗色），这里只放双主题通用的值
    radiusItem: 10,
    itemMarginInline: 12
  },
  Table: {
    tableHeaderCellSplitColor: 'transparent'
  },
  Card: {
    paddingLG: 22
  },
  Button: {
    controlHeight: 34,
    paddingContentHorizontal: 16
  },
  Layout: {
    layoutHeaderHeight: 60,
    layoutHeaderPaddingInline: 20
  }
}

const lightTheme = {
  algorithm: theme.defaultAlgorithm,
  token: {
    ...sharedToken,
    colorBgLayout: '#f4f5f9',
    colorBgContainer: '#ffffff',
    colorBgElevated: '#ffffff',
    colorBorder: '#e4e6ee',
    colorBorderSecondary: '#eef0f5',
    colorText: '#1c2333',
    colorTextSecondary: '#5f6b85'
  },
  components: {
    ...sharedComponents,
    Layout: {
      ...sharedComponents.Layout,
      // 仅作用于 header（sider 使用 theme="light"，走 colorBgContainer）
      colorBgHeader: 'rgba(255, 255, 255, 0.75)'
    },
    Table: {
      ...sharedComponents.Table,
      tableHeaderBg: '#fafbfd',
      tableHeaderTextColor: '#5f6b85',
      tableRowHoverBg: 'rgba(99, 102, 241, 0.05)'
    }
  }
}

const darkTheme = {
  algorithm: theme.darkAlgorithm,
  token: {
    ...sharedToken,
    colorBgLayout: '#0e1017',
    colorBgContainer: '#151824',
    colorBgElevated: '#1c2030',
    colorBorder: '#272d40',
    colorBorderSecondary: '#202536',
    colorText: 'rgba(255, 255, 255, 0.88)',
    colorTextSecondary: 'rgba(255, 255, 255, 0.55)'
  },
  components: {
    ...sharedComponents,
    Layout: {
      ...sharedComponents.Layout,
      colorBgHeader: 'rgba(14, 16, 23, 0.75)'
    },
    Table: {
      ...sharedComponents.Table,
      tableHeaderBg: '#191d2b',
      tableHeaderTextColor: 'rgba(255, 255, 255, 0.55)',
      tableRowHoverBg: 'rgba(99, 102, 241, 0.10)'
    }
  }
}

export function getThemeConfig(mode) {
  return mode === 'dark' ? darkTheme : lightTheme
}
