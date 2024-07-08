const routes = [
  {
    path: '/', component: () => import('pages/Login.vue')
  },
  {
    path: "/", component: () => import("layouts/MainLayout.vue"),
    children: [
      { path: "", component: () => import("pages/Index.vue") },
      { path: "/dashboard", component: () => import("pages/user/Index.vue"), meta: { requiredLogin: true } },
      { path: "/redfish", component: () => import("pages/agents/Redfish.vue"), meta: { requiredLogin: true } },
      { path: "/vsphere", component: () => import("pages/agents/Vsphere.vue"), meta: { requiredLogin: true } },
    ]
  },
  {
    path: "/:catchAll(.*)*",
    component: () => import("pages/Error404.vue"),
  },
];

export default routes;
