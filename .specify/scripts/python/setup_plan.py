#!/usr/bin/env python3
"""Create a feature plan from the configured template."""

from __future__ import annotations

import json
import sys
from pathlib import Path

from common import get_feature_paths, resolve_template


def main(argv: list[str] | None = None) -> int:
    args = list(argv if argv is not None else sys.argv[1:])
    if any(arg in {"--help", "-h"} for arg in args):
        print("Usage: setup_plan.py [--json]")
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

    paths.feature_dir.mkdir(parents=True, exist_ok=True)
    if paths.impl_plan.is_file():
        message = f"Plan already exists at {paths.impl_plan}, skipping template copy"
    else:
        template = resolve_template("plan-template", paths.repo_root)
        if template:
            paths.impl_plan.write_text(template.read_text(encoding="utf-8"), encoding="utf-8")
            message = f"Copied plan template to {paths.impl_plan}"
        else:
            paths.impl_plan.touch()
            message = "Warning: Plan template not found"
    print(message, file=sys.stderr if json_mode else sys.stdout)

    result = {
        "FEATURE_SPEC": str(paths.feature_spec),
        "IMPL_PLAN": str(paths.impl_plan),
        "SPECS_DIR": str(paths.feature_dir),
        "BRANCH": paths.current_branch,
    }
    if json_mode:
        print(json.dumps(result, ensure_ascii=False, separators=(",", ":")))
    else:
        for key, value in result.items():
            print(f"{key}: {value}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
