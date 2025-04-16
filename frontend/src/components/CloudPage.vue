<template>
  <div class="wrapper">
    <h1>My Cloud</h1>

    <!-- Навигация по пути -->
    <div class="directory-path">
      <span>Путь: </span>
      <span v-for="(segment, index) in pathSegments" :key="index">
        <span @click="goTo(index)" class="breadcrumb">
          {{ segment}}
        </span>
        <span v-if="index < pathSegments.length - 1"> / </span>
      </span>
    </div>

    <!-- Список файлов и папок -->
    <div class="file-list">
      <div
        v-for="item in items"
        :key="item.name"
        class="file-item"
        @click="handleClick(item)"
      >
        <div class="file-icon">
          <img :src="getIcon(item)" alt="icon" />
        </div>
        <div class="file-details">
          <p class="file-name">{{ item.name }}</p>
          <p v-if="item.is_file" class="file-size">{{ item.size }} bytes</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';

const route = useRoute();
const router = useRouter();

const path = ref('');
const items = ref([]);

const fetchDirectoryContents = async (targetPath = '') => {
  try {
    const query = targetPath ? `?path=${encodeURIComponent(targetPath)}` : '';
    const response = await axios.get(`/api/navigation/${query}`);
    path.value = response.data.path;
    items.value = response.data.items;
  } catch (error) {
    console.error("Ошибка загрузки содержимого директории", error);
  }
};

const routePath = computed(() => {
  const param = route.params.pathMatch;
  return Array.isArray(param) ? '/' + param.join('/') : (param ? '/' + param : '');
});

onMounted(() => {
  fetchDirectoryContents(routePath.value);
});

// следим за сменой пути при нажатии кнопок "вперёд/назад"
watch(route, () => {
  fetchDirectoryContents(routePath.value);
});

const handleClick = (item) => {
  if (item.is_dir) {
    const newPath = `${routePath.value}/${item.name}`.replace(/\/+/g, '/');
    router.push(`/cloud${newPath}`);
  } else if (item.is_file) {
    const fileUrl = `/files${path.value}/${item.name}`.replace(/\/+/g, '/');
    window.open(fileUrl, '_blank');
  }
};

const getIcon = (item) => {
  if (item.is_dir) return '/icons/folder-icon.svg';
  const type = item.name.split('.').pop().toLowerCase();
  return ['/txt', 'pdf', 'js', 'png', 'jpg'].includes(type)
    ? `/icons/${type}-icon.svg`
    : '/icons/file-icon.svg';
};

const pathSegments = computed(() => {
  return path.value.split('/').filter(Boolean);
});

const goTo = (index) => {
  const newPath = '/' + pathSegments.value.slice(0, index + 1).join('/');
  router.push(`/cloud${newPath}`);
};
</script>

<style scoped>
.breadcrumb {
  color: #42b983;
  cursor: pointer;
  font-weight: bold;
}

.breadcrumb:hover {
  text-decoration: underline;
}

.wrapper {
  padding: 20px;
  font-family: Arial, sans-serif;
  max-width: 100%;
}

.directory-path {
  margin-bottom: 10px;
  color: #aaa;
}

.back-button {
  margin-bottom: 20px;
  padding: 0.5rem 1rem;
  background-color: #333;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.back-button:disabled {
  background-color: #555;
  cursor: not-allowed;
}

.file-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 20px;
  width: 100%;
}

.file-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  background-color: #1e1e1e;
  padding: 12px;
  border-radius: 8px;
  text-align: center;
  transition: background 0.2s ease;
  cursor: pointer;
}

.file-item:hover {
  background-color: #2a2a2a;
}

.file-icon {
  width: 64px;
  height: 64px;
  margin-bottom: 8px;
}

.file-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.file-name {
  font-weight: 500;
  font-size: 0.95em;
  word-break: break-word;
  color: #e0e0e0;
}

.file-size {
  font-size: 0.8em;
  color: #a0a0a0;
}
</style>
