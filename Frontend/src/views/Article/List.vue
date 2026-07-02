<!--
  文件: views/Article/List.vue
  说明: 文章列表页
        支持：分页、按分类筛选、搜索
-->
<template>
  <div class="article-list">
    <h2 class="page__title">{{ $t('article.list') }}</h2>

    <!-- ============================================================
    搜索栏
    ============================================================ -->
    <div class="search-bar">
      <!-- v-model: 双向绑定，输入变化时 keyword 自动更新 -->
      <el-input
        v-model="keyword"
        :placeholder="$t('article.searchPlaceholder')"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #append>
          <el-button :icon="SearchIcon" @click="handleSearch" />
        </template>
      </el-input>
    </div>

    <!-- ============================================================
    文章列表
    ============================================================ -->
    <div v-loading="loading" class="articles__list">
      <el-card
        v-for="article in articles"
        :key="article.id"
        class="article__card"
        shadow="hover"
        @click="goToDetail(article)"
      >
        <div class="article__card-body">
          <!-- 封面图 -->
          <div v-if="article.cover_image" class="article__cover">
            <img :src="article.cover_image" :alt="article.title" />
          </div>
          <div class="article__content">
            <h3>{{ article.title }}</h3>
            <p>{{ article.summary || article.content.replace(/<[^>]*>/g, '').slice(0, 150) + '...' }}</p>
            <div class="article__meta">
              <el-tag
                v-if="article.category"
                size="small"
                type="primary"
              >
                {{ article.category.name }}
              </el-tag>
              <el-tag
                v-for="tag in article.tags"
                :key="tag.id"
                size="small"
                class="tag__item"
              >
                {{ tag.name }}
              </el-tag>
              <span>{{ $t('common.viewCount') }}: {{ article.view_count }}</span>
              <span>{{ $t('common.publishedAt') }} {{ formatDate(article.created_at) }}</span>
            </div>
          </div>
        </div>
      </el-card>

      <el-empty v-if="!loading && articles.length === 0" :description="$t('common.noData')" />
    </div>

    <!-- ============================================================
    分页器
    ============================================================ -->
    <div v-if="total > 0" class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 30]"
        layout="total, sizes, prev, pager, next"
        @current-change="handlePageChange"
        @size-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search as SearchIcon } from '@element-plus/icons-vue'
import { getArticles, searchArticles } from '@/api/modules/article'
import { getArticleUrl } from '@/utils'
import type { Article } from '@/types'

const router = useRouter()

// ============================================================
// 响应式数据
// ============================================================

/** 文章列表 */
const articles = ref<Article[]>([])
/** 加载状态 */
const loading = ref<boolean>(false)
/** 当前页码 */
const currentPage = ref<number>(1)
/** 每页条数 */
const pageSize = ref<number>(10)
/** 总记录数 */
const total = ref<number>(0)
/** 搜索关键词 */
const keyword = ref<string>('')

// ============================================================
// 生命周期
// ============================================================
onMounted(() => {
  fetchArticles()
})

// ============================================================
// 方法
// ============================================================

/** 获取文章列表 */
async function fetchArticles(): Promise<void> {
  loading.value = true
  try {
    const { data } = keyword.value
      ? await searchArticles({ keyword: keyword.value, page: currentPage.value, page_size: pageSize.value })
      : await getArticles({ page: currentPage.value, page_size: pageSize.value })

    if (data.data) {
      articles.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取文章列表失败:', error)
  } finally {
    loading.value = false  // finally 块无论成功失败都会执行
  }
}

/** 搜索 */
function handleSearch(): void {
  currentPage.value = 1  // 搜索时重置页码
  fetchArticles()
}

/** 翻页 */
function handlePageChange(): void {
  fetchArticles()
}

/** 跳转详情 */
function goToDetail(article: Article): void {
  router.push(getArticleUrl(article))
}

/** 格式化日期（ISO字符串 → 中文日期） */
function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN')
}
</script>

<style scoped>
.article-list {
  /* 页面容器 */
}

.page__title {
  font-size: 24px;
  margin-bottom: 20px;
}

.search-bar {
  margin-bottom: 24px;
  max-width: 500px;
}

.articles__list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.article__card {
  cursor: pointer;
}
.article__card-body {
  display: flex;
  gap: 20px;
}
.article__cover {
  width: 200px;
  height: 130px;
  flex-shrink: 0;
  border-radius: 8px;
  overflow: hidden;
}
.article__cover img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
.article__content h3 {
  font-size: 18px;
  margin-bottom: 8px;
}
.article__content p {
  color: #666;
  font-size: 14px;
  margin-bottom: 8px;
}
.article__meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #999;
}
.tag__item {
  margin-left: 0 !important;
}

.pagination {
  margin-top: 24px;
  display: flex;
  justify-content: center;
}

/* 手机端适配 */
@media (max-width: 768px) {
  .article__card-body {
    flex-direction: column;
  }
  .article__cover {
    width: 100%;
    height: 180px;
  }
}
</style>
