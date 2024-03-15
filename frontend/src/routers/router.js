import Header from '../layouts/Header.vue';
import Welcome from '../pages/Welcome.vue';
import ImageToText from '../pages/ImageToText.vue';
import LLM from '../pages/LLM.vue';
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
                path: "image-to-text",
                name: "image-to-text",
                component: ImageToText
            },
            {
                path: "llm",
                name: "llm",
                component: LLM
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