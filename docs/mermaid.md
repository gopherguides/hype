# Mermaid Diagrams

Hype supports [Mermaid](https://mermaid.js.org/) diagrams, which are automatically rendered as ASCII art. This allows you to include diagrams directly in your markdown without external image files, and keeps diagrams version-controlled as text.

## Usage

Use standard fenced code blocks with the `mermaid` language identifier:

    ```mermaid
    graph LR
        A[Start] --> B{Decision}
        B -->|Yes| C[Action]
        B -->|No| D[End]
    ```

When processed by Hype, this will be rendered as ASCII art:

```
┌──────────┐     ┌─────────────┐
│          │     │             │
│ A[Start] ├────►│ B{Decision} │
│          │     │             │
└──────────┘     └─────────────┘
```

## Supported Diagram Types

### Flowcharts/Graphs

Both `graph` and `flowchart` directives are supported with these directions:
- `LR` - Left to Right
- `TD` / `TB` - Top Down / Top to Bottom

    ```mermaid
    graph TD
        Client --> API
        API --> Auth
        API --> Database
        Auth --> Database
    ```

### Sequence Diagrams

Sequence diagrams show interactions between participants:

    ```mermaid
    sequenceDiagram
        Alice->>Bob: Hello Bob
        Bob-->>Alice: Hi Alice
    ```

## Limitations

The ASCII rendering is provided by [mermaid-ascii](https://github.com/AlexanderGrooff/mermaid-ascii), which has some limitations:

**Supported:**
- Graph flowcharts (LR, TD/TB directions)
- Sequence diagrams
- Labeled edges
- Color definitions via `classDef` (rendered as text styling in supported terminals)

**Not Supported:**
- Subgraph nesting
- Non-rectangular node shapes (diamonds render as rectangles)
- Class diagrams
- State diagrams
- Gantt charts
- Pie charts
- Diagonal arrows

## Output Format

In HTML export, mermaid diagrams are rendered as `<pre><code>` blocks with the ASCII art content.

In Markdown export, they appear as plain code blocks (without language specifier) containing the ASCII art.
