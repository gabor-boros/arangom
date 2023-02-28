package arangom

import (
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(args ...any) {
	m.Called(args)
}

func (m *MockLogger) Infof(format string, args ...any) {
	m.Called(format, args)
}

func (m *MockLogger) Error(args ...any) {
	m.Called(args)
}

func (m *MockLogger) Errorf(format string, args ...any) {
	m.Called(format, args)
}

func (m *MockLogger) Fatal(args ...any) {
	m.Called(args)
}

func (m *MockLogger) Fatalf(format string, args ...any) {
	m.Called(format, args)
}

type MockLogWriter struct {
	mock.Mock
}

func (m *MockLogWriter) WriteString(s string) (n int, err error) {
	args := m.Called(s)
	return args.Int(0), args.Error(1)
}

func TestNewDefaultLogger(t *testing.T) {
	t.Parallel()

	l := NewDefaultLogger()

	if l == nil {
		t.Errorf("Expected non-nil logger, got nil")
		return
	}

	if l.Writer != os.Stdout {
		t.Errorf("Expected writer to be os.Stdout, got %v", l.Writer)
	}

	if l.Exiter == nil {
		t.Errorf("Expected non-nil exiter, got nil")
	}
}

func TestDefaultLogger_Info(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[INFO] test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Info("test")
}

func TestDefaultLogger_Infof(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[INFO] formatted: test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Infof("formatted: %s", "test")
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[ERROR] test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Error("test")
}

func TestDefaultLogger_Errorf(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[ERROR] formatted: test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Errorf("formatted: %s", "test")
}

func TestDefaultLogger_Fatal(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[FATAL] test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
		Exiter: func(code int) {
			if code != 1 {
				t.Errorf("Expected exit code %d, got %d", 1, code)
			}
		},
	}

	l.Fatal("test")
}

func TestDefaultLogger_Fatalf(t *testing.T) {
	t.Parallel()

	w := new(MockLogWriter)
	w.On("WriteString", "[FATAL] formatted: test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
		Exiter: func(code int) {
			if code != 1 {
				t.Errorf("Expected exit code %d, got %d", 1, code)
			}
		},
	}

	l.Fatalf("formatted: %s", "test")
}
