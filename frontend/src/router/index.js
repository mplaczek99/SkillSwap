import { createRouter, createWebHistory } from "vue-router";

// Lazy-load components for performance optimization.
const Dashboard = () => import("@/components/Dashboard.vue");
const LoginForm = () => import("@/components/LoginForm.vue");
const RegisterForm = () => import("@/components/RegisterForm.vue");
const Profile = () => import("@/components/Profile.vue");
const Search = () => import("@/components/Search.vue");
const Chat = () => import("@/components/Chat.vue");
const Schedule = () => import("@/components/Schedule.vue");
const VideoUpload = () => import("@/components/VideoUpload.vue");
const VideosList = () => import("@/components/VideosList.vue");
const Transactions = () => import("@/components/Transactions.vue");
const FeedbackSystem = () => import("@/components/FeedbackSystem.vue");

// New job posting components
const JobPostings = () => import("@/components/JobPostings.vue");
const JobDetail = () => import("@/components/JobDetail.vue");
const PostJob = () => import("@/components/PostJob.vue");

const routes = [
  { path: "/", name: "Dashboard", component: Dashboard },
  { path: "/login", name: "Login", component: LoginForm },
  { path: "/register", name: "Register", component: RegisterForm },
  { path: "/profile", name: "Profile", component: Profile },
  { path: "/search", name: "Search", component: Search },
  { path: "/chat", name: "Chat", component: Chat },
  { path: "/schedule", name: "Schedule", component: Schedule },
  { path: "/upload-video", name: "VideoUpload", component: VideoUpload },
  { path: "/videos", name: "VideosList", component: VideosList },
  { path: "/transactions", name: "Transactions", component: Transactions },
  { path: "/feedback", name: "Feedback", component: FeedbackSystem },

  // Job posting routes
  { path: "/jobs", name: "JobPostings", component: JobPostings },
  { path: "/jobs/:id", name: "JobDetail", component: JobDetail },
  { path: "/post-job", name: "PostJob", component: PostJob },
  { path: "/edit-job/:id", name: "EditJob", component: PostJob },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  // Scroll to top on route navigation
  scrollBehavior() {
    return { top: 0 };
  },
});

const publicPages = ["/login", "/register"];
router.beforeEach((to, _, next) => {
  // More precise check: either exact match or followed by a slash
  const authRequired = !publicPages.some(
    (page) => to.path === page || to.path.startsWith(page + "/"),
  );
  const loggedIn = localStorage.getItem("token");
  if (authRequired && !loggedIn) {
    next("/login");
  } else {
    next();
  }
});

export default router;
