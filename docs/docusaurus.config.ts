import { themes as prismThemes } from "prism-react-renderer";
import type { Config } from "@docusaurus/types";
import type * as Preset from "@docusaurus/preset-classic";

const isDev = process.env.NODE_ENV === "dev";

const config: Config = {
  title: "Backupman",
  tagline: "Compact open-source solution for database backups",
  favicon: "img/favicon.ico",

  future: {
    v4: true,
  },

  url: "https://herytz.github.io",
  baseUrl: isDev ? "/" : "/backupman",

  organizationName: "herytz",
  projectName: "backupman",

  onBrokenLinks: "throw",
  onBrokenMarkdownLinks: "warn",

  i18n: {
    defaultLocale: "en",
    locales: ["en"],
  },

  presets: [
    [
      "classic",
      {
        docs: {
          sidebarPath: "./sidebars.ts",
        },
        blog: {
          showReadingTime: true,
          feedOptions: {
            type: ["rss", "atom"],
            xslt: true,
          },
          onInlineTags: "warn",
          onInlineAuthors: "warn",
          onUntruncatedBlogPosts: "warn",
        },
        theme: {
          customCss: "./src/css/custom.css",
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: "img/logo-banner.jpg",
    navbar: {
      title: "Backupman",
      logo: {
        alt: "Backupman Logo",
        src: "img/logo.png",
      },
      items: [
        {
          type: "docSidebar",
          sidebarId: "docsSidebar",
          position: "left",
          label: "Docs",
        },
        {
          type: "docSidebar",
          sidebarId: "referencesSidebar",
          position: "left",
          label: "References",
        },
        {
          href: "https://github.com/herytz/backupman",
          label: "GitHub",
          position: "right",
        },
      ],
    },
    footer: {
      style: "dark",
      links: [
        {
          title: "Docs",
          items: [
            {
              label: "Quickstart",
              to: "/docs/quickstart",
            },
          ],
        },
        {
          title: "More",
          items: [
            {
              label: "GitHub",
              href: "https://github.com/herytz/backupman",
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Backupman. Made with ❤️ by Hery Nirintsoa`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
    },
  } satisfies Preset.ThemeConfig,

  markdown: {
    mermaid: true,
  },

  themes: ["@docusaurus/theme-mermaid"],
};

export default config;
