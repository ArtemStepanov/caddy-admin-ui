#!/usr/bin/env python3
"""Validate task-generation prerequisites and locate the task template."""

from __future__ import annotations

import json
import sys
from pathlib import Path

from common import format_speckit_command, get_feature_paths, resolve_template


def available_docs(paths: object) -> list[str]:
    docs: list[str] = []
    for path, name in ((paths.research, "research.md"), (paths.data_model, "data-model.md")):
        if path.is_file():
            docs.append(name)
    if paths.contracts_dir.is_dir() and any(paths.contracts_dir.iterdir()):
        docs.append("contracts/")
    if paths.quickstart.is_file():
        docs.append("quickstart.md")
    return docs


def main(argv: list[str] | None = None) -> int:
    args = list(argv if argv is not None else sys.argv[1:])
    if any(arg in {"--help", "-h"} for arg in args):
        print("Usage: setup_tasks.py [--json]")
        return 0
    if any(arg != "--json" for arg in args):
        print(f"ERROR: Unknown option '{next(arg for arg in args if arg != '--json')}'", file=sys.stderr)
        return 1

    json_mode = "--json" in args
    try:
        paths = get_feature_paths(script_file=Path(__file__))
    except SystemExit:
        print("ERROR: Failed to resolve feature paths", file=sys.stderr)
        return 1
    if not paths.impl_plan.is_file():
        print(f"ERROR: plan.md not found in {paths.feature_dir}", file=sys.stderr)
        print(f"Run {format_speckit_command('plan', paths.repo_root)} first to create the implementation plan.", file=sys.stderr)
        return 1
    if not paths.feature_spec.is_file():
        print(f"ERROR: spec.md not found in {paths.feature_dir}", file=sys.stderr)
        print(f"Run {format_speckit_command('specify', paths.repo_root)} first to create the feature structure.", file=sys.stderr)
        return 1

    template = resolve_template("tasks-template", paths.repo_root)
    if not template:
        print(f"ERROR: Tasks template not found for repository root: {paths.repo_root}", file=sys.stderr)
        return 1

    result = {"FEATURE_DIR": str(paths.feature_dir), "AVAILABLE_DOCS": available_docs(paths), "TASKS_TEMPLATE": str(template)}
    if json_mode:
        print(json.dumps(result, ensure_ascii=False, separators=(",", ":")))
    else:
        print(f"FEATURE_DIR: {paths.feature_dir}")
        print(f"TASKS_TEMPLATE: {template}")
        print("AVAILABLE_DOCS:")
        for doc in result["AVAILABLE_DOCS"]:
            print(f"  ✓ {doc}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
