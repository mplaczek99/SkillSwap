import { createApp } from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";

// Import Font Awesome core and icons
import { library } from "@fortawesome/fontawesome-svg-core";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
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
} from "@fortawesome/free-solid-svg-icons";

// Add all icons to the library
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
);

const app = createApp(App);
app.component("font-awesome-icon", FontAwesomeIcon);
app.use(store).use(router).mount("#app");
