// Package browser provides web navigation capabilities using Playwright
package browser

import (
	"fmt"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

// Manager handles browser instances and tabs
type Manager struct {
	pw       *playwright.Playwright
	browser  playwright.Browser
	pages    map[string]playwright.Page // taskID -> page
	pagesMux sync.RWMutex
	options  *Options
}

// Options configures the browser manager
type Options struct {
	Headless       bool
	DefaultTimeout time.Duration
	UserAgent      string
	ViewportWidth  int
	ViewportHeight int
}

// DefaultOptions returns sensible defaults
func DefaultOptions() *Options {
	return &Options{
		Headless:       true,
		DefaultTimeout: 60 * time.Second, // Increased for complex pages
		UserAgent:      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
		ViewportWidth:  1280,
		ViewportHeight: 720,
	}
}

// NewManager creates a new browser manager
func NewManager(opts *Options) (*Manager, error) {
	if opts == nil {
		opts = DefaultOptions()
	}

	// Install Playwright browsers if needed
	err := playwright.Install(&playwright.RunOptions{
		Verbose: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to install playwright: %w", err)
	}

	// Start Playwright
	pw, err := playwright.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to start playwright: %w", err)
	}

	// Launch browser
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(opts.Headless),
	})
	if err != nil {
		pw.Stop()
		return nil, fmt.Errorf("failed to launch browser: %w", err)
	}

	return &Manager{
		pw:      pw,
		browser: browser,
		pages:   make(map[string]playwright.Page),
		options: opts,
	}, nil
}

// GetOrCreatePage gets an existing page for a task or creates a new one
func (m *Manager) GetOrCreatePage(taskID string) (playwright.Page, error) {
	m.pagesMux.Lock()
	defer m.pagesMux.Unlock()

	// Check if page exists
	if page, exists := m.pages[taskID]; exists {
		return page, nil
	}

	// Create new page
	page, err := m.browser.NewPage(playwright.BrowserNewPageOptions{
		UserAgent: playwright.String(m.options.UserAgent),
		Viewport: &playwright.Size{
			Width:  m.options.ViewportWidth,
			Height: m.options.ViewportHeight,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page: %w", err)
	}

	// Set default timeout
	page.SetDefaultTimeout(float64(m.options.DefaultTimeout.Milliseconds()))

	m.pages[taskID] = page
	return page, nil
}

// Navigate navigates to a URL in a task-specific page
func (m *Manager) Navigate(taskID, url string) (*NavigationResult, error) {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return nil, err
	}

	// Navigate
	response, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
	})
	if err != nil {
		return nil, fmt.Errorf("navigation failed: %w", err)
	}

	// Extract page info
	title, _ := page.Title()
	content, _ := page.Content()
	currentURL := page.URL()

	return &NavigationResult{
		URL:        currentURL,
		Title:      title,
		StatusCode: response.Status(),
		Content:    content,
	}, nil
}

// Click clicks an element by selector
func (m *Manager) Click(taskID, selector string) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	return page.Click(selector)
}

// ClickCoordinates clicks at specific x,y coordinates
func (m *Manager) ClickCoordinates(taskID string, x, y float64) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	return page.Mouse().Click(x, y)
}

// GetPageElements returns visible interactive elements with their info
func (m *Manager) GetPageElements(taskID string) ([]ElementInfo, error) {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return nil, err
	}

	// JavaScript to get all interactive elements with their bounds
	script := `
		Array.from(document.querySelectorAll('a, button, input, select, textarea, [onclick], [role="button"]'))
			.filter(el => {
				const rect = el.getBoundingClientRect();
				const style = window.getComputedStyle(el);
				return rect.width > 0 && rect.height > 0 &&
				       style.visibility !== 'hidden' &&
				       style.display !== 'none' &&
				       rect.top < window.innerHeight &&
				       rect.bottom > 0;
			})
			.map((el, idx) => {
				const rect = el.getBoundingClientRect();
				return {
					index: idx,
					tag: el.tagName.toLowerCase(),
					text: el.innerText?.substring(0, 50) || '',
					id: el.id || '',
					class: el.className || '',
					x: Math.round(rect.left + rect.width / 2),
					y: Math.round(rect.top + rect.height / 2),
					width: Math.round(rect.width),
					height: Math.round(rect.height)
				};
			})
			.slice(0, 50);  // Limit to first 50 elements
	`

	result, err := page.Evaluate(script)
	if err != nil {
		return nil, err
	}

	// Convert to ElementInfo slice
	var elements []ElementInfo
	if arr, ok := result.([]interface{}); ok {
		for _, item := range arr {
			if elem, ok := item.(map[string]interface{}); ok {
				info := ElementInfo{
					Index:  int(elem["index"].(float64)),
					Tag:    elem["tag"].(string),
					Text:   elem["text"].(string),
					ID:     elem["id"].(string),
					Class:  elem["class"].(string),
					X:      int(elem["x"].(float64)),
					Y:      int(elem["y"].(float64)),
					Width:  int(elem["width"].(float64)),
					Height: int(elem["height"].(float64)),
				}
				elements = append(elements, info)
			}
		}
	}

	return elements, nil
}

