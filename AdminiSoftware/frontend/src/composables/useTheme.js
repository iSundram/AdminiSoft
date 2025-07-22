
import { ref, computed, watch } from 'vue'

const THEME_KEY = 'adminisoftware-theme'
const themes = {
  default: {
    name: 'Default',
    colors: {
      primary: '#3b82f6',
      secondary: '#6b7280',
      success: '#10b981',
      warning: '#f59e0b',
      error: '#ef4444',
      background: '#ffffff',
      surface: '#f9fafb'
    }
  },
  whm: {
    name: 'WHM Style',
    colors: {
      primary: '#2563eb',
      secondary: '#374151',
      success: '#059669',
      warning: '#d97706',
      error: '#dc2626',
      background: '#f3f4f6',
      surface: '#ffffff'
    }
  },
  cpanel: {
    name: 'cPanel Style',
    colors: {
      primary: '#059669',
      secondary: '#4b5563',
      success: '#10b981',
      warning: '#f59e0b',
      error: '#ef4444',
      background: '#f0f9ff',
      surface: '#ffffff'
    }
  },
  dark: {
    name: 'Dark Mode',
    colors: {
      primary: '#3b82f6',
      secondary: '#9ca3af',
      success: '#10b981',
      warning: '#f59e0b',
      error: '#ef4444',
      background: '#111827',
      surface: '#1f2937'
    }
  }
}

const currentTheme = ref(localStorage.getItem(THEME_KEY) || 'default')

export function useTheme() {
  const theme = computed(() => themes[currentTheme.value] || themes.default)
  
  const availableThemes = computed(() => 
    Object.keys(themes).map(key => ({
      key,
      ...themes[key]
    }))
  )

  function setTheme(themeName) {
    if (themes[themeName]) {
      currentTheme.value = themeName
      localStorage.setItem(THEME_KEY, themeName)
      applyTheme(themes[themeName])
    }
  }

  function applyTheme(themeConfig) {
    const root = document.documentElement
    
    Object.entries(themeConfig.colors).forEach(([key, value]) => {
      root.style.setProperty(`--color-${key}`, value)
    })

    // Apply theme-specific classes
    document.body.className = document.body.className
      .replace(/theme-\w+/g, '')
      .concat(` theme-${currentTheme.value}`)
  }

  function getThemeForPanel(panelType) {
    const panelThemes = {
      admin: 'whm',
      reseller: 'default',
      user: 'cpanel'
    }
    
    return panelThemes[panelType] || 'default'
  }

  function autoSetThemeForPanel(panelType) {
    const recommendedTheme = getThemeForPanel(panelType)
    setTheme(recommendedTheme)
  }

  // Watch for theme changes and apply them
  watch(currentTheme, (newTheme) => {
    applyTheme(themes[newTheme])
  })

  // Apply initial theme
  applyTheme(theme.value)

  return {
    currentTheme,
    theme,
    availableThemes,
    setTheme,
    getThemeForPanel,
    autoSetThemeForPanel
  }
}
import { ref, computed, watch } from 'vue'

const theme = ref(localStorage.getItem('theme') || 'default')

export function useTheme() {
  const currentTheme = computed(() => theme.value)
  
  const themes = [
    { value: 'default', name: 'Default', description: 'Clean and modern interface' },
    { value: 'cpanel', name: 'cPanel Style', description: 'Traditional cPanel look and feel' },
    { value: 'whm', name: 'WHM Style', description: 'Professional WHM-inspired design' },
    { value: 'dark', name: 'Dark Mode', description: 'Dark theme for low-light environments' },
  ]

  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('theme', newTheme)
    applyTheme(newTheme)
  }

  const applyTheme = (themeName) => {
    // Remove existing theme classes
    const body = document.body
    themes.forEach(t => {
      body.classList.remove(`theme-${t.value}`)
    })
    
    // Add new theme class
    body.classList.add(`theme-${themeName}`)
    
    // Update meta theme-color for mobile browsers
    updateMetaThemeColor(themeName)
  }

  const updateMetaThemeColor = (themeName) => {
    const themeColors = {
      default: '#3b82f6',
      cpanel: '#059669',
      whm: '#64748b',
      dark: '#1e293b',
    }
    
    let metaThemeColor = document.querySelector('meta[name=theme-color]')
    if (!metaThemeColor) {
      metaThemeColor = document.createElement('meta')
      metaThemeColor.setAttribute('name', 'theme-color')
      document.head.appendChild(metaThemeColor)
    }
    
    metaThemeColor.setAttribute('content', themeColors[themeName] || themeColors.default)
  }

  const initTheme = () => {
    applyTheme(theme.value)
  }

  const toggleDarkMode = () => {
    const newTheme = theme.value === 'dark' ? 'default' : 'dark'
    setTheme(newTheme)
  }

  const isDarkMode = computed(() => theme.value === 'dark')

  // Watch for system theme changes
  const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)')
  const handleSystemThemeChange = (e) => {
    if (!localStorage.getItem('theme')) {
      setTheme(e.matches ? 'dark' : 'default')
    }
  }

  mediaQuery.addEventListener('change', handleSystemThemeChange)

  // Initialize with system preference if no saved theme
  if (!localStorage.getItem('theme')) {
    const systemTheme = mediaQuery.matches ? 'dark' : 'default'
    setTheme(systemTheme)
  }

  return {
    currentTheme,
    themes,
    setTheme,
    initTheme,
    toggleDarkMode,
    isDarkMode,
  }
}
