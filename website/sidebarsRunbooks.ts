import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  runbooksSidebar: [
    {
      type: 'doc',
      id: 'index',
      label: 'Runbooks',
    },
    {
      type: 'doc',
      id: 'local-development',
      label: 'Local Development',
    },
    {
      type: 'doc',
      id: 'testing',
      label: 'Testing',
    },
    {
      type: 'doc',
      id: 'debugging',
      label: 'Debugging',
    },
    {
      type: 'doc',
      id: 'quality-gates',
      label: 'Quality Gates',
    },
    {
      type: 'doc',
      id: 'on-call',
      label: 'On-Call',
    },
    {
      type: 'doc',
      id: 'stripe-fpx-card',
      label: 'Stripe FPX / Card',
    },
  ],
};

export default sidebars;
