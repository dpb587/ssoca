# `ssoca ssh agent ...`

Start an SSH agent

    Usage:
      ssoca [OPTIONS] ssh agent [agent-OPTIONS]
    
    Application Options:
          --config=              Configuration file path (default: ~/.config/ssoca/config) [$SSOCA_CONFIG]
      -e, --environment=         Environment name [$SSOCA_ENVIRONMENT]
    
    Help Options:
      -h, --help                 Show this help message
    
    [agent command options]
          -s, --service=         Service name (default: ssh) [$SSOCA_SERVICE]
              --skip-auth-retry  Skip interactive authentication retries when logged out
              --foreground       Stay in foreground
              --socket=          Socket path (ensure the directory has restricted permissions)
    

## Usage Details

This agent follows the `ssh-agent` protocol to dynamically sign a certificate when public keys are requested. Callers can then use that signed certificate for authentication and signing.

You can use this like a regular `ssh-agent` and execute the output with `eval` (which will reconfigure `SSH_AUTH_SOCK`). For example...

    $ eval `ssoca ssh agent`
    ssoca ssh agent pid 12345

Once a key is added, each request to list keys will request a signed certificate from the ssoca server.

    $ ssh-add -l
    4096 SHA256:00j1lkGyGsQWesSK+p52DzZqZk20frTza5hwqr+vGyQ /Users/me/.ssh/id_rsa (RSA)
    4096 SHA256:00j1lkGyGsQWesSK+p52DzZqZk20frTza5hwqr+vGyQ /Users/me/.ssh/id_rsa (ssoca agent) (RSA-CERT)

If you [debug](./#debugging) the extra certificate, you'll see it is short-lived and changes every time.

    $ ssh-keygen -L -f <( ssh-add -L | grep 'ssoca agent' ) | grep Valid
            Valid: from 2017-03-28T22:47:40 to 2017-03-28T22:49:45

If you pass a command as an extra argument, it will be executed with `SSH_AUTH_SOCK` and the agent will terminate when the subprocess exits.

    $ ssoca ssh agent -- env | grep SSH_AUTH_SOCK && env | grep SSH_AUTH_SOCK
    SSH_AUTH_SOCK=/var/folders/hd/981tq6w92ll4qhy5s2wjffy40000gn/T/ssoca-ssh-agent170802401/agent.sock
    SSH_AUTH_SOCK=/private/tmp/com.apple.launchd.QhrzaIS6vh/Listeners


## Client Configuration

If the `SSH_AUTH_SOCK` environment variable is defined, key and signing operations will be delegated to that process. If `SSH_AUTH_SOCK` is not configured, key and certificate management will be handled within the process.

You may find the following SSH client options useful...

 * `IdentityAgent` - configure a specific path to the agent socket
 * `PasswordAuthentication` - disable password authentication

To use an agent for a specific VM, you could use the following configuration...

    Host jumpbox.example.com
      # mkdir ~/.ssh/agent && chmod 0700 ~/.ssh/agent
      IdentityAgent ~/.ssh/agent/%h

      # buggy; this tries to automatically start a one-off agent. it currently fails if the socket already exists
      ProxyCommand ssoca ssh agent --socket=~/.ssh/agent/%h -- nc %h %p


## Workflow

The general workflow between this agent, `ssh`, `ssh-agent`, the remote SSH server, and remote ssoca server looks something like this...

<div class="wsd" wsd_style="roundgreen"><pre>
  ssh->ssh-server: hello
  ssh-server->ssh: pubkey plz
  ssh->ssoca-ssh-agent: list identities
  ssoca-ssh-agent->ssh-agent: list identities
  ssh-agent->ssoca-ssh-agent: public keys
  ssoca-ssh-agent->ssoca: sign public key
  ssoca->ssoca-ssh-agent: short-lived certificate
  ssoca-ssh-agent->ssoca-ssh-agent: add certs to list
  ssoca-ssh-agent->ssh: pubkeys + certs
  ssh->ssh-server: pubkey
  ssh-server->ssh: more plz
  ssh->ssh-server: pubkey-cert
  ssh-server->ssh: prove it
  ssh->ssoca-ssh-agent: sign proof
  ssoca-ssh-agent->ssh-agent: sign proof
  ssh-agent->ssoca-ssh-agent: signature
  ssoca-ssh-agent->ssh: signature
  ssh->ssh-server: signature
  ssh-server->ssh: you win
  ssh->ssh-server: gimme a shell
  note over ssh,ssh-server: ssh session
  ssh-server<->ssh: disconnect
</pre></div>


## Notes

Keep in mind that when forwarding any SSH agent (`-A` or `-o ForwardAgent=yes`), remote systems can access and initiate keyring operations on your local workstation. In the case of this agent, it means requests will locally be sent to the ssoca server to dynamically generate a signed certificate.


<script type="text/javascript" src="https://www.websequencediagrams.com/service.js"></script>
