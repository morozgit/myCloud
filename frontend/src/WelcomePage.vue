<template>
  <div class="welcome-page">
    <h1>Добро пожаловать!</h1>
    <p>Перейдите в облачное хранилище для просмотра файлов и папок.</p>
    <div class="cloud" @click="goToCloud">
      <img src="/src/assets/cloud.png" alt="Облако" />
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router';

const router = useRouter();

const goToCloud = async () => {
  await router.push({
    path: '/cloud',
  });
};
</script>

<style scoped>
.welcome-page {
  text-align: center;
  margin-top: 50px;
}

h1 {
  font-size: 2rem;
  color: #f0eaea;
}

p {
  font-size: 1.2rem;
  color: #666;
}

.cloud {
  display: inline-block;
  cursor: pointer;
  animation: cloud-animation 2s ease-in-out infinite;
}

.cloud img {
  width: 500px;
  transition: transform 0.3s ease;
}

.cloud:hover img {
  transform: scale(1.1);
}

@keyframes cloud-animation {
  0% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-20px);
  }
  100% {
    transform: translateY(0);
  }
}
</style>
