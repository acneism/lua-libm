import pathlib
import re
import shutil
import subprocess
import tempfile

root = pathlib.Path(__file__).resolve().parents[1]
work = root / "work"
work.mkdir(exist_ok=True)

def remove_fma(text):
    before = "math.FMA(z, atanTablePolynomial(i, z), atanTable(i*7+1))"
    after = "float64(float64(z*atanTablePolynomial(i, z))+atanTable(i*7+1))"
    assert text.count(before) == 1
    return text.replace(before, after)

def break_reduction(text):
    before = "n = int(sum) & 3"
    assert text.count(before) == 1
    return text.replace(before, "n = (int(sum)+1) & 3")

def change_table(text):
    start = text.index("const sincosTableBits =")
    match = list(re.finditer(r"\\x([0-9a-f]{2})", text[start:]))[1024]
    offset = start + match.start(1)
    value = int(match[1], 16) ^ 1
    return text[:offset] + f"{value:02x}" + text[offset + 2:]

for name, file, mutate in [("fma", "atan.go", remove_fma), ("branred", "reduce.go", break_reduction), ("table", "tables.go", change_table)]:
    with tempfile.TemporaryDirectory(prefix="mutation-", dir=work) as directory:
        target = pathlib.Path(directory)
        for source in root.glob("*.go"):
            shutil.copy2(source, target / source.name)
        shutil.copy2(root / "go.mod", target / "go.mod")
        shutil.copytree(root / "testdata", target / "testdata")
        shutil.copytree(root / "internal/grid", target / "internal/grid")
        source = target / file
        source.write_text(mutate(source.read_text()))
        result = subprocess.run(["go", "test", ".", "-run", "^TestRedisGrid$", "-count=1"], cwd=target, stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True)
        (work / f"mutation-{name}.log").write_text(result.stdout)
        if result.returncode == 0 or "--- FAIL: TestRedisGrid" not in result.stdout:
            raise RuntimeError(f"{name}: mutation did not fail TestRedisGrid\n{result.stdout}")
        print(f"{name} -> TestRedisGrid: caught")
