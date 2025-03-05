import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import axios from "axios";

// Import Font Awesome core and icons
import { library } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

// Import solid icons
import {
  faUser,
  faSignOutAlt,
  faHome,
  faSearch,
  faComments,
  faEnvelope,
  faEye,
  faEyeSlash,
  faCode,
  faFilter,
  faChevronUp,
  faChevronDown,
  faExclamationCircle,
  faTimes,
  faHistory,
  faStar,
  faStarHalfAlt,
  faArrowRight,
  faBellSlash,
  faMusic,
  faUtensils,
  faPalette,
  faPenFancy,
  faGuitar,
  faLanguage,
  faCog,
  faVideo,
  faUpload,
  faFilm,
  faSpinner,
  faCheckCircle,
  faPlay,
  faPlayCircle,
  faDownload,
  faCoins,
  faCalendarAlt,
  faGraduationCap,
  faExchangeAlt,
  faComment,
  faCommentAlt,
  faPaperPlane,
  faSync,
  faArrowUp,
  faArrowDown,
  // New icons for job posting features
  faBriefcase,
  faBuilding,
  faMapMarkerAlt,
  faExternalLinkAlt,
  faBookmark,
  faMoneyBillAlt,
  faPlus,
  faPlusCircle,
  faCheck,
  faArrowLeft,
  faUndo,
} from "@fortawesome/free-solid-svg-icons";

// Add all solid icons to the library
library.add(
  faUser,
  faSignOutAlt,
  faHome,
  faSearch,
  faComments,
  faEnvelope,
  faEye,
  faEyeSlash,
  faCode,
  faFilter,
  faChevronUp,
  faChevronDown,
  faExclamationCircle,
  faTimes,
  faHistory,
  faStar,
  faStarHalfAlt,
  faArrowRight,
  faBellSlash,
  faMusic,
  faUtensils,
  faPalette,
  faPenFancy,
  faGuitar,
  faLanguage,
  faCog,
  faVideo,
  faUpload,
  faFilm,
  faSpinner,
  faCheckCircle,
  faPlay,
  faPlayCircle,
  faDownload,
  faCoins,
  faCalendarAlt,
  faGraduationCap,
  faExchangeAlt,
  faComment,
  faCommentAlt,
  faPaperPlane,
  faSync,
  faArrowUp,
  faArrowDown,
  // New icons
  faBriefcase,
  faBuilding,
  faMapMarkerAlt,
  faExternalLinkAlt,
  faBookmark,
  faMoneyBillAlt,
  faPlus,
  faPlusCircle,
  faCheck,
  faArrowLeft,
  faUndo,
);

// Set up axios interceptors to add the auth token to all requests
axios.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  },
);

// Handle 401 responses globally
axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (
      error.response &&
      error.response.status === 401 &&
      router.currentRoute.value.path !== "/login"
    ) {
      console.log("Session expired, redirecting to login");
      // Token is expired or invalid
      store.commit("logout");
      router.push("/login");
    }
    return Promise.reject(error);
  },
);

const app = createApp(App);
app.component("font-awesome-icon", FontAwesomeIcon);

// Initialize the store from localStorage before mounting the app
store.dispatch("initializeStore");

app.use(store).use(router).mount("#app");
