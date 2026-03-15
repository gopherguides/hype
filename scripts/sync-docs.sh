#!/usr/bin/env bash
set -euo pipefail

# sync-docs.sh — Sync documentation from the hype repo to hypemd.dev
#
# Usage: ./scripts/sync-docs.sh <output-dir>
#
# This script processes docs from the hype repo and outputs them
# as hype blog content (with frontmatter) to the specified directory.

OUTPUT_DIR="${1:?Usage: sync-docs.sh <output-dir>}"
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Doc-to-slug mapping: source_file|slug|title|needs_hype_render
DOCS=(
  "docs/installation.md|installation|Installation|no"
  "docs/cli-reference.md|cli-reference|CLI Reference|no"
  "docs/quick-reference.md|quick-reference|Quick Reference|no"
  "docs/quickstart/hype.md|quickstart|Quick Start Guide|yes"
  "docs/html-export.md|html-export|HTML Export|no"
  "docs/preview.md|preview|Live Preview|no"
  "docs/mermaid.md|mermaid|Mermaid Diagrams|no"
  "docs/marked.md|marked|Marked 2 Integration|no"
  "docs/slides.md|slides|Slides & Presentations|no"
  "docs/blog/hype.md|blog|Blog Generator|yes"
)

PUBLISHED="03/15/2026"
AUTHOR="Gopher Guides"

wrap_frontmatter() {
  local title="$1"
  local slug="$2"
  local tags="$3"
  local seo="$4"
  local content="$5"

  cat <<EOF
# ${title}

<details>
slug: docs/${slug}
published: ${PUBLISHED}
author: ${AUTHOR}
seo_description: ${seo}
tags: ${tags}
</details>

${content}
EOF
}

generate_seo() {
  local title="$1"
  echo "Hype documentation — ${title}. Learn how to use this feature in the Hype dynamic Markdown engine."
}

generate_tags() {
  local slug="$1"
  echo "docs, ${slug}, hype"
}

for entry in "${DOCS[@]}"; do
  IFS='|' read -r src slug title needs_render <<< "$entry"

  src_path="${REPO_ROOT}/${src}"
  if [[ ! -f "$src_path" ]]; then
    echo "WARNING: Source file not found: ${src_path}, skipping" >&2
    continue
  fi

  # Use flat directory names (docs-slug) since hype blog only discovers one level deep
  flat_dir="docs-${slug}"
  dest_dir="${OUTPUT_DIR}/${flat_dir}"
  mkdir -p "$dest_dir"

  if [[ "$needs_render" == "yes" ]]; then
    src_dir="$(dirname "$src_path")"
    if rendered=$(cd "$src_dir" && hype export -format markdown -f "$(basename "$src_path")" 2>/dev/null); then
      content="$rendered"
    else
      echo "NOTE: hype export failed for ${src}, using raw content" >&2
      content=$(cat "$src_path")
    fi
  else
    content=$(cat "$src_path")
  fi

  # Fix image paths: convert repo-root-relative paths to be relative to the doc's source dir.
  # For example, docs/blog/hype.md references "docs/blog/images/..." but images are at "images/..."
  src_dir_rel="$(dirname "$src")"
  content=$(echo "$content" | sed "s|(${src_dir_rel}/|(|g")

  # Copy any images directory from the source doc's location to the output
  src_dir_abs="${REPO_ROOT}/$(dirname "$src")"
  if [[ -d "${src_dir_abs}/images" ]]; then
    cp -r "${src_dir_abs}/images" "${dest_dir}/images"
  fi

  # Escape ALL hype-specific tags so the blog builder won't try to process them.
  # Hype's parser processes tags even inside markdown code fences, so we must
  # escape them everywhere. We use HTML entities for the angle brackets.
  content=$(echo "$content" | sed -E \
    -e 's/<(code src=)/\&lt;\1/g' \
    -e 's/<(go (src|doc|run))/\&lt;\1/g' \
    -e 's/<(cmd exec)/\&lt;\1/g' \
    -e 's/<(include src)/\&lt;\1/g' \
    -e 's/<\/(code|go|cmd|include)>/\&lt;\/\1\&gt;/g')

  # Strip leading blank lines, then remove the first H1 title since we add our own
  content=$(echo "$content" | sed '/./,$!d' | sed '1{/^# /d;}')

  seo=$(generate_seo "$title")
  tags=$(generate_tags "$slug")

  wrap_frontmatter "$title" "$slug" "$tags" "$seo" "$content" > "${dest_dir}/module.md"

  echo "Synced: ${src} -> ${dest_dir}/module.md"
done

echo "Done. Synced ${#DOCS[@]} docs to ${OUTPUT_DIR}"
