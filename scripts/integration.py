import argparse
import json
import pathlib
import shutil
import subprocess
import tempfile

parser = argparse.ArgumentParser()
parser.add_argument("--lua", required=True)
arguments = parser.parse_args()
source = pathlib.Path(arguments.lua).resolve()
libm = pathlib.Path(__file__).resolve().parents[1]

with tempfile.TemporaryDirectory(prefix="lua-libm-integration-") as directory:
    target = pathlib.Path(directory) / "lua"
    shutil.copytree(source, target, ignore=shutil.ignore_patterns(".git", "work", ".integration"))
    corpus = target / "testdata/snippets.json"
    snippets = json.loads(corpus.read_text())
    precise = json.loads((target / "testdata/libm.json").read_text())
    corpus.write_text(json.dumps(snippets + precise, ensure_ascii=False))
    test = target / "corpus_test.go"
    content = test.read_text()
    assert content.count("NewState(Options{HookEvery:") == 1
    content = content.replace("import (", 'import (\nlibm "github.com/acneism/lua-libm"', 1)
    content = content.replace("NewState(Options{HookEvery:", "NewState(Options{Math: &Math{Sin: libm.Sin, Cos: libm.Cos, Tan: libm.Tan, Asin: libm.Asin, Acos: libm.Acos, Atan: libm.Atan, Atan2: libm.Atan2}, HookEvery:", 1)
    test.write_text(content)
    subprocess.run(["go", "mod", "edit", "-require=github.com/acneism/lua-libm@v0.0.0", f"-replace=github.com/acneism/lua-libm={libm}"], cwd=target, check=True)
    subprocess.run(["go", "test", ".", "-run", "^TestCorpus$", "-count=1"], cwd=target, check=True)
    print(f"integration: {len(snippets)} corpus snippets and {len(precise)} exact libm snippets")
