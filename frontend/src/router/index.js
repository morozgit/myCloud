import { createRouter, createWebHistory } from 'vue-router';
import WelcomePage from '../WelcomePage.vue';
import CloudPage from '@/components/CloudPage.vue';

const routes = [
  {
    path: '/',
    name: 'WelcomePage',
    component: WelcomePage,
  },
  {
    path: '/cloud/:pathMatch(.*)*',  // <-- универсальный маршрут
    name: 'CloudPage',
    component: CloudPage,
    props: true, // позволяет передавать :pathMatch как пропс
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
