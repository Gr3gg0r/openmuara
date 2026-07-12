import '@testing-library/jest-dom/vitest';
import { cleanup, configure } from '@testing-library/preact';
import { afterEach } from 'vitest';

configure({ testIdAttribute: 'data-testid' });
afterEach(cleanup);
