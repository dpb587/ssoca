# `ssoca env set ...`

Configure the connection to an environment (aliases: add)

    Usage:
      ssoca [OPTIONS] env set [set-OPTIONS] [URI]
    
    Application Options:
          --config=      Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment= Environment name [$SSOCA_ENVIRONMENT]
          --log-level=   Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help         Show this help message
    
    [set command options]
          -s, --service= Service name (default: env) [$SSOCA_SERVICE]
              --ca-cert= Environment CA certificate path
    
    [set command arguments]
      URI:               Environment URL
    
