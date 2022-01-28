SELECT body,
    command,
    exit,
    go_version,
    src,
    sum
FROM cmd_cache
WHERE command = ?
    and exit = ?
    and go_version = ?
    and src = ?
    and sum = ?