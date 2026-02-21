import { useTheme } from 'next-themes'

function Logo() {
  const { resolvedTheme } = useTheme()
  const src = resolvedTheme === 'dark' ? '/logo-white.png' : '/logo-black.png'
  return <img src={src} alt="lit" style={{ height: 30 }} />
}

export default {
  logo: <Logo />,
  project: {
    link: 'https://github.com/tracewayapp/lit'
  },
  docsRepositoryBase: 'https://github.com/tracewayapp/lit/tree/main/docs',
  useNextSeoProps() {
    return {
      titleTemplate: '%s – lit'
    }
  },
  head: (
    <>
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <meta property="og:title" content="lit" />
      <meta property="og:description" content="Lightweight Go library for simplified database operations" />
    </>
  ),
  primaryHue: 205,
  darkMode: true,
  nextThemes: {
    defaultTheme: 'dark',
    enableSystem: false,
  },
  sidebar: {
    defaultMenuCollapseLevel: 1,
    toggleButton: true
  },
  toc: {
    backToTop: true
  },
  editLink: {
    text: 'Edit this page on GitHub →'
  },
  feedback: {
    content: 'Question? Give us feedback →',
    labels: 'feedback'
  },
  footer: {
    text: (
      <span>
        MIT {new Date().getFullYear()} ©{' '}
        <a href="https://github.com/tracewayapp" target="_blank">
          Traceway
        </a>
      </span>
    )
  }
}
