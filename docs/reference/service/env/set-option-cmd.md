---
title: ssoca env set-option
aliases:
- /service/env/set-option-cmd/
---

# `ssoca env set-option ...`

Set a local client option in the environment

    Usage:
      ssoca [OPTIONS] env set-option [set-option-OPTIONS] [NAME] [VALUE]

    Application Options:
          --config=      Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment= Environment name [$SSOCA_ENVIRONMENT]
          --log-level=   Log level (default: WARN) [$SSOCA_LOG_LEVEL]

    Help Options:
      -h, --help         Show this help message

    [set-option command options]
          -s, --service= Service name (default: env) [$SSOCA_SERVICE]
              --ca-cert= Environment CA certificate path

    [set-option command arguments]
      NAME:              Client option name
      VALUE:             Client option value (parsed as YAML)
