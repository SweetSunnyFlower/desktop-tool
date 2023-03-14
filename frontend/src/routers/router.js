import Header from '../layouts/Header.vue';
import Welcome from '../pages/Welcome.vue';
import Replace from '../pages/Replace.vue';
import Login from '../pages/Login.vue';

const routes = [
    {
        path:"/login",
        name:"login",
        component: Login,
    },
    {
        path:"/admin",
        component:Header,
        children:[
            {
                path: "welcome",
                name:"welcome",
                component: Welcome,
            },
            {
                path: "replace",
                name: "replace",
                component: Replace
            }
        ]
    },
    {
        path: "/:catchAll(.*)",
        hidden: true,
        component: Login
    }
];

export default routes;