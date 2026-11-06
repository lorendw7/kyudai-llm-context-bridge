"""Check local links, website structure, and common accidental secret patterns.

Uses only Python's standard library. This is not a comprehensive security audit.
Pass --history to also inspect every reachable Git blob for common credentials
and prohibited file types; matches are reported without printing their contents.
"""

from html.parser import HTMLParser
from pathlib import Path
import re
import subprocess
import sys
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[1]
SITE = ROOT / "website"
REPO_URL = "https://github.com/lorendw7/kyudai-llm-context-bridge/blob/main/"
ERRORS = []
SECRET = re.compile(
    r"(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{50,}"
    r"|sk-(?:proj-|ant-)?[A-Za-z0-9_-]{32,}|AKIA[A-Z0-9]{16}"
    r"|-----BEGIN (?:RSA |EC |OPENSSH )?PRIVATE KEY-----)"
)
FORBIDDEN = {".pem", ".key", ".db", ".sqlite", ".sqlite3", ".gguf", ".safetensors"}
PRIVATE_DIRS = {"private", "secrets", "receipts", "node_modules", ".venv"}


def git(*args):
    return subprocess.check_output(["git", *args], cwd=ROOT)


class Page(HTMLParser):
    def __init__(self, text):
        super().__init__()
        self.ids = set()
        self.duplicates = set()
        self.links = []
        self.h1 = 0
        self.lang = None
        self.has_title = False
        self.steps = []
        self.feed(text)

    def handle_starttag(self, tag, attrs):
        attrs = dict(attrs)
        if "id" in attrs:
            if attrs["id"] in self.ids:
                self.duplicates.add(attrs["id"])
            self.ids.add(attrs["id"])
        if tag == "html":
            self.lang = attrs.get("lang")
        if tag == "h1":
            self.h1 += 1
        if tag == "title":
            self.has_title = True
        if tag == "button" and "data-step" in attrs:
            self.steps.append(attrs["data-step"])
        if tag == "img" and "alt" not in attrs:
            ERRORS.append("An image is missing alt text")
        for key in ("href", "src"):
            if key in attrs:
                self.links.append(attrs[key])


def check_link(source, link):
    if link.startswith(REPO_URL):
        target = ROOT / unquote(urlsplit(link[len(REPO_URL):]).path)
        if not target.exists():
            ERRORS.append(f"{source.relative_to(ROOT)}: missing repository target {target.relative_to(ROOT)}")
        return
    parts = urlsplit(link)
    if parts.scheme or parts.netloc:
        return  # External availability is checked separately during publication.
    target = (source.parent / unquote(parts.path)).resolve() if parts.path else source
    if not target.is_relative_to(ROOT) or not target.exists():
        ERRORS.append(f"{source.relative_to(ROOT)}: broken local link {link}")
        return
    if source.is_relative_to(SITE) and not target.is_relative_to(SITE):
        ERRORS.append(f"{source.relative_to(ROOT)}: link escapes the deployed website: {link}")
    if parts.fragment and target.suffix == ".html":
        page = Page(target.read_text(encoding="utf-8-sig"))
        if unquote(parts.fragment) not in page.ids:
            ERRORS.append(f"{source.relative_to(ROOT)}: missing HTML anchor {link}")


def inspect_content(name, content):
    path = Path(name)
    if path.suffix.lower() in FORBIDDEN or PRIVATE_DIRS.intersection(path.parts):
        ERRORS.append(f"Prohibited file: {name}")
    if path.name == ".env" or (path.name.startswith(".env.") and path.name != ".env.example"):
        ERRORS.append(f"Environment file: {name}")
    if SECRET.search(content.decode("utf-8", errors="replace")):
        ERRORS.append(f"Potential credential in {name} (value suppressed)")


paths = sorted(set(git("ls-files", "--cached", "--others", "--exclude-standard", "-z").decode().split("\0")) - {""})
checked = 0
for name in paths:
    file = ROOT / name
    if not file.is_file():
        continue  # Files moved in the working tree before staging.
    content = file.read_bytes()
    inspect_content(name, content)
    checked += 1
    if file.suffix not in {".md", ".html"}:
        continue
    text = content.decode("utf-8-sig")
    if file.suffix == ".html":
        page = Page(text)
        if page.lang not in {"en", "ja"} or page.h1 != 1 or not page.has_title or page.duplicates:
            ERRORS.append(f"{name}: invalid language, heading, title, or duplicate IDs")
        if page.steps != ["0", "1", "2"]:
            ERRORS.append(f"{name}: expected three synthetic handoff stages")
        links = page.links
    else:
        # Repository Markdown uses inline links and a small HTML masthead.
        links = re.findall(r"\[[^\]]*\]\(([^\s)]+)\)", text) + Page(text).links
    for link in links:
        check_link(file, link)

if "--history" in sys.argv:
    count = 0
    for line in git("rev-list", "--objects", "--all").decode().splitlines():
        oid, _, name = line.partition(" ")
        if name and git("cat-file", "-t", oid).strip() == b"blob":
            inspect_content(name, git("cat-file", "blob", oid))
            count += 1
    print(f"Inspected {count} historical blobs for common credentials and prohibited file types.")

if ERRORS:
    print("\n".join(sorted(set(ERRORS))))
    sys.exit(1)
print(f"PASS: {checked} files; local documentation/site links, HTML structure, common credential and file checks.")
print("External URL availability, visual behavior, and research claims require separate review.")
