package arangom

import (
	"os"
	"testing"

	"go.uber.org/mock/gomock"
)

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

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[INFO] test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Info("test")
}

func TestDefaultLogger_Infof(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[INFO] formatted: test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Infof("formatted: %s", "test")
}

func TestDefaultLogger_Error(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[ERROR] test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Error("test")
}

func TestDefaultLogger_Errorf(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[ERROR] formatted: test\n").Return(0, nil)

	l := &DefaultLogger{
		Writer: w,
	}

	l.Errorf("formatted: %s", "test")
}

func TestDefaultLogger_Fatal(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[FATAL] test\n").Return(0, nil)

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

	ctrl := gomock.NewController(t)
	w := NewMockLogWriter(ctrl)
	w.EXPECT().WriteString("[FATAL] formatted: test\n").Return(0, nil)

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
