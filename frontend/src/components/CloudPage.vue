<template>
  <div class="wrapper">
    <h1 @click="goToCloud" class="cloud-button">My Cloud</h1>

    <!-- Навигация по пути -->
    <div class="directory-path">
      <span>Путь: </span>
      <span v-for="(segment, index) in pathSegments" :key="index">
        <span @click="goTo(index)" class="breadcrumb">
          {{ segment }}
        </span>
        <span v-if="index < pathSegments.length - 1"> / </span>
      </span>
    </div>
  <div class="actions-bar">
    <!-- Кнопка загрузки -->
    <div class="upload-container">
      <button @click.stop="triggerFileUpload" class="upload-button">
        <img src="/icons/upload-icon.svg" alt="Upload" class="upload-icon" />
      </button>
      <span class="upload-text">Загрузить файл</span>
      <input ref="fileInput" type="file" @change="handleFileUpload" style="display: none;" />
    </div>

    <!-- Кнопка сортировки -->
    <div class="sort-bar">
      <span class="sort-label" @click="toggleSortOrder" style="cursor: pointer;">
        <span class="upload-text">Сортировать по имени  </span>
        <span v-if="sortAscending">⬆️</span>
        <span v-else>⬇️</span>
      </span>
    </div>
  </div>
    <!-- Список файлов и папок -->
    <div class="file-list">
      <div
        v-for="item in sortedItems"
        :key="item.name"
        class="file-item"
        @click="handleClick(item)"
      >
        <!-- Кнопка удаления -->
        <button @click.stop="deleteFileWrapper(item)" class="delete-button">
          <img src="/icons/delete-icon.svg" alt="Delete" class="delete-icon" />
        </button>
        <div class="file-icon">
          <img :src="getIcon(item)" alt="icon" />
        </div>
        <div class="file-details">
          <p class="file-name">{{ item.name }}</p>
          <p class="file-size">
            <span v-if="item.is_file">{{ item.size }} bytes</span>
            <span v-else>{{ item.children_count !== null ? item.children_count + ' объектов' : 'Папка' }}</span>
          </p>
          <!-- Кнопка скачивания -->
          <button @click.stop="downloadFileWrapper(item)" class="download-button">
            <img src="/icons/download-icon.svg" alt="Download" class="download-icon" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import axios from 'axios';
import { useDelete } from '@/composables/useDelete';
import { useDownload } from '@/composables/useDownload';
import { uploadFile } from '@/composables/fileUploader';
import './CloudView.css';

const route = useRoute();
const router = useRouter();
const { deleteFile } = useDelete();
const { downloadFile } = useDownload();

const path = ref('');
const items = ref([]);
const fileInput = ref(null);
const apiBase = import.meta.env.VITE_API_BASE;

// Состояние сортировки
const sortAscending = ref(true);
const toggleSortOrder = () => {
  sortAscending.value = !sortAscending.value;
};

// Сортированный список
const sortedItems = computed(() => {
  return [...(items.value || [])].sort((a, b) => {
    if (a.is_file !== b.is_file) {
      return a.is_file ? 1 : -1;
    }
    const result = a.name.localeCompare(b.name);
    return sortAscending.value ? result : -result;
  });
});

const triggerFileUpload = () => {
  fileInput.value.click();
};

const handleFileUpload = async (event) => {
  const file = event.target.files[0];
  if (!file) return;

  try {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('path', path.value);
    await uploadFile(formData);
    await fetchDirectoryContents(path.value);
  } catch (error) {
    console.error('Ошибка при загрузке файла:', error);
  }
};

const fetchDirectoryContents = async (targetPath = '') => {
  try {
    const query = targetPath ? `?path=${encodeURIComponent(targetPath)}` : '';
    const response = await axios.get(`${apiBase}/navigation/${query}`);
    path.value = response.data.path;
    console.log("path.value", response)
    items.value = response.data.items;
  } catch (error) {
    console.error('Ошибка загрузки содержимого директории', error);
  }
};

const routePath = computed(() => {
  const param = route.params.pathMatch;
  return Array.isArray(param) ? '/' + param.join('/') : (param ? '/' + param : '');
});

onMounted(() => {
  fetchDirectoryContents(routePath.value);
});

watch(() => route.path, () => {
  fetchDirectoryContents(routePath.value);
});

const handleClick = (item) => {
  if (item.is_dir) {
    const newPath = `${routePath.value}/${item.name}`.replace(/\/+/g, '/');
    router.push(`/cloud${newPath}`);
  } else if (item.is_file) {
    const fileUrl = `${apiBase}/files${path.value}/${item.name}`.replace(/\/+/g, '/');
    window.open(fileUrl, '_blank');
  }
};

const deleteFileWrapper = async (item) => {
  const fullPath = `${path.value}/${item.name}`.replace(/\/+/g, '/');
  try {
    await deleteFile({ ...item, path: path.value });
    items.value = items.value.filter(i => i.name !== item.name);
  } catch (error) {
    console.error('Ошибка при удалении:', error);
  }
};

const downloadFileWrapper = (item) => {
  const fullPath = `${path.value}/${item.name}`.replace(/\/+/g, '/');
  downloadFile({ ...item, path: fullPath });
};

const getIcon = (item) => {
  if (item.is_dir) return '/mycloud/icons/folder-icon.svg';
  if (item.is_file) return '/mycloud/icons/file-icon.svg';
};

const pathSegments = computed(() => (path.value || '').split('/').filter(Boolean));

const goTo = (index) => {
  const newPath = '/' + pathSegments.value.slice(0, index + 1).join('/');
  router.push(`/cloud${newPath}`);
};

const goToCloud = () => {
  router.push('/cloud');
};
</script>
