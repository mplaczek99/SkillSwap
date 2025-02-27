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
  faArrowRight,
  faBellSlash,
  faMusic,
  faUtensils,
  faPalette,
  faPenFancy,
  faGuitar,
  faLanguage,
  faCog,
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
  faArrowRight,
  faBellSlash,
  faMusic,
  faUtensils,
  faPalette,
  faPenFancy,
  faGuitar,
  faLanguage,
  faCog,
);

const app = createApp(App);
app.component("font-awesome-icon", FontAwesomeIcon);
app.use(store).use(router).mount("#app");
