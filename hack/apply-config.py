#!/usr/bin/env python3
import argparse
import pathlib
import re
import subprocess

ADMIN_CONFIG_PATTERN = r'^admin_.*\.yaml$'
ADMIN_CONFIG = re.compile(ADMIN_CONFIG_PATTERN)

CONFIG_PATTERN = r'^cfg_.*\.yaml$'
CONFIG = re.compile(CONFIG_PATTERN)


def _is_admin_config(file_path):
    if not file_path.is_file():
        return False

    return ADMIN_CONFIG.match(file_path.name) is not None


def _is_standard_config(file_path):
    if not file_path.is_file():
        return False

    return CONFIG.match(file_path.name) is not None


def apply_config(config_path, user=None, dry=True):
    oc_apply = ["oc", "apply", "-f", str(config_path)]

    if dry:
        oc_apply.append("--dry-run")

    if user is not None:
        oc_apply.extend(["--as", user])

    p = subprocess.run(oc_apply, capture_output=True)
    try:
        p.check_returncode()
        print(f"{' '.join(oc_apply)}: OK")
    except subprocess.CalledProcessError:
        print(f"[ERROR] {' '.join(oc_apply)}: Failed")
        print(f"=== OUTPUT ===\n{p.stdout.decode('utf-8')}\n=== STDERR ===\n{p.stderr.decode('utf-8')}")
        raise


def apply_dir(directory, level, user, dry=True):
    for subdir in sorted([x for x in directory.iterdir() if x.is_dir()]):
        apply_dir(subdir, level, user)

    admin_files = sorted([x for x in directory.iterdir() if _is_admin_config(x)])
    standard_files = sorted([x for x in directory.iterdir() if _is_standard_config(x)])

    failures = False

    if level in ("admin", "all"):
        print(f"Applying admin config in: {directory}")
        for f in admin_files:
            try:
                apply_config(f, user, dry)
            except subprocess.CalledProcessError:
                # Attempt to continue with the rest
                failures = True

    if level in ("standard", "all"):
        print(f"Applying config in: {directory}")
        for f in standard_files:
            try:
                apply_config(f, user, dry)
            except subprocess.CalledProcessError:
                # Attempt to continue with the rest
                failures = True

    if failures:
        raise Exception("[ERROR] Some config failed to apply")


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--confirm", dest="confirm", action="store_true", default=False)
    parser.add_argument("--level", choices=("standard", "admin", "all"), default="standard")
    parser.add_argument("--as", dest="user")
    parser.add_argument("directory", nargs=1, type=pathlib.Path)
    args = parser.parse_args()

    try:
        apply_dir(args.directory[0], args.level, args.user, not args.confirm)
    except BaseException as err:
        print(err)
        return 1

    print("Success!")
    return 1


if __name__ == "__main__":
    main()
