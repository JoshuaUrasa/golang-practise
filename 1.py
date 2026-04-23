from pathlib import Path


def ask_yes_no(prompt: str) -> bool:
    while True:
        value = input(f"{prompt} (y/n): ").strip().lower()
        if value in ("y", "yes"):
            return True
        if value in ("n", "no"):
            return False
        print("Please enter y or n.")


def ensure_dirs(root: Path, folders: list[str]) -> None:
    for folder in folders:
        path = root / folder
        path.mkdir(parents=True, exist_ok=True)
        print(f"[+] {folder}/")


def create_go_project_structure() -> None:
    root = Path.cwd()
    project_name = root.name

    print(f"\nProject root: {project_name}\n")

    api_cmd_name = input("API cmd folder name [api]: ").strip().lower() or "api"
    include_worker = ask_yes_no("Include worker folder?")
    include_uploads = ask_yes_no("Include uploads folder?")
    include_deployments = ask_yes_no("Include deployments folder?")

    folders = [
        f"cmd/{api_cmd_name}",
        "internal/auth",
        "internal/user",
        "internal/feature",
        "internal/middleware",
        "internal/platform/config",
        "internal/platform/database",
        "internal/platform/logger",
        "internal/platform/server",
        "internal/platform/storage",
        "pkg",
        "api/openapi",
        "configs",
        "scripts",
        "test/integration",
        "test/testdata",
    ]

    if include_worker:
        folders.append("cmd/worker")

    if include_uploads:
        folders.append("uploads")

    if include_deployments:
        folders.extend([
            "deployments/docker",
            "deployments/kubernetes",
        ])

    print("Creating folders...\n")
    ensure_dirs(root, folders)

    print("\nDone.\n")
    print("Recommended next steps:")
    print("  1. Rename internal/feature to your real domain name")
    print("     examples: expense, booking, order, product")
    print("  2. Run: go mod init github.com/yourname/" + project_name)
    print("  3. Start with cmd/, internal/platform/, and internal/auth/")


if __name__ == "__main__":
    create_go_project_structure()