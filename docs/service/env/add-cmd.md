# `ssoca env add ...`

Add a new environment

    Usage:
      ssoca [OPTIONS] env add [add-OPTIONS] [URI]
    
    Application Options:
          --config=      Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment= Environment name [$SSOCA_ENVIRONMENT]
          --log-level=   Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help         Show this help message
    
    [add command options]
          -s, --service= Service name (default: env) [$SSOCA_SERVICE]
              --ca-cert= Environment CA certificate path
    
    [add command arguments]
      URI:               Environment URL
    
