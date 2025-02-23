// jest.setup.js
import { config } from '@vue/test-utils';

// Stub the Font Awesome component to avoid resolution warnings in tests.
config.global.stubs['font-awesome-icon'] = true;

