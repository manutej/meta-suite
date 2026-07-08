# Meta-Suite — Workflow Diagrams

Hand-drawn (Excalidraw-aesthetic) diagrams documenting how the **Meta-Suite** plugin
works. Ink `#0C0B0A` on warm paper `#F4EFE7`; MAKE nodes accented gold `#E2B65A`,
KEEP nodes accented teal `#43C9BC`, feedback flow in blue `#58A6FF`.

## Diagrams

### 1. `workflow` — The meta-meta space
**Two axes that close a loop.**

- **MAKE axis** (gold): a left-to-right chain
  `meta-prompting → meta-planning → meta-design → meta-mvp`.
- **KEEP axis** (teal): `meta-review → meta-compound → meta-context`.
- **meta-craft** is drawn as a banner/standard spanning across both axes — the
  quality bar carried through every step.
- A curved **blue feedback arrow** loops from the KEEP axis back up to
  `meta-prompting`: *learnings + context flow back → warmer, sharper next round.*

### 2. `skill-anatomy` — Anatomy of one meta-skill
**The universal contract every meta-suite skill follows**, as a vertical flow:

1. **Silent intake** — auto-classify the request, never ask.
2. **Run typed reasoning INTERNALLY** — the category-theory machinery stays hidden.
3. **Emit a PLAIN-language deliverable** — no jargon, just the useful result.
4. **Offer to reveal internal work** — only if the user asks.

A side box (`references/internal-model.md`, *hidden unless asked*) holds the internal
model that step 2 draws on.

## Files

| File | Role |
|------|------|
| `workflow.excalidraw` | **Editable source** — open at [excalidraw.com](https://excalidraw.com) (File → Open) |
| `workflow.svg` | Rendered export, embedded in the main README |
| `skill-anatomy.excalidraw` | **Editable source** — open at [excalidraw.com](https://excalidraw.com) |
| `skill-anatomy.svg` | Rendered export, embedded in the main README |

The `.excalidraw` files are the editable sources — open them at
[excalidraw.com](https://excalidraw.com) to revise. The `.svg` files are the rendered
exports that get embedded in the main project README. When you change a diagram, edit
the `.excalidraw` source and re-export the matching `.svg`.
