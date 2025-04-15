<template>
  <div class="wrapper">
    <h1>My Cloud</h1>

    <div class="directory-path">
      <p>Путь: <strong>{{ path }}</strong></p>
    </div>

    <button @click="goBack" class="back-button" :disabled="isRoot">⬅ Назад</button>

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
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';

const path = ref('');
const items = ref([]);

const isRoot = computed(() => path.value === '' || path.value === '/');

const fetchDirectoryContents = async (targetPath = '') => {
  try {
    const response = await axios.get('/api/navigation/', {
      params: { path: targetPath }
    });
    path.value = response.data.path.replace(/\/+$/, ''); // убираем лишние слэши в конце
    items.value = response.data.items;
  } catch (error) {
    console.error("Ошибка загрузки содержимого директории", error);
  }
};

onMounted(() => {
  fetchDirectoryContents();
});

const getIcon = (item) => {
  if (item.is_dir) {
    return '/icons/folder-icon.svg';
  } else {
    const fileType = item.name.split('.').pop().toLowerCase();
    const knownTypes = ['txt', 'pdf', 'js', 'png', 'jpg'];
    return knownTypes.includes(fileType)
      ? `/icons/${fileType}-icon.svg`
      : '/icons/file-icon.svg';
  }
};

const handleClick = (item) => {
  if (item.is_dir) {
    const newPath = `${path.value}/${item.name}`.replace(/\/+/g, '/');
    fetchDirectoryContents(newPath);
  } else if (item.is_file) {
    const fileUrl = `/files${path.value}/${item.name}`.replace(/\/+/g, '/');
    window.open(fileUrl, '_blank');
  }
};

const goBack = () => {
  const segments = path.value.split('/').filter(Boolean);
  segments.pop(); // убираем последнюю папку
  const newPath = '/' + segments.join('/');
  fetchDirectoryContents(newPath);
};
</script>

<style scoped>
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
