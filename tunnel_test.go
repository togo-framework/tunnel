package tunnel

import (
	"context"
	"testing"
)

func TestNormalizeAddr(t *testing.T) {
	cases := map[string]string{
		"8080":           "127.0.0.1:8080",
		":8080":          "127.0.0.1:8080",
		"localhost:3000": "localhost:3000",
		"0.0.0.0:80":     "0.0.0.0:80",
		"":               "127.0.0.1:80",
	}
	for in, want := range cases {
		if got := NormalizeAddr(in); got != want {
			t.Errorf("NormalizeAddr(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPortOf(t *testing.T) {
	cases := map[string]string{
		"8080":           "8080",
		":3000":          "3000",
		"localhost:5173": "5173",
	}
	for in, want := range cases {
		if got := PortOf(in); got != want {
			t.Errorf("PortOf(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestLogDriverRegistered(t *testing.T) {
	found := false
	for _, n := range Drivers() {
		if n == "log" {
			found = true
		}
	}
	if !found {
		t.Fatal("log driver not registered")
	}
}

func TestBuildLogDriver(t *testing.T) {
	svc, err := Build("log", nil)
	if err != nil {
		t.Fatalf("Build(log): %v", err)
	}
	if svc.Driver() != "log" {
		t.Errorf("driver = %q, want log", svc.Driver())
	}
	url, err := svc.Start(context.Background(), "8080")
	if err != nil {
		t.Fatalf("Start: %v", err)
	}
	if url == "" {
		t.Error("empty url")
	}
	st, err := svc.Status(context.Background())
	if err != nil || !st.Running || st.Driver != "log" {
		t.Errorf("status = %+v err=%v", st, err)
	}
	if err := svc.Stop(context.Background()); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}

func TestBuildUnknownDriver(t *testing.T) {
	if _, err := Build("nope", nil); err == nil {
		t.Fatal("expected error for unknown driver")
	}
}
