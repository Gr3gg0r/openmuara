import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

const sidebars: SidebarsConfig = {
  docsSidebar: [
    'intro',
    'quickstart',
    'architecture',
    {
      type: 'category',
      label: 'Deployment',
      items: ['operations', 'hosted-testing', 'security'],
    },
    {
      type: 'category',
      label: 'Providers',
      link: {type: 'doc', id: 'providers'},
      items: [
        'providers/default',
        'providers/fawry',
        'providers/stripe',
        'providers/billplz',
        'providers/toyyibpay',
        'providers/senangpay',
        'providers/ipay88',
      ],
    },
    {
      type: 'category',
      label: 'Reference',
      items: [
        'cli',
        'errors',
        'webhooks',
        'provider-contract',
        'mkp-billing-requirements',
      ],
    },
    {
      type: 'category',
      label: 'Migration',
      items: [
        'migration/openmuara-to-openmuara',
        'migration/provider-manifests',
      ],
    },
    {
      type: 'category',
      label: 'Community',
      items: ['contributing', 'contributing-providers'],
    },
    {
      type: 'link',
      label: 'Runbooks',
      href: '/runbooks',
    },
  ],
};

export default sidebars;
