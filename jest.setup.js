global.fetch = require('node-fetch');

import { config } from '@vue/test-utils';
config.global.stubs['font-awesome-icon'] = true;

