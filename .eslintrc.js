module.exports = {
  root: true,
  env: {
    node: true,
    jest: true
  },
  extends: [
    'plugin:vue/essential',
    'eslint:recommended',
    'prettier'
  ],
  rules: {
    'vue/multi-word-component-names': 'off',
    'no-useless-catch': 'off'
  },
  overrides: [
    {
      files: ['tests/**/*.spec.js'],
      env: {
        jest: true
      }
    }
  ]
};

