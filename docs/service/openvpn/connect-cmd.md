# `ssoca openvpn connect ...`

Connect to a remote OpenVPN server

    Usage:
      ssoca [OPTIONS] openvpn connect [connect-OPTIONS] [EXTRA...]
    
    Application Options:
          --config=                 Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=            Environment name [$SSOCA_ENVIRONMENT]
          --log-level=              Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                    Show this help message
    
    [connect command options]
          -s, --service=            Service name (default: openvpn) [$SSOCA_SERVICE]
              --skip-auth-retry     Skip interactive authentication retries when logged out
              --exec=               Path to the openvpn binary
              --reconnect           Reconnect on connection disconnects
              --static-certificate  Write a static certificate in the configuration instead of dynamic renewals
              --sudo                Execute openvpn with sudo
    
    [connect command arguments]
      EXTRA:                        Additional arguments to pass to openvpn
    
