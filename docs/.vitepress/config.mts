import { defineConfig, HeadConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({

  title: "GENG",
  description: "A tool to generate golang web project.",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Getting Started", link: "/getting-started" },
      { text: "Guide", link: "/how-to-use" },
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
          { text: "Generate a project", link: "/generate-project" },
          { text: "Generate a project using config", link: "/generate-project-config" },
          { text: "Start the project", link: "/start-project" },
          { text: "Project Folder structure", link: "/folder-structure" },
          { text: "Generate a module", link: "/generate-module" },
          {
            text: "Generate a infrastructure",
            link: "/generate-infrastructure",
          },
        ],
      },
      {
        text: "Task",
        collapsed: false,
        items: [{ text: "Roadmap", link: "/roadmap" }],
      },
    ],

    socialLinks: [{ icon: "github", link: "https://github.com/mukezhz/geng" }],
    search: {
      provider: "local",
    },
  },
  transformHead: ({ pageData }) => {
    const head: HeadConfig[] = [];
    const title = pageData.frontmatter.title || "GenG";
    const description =
      pageData.frontmatter.description ||
      "GENG - A tool to generate golang web project.";
    const ogImage =
      pageData.frontmatter.ogImage ||
      "https://github.com/mukezhz/geng/assets/43813670/da07d8cc-8896-4a13-9b31-099958e65cb4";

    head.push(["meta", { property: "og:title", content: title }]);
    head.push(["meta", { property: "og:description", content: description }]);
    head.push(["meta", { property: "og:image", content: ogImage }]);

    head.push(["meta", { name: "twitter:card", content: "geng" }]);
    head.push(["meta", { name: "twitter:image", content: ogImage }]);
    head.push(["meta", { name: "twitter:title", content: title }]);
    head.push(["meta", { name: "twitter:description", content: description }]);

    return head;
  },
});
