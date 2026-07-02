<!--
  文件: views/Home/index.vue
  说明: 首页
        展示：Banner区域、三大技术专栏卡片、最新文章、热门文章
        这是用户访问网站看到的第一个页面

  后续拓展：替换为真实 API 数据、添加轮播图、添加更多内容区域
-->
<template>
  <div class="home">
    <!-- ============================================================
    Hero Banner（顶部横幅区域）
    ============================================================ -->
    <section class="home__hero">
      <img src="/panda.png" alt="OpenPanda" class="hero__logo" />
      <!-- <h1 class="hero__title">{{ $t('homePage.title') }}</h1> -->
      <p class="hero__subtitle">{{ $t('homePage.subtitle') }}</p>
      <p class="hero__desc">{{ $t('homePage.description') }}</p>
    </section>

    <!-- ============================================================
    三大技术专栏卡片
    ============================================================ -->
    <section class="home__categories">
      <h2 class="section__title">{{ $t('common.categories') }}</h2>
      <div class="categories__grid">
        <!-- 使用 el-card 展示每个分类 -->
        <el-card
          v-for="cat in categories"
          :key="cat.slug"
          class="category__card"
          shadow="hover"
          @click="router.push(`/category/${cat.slug}`)"
        >
          <h3>{{ cat.name }}</h3>
          <p>{{ cat.description }}</p>
        </el-card>
      </div>
    </section>

    <!-- ============================================================
    最新文章列表
    ============================================================ -->
    <section class="home__latest">
      <h2 class="section__title">{{ $t('homePage.latestArticles') }}</h2>
      <div class="articles__list">
        <!-- v-for 遍历文章列表渲染 -->
        <el-card
          v-for="article in latestArticles"
          :key="article.id"
          class="article__card"
          shadow="hover"
          @click="goToArticle(article)"
        >
          <div class="article__info">
            <h3>{{ article.title }}</h3>
            <p class="article__summary">{{ article.summary || article.content.slice(0, 100) + '...' }}</p>
            <div class="article__meta">
              <span v-if="article.category">{{ article.category.name }}</span>
              <span>{{ $t('common.viewCount') }}: {{ article.view_count }}</span>
            </div>
          </div>
        </el-card>

        <!-- 数据为空时显示 -->
        <el-empty v-if="latestArticles.length === 0" :description="$t('common.noData')" />
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
// ============================================================
// 导入
// ============================================================
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getArticles, getCategories } from '@/api/modules/article'
import { getArticleUrl } from '@/utils'
import type { Article, Category } from '@/types'

const router = useRouter()

// ============================================================
// 响应式数据
// ============================================================

/** 专栏分类（从 API 动态加载，与专栏管理保持一致） */
const categories = ref<Category[]>([])

/** 最新文章列表 */
const latestArticles = ref<Article[]>([])

/** 热门文章列表（首页暂未显示，预留后续使用） */
// const hotArticles = ref<Article[]>([]) // 取消注释即可使用

// ============================================================
// 生命周期钩子
// onMounted: 组件挂载到 DOM 后执行（类似于以前的 mounted）
// 适合在此处发起 API 请求加载数据
// ============================================================
onMounted(async () => {
  await Promise.all([fetchCategories(), fetchLatestArticles()])
})

// ============================================================
// 方法
// ============================================================

/** 从 API 加载专栏分类（与专栏管理页面数据一致） */
async function fetchCategories(): Promise<void> {
  try {
    const { data } = await getCategories()
    if (data.data) {
      categories.value = data.data
    }
  } catch {
    console.error('获取分类列表失败')
  }
}

/** 获取最新文章 */
async function fetchLatestArticles(): Promise<void> {
  try {
    // 发送 API 请求，解构取出 data.data.list
    const { data } = await getArticles({ page: 1, page_size: 6 })
    // data.data 是 ApiResponse<PaginatedData<Article>>
    // data.data.data.list 才是文章数组
    if (data.data && data.data.list) {
      latestArticles.value = data.data.list
    }
  } catch (error) {
    // 接口未启动时静默失败，页面显示空状态
    console.error('获取文章列表失败:', error)
  }
}

/** 跳转到文章详情页 */
function goToArticle(article: Article): void {
  router.push(getArticleUrl(article))
}
</script>

<style scoped>
.home {
  /* 页面通用样式 */
}

/* Hero Banner */
.home__hero {
  text-align: center;
  padding: 60px 20px;
  background: linear-gradient(135deg, #c8754a 0%, #d4946e 100%);
  border-radius: 12px;
  color: #fff;
  margin-bottom: 40px;
}
.hero__logo {
  width: 100px;
  height: 100px;
  border-radius: 18px;
  margin-bottom: 20px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.2);
}
.hero__title {
  font-size: 42px;
  font-weight: bold;
  margin-bottom: 12px;
}
.hero__subtitle {
  font-size: 20px;
  opacity: 0.9;
  margin-bottom: 8px;
}
.hero__desc {
  font-size: 14px;
  opacity: 0.7;
}

/* 通用区块标题 */
.section__title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 24px;
  padding-bottom: 12px;
  border-bottom: 2px solid #c8754a;
}

/* 分类卡片网格 */
.categories__grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
}
.category__card {
  cursor: pointer;
  transition: transform 0.2s;
}
.category__card:hover {
  transform: translateY(-4px);
}
.category__card h3 {
  font-size: 18px;
  margin-bottom: 8px;
  color: #c8754a;
}
.category__card p {
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}
.category__link {
  font-size: 13px;
  font-weight: 700;
  color: #9e9e9d;
}

/* 文章列表 */
.articles__list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}
.article__card {
  cursor: pointer;
  transition: transform 0.2s;
}
.article__card:hover {
  transform: translateX(4px);
}
.article__info h3 {
  font-size: 18px;
  margin-bottom: 8px;
}
.article__summary {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.article__meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: #999;
}

/* 响应式：手机端 */
@media (max-width: 768px) {
  .home__hero {
    padding: 40px 16px;
  }
  .hero__title {
    font-size: 28px;
  }
  .hero__subtitle {
    font-size: 16px;
  }
}
</style>
