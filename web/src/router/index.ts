import { createRouter, createWebHistory, RouteRecordRaw } from "vue-router";
import {injectGuards} from "@/router/guards";
import Layouts from "@/layouts/index.vue"


export const constantRoutes: RouteRecordRaw[] = [
    {
        path: "/login",
        component: () => import("@/views/login/index.vue"),
    },
    {
      path: "/signup",
      component: () => import("@/views/signup/index.vue")
    },
    {
        path: "/upload",
        component: () => import("@/views/upload/index.vue")
    },
    {
        path: "/files",
        component: () => import("@/views/files/index.vue")
    },
    {
        path: "/d",
        component: () => import("@/views/download/index.vue")
    },
    {
        path: "/settings",
        component: () => import("@/layouts/Settings.vue"),
        redirect: "/settings/profile",
        children: [
            {
                path: "profile",
                component: () => import("@/views/settings/profile/index.vue")
            },
            {
                path: "security",
                component: () => import("@/views/settings/security/index.vue")
            },
        ],
    },
    {
        path: "/admin",
        component: () => import("@/layouts/Settings.vue"),
        children: [
            {
                path: "/admin",
                component: () => import("@/views/admin/general/index.vue")
            },
            {
                path: "users",
                component: () => import("@/views/admin/users/index.vue")
            },
        ],
    },
];

export const asyncRoutes: RouteRecordRaw[] = [];

const router = createRouter({
    history: createWebHistory(),
    scrollBehavior: (to, from, savedPosition) => {
        if (savedPosition) {
            return savedPosition
        } else {
            return { top: 0, left: 0 }
        }
    },
    routes: constantRoutes,
});

injectGuards(router)

export {router, router as default};

