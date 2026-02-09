// Package editor provides multi-file edit tracking and rollback.
package editor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Session represents an editing session
type Session struct {
	ID          string        `json:"id"`
	StartedAt   time.Time     `json:"started_at"`
	Description string        `json:"description"`
	Edits       []EditRecord  `json:"edits"`
	Status      SessionStatus `json:"status"`
	Checkpoint  string        `json:"checkpoint,omitempty"` // Git commit hash
}

// SessionStatus represents the session state
type SessionStatus string

const (
	StatusActive     SessionStatus = "active"
	StatusCompleted  SessionStatus = "completed"
	StatusRolledBack SessionStatus = "rolled_back"
)

// EditRecord represents a single file edit
type EditRecord struct {
	ID          string    `json:"id"`
	FilePath    string    `json:"file_path"`
	Operation   string    `json:"operation"` // "create", "modify", "delete"
	OldContent  string    `json:"old_content,omitempty"`
	NewContent  string    `json:"new_content,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Applied     bool      `json:"applied"`
	Description string    `json:"description"`
}

// Manager manages edit sessions
type Manager struct {
	projectRoot string
	storagePath string
	sessions    map[string]*Session
	current     *Session
}

// NewManager creates a new edit manager
func NewManager(projectRoot, storagePath string) *Manager {
	absRoot, _ := filepath.Abs(projectRoot)
	return &Manager{
		projectRoot: absRoot,
		storagePath: storagePath,
		sessions:    make(map[string]*Session),
	}
}

func (m *Manager) validatePath(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	rel, err := filepath.Rel(m.projectRoot, absPath)
	if err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("path escapes project root: %s", path)
	}
	return nil
}

// StartSession begins a new editing session
func (m *Manager) StartSession(description string) *Session {
	session := &Session{
		ID:          fmt.Sprintf("session_%d", time.Now().UnixNano()),
		StartedAt:   time.Now(),
		Description: description,
		Edits:       make([]EditRecord, 0),
		Status:      StatusActive,
	}

	m.sessions[session.ID] = session
	m.current = session

	return session
}

// CurrentSession returns the current active session
func (m *Manager) CurrentSession() *Session {
	return m.current
}

// RecordEdit records a file edit without applying it
func (m *Manager) RecordEdit(filePath, operation, oldContent, newContent, description string) *EditRecord {
	if m.current == nil {
		m.StartSession("Auto-session")
	}

	edit := EditRecord{
		ID:          fmt.Sprintf("edit_%d", time.Now().UnixNano()),
		FilePath:    filePath,
		Operation:   operation,
		OldContent:  oldContent,
		NewContent:  newContent,
		Timestamp:   time.Now(),
		Applied:     false,
		Description: description,
	}

	m.current.Edits = append(m.current.Edits, edit)
	return &edit
}

// ApplyEdit applies a single edit
func (m *Manager) ApplyEdit(editID string) error {
	if m.current == nil {
		return fmt.Errorf("no active session")
	}

	for i, edit := range m.current.Edits {
		if edit.ID == editID && !edit.Applied {
			if err := m.applyFileEdit(&edit); err != nil {
				return err
			}
			m.current.Edits[i].Applied = true
			return nil
		}
	}

	return fmt.Errorf("edit not found or already applied: %s", editID)
}

// ApplyAll applies all pending edits in the session
func (m *Manager) ApplyAll() error {
	if m.current == nil {
		return fmt.Errorf("no active session")
	}

	for i, edit := range m.current.Edits {
		if !edit.Applied {
			if err := m.applyFileEdit(&edit); err != nil {
				return fmt.Errorf("failed to apply edit %s: %w", edit.ID, err)
			}
			m.current.Edits[i].Applied = true
		}
	}

	return nil
}

// RollbackEdit rolls back a single applied edit
func (m *Manager) RollbackEdit(editID string) error {
	if m.current == nil {
		return fmt.Errorf("no active session")
	}

	for i, edit := range m.current.Edits {
		if edit.ID == editID && edit.Applied {
			if err := m.rollbackFileEdit(&edit); err != nil {
				return err
			}
			m.current.Edits[i].Applied = false
			return nil
		}
	}

	return fmt.Errorf("edit not found or not applied: %s", editID)
}

// RollbackAll rolls back all applied edits in reverse order
func (m *Manager) RollbackAll() error {
	if m.current == nil {
		return fmt.Errorf("no active session")
	}

	// Rollback in reverse order
	for i := len(m.current.Edits) - 1; i >= 0; i-- {
		edit := &m.current.Edits[i]
		if edit.Applied {
			if err := m.rollbackFileEdit(edit); err != nil {
				return fmt.Errorf("failed to rollback edit %s: %w", edit.ID, err)
			}
			edit.Applied = false
		}
	}

	m.current.Status = StatusRolledBack
	return nil
}

// CompleteSession marks the session as completed
func (m *Manager) CompleteSession() error {
	if m.current == nil {
		return fmt.Errorf("no active session")
	}

	// Apply all pending edits
	if err := m.ApplyAll(); err != nil {
		return err
	}

	m.current.Status = StatusCompleted

	// Save session
	if err := m.SaveSession(m.current.ID); err != nil {
		return err
	}

	m.current = nil
	return nil
}

// GetPendingEdits returns all unapplied edits
func (m *Manager) GetPendingEdits() []EditRecord {
	if m.current == nil {
		return nil
	}

	var pending []EditRecord
	for _, edit := range m.current.Edits {
		if !edit.Applied {
			pending = append(pending, edit)
		}
	}
	return pending
}

// GetDiff returns a diff-style summary of all edits
func (m *Manager) GetDiff() string {
	if m.current == nil {
		return ""
	}

	var result string
	for _, edit := range m.current.Edits {
		status := "[ ]"
		if edit.Applied {
			status = "[x]"
		}
		result += fmt.Sprintf("%s %s: %s - %s\n",
			status, edit.Operation, edit.FilePath, edit.Description)
	}
	return result
}

// SaveSession saves session to disk
func (m *Manager) SaveSession(sessionID string) error {
	session, exists := m.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found: %s", sessionID)
	}

	// Create directory
	dir := filepath.Join(m.storagePath, "sessions")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Save as JSON
	filePath := filepath.Join(dir, sessionID+".json")
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}

// LoadSession loads a session from disk
func (m *Manager) LoadSession(sessionID string) (*Session, error) {
	filePath := filepath.Join(m.storagePath, "sessions", sessionID+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	m.sessions[session.ID] = &session
	return &session, nil
}

// applyFileEdit applies a file edit to disk
func (m *Manager) applyFileEdit(edit *EditRecord) error {
	if err := m.validatePath(edit.FilePath); err != nil {
		return err
	}

	switch edit.Operation {
	case "create":
		dir := filepath.Dir(edit.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return os.WriteFile(edit.FilePath, []byte(edit.NewContent), 0644)

	case "modify":
		return os.WriteFile(edit.FilePath, []byte(edit.NewContent), 0644)

	case "delete":
		return os.Remove(edit.FilePath)

	default:
		return fmt.Errorf("unknown operation: %s", edit.Operation)
	}
}

// rollbackFileEdit reverses a file edit
func (m *Manager) rollbackFileEdit(edit *EditRecord) error {
	if err := m.validatePath(edit.FilePath); err != nil {
		return err
	}

	switch edit.Operation {
	case "create":
		// Rollback create = delete
		return os.Remove(edit.FilePath)

	case "modify":
		// Rollback modify = restore old content
		return os.WriteFile(edit.FilePath, []byte(edit.OldContent), 0644)

	case "delete":
		// Rollback delete = restore file
		dir := filepath.Dir(edit.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		return os.WriteFile(edit.FilePath, []byte(edit.OldContent), 0644)

	default:
		return fmt.Errorf("unknown operation: %s", edit.Operation)
	}
}
