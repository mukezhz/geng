import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "GENG",
  description: "A tool to generate golang web project.",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Getting Started", link: "/getting-started" },
    ],

    sidebar: [
      {
        text: "Introduction",
        collapsed: false,
        items: [
          { text: "What is Geng?", link: "/introduction" },
          { text: "Getting Started.", link: "/getting-started" },
        ],
      },
      {
        text: "Guide",
        collapsed: false,
        items: [
          { text: "How to use?", link: "/how-to-use" },
          // { text: "Getting Started", link: "/getting-started" },
        ],
      },
      {
        text: "Task",
        collapsed: false,
        items: [
          { text: "Roadmap", link: "/roadmap" },
        ],
      },
    ],

    socialLinks: [
      { icon: "github", link: "https://github.com/mukezhz/geng" },
    ],
    search: {
      provider: "local",
    },
  },
});
