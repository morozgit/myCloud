
<template>
  <div class="welcome-page">
    <h1>Добро пожаловать!</h1>
    <p>Перейдите в облачное хранилище для просмотра файлов и папок.</p>
    <button @click="goToCloud">Войти в облако</button>
  </div>
</template>
<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goToCloud = async () => {
  router.push({
      path: '/cloud',
    });
};
</script>

<style scoped>
.welcome-page {
  text-align: center;
  margin-top: 50px;
}

button {
  background-color: #4CAF50;
  color: white;
  padding: 10px 20px;
  border: none;
  cursor: pointer;
  font-size: 16px;
}

button:hover {
  background-color: #45a049;
}
</style>
