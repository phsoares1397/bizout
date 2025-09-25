import { createRouter, createWebHistory } from "vue-router"
import Login from "../views/Login.vue"
import Register from "../views/Register.vue"
import ForgotPassword from "../views/ForgotPassword.vue"
import ProfileMenu from "../views/ProfileMenu_Sidebar.vue"
import UserImage from "../views/UserImage.vue"
import MeusDados from "../views/MeusDados.vue"
import Config from "../views/Config.vue"

const routes = [
    { path: "/", redirect: "/login" },
    { path: "/login", name: "Login", component: Login },
    { path: "/register", name: "Register", component: Register },
    { path: "/forgot-password", name: "ForgotPassword", component: ForgotPassword },
    { path: "/perfil", name: "Perfil", component: ProfileMenu },
    { path: "/image", name: "UserImage", component: UserImage },
    { path: "/meus-dados", name: "MeusDados", component: MeusDados },
    { path: "/config", name: "Configurações", component: Config },
]

const router = createRouter({
    history: createWebHistory('/conta/'),
    routes,
})

export default router
