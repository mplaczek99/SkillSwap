module.exports = {
  moduleFileExtensions: ["js", "json", "vue"],
  transform: {
    "^.+\\.vue$": "vue-jest",
    "^.+\\.js$": "babel-jest",
  },
  // Transform these ESM packages so Jest can process them:
  transformIgnorePatterns: [
    "/node_modules/(?!(compromise|grad-school|efrt|suffix-thumb)/)",
  ],
  moduleNameMapper: {
    "^@/(.*)$": "<rootDir>/src/$1",
    // More specific SVG file mock pattern
    "\\.(svg)$": "<rootDir>/tests/mocks/svgMock.js",
    // Also mock CSS imports
    "\\.(css|less)$": "<rootDir>/tests/mocks/styleMock.js",
  },
  testEnvironment: "jsdom",
  setupFilesAfterEnv: ["<rootDir>/jest.setup.js"],
};
