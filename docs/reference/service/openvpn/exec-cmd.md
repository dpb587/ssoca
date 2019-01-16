---
title: ssoca openvpn exec
aliases:
- /service/openvpn/exec-cmd/
---

# `ssoca openvpn exec ...`

Execute openvpn to connect to the remote server (aliases: connect)

    Usage:
      ssoca [OPTIONS] openvpn [openvpn-OPTIONS] exec [exec-OPTIONS] [EXTRA...]
    
    Application Options:
          --config=                 Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=            Environment name [$SSOCA_ENVIRONMENT]
          --log-level=              Log level (default: WARN) [$SSOCA_LOG_LEVEL]
    
    Help Options:
      -h, --help                    Show this help message
    
    [openvpn command options]
    
        Establish OpenVPN connections to remote servers:
          -s, --service=            Service name (default: openvpn) [$SSOCA_SERVICE]
    
    [exec command options]
              --skip-auth-retry     Skip interactive authentication retries when logged out
              --exec=               Path to the openvpn binary
              --reconnect           Reconnect on connection disconnects
              --static-certificate  Write a static certificate in the configuration instead of dynamic renewals
              --sudo                Execute openvpn with sudo
    
    [exec command arguments]
      EXTRA:                        Additional arguments to pass to openvpn
    
