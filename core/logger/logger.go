package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

var (
	mu        sync.Mutex
	baseDir   string
	logFile   *os.File
	LoggerErr *log.Logger
)

// Init inicializa el sistema de logs.
//
// appName: nombre para la carpeta (ej: "GoFinder").
// Si falla devuelve error.
func Init(appName string) error {
	mu.Lock()
	defer mu.Unlock()

	cfgDir, err := os.UserConfigDir()
	if err != nil || cfgDir == "" {
		home, herr := os.UserHomeDir()
		if herr != nil {
			return fmt.Errorf("no se pudo obtener carpeta config/home: %v, %v", err, herr)
		}
		cfgDir = filepath.Join(home, ".config")
	}

	baseDir = filepath.Join(cfgDir, appName)
	logDir := filepath.Join(baseDir, "logs")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		return fmt.Errorf("mkdir logs: %w", err)
	}

	logPath := filepath.Join(logDir, "errors.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	logFile = f

	// Solo errores → stderr + fichero
	multi := io.MultiWriter(os.Stderr, logFile)
	LoggerErr = log.New(multi, "[ERROR] ", log.LstdFlags|log.Lmicroseconds)

	logSystemInfo()

	return nil
}

// Close cierra el fichero de logs (llamar al salir)
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	if logFile != nil {
		if err := logFile.Sync(); err != nil {
			_ = logFile.Close()
			return err
		}
		return logFile.Close()
	}
	return nil
}

// CatchPanic captura un panic y lo registra (usar con defer en main y en goroutines)
func CatchPanic() {
	if r := recover(); r != nil {
		msg := fmt.Sprintf("PANIC: %v\n\nSTACK:\n%s", r, debug.Stack())
		if LoggerErr != nil {
			LoggerErr.Println(msg)
		} else {
			log.Println(msg)
		}
		if err := WriteCrashFile(fmt.Sprintf("%v", r), debug.Stack()); err != nil && LoggerErr != nil {
			LoggerErr.Printf("error escribiendo crash file: %v", err)
		}
	}
}

// GoSafe ejecuta una goroutine protegiéndola con recover y registro
func GoSafe(fn func()) {
	go func() {
		defer CatchPanic()
		fn()
	}()
}

// writeCrashFile guarda un archivo separado con stack y metadata
func WriteCrashFile(reason string, stack []byte) error {
	var f *os.File
	var err error

	// Prefer writing under configured baseDir (user config dir) when available.
	if baseDir != "" {
		crashDir := filepath.Join(baseDir, "crashes")
		if err := os.MkdirAll(crashDir, 0o700); err != nil {
			return err
		}
		name := fmt.Sprintf("crash-%s.log", time.Now().Format("20060102-150405"))
		path := filepath.Join(crashDir, name)

		// Create the file exclusively with restrictive permissions to avoid TOCTOU in public dirs.
		f, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
		if err != nil {
			return err
		}
	} else {
		// If logging system not initialized, create a secure temporary file using CreateTemp
		// which generates an unpredictable filename in the system temp dir.
		f, err = os.CreateTemp("", "gofinder-crash-*.log")
		if err != nil {
			return err
		}
		// Try to restrict permissions where supported.
		_ = f.Chmod(0o600)
	}
	defer f.Close()

	u, _ := user.Current()
	host, _ := os.Hostname()
	_, _ = fmt.Fprintf(f, "Time: %s\nUser: %s\nHost: %s\nPID: %d\nOS: %s %s\nGo: %s\n\nReason: %s\n\nStack:\n%s\n",
		time.Now().Format(time.RFC3339), u.Username, host, os.Getpid(), runtime.GOOS, runtime.GOARCH, runtime.Version(), reason, stack)

	return nil
}

func logSystemInfo() {
	u, _ := user.Current()
	host, _ := os.Hostname()
	LoggerErr.Printf("SYSINFO user=%s host=%s pid=%d go=%s os=%s arch=%s",
		u.Username, host, os.Getpid(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
