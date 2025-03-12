// Centralized event bus implementation using mitt library
import mitt from "mitt";

// Create a singleton instance of the event bus
const emitter = mitt();

// Export the event bus singleton
export default emitter;
