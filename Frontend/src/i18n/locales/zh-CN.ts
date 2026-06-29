// ============================================================
// 文件: i18n/locales/zh-CN.ts
// 说明: 中文语言包
//       所有中文文案集中管理在此文件中
//       拓展方式：新增页面时在此文件对应位置添加文案即可
// ============================================================

// 导出中文语言包对象
export default {
  // --- 通用 ---
  common: {
    home: '首页',
    articles: '技术文章',
    categories: '技术专栏',
    about: '关于',
    search: '搜索',
    loading: '加载中...',
    noData: '暂无数据',
    loadMore: '加载更多',
    viewCount: '阅读',
    publishedAt: '发布于',
  },

  // --- 导航菜单 ---
  nav: {
    embeddedLinux: '嵌入式Linux',
    hardwareDesign: '硬件电路设计',
    mcuDevelopment: '单片机开发',
    hotArticles: '热门文章',
    latestArticles: '最新文章',
    writeArticle: '写文章',
  },

  // --- 首页 ---
  homePage: {
    title: '开源熊猫',
    subtitle: '嵌入式与硬件开发技术交流平台',
    description: '嵌入式Linux、硬件电路设计、单片机开发、FPGA开发等技术交流与学习平台，分享技术文章、项目实战经验和开发技巧',
    hotArticles: '热门文章',
    latestArticles: '最新文章',
  },

  // --- 文章 ---
  article: {
    detail: '文章详情',
    list: '文章列表',
    tags: '标签',
    category: '分类',
    relatedArticles: '相关文章',
    backToList: '返回列表',
    searchPlaceholder: '搜索技术文章...',
    searchResult: '搜索结果',
    noSearchResult: '未找到相关文章',
  },

  // --- 分类 ---
  category: {
    embeddedLinux: '嵌入式Linux',
    embeddedLinuxDesc: 'Linux环境搭建、驱动开发、系统移植实操记录',
    hardwareDesign: '硬件电路设计',
    hardwareDesignDesc: '原理图设计、PCB布局、硬件调试技巧',
    mcuDevelopment: '单片机开发',
    mcuDevelopmentDesc: 'STM32、ESP32等单片机开发教程与项目实战',
  },

  // --- 页脚 ---
  footer: {
    copyright: '© 2026 开源熊猫 - 嵌入式技术交流平台',
    poweredBy: 'Powered by Vue3 + Gin + PostgreSQL',
  },
}
