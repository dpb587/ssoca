# `ssoca openvpn connect ...`

Connect to a remote OpenVPN server

    Usage:
      ssoca [OPTIONS] openvpn connect [connect-OPTIONS] [EXTRA...]
    
    Application Options:
          --config=        Configuration file path (default: ~/.ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=   Environment name [$SSOCA_ENVIRONMENT]
    
    Help Options:
      -h, --help           Show this help message
    
    [connect command options]
          -s, --service=   Service name (default: openvpn) [$SSOCA_SERVICE]
              --exec=      Path to the openvpn binary
              --reconnect  Reconnect on connection disconnects
              --sudo       Execute openvpn with sudo
    
    [connect command arguments]
      EXTRA:               Additional arguments to pass to openvpn
    
