import { createApp } from 'vue';
import App from './App.vue'; // Это может быть App.vue, который использует router-view
import router from './router';
import toastr from 'toastr';
import 'toastr/build/toastr.min.css';

toastr.options = {
  closeButton: true,
  positionClass: 'toast-bottom-right',
  timeOut: '3000',
};

// Импортируем файл стилей
import './style.css';

createApp(App)
  .use(router)
  .mount('#app');
