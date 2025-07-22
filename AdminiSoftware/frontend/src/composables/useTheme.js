
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
