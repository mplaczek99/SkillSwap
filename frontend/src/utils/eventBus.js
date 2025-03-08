// Centralized event bus implementation using mitt library
import mitt from "mitt";

const emitter = mitt();

// Export singleton instance of the event bus
export default emitter;

// Add backward compatibility layer for components still using $root.$emit
export function setupRootCompatibility(app) {
  // Create a global property that proxies to the event bus
  app.config.globalProperties.$root = {
    $emit: (event, ...args) => emitter.emit(event, ...args),
    $on: (event, callback) => {
      emitter.on(event, callback);
      // Return a function to remove the event listener
      return () => emitter.off(event, callback);
    },
    $off: (event, callback) => emitter.off(event, callback),
  };
}
