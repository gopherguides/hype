SELECT body,
    command,
    exit,
    go_version,
    src,
    tag,
    sum
FROM cmd_cache
WHERE command = ?
    and exit = ?
    and go_version = ?
    and src = ?
    and tag = ?
    and sum = ?