INSERT INTO cmd_cache (
        body,
        command,
        exit,
        go_version,
        src,
        tag,
        sum
    )
VALUES (
        :body,
        :command,
        :exit,
        :go_version,
        :src,
        :tag,
        :sum
    )