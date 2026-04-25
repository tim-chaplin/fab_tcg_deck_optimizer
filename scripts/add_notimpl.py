"""Phase-2 helper: add NotImplemented{} markers to existing card files that aren't fully
modelled. Drops the legacy "Simplification:" doc paragraph (now captured by the marker) and
inserts a per-printing `// not implemented: <quirk>` line plus the NotImplemented method
between GoAgain and Play.

Usage:

  python scripts/add_notimpl.py <relative_or_absolute_file_path> "<note>"

Idempotent: rewrites a stale `// not implemented:` line if one already sits above
NotImplemented; does nothing if the file already matches the desired shape.
"""
import os
import re
import sys


def transform(src, note):
    # 1. Strip the leading-comment "Simplification:" paragraph from the doc block. The paragraph
    #    starts at a line beginning "// Simplification:" and includes any continuation `//` lines
    #    that follow. The blank `//` separator immediately above is also removed.
    src = re.sub(
        r"//\n// Simplification:[^\n]*(?:\n//[^\n]*)*\n",
        "",
        src,
        count=1,
    )

    # 2. For each printing's struct, insert (or refresh) the `// not implemented:` line and the
    #    NotImplemented method between the GoAgain method and the Play method. Match GoAgain by
    #    its full func line; insert after that line and any trailing whitespace, before Play.
    pattern = re.compile(
        r"^(?P<indent>[ \t]*)func \((?P<recv>[A-Za-z0-9]+)\) GoAgain\(\) bool[^\n]*\n"
        r"(?P<between>(?:[ \t]*//[^\n]*\n)?"
        r"(?:[ \t]*func \([A-Za-z0-9 ]+\) NotImplemented\(\)[^\n]*\n)?"
        r"(?:[ \t]*func \([A-Za-z0-9 ]+\) NotSilverAgeLegal\(\)[^\n]*\n)?)"
        r"(?P<rest>[ \t]*func \([ \t]*[A-Za-z]?[ \t]*[A-Za-z0-9]+\) Play\()",
        re.MULTILINE,
    )

    def repl(m):
        indent = m.group("indent")
        recv = m.group("recv")
        # Strip any prior NotImplemented + stale `// not implemented:` line from `between`.
        between = m.group("between")
        between = re.sub(r"^[ \t]*// not implemented:[^\n]*\n", "", between, flags=re.MULTILINE)
        between = re.sub(r"^[ \t]*func \([A-Za-z0-9 ]+\) NotImplemented\(\)[^\n]*\n", "", between, flags=re.MULTILINE)
        # Build replacement: GoAgain line + (any preserved markers, e.g. NotSilverAgeLegal) +
        # not-implemented comment + NotImplemented method + the Play func line we matched.
        goagain_line = m.group(0)[: m.group(0).find("\n") + 1]
        notimpl = f"{indent}// not implemented: {note}\n{indent}func ({recv}) NotImplemented()             {{}}\n"
        return goagain_line + between + notimpl + m.group("rest")

    new = pattern.sub(repl, src)
    return new


def main():
    if len(sys.argv) != 3:
        print("usage: add_notimpl.py <file> <note>", file=sys.stderr)
        sys.exit(2)
    path = sys.argv[1]
    note = sys.argv[2]
    with open(path, encoding="utf-8") as f:
        src = f.read()
    new = transform(src, note)
    if new == src:
        print(f"no change: {path}")
        return
    with open(path, "w", encoding="utf-8", newline="\n") as f:
        f.write(new)
    print(f"rewrote: {path}")


if __name__ == "__main__":
    main()
