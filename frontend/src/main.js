import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

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

const app = createApp(App);
app.component("font-awesome-icon", FontAwesomeIcon);

// Initialize the store from localStorage before mounting the app
store.dispatch("initializeStore");

app.use(store).use(router).mount("#app");
