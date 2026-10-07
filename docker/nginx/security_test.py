#!/usr/bin/env python3
"""Exercise the real proxy image in an isolated network, without touching the stack."""

import argparse
import json
import subprocess
import tempfile
import time
import urllib.error
import urllib.request
import uuid
from pathlib import Path


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def docker(*args, check=True):
    return subprocess.run(
        ["docker", *args], check=check, capture_output=True, text=True
    )


def response(port, path, headers=None):
    opener = urllib.request.build_opener(urllib.request.ProxyHandler({}), NoRedirect())
    request = urllib.request.Request(f"http://127.0.0.1:{port}{path}", headers=headers or {})
    try:
        result = opener.open(request, timeout=3)
    except urllib.error.HTTPError as error:
        result = error
    with result:
        return result.status, result.headers, result.read().decode()


def wait_ready(name, port):
    for _ in range(50):
        try:
            if response(port, "/api/echo")[0] == 200:
                return
        except (OSError, urllib.error.URLError):
            pass
        time.sleep(0.1)
    raise AssertionError(f"proxy failed to start: {docker('logs', name).stdout}")


def check_headers(port):
    for path, expected in [("/", 200), ("/api/echo", 200), ("/api/private", 401),
                           ("/assets/missing.js", 404), ("/api/docs", 200)]:
        status, headers, _ = response(port, path)
        assert status == expected, (path, status)
        for field, value in [("X-Content-Type-Options", "nosniff"),
                             ("Referrer-Policy", "strict-origin-when-cross-origin"),
                             ("X-Frame-Options", "DENY")]:
            assert headers.get_all(field) == [value], (path, field, headers.get_all(field))
        assert headers.get("X-Powered-By") is None
        assert headers.get("Strict-Transport-Security") is None
        if path.startswith("/api/"):
            assert headers.get("Cache-Control") == "no-store"
        if path == "/api/docs":
            assert headers.get_all("Content-Security-Policy") == ["default-src 'none'"]
        if path == "/":
            assert "frame-ancestors 'none'" in headers["Content-Security-Policy"]
    status, headers, _ = response(port, "/api", {"Host": "public.example:443",
                                               "X-Forwarded-Proto": "https"})
    assert status == 308 and headers["Location"] == "/api/", headers


def check_forwarding(port, trusted):
    headers = {"X-Forwarded-For": "192.0.2.10, 198.51.100.20",
               "X-Forwarded-Proto": "https", "Forwarded": "for=192.0.2.99;proto=https",
               "CF-Connecting-IP": "192.0.2.99", "X-Forwarded-Host": "evil.example",
               "X-Forwarded-Port": "443"}
    first = json.loads(response(port, "/api/echo", headers)[2])
    assert first["scheme"] == ("https" if trusted else "http"), first
    assert (first["ip"] == "198.51.100.20") == trusted, first
    assert first["ip"] == first["real_ip"], first
    assert all(first[field] == "" for field in ["forwarded", "cf_ip", "host", "port"]), first
    headers["X-Forwarded-For"] = "198.51.100.21"
    second = json.loads(response(port, "/api/echo", headers)[2])
    assert (first["ip"] != second["ip"]) == trusted
    headers["X-Forwarded-Proto"] = "https,http"
    assert json.loads(response(port, "/api/echo", headers)[2])["scheme"] == "http"
    # Both the API and frontend proxy paths must normalize forwarding fields.
    headers["X-Forwarded-Proto"] = "https"
    assert json.loads(response(port, "/echo", headers)[2])["scheme"] == ("https" if trusted else "http")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--image", default="ent-go-vue-nginx")
    args = parser.parse_args()
    prefix = "ent-proxy-security-" + uuid.uuid4().hex[:10]
    network, fixture = prefix + "-net", prefix + "-fixture"
    containers = []
    with tempfile.TemporaryDirectory(prefix=prefix) as directory:
        config = Path(directory) / "fixture.conf"
        config.write_text('''server {
    listen 8080;
    add_header X-Content-Type-Options nosniff always;
    add_header Referrer-Policy strict-origin-when-cross-origin always;
    add_header X-Frame-Options SAMEORIGIN always;
    add_header X-Powered-By fixture always;
    add_header Cache-Control no-store always;
    location = /api/docs {
        add_header Content-Security-Policy "default-src 'none'" always;
        add_header Cache-Control no-store always;
        return 200 '<html>docs</html>';
    }
    location = /api/private { return 401; }
    location = /assets/missing.js { return 404; }
    location ~ ^/(api/)?echo$ {
        default_type application/json;
        return 200 '{"ip":"$http_x_forwarded_for","real_ip":"$http_x_real_ip","scheme":"$http_x_forwarded_proto","forwarded":"$http_forwarded","cf_ip":"$http_cf_connecting_ip","host":"$http_x_forwarded_host","port":"$http_x_forwarded_port"}';
    }
    location / { return 200 '<html>fixture</html>'; }
}
''')
        docker("network", "create", network)
        try:
            gateway = json.loads(docker("network", "inspect", network).stdout)[0]["IPAM"]["Config"][0]["Gateway"]
            containers.append(fixture)
            docker("run", "-d", "--name", fixture, "--network", network,
                   "--network-alias", "backend", "--network-alias", "frontend",
                   "--mount", f"type=bind,source={config},target=/etc/nginx/conf.d/default.conf,readonly",
                   "--entrypoint", "nginx", args.image, "-g", "daemon off;")
            for trusted in [False, True]:
                name = prefix + ("-trusted" if trusted else "-untrusted")
                containers.append(name)
                docker("run", "-d", "--name", name, "--network", network,
                       "-p", "127.0.0.1::8080", "-e",
                       f"TRUSTED_EDGE_CIDRS={gateway}/32, ::1/128" if trusted else "TRUSTED_EDGE_CIDRS=",
                       args.image)
                port = docker("port", name, "8080/tcp").stdout.strip().rsplit(":", 1)[1]
                wait_ready(name, port)
                check_headers(port)
                check_forwarding(port, trusted)
                if trusted:
                    # A different container is outside the trusted gateway /32.
                    result = docker("exec", fixture, "wget", "-qO-",
                                    "--header=X-Forwarded-For: 198.51.100.20",
                                    "--header=X-Forwarded-Proto: https", f"http://{name}:8080/api/echo")
                    data = json.loads(result.stdout)
                    assert data["scheme"] == "http" and data["ip"] != "198.51.100.20", data
                    # Proxy-generated errors also retain the common headers.
                    request = urllib.request.Request(f"http://127.0.0.1:{port}/api/echo",
                                                     data=b"x" * 17000, method="POST")
                    try:
                        urllib.request.build_opener(urllib.request.ProxyHandler({})).open(request, timeout=3)
                        raise AssertionError("oversized request accepted")
                    except urllib.error.HTTPError as error:
                        assert error.code == 413
                        assert error.headers.get_all("X-Content-Type-Options") == ["nosniff"]
                        error.close()
            for invalid in ["0.0.0.0/0", "::/0", "127.0.0.1/32; include /tmp/evil;", "999.1.1.1/32"]:
                result = docker("run", "--rm", "-e", f"TRUSTED_EDGE_CIDRS={invalid}",
                                args.image, "nginx", "-t", check=False)
                assert result.returncode != 0, "unsafe proxy trust accepted"
            print("PASS: relative redirects, header uniqueness/CSP/cache/errors, trusted and untrusted forwarding, distinct visitor IPs, invalid trust rejection")
        finally:
            for name in reversed(containers):
                docker("rm", "-f", name, check=False)
            docker("network", "rm", network, check=False)


if __name__ == "__main__":
    main()