// Type types text into an element
func (m *Manager) Type(taskID, selector, text string) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	return page.Fill(selector, text)
}

// Screenshot takes a screenshot of the current page
func (m *Manager) Screenshot(taskID, path string) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path: playwright.String(path),
		Type: playwright.ScreenshotTypePng,
	})
	return err
}

// ScreenshotOptimized takes an AI-optimized screenshot (lower resolution, compressed)
func (m *Manager) ScreenshotOptimized(taskID, path string) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	// Set viewport to lower resolution for AI processing
	if err := page.SetViewportSize(800, 600); err != nil {
		return err
	}

	// Take screenshot
	_, err = page.Screenshot(playwright.PageScreenshotOptions{
		Path:    playwright.String(path),
		Type:    playwright.ScreenshotTypeJpeg,
		Quality: playwright.Int(60), // Compressed quality
	})

	// Restore original viewport
	page.SetViewportSize(m.options.ViewportWidth, m.options.ViewportHeight)

	return err
}

// GetText extracts text from an element
func (m *Manager) GetText(taskID, selector string) (string, error) {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return "", err
	}

	element, err := page.QuerySelector(selector)
	if err != nil || element == nil {
		return "", fmt.Errorf("element not found: %s", selector)
	}

	return element.TextContent()
}

// WaitForSelector waits for an element to appear
func (m *Manager) WaitForSelector(taskID, selector string, timeout time.Duration) error {
	page, err := m.GetOrCreatePage(taskID)
	if err != nil {
		return err
	}

	_, err = page.WaitForSelector(selector, playwright.PageWaitForSelectorOptions{
		Timeout: playwright.Float(float64(timeout.Milliseconds())),
	})
	return err
}

// ClosePage closes a specific task page
func (m *Manager) ClosePage(taskID string) error {
	m.pagesMux.Lock()
	defer m.pagesMux.Unlock()

	page, exists := m.pages[taskID]
	if !exists {
		return nil // Already closed
	}

	if err := page.Close(); err != nil {
		return err
	}

	delete(m.pages, taskID)
	return nil
}

// CloseAllPages closes all open pages
func (m *Manager) CloseAllPages() error {
	m.pagesMux.Lock()
	defer m.pagesMux.Unlock()

	for taskID, page := range m.pages {
		page.Close()
		delete(m.pages, taskID)
	}
	return nil
}

// Close shuts down the browser and playwright
func (m *Manager) Close() error {
	m.CloseAllPages()

	if m.browser != nil {
		m.browser.Close()
	}

	if m.pw != nil {
		return m.pw.Stop()
	}

	return nil
}

// GetActiveTasks returns list of active task IDs
func (m *Manager) GetActiveTasks() []string {
	m.pagesMux.RLock()
	defer m.pagesMux.RUnlock()

	tasks := make([]string, 0, len(m.pages))
	for taskID := range m.pages {
		tasks = append(tasks, taskID)
	}
	return tasks
}

// NavigationResult contains the result of a navigation
type NavigationResult struct {
	URL        string
	Title      string
	StatusCode int
	Content    string
}

// ElementInfo contains information about a page element
type ElementInfo struct {
	Index  int    `json:"index"`
	Tag    string `json:"tag"`
	Text   string `json:"text"`
	ID     string `json:"id"`
	Class  string `json:"class"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}
