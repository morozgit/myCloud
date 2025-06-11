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
    path: '/cloud/:pathMatch(.*)*',
    name: 'CloudPage',
    component: CloudPage,
    props: true,
  },
];

const router = createRouter({
  history: createWebHistory('/mycloud/'),
  routes,
});

export default router;
