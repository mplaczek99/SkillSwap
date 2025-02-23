module.exports = {
  root: true,
  env: {
    node: true,
    jest: true
  },
  extends: [
    'plugin:vue/essential',
    'eslint:recommended',
    'prettier' // Use the base Prettier config instead of @vue/prettier
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

