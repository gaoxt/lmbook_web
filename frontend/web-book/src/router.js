import Vue from 'vue';
import VueRouter from 'vue-router';

import BookDetail from "./views/BookDetail.vue";
import BookList from "./views/BookList.vue";

Vue.use(VueRouter);

const router = new VueRouter({
    mode: 'history',
    routes: [
        { path: '/', component: BookList },
        { path: '/detail/:id', component: BookDetail, props: (route) => ({ id: route.query.id }) },
    ]
});

export default router;