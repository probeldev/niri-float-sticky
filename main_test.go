package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const testDaemonEnv = "NIRI_FLOAT_STICKY_TEST_DAEMON"

func TestMain(m *testing.M) {
	if os.Getenv(testDaemonEnv) == "1" {
		// Run the real entry point in a subprocess, without the test flags.
		os.Args = os.Args[:1]
		main()
		os.Exit(0)
	}
	os.Exit(m.Run())
}

func TestDaemonExitsOnStreamEnd(t *testing.T) {
	for _, tc := range []struct {
		name      string
		malformed bool
	}{
		{name: "EOF"},
		{name: "decode_error", malformed: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// Keep paths short enough for Unix sockets. Each subprocess gets
			// isolated sockets and must never touch the real desktop session.
			runtimeDir, err := os.MkdirTemp("", "nfs-test-")
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() {
				if err := os.RemoveAll(runtimeDir); err != nil {
					t.Error(err)
				}
			})

			socketPath := filepath.Join(runtimeDir, "niri.sock")
			listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: socketPath, Net: "unix"})
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { listener.Close() })
			if err := listener.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}

			executable, err := os.Executable()
			if err != nil {
				t.Fatal(err)
			}
			cmd := exec.Command(executable)
			cmd.Env = append(os.Environ(),
				testDaemonEnv+"=1",
				"NIRI_SOCKET="+socketPath,
				"XDG_RUNTIME_DIR="+runtimeDir,
			)
			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = &output
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			done := make(chan struct{})
			var waitErr error
			go func() {
				waitErr = cmd.Wait()
				close(done)
			}()
			t.Cleanup(func() {
				// A regression must not leave a busy-looping daemon behind.
				cmd.Process.Kill()
				<-done
				if t.Failed() {
					t.Logf("daemon output:\n%s", output.String())
				}
			})

			conn, err := listener.AcceptUnix()
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { conn.Close() })
			if err := conn.SetDeadline(time.Now().Add(5 * time.Second)); err != nil {
				t.Fatal(err)
			}
			var request string
			if err := json.NewDecoder(conn).Decode(&request); err != nil {
				t.Fatal(err)
			}
			if request != "EventStream" {
				t.Fatalf("request = %q, want EventStream", request)
			}

			if _, err := io.WriteString(conn, "{\"WorkspacesChanged\":{\"workspaces\":[]}}\n"); err != nil {
				t.Fatal(err)
			}
			select {
			case <-done:
				t.Fatalf("daemon exited while the event stream was still open: %v", waitErr)
			case <-time.After(100 * time.Millisecond):
			}

			if tc.malformed {
				// Keep the connection open: the decoder error itself must
				// close the event channel and terminate the daemon.
				if _, err := io.WriteString(conn, "not JSON\n"); err != nil {
					t.Fatal(err)
				}
			} else if err := conn.Close(); err != nil {
				t.Fatal(err)
			}

			select {
			case <-done:
				if waitErr != nil {
					t.Fatalf("daemon did not exit successfully: %v", waitErr)
				}
			case <-time.After(3 * time.Second):
				t.Fatal("daemon did not exit after the event stream ended")
			}
		})
	}
}
