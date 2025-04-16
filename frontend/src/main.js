import { createApp } from 'vue';
import App from './App.vue'; // Это может быть App.vue, который использует router-view
import router from './router';

// Импортируем файл стилей
import './style.css';

createApp(App)
  .use(router)
  .mount('#app');