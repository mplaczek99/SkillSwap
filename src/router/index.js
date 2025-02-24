import { createRouter, createWebHistory } from "vue-router";

// Lazy-load components for performance optimization
const Dashboard = () => import("@/components/Dashboard.vue");
const LoginForm = () => import("@/components/LoginForm.vue");
const RegisterForm = () => import("@/components/RegisterForm.vue");
const Profile = () => import("@/components/Profile.vue");
const Search = () => import("@/components/Search.vue");
const Chat = () => import("@/components/Chat.vue");

const routes = [
  { path: "/", name: "Dashboard", component: Dashboard },
  { path: "/login", name: "Login", component: LoginForm },
  { path: "/register", name: "Register", component: RegisterForm },
  { path: "/profile", name: "Profile", component: Profile },
  { path: "/search", name: "Search", component: Search },
  { path: "/chat", name: "Chat", component: Chat },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard for protected routes
router.beforeEach((to, from, next) => {
  const publicPages = ["/login", "/register"];
  const authRequired = !publicPages.includes(to.path);
  const loggedIn = localStorage.getItem("token");

  if (authRequired && !loggedIn) {
    next("/login");
  } else {
    next();
  }
});

export default router;

