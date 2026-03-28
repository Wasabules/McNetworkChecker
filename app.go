package main

import (
	"McNetworkChecker/checker"
	"context"
	"sync"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	cancelDiag context.CancelFunc

	stepMu      sync.Mutex
	cancelStep  context.CancelFunc
	stepSkipped bool
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetSystemLocale detects the OS language and returns a 2-letter code.
func (a *App) GetSystemLocale() string {
	return checker.DetectOSLocale()
}

// SetLocale changes the language for backend logs and reports.
func (a *App) SetLocale(code string) {
	checker.SetLocale(code)
}

func (a *App) RunDiagnostic(address string) {
	a.stepMu.Lock()
	if a.cancelDiag != nil {
		a.cancelDiag()
	}
	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelDiag = cancel
	a.cancelStep = nil
	a.stepSkipped = false
	a.stepMu.Unlock()
	go a.runDiagnostic(ctx, address)
}

func (a *App) StopDiagnostic() {
	a.stepMu.Lock()
	defer a.stepMu.Unlock()
	if a.cancelStep != nil {
		a.cancelStep()
	}
	if a.cancelDiag != nil {
		a.cancelDiag()
		a.cancelDiag = nil
	}
}

func (a *App) SkipStep() {
	a.stepMu.Lock()
	defer a.stepMu.Unlock()
	if a.cancelStep != nil {
		a.stepSkipped = true
		a.cancelStep()
	}
}

func (a *App) emitStep(step, status string, result interface{}) {
	wailsRuntime.EventsEmit(a.ctx, "check:update", map[string]interface{}{
		"step": step, "status": status, "result": result,
	})
}

func (a *App) emitSkipped(step, reason string) {
	wailsRuntime.EventsEmit(a.ctx, "check:update", map[string]interface{}{
		"step": step, "status": "skipped", "reason": reason,
	})
}

func (a *App) logFor(step string) checker.LogFunc {
	return func(line string) {
		wailsRuntime.EventsEmit(a.ctx, "check:log", map[string]interface{}{
			"step": step, "line": line,
		})
	}
}
