"""Explicit local startup verification using synthetic configuration only.

Run with Python's standard library. Builds into a temporary directory, requires
127.0.0.1:8080 to be free, and stops only the child process it starts.
"""

import json
import os
from pathlib import Path
import socket
import subprocess
import tempfile
import time
import urllib.error
import urllib.request

ROOT = Path(__file__).resolve().parents[1]
TOKEN = "synthetic-startup-token"
FLAGS = subprocess.CREATE_NO_WINDOW if os.name == "nt" else 0


def assert_port_free():
    with socket.socket() as probe:
        probe.bind(("127.0.0.1", 8080))


def main():
    assert_port_free()
    env = {key: value for key, value in os.environ.items()
           if not key.upper().startswith("BRIDGE_")}
    env.update(GOTOOLCHAIN="local", GOPROXY="off")
    fixture = dict(env, BRIDGE_TOKEN=TOKEN, BRIDGE_LAB_ID="synthetic-lab",
                   BRIDGE_USER_ID="synthetic-user",
                   BRIDGE_MODEL_IDS="synthetic-model-b,synthetic-model-a")
    with tempfile.TemporaryDirectory(prefix="bridge-startup-") as temporary:
        binary = Path(temporary) / ("bridge.exe" if os.name == "nt" else "bridge")
        subprocess.run(["go", "build", "-o", str(binary), "./cmd/bridge"],
                       cwd=ROOT, env=env, check=True, timeout=120,
                       creationflags=FLAGS)
        for key in ("BRIDGE_TOKEN", "BRIDGE_LAB_ID", "BRIDGE_USER_ID", "BRIDGE_MODEL_IDS"):
            invalid = dict(fixture)
            invalid.pop(key)
            result = subprocess.run([str(binary)], cwd=ROOT, env=invalid,
                                    capture_output=True, text=True, timeout=10,
                                    creationflags=FLAGS)
            assert result.returncode != 0, f"missing {key} did not fail"
            output = result.stdout + result.stderr
            assert "Starting Bridge" not in output, "invalid config reached startup"
            assert TOKEN not in output, "startup error exposed token"
            assert_port_free()
        print("PASS: four missing configuration values fail before listening")

        opener = urllib.request.build_opener(urllib.request.ProxyHandler({}))

        def request(path, method="GET", data=None, authenticated=True, expected=200):
            headers = {"Content-Type": "application/json"}
            if authenticated:
                headers["Authorization"] = "Bearer " + TOKEN
            payload = data if isinstance(data, bytes) else (
                json.dumps(data).encode() if data is not None else None)
            req = urllib.request.Request("http://127.0.0.1:8080" + path,
                                         data=payload, headers=headers, method=method)
            try:
                response = opener.open(req, timeout=3)
            except urllib.error.HTTPError as error:
                response = error
            with response:
                body = response.read()
                assert response.status == expected, (path, response.status, expected)
                assert response.headers.get_content_type() == "application/json"
                return json.loads(body), response.headers

        with open(Path(temporary) / "server.log", "w+", encoding="utf-8") as log:
            child = subprocess.Popen([str(binary)], cwd=ROOT, env=fixture,
                                     stdout=log, stderr=log, creationflags=FLAGS)
            try:
                deadline = time.monotonic() + 10
                while True:
                    assert child.poll() is None, "server exited before readiness"
                    try:
                        health, _ = request("/health", authenticated=False)
                        assert health == {"status": "ok", "service": "context-bridge"}
                        break
                    except urllib.error.URLError:
                        if time.monotonic() >= deadline:
                            raise AssertionError("server readiness timed out")
                        time.sleep(0.1)
                models, _ = request("/v1/models")
                assert models == {"models": ["synthetic-model-a", "synthetic-model-b"]}
                listing, _ = request("/v1/sessions")
                assert listing == {"sessions": []}
                for path, method in (("/v1/models", "GET"), ("/v1/sessions", "GET"),
                                     ("/v1/sessions", "POST"), ("/v1/sessions/missing", "GET")):
                    error, headers = request(path, method, authenticated=False, expected=401)
                    assert error == {"code": "unauthorized"}
                    assert headers["WWW-Authenticate"] == "Bearer"
                created = []
                for _ in range(2):
                    item, headers = request("/v1/sessions", "POST",
                                            {"schema_version": "1", "policy_id": "full-history-v1"},
                                            expected=201)
                    assert len(item["id"]) == 32
                    assert headers["Location"] == "/v1/sessions/" + item["id"]
                    fetched, _ = request(headers["Location"])
                    assert fetched == item
                    created.append(item)
                assert created[0]["id"] != created[1]["id"]
                for payload, status in (
                    ({"schema_version": "1", "policy_id": "full-history-v1", "user_id": "forged"}, 422),
                    (b"{}{}", 422),
                    (b" " * 4097, 413),
                ):
                    request("/v1/sessions", "POST", payload, expected=status)
                missing, _ = request("/v1/sessions/synthetic-missing", expected=404)
                assert missing == {"code": "session_not_found"}
                listing, _ = request("/v1/sessions")
                assert sorted(listing["sessions"], key=lambda x: x["id"]) == sorted(created, key=lambda x: x["id"])
                print("PASS: real localhost health/models/auth/create/get/list/invalid-input/size-limit checks")
            finally:
                if child.poll() is None:
                    child.terminate()
                    try:
                        child.wait(timeout=5)
                    except subprocess.TimeoutExpired:
                        child.kill()
                        child.wait(timeout=5)
                log.flush()
                log.seek(0)
                assert TOKEN not in log.read(), "server log exposed token"
        with socket.socket() as probe:
            probe.settimeout(1)
            assert probe.connect_ex(("127.0.0.1", 8080)) != 0, "listener remains after cleanup"
        print("PASS: child stopped, listener closed, temporary build/logs removed on exit")


if __name__ == "__main__":
    main()
