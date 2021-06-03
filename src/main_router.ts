import {createRouter, createWebHistory} from 'vue-router'
import Events from './pages/Events.vue'
import Layout from './pages/Layout.vue'

const routes = [
    {
        path: '',
        component: Layout,
        children: [
            {
                path: '',
                name: 'dashboard',
                component: Events,
                meta: {title: 'Events'},
            },
            {
                path: '/search',
                name: 'search',
                component: Events,
                meta: {title: 'Search'},
            },
            {
                path: '/events/:informerName',
                name: 'events',
                component: Events,
                meta: {title: 'Events'},
            },
        ],
    },
    {
        path: '/api:pathMatch(.*)*',
        component: () => import('./pages/errors/RenderedByBackendPage.vue'),
        meta: {title: 'Not Found'}
    },
    {
        path: '/:pathMatch(.*)*',
        component: () => import('./pages/errors/404.vue'),
        meta: {title: 'Not Found'},
    },
]


const router = createRouter({
    history: createWebHistory(),
    routes: routes,
})

router.beforeEach((to, from, next) => {
    if (typeof to.meta.title === "string") {
        document.title = to.meta.title
    } else {
        document.title = (to.name ? to.name.toString() : 'Untitled')
    }
    next()
})

export default router