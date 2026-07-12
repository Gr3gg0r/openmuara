import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'OpenMuara',
  tagline: 'Local-first payment virtualization for developers',
  favicon: 'img/logo.svg',

  future: {
    v4: true,
  },

  url: 'https://gr3gg0r.github.io',
  baseUrl: '/openmuara/',

  organizationName: 'Gr3gg0r',
  projectName: 'openmuara',

  onBrokenLinks: 'throw',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          path: '../docs',
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/Gr3gg0r/openmuara/tree/main/',
          routeBasePath: 'docs',
          exclude: [
            'cli-schemas/**',
          ],
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  plugins: [
    [
      '@docusaurus/plugin-content-docs',
      {
        id: 'runbooks',
        path: '../runbooks',
        routeBasePath: 'runbooks',
        sidebarPath: './sidebarsRunbooks.ts',
        editUrl: 'https://github.com/Gr3gg0r/openmuara/tree/main/',
      } satisfies import('@docusaurus/plugin-content-docs').Options,
    ],
    [
      '@easyops-cn/docusaurus-search-local',
      {
        hashed: true,
        language: ['en'],
        indexDocs: true,
        indexBlog: false,
        docsRouteBasePath: ['/docs', '/runbooks'],
      },
    ],
  ],

  themeConfig: {
    image: 'img/logo.svg',
    colorMode: {
      respectPrefersColorScheme: true,
    },
    navbar: {
      title: 'OpenMuara',
      logo: {
        alt: 'OpenMuara Logo',
        src: 'img/logo.svg',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          to: '/docs/providers',
          label: 'Providers',
          position: 'left',
        },
        {
          to: '/runbooks',
          label: 'Runbooks',
          position: 'left',
        },
        {
          href: 'https://github.com/Gr3gg0r/openmuara',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {label: 'Quickstart', to: '/docs/quickstart'},
            {label: 'Architecture', to: '/docs/architecture'},
            {label: 'Hosted Testing', to: '/docs/hosted-testing'},
          ],
        },
        {
          title: 'Reference',
          items: [
            {label: 'Security', to: '/docs/security'},
            {label: 'Operations', to: '/docs/operations'},
            {label: 'OpenAPI', href: 'pathname:///openapi.yaml'},
          ],
        },
        {
          title: 'More',
          items: [
            {label: 'GitHub', href: 'https://github.com/Gr3gg0r/openmuara'},
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} OpenMuara. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
